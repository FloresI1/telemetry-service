package main

import (
	"forinter/database"
	"forinter/handler"
	"log"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {
	// Инициализация базы данных
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Инициализация Gin
	r := gin.Default()

	// Обработчик для маршрута /track
	r.POST("/track", handler.TrackHandler(db))

	// Запуск сервера
	log.Printf("Server is running on http://localhost%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
