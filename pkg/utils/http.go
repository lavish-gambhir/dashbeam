package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
)

type SuccessResponse struct {
	Status    string    `json:"status"`
	Data      any       `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type ErrorResponse struct {
	Status    string    `json:"status"`
	Error     string    `json:"error"`
	Code      string    `json:"code,omitempty"`
	Details   any       `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func WriteJSONSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := SuccessResponse{
		Status:    "success",
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	json.NewEncoder(w).Encode(response)
}

func WriteJSONSuccessWithStatus(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SuccessResponse{
		Status:    "success",
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	json.NewEncoder(w).Encode(response)
}

func WriteJSONError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Status:    "error",
		Error:     err.Error(),
		Timestamp: time.Now().UTC(),
	}

	if appErr, ok := err.(*apperr.Error); ok {
		response.Code = string(appErr.Code)
		response.Details = appErr.Details
	}

	json.NewEncoder(w).Encode(response)
}

func WriteJSONErrorWithDetails(w http.ResponseWriter, err error, details any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Status:    "error",
		Error:     err.Error(),
		Details:   details,
		Timestamp: time.Now().UTC(),
	}

	if appErr, ok := err.(*apperr.Error); ok {
		response.Code = string(appErr.Code)
	}

	json.NewEncoder(w).Encode(response)
}
