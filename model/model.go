package model

import "time"

// Telemetry представляет структуру данных для телеметрии.
type Telemetry struct {
	UserID     int    `json:"user_id"`
	ScreenName string `json:"screen_name"`
	ActionName string `json:"action_name"`
}

// Timestamp возвращает текущее время.
func (t *Telemetry) Timestamp() time.Time {
	return time.Now()
}
