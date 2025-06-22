package handler

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	wsManager "github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WebSocketHandler struct {
	manager     *wsManager.Manager
	roomUsecase *usecase.RoomUsecase
	userUsecase *usecase.UserUsecase
}

func NewWebSocketHandler(manager *wsManager.Manager, roomUsecase *usecase.RoomUsecase, userUsecase *usecase.UserUsecase) *WebSocketHandler {
	return &WebSocketHandler{
		manager:     manager,
		roomUsecase: roomUsecase,
		userUsecase: userUsecase,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// ユーザー名をクエリパラメータから取得
	username := c.QueryParam("username")
	if username == "" {
		log.Error().Msg("Username is required for WebSocket connection")
		// WebSocketアップグレード前なので、まだJSONレスポンスが可能
		return c.JSON(400, map[string]string{"error": "username is required"})
	}

	// データベースからユーザーを検索
	user, err := h.userUsecase.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("User not found for WebSocket connection")
		// WebSocketアップグレード前なので、まだJSONレスポンスが可能
		return c.JSON(404, map[string]string{"error": "user not found"})
	}

	userID := int(user.ID)

	// WebSocket接続をアップグレード（CORS対応のオプション追加）
	conn, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{
		Subprotocols:   []string{"echo"},
		OriginPatterns: []string{"http://localhost:5173"},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade WebSocket connection")
		// WebSocketアップグレード後はJSONレスポンスできないため、errorを返すのみ
		return err
	}

	// クライアントIDを生成
	clientID := uuid.New().String()

	// コンテキストとキャンセル関数を作成
	ctx, cancel := context.WithCancel(context.Background())

	// クライアントをマネージャーに登録
	h.manager.AddClient(clientID, userID, conn, cancel)

	// 接続完了メッセージを送信
	if err := h.SendConnectionEvent(userID, clientID, "Connected successfully", time.Now().Unix()); err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to send welcome message")
	}

	// 接続を維持し、切断を監視
	go h.handleConnection(ctx, conn, clientID)

	return nil
}

func (h *WebSocketHandler) handleConnection(ctx context.Context, conn *websocket.Conn, clientID string) {
	defer func() {
		h.manager.RemoveClient(clientID)
		conn.Close(websocket.StatusNormalClosure, "Connection closed")
	}()

	// 一方向通信なので、クライアントからのメッセージは基本的に受信しない
	// ただし、接続の生存確認のためにpingを監視
	for {
		select {
		case <-ctx.Done():
			log.Info().Str("client_id", clientID).Msg("WebSocket connection context cancelled")
			return
		default:
			// タイムアウト付きでメッセージを読み取り（主にpingフレーム用）
			ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
			_, _, err := conn.Read(ctx)
			cancel()

			if err != nil {
				// 通常の切断またはタイムアウトの場合
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway {
					log.Info().Str("client_id", clientID).Msg("WebSocket connection closed normally")
				} else {
					log.Warn().Err(err).Str("client_id", clientID).Msg("WebSocket connection error")
				}
				return
			}
		}
	}
}

// 全クライアントにメッセージを送信
func (h *WebSocketHandler) BroadcastToAll(event string, content interface{}) {
	h.manager.NotifyAll(event, content)
}

// room参加者全員にメッセージを送信
func (h *WebSocketHandler) BroadcastToRoom(roomID int, event string, content interface{}) {
	h.manager.NotifyRoom(roomID, event, content)
}

// room未参加者全員にメッセージを送信
func (h *WebSocketHandler) BroadcastToNonRoomMembers(event string, content interface{}) {
	h.manager.NotifyNonRoomMembers(event, content)
}

// 特定のユーザーにメッセージを送信
func (h *WebSocketHandler) SendToUser(userID int, event string, content interface{}) error {
	return h.manager.NotifyUser(userID, event, content)
}

// ユーザーをroomに参加させる
func (h *WebSocketHandler) JoinRoom(userID, roomID int) error {
	return h.manager.JoinRoom(userID, roomID)
}

