package repositories

import (
	"bellesquare-be/internal/models"
	"context"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateProfile(ctx context.Context, user *models.User) error
	UpgradeGuestToMember(ctx context.Context, req *models.UpgradeRequest) error
}

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepository {
	return &UserRepo{
		pool: pool,
	}
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	slog.Debug("Adding a new user to the db", "username", user.UserName, "email", user.Email)

	newID, err := uuid.NewV7()
	if err != nil {
		slog.Error("Error generating uuid")
		return fmt.Errorf("failed to generate uuid: %w", err)
	}
	user.ID = newID

	slog.Info("ID generated for user", "email: ", user.Email)

	qry := `
		INSERT INTO users (id, username, password_hash, first_name, last_name, email, venmo_handle, zelle_handle, preferred_color, is_guest)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at
	`

	err = ur.pool.QueryRow(ctx, qry,
		user.ID, user.UserName, user.PasswordHash, user.FirstName, user.LastName,
		user.Email, user.VenmoHandle, user.ZelleHandle, user.PreferredColor, user.IsGuest,
	).Scan(&user.CreatedAt)

	if err != nil {
		return fmt.Errorf("create user query failed: %w", err)
	}

	slog.Info("Successfully added user.", "email", user.Email, "guest?", user.IsGuest)

	return nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	slog.Debug("Fetching user by ID", "id", id)

	user := new(models.User)

	slog.Debug("Executing query to fetch user by ID")
	qry := `
		SELECT id, username, first_name, last_name, password_hash, email, is_guest, venmo_handle, zelle_handle, preferred_color, created_at
		FROM users
		WHERE id = $1
	`

	err := ur.pool.QueryRow(ctx, qry, id).Scan(
		&user.ID,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Email,
		&user.IsGuest,
		&user.VenmoHandle,
		&user.ZelleHandle,
		&user.PreferredColor,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by id failed: %w", err)
	}

	slog.Info("Successfully fetched user by ID", "email", user.Email)
	return user, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	slog.Info("Fetching user by email", "email", email)

	var u models.User

	slog.Debug("Executing query to fetch user by email")
	query := `
		SELECT 
			id, username, first_name, last_name, email, 
			password_hash, venmo_handle, zelle_handle, 
			preferred_color, is_guest 
		FROM users WHERE email = $1`
	err := ur.pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.UserName,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PasswordHash,
		&u.VenmoHandle,
		&u.ZelleHandle,
		&u.PreferredColor,
		&u.IsGuest,
	)
	if err != nil {
		slog.Error("Error fetching user by email", "email", email, "error", err)
		return nil, fmt.Errorf("get user by email failed: %w", err)
	}

	slog.Info("Successfully retrieved user by email")
	return &u, nil
}

func (ur *UserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error {
	slog.Debug("Updating password for user", "user_id", id)
	qry := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`

	slog.Debug("Executing password update query")
	result, err := ur.pool.Exec(ctx, qry, newHash, id)
	if err != nil {
		return fmt.Errorf("update password query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}

	slog.Info("Password updated successfully for user", "user_id", id)
	return nil
}

func (ur *UserRepo) UpdateProfile(ctx context.Context, user *models.User) error {
	slog.Info("Updating user profile", "user_id", user.ID)
	qry := `
		UPDATE users SET 
			first_name = $1, 
			last_name = $2, 
			venmo_handle = $3, 
			zelle_handle = $4,
			preferred_color = $5,
			is_guest = $6,
			updated_at = NOW() 
		WHERE id = $7`

	slog.Debug("Executing profile update query")
	result, err := ur.pool.Exec(ctx, qry, user.FirstName, user.LastName, user.VenmoHandle, user.ZelleHandle, user.PreferredColor, user.IsGuest, user.ID)
	if err != nil {
		slog.Error("Failed to execute profile update query", "user_id", user.ID, "error", err)
		return fmt.Errorf("update profile query failed: %w", err)
	}

	slog.Debug("Profile update query executed", "user_id", user.ID, "rows_affected", result.RowsAffected())
	if result.RowsAffected() == 0 {
		slog.Warn("No user found to update", "user_id", user.ID)
		return fmt.Errorf("no user found with id %s", user.ID)
	}

	slog.Info("User profile updated successfully", "user_id", user.ID)
	return nil
}

// NOTE: This is safer than re-using UpdateProfile
func (ur *UserRepo) UpgradeGuestToMember(ctx context.Context, req *models.UpgradeRequest) error {
	slog.Info("Upgrading guest to member", "user_id", req.UserID)

	slog.Debug("Executing guest upgrade query")
	qry := `
        UPDATE users 
        SET 
            password_hash = $1, 
            is_guest = false, 
            updated_at = NOW() 
        WHERE id = $2 AND is_guest = true` // CRITICAL: Only works if they are currently a guest

	result, err := ur.pool.Exec(ctx, qry, req.PasswordHash, req.UserID)
	if err != nil {
		slog.Error("Failed to execute guest upgrade query", "user_id", req.UserID, "error", err)
		return fmt.Errorf("upgrade query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		// This tells us the user was already a member or doesn't exist
		slog.Warn("Account upgrade ineligible for user", "email", req.Email)
		return fmt.Errorf("account upgrade ineligible for user %s", req.Email)
	}

	slog.Info("Successfully upgraded guest to member", "user_id", req.UserID, "email", req.Email)
	return nil
}
