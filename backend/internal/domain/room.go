package domain

import (
	"fmt"
	"math/rand"
	"sort"
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
	// バージョンごとの変更履歴を記録
	ChangeHistory map[int][]Matches // version -> 変更された行/列のリスト
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
	Linetype string //"row", "col", "diagonal_main", "diagonal_anti"
	Index    int
	// 新仕様：該当するマス位置のリスト (row, col)
	Positions []Position
}

// マス位置を表す構造体
type Position struct {
	Row int
	Col int
}

// 新規盤面の作成
func NewBoard() GameBoard {
	size := 4
	gb := &GameBoard{
		Version:       1,
		Board:         make([][]int, size),
		Size:          size,
		ChangeHistory: make(map[int][]Matches),
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

// AttemptMoveWithVersion はバージョンを考慮した細かい衝突検出付きの処理（新仕様）
func AttemptMoveWithVersion(gb *GameBoard, expression string, submittedVersion int) (bool, string, int) {
	matches, found := FindAllMatchingLinesWithSets(gb, expression)
	if !found {
		return false, "エラー: その計算式で使える数字の組み合わせは、盤面上に見つかりません。", 0
	}

	// バージョン衝突チェック（細かい衝突検出）
	hasConflict, conflictMsg := gb.CheckConflictWithPositions(submittedVersion, matches)
	if hasConflict {
		return false, conflictMsg, 0
	}

	// 新しいRPN専用計算システムを使用
	calculator := NewFormulaCalculator()
	evalResult, err := calculator.EvaluateFormula(expression)
	if err != nil {
		return false, fmt.Sprintf("エラー: 無効な数式です (%s)", err.Error()), 0
	}

	// 結果が10かどうかをチェック
	if !calculator.CheckTarget10(evalResult) {
		// solvePoland.tsと同じ形式でより詳細な結果を返す
		resultType := calculator.CheckResultType(evalResult)
		switch resultType {
		case "Not an integer":
			return false, "エラー: 計算結果が整数になりません。", 0
		case "Not 10":
			return false, fmt.Sprintf("エラー: 計算結果が10になりません。(結果: %.0f)", evalResult), 0
		default:
			return false, fmt.Sprintf("エラー: 計算結果が10になりません。(結果: %.6f)", evalResult), 0
		}
	}

	// 検証をクリアしたら盤面を更新（新仕様）
	gb.UpdateLinesWithPositions(matches)

	// 成功時は true と空のメッセージ、マッチ数を返す
	return true, "", len(matches)
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

// UpdateLinesWithPositions 新仕様：複数のマス位置を直接更新
func (gb *GameBoard) UpdateLinesWithPositions(matches []Matches) error {
	// 更新対象のマス位置を重複なしで収集
	positionSet := make(map[string]bool)
	var allPositions []Position

	for _, match := range matches {
		for _, pos := range match.Positions {
			key := fmt.Sprintf("%d_%d", pos.Row, pos.Col)
			if !positionSet[key] {
				positionSet[key] = true
				allPositions = append(allPositions, pos)
			}
		}
	}

	// 収集したマス位置を更新
	for _, pos := range allPositions {
		if pos.Row < 0 || pos.Row >= gb.Size || pos.Col < 0 || pos.Col >= gb.Size {
			return fmt.Errorf("無効なマス位置: (%d, %d)", pos.Row, pos.Col)
		}
		gb.Board[pos.Row][pos.Col] = rand.Intn(9) + 1
	}

	gb.Version++
	// 変更履歴を記録（新仕様用）
	gb.ChangeHistory[gb.Version] = matches
	return nil
}

// CheckConflictWithPositions はマス位置ベースでの衝突検出（新仕様）
func (gb *GameBoard) CheckConflictWithPositions(submittedVersion int, formulaMatches []Matches) (bool, string) {
	// 提出されたバージョンが現在より新しい場合はエラー
	if submittedVersion > gb.Version {
		return true, fmt.Sprintf("無効なバージョンです: 提出バージョン%d > 現在バージョン%d", submittedVersion, gb.Version)
	}

	// 提出されたバージョンが現在と同じ場合は衝突なし
	if submittedVersion == gb.Version {
		return false, ""
	}

	// 提出バージョン以降に変更されたマス位置を収集
	changedPositions := make(map[string]bool)
	for version := submittedVersion + 1; version <= gb.Version; version++ {
		if changes, exists := gb.ChangeHistory[version]; exists {
			for _, change := range changes {
				for _, pos := range change.Positions {
					key := fmt.Sprintf("%d_%d", pos.Row, pos.Col)
					changedPositions[key] = true
				}
			}
		}
	}

	// 数式で使用されるマス位置と変更されたマス位置の重複をチェック
	for _, match := range formulaMatches {
		for _, pos := range match.Positions {
			key := fmt.Sprintf("%d_%d", pos.Row, pos.Col)
			if changedPositions[key] {
				return true, fmt.Sprintf("エラー: マス位置(%d,%d)は他のプレイヤーによって更新されています",
					pos.Row, pos.Col)
			}
		}
	}

	return false, ""
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

// 入力された数式の計算（新しい安全な実装）
func EvaluateExpression(expression string) (float64, error) {
	calculator := NewFormulaCalculator()
	return calculator.EvaluateFormula(expression)
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

// 新仕様：縦横斜め + 2×2ブロックの4つの数字の組を判定（順序順不同）
func FindAllMatchingLinesWithSets(gb *GameBoard, expression string) ([]Matches, bool) {
	// 数式から数字を抽出してソート（順序順不同対応）
	formulaNumbers, err := ExtractAndSortNumbers(expression)
	if err != nil {
		return nil, false
	}

	var matches []Matches

	// 行をチェック (4つの行)
	for i := 0; i < gb.Size; i++ {
		rowNumbers := make([]int, gb.Size)
		positions := make([]Position, gb.Size)
		for j := 0; j < gb.Size; j++ {
			rowNumbers[j] = gb.Board[i][j]
			positions[j] = Position{Row: i, Col: j}
		}

		// 数字をソートして比較
		sortedRow := make([]int, len(rowNumbers))
		copy(sortedRow, rowNumbers)
		sort.Ints(sortedRow)

		if arraysEqual(formulaNumbers, sortedRow) {
			matches = append(matches, Matches{
				Linetype:  "row",
				Index:     i,
				Positions: positions,
			})
		}
	}

	// 列をチェック (4つの列)
	for j := 0; j < gb.Size; j++ {
		colNumbers := make([]int, gb.Size)
		positions := make([]Position, gb.Size)
		for i := 0; i < gb.Size; i++ {
			colNumbers[i] = gb.Board[i][j]
			positions[i] = Position{Row: i, Col: j}
		}

		// 数字をソートして比較
		sortedCol := make([]int, len(colNumbers))
		copy(sortedCol, colNumbers)
		sort.Ints(sortedCol)

		if arraysEqual(formulaNumbers, sortedCol) {
			matches = append(matches, Matches{
				Linetype:  "col",
				Index:     j,
				Positions: positions,
			})
		}
	}

	// 主対角線をチェック（左上から右下）
	mainDiagNumbers := make([]int, gb.Size)
	mainDiagPositions := make([]Position, gb.Size)
	for i := 0; i < gb.Size; i++ {
		mainDiagNumbers[i] = gb.Board[i][i]
		mainDiagPositions[i] = Position{Row: i, Col: i}
	}
	sortedMainDiag := make([]int, len(mainDiagNumbers))
	copy(sortedMainDiag, mainDiagNumbers)
	sort.Ints(sortedMainDiag)

	if arraysEqual(formulaNumbers, sortedMainDiag) {
		matches = append(matches, Matches{
			Linetype:  "diagonal_main",
			Index:     0,
			Positions: mainDiagPositions,
		})
	}

	// 反対角線をチェック（右上から左下）
	antiDiagNumbers := make([]int, gb.Size)
	antiDiagPositions := make([]Position, gb.Size)
	for i := 0; i < gb.Size; i++ {
		antiDiagNumbers[i] = gb.Board[i][gb.Size-1-i]
		antiDiagPositions[i] = Position{Row: i, Col: gb.Size - 1 - i}
	}
	sortedAntiDiag := make([]int, len(antiDiagNumbers))
	copy(sortedAntiDiag, antiDiagNumbers)
	sort.Ints(sortedAntiDiag)

	if arraysEqual(formulaNumbers, sortedAntiDiag) {
		matches = append(matches, Matches{
			Linetype:  "diagonal_anti",
			Index:     0,
			Positions: antiDiagPositions,
		})
	}

	// 2×2ブロックをチェック（数独風の領域判定）
	blocks := []struct {
		name     string
		startRow int
		startCol int
	}{
		{"block_top_left", 0, 0},     // 左上ブロック
		{"block_top_right", 0, 2},    // 右上ブロック
		{"block_bottom_left", 2, 0},  // 左下ブロック
		{"block_bottom_right", 2, 2}, // 右下ブロック
	}

	for blockIndex, block := range blocks {
		blockNumbers := make([]int, 4)
		blockPositions := make([]Position, 4)
		index := 0

		// 2×2のブロック内の数字を収集
		for i := block.startRow; i < block.startRow+2; i++ {
			for j := block.startCol; j < block.startCol+2; j++ {
				blockNumbers[index] = gb.Board[i][j]
				blockPositions[index] = Position{Row: i, Col: j}
				index++
			}
		}

		// 数字をソートして比較
		sortedBlock := make([]int, len(blockNumbers))
		copy(sortedBlock, blockNumbers)
		sort.Ints(sortedBlock)

		if arraysEqual(formulaNumbers, sortedBlock) {
			matches = append(matches, Matches{
				Linetype:  block.name,
				Index:     blockIndex,
				Positions: blockPositions,
			})
		}
	}

	return matches, len(matches) > 0
}

// 数式から数字を抽出してソートする（新しい安全な実装）
func ExtractAndSortNumbers(expression string) ([]int, error) {
	calculator := NewFormulaCalculator()
	numbers, err := calculator.ValidateFormulaNumbers(expression)
	if err != nil {
		return nil, err
	}

	// 数字をソートして返す
	sort.Ints(numbers)
	return numbers, nil
}

// 配列が等しいかチェック
func arraysEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
