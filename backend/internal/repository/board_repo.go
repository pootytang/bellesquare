package repositories

import (
	"bellesquare-be/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardRepository interface {
	CreateBoard(ctx context.Context, board *models.Board, payoutRepo PayoutRepository, amounts map[models.PayoutPeriod]float64) error
	GetBoardByBoardID(ctx context.Context, id uuid.UUID) (*models.BoardDetailsResponse, error)
	GetBoardsByUserID(ctx context.Context, userID uuid.UUID) ([]models.BoardSummary, error)
	GetBoardWithSquares(ctx context.Context, id uuid.UUID) (*models.BoardWithSquares, error)
	ClaimSquaresBatch(ctx context.Context, boardID uuid.UUID, userID uuid.UUID, color string, selections []models.SquareCoords) error
	UpdateAxisNumbers(ctx context.Context, boardID uuid.UUID, home, away []int) error
	UpdateBoardStatus(ctx context.Context, boardID uuid.UUID, status models.BoardStatus) error
	ToggleUserPaidStatus(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	GetSquaresByBoardAndUser(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) ([]models.SquareUpdate, error)
	UpdateQuarterScore(ctx context.Context, boardID uuid.UUID, locked bool, quarter, home, away int) error
	GetJoinedBoards(ctx context.Context, userID uuid.UUID) ([]models.BoardSummary, error)
	UpdateUserPaymentStatus(ctx context.Context, boardID, userID uuid.UUID, status models.PaymentStatus) error
}

type BoardRepo struct {
	pool *pgxpool.Pool
}

func NewBoardRepo(pool *pgxpool.Pool) BoardRepository {
	slog.Info("Initializing Board Repository")
	return &BoardRepo{pool: pool}
}

func (r *BoardRepo) CreateBoard(ctx context.Context, board *models.Board, payoutRepo PayoutRepository, amounts map[models.PayoutPeriod]float64) error {
	slog.Info("Starting atomic transaction for board creation", "title", board.Title)

	slog.Debug("Starting transaction to create board and squares")

	// 1. Start the Transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	// Defer a rollback. If Commit() is called, the rollback does nothing.
	defer tx.Rollback(ctx)

	// 2. Generate Board ID and Insert Parent Row
	slog.Info("Querying the DB for board creation")
	boardID, _ := uuid.NewV7()
	board.ID = boardID

	// 3. Initialize the Batch
	batch := &pgx.Batch{}

	// Queue the Board Insert
	slog.Debug("Queing the batch Board insert")
	batch.Queue(`
		INSERT INTO boards (id, title, sport, status, price_per_square, creator_id, home_team_id, away_team_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`,
		board.ID, board.Title, board.Sport, board.Status,
		board.PricePerSquare, board.CreatorID, board.HomeTeamID, board.AwayTeamID,
	)

	slog.Debug("Queing up the squares inserts for the board", "board id", board.ID)
	// Queue 100 Square Inserts
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			sID, _ := uuid.NewV7()
			batch.Queue(`
				INSERT INTO squares (id, board_id, row_index, col_index)
				VALUES ($1, $2, $3, $4)`,
				sID, board.ID, r, c,
			)
		}
	}

	// 4. Send the Batch
	slog.Info("Sending the batch to the database")
	br := tx.SendBatch(ctx, batch)

	// 5. Check the Board Insert Result (the first queued item)
	slog.Debug("Checking the board insert result")
	if err := br.QueryRow().Scan(&board.CreatedAt, &board.UpdatedAt); err != nil {
		br.Close()
		return fmt.Errorf("batch board insert failed: %w", err)
	}

	// Close batch results (clears the 100 square results)
	slog.Debug("Closing the batch results to clear the square insert results")
	if err := br.Close(); err != nil {
		return err
	}

	// 6. Check the 100 Square Insert Results
	// slog.Debug("Checking the square insert results")
	// for i := 0; i < 100; i++ {
	// 	_, err := br.Exec()
	// 	if err != nil {
	// 		br.Close()
	// 		return fmt.Errorf("batch square insert failed at index %d: %w", i, err)
	// 	}
	// }

	// 7. Call Payout Repo inside the same Transaction
	if err := payoutRepo.InitializeBoardPayouts(ctx, tx, board.ID, amounts); err != nil {
		slog.Error("Failed to initialize payouts for board", "board_id", board.ID, "error", err)
		return fmt.Errorf("payout initialization failed: %w", err)
	}

	// 8. Commit the Transaction
	slog.Info("Committing the transaction for board creation")
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Successfully created board and all squares", "board_id", board.ID)
	return nil
}

