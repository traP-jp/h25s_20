package handler

import (
	"net/http"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/openapi/models"
	"github.com/labstack/echo/v4"
)

// GetRooms returns a list of all rooms
func (h *Handler) GetRooms(c echo.Context) error {
	rooms := h.roomUsecase.GetRooms()
	return c.JSON(http.StatusOK, rooms)
}

// PostRoomsRoomIdActions performs an action on a specific room
func (h *Handler) PostRoomsRoomIdActions(c echo.Context, roomId int) error {
	var req models.PostRoomsRoomIdActionsJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// TODO: 実際のユーザー認証を実装した後に、ユーザー情報を取得する
	// 現在はテスト用のプレイヤーを使用
	var mockPlayer = domain.Player{
		ID:       1,
		UserName: "testuser",
	}

	switch req.Action {
	case models.JOIN:
		_, err := h.roomUsecase.AddPlayerToRoom(roomId, mockPlayer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to join room: " + err.Error(),
			})
		}
		return c.NoContent(http.StatusNoContent)

	case models.READY:
		_, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, mockPlayer.ID, true)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update ready status: " + err.Error(),
			})
		}
		return c.NoContent(http.StatusNoContent)

	case models.CANCEL:
		_, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, mockPlayer.ID, false)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to cancel ready status: " + err.Error(),
			})
		}
		return c.NoContent(http.StatusNoContent)

	case models.START:
		room, err := h.roomUsecase.GetRoomByID(roomId)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Room not found",
			})
		}

		// ゲーム開始権限チェック（最初のプレイヤーのみが開始可能）
		firstPlayer := room.GetFirstPlayer()
		if firstPlayer == nil || firstPlayer.ID != mockPlayer.ID {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Only the first player can start the game",
			})
		}

		_, err = h.roomUsecase.StartGame(roomId)
		if err != nil {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Cannot start game: " + err.Error(),
			})
		}

		// カウントダウンと最初のボード生成を別のgoroutineで実行
		go h.handleGameStart(roomId)

		return c.NoContent(http.StatusNoContent)

	case models.ABORT:
		_, err := h.roomUsecase.AbortGame(roomId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to abort game: " + err.Error(),
			})
		}
		return c.NoContent(http.StatusNoContent)

	case models.CLOSERESULT:
		_, err := h.roomUsecase.CloseResult(roomId, mockPlayer.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to close result: " + err.Error(),
			})
		}
		return c.NoContent(http.StatusNoContent)

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid action",
		})
	}
}

// PostRoomsRoomIdFormulas submits a formula for calculation
func (h *Handler) PostRoomsRoomIdFormulas(c echo.Context, roomId int) error {
	var req models.PostRoomsRoomIdFormulasJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// TODO: 実際のユーザー認証を実装した後に、ユーザー情報を取得する
	mockPlayer := domain.Player{
		ID:       1,
		UserName: "testuser",
	}

	room, err := h.roomUsecase.GetRoomByID(roomId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Room not found",
		})
	}

	// プレイヤーがルームに参加しているかチェック
	playerInRoom := false
	for _, player := range room.Players {
		if player.ID == mockPlayer.ID {
			playerInRoom = true
			break
		}
	}

	if !playerInRoom {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Player is not in this room",
		})
	}

	// ゲームが進行中かチェック
	if room.State != domain.StateGameInProgress {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Game is not in progress",
		})
	}

	// 現在のゲームボードを取得
	if len(room.GameBoards) == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "No game board available",
		})
	}

	currentBoard := &room.GameBoards[len(room.GameBoards)-1]

	// 数式を検証し、盤面を更新
	domain.AttemptMove(currentBoard, req.Formula)

	// 盤面データを1次元配列に変換
	content := make([]int, 0, currentBoard.Size*currentBoard.Size)
	for i := 0; i < currentBoard.Size; i++ {
		for j := 0; j < currentBoard.Size; j++ {
			content = append(content, currentBoard.Board[i][j])
		}
	}

	// TODO: 実際のスコア計算ロジックを実装
	gainScore := 10

	response := models.Board{
		Content:   content,
		Version:   currentBoard.Version,
		GainScore: gainScore,
	}

	return c.JSON(http.StatusOK, response)
}

// GetRoomsRoomIdResult returns the results of a specific room
func (h *Handler) GetRoomsRoomIdResult(c echo.Context, roomId int) error {
	// TODO: 実際のユーザー認証を実装した後に、ユーザー情報を取得する
	mockPlayer := domain.Player{
		ID:       1,
		UserName: "testuser",
	}

	room, err := h.roomUsecase.GetRoomByID(roomId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Room not found",
		})
	}

	// プレイヤーがルームに参加しているかチェック
	playerInRoom := false
	for _, player := range room.Players {
		if player.ID == mockPlayer.ID {
			playerInRoom = true
			break
		}
	}

	if !playerInRoom {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Player is not in this room",
		})
	}

	// 結果を構築
	var results []models.RoomResultItem
	for _, player := range room.Players {
		results = append(results, models.RoomResultItem{
			User:  player.UserName,
			Score: player.Score,
		})
	}

	return c.JSON(http.StatusOK, results)
}

// handleGameStart はゲーム開始時のカウントダウンと最初のボード生成・送信を処理する
func (h *Handler) handleGameStart(roomID int) {
	// カウントダウンを開始する通知を送信
	if h.notificationService != nil {
		h.notificationService.NotifyRoom(roomID, "countdown_start", map[string]interface{}{
			"message":   "Game starting in 3 seconds",
			"countdown": 3,
		})
	}

	// 3秒のカウントダウン
	for i := 3; i > 0; i-- {
		time.Sleep(1 * time.Second)
		if h.notificationService != nil {
			h.notificationService.NotifyRoom(roomID, "countdown", map[string]interface{}{
				"count": i,
			})
		}
	}

	// カウントダウン完了後、ゲームを実際に開始
	_, err := h.roomUsecase.CompleteCountdown(roomID)
	if err != nil {
		return
	}

	// 新しいボードを生成
	newBoard := domain.NewBoard()

	// ボードをroomに追加
	_, err = h.roomUsecase.UpdateGameBoard(roomID, newBoard)
	if err != nil {
		return
	}

	// ボードデータを1次元配列に変換
	content := make([]int, 0, newBoard.Size*newBoard.Size)
	for i := 0; i < newBoard.Size; i++ {
		for j := 0; j < newBoard.Size; j++ {
			content = append(content, newBoard.Board[i][j])
		}
	}

	// ゲーム開始とボード情報を送信
	if h.notificationService != nil {
		h.notificationService.NotifyRoom(roomID, "game_start", map[string]interface{}{
			"message": "Game started!",
			"board": map[string]interface{}{
				"content": content,
				"version": newBoard.Version,
				"size":    newBoard.Size,
			},
		})
	}
}
