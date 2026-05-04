package handlers

import (
	"bellesquare-be/internal/models"
	repositories "bellesquare-be/internal/repository"
	"bellesquare-be/internal/utils"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	userRepo repositories.UserRepository
	secret   []byte
}

func NewAuthHandler(u repositories.UserRepository, secret string) *AuthHandler {
	return &AuthHandler{
		userRepo: u,
		secret:   []byte(secret),
	}
}

func (h *AuthHandler) Login(c *echo.Context) error {
	slog.Info("Login attempt")
	// 1. Initialize the user struct
	req := new(models.RegisterUserRequest)

	// 2. Bind the form data directly to the struct
	// This maps "user_name" to u.UserName based on the json tags
	slog.Debug("Binding the request")
	if err := c.Bind(req); err != nil {
		slog.Error("Failed to bind form data", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
	}

	// we return "Invalid credentials" for all errors to avoid giving hints to attackers about which part of the login failed (email vs password)
	// This is a common security practice to prevent user enumeration attacks.
	slog.Info("Validating email")
	if req.Email == "" {
		slog.Error("Missing email")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid credentials"})
	}

	// Use Repo to check if user is in DB
	slog.Info("Checking user email in DB", "email", req.Email)
	user, err := h.userRepo.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		slog.Error("User not found or DB error", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	if !user.IsGuest {
		slog.Info("Verifying password for non-guest user")
		if req.Password == "" {
			slog.Error("Missing password")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid credentials"})
		}

		match, err := utils.VerifyPassword(req.Password, user.PasswordHash)
		if err != nil || !match {
			slog.Warn("Login failed: password mismatch")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
	} else {
		slog.Info("Guest login bypass", "email", user.Email)
	}

	// Generate JWT (Handler logic)
	slog.Debug("Generating JWT token for user", "user_id", user.Email)
	token, err := utils.GenerateToken(user, h.secret)
	if err != nil {
		slog.Error("Failed to generate JWT token", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	slog.Debug("Token generated", "token", token)

	slog.Info("Login successful", "email", user.Email)
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) Logout(c *echo.Context) error {
	slog.Info("Logout attempt")
	cookie := &http.Cookie{
		Name:     "at",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Deletes the cookie immediately
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	slog.Info("Logout successful")
	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) Register(c *echo.Context) error {
	slog.Info("Registration attempt")
	// 1. Initialize the user struct
	req := new(models.RegisterUserRequest)
	var hash string
	var err error

	// 2. Bind the form data directly to the struct
	// This maps "user_name" to u.UserName based on the json tags
	if err = c.Bind(req); err != nil {
		slog.Error("Failed to bind form data", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
	}

	// 3. Checking if this is a guest as a guest does not need a password
	slog.Info("Is this a guest user?")
	if !req.IsGuest {
		slog.Info("Not a Guest - Checking for password in form data")
		if req.Password == "" {
			slog.Warn("Registration failed: password is missing")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password is required"})
		}
		slog.Debug("Password provided, hashing password")
		hash, err = utils.HashPassword(req.Password)
		if err != nil {
			slog.Error("Failed to hash password", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to secure password"})
		}
	} else {
		// Guest users get a placeholder hash or empty string
		slog.Info("This is a guest user")
		hash = "GUEST_ACCOUNT_NO_PASSWORD"
	}

	// 4. Map DTO to the User Model and Save to Database via the Repository
	slog.Info("Creating the user account")
	u := &models.User{
		Email:          req.Email,
		UserName:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		IsGuest:        req.IsGuest,
		VenmoHandle:    req.VenmoHandle,
		ZelleHandle:    req.ZelleHandle,
		PreferredColor: req.PreferredColor,
		PasswordHash:   hash,
	}

	slog.Debug("Saving new user to database", "email", u.Email)
	if err := h.userRepo.CreateUser(c.Request().Context(), u); err != nil {
		slog.Error("Failed to create user", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not save user: " + err.Error()})
	}

	// 5. Return the created user (minus the password hash)
	slog.Info("User created successfully", "email", u.Email)
	return c.JSON(http.StatusCreated, u)
}

func (h *AuthHandler) Me(c *echo.Context) error {
	slog.Info("Fetching current user info")

	// 1. Get the userID from the token
	slog.Debug("Grabbing the users ID from the token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get the user id from the token", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid User ID format"})
	}

	// 2. Fetch fresh user data from DB
	slog.Debug("Fetching user data from DB", "user_id", userID)
	user, err := h.userRepo.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		slog.Error("User not found or DB error", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	slog.Info("User data fetched successfully", "user_id", user.ID)
	return c.JSON(http.StatusOK, user)
}
