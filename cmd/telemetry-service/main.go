package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"telemetry-service/config"
	"telemetry-service/internal/database"
	"telemetry-service/internal/handler"
	"time"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Error loading config:", err)
		return
	}
	db, err := database.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	r := gin.Default()
	r.POST("/track", handler.TrackHandler(db))
	
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	log.Printf("Server is starting and listening on %s...\n", port)

	// Запуск сервера
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server:", err)
		}
	}()

	// Ожидание сигнала завершения работы программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.Close(); err != nil {
		log.Println("Error closing database:", err)
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown error:", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}