package calculation

import "errors"

var (
	// Неправильный формат или не вычисляется
	ErrInvalidExpression = errors.New("invalid expression")
	// Деление на ноль
	ErrDivisionByZero = errors.New("division by zero")
	// Неправильный оператор
	ErrInvalidOperator = errors.New("invalid operator")
	// Неправильные скобки
	ErrMismatchedParens = errors.New("mismatched parentheses")
	// Неправильный символ
	ErrInvalidCharacter = errors.New("invalid character")
)
