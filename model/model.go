package model

// Telemetry представляет структуру данных для телеметрии.
type Telemetry struct {
	UserID   int `json:"user_id"`
	ScreenName string `json:"screen_name"`
	ActionName string `json:"action_name"`
	DeviceID int `json:"device_id"`
}
