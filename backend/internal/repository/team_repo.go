package repositories

import (
	"bellesquare-be/internal/models"
	"bellesquare-be/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn" // Required for Postgres error codes
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTeamInUse = errors.New("cannot delete team: it is currently assigned to one or more boards")

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamsBySport(ctx context.Context, sport models.SportEnum, creator uuid.UUID) ([]models.Team, error)
	GetTeamByID(ctx context.Context, id, creator uuid.UUID) (*models.Team, error)
	GetTeamByFullIdentity(ctx context.Context, creator, city, mascot string, sport models.SportEnum) (*models.Team, error)
	GetTeamCounts(ctx context.Context, userID uuid.UUID) (models.TeamCounts, error)
	// UpdateTeamLogo(ctx context.Context, id uuid.UUID, logoURL string) error
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id, creatorid string) error
}

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) TeamRepository {
	slog.Debug("Creating a Team repo")
	return &TeamRepo{pool: pool}
}

func (r *TeamRepo) CreateTeam(ctx context.Context, team *models.Team) error {
	slog.Debug("Creating team in the db", "team", team)

	slog.Debug("Normalizing the data")
	team.City = utils.NormalizeString(team.City)
	team.Mascot = utils.NormalizeString(team.Mascot)

	newID, err := uuid.NewV7()
	if err != nil {
		slog.Error("Problem generating a UUID", "error", err)
		return fmt.Errorf("failed to generate uuid: %w", err)
	}
	team.ID = newID

	qry := `
		INSERT INTO teams (id, creator_id, city, mascot, sport, team_logo_url, primary_color, secondary_color)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at`

	// pgx handles team.Sport (Enum) and team.TeamLogoURL (*string) automatically
	slog.Debug("Inserting team into the db", "team", team)
	err = r.pool.QueryRow(ctx, qry,
		team.ID,
		team.CreatorID,
		team.City,
		team.Mascot,
		team.Sport,
		team.TeamLogoURL,
		team.PrimaryColor,
		team.SecondaryColor,
	).Scan(&team.CreatedAt)

	if err != nil {
		// 1. Check if the error is a Postgres error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 2. 23505 is the code for unique_violation
			if pgErr.Code == "23505" {
				slog.Warn("Unique Violation in teams table. Duplicate entry found",
					"city", team.City,
					"mascot", team.Mascot,
				)
				return fmt.Errorf("a team with this city and mascot already exists")
			}
		}
		slog.Error("Problem inserting team", "error", err)
		return fmt.Errorf("failed to insert team: %w", err)
	}

	slog.Debug("Successfully inserted team into the db", "team", team)
	return nil
}

func (r *TeamRepo) GetTeamsBySport(ctx context.Context, sport models.SportEnum, creator uuid.UUID) ([]models.Team, error) {
	slog.Debug("Retrieving teams by sport", "sport", sport)

	qry := `
		SELECT id, creator_id, city, mascot, sport, team_logo_url, primary_color, secondary_color, created_at 
		FROM teams 
		WHERE sport = $1 AND creator_id = $2
		ORDER BY mascot ASC`

	rows, err := r.pool.Query(ctx, qry, sport, creator)
	if err != nil {
		slog.Error("Querying teams failed", "error", err)
		return nil, fmt.Errorf("query teams failed: %w", err)
	}
	defer rows.Close()

	slog.Info("retrieved teams by sport")
	slog.Debug("Creating a teams list")

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		// SportEnum.Scan() will be called automatically here
		err := rows.Scan(&t.ID, &t.CreatorID, &t.City, &t.Mascot, &t.Sport, &t.TeamLogoURL, &t.PrimaryColor, &t.SecondaryColor, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan team failed: %w", err)
		}
		teams = append(teams, t)
	}

	slog.Info("Successfully retrieved the teams list")
	return teams, nil
}

func (r *TeamRepo) GetTeamByID(ctx context.Context, id, creator uuid.UUID) (*models.Team, error) {
	slog.Debug("Retrieving team by ID", "id", id, "creatorID", creator)
	t := &models.Team{}

	slog.Debug("Querying the DB", "id", id)
	qry := `
		SELECT 
			id, creator_id, city, mascot, sport, team_logo_url, primary_color, secondary_color, COALESCE(tertiary_color, ''), created_at 
		FROM 
			teams WHERE id = $1 AND creator_id = $2`

	err := r.pool.QueryRow(ctx, qry, id, creator).Scan(
		&t.ID,
		&t.CreatorID,
		&t.City,
		&t.Mascot,
		&t.Sport,
		&t.TeamLogoURL,
		&t.PrimaryColor,
		&t.SecondaryColor,
		&t.TertiaryColor,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Team not found or unauthorized access", "id", id, "creator", creator)
			return nil, nil // Or a specific ErrNotFound to handle in the handler
		}
		slog.Error("Database failure during GetTeamByID", "error", err)
		return nil, fmt.Errorf("get team by id failed: %w", err)
	}

	slog.Info("Successfully retrieved team by ID", "id", id)
	return t, nil
}

