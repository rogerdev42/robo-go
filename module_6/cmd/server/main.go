package main

import (
	"fmt"
	"module_6/cmd/server/app"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл в development режиме
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	// Создаем приложение
	application, err := app.New()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer application.Close()

	// Запускаем приложение
	if err := application.Run(); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}
}