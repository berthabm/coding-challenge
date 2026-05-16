package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/interseguros/challenge/api-go/internal/models"
	"github.com/interseguros/challenge/api-go/internal/services"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

// MatrixController handles HTTP for matrix/QR endpoints (thin layer).
type MatrixController struct {
	service services.MatrixService
	logger  utils.Logger
}

// NewMatrixController creates a matrix HTTP handler group.
func NewMatrixController(service services.MatrixService, logger utils.Logger) *MatrixController {
	return &MatrixController{service: service, logger: logger}
}

// ProcessQR handles POST /api/qr — returns the assembled QRResponse.
func (c *MatrixController) ProcessQR(ctx *fiber.Ctx) error {
	var req models.MatrixRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn("failed to parse body", "error", err)
		return utils.ErrInvalidMatrix
	}

	c.logger.Info("POST /api/qr received")

	result, err := c.service.ProcessMatrixQR(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

// DecomposeAndForward handles POST /api/v1/matrices/qr (legacy wrapper).
func (c *MatrixController) DecomposeAndForward(ctx *fiber.Ctx) error {
	var req models.MatrixRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrInvalidMatrix
	}

	result, err := c.service.ProcessMatrixQR(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// Health handles GET /health.
func (c *MatrixController) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "ok",
		"service": "api-go",
	})
}
