package calculation

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero    = errors.New("division by zero")
	ErrInvalidOperator   = errors.New("invalid operator")
	ErrMismatchedParens  = errors.New("mismatched parentheses")
	ErrInvalidCharacter  = errors.New("invalid character")
)
