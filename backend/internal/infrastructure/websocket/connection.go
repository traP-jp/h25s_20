package websocket

import (
	"context"
	"sync"

	"github.com/coder/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
)

// connection implements domain.WebSocketConnection
type connection struct {
	conn   *websocket.Conn
	userID string
	roomID string
	mu     sync.RWMutex
}

// NewConnection creates a new WebSocket connection wrapper
func NewConnection(conn *websocket.Conn, userID string) domain.WebSocketConnection {
	return &connection{
		conn:   conn,
		userID: userID,
		roomID: "",
	}
}

// ReadMessage reads a message from the WebSocket connection
func (c *connection) ReadMessage(ctx context.Context) (messageType int, data []byte, err error) {
	msgType, data, err := c.conn.Read(ctx)
	return int(msgType), data, err
}

// WriteMessage writes a message to the WebSocket connection
func (c *connection) WriteMessage(ctx context.Context, messageType int, data []byte) error {
	return c.conn.Write(ctx, websocket.MessageType(messageType), data)
}

// Close closes the WebSocket connection
func (c *connection) Close() error {
	return c.conn.CloseNow()
}

// GetUserID returns the user ID associated with this connection
func (c *connection) GetUserID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userID
}

// GetRoomID returns the room ID associated with this connection
func (c *connection) GetRoomID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.roomID
}

// SetRoomID sets the room ID for this connection
func (c *connection) SetRoomID(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.roomID = roomID
}
