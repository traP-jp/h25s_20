package domain

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// FormulaCalculator は安全な数式計算システム
type FormulaCalculator struct{}

// NewFormulaCalculator creates a new formula calculator
func NewFormulaCalculator() *FormulaCalculator {
	return &FormulaCalculator{}
}

// EvaluateFormula は数式を安全に評価し、結果が10かどうかを確認
func (fc *FormulaCalculator) EvaluateFormula(expression string) (float64, error) {
	// 入力をサニタイズ
	expression = strings.ReplaceAll(expression, " ", "")

	// 基本的な文字チェック
	if !regexp.MustCompile(`^[0-9+\-*/()]*$`).MatchString(expression) {
		return 0, fmt.Errorf("式に無効な文字が含まれています")
	}

	// 数字の数をチェック（1-9の数字が4つ必要）
	numbers := regexp.MustCompile(`[1-9]`).FindAllString(expression, -1)
	if len(numbers) != 4 {
		return 0, fmt.Errorf("数式には1-9の数字が4つ必要です")
	}

	// 逆ポーランド記法に変換
	rpn, err := fc.convertToRPN(expression)
	if err != nil {
		return 0, fmt.Errorf("数式の変換に失敗しました: %w", err)
	}

	// 逆ポーランド記法で計算
	result, err := fc.calculateRPN(rpn)
	if err != nil {
		return 0, fmt.Errorf("計算に失敗しました: %w", err)
	}

	return result, nil
}

// convertToRPN は中置記法を逆ポーランド記法に変換
func (fc *FormulaCalculator) convertToRPN(expression string) ([]string, error) {
	var output []string
	var operators []string

	// 演算子の優先順位
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	isOperator := func(token string) bool {
		_, exists := precedence[token]
		return exists
	}

	i := 0
	for i < len(expression) {
		char := string(expression[i])

		if regexp.MustCompile(`\d`).MatchString(char) {
			// 数字をまとめて読み取る
			num := char
			for i+1 < len(expression) && regexp.MustCompile(`\d`).MatchString(string(expression[i+1])) {
				i++
				num += string(expression[i])
			}
			output = append(output, num)
		} else if char == "(" {
			operators = append(operators, char)
		} else if char == ")" {
			// '(' が出てくるまで演算子を出力
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 || operators[len(operators)-1] != "(" {
				return nil, fmt.Errorf("括弧が不正です")
			}
			operators = operators[:len(operators)-1] // '(' を取り除く
		} else if isOperator(char) {
			for len(operators) > 0 &&
				isOperator(operators[len(operators)-1]) &&
				precedence[operators[len(operators)-1]] >= precedence[char] {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)
		} else {
			return nil, fmt.Errorf("無効なトークン: %s", char)
		}

		i++
	}

	// 残っている演算子をすべて出力
	for len(operators) > 0 {
		op := operators[len(operators)-1]
		if op == "(" || op == ")" {
			return nil, fmt.Errorf("括弧が不正です")
		}
		output = append(output, op)
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// calculateRPN は逆ポーランド記法で計算
func (fc *FormulaCalculator) calculateRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		if regexp.MustCompile(`^\d+$`).MatchString(token) {
			// 数字
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("数字の変換に失敗: %w", err)
			}
			stack = append(stack, num)
		} else {
			// 演算子
			if len(stack) < 2 {
				return 0, fmt.Errorf("不正な逆ポーランド記法")
			}

			second := stack[len(stack)-1]
			first := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = first + second
			case "-":
				result = first - second
			case "*":
				result = first * second
			case "/":
				if second == 0 {
					return 0, fmt.Errorf("ゼロ除算エラー")
				}
				result = first / second
			default:
				return 0, fmt.Errorf("無効な演算子: %s", token)
			}

			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("不正な数式")
	}

	return stack[0], nil
}

// ValidateFormulaNumbers は数式で使用される数字を抽出・検証
func (fc *FormulaCalculator) ValidateFormulaNumbers(expression string) ([]int, error) {
	// 数字を抽出（1-9のみ）
	re := regexp.MustCompile(`[1-9]`)
	numberStrings := re.FindAllString(expression, -1)

	if len(numberStrings) != 4 {
		return nil, fmt.Errorf("数式には1-9の数字が4つ必要です")
	}

	numbers := make([]int, len(numberStrings))
	for i, numStr := range numberStrings {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("数字の変換に失敗しました: %v", err)
		}
		if num < 1 || num > 9 {
			return nil, fmt.Errorf("数字は1-9の範囲である必要があります")
		}
		numbers[i] = num
	}

	return numbers, nil
}

// CheckTarget10 は計算結果が10かどうかを厳密にチェック
func (fc *FormulaCalculator) CheckTarget10(result float64) bool {
	const epsilon = 1e-9
	return math.Abs(result-10) < epsilon
}

// IsValidFormulaPattern は数式パターンが有効かチェック
func (fc *FormulaCalculator) IsValidFormulaPattern(expression string) bool {
	// 基本的な構造チェック
	if len(expression) == 0 {
		return false
	}

	// 括弧のバランスチェック
	balance := 0
	for _, char := range expression {
		if char == '(' {
			balance++
		} else if char == ')' {
			balance--
		}
		if balance < 0 {
			return false
		}
	}

	return balance == 0
}

// GetInvalidCombinations は10が作れない数字の組み合わせを返す
func (fc *FormulaCalculator) GetInvalidCombinations() []string {
	// フロントエンドのinvalid_list.txtから移植
	return []string{
		"1111", "1112", "1113", "1122", "1159", "1169", "1177", "1178", "1179", "1188",
		"1399", "1444", "1499", "1666", "1667", "1677", "1699", "1777", "2257", "3444",
		"3669", "3779", "3999", "4444", "4459", "4477", "4558", "4899", "4999", "5668",
		"5788", "5799", "5899", "6666", "6667", "6677", "6777", "6778", "6888", "6899",
		"6999", "7777", "7788", "7789", "7799", "7888", "7999", "8899",
	}
}

// IsImpossibleCombination は不可能な数字の組み合わせかチェック
func (fc *FormulaCalculator) IsImpossibleCombination(numbers []int) bool {
	if len(numbers) != 4 {
		return false
	}

	// 数字をソートして文字列に変換
	sorted := make([]int, len(numbers))
	copy(sorted, numbers)

	// 简単なソート（バブルソート）
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	combination := fmt.Sprintf("%d%d%d%d", sorted[0], sorted[1], sorted[2], sorted[3])

	invalid := fc.GetInvalidCombinations()
	for _, inv := range invalid {
		if combination == inv {
			return true
		}
	}

	return false
}
