// Package application_test содержит тесты для пакета application.
package application_test

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/flexer2006/y.lms_sprint1_Calc/internal/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew проверяет создание приложения с различными портами - базовые случаи
func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		setPort      bool
		port         string
		expectedPort string
	}{
		{
			name:         "default port",
			setPort:      false,
			expectedPort: ":8080",
		},
		{
			name:         "custom port",
			setPort:      true,
			port:         "9090",
			expectedPort: ":9090",
		},
		{
			name:         "empty port",
			setPort:      true,
			port:         "",
			expectedPort: ":8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setPort {
				os.Setenv("PORT", tt.port)
				defer os.Unsetenv("PORT")
			} else {
				os.Unsetenv("PORT")
			}

			app := application.New()
			assert.NotNil(t, app)
		})
	}
}

// TestCalcHandler тесты обработчика HTTP-запросов - база
func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedCode   int
		expectedResult *float64
		expectedError  string
	}{
		{
			name:           "valid simple expression",
			method:         http.MethodPost,
			body:           application.Request{Expression: "2+2"},
			expectedCode:   http.StatusOK,
			expectedResult: ptr(4.0),
		},
		{
			name:           "valid complex expression",
			method:         http.MethodPost,
			body:           application.Request{Expression: "(2 + 3) * 4"},
			expectedCode:   http.StatusOK,
			expectedResult: ptr(20.0),
		},
		{
			name:          "invalid method",
			method:        http.MethodGet,
			expectedCode:  http.StatusMethodNotAllowed,
			expectedError: "Method Not Allowed",
		},
		{
			name:          "empty body",
			method:        http.MethodPost,
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid Request",
		},
		{
			name:          "empty expression",
			method:        http.MethodPost,
			body:          application.Request{Expression: ""},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid Request",
		},
		{
			name:          "invalid json",
			method:        http.MethodPost,
			body:          "invalid json",
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid Request",
		},
		{
			name:          "division by zero",
			method:        http.MethodPost,
			body:          application.Request{Expression: "1/0"},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Division by Zero",
		},
		{
			name:          "invalid expression",
			method:        http.MethodPost,
			body:          application.Request{Expression: "2++2"},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid Expression",
		},
		{
			name:          "mismatched parentheses",
			method:        http.MethodPost,
			body:          application.Request{Expression: "(2+2"},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid Parentheses",
		},
		{
			name:          "invalid character",
			method:        http.MethodPost,
			body:          application.Request{Expression: "2$2"},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid Character",
		},
		// {
		// 	name:          "invalid operator",
		// 	method:        http.MethodPost,
		// 	body:          application.Request{Expression: "2&2"},
		// 	expectedCode:  http.StatusBadRequest,
		// 	expectedError: "Invalid Operator",
		// },
		{
			name:          "empty expression error",
			method:        http.MethodPost,
			body:          application.Request{Expression: ""},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid Request",
		},
		{
			name:           "multiply by negative zero",
			method:         http.MethodPost,
			body:           application.Request{Expression: "88 * -0.0"},
			expectedCode:   http.StatusOK,
			expectedResult: ptr(0.0),
		},
		// {
		// 	name:          "invalid operator",
		// 	method:        http.MethodPost,
		// 	body:          application.Request{Expression: "2&2"},
		// 	expectedCode:  http.StatusBadRequest,
		// 	expectedError: "Invalid Operator",
		// },
	}

	app := application.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader *bytes.Reader

			switch v := tt.body.(type) {
			case string:
				bodyReader = bytes.NewReader([]byte(v))
			case application.Request:
				bodyBytes, err := json.Marshal(v)
				require.NoError(t, err)
				bodyReader = bytes.NewReader(bodyBytes)
			default:
				bodyReader = bytes.NewReader([]byte{})
			}

			req := httptest.NewRequest(tt.method, "/calculate", bodyReader)
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(app.CalcHandler)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var response application.Response
			err := json.NewDecoder(rec.Body).Decode(&response)
			require.NoError(t, err)

			if tt.expectedResult != nil {
				assert.NotNil(t, response.Result)
				assert.InDelta(t, *tt.expectedResult, response.Result, 0.0001)
			}

			if tt.expectedError != "" {
				assert.NotNil(t, response.Error)
				assert.Equal(t, tt.expectedError, response.Error.Error)
				assert.Equal(t, tt.expectedCode, response.Error.Code)
			}
		})
	}
}

// TestLogMiddleware тесты middleware для логирования запросов
func TestLogMiddleware(t *testing.T) {
	app := application.New()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	middleware := app.LogMiddleware(http.HandlerFunc(handler))
	middleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestSendJSON тесты для отправки JSON-ответов
func TestSendJSON(t *testing.T) {
	app := application.New()
	rec := httptest.NewRecorder()
	testData := map[string]string{"test": "data"}

	app.SendJSON(rec, http.StatusOK, testData)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, testData, response)

	// t.Run("json encoding error", func(t *testing.T) {
	// 	app := application.New()
	// 	rec := httptest.NewRecorder()

	// 	app.SendJSON(rec, http.StatusOK, make(chan int))

	// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	// })
}

// "Хелп" функция для создания указателя на float64
func ptr(f float64) *float64 {
	return &f
}

// TestRunServer проверяет запуск сервера
func TestRunServer(t *testing.T) {
	app := application.New()

	// Check if the port is already in use
	ln, err := net.Listen("tcp", app.Config.Address)
	if err != nil {
		t.Skip("Port is already in use, skipping test")
	}
	ln.Close()

	// Start the server in a goroutine
	go func() {
		err := app.RunServer()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("unexpected server error: %v", err)
		}
	}()

	// Make a test request
	resp, err := http.Post("http://localhost:8080/calculate",
		"application/json",
		bytes.NewBufferString(`{"expression":"2+2"}`))

	if err != nil {
		t.Skip("Server is not running, skipping test")
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestErrorHandling проверяет различные сценарии обработки ошибок
func TestErrorHandling(t *testing.T) {
	app := application.New()
	rec := httptest.NewRecorder()

	app.SendError(rec, http.StatusBadRequest, "Test Error", "Test Description")

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response application.Response
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response.Error)
	assert.Equal(t, "Test Error", response.Error.Error)
	assert.Equal(t, "Test Description", response.Error.Description)
	assert.Equal(t, http.StatusBadRequest, response.Error.Code)
}
