package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/auth"
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

	// 認証されたユーザー情報を取得
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	player := domain.Player{
		ID:       int(user.UserID),
		UserName: user.Username,
	}

	switch req.Action {
	case models.JOIN:
		updatedRoom, err := h.roomUsecase.AddPlayerToRoom(roomId, player)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to join room: " + err.Error(),
			})
		}

		// WebSocketでルームに参加
		if h.WebSocketHandler != nil {
			err = h.WebSocketHandler.JoinRoom(player.ID, roomId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to join WebSocket room: " + err.Error(),
				})
			}
		}

		// ルーム情報を構築
		var playerInfos []wsManager.PlayerInfo
		for _, p := range updatedRoom.Players {
			playerInfos = append(playerInfos, wsManager.ConvertToPlayerInfo(
				p.ID, p.UserName, p.IsReady, p.HasClosedResult, p.Score,
			))
		}

		roomInfo := wsManager.ConvertToRoomInfo(
			updatedRoom.ID,
			updatedRoom.Name,
			updatedRoom.State.String(),
			updatedRoom.IsOpened,
			playerInfos,
		)

		// WebSocketでルーム全員に通知（ルーム情報付き）
		if h.WebSocketHandler != nil {
			h.WebSocketHandler.SendPlayerJoinedEventToRoom(player.ID, player.UserName, roomInfo)
		}
		return c.NoContent(http.StatusNoContent)

	case models.READY:
		updatedRoom, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, player.ID, true)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update ready status: " + err.Error(),
			})
		}
		// WebSocketでルーム全員に通知
		h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventPlayerReady, player.ID, player.UserName)

		// 全員準備完了チェック
		if updatedRoom.AreAllPlayersReady() && len(updatedRoom.Players) > 0 {
			h.WebSocketHandler.SendPlayerAllReadyEventToRoom(roomId, "All players are ready!")
		}

		return c.NoContent(http.StatusNoContent)

	case models.CANCEL:
		updatedRoom, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, player.ID, false)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to cancel ready status: " + err.Error(),
			})
		}

		h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventPlayerCanceled, player.ID, player.UserName)

		// 全員準備完了状態から変更されたかチェック（必要に応じて通知）
		if !updatedRoom.AreAllPlayersReady() {
			// 必要に応じて「準備完了状態解除」イベントを送信
			// 現在は特別な処理なし
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
		if firstPlayer == nil || firstPlayer.ID != player.ID {
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
		// StartGame()が成功した時点でStateがCountdownに変わっているため、
		// 重複実行は防止されている
		go h.handleGameStart(roomId)
		return c.NoContent(http.StatusNoContent)

	case models.ABORT:
		// WebSocketからルームを退出
		if h.WebSocketHandler != nil {
			err := h.WebSocketHandler.LeaveRoom(player.ID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to leave WebSocket room: " + err.Error(),
				})
			}
		}

		// プレイヤーをルームから削除
		updatedRoom, err := h.roomUsecase.RemovePlayerFromRoom(roomId, player.ID)
		if err != nil {
			// プレイヤーが見つからない場合でもエラーにしない（既に退出済みの可能性）
		}

		// ルーム情報を構築（退出後の状態）
		if updatedRoom != nil && h.WebSocketHandler != nil {
			var playerInfos []wsManager.PlayerInfo
			for _, p := range updatedRoom.Players {
				playerInfos = append(playerInfos, wsManager.ConvertToPlayerInfo(
					p.ID, p.UserName, p.IsReady, p.HasClosedResult, p.Score,
				))
			}

			roomInfo := wsManager.ConvertToRoomInfo(
				updatedRoom.ID,
				updatedRoom.Name,
				updatedRoom.State.String(),
				updatedRoom.IsOpened,
				playerInfos,
			)

			// WebSocketでルーム全員に通知（ルーム情報付き）
			h.WebSocketHandler.SendPlayerLeftEventToRoom(player.ID, player.UserName, roomInfo)
		}

		return c.NoContent(http.StatusNoContent)

	case models.CLOSERESULT:
		_, err := h.roomUsecase.CloseResult(roomId, player.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to close result: " + err.Error(),
			})
		}

		// WebSocketでルーム全員に通知
		h.WebSocketHandler.SendPlayerEventToRoom(roomId, wsManager.EventResultClosed, player.ID, player.UserName)
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

	// 認証されたユーザー情報を取得
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	player := domain.Player{
		ID:       int(user.UserID),
		UserName: user.Username,
	}

	room, err := h.roomUsecase.GetRoomByID(roomId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Room not found",
		})
	}

	// プレイヤーがルームに参加しているかチェック
	playerInRoom := false
	for _, roomPlayer := range room.Players {
		if roomPlayer.ID == player.ID {
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
	board, gainScore, err := h.roomUsecase.ApplyFormulaWithVersion(roomId, player.ID, req.Formula, req.Version)
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
		h.WebSocketHandler.SendBoardUpdateEventTyped(roomId, player.ID, player.UserName, boardData, gainScore)
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
	// 認証されたユーザー情報を取得
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	player := domain.Player{
		ID:       int(user.UserID),
		UserName: user.Username,
	}

	room, err := h.roomUsecase.GetRoomByID(roomId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Room not found",
		})
	}

	// プレイヤーがルームに参加しているかチェック
	playerInRoom := false
	for _, roomPlayer := range room.Players {
		if roomPlayer.ID == player.ID {
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
	for _, roomPlayer := range room.Players {
		results = append(results, models.RoomResultItem{
			User:  roomPlayer.UserName,
			Score: roomPlayer.Score,
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
		if h.WebSocketHandler != nil {
			h.WebSocketHandler.SendCountdownEventToRoom(roomID, i)
		}
		time.Sleep(1 * time.Second)
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

	// 120秒タイマーを開始（別goroutineで実行）
	go h.handleGameTimer(roomID)
}

// handleGameTimer は120秒のゲームタイマーとラスト10秒のカウントダウンを処理する
func (h *Handler) handleGameTimer(roomID int) {
	// 110秒待機（120秒 - 10秒のカウントダウン）
	time.Sleep(110 * time.Second)

	// ゲームがまだ進行中かチェック
	room, err := h.roomUsecase.GetRoomByID(roomID)
	if err != nil || room.State != domain.StateGameInProgress {
		return // ゲームが既に終了している場合は何もしない
	}

	// ラスト10秒のカウントダウン開始を通知
	if h.WebSocketHandler != nil {
		h.WebSocketHandler.SendCountdownStartEventToRoom(roomID, "Game ending in 10 seconds", 10)
	}

	// 10秒のカウントダウン
	for i := 10; i > 0; i-- {

		// 各秒でゲームがまだ進行中かチェック
		room, err := h.roomUsecase.GetRoomByID(roomID)
		if err != nil || room.State != domain.StateGameInProgress {
			return // ゲームが既に終了している場合は中断
		}

		if h.WebSocketHandler != nil {
			h.WebSocketHandler.SendCountdownEventToRoom(roomID, i)
		}

		time.Sleep(1 * time.Second)
	}

	// タイマー終了、ゲームを終了する
	_, err = h.roomUsecase.EndGame(roomID)
	if err != nil {
		return
	}

	// ゲーム終了をWebSocketで通知
	if h.WebSocketHandler != nil {
		h.WebSocketHandler.SendGameEndEventToRoom(roomID, "Time's up! Game ended.")
	}
}
