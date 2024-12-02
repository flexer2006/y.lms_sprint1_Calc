package calculation

import (
	"strconv"
	"unicode"
)

type operator struct {
	precedence int                                 // Приоритет операции (1 для сложения и вычитания, 2 для умножения и деления)
	operation  func(a, b float64) (float64, error) // Операция
}

// operators определяет поддерживаемые математические операции калькулятора.
var operators = map[rune]operator{
	'+': {1, func(a, b float64) (float64, error) { return a + b, nil }},
	'-': {1, func(a, b float64) (float64, error) { return a - b, nil }},
	'*': {2, func(a, b float64) (float64, error) { return a * b, nil }},
	'/': {2, func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	}},
}

// Calculator хранит состояние вычислений
type Calculator struct {
	numbers    []float64 // Список чисел
	operations []rune    // Список операций
}

// NewCalculator создает новый экземпляр калькулятора
func NewCalculator() *Calculator {
	return &Calculator{
		numbers:    make([]float64, 0, 8),
		operations: make([]rune, 0, 8),
	}
}

// applyOperation применяет операцию к двум последними числами в стеке
func (c *Calculator) applyOperation() error {
	if len(c.numbers) < 2 || len(c.operations) == 0 {
		return ErrInvalidExpression
	}

	b, a := c.numbers[len(c.numbers)-1], c.numbers[len(c.numbers)-2]
	c.numbers = c.numbers[:len(c.numbers)-2]

	op := c.operations[len(c.operations)-1]
	c.operations = c.operations[:len(c.operations)-1]

	operator, exists := operators[op]
	if !exists {
		return ErrInvalidOperator
	}

	result, err := operator.operation(a, b)
	if err != nil {
		return err
	}

	c.numbers = append(c.numbers, result)
	return nil
}

// parseNumber читает число из строки и добавляет его в стек чисел
func (c *Calculator) parseNumber(expression string, startIndex int) (int, error) {
	endIndex := startIndex
	for endIndex < len(expression) && (unicode.IsDigit(rune(expression[endIndex])) || rune(expression[endIndex]) == '.') {
		endIndex++
	}

	number, err := strconv.ParseFloat(expression[startIndex:endIndex], 64)
	if err != nil {
		return endIndex, ErrInvalidExpression
	}

	c.numbers = append(c.numbers, number)
	return endIndex - 1, nil
}

// Calc вычисляет значение математического выражения
func Calc(expression string) (float64, error) {
	calc := NewCalculator()
	expectNumber := true

	for i := 0; i < len(expression); i++ {
		currentChar := rune(expression[i])

		switch {
		case unicode.IsDigit(currentChar) || currentChar == '.':
			var err error
			i, err = calc.parseNumber(expression, i)
			if err != nil {
				return 0, err
			}
			expectNumber = false

		case currentChar == '(':
			calc.operations = append(calc.operations, currentChar)
			expectNumber = true

		case currentChar == ')':
			if expectNumber {
				return 0, ErrInvalidExpression
			}
			err := calc.handleClosingParenthesis()
			if err != nil {
				return 0, err
			}
			expectNumber = false

		case currentChar == '-' && expectNumber:
			if i+1 >= len(expression) {
				return 0, ErrInvalidExpression
			}

			calc.numbers = append(calc.numbers, 0)
			calc.operations = append(calc.operations, '-')
			expectNumber = true

		case isOperator(currentChar):
			if expectNumber {
				return 0, ErrInvalidExpression
			}
			err := calc.handleOperator(currentChar)
			if err != nil {
				return 0, err
			}
			expectNumber = true

		case !unicode.IsSpace(currentChar):
			return 0, ErrInvalidCharacter

		case unicode.IsSpace(currentChar):
			continue
		}
	}

	if expectNumber {
		return 0, ErrInvalidExpression
	}

	for _, op := range calc.operations {
		if op == '(' {
			return 0, ErrMismatchedParens
		}
	}
	// Возвращает окончательный результат
	return calc.calculateFinal()
}

// handleClosingParenthesis обрабатывает закрывающую скобку
func (c *Calculator) handleClosingParenthesis() error {
	for len(c.operations) > 0 && c.operations[len(c.operations)-1] != '(' {
		if err := c.applyOperation(); err != nil {
			return err
		}
	}

	if len(c.operations) == 0 {
		return ErrMismatchedParens
	}

	c.operations = c.operations[:len(c.operations)-1]
	return nil
}

// handleOperator обрабатывает оператор
func (c *Calculator) handleOperator(op rune) error {
	currentOp := operators[op]
	for len(c.operations) > 0 {
		lastOp := c.operations[len(c.operations)-1]
		if lastOp == '(' || operators[lastOp].precedence < currentOp.precedence {
			break
		}
		if err := c.applyOperation(); err != nil {
			return err
		}
	}
	c.operations = append(c.operations, op)
	return nil
}

// calculateFinal вычисляет окончательный результат
func (c *Calculator) calculateFinal() (float64, error) {
	for len(c.operations) > 0 {
		if err := c.applyOperation(); err != nil {
			return 0, err
		}
	}

	if len(c.numbers) != 1 {
		return 0, ErrInvalidExpression
	}

	return c.numbers[0], nil
}

// isOperator проверяет, является ли символ оператором
func isOperator(ch rune) bool {
	_, exists := operators[ch]
	return exists
}
