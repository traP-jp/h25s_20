package usecase

import (
	"fmt"
	"sync"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/openapi/models"
)

type RoomUsecase struct {
	rooms map[int]*domain.Room
	mutex sync.RWMutex
}

func NewRoomUsecase() *RoomUsecase {
	usecase := &RoomUsecase{
		rooms: make(map[int]*domain.Room),
		mutex: sync.RWMutex{},
	}

	// 10個のroomを初期化
	usecase.initializeRooms()

	return usecase
}

// 10個のroomを初期化する
func (r *RoomUsecase) initializeRooms() {
	for i := 1; i <= 10; i++ {
		room := &domain.Room{
			ID:         i,
			Name:       fmt.Sprintf("Room %d", i),
			GameBoards: []domain.GameBoard{domain.NewBoard()},
			IsOpened:   i%2 == 1, // 奇数のroomはオープン、偶数はクローズ
			Players:    []domain.Player{},
			ResultLog:  []domain.Result{},
		}
		r.rooms[i] = room
	}
}

// すべてのroomを取得
func (r *RoomUsecase) GetRooms() []models.Room {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var rooms []models.Room
	for _, domainRoom := range r.rooms {
		// domain.RoomからAPI用のmodels.Roomに変換
		users := make([]string, len(domainRoom.Players))
		for i, player := range domainRoom.Players {
			users[i] = player.UserName // Playerに Name フィールドがあると仮定
		}

		apiRoom := models.Room{
			RoomId:   domainRoom.ID,
			RoomName: domainRoom.Name,
			Users:    users,
			IsOpened: domainRoom.IsOpened,
		}
		rooms = append(rooms, apiRoom)
	}

	return rooms
}

// roomIDでroomを取得
func (r *RoomUsecase) GetRoomByID(roomID int) (*domain.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	return room, nil
}

func (r *RoomUsecase) AddPlayerToRoom(roomID int, player domain.Player) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	// プレイヤーがすでに存在するかチェック
	for _, p := range room.Players {
		if p.ID == player.ID {
			return nil, fmt.Errorf("player with ID %s already exists in room %d", player.ID, roomID)
		}
	}

	room.Players = append(room.Players, player)
	return room, nil
}
