package notification

import (
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
)

type Service struct {
	roomUsecase *usecase.RoomUsecase
	wsManager   *websocket.Manager
}

func NewService(roomUsecase *usecase.RoomUsecase, wsManager *websocket.Manager) *Service {
	return &Service{
		roomUsecase: roomUsecase,
		wsManager:   wsManager,
	}
}

// 全クライアントに通知
func (s *Service) NotifyAll(event string, content interface{}) {
	message := websocket.NotificationMessage{
		Event:   event,
		Content: content,
	}
	s.wsManager.BroadcastToAll(message)
}

// 特定roomの参加者全員に通知
func (s *Service) NotifyRoom(roomID int, event string, content interface{}) {
	message := websocket.NotificationMessage{
		Event:   event,
		Content: content,
	}
	s.wsManager.BroadcastToRoom(roomID, message)
}

// room未参加者全員に通知
func (s *Service) NotifyNonRoomMembers(event string, content interface{}) {
	message := websocket.NotificationMessage{
		Event:   event,
		Content: content,
	}
	s.wsManager.BroadcastToNonRoomMembers(message)
}

// 特定ユーザーに通知
func (s *Service) NotifyUser(userID int, event string, content interface{}) error {
	message := websocket.NotificationMessage{
		Event:   event,
		Content: content,
	}
	return s.wsManager.SendToUser(userID, message)
}

// 接続状況の取得
func (s *Service) GetConnectionStats() map[string]interface{} {
	return map[string]interface{}{
		"total_connected": s.wsManager.GetClientCount(),
	}
}

// 特定roomの接続状況
func (s *Service) GetRoomConnectionStats(roomID int) map[string]interface{} {
	return map[string]interface{}{
		"room_id":         roomID,
		"connected_count": s.wsManager.GetRoomClientCount(roomID),
	}
}
