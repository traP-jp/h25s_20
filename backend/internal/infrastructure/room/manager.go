package room

import (
	"fmt"
	"sync"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/rs/zerolog/log"
)

// manager implements domain.RoomManager
type manager struct {
	rooms  map[int]*domain.Room
	nextID int
	mu     sync.RWMutex
}

// NewManager creates a new room manager
func NewManager() domain.RoomManager {
	return &manager{
		rooms:  make(map[int]*domain.Room),
		nextID: 1,
	}
}

// CreateRoom creates a new room
func (m *manager) CreateRoom(name string) (*domain.Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if name == "" {
		name = fmt.Sprintf("Room %d", m.nextID)
	}

	room := &domain.Room{
		ID:        m.nextID,
		Name:      name,
		Boards:    make([]domain.Board, 0),
		IsOpened:  true,
		Players:   make([]domain.Player, 0),
		ResultLog: make([]domain.Result, 0),
	}

	m.rooms[m.nextID] = room
	roomID := m.nextID
	m.nextID++

	log.Info().Int("roomID", roomID).Str("name", name).Msg("Room created")
	return room, nil
}

// GetRoom retrieves a room by ID
func (m *manager) GetRoom(roomID int) (*domain.Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room %d not found", roomID)
	}

	return room, nil
}

// UpdateRoom updates an existing room
func (m *manager) UpdateRoom(roomID int, room *domain.Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rooms[roomID]; !exists {
		return fmt.Errorf("room %d not found", roomID)
	}

	m.rooms[roomID] = room
	log.Debug().Int("roomID", roomID).Msg("Room updated")
	return nil
}

// DeleteRoom deletes a room
func (m *manager) DeleteRoom(roomID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rooms[roomID]; !exists {
		return fmt.Errorf("room %d not found", roomID)
	}

	delete(m.rooms, roomID)
	log.Info().Int("roomID", roomID).Msg("Room deleted")
	return nil
}

// ListRooms returns all rooms
func (m *manager) ListRooms() ([]*domain.Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rooms := make([]*domain.Room, 0, len(m.rooms))
	for _, room := range m.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// JoinRoom adds a player to a room
func (m *manager) JoinRoom(roomID int, player domain.Player) (*domain.Room, error) {
	m.mu.RLock()
	room, exists := m.rooms[roomID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("room %d not found", roomID)
	}

	if !room.IsOpened {
		return nil, fmt.Errorf("room %d is closed", roomID)
	}

	player.JoinedAt = time.Now()
	if err := room.AddPlayer(player); err != nil {
		return nil, fmt.Errorf("failed to join room: %w", err)
	}

	log.Info().Int("roomID", roomID).Str("playerID", player.ID).Str("playerName", player.Name).Msg("Player joined room")
	return room, nil
}

// LeaveRoom removes a player from a room
func (m *manager) LeaveRoom(roomID int, playerID string) (*domain.Room, error) {
	m.mu.RLock()
	room, exists := m.rooms[roomID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("room %d not found", roomID)
	}

	if err := room.RemovePlayer(playerID); err != nil {
		return nil, fmt.Errorf("failed to leave room: %w", err)
	}

	// Delete room if empty
	if room.PlayerCount() == 0 {
		m.DeleteRoom(roomID)
		log.Info().Int("roomID", roomID).Msg("Empty room deleted")
	}

	log.Info().Int("roomID", roomID).Str("playerID", playerID).Msg("Player left room")
	return room, nil
}
