package domain

import (
	"testing"
)

func TestFormulaCalculator_EvaluateFormula(t *testing.T) {
	calculator := NewFormulaCalculator()

	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedError  bool
		expectedType   string
	}{
		// 正常系 - 10になるケース
		{
			name:           "Simple addition to 10",
			expression:     "1234+++",
			expectedResult: 10,
			expectedError:  false,
			expectedType:   "10",
		},
		{
			name:           "Complex expression to 10",
			expression:     "19+2*+",
			expectedResult: 10,
			expectedError:  false,
			expectedType:   "10",
		},
		// 正常系 - 10以外の整数
		{
			name:           "Not 10 but integer",
			expression:     "12-34*+",
			expectedResult: 11, // (1-2)+(3*4) = -1+12 = 11
			expectedError:  false,
			expectedType:   "Not 10",
		},
		// 正常系 - 整数でない結果
		{
			name:           "Invalid input",
			expression:     "34+5/",
			expectedResult: 1.4, // (3+4)/5 = 7/5 = 1.4
			expectedError:  false,
			expectedType:   "Invalid input",
		},
		// エラーケース - 長さが違う
		{
			name:          "Invalid length - too short",
			expression:    "12+34+",
			expectedError: true,
		},
		{
			name:          "Invalid length - too long",
			expression:    "12+34+5+",
			expectedError: true,
		},
		// エラーケース - 無効な文字
		{
			name:          "Invalid character - contains 0",
			expression:    "10+2+3+",
			expectedError: true,
		},
		{
			name:          "Invalid character - contains parentheses",
			expression:    "(1+2)+3+",
			expectedError: true,
		},
		// エラーケース - 数字の数が間違い
		{
			name:          "Too few numbers",
			expression:    "12+++++",
			expectedError: true,
		},
		{
			name:          "Too many numbers",
			expression:    "12345++",
			expectedError: true,
		},
		// エラーケース - 無効なRPNパターン
		{
			name:          "Invalid RPN pattern",
			expression:    "+123456",
			expectedError: true,
		},
		{
			name:          "Invalid RPN pattern - operators first",
			expression:    "++12345",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.EvaluateFormula(tt.expression)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("Expected result %f, got %f", tt.expectedResult, result)
			}

			// CheckResultType のテスト
			if tt.expectedType != "" {
				resultType := calculator.CheckResultType(result)
				if resultType != tt.expectedType {
					t.Errorf("Expected result type %s, got %s", tt.expectedType, resultType)
				}
			}
		})
	}
}

func TestFormulaCalculator_IsValidRPNPattern(t *testing.T) {
	calculator := NewFormulaCalculator()

	tests := []struct {
		name       string
		expression string
		expected   bool
	}{
		// 有効なパターン
		{"xxxxooo pattern", "1234+++", true},
		{"xxxoxoo pattern", "123+4*+", true},
		{"xxxooxo pattern", "123++4*", true},
		{"xxoxxoo pattern", "12+34*+", true},
		{"xxoxoxo pattern", "12+3+4+", true},
		// 無効なパターン
		{"operators first", "+++1234", false},
		{"mixed invalid", "1+2+3+4", false},
		{"too many operators", "12+++++", false},
		{"too few operators", "1234+", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.isValidRPNPattern(tt.expression)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for expression %s", tt.expected, result, tt.expression)
			}
		})
	}
}

func TestFormulaCalculator_CalculateRPN(t *testing.T) {
	calculator := NewFormulaCalculator()

	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedError  bool
	}{
		{"Simple addition", "12+", 3, false},
		{"Simple subtraction", "12-", -1, false},
		{"Simple multiplication", "12*", 2, false},
		{"Simple division", "12/", 0.5, false},
		{"Complex expression", "12+34*+", 15, false}, // (1+2)+(3*4) = 3+12 = 15
		{"Division by zero", "10/", 0, true},
		{"Insufficient operands", "+12", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.calculateRPN(tt.expression)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("Expected result %f, got %f", tt.expectedResult, result)
			}
		})
	}
}

// solvePoland.tsとの互換性テスト
func TestFormulaCalculator_SolvePolandCompatibility(t *testing.T) {
	calculator := NewFormulaCalculator()

	// solvePoland.tsの例と同じテストケース
	tests := []struct {
		rpnExpression  string
		expectedResult string
	}{
		{"12+34*+", "Not 10"},      // (1+2)+(3*4) = 3+12 = 15 -> "Not 10"
		{"1234+++", "10"},          // 1+2+3+4 = 10 -> "10"
		{"34+5/", "Invalid input"}, // (3+4)/5 = 7/5 = 1.4 -> "Invalid input"
		{"12-34*+", "Not 10"},      // (1-2)+(3*4) = -1+12 = 11 -> "Not 10"
	}

	for _, tt := range tests {
		t.Run(tt.rpnExpression, func(t *testing.T) {
			result, err := calculator.EvaluateFormula(tt.rpnExpression)
			if err != nil {
				if tt.expectedResult == "Invalid input" {
					return // Expected error
				}
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultType := calculator.CheckResultType(result)
			if resultType != tt.expectedResult {
				t.Errorf("Expected %s, got %s for expression %s (result: %f)",
					tt.expectedResult, resultType, tt.rpnExpression, result)
			}
		})
	}
}
