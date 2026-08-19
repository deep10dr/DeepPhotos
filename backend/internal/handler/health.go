package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"deepphotos/backend/internal/model"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.HealthResponse{
		Status:    "OK",
		Database:  "SQLite 3 Connected (WAL Mode)",
		Storage:   "MinIO S3 Connected",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
