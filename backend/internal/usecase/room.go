package usecase

import (
	"fmt"
	"sync"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/openapi/models"
	"github.com/rs/zerolog/log"
)

type RoomUsecase struct {
	rooms      map[int]*domain.Room
	mutex      sync.RWMutex
	gameTimers map[int]bool // ゲームタイマー重複実行防止用
	timerMutex sync.Mutex   // gameTimers用の専用mutex
}

func NewRoomUsecase() *RoomUsecase {
	usecase := &RoomUsecase{
		rooms:      make(map[int]*domain.Room),
		mutex:      sync.RWMutex{},
		gameTimers: make(map[int]bool),
		timerMutex: sync.Mutex{},
	}

	// 10個のroomを初期化
	usecase.initializeRooms()

	return usecase
}

// 10個のroomを初期化する
func (r *RoomUsecase) initializeRooms() {
	for i := 1; i <= 10; i++ {
		room := domain.NewRoom(i, fmt.Sprintf("Room %d", i))
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

	// 新しいプレイヤーの接続状態を初期化
	player.IsConnected = true
	player.LastSeenAt = nil

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
		// 全員がreadyになったらroomをクローズ
		room.IsOpened = false
	} else if !room.AreAllPlayersReady() && room.State == domain.StateAllReady {
		room.TransitionTo(domain.StateWaitingForPlayers)
		// 状態がWaitingForPlayersに戻った場合、部屋を再度開放
		room.IsOpened = true
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

	// プレイヤーを見つけて削除（安全性向上）
	playerFound := false
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].ID == playerID {
			// 境界チェック付きでスライスから削除
			if i < len(room.Players)-1 {
				room.Players = append(room.Players[:i], room.Players[i+1:]...)
			} else {
				room.Players = room.Players[:i]
			}
			playerFound = true
			break
		}
	}

	if !playerFound {
		return nil, fmt.Errorf("player with ID %d not found in room %d", playerID, roomID)
	}

	// 参加者が0人になった場合、ルームをリセット
	if len(room.Players) == 0 {
		// 状態に関わらず、ルームを完全に初期状態に戻す
		room.State = domain.StateWaitingForPlayers
		room.GameBoards = []domain.GameBoard{domain.NewBoard()}
		room.ResultLog = []domain.Result{}
		room.IsOpened = true
		room.LastCorrectPlayerID = 0
		room.StreakCount = 0
		// プレイヤーリストは既に空なので、個々のリセットは不要
	} else {
		// まだプレイヤーがいる場合、READY状態をチェック
		if !room.AreAllPlayersReady() && room.State == domain.StateAllReady {
			room.TransitionTo(domain.StateWaitingForPlayers)
			// 状態がWaitingForPlayersに戻った場合、部屋を再度開放
			room.IsOpened = true
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

	// ゲーム終了時にタイマーを確実に停止
	r.StopGameTimer(roomID)

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

	// 連続正解数とスコア計算を原子的に実行
	var gainScore int
	playerFound := false

	// プレイヤーを見つけてスコア更新を1回だけ実行
	for i := range room.Players {
		if room.Players[i].ID == playerID {
			// 連続正解数をカウント
			if room.LastCorrectPlayerID == playerID {
				room.StreakCount++
			} else {
				room.StreakCount = 1
				room.LastCorrectPlayerID = playerID
			}

			// スコア計算: 消した組数 * (5+5*"連続正解数")点
			gainScore = matchCount * (5 + 5*room.StreakCount)

			// プレイヤーのスコアに加算（1回のみ保証）
			room.Players[i].Score += gainScore
			playerFound = true
			break
		}
	}

	if !playerFound {
		return nil, 0, fmt.Errorf("player with ID %d not found in room", playerID)
	}

	// データレース回避のためGameBoardのディープコピーを返す
	safeBoard := &domain.GameBoard{
		Version:       currentBoard.Version,
		Size:          currentBoard.Size,
		Board:         make([][]int, currentBoard.Size),
		ChangeHistory: make(map[int][]domain.Matches),
	}

	// 盤面データをコピー
	for i := 0; i < currentBoard.Size; i++ {
		safeBoard.Board[i] = make([]int, currentBoard.Size)
		copy(safeBoard.Board[i], currentBoard.Board[i])
	}

	// 変更履歴もコピー（必要に応じて）
	for k, v := range currentBoard.ChangeHistory {
		safeBoard.ChangeHistory[k] = v
	}

	return safeBoard, gainScore, nil
}

// SetPlayerDisconnected marks a player as disconnected but keeps them in the room
func (r *RoomUsecase) SetPlayerDisconnected(roomID int, playerID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	// プレイヤーを見つけて切断状態に設定
	playerFound := false
	now := time.Now()
	for i, player := range room.Players {
		if player.ID == playerID {
			room.Players[i].IsConnected = false
			room.Players[i].LastSeenAt = &now
			playerFound = true

			log.Info().
				Int("room_id", roomID).
				Int("player_id", playerID).
				Str("player_name", player.UserName).
				Msg("Player marked as disconnected in room")
			break
		}
	}

	if !playerFound {
		return nil, fmt.Errorf("player with ID %d not found in room %d", playerID, roomID)
	}

	return room, nil
}

// SetPlayerReconnected marks a player as reconnected
func (r *RoomUsecase) SetPlayerReconnected(roomID int, playerID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	// プレイヤーを見つけて再接続状態に設定
	playerFound := false
	for i, player := range room.Players {
		if player.ID == playerID {
			room.Players[i].IsConnected = true
			room.Players[i].LastSeenAt = nil
			playerFound = true

			log.Info().
				Int("room_id", roomID).
				Int("player_id", playerID).
				Str("player_name", player.UserName).
				Msg("Player reconnected to room")
			break
		}
	}

	if !playerFound {
		return nil, fmt.Errorf("player with ID %d not found in room %d", playerID, roomID)
	}

	return room, nil
}

// RemoveDisconnectedPlayer removes a player who has been disconnected for too long
func (r *RoomUsecase) RemoveDisconnectedPlayer(roomID int, playerID int) (*domain.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %d not found", roomID)
	}

	// プレイヤーを見つけて削除（安全性向上）
	playerFound := false
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].ID == playerID && !room.Players[i].IsConnected {
			playerName := room.Players[i].UserName
			// 境界チェック付きでスライスから削除
			if i < len(room.Players)-1 {
				room.Players = append(room.Players[:i], room.Players[i+1:]...)
			} else {
				room.Players = room.Players[:i]
			}
			playerFound = true

			log.Info().
				Int("room_id", roomID).
				Int("player_id", playerID).
				Str("player_name", playerName).
				Msg("Disconnected player permanently removed from room")
			break
		}
	}

	if !playerFound {
		return nil, fmt.Errorf("disconnected player with ID %d not found in room %d", playerID, roomID)
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
		room.IsOpened = true // ルームを再度開放
	} else {
		// まだプレイヤーがいる場合、READY状態をチェック（接続中のプレイヤーのみ）
		if !r.areConnectedPlayersReady(room) && room.State == domain.StateAllReady {
			room.TransitionTo(domain.StateWaitingForPlayers)
			// 状態がWaitingForPlayersに戻った場合、部屋を再度開放
			room.IsOpened = true
		}
	}

	return room, nil
}

