package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	wsManager "github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
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

		// WebSocketでルームに参加
		if h.WebSocketHandler != nil {
			err = h.WebSocketHandler.JoinRoom(mockPlayer.ID, roomId)
			if err != nil {
				// WebSocketのjoinに失敗してもHTTPレスポンスはエラーにしない（ログ出力のみ）
				// WebSocket接続がない場合もあるため
			}
		}

		// WebSocketでルーム全員に通知
		if h.WebSocketHandler != nil {
			h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventPlayerJoined, mockPlayer.ID, mockPlayer.UserName)
		}
		return c.NoContent(http.StatusNoContent)

	case models.READY:
		_, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, mockPlayer.ID, true)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update ready status: " + err.Error(),
			})
		}
		// WebSocketでルーム全員に通知
		h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventPlayerReady, mockPlayer.ID, mockPlayer.UserName)
		return c.NoContent(http.StatusNoContent)

	case models.CANCEL:
		_, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, mockPlayer.ID, false)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to cancel ready status: " + err.Error(),
			})
		}

		h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventPlayerCanceled, mockPlayer.ID, mockPlayer.UserName)
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
		// WebSocketでルーム全員にゲーム開始を通知
		if h.WebSocketHandler != nil {
			h.WebSocketHandler.SendGameStartEventToRoom(roomId, "Game has started")
		}

		// カウントダウンと最初のボード生成を別のgoroutineで実行
		go h.handleGameStart(roomId)
		return c.NoContent(http.StatusNoContent)

	case models.ABORT:
		// WebSocketからルームを退出
		if h.WebSocketHandler != nil {
			err := h.WebSocketHandler.LeaveRoom(mockPlayer.ID)
			if err != nil {
				// WebSocketのleaveに失敗してもHTTPレスポンスはエラーにしない（ログ出力のみ）
			}
		}

		// プレイヤーをルームから削除
		_, err := h.roomUsecase.RemovePlayerFromRoom(roomId, mockPlayer.ID)
		if err != nil {
			// プレイヤーが見つからない場合でもエラーにしない（既に退出済みの可能性）
		}

		// WebSocketでの通知
		if h.WebSocketHandler != nil {
			h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventPlayerLeft, mockPlayer.ID, mockPlayer.UserName)
		}

		return c.NoContent(http.StatusNoContent)

	case models.CLOSERESULT:
		_, err := h.roomUsecase.CloseResult(roomId, mockPlayer.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to close result: " + err.Error(),
			})
		}

		// WebSocketでルーム全員に通知
		h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventResultClosed, mockPlayer.ID, mockPlayer.UserName)
		return c.NoContent(http.StatusNoContent)

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid action",
		})
	}
}

// PostRoomsRoomIdFormulas submits a formula for calculation
func (h *Handler) PostRoomsRoomIdFormulas(c echo.Context, roomId int) error {
	var req models.PostRoomsRoomIdFormulasJSONBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
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

	// バージョン付きの細かい衝突検出を使用
	board, gainScore, err := h.roomUsecase.ApplyFormulaWithVersion(roomId, mockPlayer.ID, req.Formula, req.Version)
	if err != nil {
		// 衝突エラーの場合は409を返す
		if strings.Contains(err.Error(), "他のプレイヤーによって更新されています") ||
			strings.Contains(err.Error(), "無効なバージョンです") {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		// その他のエラーは400を返す
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// 盤面データを1次元配列に変換
	content := make([]int, 0, board.Size*board.Size)
	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			content = append(content, board.Board[i][j])
		}
	}

	// 成功時はWebSocketでルーム全体に盤面更新を通知
	if h.WebSocketHandler != nil {
		boardData := wsManager.BoardData{
			Content: content,
			Version: board.Version,
			Size:    board.Size,
		}
		h.WebSocketHandler.SendBoardUpdateEventTyped(roomId, mockPlayer.ID, mockPlayer.UserName, boardData, gainScore)
	}

	// HTTPレスポンス（提出者に対する結果）
	return c.JSON(http.StatusOK, models.Board{
		Content:   content,
		Version:   board.Version,
		GainScore: gainScore,
	})
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
	if h.WebSocketHandler != nil {
		h.WebSocketHandler.SendCountdownStartEventToRoom(roomID, "Game starting in 3 seconds", 3)
	}

	// 3秒のカウントダウン
	for i := 3; i > 0; i-- {
		time.Sleep(1 * time.Second)
		if h.WebSocketHandler != nil {
			h.WebSocketHandler.SendCountdownEventToRoom(roomID, i)
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
	if h.WebSocketHandler != nil {
		boardData := wsManager.BoardData{
			Content: content,
			Version: newBoard.Version,
			Size:    newBoard.Size,
		}
		h.WebSocketHandler.SendGameStartBoardEventToRoom(roomID, "Game started!", boardData)
	}
}
