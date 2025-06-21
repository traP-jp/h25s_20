package domain

type Room struct {
	Id       int
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
	Id     int
	Time   string
	Scores []Playerscore
}

type PlayerScore struct {
	Id       int
	PlayerId string
	Socre    int
}
