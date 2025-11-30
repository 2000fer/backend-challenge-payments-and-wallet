package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal/api"
	"github.com/2000fer/backend-challenge-payments-and-wallet/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Set Gin to release mode in production
	gin.SetMode(cfg.GinMode)

	// Initialize Gin router
	r := api.Init(cfg)

	// Server configuration
	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		slog.InfoContext(context.Background(), fmt.Sprintf("Server starting on port %s", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(context.Background(), "Could not start server", "error", err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.InfoContext(context.Background(), "Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Server forced to shutdown", "error", err.Error())
	}

	slog.InfoContext(context.Background(), "Server exiting")
}
