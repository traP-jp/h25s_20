package domain

type Player struct {
	ID              int
	UserName        string
	IsReady         bool
	HasClosedResult bool // 結果表示を閉じたかどうか
	Score           int
}
