package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"telemetry-service/internal/database"
	"telemetry-service/internal/handler"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const port = ":8080"

func main() {
	// Попытка загрузить .env из корневой директории
	envPath := filepath.Join("..", "telemetry-service", ".env")
	if err := loadEnv(envPath); err != nil {
		log.Fatal("Error loading .env file:", err)

	}

	// Получение значения переменной окружения GIN_MODE
	ginMode := os.Getenv("GIN_MODE")

	// Установка режима работы Gin в соответствии с переменной окружения
	switch ginMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		log.Fatal("Invalid GIN_MODE. Use 'release', 'debug', or 'test'")
	}

	// Инициализация базы данных
	db, err := database.InitDB()
	if err != nil {
		log.Println("Failed to initialize database:", err)
		return
	}

	defer db.Close()

	log.Println("Database initialized successfully")

	// Инициализация Gin
	r := gin.Default()

	// Обработчик для маршрута /track
	r.POST("/track", handler.TrackHandler(db))

	// Создание http.Server для Gin
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	log.Printf("Server is starting and listening on %s...\n", port)

	// Запуск сервера в отдельной goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server:", err)
		}
	}()

	// Ожидание сигнала завершения работы программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Ждем сигнала завершения

	log.Println("Shutting down server...")

	// Контекст с таймаутом для завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Закрываем сервер
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown error:", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
func loadEnv(envPath string) error {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	return nil
}