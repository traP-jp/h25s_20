package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	wsManager "github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WebSocketHandler struct {
	manager *wsManager.Manager
}

func NewWebSocketHandler(manager *wsManager.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		manager: manager,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// ユーザーIDをクエリパラメータから取得
	userIDStr := c.QueryParam("user_id")
	if userIDStr == "" {
		return c.JSON(400, map[string]string{"error": "user_id is required"})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid user_id"})
	}

	// WebSocket接続をアップグレード
	conn, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade WebSocket connection")
		return err
	}

	// クライアントIDを生成
	clientID := uuid.New().String()

	// コンテキストとキャンセル関数を作成
	ctx, cancel := context.WithCancel(context.Background())

	// クライアントをマネージャーに登録
	h.manager.AddClient(clientID, userID, conn, cancel)

	// 接続完了メッセージを送信
	welcomeMessage := wsManager.NotificationMessage{
		Event: "connection",
		Content: map[string]interface{}{
			"client_id": clientID,
			"user_id":   userID,
			"message":   "Connected successfully",
			"timestamp": time.Now().Unix(),
		},
	}

	if err := h.manager.SendToUser(userID, welcomeMessage); err != nil {
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
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
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
	message := wsManager.NotificationMessage{
		Event:   event,
		Content: content,
	}
	h.manager.BroadcastToAll(message)
}

// room参加者全員にメッセージを送信
func (h *WebSocketHandler) BroadcastToRoom(roomID int, event string, content interface{}) {
	message := wsManager.NotificationMessage{
		Event:   event,
		Content: content,
	}
	h.manager.BroadcastToRoom(roomID, message)
}

// room未参加者全員にメッセージを送信
func (h *WebSocketHandler) BroadcastToNonRoomMembers(event string, content interface{}) {
	message := wsManager.NotificationMessage{
		Event:   event,
		Content: content,
	}
	h.manager.BroadcastToNonRoomMembers(message)
}

// 特定のユーザーにメッセージを送信
func (h *WebSocketHandler) SendToUser(userID int, event string, content interface{}) error {
	message := wsManager.NotificationMessage{
		Event:   event,
		Content: content,
	}
	return h.manager.SendToUser(userID, message)
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
