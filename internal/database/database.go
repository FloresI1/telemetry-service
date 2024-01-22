package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"telemetry-service/model"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// InitDB инициализирует базу данных.
func InitDB() (*sql.DB, error) {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Получение значений переменных окружения для подключения к PostgreSQL
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Проверка наличия всех обязательных переменных окружения
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, errors.New("missing required environment variables")
	}

	// Формирование строки подключения
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Проверка соединения с базой данных
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS telemetry (
			id SERIAL PRIMARY KEY,
			user_id INTEGER,
			screen_name VARCHAR(255),
			action_name VARCHAR(255),
			device_id INTEGER,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return nil, fmt.Errorf("error creating telemetry table: %w", err)
	}

	return db, nil
}

// InsertTelemetry вставляет телеметрические данные в базу данных.
func InsertTelemetry(ctx context.Context, db *sql.DB, t model.Telemetry) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO telemetry (user_id, screen_name, action_name, device_id)
		VALUES ($1, $2, $3, $4)`,
		t.UserID, t.ScreenName, t.ActionName, t.DeviceID)
	if err != nil {
		return fmt.Errorf("error inserting telemetry data: %w", err)
	}
	return nil
}