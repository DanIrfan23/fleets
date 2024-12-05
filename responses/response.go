package responses

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Code    int
	Message string
	Data    T
}

func SuccessResponse[T any](w http.ResponseWriter, statusCode int, message string, data T) {
	response := Response[T]{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := Response[string]{
		Code:    statusCode,
		Message: message,
		Data:    "",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func ErrorValidationResponse[T any](w http.ResponseWriter, statusCode int, message string, data T) {
	response := Response[T]{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
