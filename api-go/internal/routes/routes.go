// Package routes registers HTTP routes and wires the dependency graph.
// This is the only place that should know concrete implementations (DI root).
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/interseguros/challenge/api-go/internal/clients"
	"github.com/interseguros/challenge/api-go/internal/config"
	"github.com/interseguros/challenge/api-go/internal/controllers"
	"github.com/interseguros/challenge/api-go/internal/middleware"
	"github.com/interseguros/challenge/api-go/internal/services"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

// Setup builds the Fiber app with middleware, routes, and injected services.
func Setup(cfg *config.Config, logger utils.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(middleware.RequestLogger(logger))

	// CORS — must be registered before any route
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Dependency injection
	qrSvc := services.NewQRService()
	statsClient := clients.NewStatsClient(cfg, logger)
	matrixSvc := services.NewMatrixService(qrSvc, statsClient, logger)
	matrixCtrl := controllers.NewMatrixController(matrixSvc, logger)
	authCtrl := controllers.NewAuthController(cfg, logger)

	// Public endpoints
	app.Get("/health", matrixCtrl.Health)
	app.Post("/api/auth/login", authCtrl.Login)

	// Protected endpoints — require valid JWT Bearer token
	api := app.Group("/api", middleware.JWTAuth(cfg.JWTSecret))
	api.Post("/qr", matrixCtrl.ProcessQR)

	// Legacy v1 route (same use case, wrapped response)
	v1 := app.Group("/api/v1", middleware.JWTAuth(cfg.JWTSecret))
	matrices := v1.Group("/matrices")
	matrices.Post("/qr", matrixCtrl.DecomposeAndForward)

	return app
}
