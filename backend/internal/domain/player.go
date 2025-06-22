package domain

import "time"

type Player struct {
	ID              int
	UserName        string
	IsReady         bool
	HasClosedResult bool // 結果表示を閉じたかどうか
	Score           int
	IsConnected     bool       // WebSocket接続状態
	LastSeenAt      *time.Time // 最後に確認された時刻（切断時に設定）
}
