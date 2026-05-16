package utils

import "net/http"

// AppError represents a domain/HTTP error with a stable code for clients.
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Common application errors (extend as business rules are implemented).
var (
	ErrInvalidMatrix = &AppError{
		Code:       "INVALID_MATRIX",
		Message:    "matrix payload is invalid",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrQRDecomposition = &AppError{
		Code:       "QR_DECOMPOSITION_FAILED",
		Message:    "QR decomposition could not be performed",
		HTTPStatus: http.StatusUnprocessableEntity,
	}
	ErrDownstreamAPI = &AppError{
		Code:       "DOWNSTREAM_UNAVAILABLE",
		Message:    "statistics service is unavailable",
		HTTPStatus: http.StatusBadGateway,
	}
	ErrInvalidCredentials = &AppError{
		Code:       "INVALID_CREDENTIALS",
		Message:    "invalid username or password",
		HTTPStatus: http.StatusUnauthorized,
	}
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "missing or invalid authorization token",
		HTTPStatus: http.StatusUnauthorized,
	}
)

// ErrorResponse is the standard JSON error envelope for REST clients.
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   *AppError `json:"error"`
}
