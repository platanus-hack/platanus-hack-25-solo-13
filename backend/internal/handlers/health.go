package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/platanus-hack-25/lumera_app/internal/db"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Database  string    `json:"database"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API and database connection
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse "API health status"
// @Router /api/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Database:  "disconnected",
		Timestamp: time.Now(),
	}

	// Check database connection
	if err := db.Ping(); err == nil {
		response.Database = "connected"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