// ユーザーをroomから退出させる
func (h *WebSocketHandler) LeaveRoom(userID int) error {
	return h.manager.LeaveRoom(userID)
}

func (h *WebSocketHandler) GetConnectedClients() int {
	return h.manager.GetClientCount()
}

func (h *WebSocketHandler) GetRoomConnectedClients(roomID int) int {
	return h.manager.GetRoomClientCount(roomID)
}

// WebSocketイベント送信のヘルパーメソッド群（HTTP専用方針対応）

// SendStandardEvent sends a standardized event to the room
func (h *WebSocketHandler) SendStandardEvent(roomID int, event string, userID int, userName string, message string, data interface{}) {
	content := wsManager.StandardEventContent{
		UserID:   userID,
		UserName: userName,
		RoomID:   roomID,
		Message:  message,
		Data:     data,
	}
	h.BroadcastToRoom(roomID, event, content)
}

// SendBoardUpdateEvent sends a board update event to the room
// 非推奨: SendBoardUpdateEventTypedを使用してください
func (h *WebSocketHandler) SendBoardUpdateEvent(roomID int, userID int, userName string, board interface{}, gainScore int) {
	content := wsManager.BoardUpdateContent{
		StandardEventContent: wsManager.StandardEventContent{
			UserID:   userID,
			UserName: userName,
			RoomID:   roomID,
			Message:  "Board updated",
		},
		Board:     board,
		GainScore: gainScore,
	}
	h.BroadcastToRoom(roomID, "board_updated", content)
}

// 新しい統一されたイベント送信メソッド群

// SendConnectionEvent sends a connection event to a user
func (h *WebSocketHandler) SendConnectionEvent(userID int, clientID string, message string, timestamp int64) error {
	event := wsManager.NewConnectionEvent(clientID, userID, message, timestamp)
	return h.manager.SendEventToUser(userID, event)
}

// SendPlayerEventToRoom sends a player event to all room members
func (h *WebSocketHandler) SendPlayerEventToRoom(roomID int, eventType string, userID int, userName string) {
	event := wsManager.NewPlayerEvent(eventType, userID, userName, roomID)
	h.manager.SendEventToRoom(roomID, event)
}

// SendGameStartEventToRoom sends a game start event to all room members
func (h *WebSocketHandler) SendGameStartEventToRoom(roomID int, message string) {
	event := wsManager.NewGameStartEvent(roomID, message)
	h.manager.SendEventToRoom(roomID, event)
}

// SendCountdownStartEventToRoom sends a countdown start event to all room members
func (h *WebSocketHandler) SendCountdownStartEventToRoom(roomID int, message string, countdown int) {
	event := wsManager.NewCountdownStartEvent(roomID, message, countdown)
	h.manager.SendEventToRoom(roomID, event)
}

// SendCountdownEventToRoom sends a countdown event to all room members
func (h *WebSocketHandler) SendCountdownEventToRoom(roomID int, count int) {
	event := wsManager.NewCountdownEvent(roomID, count)
	h.manager.SendEventToRoom(roomID, event)
}

// SendBoardUpdateEventTyped sends a typed board update event to all room members
func (h *WebSocketHandler) SendBoardUpdateEventTyped(roomID int, userID int, userName string, board wsManager.BoardData, gainScore int) {
	event := wsManager.NewBoardUpdateEvent(userID, userName, roomID, board, gainScore)
	h.manager.SendEventToRoom(roomID, event)
}

// SendGameStartBoardEventToRoom sends a game start event with board data to all room members
func (h *WebSocketHandler) SendGameStartBoardEventToRoom(roomID int, message string, board wsManager.BoardData) {
	event := wsManager.NewGameStartBoardEvent(roomID, message, board)
	h.manager.SendEventToRoom(roomID, event)
}

// SendGameEndEventToRoom sends a game end event to all room members
func (h *WebSocketHandler) SendGameEndEventToRoom(roomID int, message string) {
	event := wsManager.NewGameEndEvent(roomID, message)
	h.manager.SendEventToRoom(roomID, event)
}
