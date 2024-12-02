package application

import "errors"

var (
	ErrEmptyExpression = errors.New("expression is required")
	ErrInvalidJSON     = errors.New("invalid JSON format")
)

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Error       string `json:"error"`
	Code        int    `json:"code"`
	Description string `json:"description,omitempty"`
}
