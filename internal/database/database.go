package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"telemetry-service/model"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const telemetryTableName = "telemetry"

// InitDB инициализирует базу данных.
func InitDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Проверка наличия всех обязательных переменных окружения
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, errors.New("missing required environment variables")
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
	log.Println("Not pinged")
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + telemetryTableName + ` (
			id SERIAL PRIMARY KEY,
			user_id INTEGER,
			screen_name VARCHAR(255),
			action_name VARCHAR(255),
			device_id INTEGER,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		log.Fatalf("error creating %s table: %v", telemetryTableName, err)
	}

	log.Println("Database initialized successfully")

	return db, nil
}

// InsertTelemetry вставляет телеметрические данные в базу данных.
func InsertTelemetry(ctx context.Context, db *sql.DB, t model.Telemetry) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO `+telemetryTableName+` (user_id, screen_name, action_name, device_id)
		VALUES ($1, $2, $3, $4)`,
		t.UserID, t.ScreenName, t.ActionName, t.DeviceID)
	if err != nil {
		return fmt.Errorf("error inserting telemetry data: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Println("Telemetry data inserted successfully")
	return nil
}
