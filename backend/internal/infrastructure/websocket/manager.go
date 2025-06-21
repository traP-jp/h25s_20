package websocket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
	"github.com/rs/zerolog/log"
)

type Client struct {
	ID     string
	UserID int  // ユーザーID
	RoomID *int // 参加しているroomのID（未参加の場合はnil）
	Conn   *websocket.Conn
	Cancel context.CancelFunc
}

type Manager struct {
	clients     map[string]*Client
	userClients map[int]*Client   // UserID -> Client のマッピング
	roomClients map[int][]*Client // RoomID -> []*Client のマッピング
	mutex       sync.RWMutex
}

type NotificationMessage struct {
	Event   string      `json:"event"`
	Content interface{} `json:"content"`
}

// WebSocketイベントの標準化された構造体
type StandardEventContent struct {
	UserID   int         `json:"user_id,omitempty"`
	UserName string      `json:"user_name,omitempty"`
	RoomID   int         `json:"room_id,omitempty"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// ボード更新用の専用構造体
type BoardUpdateContent struct {
	StandardEventContent
	Board     interface{} `json:"board"`
	GainScore int         `json:"gain_score"`
}

func NewManager() *Manager {
	return &Manager{
		clients:     make(map[string]*Client),
		userClients: make(map[int]*Client),
		roomClients: make(map[int][]*Client),
	}
}

func (m *Manager) AddClient(clientID string, userID int, conn *websocket.Conn, cancel context.CancelFunc) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client := &Client{
		ID:     clientID,
		UserID: userID,
		RoomID: nil,
		Conn:   conn,
		Cancel: cancel,
	}

	m.clients[clientID] = client
	m.userClients[userID] = client

	log.Info().
		Str("client_id", clientID).
		Int("user_id", userID).
		Msg("WebSocket client connected")
}

func (m *Manager) RemoveClient(clientID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.clients[clientID]
	if !exists {
		return
	}

	// userClientsから削除
	delete(m.userClients, client.UserID)

	// roomClientsから削除
	if client.RoomID != nil {
		roomID := *client.RoomID
		if clients, exists := m.roomClients[roomID]; exists {
			// スライスから該当クライアントを削除
			for i, c := range clients {
				if c.ID == clientID {
					m.roomClients[roomID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			// room内にクライアントがいなくなった場合はマップエントリも削除
			if len(m.roomClients[roomID]) == 0 {
				delete(m.roomClients, roomID)
			}
		}
	}

	// クライアントを削除
	client.Cancel()
	delete(m.clients, clientID)

	log.Info().
		Str("client_id", clientID).
		Int("user_id", client.UserID).
		Msg("WebSocket client disconnected")
}

func (m *Manager) JoinRoom(userID, roomID int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.userClients[userID]
	if !exists {
		log.Warn().Int("user_id", userID).Msg("User not connected via WebSocket")
		return nil // WebSocket接続がない場合はエラーにしない
	}

	// 既に別のroomに参加している場合は先に退出
	if client.RoomID != nil {
		m.leaveRoomInternal(client)
	}

	// 新しいroomに参加
	client.RoomID = &roomID
	m.roomClients[roomID] = append(m.roomClients[roomID], client)

	log.Info().
		Int("user_id", userID).
		Int("room_id", roomID).
		Msg("User joined room via WebSocket")

	return nil
}

func (m *Manager) LeaveRoom(userID int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.userClients[userID]
	if !exists || client.RoomID == nil {
		return nil
	}

	m.leaveRoomInternal(client)

	log.Info().
		Int("user_id", userID).
		Int("room_id", *client.RoomID).
		Msg("User left room via WebSocket")

	return nil
}

func (m *Manager) leaveRoomInternal(client *Client) {
	if client.RoomID == nil {
		return
	}

	roomID := *client.RoomID
	if clients, exists := m.roomClients[roomID]; exists {
		for i, c := range clients {
			if c.ID == client.ID {
				m.roomClients[roomID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(m.roomClients[roomID]) == 0 {
			delete(m.roomClients, roomID)
		}
	}

	client.RoomID = nil
}

// 全クライアントに通知
func (m *Manager) NotifyAll(event string, content interface{}) {
	message := NotificationMessage{
		Event:   event,
		Content: content,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return
	}

	m.mutex.RLock()
	clients := make([]*Client, 0, len(m.clients))
	for _, client := range m.clients {
		clients = append(clients, client)
	}
	m.mutex.RUnlock()

	for _, client := range clients {
		go m.sendToClient(client.ID, client, data)
	}

	log.Info().
		Str("event", event).
		Int("client_count", len(clients)).
		Msg("Broadcasted message to all clients")
}

// 特定roomの参加者全員に通知
func (m *Manager) NotifyRoom(roomID int, event string, content interface{}) {
	message := NotificationMessage{
		Event:   event,
		Content: content,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return
	}

	m.mutex.RLock()
	clients, exists := m.roomClients[roomID]
	if !exists {
		m.mutex.RUnlock()
		log.Warn().Int("room_id", roomID).Msg("No clients in room")
		return
	}
	// コピーを作成
	clientsCopy := make([]*Client, len(clients))
	copy(clientsCopy, clients)
	m.mutex.RUnlock()

	for _, client := range clientsCopy {
		go m.sendToClient(client.ID, client, data)
	}

	log.Info().
		Str("event", event).
		Int("room_id", roomID).
		Int("client_count", len(clientsCopy)).
		Msg("Broadcasted message to room clients")
}

// room未参加者全員に通知
func (m *Manager) NotifyNonRoomMembers(event string, content interface{}) {
	message := NotificationMessage{
		Event:   event,
		Content: content,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return
	}

	m.mutex.RLock()
	clients := make([]*Client, 0)
	for _, client := range m.clients {
		if client.RoomID == nil {
			clients = append(clients, client)
		}
	}
	m.mutex.RUnlock()

	for _, client := range clients {
		go m.sendToClient(client.ID, client, data)
	}

	log.Info().
		Str("event", event).
		Int("client_count", len(clients)).
		Msg("Broadcasted message to non-room clients")
}

// 特定ユーザーに通知
func (m *Manager) NotifyUser(userID int, event string, content interface{}) error {
	message := NotificationMessage{
		Event:   event,
		Content: content,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	m.mutex.RLock()
	client, exists := m.userClients[userID]
	m.mutex.RUnlock()

	if !exists {
		log.Warn().Int("user_id", userID).Msg("User not connected via WebSocket")
		return nil
	}

	go m.sendToClient(client.ID, client, data)

	log.Info().
		Int("user_id", userID).
		Str("event", event).
		Msg("Sent message to user")

	return nil
}

// 低レベルAPI: 複数クライアントに直接メッセージ送信
func (m *Manager) SendToClients(clients []*Client, data []byte) {
	for _, client := range clients {
		go m.sendToClient(client.ID, client, data)
	}
}

// 低レベルAPI: 特定クライアントの取得
func (m *Manager) GetAllClients() []*Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	clients := make([]*Client, 0, len(m.clients))
	for _, client := range m.clients {
		clients = append(clients, client)
	}
	return clients
}

func (m *Manager) GetClientsInRoom(roomID int) []*Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	clients, exists := m.roomClients[roomID]
	if !exists {
		return []*Client{}
	}

	// コピーを返す
	result := make([]*Client, len(clients))
	copy(result, clients)
	return result
}

func (m *Manager) GetClientsNotInRoom() []*Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	clients := make([]*Client, 0)
	for _, client := range m.clients {
		if client.RoomID == nil {
			clients = append(clients, client)
		}
	}
	return clients
}

func (m *Manager) GetClientByUser(userID int) *Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.userClients[userID]
}

func (m *Manager) sendToClient(clientID string, client *Client, data []byte) {
	err := client.Conn.Write(context.Background(), websocket.MessageText, data)
	if err != nil {
		log.Error().Err(err).Str("client_id", clientID).Msg("Failed to send message to client")
		m.RemoveClient(clientID)
	}
}

func (m *Manager) GetClientCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.clients)
}

func (m *Manager) GetRoomClientCount(roomID int) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	clients, exists := m.roomClients[roomID]
	if !exists {
		return 0
	}
	return len(clients)
}

// 接続状況の取得
func (m *Manager) GetConnectionStats() map[string]interface{} {
	return map[string]interface{}{
		"total_connected": m.GetClientCount(),
	}
}

// 特定roomの接続状況
func (m *Manager) GetRoomConnectionStats(roomID int) map[string]interface{} {
	return map[string]interface{}{
		"room_id":         roomID,
		"connected_count": m.GetRoomClientCount(roomID),
	}
}
