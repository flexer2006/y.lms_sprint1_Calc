package application

import "errors"

var (
	// Пустой запрос
	ErrEmptyExpression = errors.New("expression is required")
	// Неправильный JSON формат
	ErrInvalidJSON = errors.New("invalid JSON format")
)

type ErrorResponse struct {
	Error       string `json:"error"`                 // Сообщение об ошибке
	Code        int    `json:"code"`                  // Код ошибки
	Description string `json:"description,omitempty"` // Описание ошибки (необязательное поле)
}
