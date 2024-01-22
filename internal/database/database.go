package database

import (
	"database/sql"
	"telemetry-service/model"

	_ "github.com/lib/pq"
)

const connectionString = "user=postgres password=1 dbname=telemetrydb sslmode=disable"

// InitDB инициализирует базу данных.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS telemetry (
			id SERIAL PRIMARY KEY,
			user_id INTEGER,
			screen_name VARCHAR(255),
			action_name VARCHAR(255),
			device_id	INTEGER,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InsertTelemetry вставляет телеметрические данные в базу данных.
func InsertTelemetry(db *sql.DB, t model.Telemetry) error {
	_, err := db.Exec(`
		INSERT INTO telemetry (user_id, screen_name, action_name, device_id)
		VALUES ($1, $2, $3, $4)`,
		t.UserID, t.ScreenName, t.ActionName, t.DeviceID)
	return err
}
