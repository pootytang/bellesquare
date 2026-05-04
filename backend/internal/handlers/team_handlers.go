package handlers

import (
	"bellesquare-be/internal/models"
	repositories "bellesquare-be/internal/repository"
	"bellesquare-be/internal/utils"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v5"
)

type TeamHandler struct {
	teamRepo repositories.TeamRepository
}

func NewTeamHandler(repo repositories.TeamRepository) *TeamHandler {
	slog.Debug("Creating a Team handler")
	return &TeamHandler{teamRepo: repo}
}

func (th *TeamHandler) AddTeamHandler(c *echo.Context) error {
	slog.Info("Received request to add a new team")

	slog.Debug("Grabbing the users ID from the token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get the user id from the token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	req := new(models.CreateTeamRequest)
	if err := c.Bind(req); err != nil {
		slog.Warn("Problem with team data", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team data"})
	}
	req.CreatorID = userID

	// Basic validation
	slog.Debug("Validating the team data")
	if !req.Sport.IsValid() {
		slog.Warn("Invalid sport provided for team creation", "sport", req.Sport)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Sport is required"})
	}

	if req.City == "" || req.Mascot == "" {
		slog.Warn("Missing required fields for team creation")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "City and Mascot are required"})
	}

	slog.Debug("Creating a team object in memory from the request")
	team := &models.Team{
		City:           req.City,
		CreatorID:      req.CreatorID,
		Mascot:         req.Mascot,
		Sport:          req.Sport,
		TeamLogoURL:    req.TeamLogoURL,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
	}

	slog.Debug("Creating team in DB", "team", team)
	if err := th.teamRepo.CreateTeam(c.Request().Context(), team); err != nil {
		slog.Error("Failed to create team in the database", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save team"})
	}

	// Return the created team as a DTO
	slog.Info("Successfully created team", "team_id", team.ID)
	return c.JSON(http.StatusCreated, team.MapToDTO())
}

// GetTeamsHandler retrieves all teams by sport from the database
func (th *TeamHandler) GetTeamsHandler(c *echo.Context) error {
	slog.Info("Retrieving all teams based on sport")
	sportStr := c.QueryParam("sport")

	slog.Debug("Received request for teams", "sport", sportStr)

	slog.Debug("Grabbing the users ID from the token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get the user id from the token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// 1. Ensure the parameter exists
	if sportStr == "" {
		slog.Warn("Missing 'sport' query parameter")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Query parameter 'sport' is required (e.g., ?sport=football)",
		})
	}

	// 2. Validate the sport
	sport := models.FromString(sportStr)
	if !sport.IsValid() {
		slog.Warn("Invalid sport provided", "sport", sportStr)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("'%s' is not a supported sport", sportStr),
		})
	}

	slog.Debug("Fetching teams from the database for sport", "sport", sport)
	teams, err := th.teamRepo.GetTeamsBySport(c.Request().Context(), sport, userID)
	if err != nil {
		slog.Error("Failed to fetch teams", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// 3. Map to DTOs
	slog.Debug("Mapping team models to DTOs for response", "team count", len(teams))
	teamDTOs := make([]models.TeamDTO, 0, len(teams))
	for _, t := range teams {
		teamDTOs = append(teamDTOs, t.MapToDTO())
	}

	// 4. Return the list of teams
	slog.Info("Successfully retrieved teams", "team count", len(teamDTOs))
	return c.JSON(http.StatusOK, teamDTOs)
}

func (th *TeamHandler) UpdateTeamHandler(c *echo.Context) error {
	slog.Info("Received request to update a team")

	slog.Debug("Grabbing the team id from the params")
	idParam := c.Param("id")
	teamID, err := uuid.FromString(idParam)
	if err != nil {
		slog.Warn("Invalid team ID", "id", idParam, "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	slog.Debug("Grabbing the users ID from the token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get the user id from the token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	req := new(models.CreateTeamRequest)
	if err := c.Bind(req); err != nil {
		slog.Warn("Invalid request payload", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Strict Validation for Sport
	if !req.Sport.IsValid() {
		slog.Warn("Invalid sport provided", "sport", req.Sport)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Valid sport is required"})
	}

	ctx := c.Request().Context()

	// 1. Check if team exists
	existingTeam, err := th.teamRepo.GetTeamByID(ctx, teamID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Team not found"})
	}

	// 2. Update fields
	existingTeam.City = req.City
	existingTeam.Mascot = req.Mascot
	existingTeam.Sport = req.Sport
	existingTeam.TeamLogoURL = req.TeamLogoURL
	existingTeam.PrimaryColor = req.PrimaryColor
	existingTeam.SecondaryColor = req.SecondaryColor

	// 3. Persist
	if err := th.teamRepo.UpdateTeam(ctx, existingTeam); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update team"})
	}

	// 4. Return the updated Team as a DTO
	return c.JSON(http.StatusOK, existingTeam.MapToDTO())
}

func (th *TeamHandler) DeleteTeamHandler(c *echo.Context) error {
	slog.Info("Received a request to delete a team")
	idStr := c.Param("id")

	// 1. Parse the ID (Ensure it's a valid UUID)
	// If you're using google/uuid:
	// id, err := uuid.Parse(idStr)
	// if err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid team id"}) }

	//1. Get the users id from the token
	slog.Debug("Grabbing the users ID from the token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get the user id from the token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// 2. Call the repository to delete
	slog.Debug("Querying the DB to Delete team", "id", idStr, "creator id", userID)
	err = th.teamRepo.DeleteTeam(c.Request().Context(), idStr, userID.String())
	if err != nil {
		slog.Error("Problem deleting team", "error", err)
		if errors.Is(err, repositories.ErrTeamInUse) {
			// Return a 409 Conflict or 400 Bad Request with the specific message
			slog.Error("DB Violation error")
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "This team is currently assigned to a board and cannot be deleted.",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "An internal error occurred while trying to delete the team.",
		})
	}

	slog.Info("Successfully deleted team. Returning 204")
	return c.NoContent(http.StatusNoContent) // 204 No Content is standard for successful DELETE
}

func (th *TeamHandler) GetTeam(c *echo.Context) error {
	slog.Info("GetTeam called")

	slog.Debug("Grabbing the users ID from the token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get the user id from the token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	slog.Debug("Converting the team ID retrieved from the params to a uuid")
	teamID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		slog.Warn("Unable to get team ID")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid team ID format",
		})
	}

	slog.Info("Retrieved team from DB", "id", teamID)
	team, err := th.teamRepo.GetTeamByID(c.Request().Context(), teamID, userID)
	if err != nil {
		slog.Error("Unable to retrieve team from DB", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Team not found",
		})
	}

	slog.Debug("Mapping to team DTO")
	dto := models.TeamDTO{
		ID:             team.ID.String(),
		FullName:       team.FullName(),
		City:           team.City,
		Mascot:         team.Mascot,
		Sport:          team.Sport.String(),
		TeamLogoURL:    team.TeamLogoURL,
		PrimaryColor:   team.PrimaryColor,
		SecondaryColor: team.SecondaryColor,
	}

	slog.Info("Successfully retrieved team")
	return c.JSON(http.StatusOK, dto)
}

func (th *TeamHandler) GetCounts(c *echo.Context) error {
	slog.Info("A request to retrieve the team counts made")

	slog.Debug("Grabbing the user id from the token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get the user id from the token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	slog.Info("Retrieving the counts")
	counts, err := th.teamRepo.GetTeamCounts(c.Request().Context(), userID)
	if err != nil {
		slog.Error("Failed to retrieve the team counts")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch counts"})
	}

	slog.Info("Successfully retrieved counts, sending the response back")
	return c.JSON(http.StatusOK, counts)
}
