package application

import "errors"

var (
	
	ErrEmptyExpression = errors.New("expression is required")
	
	ErrInvalidJSON = errors.New("invalid JSON format")
)

type ErrorResponse struct {
	Error       string `json:"error"`    
	Code        int    `json:"code"`  
	Description string `json:"description,omitempty"` 
}
// заготовочка