func (r *TeamRepo) GetTeamByFullIdentity(ctx context.Context, creator, city, mascot string, sport models.SportEnum) (*models.Team, error) {
	slog.Debug("Retrieving team by full identity", "city", city, "mascot", mascot, "sport", sport, "creatorID", creator)

	t := &models.Team{}

	// Using ILIKE for case-insensitive matching of city and mascot
	slog.Debug("Querying the DB for team by full identity")
	qry := `
		SELECT id, creator_id, city, mascot, sport, team_logo_url, created_at 
		FROM teams 
		WHERE city ILIKE $1 AND mascot ILIKE $2 AND sport = $3 AND creator_id = $4`

	err := r.pool.QueryRow(ctx, qry, city, mascot, sport, creator).Scan(
		&t.ID, &t.CreatorID, &t.City, &t.Mascot, &t.Sport, &t.TeamLogoURL, &t.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Warn("A team was not found", "error", err)
			return nil, fmt.Errorf("team not found: %s %s", city, mascot)
		}
		return nil, err
	}

	slog.Info("Successfully retrieved team by full identity", "city", city, "mascot", mascot, "sport", sport)
	return t, nil
}

// A null logoURL will set team_logo_url to NULL in the database, effectively removing the logo.
// func (r *TeamRepo) UpdateTeamLogo(ctx context.Context, id uuid.UUID, logoURL string) error {
// 	slog.Debug("Updating team logo", "team_id", id, "logo_url", logoURL)
// 	var logoPtr *string

// 	// If a string is provided, point to it.
// 	// If it's empty, logoPtr stays nil (which translates to NULL in SQL)
// 	if logoURL != "" {
// 		logoPtr = &logoURL
// 	}

// 	slog.Debug("Executing update query for team logo", "team_id", id, "logo_url", logoURL)
// 	qry := `UPDATE teams SET team_logo_url = $1 WHERE id = $2`

// 	result, err := r.pool.Exec(ctx, qry, logoPtr, id)
// 	if err != nil {
// 		slog.Error("Failed to execute update query for team logo", "team_id", id, "error", err)
// 		return fmt.Errorf("failed to update team logo: %w", err)
// 	}

// 	slog.Debug("Update query executed", "team_id", id, "rows_affected", result.RowsAffected())

// 	if result.RowsAffected() == 0 {
// 		slog.Warn("No team found to update the logo for", "team_id", id)
// 		return fmt.Errorf("no team found with id: %s", id)
// 	}

// 	slog.Info("Updated team logo", "team_id", id, "logo_url", logoURL)
// 	return nil
// }

func (r *TeamRepo) UpdateTeam(ctx context.Context, team *models.Team) error {
	slog.Info("Updating team in the db", "team", team)
	slog.Debug("Executing update query for team", "team_id", team.ID)
	query := `
        UPDATE teams 
        SET city = $1, mascot = $2, sport = $3, team_logo_url = $4, primary_color = $5, secondary_color = $6
        WHERE id = $7 AND creator_id = $8`

	result, err := r.pool.Exec(ctx, query,
		team.City,
		team.Mascot,
		team.Sport.String(), // Using the String() method for DB storage
		team.TeamLogoURL,
		team.PrimaryColor,
		team.SecondaryColor,
		team.ID,
		team.CreatorID,
	)
	if err != nil {
		slog.Warn("Problem updating team in DB")
		return fmt.Errorf("failed to update team: %w", err)
	}

	if result.RowsAffected() == 0 {
		slog.Warn("No team found to update", "team_id", team.ID, "creatorID", team.CreatorID)
		return fmt.Errorf("no team found with id: %s", team.ID)
	}

	slog.Info("Successfully updated team", "team_id", team.ID)
	return nil
}

func (r *TeamRepo) DeleteTeam(ctx context.Context, id, creatorid string) error {
	slog.Info("Deleting team from db", "id", id, "creatorID", creatorid)

	slog.Debug("Querying the DB")
	query := `DELETE FROM teams WHERE id = $1 AND creator_id = $2`

	result, err := r.pool.Exec(ctx, query, id, creatorid)
	if err != nil {
		// Check if the error is a foreign key violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			slog.Error("DB Error found")
			if pgErr.Code == "23503" { // ForeignKeyViolation
				slog.Warn("ForeignKeyViolation. Team may be in use")
				return ErrTeamInUse
			}
		}
		slog.Error("Problem executing the query to delete a team", "error", err)
		return err
	}

	// Check if a row was actually deleted
	if result.RowsAffected() == 0 {
		slog.Warn("no rows found", "id", id)
		return fmt.Errorf("no team found with id %s", id)
	}

	slog.Info("Successfully deleted team from DB")
	return nil
}

func (r *TeamRepo) GetTeamCounts(ctx context.Context, userID uuid.UUID) (models.TeamCounts, error) {
	slog.Info("Grabbing Team Counts by sport from the DB")
	counts := models.TeamCounts{}

	// We can get all counts in a single query using conditional aggregation
	slog.Debug("Querying the db")
	qry := `
        SELECT 
            COUNT(*) FILTER (WHERE sport = 'football') as football,
            COUNT(*) FILTER (WHERE sport = 'basketball') as basketball,
            COUNT(*) FILTER (WHERE sport = 'baseball') as baseball,
            COUNT(*) FILTER (WHERE sport = 'soccer') as soccer,
            COUNT(*) FILTER (WHERE sport = 'hockey') as hockey
        FROM teams 
        WHERE creator_id = $1`

	err := r.pool.QueryRow(ctx, qry, userID).Scan(
		&counts.Football,
		&counts.Basketball,
		&counts.Baseball,
		&counts.Soccer,
		&counts.Hockey,
	)

	if err != nil {
		slog.Error("Unable to get team counts from the db", "error", err)
		return models.TeamCounts{}, err
	}

	slog.Info("Successfully retrieved the counts")
	return counts, err
}
