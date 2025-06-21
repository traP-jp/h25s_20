package usecase

import (
	"errors"
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
		room := domain.NewRoom(i, fmt.Sprintf("Room %d", i))
		if i%2 == 0 { // 偶数のroomはクローズ
			room.IsOpened = false
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
			return nil, fmt.Errorf("player with ID %d already exists in room %d", player.ID, roomID)
		}
	}

	room.Players = append(room.Players, player)
	return room, nil
}

func (r *RoomUsecase) UpdatePlayerReadyStatus(roomID int, playerID int, isReady bool) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	// プレイヤーを見つけて更新
	playerFound := false
	for i, player := range room.Players {
		if player.ID == playerID {
			room.Players[i].IsReady = isReady
			playerFound = true
			break
		}
	}

	if !playerFound {
		return nil, fmt.Errorf("player with ID %d not found in room %d", playerID, roomID)
	}

	// 全員がREADYになったら状態を更新
	if room.AreAllPlayersReady() && room.State == domain.StateWaitingForPlayers {
		room.TransitionTo(domain.StateAllReady)
	} else if !room.AreAllPlayersReady() && room.State == domain.StateAllReady {
		room.TransitionTo(domain.StateWaitingForPlayers)
	}

	return room, nil
}

// StartGame starts the game for the specified room
func (r *RoomUsecase) StartGame(roomID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	err := room.StartGame()
	if err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}

	return room, nil
}

// AbortGame aborts the game for the specified room
func (r *RoomUsecase) AbortGame(roomID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	err := room.AbortGame()
	if err != nil {
		return nil, fmt.Errorf("failed to abort game: %w", err)
	}

	return room, nil
}

// CloseResult closes the result display for a player in the specified room
func (r *RoomUsecase) CloseResult(roomID int, playerID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	err := room.CloseResult(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to close result: %w", err)
	}

	return room, nil
}

// 数式を受け取りスコアとボードの更新を行う
func (r *RoomUsecase) ApplyFormula(roomID int, playerID int, formula string) (*domain.GameBoard, int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	//ルームが存在するかをチェック
	room, exists := r.rooms[roomID]
	if !exists {
		return nil, 0, fmt.Errorf("room with ID %d not found", roomID)
	}

	// プレイヤーが参加しているかをチェック
	playerInRoom := false
	for _, p := range room.Players {
		if p.ID == playerID {
			playerInRoom = true
			break
		}
	}
	if !playerInRoom {
		return nil, 0, fmt.Errorf("player is not in this room")
	}

	// ゲーム進行中かをチェック
	if room.State != domain.StateGameInProgress {
		return nil, 0, fmt.Errorf("game is not in progress")
	}

	if len(room.GameBoards) == 0 {
		return nil, 0, fmt.Errorf("no game board available")
	}
	currentBoard := &room.GameBoards[len(room.GameBoards)-1]

	// AttemptMoveを呼び出し、成否を受け取る
	success, message := domain.AttemptMove(currentBoard, formula)

	// もし、AttemptMoveの結果が失敗だったら元のボードを送信する
	if !success {
		return currentBoard, 0, errors.New(message)
	}

	// --- 成功した場合のみ、以下の処理を行う ---
	gainScore := 10
	for i := range room.Players {
		if room.Players[i].ID == playerID {
			room.Players[i].Score += gainScore
			break
		}
	}

	return currentBoard, gainScore, nil
}
