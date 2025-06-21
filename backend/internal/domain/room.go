package domain

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
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

// 回答確認後一致した行と列の情報を保持
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
func (gb *GameBoard) UpdateLine(LineType string, Index int) error {
	if !(0 <= Index && Index < gb.Size) {
		return fmt.Errorf("インデックスは0から%dの間で指定してください", gb.Size-1)
	}
	switch LineType {
	case "row":
		gb.PopulateRow(Index)
		return nil
	case "col":
		gb.PopulateColumn(Index)
		return nil
	default:
		return fmt.Errorf("lineTypeは'row'または'col'である必要があります")
	}
}

// Matches内に保存されている行、列を更新する
func (gb *GameBoard) UpdateLines(matches []Match) error {
	for _, match := range matches {
		err := gb.UpdateLine(match.LineType, match.Index)
		if err != nil {
			return fmt.Errorf("更新中にエラーが発生しました (%s %d): %w", match.LineType, match.Index, err)
		}
	}
	
	gb.Version++
	return nil
}

// 入力された数式に含まれる数字が指定された盤面上の行または列に一致するかを判定する
func ValidateExpressionNumbers(Expression string, BoardLine []int) (bool, error) {
	// 数式からすべての数字を文字列として抽出する
	re := regexp.MustCompile(`\d+`) // 1文字以上の数字にマッチする正規表現
	NumStringsInExpr := re.FindAllString(Expression, -1)

	if len(NumStringsInExpr) == 0 {
		return false, fmt.Errorf("式に数字が含まれていません")
	}
	// 盤面の行にある数字の出現回数を数える
	BoardCounts := make(map[int]int)
	for _, num := range BoardLine {
		BoardCounts[num]++
	}
	// 数式にある数字の出現回数を数える
	ExprCounts := make(map[int]int)
	for _, s := range NumStringsInExpr {
		num, err := strconv.Atoi(s)
		if err != nil { // 数字の変換に失敗した場合、エラーを返す
			return false, fmt.Errorf("数字の変換に失敗しました: %v", err)
		}
		ExprCounts[num]++
	}
	// 数式の数字が、盤面の数字の個数のに一致するかチェック
	for NumInExpr, CountInExpr := range ExprCounts {
		CountInBoard, ok := BoardCounts[NumInExpr]
		// 盤面に存在しない数字が数式に含まれている場合、falseを返す
		// また、数式の数字の出現回数が盤面の数字の出現回数にあわない場合もfalseを返す
		if !ok || CountInExpr != CountInBoard {
			return false, nil
		}
	}
	return true, nil
}

// 盤面すべての行,列についてValidateExpressionNumbersを実行
func FindAllMatchingLines(gb *GameBoard, expression string) ([]Match, bool) {
	// 見つかったマッチを格納するためのスライスを初期化
	var matches []Match

	// --- すべての行をチェック ---
	for i := 0; i < gb.Size; i++ {
		rowLine, _ := gb.GetLine("row", i)
		isValid, err := ValidateExpressionNumbers(expression, rowLine)
		if err == nil && isValid {
			// 見つかった情報をMatch構造体としてスライスに追加
			matches = append(matches, Match{LineType: "row", Index: i})
		}
	}
	// --- すべての列をチェック ---
	for i := 0; i < gb.Size; i++ {
		colLine, _ := gb.GetLine("col", i)
		isValid, err := ValidateExpressionNumbers(expression, colLine)

		if err == nil && isValid {
			matches = append(matches, Match{LineType: "col", Index: i})
		}
	}
	if len(matches) == 0 {
		return nil, false
	}
	// ループがすべて終わった後、見つかったマッチのリストを返す
	return matches, true
}

//指定された行または列を取得
func (gb *GameBoard) GetLine(lineType string, index int) ([]int, error) {
	if lineType == "row" {
		if index < 0 || index >= gb.Size {
			return nil, fmt.Errorf("row index out of range")
		}
		return append([]int{}, gb.Board[index]...), nil
	} else if lineType == "col" {
		if index < 0 || index >= gb.Size {
			return nil, fmt.Errorf("column index out of range")
		}
		col := make([]int, gb.Size)
		for i := 0; i < gb.Size; i++ {
			col[i] = gb.Board[i][index]
		}
		return col, nil
	}
	return nil, fmt.Errorf("invalid line type: %s", lineType)
}
