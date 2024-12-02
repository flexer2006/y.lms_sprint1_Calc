package calculation_test

import (
	"testing"

	"github.com/flexer2006/y.lms_sprint1_Calc/pkg/calculation"

	"github.com/stretchr/testify/assert"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
		err      error
	}{
		// Базовые операции
		{"simple addition", "2 + 2", 4, false, nil},
		{"simple subtraction", "5 - 3", 2, false, nil},
		{"simple multiplication", "4 * 2", 8, false, nil},
		{"simple division", "8 / 2", 4, false, nil},

		// Сложные выражения
		{"complex expression", "2 + 2 * 2", 6, false, nil},
		{"parentheses", "(2 + 2) * 2", 8, false, nil},
		{"nested parentheses", "((2 + 3) * 2) + 1", 11, false, nil},
		{"multiple operations", "1 + 2 + 3 + 4", 10, false, nil},

		// Десятичные числа
		{"decimal numbers", "1.5 + 2.5", 4.0, false, nil},
		{"complex decimals", "1.5 * 2.5 + 3.5", 7.25, false, nil},

		// Отрицательные результаты
		{"negative result", "2 - 5", -3, false, nil},
		{"negative in parentheses", "2 * (5 - 8)", -6, false, nil},

		// Обработка пробелов
		{"spaces handling", "  2  +  2  ", 4, false, nil},
		{"no spaces", "2+2", 4, false, nil},

		// Ошибочные случаи
		{"division by zero", "5 / 0", 0, true, calculation.ErrDivisionByZero},
		{"invalid expression", "2 + ", 0, true, calculation.ErrInvalidExpression},
		{"invalid character", "2 $ 2", 0, true, calculation.ErrInvalidCharacter},
		{"mismatched parentheses", "(2 + 2", 0, true, calculation.ErrMismatchedParens},
		{"empty expression", "", 0, true, calculation.ErrInvalidExpression},
		{"double operators", "2 ++ 2", 0, true, calculation.ErrInvalidExpression},
		{"invalid number format", "2.2.2 + 1", 0, true, calculation.ErrInvalidExpression},

		// Дополнительные тесты на скобки
		{"empty parentheses", "()", 0, true, calculation.ErrInvalidExpression},
		{"missing opening parenthesis", "1 + 2)", 0, true, calculation.ErrMismatchedParens},
		{"missing closing parenthesis", "(1 + 2", 0, true, calculation.ErrMismatchedParens},
		{"multiple missing parentheses", "((1 + 2)", 0, true, calculation.ErrMismatchedParens},

		// Дополнительные тесты на пробелы
		{"tabs and spaces", "\t1 \t+\t 2\t", 3, false, nil},
		{"multiple spaces", "1     +     2", 3, false, nil},

		// Тесты на последовательные операции
		{"consecutive operations", "2 * 3 * 4", 24, false, nil},
		{"mixed operations", "2 * 3 + 4 * 5", 26, false, nil},

		// Тесты на десятичные числа
		{"leading zero decimal", "0.5 + 0.3", 0.8, false, nil},
		{"multiple decimals", "1.5 * 2.5 * 3.5", 13.125, false, nil},

		// Тесты на унарный минус
		{"simple unary minus", "-5", -5, false, nil},
		{"unary minus with parentheses", "-(2 + 3)", -5, false, nil},
		{"multiple unary minus", "--5", 5, false, nil},
		{"unary minus in expression", "2 * (-3)", -6, false, nil},
		{"unary minus with spaces", "- 5", -5, false, nil},
		{"complex unary minus", "(-2) * (-3)", 6, false, nil},
		{"unary minus with decimals", "-2.5", -2.5, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculation.Calc(tt.input)

			if tt.hasError {
				assert.Error(t, err)
				if tt.err != nil {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 0.0001)
			}
		})
	}
}

func TestCalculator_Complex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{
			name:     "complex nested parentheses",
			input:    "((1 + 2) * (3 + 4)) / 2",
			expected: 10.5,
			hasError: false,
		},
		{
			name:     "multiple operations with precedence",
			input:    "1 + 2 * 3 - 4 / 2",
			expected: 5,
			hasError: false,
		},
		{
			name:     "decimal arithmetic",
			input:    "1.5 * 2.5 + 3.75 / 1.25",
			expected: 6.75,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculation.Calc(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 0.0001)
			}
		})
	}
}

func TestCalculator_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
		err      error
	}{
		{
			name:     "multiple decimal points",
			input:    "1.2.3 + 4",
			hasError: true,
			err:      calculation.ErrInvalidExpression,
		},
		{
			name:     "only parentheses",
			input:    "()",
			hasError: true,
			err:      calculation.ErrInvalidExpression,
		},
		{
			name:     "multiple closing parentheses",
			input:    "(1 + 2))",
			hasError: true,
			err:      calculation.ErrMismatchedParens,
		},
		{
			name:     "only operators",
			input:    "+++",
			hasError: true,
			err:      calculation.ErrInvalidExpression,
		},
		{
			name:     "invalid characters between numbers",
			input:    "1 @ 2",
			hasError: true,
			err:      calculation.ErrInvalidCharacter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calculation.Calc(tt.input)
			assert.Error(t, err)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			}
		})
	}
}

func TestCalculator_LargeNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{
			name:     "large multiplication",
			input:    "999999 * 999999",
			expected: 999998000001,
			hasError: false,
		},
		{
			name:     "precise division",
			input:    "1000000 / 3",
			expected: 333333.3333333333,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculation.Calc(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 0.0001)
			}
		})
	}
}

func TestCalculator_Stress(t *testing.T) {
	// Тест на длинное выражение
	longExpr := "1"
	expected := 1.0
	for i := 0; i < 100; i++ {
		longExpr += " + 1"
		expected += 1.0
	}

	result, err := calculation.Calc(longExpr)
	assert.NoError(t, err)
	assert.InDelta(t, expected, result, 0.0001)
}
