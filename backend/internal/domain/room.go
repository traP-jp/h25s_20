package domain

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/Knetic/govaluate"
)

// RoomState represents the game state of a room
type RoomState int

const (
	StateWaitingForPlayers RoomState = iota // 募集中 (全員のready待ち)
	StateAllReady                           // 全員READY済 (ユーザーの列の先頭のユーザーがstartを押した)
	StateCountdown                          // START(カウントダウン中)
	StateGameInProgress                     // 実際にゲーム開始(盤面情報配信)
	StateGameEnded                          // ゲーム終了(全員結果表示を閉じるのまち)
)

// String returns the string representation of RoomState
func (rs RoomState) String() string {
	switch rs {
	case StateWaitingForPlayers:
		return "WaitingForPlayers"
	case StateAllReady:
		return "AllReady"
	case StateCountdown:
		return "Countdown"
	case StateGameInProgress:
		return "GameInProgress"
	case StateGameEnded:
		return "GameEnded"
	default:
		return "Unknown"
	}
}

type Room struct {
	ID         int
	Name       string
	GameBoards []GameBoard
	IsOpened   bool
	Players    []Player
	ResultLog  []Result
	State      RoomState // ステートマシンの現在の状態
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
type Matches struct {
	Linetype string //"row"または"col"
	Index    int
}

// 新規盤面の作成
func NewBoard() GameBoard {
	size := 4
	gb := &GameBoard{
		Version: 1,
		Board:   make([][]int, size),
		Size:    size,
	}
	//盤面の初期化
	for i := range gb.Board {
		gb.Board[i] = make([]int, size)
	}
	//盤面全体を1から9のランダムな整数で埋める
	for i := 0; i < size; i++ {
		gb.PopulateRow(i)
	}
	return *gb
}

func AttemptMove(gb *GameBoard, expression string) (bool, string) {

	matches, _ := FindAllMatchingLines(gb, expression)
	if len(matches) == 0 {
		return false, "エラー: その計算式で使える数字の組み合わせは、盤面上に見つかりません。"
	}

	evalResult, err := EvaluateExpression(expression)
	if err != nil {
		return false, "エラー: 計算ができませんでした"
	}

	const epsilon = 1e-9
	if math.Abs(evalResult-10) > epsilon {
		return false, fmt.Sprintf("エラー: 計算結果が10になりません。(結果: %v)", evalResult)
	}

	// 検証をクリアしたら盤面を更新
	gb.UpdateLines(matches)

	// 成功時は true と空のメッセージを返す
	return true, ""
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
func (gb *GameBoard) UpdateLine(linetype string, index int) error {
	if !(0 <= index && index < gb.Size) {
		return fmt.Errorf("インデックスは0から%dの間で指定してください", gb.Size-1)
	}
	switch linetype {
	case "row":
		gb.PopulateRow(index)
		return nil
	case "col":
		gb.PopulateColumn(index)
		return nil
	default:
		return fmt.Errorf("linetypeは'row'または'col'である必要があります")
	}
}

// Matches内に保存されている行、列を更新する
func (gb *GameBoard) UpdateLines(matches []Matches) error {
	for _, match := range matches {
		err := gb.UpdateLine(match.Linetype, match.Index)
		if err != nil {
			return fmt.Errorf("更新中にエラーが発生しました (%s %d): %w", match.Linetype, match.Index, err)
		}
	}

	gb.Version++
	return nil
}

// 入力された数式に含まれる数字が指定された盤面上の行または列に一致するかを判定する
func ValidateExpressionNumbers(expression string, boardLine []int) (bool, error) {
	// 数式からすべての数字を文字列として抽出する
	re := regexp.MustCompile(`\d+`) // 1文字以上の数字にマッチする正規表現
	numStringsInExpr := re.FindAllString(expression, -1)

	if len(numStringsInExpr) == 0 {
		return false, fmt.Errorf("式に数字が含まれていません")
	}
	// 盤面の行にある数字の出現回数を数える
	boardCounts := make(map[int]int)
	for _, num := range boardLine {
		boardCounts[num]++
	}
	// 数式にある数字の出現回数を数える
	exprCounts := make(map[int]int)
	for _, s := range numStringsInExpr {
		num, err := strconv.Atoi(s)
		if err != nil { // 数字の変換に失敗した場合、エラーを返す
			return false, fmt.Errorf("数字の変換に失敗しました: %v", err)
		}
		exprCounts[num]++
	}
	// 数式の数字が、盤面の数字の個数のに一致するかチェック
	for NumInExpr, countInExpr := range exprCounts {
		countInBoard, ok := boardCounts[NumInExpr]
		// 盤面に存在しない数字が数式に含まれている場合、falseを返す
		// また、数式の数字の出現回数が盤面の数字の出現回数にあわない場合もfalseを返す
		if !ok || countInExpr != countInBoard {
			return false, nil
		}
	}
	return true, nil
}

// 盤面すべての行,列についてValidateExpressionNumbersを実行
func FindAllMatchingLines(gb *GameBoard, expression string) ([]Matches, bool) {
	// 見つかったマッチを格納するためのスライスを初期化
	var matches []Matches

	// すべての行をチェック
	for i := 0; i < gb.Size; i++ {
		rowLine, _ := gb.GetLine("row", i)
		isValid, err := ValidateExpressionNumbers(expression, rowLine)
		if err == nil && isValid {
			// 見つかった情報をMatch構造体としてスライスに追加
			matches = append(matches, Matches{Linetype: "row", Index: i})
		}
	}
	// すべての列をチェック
	for i := 0; i < gb.Size; i++ {
		colLine, _ := gb.GetLine("col", i)
		isValid, err := ValidateExpressionNumbers(expression, colLine)

		if err == nil && isValid {
			matches = append(matches, Matches{Linetype: "col", Index: i})
		}
	}
	if len(matches) == 0 {
		return nil, false
	}
	// ループがすべて終わった後、見つかったマッチのリストを返す
	return matches, true
}

// 指定された行または列を取得
func (gb *GameBoard) GetLine(linetype string, index int) ([]int, error) {
	if linetype == "row" {
		if index < 0 || index >= gb.Size {
			return nil, fmt.Errorf("row index out of range")
		}
		return append([]int{}, gb.Board[index]...), nil
	} else if linetype == "col" {
		if index < 0 || index >= gb.Size {
			return nil, fmt.Errorf("column index out of range")
		}
		col := make([]int, gb.Size)
		for i := 0; i < gb.Size; i++ {
			col[i] = gb.Board[i][index]
		}
		return col, nil
	}
	return nil, fmt.Errorf("invalid line type: %s", linetype)
}

// 入力された数式の計算
func EvaluateExpression(expression string) (float64, error) {
	eval, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, fmt.Errorf("無効な数式です: %w", err)
	}

	result, err := eval.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("式の計算に失敗しました: %w", err)
	}

	if val, ok := result.(float64); ok {
		return val, nil
	}
	return 0, fmt.Errorf("計算結果を数値に変換できませんでした")
}

