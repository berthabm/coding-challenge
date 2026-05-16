// Package main is the composition root: wires config, logger, routes, and starts Fiber.
// Business logic lives in internal/; this file only bootstraps the application.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/interseguros/challenge/api-go/internal/config"
	"github.com/interseguros/challenge/api-go/internal/routes"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger := utils.NewLogger(cfg.LogLevel)
	app := routes.Setup(cfg, logger)

	go func() {
		addr := cfg.Host + ":" + cfg.Port
		logger.Info("server starting", "addr", addr, "env", cfg.Env)
		if err := app.Listen(addr); err != nil {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down")
	_ = app.Shutdown()
}
