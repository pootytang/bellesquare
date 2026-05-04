package utils

import (
	"bellesquare-be/internal/models"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type UserProfile struct {
	ID       uuid.UUID
	Initials string
	IsGuest  bool
}

func HashPassword(password string) (string, error) {
	slog.Info("Hashing password")
	// Argon2id Recommended Parameters (can be adjusted based on server hardware)
	time := uint32(1)
	memory := uint32(64 * 1024) // 64 MB
	threads := uint8(4)
	keyLen := uint32(32)
	saltLen := uint32(16)

	slog.Debug("Generating the salt")
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Return a string that includes the salt and parameters so we can verify later
	// Format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	slog.Debug("Encoding the hash and salt")
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, time, threads, b64Salt, b64Hash)

	slog.Debug("Final encoded hash", "hash", encodedHash)
	slog.Info("Password hashing complete")
	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	slog.Info("Verifying password")
	// 1. Split the string into parts
	// Format: $argon2id$v=19$m=65536,t=1,p=4$SALT$HASH
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, ErrInvalidHash
	}

	slog.Debug("Parsed hash parts", "parts", parts)
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, ErrIncompatibleVersion
	}

	slog.Debug("Extracted version from hash", "version", version)
	var memory, iterations uint32
	var threads uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &threads)
	if err != nil {
		return false, err
	}

	// 2. Decode the Salt and the Original Hash from Base64
	slog.Debug("Extracted parameters from hash", "memory", memory, "iterations", iterations, "threads", threads)
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	slog.Debug("Decoded salt from hash", "salt", salt)
	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	keyLen := uint32(len(decodedHash))

	// 3. Hash the provided password using the parameters extracted from the stored hash
	slog.Debug("Hashing the provided password with extracted parameters")
	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, threads, keyLen)

	// 4. Use subtle.ConstantTimeCompare to prevent timing attacks
	slog.Debug("Comparing hashes", "decodedHash", decodedHash, "comparisonHash", comparisonHash)
	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1 {
		slog.Info("Password verification successful")
		return true, nil
	}

	slog.Debug("Hashes do not match")
	slog.Info("Password verification failed", "password", password, "encodedHash", encodedHash)
	return false, nil
}

// generateToken creates a signed JWT for the user
func GenerateToken(user *models.User, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub":        user.ID.String(),
		"given_name": user.FirstName,
		"sur_name":   user.LastName,
		"is_guest":   user.IsGuest,
		"exp":        time.Now().Add(time.Hour * 72).Unix(), // Valid for 3 days
		"iat":        time.Now().Unix(),                     // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GetUserIDFromJWT(c *echo.Context) (uuid.UUID, error) {
	slog.Info("Retrieving ID from token")

	slog.Debug("Retrieving the claims from the context")
	claims, err := getClaims(c)
	if err != nil {
		slog.Error("Problem retrieving claims from context")
		return uuid.Nil, fmt.Errorf("Problem retrieving claims")
	}

	slog.Debug("Grabbing the subject from the token")
	subject, _ := claims["sub"].(string)
	creatorID, err := uuid.FromString(subject)
	if err != nil {
		slog.Debug("Unable to parse the subject from the claims")
		return uuid.Nil, fmt.Errorf("Unable to get the creatorID from the subject claim")
	}

	slog.Info("Successfully retrieved the creatorID", "id", creatorID)
	return creatorID, nil
}

func GetUserProfileFromJWT(c *echo.Context) (UserProfile, error) {
	slog.Info("Grabbing user profile from JWT")

	slog.Debug("retrieving the claims")
	claims, err := getClaims(c)
	if err != nil {
		slog.Error("Unable to get claims from context")
		return UserProfile{}, fmt.Errorf("Unable to retrieve claims")
	}

	// 1. Get ID
	sub, _ := claims["sub"].(string)
	id, err := uuid.FromString(sub)
	if err != nil {
		return UserProfile{}, fmt.Errorf("invalid sub claim")
	}

	// 2. Get Names & Generate Initials
	// Note: match the keys used in your GenerateToken function
	fName, _ := claims["given_name"].(string)
	lName, _ := claims["sur_name"].(string)
	is_guest, _ := claims["is_guest"].(bool)

	slog.Debug(("Creating initials"))
	var initials string
	if len(fName) > 0 {
		initials += string(fName[0])
	}
	if len(lName) > 0 {
		initials += string(lName[0])
	}
	initials = strings.ToUpper(initials)

	if initials == "" {
		initials = "??"
	}

	slog.Info("Returning the users simple profile", "ID", id, "Initials", initials, "Is Guest", is_guest)
	return UserProfile{
		ID:       id,
		Initials: initials,
		IsGuest:  is_guest,
	}, nil
}

func getClaims(c *echo.Context) (jwt.MapClaims, error) {
	slog.Info("Grabbing the claims from the context")

	slog.Debug("retrieving the token")
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		slog.Error("Unable to retrieve token")
		return jwt.MapClaims{}, fmt.Errorf("token not found in context")
	}

	slog.Debug("retrieving the claims out of the token")
	claims, ok := token.Claims.(jwt.MapClaims)

	return claims, nil
}
