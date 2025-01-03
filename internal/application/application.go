// Package application предоставляет HTTP-сервер для вычисления математических выражений.
package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flexer2006/y.lms_sprint1_Calc/pkg/calculation"
)

type Config struct {
	Address string
	Logger  *log.Logger
}

type Application struct {
	Config *Config // Измените с config на Config
	Logger *log.Logger
}
type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64        `json:"result"`
	Error  *ErrorResponse `json:"error,omitempty"`
}

func New() *Application {
	logger := log.New(os.Stdout, "[CALC] ", log.LstdFlags|log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Application{
		Config: &Config{
			Address: fmt.Sprintf(":%s", port),
			Logger:  logger,
		},
		Logger: logger,
	}
}

func (app *Application) LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Printf("Request: Method=%s Path=%s RemoteAddr=%s",
			r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}

func (app *Application) CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.SendError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.SendError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if req.Expression == "" {
		app.SendError(w, http.StatusBadRequest, "Expression is required")
		return
	}

	result, err := calculation.Calc(req.Expression)
	if err != nil {
		app.handleCalculationError(w, err)
		return
	}

	app.Logger.Printf("Calculated result: %f", result)

	app.SendJSON(w, http.StatusOK, Response{
		Result: result,
	})
}

func (app *Application) handleCalculationError(w http.ResponseWriter, err error) {
	switch err {
	case calculation.ErrInvalidExpression:
		app.SendError(w, http.StatusBadRequest, "Expression is not valid")

	case calculation.ErrInvalidCharacter:
		app.SendError(w, http.StatusBadRequest, "Expression is not valid")

	case calculation.ErrMismatchedParens:
		app.SendError(w, http.StatusBadRequest, "Expression is not valid")

	case calculation.ErrDivisionByZero:
		app.SendError(w, http.StatusUnprocessableEntity, "Division by Zero")

	case calculation.ErrInvalidOperator:
		app.SendError(w, http.StatusBadRequest, "Expression is not valid")

	default:
		app.SendError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func (app *Application) SendError(w http.ResponseWriter, code int, message string) {
	app.Logger.Printf("Error: %s (Code: %d)", message, code)

	response := map[string]string{
		"error": message,
	}

	app.SendJSON(w, code, response)
}

func (app *Application) SendJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.Logger.Printf("Error encoding response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *Application) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/calculate", app.LogMiddleware(app.CalcHandler))

	app.Logger.Printf("Starting server on %s", app.Config.Address)
	return http.ListenAndServe(app.Config.Address, mux)
}