func (r *BoardRepo) GetBoardByBoardID(ctx context.Context, id uuid.UUID) (*models.BoardDetailsResponse, error) {
	slog.Info("---------- GET BOARD BY BOARD ID CALLED ----------")
	slog.Debug("Fetching board with scores", "board_id", id)
	b := &models.Board{}
	var home, away models.TeamSummary
	var scores models.BoardScores
	var squaresJSON []byte  // Temporary holder for the aggregated squares
	var quartersJSON []byte // Temporary holder for scores by quarter

	slog.Info("Querying the DB for board retrieval", "board_id", id)
	qry := `
        WITH current_scores AS (
            SELECT COALESCE(json_agg(qs ORDER BY quarter), '[]'::json) as quarters_list
            FROM (
                SELECT quarter, home_score, away_score, is_locked
                FROM board_scores 
                WHERE board_id = $1
            ) qs
        ),
        aggregated_squares AS (
            SELECT COALESCE(json_agg(sq), '[]'::json) as squares_list
            FROM (
                SELECT bs.row_index, bs.col_index, bs.user_id, bs.is_paid, bs.payment_status,
                    u.first_name, u.last_name, 
                    u.preferred_color as user_color
                FROM squares bs
                LEFT JOIN users u ON bs.user_id = u.id
                WHERE bs.board_id = $1
            ) sq
        )
        SELECT 
            b.id, b.title, b.sport, b.status, b.price_per_square, 
            b.home_axis_numbers, b.away_axis_numbers, 
            b.creator_id, b.home_team_id, b.away_team_id, 
            b.created_at, b.updated_at,
            ht.id, ht.city, ht.mascot, ht.team_logo_url, ht.primary_color, ht.secondary_color,
            at.id, at.city, at.mascot, at.team_logo_url, at.primary_color, at.secondary_color,
            cs.quarters_list, -- Scan individual quarters instead of sums
            asq.squares_list,
			u.venmo_handle,
        	u.zelle_handle
        FROM boards b
        JOIN teams ht ON b.home_team_id = ht.id
        JOIN teams at ON b.away_team_id = at.id
		JOIN users u ON b.creator_id = u.id
        CROSS JOIN current_scores cs
        CROSS JOIN aggregated_squares asq
        WHERE b.id = $1`

	var creatorVenmo, creatorZelle string
	err := r.pool.QueryRow(ctx, qry, id).Scan(
		&b.ID, &b.Title, &b.Sport, &b.Status, &b.PricePerSquare,
		&b.HomeAxisNumbers, &b.AwayAxisNumbers,
		&b.CreatorID, &b.HomeTeamID, &b.AwayTeamID,
		&b.CreatedAt, &b.UpdatedAt,
		&home.ID, &home.City, &home.Mascot, &home.TeamLogoURL, &home.PrimaryColor, &home.SecondaryColor,
		&away.ID, &away.City, &away.Mascot, &away.TeamLogoURL, &away.PrimaryColor, &away.SecondaryColor,
		&quartersJSON,
		&squaresJSON,
		&creatorVenmo,
		&creatorZelle,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Query returned zero rows")
			return nil, err // Handlers should check for pgx.ErrNoRows to return 404
		}

		slog.Error("Query failed", "error", err)
		return nil, fmt.Errorf("failed to get hydrated board: %w", err)
	}

	// Initialize the 10x10 grid
	slog.Info("Initializing the grid")
	owners := make(map[string]models.SquareOwnerDetails)
	var grid [10][10]models.Square
	for rIdx := 0; rIdx < 10; rIdx++ {
		for cIdx := 0; cIdx < 10; cIdx++ {
			grid[rIdx][cIdx] = models.Square{RowIndex: rIdx, ColIndex: cIdx}
		}
	}

	// Unmarshal the JSON stuff
	slog.Info("Getting the list of Quarters completed")
	var quarterList []models.UpdateScoreRequest // Assuming you have this struct in your models
	if err := json.Unmarshal(quartersJSON, &quarterList); err != nil {
		slog.Error("Failed to unmarshal quarters", "error", err)
		return nil, err
	}

	if len(quarterList) == 0 {
		for i := 1; i <= 4; i++ {
			quarterList = append(quarterList, models.UpdateScoreRequest{Quarter: i})
		}
	}

	slog.Info("Grabbing the home and away total score")
	var totalHome, totalAway int
	for _, q := range quarterList {
		// RUNNING TOTALS LOGIC:
		// The highest score entered in any quarter is the current total.
		if q.HomeScore > totalHome {
			totalHome = q.HomeScore
		}
		if q.AwayScore > totalAway {
			totalAway = q.AwayScore
		}
	}

	// Populate the scores struct for the response
	scores = models.BoardScores{
		Home:     totalHome,
		Away:     totalAway,
		Quarters: quarterList,
	}

	slog.Debug("Unmarshalling the json to squares slice")
	var flatSquares []struct {
		RowIndex      int                  `json:"row_index"`
		ColIndex      int                  `json:"col_index"`
		UserID        *uuid.UUID           `json:"user_id"`
		IsPaid        bool                 `json:"is_paid"`
		FirstName     string               `json:"first_name"`
		LastName      string               `json:"last_name"`
		UserColor     string               `json:"user_color"`
		PaymentStatus models.PaymentStatus `json:"payment_status"`
	}

	if err := json.Unmarshal(squaresJSON, &flatSquares); err != nil {
		slog.Error("Unable to unmarshal json", "error", err)
		return nil, fmt.Errorf("failed to unmarshal squares: %w", err)
	}

	// 3. Populate the grid with owner names/details
	slog.Debug("Populating grid with details")
	slog.Debug("----- DATA FROM JSON -----")
	for _, s := range flatSquares {
		if s.UserID != nil {
			key := fmt.Sprintf("%d-%d", s.RowIndex, s.ColIndex)
			owners[key] = models.SquareOwnerDetails{
				FirstName: s.FirstName,
				LastName:  s.LastName,
				Color:     s.UserColor,
			}

			// Calculate the initials
			var initials string
			if s.FirstName != "" && s.LastName != "" {
				initials = strings.ToUpper(string(s.FirstName[0]) + string(s.LastName[0]))
			}
			// Populate the Pure Square for the grid
			grid[s.RowIndex][s.ColIndex] = models.Square{
				RowIndex:      s.RowIndex,
				ColIndex:      s.ColIndex,
				UserID:        s.UserID,
				IsPaid:        s.IsPaid,
				PaymentStatus: s.PaymentStatus,
				UserColor:     s.UserColor,
				UserInitials:  initials,
			}

			if s.PaymentStatus == models.Pending {
				slog.Debug("PENDING PAYMENT SQUARE", "Row", s.RowIndex, "Col", s.ColIndex, "UserID", s.UserID, "FirstName", s.FirstName, "LastName", s.LastName)
			}
			// slog.Debug("Added to the Owners map", "key", key, "First Name", owners[key].FirstName, "Last Name", owners[key].LastName, "User Color", owners[key].Color)
		}
	}

	home.FullName = home.City + " " + home.Mascot
	away.FullName = away.City + " " + away.Mascot

	slog.Info("Returning the board details")
	return &models.BoardDetailsResponse{
		Board:        models.BoardWithTeams{Board: b, HomeTeam: home, AwayTeam: away},
		Scores:       scores,
		Squares:      grid,
		SquareOwners: owners,
		CreatorVenmo: creatorVenmo,
		CreatorZelle: creatorZelle,
	}, nil
}

