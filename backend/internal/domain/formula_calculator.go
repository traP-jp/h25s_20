package domain

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// FormulaCalculator は逆ポーランド記法専用の安全な数式計算システム
type FormulaCalculator struct {
	// 正規表現を事前にコンパイルしてパフォーマンス向上
	validCharsRe *regexp.Regexp
	digitRe      *regexp.Regexp
	operatorRe   *regexp.Regexp
}

// NewFormulaCalculator creates a new RPN-only formula calculator
func NewFormulaCalculator() *FormulaCalculator {
	return &FormulaCalculator{
		validCharsRe: regexp.MustCompile(`^[123456789+\-*/]*$`),
		digitRe:      regexp.MustCompile(`[1-9]`),
		operatorRe:   regexp.MustCompile(`^[+\-*/]$`),
	}
}

// EvaluateFormula は逆ポーランド記法の数式を安全に評価し、結果が10かどうかを確認
// solvePoland.tsと同じ仕様で実装
func (fc *FormulaCalculator) EvaluateFormula(expression string) (float64, error) {
	// 入力をサニタイズ（空白除去）
	expression = strings.ReplaceAll(expression, " ", "")

	// solvePoland.tsと同じバリデーション
	// 1. 長さチェック（7文字固定）
	if len(expression) != 7 {
		return 0, fmt.Errorf("Invalid input")
	}

	// 2. 使用可能文字チェック（1-9と四則演算子のみ）
	if !fc.validCharsRe.MatchString(expression) {
		return 0, fmt.Errorf("Invalid input")
	}

	// 3. 数字の数をチェック（4つ必要）
	numbers := fc.digitRe.FindAllString(expression, -1)
	if len(numbers) != 4 {
		return 0, fmt.Errorf("Invalid input")
	}

	// 4. RPNパターンチェック（solvePoland.tsと同じ）
	if !fc.isValidRPNPattern(expression) {
		return 0, fmt.Errorf("Invalid input")
	}

	// 5. RPN計算実行
	result, err := fc.calculateRPN(expression)
	if err != nil {
		return 0, fmt.Errorf("Invalid input")
	}

	return result, nil
}

// isValidRPNPattern はRPNパターンが有効かチェック（solvePoland.tsと同じ）
func (fc *FormulaCalculator) isValidRPNPattern(expression string) bool {
	// 文字を 'x'（数字）と 'o'（演算子）にマッピング
	pattern := ""
	for _, char := range expression {
		if fc.digitRe.MatchString(string(char)) {
			pattern += "x"
		} else if fc.operatorRe.MatchString(string(char)) {
			pattern += "o"
		} else {
			return false
		}
	}

	// 有効なRPNパターン（solvePoland.tsと同じ）
	validPatterns := []string{
		"xxxxooo", // 1234+*-
		"xxxoxoo", // 123+4*-
		"xxxooxo", // 123++4-
		"xxoxxoo", // 12+34*-
		"xxoxoxo", // 12+3+4-
	}

	for _, validPattern := range validPatterns {
		if pattern == validPattern {
			return true
		}
	}

	return false
}

// calculateRPN は逆ポーランド記法で計算（solvePoland.tsのcalc_polandと同じ）
func (fc *FormulaCalculator) calculateRPN(expression string) (float64, error) {
	var stack []float64

	for _, char := range expression {
		charStr := string(char)

		if fc.digitRe.MatchString(charStr) {
			// 数字をスタックにプッシュ
			num, err := strconv.ParseFloat(charStr, 64)
			if err != nil {
				return 0, fmt.Errorf("数字の変換に失敗: %w", err)
			}
			stack = append(stack, num)
		} else if fc.operatorRe.MatchString(charStr) {
			// 演算子処理
			if len(stack) < 2 {
				return 0, fmt.Errorf("Invalid RPN")
			}

			// スタックから2つの値を取得（順序注意：solvePoland.tsと同じ）
			second := stack[len(stack)-1]
			first := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch charStr {
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
				return 0, fmt.Errorf("無効な演算子: %s", charStr)
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
	numberStrings := fc.digitRe.FindAllString(expression, -1)

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

// CheckResultType は結果の種類を判定（solvePoland.tsと同じロジック）
func (fc *FormulaCalculator) CheckResultType(result float64) string {
	resultInt := int(math.Round(result))
	if math.Abs(result-float64(resultInt)) < 1e-9 {
		if resultInt == 10 {
			return "10"
		} else {
			return "Not 10"
		}
	} else {
		return "Not an integer"
	}
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

	// 簡単なソート（バブルソート）
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
