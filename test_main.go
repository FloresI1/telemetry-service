package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"telemetry-service/model"
)



func generateRandomTelemetry() model.Telemetry {
	return model.Telemetry{
		UserID:     rand.Intn(1000) + 1,        
		ScreenName: getRandomString(5),          
		ActionName: getRandomString(5),          
		DeviceID:   rand.Intn(10) + 1,           
	}
}

func getRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func sendTelemetry(telemetry model.Telemetry, url string) error {
	jsonData, err := json.Marshal(telemetry)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("HTTP Status Code: %s\n", resp.Status)
	return nil
}

func main() {

	// Количество отправляемых сообщений
	messageCount := 1000000

	for i := 0; i < messageCount; i++ {
		telemetry := generateRandomTelemetry()

		err := sendTelemetry(telemetry, "http://localhost:8080/track")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
