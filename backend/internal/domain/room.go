package domain

import (
	"fmt"
	"sync"
)

type Room struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Boards    []Board  `json:"boards"`
	IsOpened  bool     `json:"isOpened"`
	Players   []Player `json:"players"`
	ResultLog []Result `json:"resultLog"`
	mu        sync.RWMutex
}

type Board struct {
	Version int   `json:"version"`
	Board   []int `json:"board"`
}

type Result struct {
	ID     int           `json:"id"`
	Time   string        `json:"time"`
	Scores []PlayerScore `json:"scores"`
}

type PlayerScore struct {
	ID       int    `json:"id"`
	PlayerId string `json:"playerId"`
	Score    int    `json:"score"`
}

// Room methods for thread-safe operations
func (r *Room) AddPlayer(player Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if player already exists
	for _, p := range r.Players {
		if p.ID == player.ID {
			return fmt.Errorf("player %s already in room", player.ID)
		}
	}

	r.Players = append(r.Players, player)
	return nil
}

func (r *Room) RemovePlayer(playerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, player := range r.Players {
		if player.ID == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("player %s not found in room", playerID)
}

func (r *Room) GetPlayer(playerID string) *Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, player := range r.Players {
		if player.ID == playerID {
			// Return a copy to prevent external modification
			playerCopy := player
			return &playerCopy
		}
	}
	return nil
}

func (r *Room) GetPlayers() []Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]Player, len(r.Players))
	copy(players, r.Players)
	return players
}

func (r *Room) UpdatePlayerReady(playerID string, isReady bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, player := range r.Players {
		if player.ID == playerID {
			r.Players[i].IsReady = isReady
			return nil
		}
	}

	return fmt.Errorf("player %s not found in room", playerID)
}

func (r *Room) PlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Players)
}

func (r *Room) AddBoard(board Board) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Boards = append(r.Boards, board)
}

func (r *Room) GetLatestBoard() *Board {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.Boards) == 0 {
		return nil
	}

	// Return a copy of the latest board
	latest := r.Boards[len(r.Boards)-1]
	return &latest
}

// RoomManager interface for managing rooms
type RoomManager interface {
	CreateRoom(name string) (*Room, error)
	GetRoom(roomID int) (*Room, error)
	UpdateRoom(roomID int, room *Room) error
	DeleteRoom(roomID int) error
	ListRooms() ([]*Room, error)
	JoinRoom(roomID int, player Player) (*Room, error)
	LeaveRoom(roomID int, playerID string) (*Room, error)
}
