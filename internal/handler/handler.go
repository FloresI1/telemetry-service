package handler

import (
	"database/sql"
	"telemetry-service/internal/database"
	"telemetry-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TrackHandler обрабатывает запросы на /track.
func TrackHandler(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var t model.Telemetry
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := database.InsertTelemetry(db, t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Telemetry received successfully"})
	}
}
