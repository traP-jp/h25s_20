package domain

import (
	"fmt"
	"math/rand"
)

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

//回答確認後一致した行と列の情報を保持
type Match struct {
	LineType string //"row"または"col"
	Index    int	
}

// 新規盤面の作成
func NewBoard() GameBoard {
	size := 4
	gb := &GameBoard{
		Version: 1,
		Board:   make([][]int, size),
	}
	//盤面の初期化
	for i := range gb.Board {
		gb.Board[i] = make([]int, size)
	}
	//盤面全体を1から9のランダムな整数で埋める
	for i := range gb.Size {
		gb.PopulateRow(i)
	}
	return *gb
}

// 指定の列を1から9のランダムな整数で埋める
func (gb GameBoard) PopulateRow(row int) {
	for i := 0; i < gb.Size; i++ {
		gb.Board[row][i] = rand.Intn(9) + 1 //1-9の乱数
	}
}

// 指定の行を1から9のランダムな整数で埋める
func (gb GameBoard) PopulateColumn(col int) {
	for i := 0; i < gb.Size; i++ {
		gb.Board[i][col] = rand.Intn(9) + 1 // 1-9の乱数
	}
}

// UpdateLine は指定された行または列を新しいランダムな数字で更新します。
func (gb *GameBoard) UpdateLine(lineType string, index int) error {
	if !(0 <= index && index < gb.Size) {
		return fmt.Errorf("インデックスは0から%dの間で指定してください", gb.Size-1)
	}
	switch lineType {
	case "row":
		gb.PopulateRow(index)
		return nil
	case "col":
		gb.PopulateColumn(index)
		return nil
	default:
		return fmt.Errorf("lineTypeは'row'または'col'である必要があります")
	}
}

//Matches内に保存されている行、列を更新する
func (gb *GameBoard) UpdateLines(matches []Match) error {
	for _, match := range matches {
		err := gb.UpdateLine(match.LineType, match.Index)
		if err != nil {
			return fmt.Errorf("更新中にエラーが発生しました (%s %d): %w", match.LineType, match.Index, err)
		}
	}
	return nil
}
