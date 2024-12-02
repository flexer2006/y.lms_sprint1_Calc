package main

import (
	"log"

	"github.com/flexer2006/y.lms_sprint1_Calc/internal/application"
)

// старт сервера, иначе выбрасываем ошибку
func main() {
	app := application.New()
	if err := app.RunServer(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
