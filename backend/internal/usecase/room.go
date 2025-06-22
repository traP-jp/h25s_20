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
		users := make([]models.User, len(domainRoom.Players))
		for i, player := range domainRoom.Players {
			users[i] = models.User{
				Username: player.UserName,
				IsReady:  player.IsReady,
			}
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

// CompleteCountdown completes the countdown and transitions to game in progress
func (r *RoomUsecase) CompleteCountdown(roomID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	err := room.CompleteCountdown()
	if err != nil {
		return nil, fmt.Errorf("failed to complete countdown: %w", err)
	}

	return room, nil
}

// UpdateGameBoard updates the game board for the specified room
func (r *RoomUsecase) UpdateGameBoard(roomID int, newBoard domain.GameBoard) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	room.GameBoards = append(room.GameBoards, newBoard)
	return room, nil
}

// RemovePlayerFromRoom removes a player from the specified room
func (r *RoomUsecase) RemovePlayerFromRoom(roomID int, playerID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	// プレイヤーを見つけて削除
	playerFound := false
	for i, player := range room.Players {
		if player.ID == playerID {
			// スライスから削除
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			playerFound = true
			break
		}
	}

	if !playerFound {
		return nil, fmt.Errorf("player with ID %d not found in room %d", playerID, roomID)
	}

	// 参加者が0人になった場合、ルームをリセット
	if len(room.Players) == 0 {
		err := room.ResetRoom()
		if err != nil {
			// StateGameEndedでない場合は強制的にWaitingForPlayersに戻す
			room.State = domain.StateWaitingForPlayers
		}
		// ゲームボードもリセット
		room.GameBoards = []domain.GameBoard{domain.NewBoard()}
		room.ResultLog = []domain.Result{}
	} else {
		// まだプレイヤーがいる場合、READY状態をチェック
		if !room.AreAllPlayersReady() && room.State == domain.StateAllReady {
			room.TransitionTo(domain.StateWaitingForPlayers)
		}
	}

	return room, nil
}

// EndGame ends the game for the specified room
func (r *RoomUsecase) EndGame(roomID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	err := room.EndGame()
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}

	return room, nil
}

// ApplyFormulaWithVersion はバージョンを考慮した細かい衝突検出付きの数式適用
func (r *RoomUsecase) ApplyFormulaWithVersion(roomID int, playerID int, formula string, submittedVersion int) (*domain.GameBoard, int, error) {
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

	// バージョン付きの細かい衝突検出を実行
	success, errMessage, matchCount := domain.AttemptMoveWithVersion(currentBoard, formula, submittedVersion)

	if !success {
		return nil, 0, fmt.Errorf("%s", errMessage)
	}

	// 連続正解数をカウント
	if room.LastCorrectPlayerID == playerID {
		room.StreakCount++
	} else {
		room.StreakCount = 1
		room.LastCorrectPlayerID = playerID
	}

	// スコア計算: 消した組数 * (5+5*"連続正解数")点
	gainScore := matchCount * (5 + 5*room.StreakCount)

	// プレイヤーのスコアに加算
	for i := range room.Players {
		if room.Players[i].ID == playerID {
			room.Players[i].Score += gainScore
			break
		}
	}

	return currentBoard, gainScore, nil
}
