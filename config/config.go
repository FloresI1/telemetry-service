package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// LoadConfig загружает конфигурацию.
func LoadConfig() error {
	envPath := filepath.Join("..", "telemetry-service", ".env")
	if err := loadEnv(envPath); err != nil {
		return err
	}

	ginMode := os.Getenv("GIN_MODE")

	switch ginMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		return errors.New("invalid GIN_MODE. Use 'release', 'debug', or 'test'")
	}

	return nil
}

func loadEnv(envPath string) error {
	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	return nil
}
