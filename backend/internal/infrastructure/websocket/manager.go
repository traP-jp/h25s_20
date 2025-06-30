package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/rs/zerolog/log"
)

// ユーザーの状態情報を保存するための構造体
type UserState struct {
	UserID      int         `json:"user_id"`
	RoomID      *int        `json:"room_id"`
	LastSeenAt  time.Time   `json:"last_seen_at"`
	DeleteTimer *time.Timer `json:"-"` // JSONには含めない
}

type Client struct {
	ID     string
	UserID int  // ユーザーID
	RoomID *int // 参加しているroomのID（未参加の場合はnil）
	Conn   *websocket.Conn
	Cancel context.CancelFunc
}

// RoomUsecaseのインターフェース定義（循環importを避けるため）
type RoomUsecaseInterface interface {
	SetPlayerDisconnected(roomID int, playerID int) (*domain.Room, error)
	SetPlayerReconnected(roomID int, playerID int) (*domain.Room, error)
	RemoveDisconnectedPlayer(roomID int, playerID int) (*domain.Room, error)
	GetRoomByID(roomID int) (*domain.Room, error)
	RemovePlayerFromRoom(roomID int, playerID int) (*domain.Room, error)
}

type Manager struct {
	clients           map[string]*Client
	userClients       map[int]*Client    // UserID -> Client のマッピング
	roomClients       map[int][]*Client  // RoomID -> []*Client のマッピング
	disconnectedUsers map[int]*UserState // 切断されたユーザーの状態を一時保存
	mutex             sync.RWMutex
	deleteTimeout     time.Duration        // ユーザー削除までのタイムアウト時間
	roomUsecase       RoomUsecaseInterface // RoomUsecaseとの連携用
}

// 後方互換性のため残す（非推奨）
type NotificationMessage struct {
	Event   string      `json:"event"`
	Content interface{} `json:"content"`
}

