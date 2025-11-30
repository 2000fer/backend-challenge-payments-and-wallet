package routes

import (
	"github.com/2000fer/backend-challenge-payments-and-wallet/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	// Initialize Gin with default middleware
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ping", handlers.Ping)

	// API v1 routes
	_ = r.Group("/api/v1")

	return r
}
