package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/interseguros/challenge/api-go/internal/config"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

// AuthController handles authentication endpoints.
type AuthController struct {
	cfg    *config.Config
	logger utils.Logger
}

// NewAuthController creates an AuthController with the given config.
func NewAuthController(cfg *config.Config, logger utils.Logger) *AuthController {
	return &AuthController{cfg: cfg, logger: logger}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login handles POST /api/auth/login.
// Validates credentials from env vars and returns a signed JWT on success.
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req loginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrInvalidCredentials
	}

	if req.Username != c.cfg.AuthUsername || req.Password != c.cfg.AuthPassword {
		c.logger.Warn("failed login attempt", "username", req.Username)
		return utils.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(c.cfg.JWTSecret))
	if err != nil {
		c.logger.Error("token signing failed", "error", err)
		return fiber.ErrInternalServerError
	}

	c.logger.Info("login successful", "username", req.Username)
	return ctx.JSON(fiber.Map{"token": signed})
}
