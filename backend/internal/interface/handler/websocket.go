package handler

import (
	"context"
	"strconv"
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
}

func NewWebSocketHandler(manager *wsManager.Manager, roomUsecase *usecase.RoomUsecase) *WebSocketHandler {
	return &WebSocketHandler{
		manager:     manager,
		roomUsecase: roomUsecase,
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

	if err := h.manager.NotifyUser(userID, welcomeMessage.Event, welcomeMessage.Content); err != nil {
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

func (h *WebSocketHandler) HandleFormulaEvent(userID, roomID int, formula string) {
	board, gainScore, err := h.roomUsecase.ApplyFormula(roomID, userID, formula)
	if err != nil {
		h.SendToUser(userID, "formula_error", err.Error())
		return
	}
	// 成功時はルーム全体に通知
	h.BroadcastToRoom(roomID, "formula_applied", map[string]interface{}{
		"board":      board,
		"gain_score": gainScore,
		"user_id":    userID,
	})
}
