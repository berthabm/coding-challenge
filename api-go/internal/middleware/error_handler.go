package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

// ErrorHandler maps AppError and unknown errors to consistent JSON responses.
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var appErr *utils.AppError

	if errors.As(err, &appErr) {
		code = appErr.HTTPStatus
		return c.Status(code).JSON(utils.ErrorResponse{
			Success: false,
			Error:   appErr,
		})
	}

	// Fiber body parser and similar may return *fiber.Error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	return c.Status(code).JSON(utils.ErrorResponse{
		Success: false,
		Error: &utils.AppError{
			Code:       "INTERNAL_ERROR",
			Message:    err.Error(),
			HTTPStatus: code,
		},
	})
}
