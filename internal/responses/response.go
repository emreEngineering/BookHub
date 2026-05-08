package responses

import (
	"encoding/json"
	"net/http"
)

// Cevap formatımız
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty Eğer Data boşsa JSON’a data alanını ekleme.
}

func JSON(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(response)
}

func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	JSON(w, statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}
