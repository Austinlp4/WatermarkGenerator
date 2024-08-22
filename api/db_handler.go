package api

import (
	"net/http"

	"watermark-generator/db"
)

// type Handler struct {
//     DB *mongo.Client
// }

func (h *WatermarkHandler) TestDBConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to test DB connection
	// For example:
	err := db.TestConnection()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Database connection successful"))
}
