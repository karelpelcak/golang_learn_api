package http

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

func SendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	SendJSON(w, response, http.StatusOK)
}

func SendError(w http.ResponseWriter, message string, code int) {
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}
	statusCode := http.StatusInternalServerError
	switch code {
	case 400:
		statusCode = http.StatusBadRequest
	case 404:
		statusCode = http.StatusNotFound
	case 409:
		statusCode = http.StatusConflict
	case 500:
		statusCode = http.StatusInternalServerError
	}
	SendJSON(w, response, statusCode)
}
