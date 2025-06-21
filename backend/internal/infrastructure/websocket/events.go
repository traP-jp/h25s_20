package websocket

// WebSocketイベント名の定数定義
const (
	// 接続関連
	EventConnection = "connection"

	// プレイヤー関連
	EventPlayerJoined   = "player_joined"
	EventPlayerReady    = "player_ready"
	EventPlayerCanceled = "player_canceled"
	EventPlayerLeft     = "player_left"

	// ゲーム関連
	EventGameStarted    = "game_started"
	EventGameStart      = "game_start"
	EventCountdownStart = "countdown_start"
	EventCountdown      = "countdown"
	EventBoardUpdated   = "board_updated"
	EventResultClosed   = "result_closed"
)

// 統一されたWebSocketイベントの基本構造
type WebSocketEvent struct {
	Event   string       `json:"event"`
	Content EventContent `json:"content"`
}

// イベントコンテンツの基本インターface
type EventContent interface {
	GetEventType() string
}

// 基本的なイベントコンテンツ
type BaseEventContent struct {
	UserID    int    `json:"user_id,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	RoomID    int    `json:"room_id,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (b BaseEventContent) GetEventType() string {
	return "base"
}

// 接続イベント用
type ConnectionEventContent struct {
	BaseEventContent
	ClientID string `json:"client_id"`
}

func (c ConnectionEventContent) GetEventType() string {
	return "connection"
}

// プレイヤーアクション用
type PlayerEventContent struct {
	BaseEventContent
}

func (p PlayerEventContent) GetEventType() string {
	return "player"
}

// ゲーム開始用
type GameStartEventContent struct {
	BaseEventContent
}

func (g GameStartEventContent) GetEventType() string {
	return "game_start"
}

// カウントダウン用
type CountdownEventContent struct {
	BaseEventContent
	Count     int `json:"count,omitempty"`
	Countdown int `json:"countdown,omitempty"`
}

func (c CountdownEventContent) GetEventType() string {
	return "countdown"
}

// ボード更新用
type BoardUpdateEventContent struct {
	BaseEventContent
	Board     BoardData `json:"board"`
	GainScore int       `json:"gain_score"`
}

type BoardData struct {
	Content []int `json:"content"`
	Version int   `json:"version"`
	Size    int   `json:"size"`
}

func (b BoardUpdateEventContent) GetEventType() string {
	return "board_update"
}

// ゲーム開始時のボード送信用
type GameStartBoardEventContent struct {
	BaseEventContent
	Board BoardData `json:"board"`
}

func (g GameStartBoardEventContent) GetEventType() string {
	return "game_start_board"
}

// イベント作成のヘルパー関数群

func NewConnectionEvent(clientID string, userID int, message string, timestamp int64) WebSocketEvent {
	return WebSocketEvent{
		Event: EventConnection,
		Content: ConnectionEventContent{
			BaseEventContent: BaseEventContent{
				UserID:    userID,
				Message:   message,
				Timestamp: timestamp,
			},
			ClientID: clientID,
		},
	}
}

func NewPlayerEvent(eventType string, userID int, userName string, roomID int) WebSocketEvent {
	return WebSocketEvent{
		Event: eventType,
		Content: PlayerEventContent{
			BaseEventContent: BaseEventContent{
				UserID:   userID,
				UserName: userName,
				RoomID:   roomID,
			},
		},
	}
}

func NewGameStartEvent(roomID int, message string) WebSocketEvent {
	return WebSocketEvent{
		Event: EventGameStarted,
		Content: GameStartEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
		},
	}
}

func NewCountdownStartEvent(roomID int, message string, countdown int) WebSocketEvent {
	return WebSocketEvent{
		Event: EventCountdownStart,
		Content: CountdownEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
			Countdown: countdown,
		},
	}
}

func NewCountdownEvent(roomID int, count int) WebSocketEvent {
	return WebSocketEvent{
		Event: EventCountdown,
		Content: CountdownEventContent{
			BaseEventContent: BaseEventContent{
				RoomID: roomID,
			},
			Count: count,
		},
	}
}

func NewBoardUpdateEvent(userID int, userName string, roomID int, board BoardData, gainScore int) WebSocketEvent {
	return WebSocketEvent{
		Event: EventBoardUpdated,
		Content: BoardUpdateEventContent{
			BaseEventContent: BaseEventContent{
				UserID:   userID,
				UserName: userName,
				RoomID:   roomID,
				Message:  "Board updated",
			},
			Board:     board,
			GainScore: gainScore,
		},
	}
}

func NewGameStartBoardEvent(roomID int, message string, board BoardData) WebSocketEvent {
	return WebSocketEvent{
		Event: EventGameStart,
		Content: GameStartBoardEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
			Board: board,
		},
	}
}
