package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

type PaymentGatewayService interface {
	CreatePayment(paymentRequest internal.PaymentRequest) (string, error)
}

type CreatePaymentRequest struct {
	UserID uint64  `json:"user_id"`
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
}

type CreatePaymentResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id,omitempty"`
	Error         string `json:"error,omitempty"`
}

func CreatePayment(paymentsService PaymentGatewayService) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestParams, err := extractRequestParams(c)
		if err != nil {
			handleCreatePaymentError(c, err)
			return
		}

		paymentRequest := internal.PaymentRequest{
			UserID: requestParams.UserID,
			Method: requestParams.Method,
			Amount: requestParams.Amount,
		}
		transactionID, err := paymentsService.CreatePayment(paymentRequest)
		if err != nil {
			handleCreatePaymentError(c, err)
			return
		}

		c.JSON(http.StatusOK, CreatePaymentResponse{
			Status:        internal.PaymentStatusSuccess,
			TransactionID: transactionID,
		})
	}
}

func handleCreatePaymentError(c *gin.Context, err error) {
	slog.ErrorContext(c.Request.Context(), err.Error())
	errorStatusCode := http.StatusInternalServerError
	// TODO: for each error type send it to telemetry service

	if errors.Is(err, ErrInvalidRequest) {
		errorStatusCode = http.StatusBadRequest
	}

	c.JSON(errorStatusCode, CreatePaymentResponse{
		Status: internal.PaymentStatusFailed,
		Error:  err.Error(),
	})
}

func extractRequestParams(c *gin.Context) (CreatePaymentRequest, error) {
	var requestParams CreatePaymentRequest
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		return CreatePaymentRequest{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	userID := c.Param("user_id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return CreatePaymentRequest{}, fmt.Errorf("%w: invalid user id %s", ErrInvalidRequest, userID)
	}
	requestParams.UserID = uint64(userIDInt)

	if !slices.Contains(internal.ValidPaymentMethods, requestParams.Method) {
		return CreatePaymentRequest{}, fmt.Errorf("%w: invalid payment method %s", ErrInvalidRequest, requestParams.Method)
	}

	if requestParams.Amount <= 0 {
		return CreatePaymentRequest{}, fmt.Errorf("%w: invalid amount %.2f", ErrInvalidRequest, requestParams.Amount)
	}
	return requestParams, nil
}
