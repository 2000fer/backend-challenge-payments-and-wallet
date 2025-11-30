package api

import (
	"github.com/2000fer/backend-challenge-payments-and-wallet/internal/handlers"
	"github.com/2000fer/backend-challenge-payments-and-wallet/internal/services"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	WalletService := services.NewWalletService()
	PaymentGatewayService := services.NewPaymentGatewayService()

	// Initialize Gin with default middleware
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ping", handlers.Ping)

	// API v1 routes
	apiV1 := r.Group("/api/v1")
	apiV1.POST("/wallets/:user_id/payments", handlers.CreatePayment(PaymentGatewayService))
	apiV1.GET("/wallets/:user_id/balance", handlers.GetBalance(WalletService))
	apiV1.GET("/wallets/:user_id/transactions", handlers.GetTransactions(WalletService))

	return r
}