// GetDisconnectedPlayers returns all disconnected players in all rooms
func (r *RoomUsecase) GetDisconnectedPlayers() map[int][]domain.Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	disconnectedPlayers := make(map[int][]domain.Player)

	for roomID, room := range r.rooms {
		var roomDisconnected []domain.Player
		for _, player := range room.Players {
			if !player.IsConnected {
				roomDisconnected = append(roomDisconnected, player)
			}
		}
		if len(roomDisconnected) > 0 {
			disconnectedPlayers[roomID] = roomDisconnected
		}
	}

	return disconnectedPlayers
}

// areConnectedPlayersReady checks if all connected players are ready
func (r *RoomUsecase) areConnectedPlayersReady(room *domain.Room) bool {
	if len(room.Players) == 0 {
		return false
	}

	connectedCount := 0
	readyCount := 0

	for _, player := range room.Players {
		if player.IsConnected {
			connectedCount++
			if player.IsReady {
				readyCount++
			}
		}
	}

	// 接続中のプレイヤーが1人以上いて、全員がREADY
	return connectedCount > 0 && connectedCount == readyCount
}

// CanStartGameTimer ゲームタイマーが開始可能かチェック（重複実行防止）
func (r *RoomUsecase) CanStartGameTimer(roomID int) bool {
	r.timerMutex.Lock()
	defer r.timerMutex.Unlock()

	if r.gameTimers[roomID] {
		return false // 既にタイマーが実行中
	}
	r.gameTimers[roomID] = true
	return true
}

// StopGameTimer ゲームタイマーを停止状態にマーク
func (r *RoomUsecase) StopGameTimer(roomID int) {
	r.timerMutex.Lock()
	defer r.timerMutex.Unlock()
	delete(r.gameTimers, roomID)
}
