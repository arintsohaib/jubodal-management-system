package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message"`
	Details   []FieldError      `json:"details,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// FieldError represents a field-specific validation error
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSON writes JSON response
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Success writes successful JSON response
func Success(w http.ResponseWriter, data interface{}, message string) {
	JSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Created writes 201 Created JSON response
func Created(w http.ResponseWriter, data interface{}, message string) {
	JSON(w, http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Error writes error JSON response
func Error(w http.ResponseWriter, statusCode int, errorCode, message string, requestID string) {
	JSON(w, statusCode, ErrorResponse{
		Error:     errorCode,
		Message:   message,
		RequestID: requestID,
	})
}

// BadRequest writes 400 error
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "bad_request", message, "")
}

// Unauthorized writes 401 error
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "unauthorized", message, "")
}

// Forbidden writes 403 error
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, "forbidden", message, "")
}

// NotFound writes 404 error
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "not_found", message, "")
}

// Conflict writes 409 error
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "conflict", message, "")
}

// InternalError writes 500 error
func InternalError(w http.ResponseWriter, message, requestID string) {
	Error(w, http.StatusInternalServerError, "internal_server_error", message, requestID)
}
