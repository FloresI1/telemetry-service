package main

import (
	"context"
	"forinter/database"
	"forinter/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// Создание http.Server для Gin
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Запуск сервера в отдельной goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Ожидание сигнала завершения работы программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Ждем сигнала завершения

	log.Println("Shutting down server...")

	// Создаем контекст с таймаутом для завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Закрываем сервер
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}