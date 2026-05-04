package handlers

import (
	"bellesquare-be/internal/models"
	repositories "bellesquare-be/internal/repository"
	"bellesquare-be/internal/utils"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (uh *UserHandler) AddUserHandler(c *echo.Context) error {
	slog.Info("Received request to add a new user")

	// 1. Initialize the model and bind JSON/Form data
	slog.Info("Creating a user object in memory from the request")
	req := new(models.RegisterUserRequest)
	if err := c.Bind(req); err != nil {
		slog.Error("Failed to bind user request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// If it's a guest, we store a placeholder or empty string
	slog.Info("Checking if the user is a guest")
	var hashedPassword string
	var err error
	if req.IsGuest {
		slog.Info("User is a guest, skipping password hashing")
		hashedPassword = "GUEST_ACCOUNT"
	} else {
		// Validate required fields
		slog.Debug("Validating the fields")
		if req.Username == "" || req.Password == "" || req.Email == "" {
			slog.Warn("Missing required fields for user creation")
			return c.JSON(http.StatusBadRequest, map[string]string{
				"warning": "Missing required fields: username, password, fullname, email",
			})
		}
		// Argon2 Hashing (using the utility we discussed)
		slog.Debug("Hashing the password using")
		hashedPassword, err = utils.HashPassword(req.Password)
		if err != nil {
			slog.Error("Security failure during hashing", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal security error"})
		}
	}

	// Map DTO to the clean Model
	slog.Debug("Convert request to a user object")
	user := &models.User{
		UserName:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		VenmoHandle:    req.VenmoHandle,
		ZelleHandle:    req.ZelleHandle,
		PreferredColor: req.PreferredColor,
		IsGuest:        req.IsGuest,
	}

	// Create the user using the repository method
	err = uh.userRepo.CreateUser(c.Request().Context(), user)
	if err != nil {
		slog.Error("Failed to create user in database", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	// Return the fully hydrated user (including ID and CreatedAt, json:"-" tags will hide the password hash)
	slog.Info("Successfully created user", "id", user.ID, "email", user.Email)
	return c.JSON(http.StatusCreated, user)
}

func (uh *UserHandler) UpdatePasswordHandler(c *echo.Context) error {
	slog.Info("Received request to update user password")

	slog.Debug("Parsing the update password request")
	req := new(models.UpdatePasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// 1. Fetch the existing user to get the current hash
	slog.Debug("Fetching user by ID", "id", req.UserID)
	ctx := c.Request().Context()
	user, err := uh.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// 2. Verify the old password (we'll need a Verify utility for Argon2)
	slog.Debug("Verifying the old password")
	match, err := utils.VerifyPassword(req.OldPassword, user.PasswordHash)
	if err != nil || !match {
		slog.Warn("Old password verification failed", "user_id", req.UserID)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Current password incorrect"})
	}

	// 3. Hash the NEW password
	slog.Debug("Hashing the new password")
	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		slog.Error("Failed to hash new password", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Hashing failed"})
	}

	// 4. Update via Repo
	slog.Debug("Updating the password in the database")
	if err := uh.userRepo.UpdatePassword(ctx, user.ID, newHash); err != nil {
		slog.Error("Failed to update password in database", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
	}

	slog.Info("Password updated successfully", "user_id", user.ID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Password updated successfully"})
}

func (uh *UserHandler) UpdateProfile(c *echo.Context) error {
	slog.Info("Updating users profile")

	// 1. Identification (Extract ID from JWT set by middleware)
	slog.Debug("Extracting user ID from JWT token")
	// userToken := c.Get("user").(*jwt.Token)
	// claims := userToken.Claims.(jwt.MapClaims)
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Invalid user ID in token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// 2. Data Extraction
	// Use a specific struct for updates to prevent users from
	// trying to overwrite their own 'id' or 'created_at' fields.
	slog.Info("Binding update profile request data")
	req := new(models.UpdateRequest)
	if err := c.Bind(req); err != nil {
		slog.Error("Failed to bind update profile request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// 3. Map DTO to model (fetch existing user first to preserve unchanged fields)
	slog.Debug("Fetching existing user data for profile update", "user_id", userID)
	user, err := uh.userRepo.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		slog.Error("Failed to fetch existing user data", "user_id", userID, "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	user.UserName = req.Username // client side can decide to allow username changes or not
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.VenmoHandle = req.VenmoHandle
	user.ZelleHandle = req.ZelleHandle
	user.PreferredColor = req.PreferredColor

	// 4. Database Operation
	slog.Debug("Updating user profile in database", "user_id", userID)
	err = uh.userRepo.UpdateProfile(c.Request().Context(), user)
	if err != nil {
		slog.Error("Failed to update user profile in database", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	slog.Info("User profile updated successfully", "user_id", userID)
	return c.NoContent(http.StatusOK)
}

func (uh *UserHandler) UpgradeAccount(c *echo.Context) error {
	slog.Info("---------- UPGRADE GUEST ACCOUNT TO MEMBER ----------")

	// Identification (Extract ID from JWT set by middleware)
	slog.Debug("Extracting user ID from JWT token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Invalid user ID in token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	req := new(models.UpgradeRequest)
	req.UserID = userID

	// Bind the email and password from the request body
	slog.Info("Binding update profile request data")
	if err := c.Bind(req); err != nil {
		slog.Error("Failed to bind to the upgrade request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	slog.Debug("Validating password")
	if req.Password == "" {
		slog.Warn("Upgrade failed due to empty password", "user_id", userID)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password cannot be empty"})
	}

	slog.Debug("Hashing password")
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid password"})
	}

	req.PasswordHash = hashed

	// Call the specific, safe method
	err = uh.userRepo.UpgradeGuestToMember(c.Request().Context(), req)
	if err != nil {
		slog.Error("Failed to upgrade guest account", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Upgrade failed"})
	}

	slog.Info("Guest account upgraded successfully", "email", req.Email)
	return c.NoContent(http.StatusOK)
}
