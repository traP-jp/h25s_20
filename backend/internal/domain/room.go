package domain

type Room struct {
	ID       int
	Name     string
	Boards   []Board
	IsOpened bool
	Players  []Player
	Result   []Result
}

type Board struct {
	Version int
	Board   []int
}

type Result struct {
	ID     int
	Time   string
	Scores []PlayerScore
}

type PlayerScore struct {
	ID       int
	PlayerId string
	Socre    int
}
