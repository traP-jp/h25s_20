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
	EventPlayerAllReady = "player_all_ready"

	// ルーム関連
	EventRoomClosed = "room_closed"

	// ゲーム関連
	EventGameStarted        = "game_started"
	EventGameStart          = "game_start"
	EventCountdownStartGame = "countdown_start_game" // ゲーム開始時のカウントダウン開始
	EventCountdownGame      = "countdown_game"       // ゲーム開始時のカウントダウン
	EventCountdownStartEnd  = "countdown_start_end"  // ゲーム終了時のカウントダウン開始
	EventCountdownEndGame   = "countdown_end_game"   // ゲーム終了時のカウントダウン
	EventBoardUpdated       = "board_updated"
	EventResultClosed       = "result_closed"
	EventGameEnded          = "game_ended"
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

// プレイヤー参加イベント用（ルーム情報付き）
type PlayerJoinedEventContent struct {
	BaseEventContent
	Room RoomInfo `json:"room"`
}

func (p PlayerJoinedEventContent) GetEventType() string {
	return "player_joined"
}

// プレイヤー退出イベント用（ルーム情報付き）
type PlayerLeftEventContent struct {
	BaseEventContent
	Room RoomInfo `json:"room"`
}

func (p PlayerLeftEventContent) GetEventType() string {
	return "player_left"
}

// ルーム情報
type RoomInfo struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	State    string       `json:"state"`
	IsOpened bool         `json:"is_opened"`
	Players  []PlayerInfo `json:"players"`
}

// プレイヤー情報
type PlayerInfo struct {
	ID              int    `json:"id"`
	UserName        string `json:"user_name"`
	IsReady         bool   `json:"is_ready"`
	HasClosedResult bool   `json:"has_closed_result"`
	Score           int    `json:"score"`
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

func NewPlayerJoinedEvent(userID int, userName string, room RoomInfo) WebSocketEvent {
	return WebSocketEvent{
		Event: EventPlayerJoined,
		Content: PlayerJoinedEventContent{
			BaseEventContent: BaseEventContent{
				UserID:   userID,
				UserName: userName,
				RoomID:   room.ID,
			},
			Room: room,
		},
	}
}

func NewPlayerLeftEvent(userID int, userName string, room RoomInfo) WebSocketEvent {
	return WebSocketEvent{
		Event: EventPlayerLeft,
		Content: PlayerLeftEventContent{
			BaseEventContent: BaseEventContent{
				UserID:   userID,
				UserName: userName,
				RoomID:   room.ID,
			},
			Room: room,
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

func NewCountdownStartGameEvent(roomID int, message string, countdown int) WebSocketEvent {
	return WebSocketEvent{
		Event: EventCountdownStartGame,
		Content: CountdownEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
			Countdown: countdown,
		},
	}
}

func NewCountdownGameEvent(roomID int, count int) WebSocketEvent {
	return WebSocketEvent{
		Event: EventCountdownGame,
		Content: CountdownEventContent{
			BaseEventContent: BaseEventContent{
				RoomID: roomID,
			},
			Count: count,
		},
	}
}

func NewCountdownEndGameEvent(roomID int, count int) WebSocketEvent {
	return WebSocketEvent{
		Event: EventCountdownEndGame,
		Content: CountdownEventContent{
			BaseEventContent: BaseEventContent{
				RoomID: roomID,
			},
			Count: count,
		},
	}
}

func NewCountdownStartEndEvent(roomID int, message string, countdown int) WebSocketEvent {
	return WebSocketEvent{
		Event: EventCountdownStartEnd,
		Content: CountdownEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
			Countdown: countdown,
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

func NewGameEndEvent(roomID int, message string) WebSocketEvent {
	return WebSocketEvent{
		Event: EventGameEnded,
		Content: GameStartEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
		},
	}
}

func NewPlayerAllReadyEvent(roomID int, message string) WebSocketEvent {
	return WebSocketEvent{
		Event: EventPlayerAllReady,
		Content: PlayerEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
		},
	}
}

func NewRoomClosedEvent(roomID int, message string) WebSocketEvent {
	return WebSocketEvent{
		Event: EventRoomClosed,
		Content: PlayerEventContent{
			BaseEventContent: BaseEventContent{
				RoomID:  roomID,
				Message: message,
			},
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

// domain.RoomからRoomInfoへの変換ヘルパー関数
// この関数はインポートループを避けるため、domain.Roomを受け取らずにinterfaceを使用
type DomainRoom interface {
	GetID() int
	GetName() string
	GetState() string
	GetIsOpened() bool
	GetPlayers() []DomainPlayer
}

type DomainPlayer interface {
	GetID() int
	GetUserName() string
	GetIsReady() bool
	GetHasClosedResult() bool
	GetScore() int
}

func ConvertToRoomInfo(id int, name string, state string, isOpened bool, players []PlayerInfo) RoomInfo {
	return RoomInfo{
		ID:       id,
		Name:     name,
		State:    state,
		IsOpened: isOpened,
		Players:  players,
	}
}

func ConvertToPlayerInfo(id int, userName string, isReady bool, hasClosedResult bool, score int) PlayerInfo {
	return PlayerInfo{
		ID:              id,
		UserName:        userName,
		IsReady:         isReady,
		HasClosedResult: hasClosedResult,
		Score:           score,
	}
}