// ステートマシンの制御メソッド

// NewRoom creates a new room with initial state
func NewRoom(id int, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		GameBoards: []GameBoard{NewBoard()},
		IsOpened:   true,
		Players:    []Player{},
		ResultLog:  []Result{},
		State:      StateWaitingForPlayers,
	}
}

// CanTransitionTo checks if the room can transition to the given state
func (r *Room) CanTransitionTo(newState RoomState) bool {
	switch r.State {
	case StateWaitingForPlayers:
		return newState == StateAllReady || newState == StateWaitingForPlayers
	case StateAllReady:
		return newState == StateCountdown || newState == StateWaitingForPlayers
	case StateCountdown:
		return newState == StateGameInProgress || newState == StateWaitingForPlayers
	case StateGameInProgress:
		return newState == StateGameEnded || newState == StateWaitingForPlayers
	case StateGameEnded:
		return newState == StateWaitingForPlayers
	default:
		return false
	}
}

// TransitionTo transitions the room to the given state if valid
func (r *Room) TransitionTo(newState RoomState) error {
	if !r.CanTransitionTo(newState) {
		return fmt.Errorf("cannot transition from %s to %s", r.State.String(), newState.String())
	}
	r.State = newState
	return nil
}

// AreAllPlayersReady checks if all players in the room are ready
func (r *Room) AreAllPlayersReady() bool {
	if len(r.Players) == 0 {
		return false
	}
	for _, player := range r.Players {
		if !player.IsReady {
			return false
		}
	}
	return true
}

// AreAllPlayersClosedResult checks if all players have closed the result display
func (r *Room) AreAllPlayersClosedResult() bool {
	if len(r.Players) == 0 {
		return false
	}
	for _, player := range r.Players {
		if !player.HasClosedResult {
			return false
		}
	}
	return true
}

// GetFirstPlayer returns the first player in the room (for start game permission)
func (r *Room) GetFirstPlayer() *Player {
	if len(r.Players) == 0 {
		return nil
	}
	return &r.Players[0]
}

// CanStartGame checks if the game can be started
func (r *Room) CanStartGame() bool {
	return r.State == StateAllReady && len(r.Players) > 0
}

// StartGame starts the game by transitioning to countdown state
func (r *Room) StartGame() error {
	if !r.CanStartGame() {
		return fmt.Errorf("cannot start game in current state: %s", r.State.String())
	}
	return r.TransitionTo(StateCountdown)
}

// CompleteCountdown transitions from countdown to game in progress
func (r *Room) CompleteCountdown() error {
	if r.State != StateCountdown {
		return fmt.Errorf("room is not in countdown state")
	}
	return r.TransitionTo(StateGameInProgress)
}

// EndGame ends the current game
func (r *Room) EndGame() error {
	if r.State != StateGameInProgress {
		return fmt.Errorf("game is not in progress")
	}
	return r.TransitionTo(StateGameEnded)
}

// ResetRoom resets the room to waiting state
func (r *Room) ResetRoom() error {
	if r.State != StateGameEnded {
		return fmt.Errorf("game has not ended yet")
	}
	// Reset player ready status and result closed status
	for i := range r.Players {
		r.Players[i].IsReady = false
		r.Players[i].HasClosedResult = false
	}
	return r.TransitionTo(StateWaitingForPlayers)
}

// CloseResult closes the result display for a specific player
func (r *Room) CloseResult(playerID int) error {
	if r.State != StateGameEnded {
		return fmt.Errorf("game has not ended yet")
	}

	// Find and update the player
	playerFound := false
	for i, player := range r.Players {
		if player.ID == playerID {
			r.Players[i].HasClosedResult = true
			playerFound = true
			break
		}
	}

	if !playerFound {
		return fmt.Errorf("player with ID %d not found in room", playerID)
	}

	// If all players have closed the result, automatically reset the room
	if r.AreAllPlayersClosedResult() {
		return r.ResetRoom()
	}

	return nil
}
