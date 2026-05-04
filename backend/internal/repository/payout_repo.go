package repositories

import (
	"bellesquare-be/internal/models"
	"context"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PayoutRepository interface {
	// InitializeBoardPayouts(ctx context.Context, boardID uuid.UUID, amounts map[models.PayoutPeriod]float64) error
	InitializeBoardPayouts(ctx context.Context, tx pgx.Tx, boardID uuid.UUID, amounts map[models.PayoutPeriod]float64) error
	AwardWinner(ctx context.Context, boardID uuid.UUID, period models.PayoutPeriod, winnerID uuid.UUID) error
	GetPayoutsByBoard(ctx context.Context, boardID uuid.UUID) ([]models.Payout, error)
	UpdatePayoutAmounts(ctx context.Context, boardID uuid.UUID, amounts map[models.PayoutPeriod]float64) error
	UpsertPayouts(ctx context.Context, boardID uuid.UUID, payouts []models.Payout) error
}

type PayoutRepo struct {
	pool *pgxpool.Pool
}

func NewPayoutRepo(pool *pgxpool.Pool) PayoutRepository {
	slog.Info("Initializing Payout Repository")
	return &PayoutRepo{pool: pool}
}

// InitializeBoardPayouts sets up the 5 prize rows when a board is created
// InitializeBoardPayouts now takes a tx to be part of the atomic board creation
func (pr *PayoutRepo) InitializeBoardPayouts(ctx context.Context, tx pgx.Tx, boardID uuid.UUID, amounts map[models.PayoutPeriod]float64) error {
	slog.Info("Initializing board payouts", "board_id", boardID)

	slog.Info("Payout amounts for each period", "amounts", amounts)
	slog.Debug("Querying the DB to initialize payouts for board", "board_id", boardID)
	qry := `INSERT INTO payouts (id, board_id, period_name, amount) VALUES ($1, $2, $3, $4)`

	for period, amt := range amounts {
		id, _ := uuid.NewV7()
		if _, err := tx.Exec(ctx, qry, id, boardID, period, amt); err != nil {
			slog.Error("Failed to initialize payout for period", "period", period.String(), "error", err)
			return fmt.Errorf("failed to init payout for %s: %w", period.String(), err)
		}
	}

	slog.Info("Successfully initialized payouts for board", "board_id", boardID)
	return nil
}

// AwardWinner marks the person who won that specific period
func (pr *PayoutRepo) AwardWinner(ctx context.Context, boardID uuid.UUID, period models.PayoutPeriod, winnerID uuid.UUID) error {
	slog.Info("Awarding payout winner", "board_id", boardID, "period", period.String(), "winner_id", winnerID)

	// Only update if the associated board is LOCKED (Status 1)
	slog.Debug("Querying the DB if board is locked")
	qry := `
        UPDATE payouts 
        SET winner_user_id = $1, awarded_at = NOW() 
        WHERE board_id = $2 AND period_name = $3
        AND EXISTS (
            SELECT 1 FROM boards WHERE id = $2 AND status = 1
        )`

	res, err := pr.pool.Exec(ctx, qry, winnerID, boardID, period)
	if err != nil {
		slog.Error("Failed to award payout winner", "board_id", boardID, "period", period.String(), "winner_id", winnerID, "error", err)
		return err
	}

	slog.Debug("Checking if any rows were affected to confirm winner was awarded")
	if res.RowsAffected() == 0 {
		slog.Warn("No payout awarded - either board is not locked or payout not found", "board_id", boardID, "period", period.String())
		return fmt.Errorf("could not award winner: board is not locked or payout not found")
	}

	slog.Info("Successfully awarded payout winner", "board_id", boardID, "period", period.String(), "winner_id", winnerID)
	return nil
}

func (pr *PayoutRepo) GetPayoutsByBoard(ctx context.Context, boardID uuid.UUID) ([]models.Payout, error) {
	slog.Info("Fetching payouts for board", "board_id", boardID)

	slog.Debug("Querying the DB for payouts", "board_id", boardID)
	qry := `
        SELECT 
            p.id, p.board_id, p.period_name, p.amount, p.winner_user_id, p.awarded_at,
            CONCAT(u.first_name, ' ', u.last_name) as winner_name
        FROM payouts p
        LEFT JOIN users u ON p.winner_user_id = u.id
        WHERE p.board_id = $1 
        ORDER BY ARRAY_POSITION(ARRAY['q1', 'q2', 'q3', 'q4', 'final'], p.period_name)`

	rows, err := pr.pool.Query(ctx, qry, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payouts []models.Payout
	slog.Debug("Creating the payouts list for board", "board_id", boardID)
	for rows.Next() {
		var p models.Payout
		err := rows.Scan(&p.ID, &p.BoardID, &p.PeriodName, &p.Amount, &p.WinnerUserID, &p.AwardedAt, &p.WinnerName)
		if err != nil {
			slog.Error("Problem scanning into Payouts object", "error", err)
			return nil, err
		}
		payouts = append(payouts, p)
	}

	slog.Info("Successfully fetched payouts for board", "board_id", boardID, "num_payouts", len(payouts))
	return payouts, nil
}

func (pr *PayoutRepo) UpdatePayoutAmounts(ctx context.Context, boardID uuid.UUID, amounts map[models.PayoutPeriod]float64) error {
	slog.Info("Updating payout amounts", "board_id", boardID)

	slog.Debug("Creating the transaction")
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		slog.Error("Problem creating the transaction object", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	slog.Info("Updating db with payouts")
	qry := `UPDATE payouts SET amount = $1 WHERE board_id = $2 AND period_name = $3`

	for period, amt := range amounts {
		if _, err := tx.Exec(ctx, qry, amt, boardID, period); err != nil {
			slog.Error("Failed to update payout", "period", period.String(), "error", err)
			return err
		}
	}

	slog.Info("Committing the transaction")
	return tx.Commit(ctx)
}

func (pr *PayoutRepo) UpsertPayouts(ctx context.Context, boardID uuid.UUID, payouts []models.Payout) error {
	slog.Info("Upserting the payouts")
	// We use a transaction to ensure all payouts save or none do
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		slog.Error("Problem creating a transaction object", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	slog.Debug("Inserting the payouts")
	for _, p := range payouts {
		query := `
            INSERT INTO payouts (board_id, period_name, amount)
            VALUES ($1, $2, $3)
            ON CONFLICT (board_id, period_name) 
            DO UPDATE SET amount = EXCLUDED.amount;
        `
		_, err := tx.Exec(ctx, query, boardID, p.PeriodName, p.Amount)
		if err != nil {
			slog.Error("Problem inserting payouts", "error", err)
			return err
		}
	}

	slog.Info("Committing the transaction")
	return tx.Commit(ctx)
}
