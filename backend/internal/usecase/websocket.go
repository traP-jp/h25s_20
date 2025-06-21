package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/rs/zerolog/log"
)

type WebSocketUsecase struct {
	wsManager   domain.WebSocketManager
	roomManager domain.RoomManager
}

func NewWebSocketUsecase(wsManager domain.WebSocketManager, roomManager domain.RoomManager) *WebSocketUsecase {
	return &WebSocketUsecase{
		wsManager:   wsManager,
		roomManager: roomManager,
	}
}

// HandleConnection handles a new WebSocket connection for one-way notifications
func (w *WebSocketUsecase) HandleConnection(ctx context.Context, conn domain.WebSocketConnection) {
	userID := conn.GetUserID()
	log.Info().Str("userID", userID).Msg("Starting WebSocket connection handler")

	defer func() {
		if roomID := conn.GetRoomID(); roomID != "" {
			w.wsManager.RemoveConnection(roomID, userID)
		}
		conn.Close()
		log.Info().Str("userID", userID).Msg("WebSocket connection closed")
	}()

	// Send welcome message
	welcomeMsg := domain.WebSocketMessage{
		Type:      string(domain.MessageTypeRoomUpdate),
		UserID:    userID,
		Data:      map[string]string{"status": "connected"},
		Timestamp: time.Now().Unix(),
	}

	if err := w.sendMessage(ctx, conn, welcomeMsg); err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("Failed to send welcome message")
		return
	}

	// Keep connection alive and handle only basic connection management
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Set a read deadline to detect connection issues
			_, data, err := conn.ReadMessage(ctx)
			if err != nil {
				log.Error().Err(err).Str("userID", userID).Msg("Connection lost or read error")
				return
			}

			// For one-way notifications, we mainly ignore incoming messages
			// Only handle ping messages for connection keep-alive
			var msg domain.WebSocketMessage
			if json.Unmarshal(data, &msg) == nil && msg.Type == string(domain.MessageTypePing) {
				w.handlePing(ctx, conn, msg)
			}
		}
	}
}

// JoinRoom adds a user to a room (called from REST API)
func (w *WebSocketUsecase) JoinRoom(userID string, roomID int, userName string) error {
	conn := w.wsManager.GetUserConnection(userID)
	if conn == nil {
		return fmt.Errorf("user %s not connected", userID)
	}

	// Create player
	player := domain.Player{
		ID:      userID,
		Name:    userName,
		IsReady: false,
		Score:   0,
	}

	// Join room through room manager
	room, err := w.roomManager.JoinRoom(roomID, player)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	// Add WebSocket connection to room
	roomIDStr := strconv.Itoa(roomID)
	if err := w.wsManager.AddConnection(roomIDStr, conn); err != nil {
		// Rollback room join if WebSocket connection fails
		w.roomManager.LeaveRoom(roomID, userID)
		return fmt.Errorf("failed to join WebSocket room: %w", err)
	}

	// Notify room members
	roomUpdate := domain.WebSocketMessage{
		Type:   string(domain.MessageTypeUserJoined),
		RoomID: roomIDStr,
		UserID: userID,
		Data: map[string]interface{}{
			"action": "joined",
			"player": player,
			"room":   room,
		},
		Timestamp: time.Now().Unix(),
	}

	ctx := context.Background()
	return w.broadcastToRoom(ctx, roomIDStr, roomUpdate)
}

// LeaveRoom removes a user from a room (called from REST API)
func (w *WebSocketUsecase) LeaveRoom(userID string, roomID int) error {
	roomIDStr := strconv.Itoa(roomID)

	// Remove from room manager
	room, err := w.roomManager.LeaveRoom(roomID, userID)
	if err != nil {
		log.Warn().Err(err).Str("userID", userID).Int("roomID", roomID).Msg("Failed to leave room in room manager")
	}

	// Remove WebSocket connection from room
	if err := w.wsManager.RemoveConnection(roomIDStr, userID); err != nil {
		return fmt.Errorf("failed to leave WebSocket room: %w", err)
	}

	// Notify room members
	roomUpdate := domain.WebSocketMessage{
		Type:   string(domain.MessageTypeUserLeft),
		RoomID: roomIDStr,
		UserID: userID,
		Data: map[string]interface{}{
			"action": "left",
			"room":   room,
		},
		Timestamp: time.Now().Unix(),
	}

	ctx := context.Background()
	return w.broadcastToRoom(ctx, roomIDStr, roomUpdate)
}

// NotifyRoomUpdate notifies all users in a room about a room state change
func (w *WebSocketUsecase) NotifyRoomUpdate(roomID int, updateType string, data interface{}) error {
	roomIDStr := strconv.Itoa(roomID)

	room, err := w.roomManager.GetRoom(roomID)
	if err != nil {
		log.Warn().Err(err).Int("roomID", roomID).Msg("Room not found for notification")
		return err
	}

	update := domain.WebSocketMessage{
		Type:   updateType,
		RoomID: roomIDStr,
		Data: map[string]interface{}{
			"room":   room,
			"update": data,
		},
		Timestamp: time.Now().Unix(),
	}

	ctx := context.Background()
	return w.broadcastToRoom(ctx, roomIDStr, update)
}

// NotifyBoardUpdate notifies all users in a room about board changes
func (w *WebSocketUsecase) NotifyBoardUpdate(roomID int, board domain.Board) error {
	return w.NotifyRoomUpdate(roomID, string(domain.MessageTypeBoardUpdate), board)
}

// NotifyGameStart notifies all users in a room that the game has started
func (w *WebSocketUsecase) NotifyGameStart(roomID int) error {
	return w.NotifyRoomUpdate(roomID, string(domain.MessageTypeGameStart), map[string]string{"status": "started"})
}

// NotifyGameEnd notifies all users in a room that the game has ended
func (w *WebSocketUsecase) NotifyGameEnd(roomID int, results []domain.Result) error {
	return w.NotifyRoomUpdate(roomID, string(domain.MessageTypeGameEnd), map[string]interface{}{
		"status":  "ended",
		"results": results,
	})
}

// handlePing handles ping messages
func (w *WebSocketUsecase) handlePing(ctx context.Context, conn domain.WebSocketConnection, msg domain.WebSocketMessage) error {
	pongMsg := domain.WebSocketMessage{
		Type:      string(domain.MessageTypePong),
		UserID:    msg.UserID,
		Data:      msg.Data,
		Timestamp: time.Now().Unix(),
	}

	return w.sendMessage(ctx, conn, pongMsg)
}

// sendMessage sends a message to a specific connection
func (w *WebSocketUsecase) sendMessage(ctx context.Context, conn domain.WebSocketConnection, msg domain.WebSocketMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return conn.WriteMessage(ctx, int(websocket.MessageText), data)
}

// sendErrorMessage sends an error message to a connection
func (w *WebSocketUsecase) sendErrorMessage(ctx context.Context, conn domain.WebSocketConnection, errorMsg string) {
	msg := domain.WebSocketMessage{
		Type:      string(domain.MessageTypeError),
		UserID:    conn.GetUserID(),
		Data:      map[string]string{"error": errorMsg},
		Timestamp: time.Now().Unix(),
	}

	if err := w.sendMessage(ctx, conn, msg); err != nil {
		log.Error().Err(err).Str("userID", conn.GetUserID()).Msg("Failed to send error message")
	}
}

// broadcastToRoom broadcasts a message to all connections in a room
func (w *WebSocketUsecase) broadcastToRoom(ctx context.Context, roomID string, msg domain.WebSocketMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return w.wsManager.BroadcastToRoom(roomID, int(websocket.MessageText), data)
}