// Retrieves all 100 squares which is not good over the websocket
// TODO: Maybe remove this in favor of GetBoardsByBoardID which was doing the same but has more data.
func (r *BoardRepo) GetBoardWithSquares(ctx context.Context, id uuid.UUID) (*models.BoardWithSquares, error) {
	slog.Info("---------- GET BOARD WITH SQUARES CALLED ----------")

	// Get the Board Metadata
	slog.Debug("Fetching board with squares by ID", "board_id", id)
	boardDetails, err := r.GetBoardByBoardID(ctx, id)
	if err != nil {
		slog.Error("Unable to get board by board id", "error", err)
		return nil, err
	}

	// Get the 100 Squares
	slog.Info("Querying the DB for squares of the board", "board_id", id)
	qry := `
        SELECT 
            s.id, s.board_id, s.row_index, s.col_index, s.user_id, s.user_color, s.is_paid, s.payment_status, s.claimed_at,
            u.first_name, u.last_name
        FROM squares s
        LEFT JOIN users u ON s.user_id = u.id
        WHERE s.board_id = $1 
        ORDER BY s.row_index ASC, s.col_index ASC`

	rows, err := r.pool.Query(ctx, qry, id)
	if err != nil {
		slog.Error("Problem querying the squares table", "error", err)
		return nil, fmt.Errorf("failed to query squares: %w", err)
	}
	defer rows.Close()

	slog.Debug("Retrieving the result")
	res := &models.BoardWithSquares{
		Board:   boardDetails.Board,
		Squares: [10][10]models.Square{},
	}

	slog.Debug("Processing the retrieved squares and placing them into the 2D grid")
	for rows.Next() {
		var s models.Square
		var fName, lName *string // Use pointers because they might be NULL (unclaimed)

		err := rows.Scan(
			&s.ID, &s.BoardID, &s.RowIndex, &s.ColIndex, &s.UserID, &s.UserColor, &s.IsPaid, &s.PaymentStatus, &s.ClaimedAt,
			&fName, &lName,
		)
		if err != nil {
			return nil, err
		}

		// Calculate Initials
		if fName != nil && lName != nil && *fName != "" {
			s.UserInitials = strings.ToUpper(string((*fName)[0]) + string((*lName)[0]))
		}

		if s.RowIndex < 10 && s.ColIndex < 10 {
			res.Squares[s.RowIndex][s.ColIndex] = s
		}

		if s.PaymentStatus == models.Pending {
			slog.Debug("PENDING PAYMENT FOUND")
		}
	}

	slog.Info("Successfully retrieved board and squares", "board_id", id)
	return res, nil
}

