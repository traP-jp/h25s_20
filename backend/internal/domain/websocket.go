package domain

import (
	"context"
	"net/http"
)

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection interface {
	// ReadMessage reads a message from the WebSocket connection
	ReadMessage(ctx context.Context) (messageType int, data []byte, err error)

	// WriteMessage writes a message to the WebSocket connection
	WriteMessage(ctx context.Context, messageType int, data []byte) error

	// Close closes the WebSocket connection
	Close() error

	// GetUserID returns the user ID associated with this connection
	GetUserID() string

	// GetRoomID returns the room ID associated with this connection
	GetRoomID() string

	// SetRoomID sets the room ID for this connection
	SetRoomID(roomID string)
}

// WebSocketManager manages WebSocket connections
type WebSocketManager interface {
	// UpgradeConnection upgrades an HTTP connection to WebSocket
	UpgradeConnection(w http.ResponseWriter, r *http.Request, userID string) (WebSocketConnection, error)

	// AddConnection adds a connection to a room
	AddConnection(roomID string, conn WebSocketConnection) error

	// RemoveConnection removes a connection from a room
	RemoveConnection(roomID string, userID string) error

	// BroadcastToRoom broadcasts a message to all connections in a room
	BroadcastToRoom(roomID string, messageType int, data []byte) error

	// SendToUser sends a message to a specific user
	SendToUser(userID string, messageType int, data []byte) error

	// GetRoomConnections returns all connections in a room
	GetRoomConnections(roomID string) []WebSocketConnection

	// GetUserConnection returns a specific user's connection
	GetUserConnection(userID string) WebSocketConnection
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string      `json:"type"`
	RoomID    string      `json:"roomId,omitempty"`
	UserID    string      `json:"userId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// WebSocketMessageType defines message types for server-to-client notifications
type WebSocketMessageType string

const (
	MessageTypeRoomUpdate   WebSocketMessageType = "room_update"
	MessageTypeGameUpdate   WebSocketMessageType = "game_update"
	MessageTypeBoardUpdate  WebSocketMessageType = "board_update"
	MessageTypeGameStart    WebSocketMessageType = "game_start"
	MessageTypeGameEnd      WebSocketMessageType = "game_end"
	MessageTypeUserJoined   WebSocketMessageType = "user_joined"
	MessageTypeUserLeft     WebSocketMessageType = "user_left"
	MessageTypeNotification WebSocketMessageType = "notification"
	MessageTypeError        WebSocketMessageType = "error"
	MessageTypePing         WebSocketMessageType = "ping"
	MessageTypePong         WebSocketMessageType = "pong"
)
