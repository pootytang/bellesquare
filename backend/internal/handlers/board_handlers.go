package handlers

import (
	"bellesquare-be/internal/models"
	repositories "bellesquare-be/internal/repository"
	"bellesquare-be/internal/utils"
	"bellesquare-be/internal/ws"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

type BoardHandler struct {
	boardRepo  repositories.BoardRepository
	payoutRepo repositories.PayoutRepository
	hub        *ws.Hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin should be more restrictive in production
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewBoardHandler(br repositories.BoardRepository, pr repositories.PayoutRepository, h *ws.Hub) *BoardHandler {
	return &BoardHandler{boardRepo: br, payoutRepo: pr, hub: h}
}

// POST /api/boards
func (bh *BoardHandler) CreateBoardHandler(c *echo.Context) error {
	slog.Info("Creating a board")

	// Get the subject and is_guest from the token
	slog.Debug("retrieving the user id and guest status")
	upro, err := utils.GetUserProfileFromJWT(c)
	if err != nil {
		slog.Error("Problem trying to get creatorID from token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
	}
	slog.Info("Retrieved the creatorID and guest from token", "id", upro.ID, "guest?", upro.IsGuest)

	slog.Info("Checking if this is a guest user")
	if upro.IsGuest {
		slog.Warn("Guests are not allowed to create boards")
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Upgrade your account to create boards!",
			"code":  "GUEST_RESTRICTION",
		})
	}

	req := new(models.CreateBoardRequest)
	if err := c.Bind(req); err != nil {
		slog.Error("Failed to bind board creation request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// VALIDATIONS
	slog.Info("Performing validations")
	if !req.Sport.IsValid() {
		slog.Warn("Invalide sport", "sport", req.Sport)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unsupported or unknown sport"})
	}

	// Logic Check: Total payouts must match total pool
	slog.Info("Checking to make sure the total payout is correct")
	totalPool := req.PricePerSquare * 100
	var payoutSum float64
	for _, amt := range req.Payouts {
		payoutSum += amt
	}

	// Allow for tiny floating point variance, but basically they should match
	if (payoutSum-totalPool) > 0.01 || (totalPool-payoutSum) > 0.01 {
		slog.Error("Payouts do not match the total pool", "total_pool", totalPool, "payout_sum", payoutSum)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Payout total ($%.2f) must match pool ($%.2f)", payoutSum, totalPool),
		})
	}

	// Map to model
	board := &models.Board{
		Title:          req.Title,
		Sport:          req.Sport,
		Status:         models.BoardOpen,
		PricePerSquare: req.PricePerSquare,
		HomeTeamID:     req.HomeTeamID,
		AwayTeamID:     req.AwayTeamID,
		CreatorID:      upro.ID,
	}

	// Create atomically (passing the PayoutRepo as a dependency)
	slog.Debug("Accessing the board DB", "board", board)
	if err := bh.boardRepo.CreateBoard(c.Request().Context(), board, bh.payoutRepo, req.Payouts); err != nil {
		slog.Error("Failed to create board", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Board creation failed"})
	}

	slog.Info("Successfully created the board", "title", board.Title)
	return c.JSON(http.StatusCreated, board)
}

// GET /api/boards/:id - This was using GetBoardWithSquares but now using GetBoardByBoardID
func (bh *BoardHandler) GetBoardHandler(c *echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid board ID"})
	}

	// Use the hydrated method that already joins users and calculates scores!
	response, err := bh.boardRepo.GetBoardByBoardID(c.Request().Context(), id)
	if err != nil {
		slog.Error("Failed to fetch hydrated board", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Board not found"})
	}

	// Get Payouts
	payouts, err := bh.payoutRepo.GetPayoutsByBoard(c.Request().Context(), id)
	if err != nil {
		slog.Error("Failed to fetch payouts", "board_id", id, "error", err)
		payouts = []models.Payout{}
	}

	// Attach the payouts to your hydrated response
	response.Payouts = payouts

	slog.Info("Successfully fetched board details", "board_id", response.Board.ID)
	return c.JSON(http.StatusOK, response)
}

func (bh *BoardHandler) GetUserBoardsHandler(c *echo.Context) error {
	slog.Info("Fetching boards for the authenticated user")

	// 1. Extract the user from the JWT (assuming middleware stores it as "user")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Problem trying to get user id from token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
	}

	slog.Info("Retrieved the user id from token", "id", userID)

	// 4. Query DB using the secured ID
	slog.Debug("Querying DB for secured user", "user_id", userID)
	boards, err := bh.boardRepo.GetBoardsByUserID(c.Request().Context(), userID)
	if err != nil {
		slog.Error("Failed to fetch boards", "user_id", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Retrieval failed"})
	}

	return c.JSON(http.StatusOK, boards)
}

func (bh *BoardHandler) HandleWS(c *echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		slog.Error("WS HUB - problem converting to uuid", "error", err)
		return err
	}

	boardID := id.String()
	slog.Info("WS HUB - Retrieved board id", "id", boardID) // This is now a normalized lowercase string

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		slog.Error("WS HUB - Problem upgrading the connection", "error", err)
		return err
	}

	slog.Debug("WS HUB - creating the lock and channel")
	bh.hub.Mu.Lock()
	if bh.hub.Rooms[boardID] == nil {
		bh.hub.Rooms[boardID] = make(map[*websocket.Conn]bool)
	}
	bh.hub.Rooms[boardID][ws] = true
	bh.hub.Mu.Unlock()

	// Keep connection alive/listen for disconnect
	defer func() {
		bh.hub.Mu.Lock()
		delete(bh.hub.Rooms[boardID], ws)
		bh.hub.Mu.Unlock()
		ws.Close()
	}()

	// We don't expect messages FROM the client for now,
	// but we need to keep the loop open to detect disconnects
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			break
		}
	}
	return nil
}

