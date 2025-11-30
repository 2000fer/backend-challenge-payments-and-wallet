package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
	"github.com/gin-gonic/gin"
)

type GetTransactionsResponse struct {
	Transactions []internal.Transaction `json:"transactions"`
	Error        string                 `json:"error,omitempty"`
}

func GetTransactions(c *gin.Context) {
	userID := c.Param("user_id")
	_, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		err = fmt.Errorf("%w: invalid user_id: %w", ErrInvalidRequest, err)
		handleGetTransactionsError(c, err)
		return
	}

	c.JSON(http.StatusOK, GetTransactionsResponse{
		Transactions: []internal.Transaction{},
	})
}

func handleGetTransactionsError(c *gin.Context, err error) {
	slog.ErrorContext(c.Request.Context(), err.Error())
	errorStatusCode := http.StatusInternalServerError
	// TODO: for each error type send it to telemetry service
	// telemetry.Incr("metric_name", "error_type", errType)

	if errors.Is(err, ErrInvalidRequest) {
		errorStatusCode = http.StatusBadRequest
	}

	c.JSON(errorStatusCode, GetTransactionsResponse{
		Error: err.Error(),
	})
}