// 非推奨: 新しいWebSocketEventとイベント固有の構造体を使用してください
type StandardEventContent struct {
	UserID   int         `json:"user_id,omitempty"`
	UserName string      `json:"user_name,omitempty"`
	RoomID   int         `json:"room_id,omitempty"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// 非推奨: 新しいBoardUpdateEventContentを使用してください
type BoardUpdateContent struct {
	StandardEventContent
	Board     interface{} `json:"board"`
	GainScore int         `json:"gain_score"`
}

func NewManager() *Manager {
	return &Manager{
		clients:           make(map[string]*Client),
		userClients:       make(map[int]*Client),
		roomClients:       make(map[int][]*Client),
		disconnectedUsers: make(map[int]*UserState),
		deleteTimeout:     10 * time.Second, // ゲーム時間と同じ120秒後に削除
	}
}

// NewManagerWithTimeout creates a new Manager with custom delete timeout
func NewManagerWithTimeout(timeout time.Duration) *Manager {
	manager := NewManager()
	manager.deleteTimeout = timeout
	return manager
}

// SetRoomUsecase sets the room usecase for disconnection handling
func (m *Manager) SetRoomUsecase(roomUsecase RoomUsecaseInterface) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.roomUsecase = roomUsecase
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

	// 切断されたユーザーの状態を復元する
	if disconnectedUser, exists := m.disconnectedUsers[userID]; exists {
		// 削除タイマーをキャンセル
		if disconnectedUser.DeleteTimer != nil {
			disconnectedUser.DeleteTimer.Stop()
		}

		// 前の状態を復元
		client.RoomID = disconnectedUser.RoomID
		if client.RoomID != nil {
			// roomClientsにも再登録
			m.roomClients[*client.RoomID] = append(m.roomClients[*client.RoomID], client)

			// RoomUsecaseに再接続を通知
			if m.roomUsecase != nil {
				if _, err := m.roomUsecase.SetPlayerReconnected(*client.RoomID, userID); err != nil {
					log.Warn().Err(err).
						Int("room_id", *client.RoomID).
						Int("user_id", userID).
						Msg("Failed to notify room usecase about player reconnection")
				}
			}
		}

		// disconnectedUsersから削除
		delete(m.disconnectedUsers, userID)

		log.Info().
			Str("client_id", clientID).
			Int("user_id", userID).
			Interface("restored_room_id", client.RoomID).
			Msg("WebSocket client reconnected and state restored")
	} else {
		log.Info().
			Str("client_id", clientID).
			Int("user_id", userID).
			Msg("WebSocket client connected for the first time")
	}

	m.clients[clientID] = client
	m.userClients[userID] = client
}

func (m *Manager) RemoveClient(clientID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.clients[clientID]
	if !exists {
		return
	}

	// userClientsから削除（即座に）
	delete(m.userClients, client.UserID)

	// roomClientsから削除（即座に）
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

	// RoomUsecaseに切断を通知
	if m.roomUsecase != nil && client.RoomID != nil {
		if _, err := m.roomUsecase.SetPlayerDisconnected(*client.RoomID, client.UserID); err != nil {
			log.Warn().Err(err).
				Int("room_id", *client.RoomID).
				Int("user_id", client.UserID).
				Msg("Failed to notify room usecase about player disconnection")
		}
	}

	// ユーザーの状態を一時保存し、遅延削除タイマーを設定
	userState := &UserState{
		UserID:     client.UserID,
		RoomID:     client.RoomID,
		LastSeenAt: time.Now(),
	}

	// 削除タイマーを設定（デッドロック回避のため、別ゴルーチンで実行）
	userState.DeleteTimer = time.AfterFunc(m.deleteTimeout, func() {
		// 非同期でユーザー削除処理を実行してデッドロックを回避
		go m.handleUserDeletion(client.UserID)
	})

	m.disconnectedUsers[client.UserID] = userState

	// クライアント接続を削除
	client.Cancel()
	delete(m.clients, clientID)

	log.Info().
		Str("client_id", clientID).
		Int("user_id", client.UserID).
		Dur("delete_timeout", m.deleteTimeout).
		Msg("WebSocket client disconnected, scheduled for delayed deletion")
}

func (m *Manager) JoinRoom(userID, roomID int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.userClients[userID]
	if !exists {
		log.Warn().Int("user_id", userID).Msg("User not connected via WebSocket")
		return nil // WebSocket接続がない場合はエラーにしない
	}

	// 既に別のroomに参加している場合の処理
	if client.RoomID != nil {
		oldRoomID := *client.RoomID

		// 古い部屋でのゲーム進行状況を確認し、適切な処理を行う
		if m.roomUsecase != nil {
			if err := m.handlePreviousRoomExit(userID, oldRoomID); err != nil {
				log.Warn().Err(err).
					Int("user_id", userID).
					Int("old_room_id", oldRoomID).
					Int("new_room_id", roomID).
					Msg("Failed to handle previous room exit properly")
			}
		}

		// WebSocketレベルでの退室処理
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

// handlePreviousRoomExit handles the exit from previous room based on game state
func (m *Manager) handlePreviousRoomExit(userID, roomID int) error {
	// RoomUsecaseから部屋の状態を取得
	room, err := m.roomUsecase.GetRoomByID(roomID)
	if err != nil {
		return fmt.Errorf("failed to get room %d: %w", roomID, err)
	}

	// ゲームが進行していない状態の場合は即時削除
	switch room.State {
	case domain.StateWaitingForPlayers, domain.StateAllReady, domain.StateGameEnded:
		// ゲームが進行していない状態では即時削除
		if _, err := m.roomUsecase.RemovePlayerFromRoom(roomID, userID); err != nil {
			log.Warn().Err(err).
				Int("room_id", roomID).
				Int("user_id", userID).
				Str("room_state", room.State.String()).
				Msg("Failed to immediately remove player from non-progressing room")
			return err
		}

		log.Info().
			Int("room_id", roomID).
			Int("user_id", userID).
			Str("room_state", room.State.String()).
			Msg("Player immediately removed from non-progressing room")

	case domain.StateCountdown, domain.StateGameInProgress:
		// ゲーム進行中の場合は通常の切断処理
		if _, err := m.roomUsecase.SetPlayerDisconnected(roomID, userID); err != nil {
			log.Warn().Err(err).
				Int("room_id", roomID).
				Int("user_id", userID).
				Str("room_state", room.State.String()).
				Msg("Failed to set player as disconnected in progressing room")
			return err
		}

		log.Info().
			Int("room_id", roomID).
			Int("user_id", userID).
			Str("room_state", room.State.String()).
			Msg("Player marked as disconnected in progressing room")
	}

	return nil
}

func (m *Manager) LeaveRoom(userID int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.userClients[userID]
	if !exists || client.RoomID == nil {
		return nil
	}

	roomID := *client.RoomID
	m.leaveRoomInternal(client)

	log.Info().
		Int("user_id", userID).
		Int("room_id", roomID).
		Msg("User left room via WebSocket (connection maintained)")

	return nil
}

func (m *Manager) leaveRoomInternal(client *Client) {
	if client.RoomID == nil {
		return
	}

	roomID := *client.RoomID
	if clients, exists := m.roomClients[roomID]; exists {
		for i := 0; i < len(clients); i++ {
			if clients[i].ID == client.ID {
				// 境界チェック付きでスライスから削除
				if i < len(clients)-1 {
					m.roomClients[roomID] = append(clients[:i], clients[i+1:]...)
				} else {
					m.roomClients[roomID] = clients[:i]
				}
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

// 新しい統一されたイベント送信メソッド群

// SendEvent sends a structured WebSocketEvent to all clients
func (m *Manager) SendEvent(event WebSocketEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal WebSocketEvent")
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
		Str("event", event.Event).
		Int("client_count", len(clients)).
		Msg("Sent structured event to all clients")
}

// SendEventToRoom sends a structured WebSocketEvent to room participants
func (m *Manager) SendEventToRoom(roomID int, event WebSocketEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal WebSocketEvent")
		return
	}

	m.mutex.RLock()
	clients, exists := m.roomClients[roomID]
	if !exists {
		m.mutex.RUnlock()
		log.Warn().Int("room_id", roomID).Msg("No clients in room")
		return
	}
	clientsCopy := make([]*Client, len(clients))
	copy(clientsCopy, clients)
	m.mutex.RUnlock()

	for _, client := range clientsCopy {
		go m.sendToClient(client.ID, client, data)
	}

	log.Info().
		Str("event", event.Event).
		Int("room_id", roomID).
		Int("client_count", len(clientsCopy)).
		Msg("Sent structured event to room clients")
}

// SendEventToUser sends a structured WebSocketEvent to a specific user
func (m *Manager) SendEventToUser(userID int, event WebSocketEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	m.mutex.RLock()
	client, exists := m.userClients[userID]
	if !exists {
		m.mutex.RUnlock()
		log.Warn().Int("user_id", userID).Msg("User not connected via WebSocket")
		return nil
	}
	m.mutex.RUnlock()

	go m.sendToClient(client.ID, client, data)

	log.Info().
		Str("event", event.Event).
		Int("user_id", userID).
		Msg("Sent structured event to user")

	return nil
}

// 遅延削除システム管理機能

// GetDisconnectedUsers returns information about users scheduled for deletion
func (m *Manager) GetDisconnectedUsers() map[int]*UserState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[int]*UserState)
	for userID, state := range m.disconnectedUsers {
		result[userID] = &UserState{
			UserID:     state.UserID,
			RoomID:     state.RoomID,
			LastSeenAt: state.LastSeenAt,
			// DeleteTimer は含めない（goroutineセーフではないため）
		}
	}
	return result
}

// ForceDeleteUser immediately deletes a user from disconnected users (cancels delayed deletion)
func (m *Manager) ForceDeleteUser(userID int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if state, exists := m.disconnectedUsers[userID]; exists {
		if state.DeleteTimer != nil {
			state.DeleteTimer.Stop()
		}
		delete(m.disconnectedUsers, userID)

		log.Info().
			Int("user_id", userID).
			Msg("User force deleted from disconnected users")
		return true
	}
	return false
}

// IsUserDisconnected checks if a user is in the disconnected users list
func (m *Manager) IsUserDisconnected(userID int) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.disconnectedUsers[userID]
	return exists
}

// GetDeleteTimeout returns the current delete timeout duration
func (m *Manager) GetDeleteTimeout() time.Duration {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.deleteTimeout
}

// SetDeleteTimeout sets a new delete timeout duration
func (m *Manager) SetDeleteTimeout(timeout time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.deleteTimeout = timeout

	log.Info().
		Dur("new_timeout", timeout).
		Msg("Delete timeout updated")
}

// GetDisconnectedUserStats returns statistics about disconnected users
func (m *Manager) GetDisconnectedUserStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_disconnected": len(m.disconnectedUsers),
		"timeout_duration":   m.deleteTimeout,
		"users":              make([]map[string]interface{}, 0),
	}

	users := make([]map[string]interface{}, 0, len(m.disconnectedUsers))
	for userID, state := range m.disconnectedUsers {
		userInfo := map[string]interface{}{
			"user_id":           userID,
			"room_id":           state.RoomID,
			"last_seen_at":      state.LastSeenAt,
			"time_until_delete": m.deleteTimeout - time.Since(state.LastSeenAt),
		}
		users = append(users, userInfo)
	}
	stats["users"] = users

	return stats
}

// handleUserDeletion デッドロック回避のための非同期ユーザー削除処理
func (m *Manager) handleUserDeletion(userID int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// タイムアウト後にユーザー状態を完全削除
	if state, exists := m.disconnectedUsers[userID]; exists {
		// RoomUsecaseに永続削除を通知
		if m.roomUsecase != nil && state.RoomID != nil {
			if _, err := m.roomUsecase.RemoveDisconnectedPlayer(*state.RoomID, userID); err != nil {
				log.Warn().Err(err).
					Int("room_id", *state.RoomID).
					Int("user_id", userID).
					Msg("Failed to notify room usecase about player removal")
			}
		}

		delete(m.disconnectedUsers, userID)
		log.Info().
			Int("user_id", userID).
			Dur("after_disconnect", time.Since(state.LastSeenAt)).
			Msg("User permanently deleted after timeout")
	}
}