func (r *BoardRepo) GetBoardsByUserID(ctx context.Context, userID uuid.UUID) ([]models.BoardSummary, error) {
	slog.Info("Fetching board list for user", "user_id", userID)

	// Join the teams table twice to get details for both Home and Away teams
	qry := `
        SELECT 
            b.id, b.title, b.status, b.price_per_square,
            ht.id, ht.city, ht.mascot, ht.team_logo_url, ht.primary_color, ht.secondary_color,
            at.id, at.city, at.mascot, at.team_logo_url, at.primary_color, at.secondary_color
        FROM boards b
        INNER JOIN teams ht ON b.home_team_id = ht.id
        INNER JOIN teams at ON b.away_team_id = at.id
        WHERE b.creator_id = $1 
        ORDER BY b.created_at DESC`

	slog.Debug("Performing the join")
	rows, err := r.pool.Query(ctx, qry, userID)
	if err != nil {
		slog.Error("Failed to query boards for user", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to query the user's boards: %w", err)
	}
	defer rows.Close()

	slog.Debug("Scanning the rows into the summaries")
	summaries := []models.BoardSummary{} // Initialize with an empty slice so it JSON encodes to [] instead of null
	for rows.Next() {
		var s models.BoardSummary
		// Scan order must match the SELECT order exactly
		err := rows.Scan(
			&s.ID, &s.Title, &s.Status, &s.PricePerSquare,
			&s.HomeTeam.ID, &s.HomeTeam.City, &s.HomeTeam.Mascot, &s.HomeTeam.TeamLogoURL, &s.HomeTeam.PrimaryColor, &s.HomeTeam.SecondaryColor,
			&s.AwayTeam.ID, &s.AwayTeam.City, &s.AwayTeam.Mascot, &s.AwayTeam.TeamLogoURL, &s.AwayTeam.PrimaryColor, &s.AwayTeam.SecondaryColor,
		)
		if err != nil {
			slog.Error("Failed to scan board summary row", "error", err)
			return nil, err
		}
		summaries = append(summaries, s)
	}

	slog.Info("Successfully retrieved board summaries", "user_id", userID, "count", len(summaries))
	return summaries, nil
}

func (r *BoardRepo) ClaimSquaresBatch(ctx context.Context, boardID uuid.UUID, userID uuid.UUID, color string, selections []models.SquareCoords) error {
	slog.Info("ClaimBatch called")
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		slog.Error("error configuring a transcation object", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	slog.Debug("Querying the db")
	for _, sel := range selections {
		// We use 'user_id IS NULL' to ensure we only claim unclaimed squares
		res, err := tx.Exec(ctx, `
            UPDATE squares 
            SET user_id = $1, user_color = $2, claimed_at = NOW()
            WHERE board_id = $3 AND row_index = $4 AND col_index = $5 AND user_id IS NULL`,
			userID, color, boardID, sel.Row, sel.Col)

		if err != nil {
			slog.Error("Error querying the db", "error", err)
			return fmt.Errorf("db error: %w", err)
		}

		if res.RowsAffected() == 0 {
			slog.Warn("Square is possibly claimed already", "row", sel.Row, "col", sel.Col)
			return fmt.Errorf("square at %d,%d was already claimed", sel.Row, sel.Col)
		}
	}

	slog.Info("Committing the transaction")
	if err := tx.Commit(ctx); err != nil {
		slog.Error("Failed to Commit the transaction", "error", err)
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (r *BoardRepo) ToggleUserPaidStatus(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error {
	slog.Info("Toggling if a user has paid or not", "board id", boardID, "user id", userID)
	// We toggle the boolean using NOT.
	// This updates every square the user owns on this specific board.
	slog.Debug("Querying the db")
	qry := `
		UPDATE squares 
		SET is_paid = NOT is_paid 
		WHERE board_id = $1 AND user_id = $2`

	_, err := r.pool.Exec(ctx, qry, boardID, userID)

	slog.Info("Query executed")
	return err
}

func (r *BoardRepo) UpdateAxisNumbers(ctx context.Context, boardID uuid.UUID, home, away []int) error {
	slog.Info("Update Axis repo called")

	slog.Debug("running the query to update the axis numbers")
	query := `
        UPDATE boards 
        SET home_axis_numbers = $1, 
            away_axis_numbers = $2,
            updated_at = NOW()
        WHERE id = $3
    `
	_, err := r.pool.Exec(ctx, query, home, away, boardID)
	return err
}

// Status can be open, locked, completed
func (r *BoardRepo) UpdateBoardStatus(ctx context.Context, boardID uuid.UUID, status models.BoardStatus) error {
	slog.Info("Updating the board status", "status", status)

	slog.Debug("Querying the db")
	query := `
		UPDATE boards 
		SET status = $1, 
		    updated_at = NOW() 
		WHERE id = $2
	`

	result, err := r.pool.Exec(ctx, query, status, boardID)
	if err != nil {
		slog.Error("there was an error updating the status", "error", err)
		return fmt.Errorf("failed to update board status: %w", err)
	}

	// Optional: Check if the board actually existed
	slog.Debug("checking if the board was found")
	if result.RowsAffected() == 0 {
		slog.Warn("No board was found")
		return fmt.Errorf("no board found with id %s", boardID)
	}

	slog.Info("Successfully found and updated the Board")
	return nil
}

func (r *BoardRepo) GetSquaresByBoardAndUser(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) ([]models.SquareUpdate, error) {
	slog.Info("Retrieving squares for a user")
	// We join users to get the names for the initials calculation

	slog.Debug("Querying the db", "board id", boardID, "user id", userID)
	qry := `
    SELECT 
        s.row_index, s.col_index, s.user_id, s.user_color, s.is_paid, s.payment_status,
        u.first_name, u.last_name
    FROM squares s
    LEFT JOIN users u ON s.user_id = u.id
    WHERE s.board_id = $1 AND s.user_id = $2`

	rows, err := r.pool.Query(ctx, qry, boardID, userID)
	if err != nil {
		slog.Error("Unable to execute the query", "error", err)
		return nil, err
	}
	defer rows.Close()

	slog.Debug("Grabbing the data from the rows")
	var updates []models.SquareUpdate
	for rows.Next() {
		var su models.SquareUpdate
		var fName, lName *string

		err := rows.Scan(
			&su.Row, &su.Col, &su.UserID, &su.UserColor, &su.IsPaid, &su.PaymentStatus, &fName, &lName,
		)
		if err != nil {
			slog.Error("Unable to scan data into square object", "error", err)
			return nil, err
		}

		// Apply the initials logic so the frontend stays consistent
		if fName != nil && lName != nil && *fName != "" && *lName != "" {
			su.FirstName = *fName
			su.LastName = *lName
			su.UserInitials = strings.ToUpper(string((*fName)[0]) + string((*lName)[0]))
			// Optionally set fName/lName on the square if you added them to the struct
		}

		slog.Debug("adding square")
		updates = append(updates, su)
	}

	slog.Info("Successfully created square slice")
	return updates, nil
}

// TODO: Refactor to take an UpdateScore object instead of passing so many things
func (r *BoardRepo) UpdateQuarterScore(ctx context.Context, boardID uuid.UUID, locked bool, quarter, home, away int) error {
	slog.Info("Updating the score in the db")

	slog.Debug("querying the db", "board id", boardID, "quarter", quarter, "home score", home, "away score", away)
	qry := `
        INSERT INTO board_scores (board_id, quarter, home_score, away_score, is_locked, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
        ON CONFLICT (board_id, quarter) 
        DO UPDATE SET 
            home_score = EXCLUDED.home_score, 
            away_score = EXCLUDED.away_score, 
            is_locked = EXCLUDED.is_locked,
            updated_at = NOW()`

	_, err := r.pool.Exec(ctx, qry, boardID, quarter, home, away, locked)

	slog.Info("Query completed")
	return err
}

func (r *BoardRepo) GetJoinedBoards(ctx context.Context, userID uuid.UUID) ([]models.BoardSummary, error) {
	slog.Info("Fetching boards joined by user", "user_id", userID)

	// use DISTINCT because a user might own multiple squares on one board
	// join 'squares' to find where the user is a participant
	slog.Debug("Running the query")
	qry := `
        SELECT DISTINCT
            b.id, b.title, b.status, b.price_per_square,
            ht.id, ht.city, ht.mascot, ht.team_logo_url, ht.primary_color, ht.secondary_color,
            at.id, at.city, at.mascot, at.team_logo_url, at.primary_color, at.secondary_color,
            b.created_at
        FROM boards b
        INNER JOIN squares s ON b.id = s.board_id
        INNER JOIN teams ht ON b.home_team_id = ht.id
        INNER JOIN teams at ON b.away_team_id = at.id
        WHERE s.user_id = $1 AND b.creator_id != $1
        ORDER BY b.created_at DESC`

	rows, err := r.pool.Query(ctx, qry, userID)
	if err != nil {
		slog.Error("Failed to query joined boards", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to query joined boards: %w", err)
	}
	defer rows.Close()

	slog.Info("Scanning the rows into the summaries")
	summaries := []models.BoardSummary{}
	for rows.Next() {
		var s models.BoardSummary
		var createdAt time.Time // <--- Temp variable to catch the scan

		err := rows.Scan(
			&s.ID, &s.Title, &s.Status, &s.PricePerSquare,
			&s.HomeTeam.ID, &s.HomeTeam.City, &s.HomeTeam.Mascot, &s.HomeTeam.TeamLogoURL, &s.HomeTeam.PrimaryColor, &s.HomeTeam.SecondaryColor,
			&s.AwayTeam.ID, &s.AwayTeam.City, &s.AwayTeam.Mascot, &s.AwayTeam.TeamLogoURL, &s.AwayTeam.PrimaryColor, &s.AwayTeam.SecondaryColor,
			&createdAt,
		)
		if err != nil {
			slog.Error("Failed to scan joined board summary row", "error", err)
			return nil, err
		}
		summaries = append(summaries, s)
	}

	slog.Info("Successfully retrieved joined board summaries", "user_id", userID, "count", len(summaries))
	return summaries, nil
}

func (r *BoardRepo) UpdateUserPaymentStatus(ctx context.Context, boardID, userID uuid.UUID, status models.PaymentStatus) error {
	slog.Info("---------- UPDATE USER PAYMENT STATUS CALLED ----------")

	slog.Debug("Updating the payment status for a user's squares on a board", "board id", boardID, "user id", userID, "new status", status)
	query := `
    UPDATE squares 
    SET payment_status = $1, 
        is_paid = ($1 = 2) -- Automatically sets to true if status is 2, else false
    WHERE board_id = $2 AND user_id = $3`
	_, err := r.pool.Exec(ctx, query, int(status), boardID, userID)

	if err != nil {
		slog.Error("Failed to update user payment status", "error", err)
		return fmt.Errorf("failed to update user payment status: %w", err)
	}

	slog.Info("Successfully updated user payment status", "board id", boardID, "user id", userID, "new status", status)
	return nil
}
