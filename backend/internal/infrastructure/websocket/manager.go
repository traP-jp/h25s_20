package websocket

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/rs/zerolog/log"
)

// manager implements domain.WebSocketManager
type manager struct {
	// rooms maps room ID to user connections
	rooms map[string]map[string]domain.WebSocketConnection
	// users maps user ID to their connection
	users map[string]domain.WebSocketConnection
	mu    sync.RWMutex
}

// NewManager creates a new WebSocket manager
func NewManager() domain.WebSocketManager {
	return &manager{
		rooms: make(map[string]map[string]domain.WebSocketConnection),
		users: make(map[string]domain.WebSocketConnection),
	}
}

// UpgradeConnection upgrades an HTTP connection to WebSocket
func (m *manager) UpgradeConnection(w http.ResponseWriter, r *http.Request, userID string) (domain.WebSocketConnection, error) {
	// Basic validation
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	// Check if user is already connected
	m.mu.RLock()
	if existingConn, exists := m.users[userID]; exists {
		m.mu.RUnlock()
		// Close existing connection
		existingConn.Close()
		log.Warn().Str("userID", userID).Msg("Replacing existing WebSocket connection")
	} else {
		m.mu.RUnlock()
	}

	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,          // For development only
		OriginPatterns:     []string{"*"}, // Allow all origins for development
	}

	conn, err := websocket.Accept(w, r, options)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade connection: %w", err)
	}

	wsConn := NewConnection(conn, userID)

	m.mu.Lock()
	m.users[userID] = wsConn
	m.mu.Unlock()

	log.Info().Str("userID", userID).Msg("WebSocket connection established")

	return wsConn, nil
}

// AddConnection adds a connection to a room
func (m *manager) AddConnection(roomID string, conn domain.WebSocketConnection) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.rooms[roomID] == nil {
		m.rooms[roomID] = make(map[string]domain.WebSocketConnection)
	}

	userID := conn.GetUserID()
	m.rooms[roomID][userID] = conn
	conn.SetRoomID(roomID)

	log.Info().Str("userID", userID).Str("roomID", roomID).Msg("User joined room")
	return nil
}

// RemoveConnection removes a connection from a room
func (m *manager) RemoveConnection(roomID string, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.rooms[roomID] != nil {
		if conn, exists := m.rooms[roomID][userID]; exists {
			conn.SetRoomID("")
			delete(m.rooms[roomID], userID)

			// Remove room if empty
			if len(m.rooms[roomID]) == 0 {
				delete(m.rooms, roomID)
			}

			log.Info().Str("userID", userID).Str("roomID", roomID).Msg("User left room")
		}
	}

	// Remove from users map
	if conn, exists := m.users[userID]; exists {
		conn.Close()
		delete(m.users, userID)
	}

	return nil
}

// BroadcastToRoom broadcasts a message to all connections in a room
func (m *manager) BroadcastToRoom(roomID string, messageType int, data []byte) error {
	// Get connections with lock held to ensure consistency
	m.mu.RLock()
	connections := make([]domain.WebSocketConnection, 0)
	if room, exists := m.rooms[roomID]; exists {
		for _, conn := range room {
			connections = append(connections, conn)
		}
	}
	m.mu.RUnlock()

	if len(connections) == 0 {
		log.Debug().Str("roomID", roomID).Msg("No connections to broadcast to")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup
	failedUsers := make([]string, 0)
	var failedMu sync.Mutex

	for _, conn := range connections {
		wg.Add(1)
		go func(c domain.WebSocketConnection) {
			defer wg.Done()

			if err := c.WriteMessage(ctx, messageType, data); err != nil {
				log.Error().Err(err).Str("userID", c.GetUserID()).Str("roomID", roomID).Msg("Failed to send message to user")

				failedMu.Lock()
				failedUsers = append(failedUsers, c.GetUserID())
				failedMu.Unlock()
			}
		}(conn)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Clean up failed connections
	for _, userID := range failedUsers {
		m.RemoveConnection(roomID, userID)
	}

	if len(failedUsers) > 0 {
		log.Warn().Str("roomID", roomID).Int("failedCount", len(failedUsers)).Msg("Some connections failed during broadcast")
	}

	return nil
}

// SendToUser sends a message to a specific user
func (m *manager) SendToUser(userID string, messageType int, data []byte) error {
	m.mu.RLock()
	conn, exists := m.users[userID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.WriteMessage(ctx, messageType, data); err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("Failed to send message to user")
		// Remove failed connection
		m.RemoveConnection(conn.GetRoomID(), userID)
		return err
	}

	return nil
}

// GetRoomConnections returns all connections in a room
func (m *manager) GetRoomConnections(roomID string) []domain.WebSocketConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connections := make([]domain.WebSocketConnection, 0)
	if room, exists := m.rooms[roomID]; exists {
		for _, conn := range room {
			connections = append(connections, conn)
		}
	}

	return connections
}

// GetUserConnection returns a specific user's connection
func (m *manager) GetUserConnection(userID string) domain.WebSocketConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.users[userID]
}
