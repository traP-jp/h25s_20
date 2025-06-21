package domain

type Room struct {
	ID         int
	Name       string
	GameBoards []GameBoard
	IsOpened   bool
	Players    []Player
	ResultLog  []Result
}

type GameBoard struct {
	Version int
	Board   [][]int
	Size    int
}

type Result struct {
	ID     int
	Time   string
	Scores []PlayerScore
}

type PlayerScore struct {
	ID       int
	PlayerId string
	Score    int
}
