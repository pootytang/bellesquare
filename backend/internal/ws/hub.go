package ws

import (
	"bellesquare-be/internal/models"
	"log/slog"
	"sync"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

type Hub struct {
	// Maps BoardID -> map of ClientConnections
	// We use a nested map so we can broadcast to specific boards only
	Rooms     map[string]map[*websocket.Conn]bool
	Broadcast chan BroadcastMessage
	Mu        sync.Mutex
}

type BroadcastMessage struct {
	BoardID string      `json:"board_id"`
	Payload interface{} `json:"payload"` // This will be your Square object
}

func NewHub() *Hub {
	return &Hub{
		Rooms:     make(map[string]map[*websocket.Conn]bool),
		Broadcast: make(chan BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		msg := <-h.Broadcast
		h.Mu.Lock()
		slog.Debug("WS HUB - Hub received message from channel", "board_id", msg.BoardID)

		if clients, ok := h.Rooms[msg.BoardID]; ok {
			slog.Debug("WS HUB - Found active room", "client_count", len(clients))
			for client := range clients {
				err := client.WriteJSON(msg.Payload)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		} else {
			// IF THIS LOGS: Your BoardID string format likely doesn't match
			// how the client registered (e.g. UUID string vs raw string)
			slog.Warn("WS HUB - No active room found for board ID", "board_id", msg.BoardID)
		}
		h.Mu.Unlock()
	}
}

func (h *Hub) BroadcastSquareUpdate(boardID uuid.UUID, updatedSquares []models.SquareUpdate) {
	slog.Info("WS HUB - Hub: Broadcasting square update", "board_id", boardID, "count", len(updatedSquares))

	// This ensures we match the normalized string used in HandleWS
	roomKey := boardID.String()

	slog.Info("WS HUB - Broadcasting square update", "room_key", roomKey)

	// Construct the message to send to the channel
	msg := BroadcastMessage{
		BoardID: roomKey,
		Payload: map[string]interface{}{
			"type":    "SQUARE_UPDATE",
			"squares": updatedSquares,
		},
	}

	// Drop it into the channel; the Run() loop will handle the fan-out
	h.Broadcast <- msg
}

func (h *Hub) BroadcastScoreUpdate(boardID uuid.UUID, update models.UpdateScoreRequest) {
	slog.Info("WS HUB - Score update", "board id", boardID, "quarter", update.Quarter, "Home Score", update.HomeScore, "Away Score", update.AwayScore)
	h.Broadcast <- BroadcastMessage{
		BoardID: boardID.String(),
		Payload: map[string]interface{}{
			"type":       "SCORE_UPDATE",
			"home_score": update.HomeScore,
			"away_score": update.AwayScore,
			"quarter":    update.Quarter,
			"is_locked":  update.IsLocked,
		},
	}
}