func (bh *BoardHandler) ClaimSquareHandler(c *echo.Context) error {
	slog.Info("Processing batch square claim")

	// 1. Get User Profile from Token
	slog.Debug("Grabbing the user profile from the token")
	userProfile, err := utils.GetUserProfileFromJWT(c)
	if err != nil {
		slog.Error("Unable to find the user profile in the token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// 2. Extract Board ID
	slog.Debug("grabbing boardId")
	boardID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		slog.Warn("Unable to retrieve the board id", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid board ID"})
	}

	// 3. Bind the Request
	slog.Debug("Binding to the claimSquareRequest dto")
	req := new(models.ClaimSquareRequest)
	if err := c.Bind(req); err != nil {
		slog.Error("Failed to claim square request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// 4. Validation
	slog.Debug("Performing validations")
	if len(req.Selections) == 0 {
		slog.Debug("No squares selected")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No squares selected"})
	}

	// 5. Call the Repo Batch Method using the transaction logic
	slog.Debug("Calling the batch method")
	err = bh.boardRepo.ClaimSquaresBatch(
		c.Request().Context(),
		boardID,
		userProfile.ID,
		req.UserColor,
		req.Selections,
	)

	if err != nil {
		// Check for the specific error returned by your transaction logic
		if strings.Contains(err.Error(), "already claimed") {
			slog.Warn("Conflict: One or more squares already taken", "board_id", boardID)
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}

		slog.Error("Failed to execute batch claim", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process selection"})
	}

	// 6. Configure the payload
	slog.Debug("Setting up the payload")
	var broadcastPayload []models.SquareUpdate
	for _, sel := range req.Selections {
		broadcastPayload = append(broadcastPayload, models.SquareUpdate{
			Row:          sel.Row, // This maps to json:"row_index" in SquareUpdate
			Col:          sel.Col, // This maps to json:"col_index" in SquareUpdate
			UserID:       userProfile.ID,
			UserColor:    req.UserColor,
			UserInitials: userProfile.Initials,
			// Fields required by the SquareUpdate DTO
			FirstName: "",
			LastName:  "",
			IsPaid:    false, // New claims start as unpaid
		})
	}

	slog.Debug("Attempting to broadcast to hub")
	bh.hub.Broadcast <- ws.BroadcastMessage{
		BoardID: boardID.String(),
		Payload: broadcastPayload, // The actual square data (color, user_id, initials)
	}

	slog.Info("Batch claim successful", "board_id", boardID, "count", len(req.Selections))
	return c.JSON(http.StatusOK, map[string]string{"message": "Squares claimed successfully"})
}

func (bh *BoardHandler) GenerateAxisNumbersHandler(c *echo.Context) error {
	slog.Info("Generating Axis Numbers")
	boardID, _ := uuid.FromString(c.Param("id"))

	// 1. Security: Ensure the person clicking is the creator
	slog.Debug("retrieving the userID from the jwt")
	userID, err := utils.GetUserIDFromJWT(c)

	slog.Debug("retrieving the board", "boardID", boardID)
	boardDR, err := bh.boardRepo.GetBoardByBoardID(c.Request().Context(), boardID)
	if err != nil {
		slog.Error("Unable to retrieve the board", "boardID", boardID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to retrieve the board"})
	}

	if boardDR.Board.CreatorID != userID {
		slog.Error("Board and User mismatch", "boards creator", boardDR.Board.CreatorID, "userID", userID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only the creator can generate numbers"})
	}

	// 2. Generate Shuffled Arrays
	slog.Info("Generating the axis numbers")
	homeAxis := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	awayAxis := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	rand.Shuffle(len(homeAxis), func(i, j int) { homeAxis[i], homeAxis[j] = homeAxis[j], homeAxis[i] })
	rand.Shuffle(len(awayAxis), func(i, j int) { awayAxis[i], awayAxis[j] = awayAxis[j], awayAxis[i] })

	// 3. Save to DB
	slog.Info("Storing the axis numbers")
	slog.Debug("Axis numbers", "home", homeAxis, "away", awayAxis)
	err = bh.boardRepo.UpdateAxisNumbers(c.Request().Context(), boardID, homeAxis, awayAxis)
	if err != nil {
		slog.Warn("Failed to save the axis numbers", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save numbers"})
	}

	// 4. Broadcast via WebSocket
	// We send a specific payload type so Svelte knows to update the headers
	slog.Debug("sending numbers through websocket")
	bh.hub.Broadcast <- ws.BroadcastMessage{
		BoardID: boardID.String(),
		Payload: map[string]interface{}{
			"type": "AXIS_UPDATE",
			"home": homeAxis,
			"away": awayAxis,
		},
	}

	slog.Info("Successfully updated the axis numbers")
	return c.JSON(http.StatusOK, map[string]string{"message": "Numbers generated"})
}

func (bh *BoardHandler) UpdateStatusHandler(c *echo.Context) error {
	slog.Info("********** UPDATE BOARD STATUS CALLED **********")

	slog.Debug("Board to update", "id", c.Param("id"))
	boardID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		slog.Error("Board ID not found or is invalid")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid board ID"})
	}

	// GET THE STATUS FROM THE REQUEST
	slog.Info("Grabbing the status from the request")
	var req models.BoardStatusRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Unable to bind the request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Ensure we aren't accidentally setting 'unknown' (result of failed unmarshal)
	if req.Status.String() == "unknown" {
		slog.Warn("Request status is unknown")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid status value"})
	}

	// GET THE USER ID FROM THE TOKEN
	slog.Debug("Grabbing the user id from the JWT")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get user id from jwt token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// CHECK THAT THE BOARD IS OWNED BY THE USER
	slog.Info("grabbing the board from the db")
	board, err := bh.boardRepo.GetBoardByBoardID(c.Request().Context(), boardID)
	if err != nil {
		slog.Error("Unable to get board", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Board not found"})
	}

	slog.Debug("checking if this is the creator")
	if board.Board.CreatorID != userID {
		slog.Error("Board Creator id doesn't match the userID", "creator id", board.Board.CreatorID, "user id", userID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only the creator can lock the board"})
	}

	// UPDATE THE DB WITH THE NEW STATUS
	slog.Info("Updating the board status")
	err = bh.boardRepo.UpdateBoardStatus(c.Request().Context(), boardID, req.Status)
	if err != nil {
		slog.Error("Problem updating the status", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to lock board"})
	}

	// UTILIZE THE WEBSOCKET TO SEND THE STATUS MESSAGE
	slog.Info("Broadcasting the status through the websocket")
	bh.hub.Broadcast <- ws.BroadcastMessage{
		BoardID: boardID.String(),
		Payload: map[string]interface{}{
			"type":   "STATUS_UPDATE",
			"status": req.Status.String(),
		},
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Board locked successfully"})
}

func (bh *BoardHandler) SavePayoutsHandler(c *echo.Context) error {
	slog.Info("Saving Payouts")
	boardID, _ := uuid.FromString(c.Param("id"))
	userID, _ := utils.GetUserIDFromJWT(c) // Use your existing helper

	// 1. Verify Board Ownership
	slog.Info("verifying the owner of the board")
	board, err := bh.boardRepo.GetBoardByBoardID(c.Request().Context(), boardID)
	if err != nil || board.Board.CreatorID != userID {
		slog.Error("This is not the board owner", "error", err)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// 2. Bind the Payouts
	slog.Info("grabbing the payout data from the request")
	var req models.PayoutSaveRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Unable to bind payout data", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payout data"})
	}

	// 3. Save to DB
	slog.Info("Saving the payout data to the db")
	err = bh.payoutRepo.UpsertPayouts(c.Request().Context(), boardID, req.Payouts)
	if err != nil {
		slog.Error("Problem saving payout data", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save payouts"})
	}

	slog.Info("Successfully saved payout data")
	return c.JSON(http.StatusOK, map[string]string{"message": "Payouts updated"})
}

func (bh *BoardHandler) TogglePaidHandler(c *echo.Context) error {
	slog.Info("Toggle Paid Handler called")

	slog.Debug("Grabbing data from the request")
	boardID, _ := uuid.FromString(c.Param("boardId"))
	targetUserID, _ := uuid.FromString(c.Param("userId"))

	// safeguard
	currentUserID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get userId from token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User id not found"})
	}
	slog.Debug("Target user and Creator", "Target User ID", targetUserID, "Creator ID", currentUserID)

	// Verify Requesting User is the Creator
	slog.Debug("Checking if this is the creator")
	board, err := bh.boardRepo.GetBoardByBoardID(c.Request().Context(), boardID)
	if err != nil || board.Board.CreatorID != currentUserID {
		slog.Error("This is not the creator", "error", err)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only the creator can manage payments"})
	}

	// Perform Toggle
	slog.Info("Toggling")
	err = bh.boardRepo.ToggleUserPaidStatus(c.Request().Context(), boardID, targetUserID)
	if err != nil {
		slog.Error("Failed to toggle the payment status", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
	}

	// Fetch the squares for the target User
	slog.Info("Grabbing the updated squares for user", "target user id", targetUserID)
	updatedSquares, err := bh.boardRepo.GetSquaresByBoardAndUser(c.Request().Context(), boardID, targetUserID)
	if err != nil {
		slog.Warn("Unable to get the updated squares", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get updated squares"})
	}

	// Broadcast the message
	slog.Info("Broadcasting the update")
	bh.hub.BroadcastSquareUpdate(boardID, updatedSquares)

	slog.Info("Returning statusOK")
	return c.NoContent(http.StatusOK)
}

func (bh *BoardHandler) UpdateScoreHandler(c *echo.Context) error {
	slog.Info("---------- UPDATE SCORE HANDLER REQUESTED ---------")
	boardID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		slog.Error("Problem retrieving board id from the request", "error", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	slog.Debug("Grabbing the score update from the request")
	var req models.UpdateScoreRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Problem getting the score update from the request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// safeguard
	slog.Debug("Getting the creator from the token")
	currentUserID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get userId from token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User id not found"})
	}

	slog.Debug("Checking if this is the creator of the board", "board id", boardID, "creator id", currentUserID)
	board, err := bh.boardRepo.GetBoardByBoardID(c.Request().Context(), boardID)
	if err != nil || board.Board.CreatorID != currentUserID {
		slog.Error("This is not the creator", "error", err)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only the creator can manage payments"})
	}
	slog.Info("This is the creator of the board that is attempting to update the score")

	// Save to board_scores table (Upsert)
	slog.Info("Saving the score update")
	slog.Debug("querying the db to save the score", "board id", boardID, "quarter", req.Quarter, "home score", req.HomeScore, "away score", req.AwayScore)
	err = bh.boardRepo.UpdateQuarterScore(c.Request().Context(), boardID, req.IsLocked, req.Quarter, req.HomeScore, req.AwayScore)
	if err != nil {
		slog.Error("Unable to update the score", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB update failed"})
	}

	// Broadcast to WebSocket
	slog.Info("Sending SCORE_UPDATE broadcast")
	bh.hub.BroadcastScoreUpdate(boardID, req)

	return c.NoContent(http.StatusOK)
}

func (h *BoardHandler) GetJoinedBoards(c *echo.Context) error {
	slog.Info("---------- GET JOINED BOARDS HANDLER CALLED ---------")

	slog.Debug("Grabbing userid from token")
	userid, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get userId from token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User id not found"})
	}

	slog.Debug("Querying DB for joined boards", "user id", userid)
	boards, err := h.boardRepo.GetJoinedBoards(c.Request().Context(), userid)
	if err != nil {
		slog.Error("Failed to fetch joined boards", "user id", userid, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch joined boards"})
	}

	slog.Info("Successfully retrieved boards", "count", len(boards))
	return c.JSON(http.StatusOK, boards)
}

func (bh *BoardHandler) MarkPaymentSentHandler(c *echo.Context) error {
	slog.Info("---------- MARK PAYMENT SENT HANDLER CALLED ---------")

	slog.Debug("Grabbing board id from parameters")
	boardID, err := uuid.FromString(c.Param("boardId"))
	if err != nil {
		slog.Error("Problem retrieving board id from the request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	slog.Debug("Grabbing user id from token")
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get userId from token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User id not found"})
	}

	// Set status to 1 (Pending)
	slog.Info("Updating the payment status to pending", "board id", boardID, "user id", userID)
	err = bh.boardRepo.UpdateUserPaymentStatus(c.Request().Context(), boardID, userID, models.Pending)
	if err != nil {
		slog.Error("Failed to update user payment status", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Update failed"})
	}

	// Broadcast the update via WS so the creator's list pulses immediately
	slog.Info("Broadcasting square update", "board id", boardID, "user id", userID)
	updatedSquares, _ := bh.boardRepo.GetSquaresByBoardAndUser(c.Request().Context(), boardID, userID)
	slog.Debug("Payment status", "status", updatedSquares[0].PaymentStatus)
	bh.hub.BroadcastSquareUpdate(boardID, updatedSquares)

	return c.NoContent(http.StatusOK)
}

// Used by the creator to mark a payment as verified (status=2) or rejected (status=0)
func (bh *BoardHandler) VerifyPaymentHandler(c *echo.Context) error {
	slog.Info("---------- VERIFY PAYMENT HANDLER CALLED ---------")

	slog.Info("Grabbing board id and target user id from parameters")
	boardID, err := uuid.FromString(c.Param("boardId"))
	if err != nil {
		slog.Error("Problem retrieving board id from the request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid board ID"})
	}

	targetUserID, err := uuid.FromString(c.Param("userId"))
	if err != nil {
		slog.Error("Problem retrieving target user id from the request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// Get status from request (e.g., { "status": 2 })
	slog.Info("Retrieving the status from the request")
	var req struct {
		Status int `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		slog.Error("Unable to bind the request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Authorization: Ensure requesting user is the creator
	slog.Info("Checking if the requester is the creator of the board")
	currentUserID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		slog.Error("Unable to get userId from token", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User id not found"})
	}
	slog.Debug("Retrieved creator id from token", "id", currentUserID)

	slog.Debug("Grabbing the board from the DB")
	board, err := bh.boardRepo.GetBoardByBoardID(c.Request().Context(), boardID)
	if err != nil {
		slog.Error("Unable to get board", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Board not found"})
	}

	slog.Debug("Retrieved the board", "id", board.Board.ID)
	if board.Board.CreatorID != currentUserID {
		slog.Error("Requesting user is not the creator of the board", "user id", currentUserID, "creator id", board.Board.CreatorID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only creators can verify payments"})
	}

	// Cast the int from request to your models.PaymentStatus type
	slog.Info("Requester is the creator, proceeding with payment status update")
	err = bh.boardRepo.UpdateUserPaymentStatus(c.Request().Context(), boardID, targetUserID, models.PaymentStatus(req.Status))
	if err != nil {
		slog.Error("Failed to update user payment status", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Update failed"})
	}

	// Broadcast the update so everyone's UI stays in sync
	slog.Info("Retrieving updated squares for the target user to broadcast")
	updatedSquares, err := bh.boardRepo.GetSquaresByBoardAndUser(c.Request().Context(), boardID, targetUserID)
	if err != nil {
		slog.Warn("Failed to get updated squares", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve updated squares"})
	}

	slog.Info("Broadcasting the payment status update through the websocket")
	bh.hub.BroadcastSquareUpdate(boardID, updatedSquares)

	slog.Info("Payment status updated successfully")
	return c.NoContent(http.StatusOK)
}
