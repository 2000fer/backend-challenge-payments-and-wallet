package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrGettingBalance = errors.New("failed to get balance")
)

type WalletService interface {
	GetBalance(userID string) (float64, error)
}

type GetBalanceResponse struct {
	Balance float64 `json:"balance"`
	Error   string  `json:"error,omitempty"`
}

func GetBalance(walletService WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		_, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			err = fmt.Errorf("%w: invalid user_id: %w", ErrInvalidRequest, err)
			handleGetBalanceError(c, err)
			return
		}

		balance, err := walletService.GetBalance(userID)
		if err != nil {
			err = fmt.Errorf("%w: %w", ErrGettingBalance, err)
			handleGetBalanceError(c, err)
			return
		}

		c.JSON(http.StatusOK, GetBalanceResponse{
			Balance: balance,
		})
	}
}

func handleGetBalanceError(c *gin.Context, err error) {
	slog.ErrorContext(c.Request.Context(), err.Error())
	errorStatusCode := http.StatusInternalServerError
	// TODO: for each error type send it to telemetry service
	// telemetry.Incr("metric_name", "error_type", errType)

	if errors.Is(err, ErrInvalidRequest) {
		errorStatusCode = http.StatusBadRequest
	}

	if errors.Is(err, ErrGettingBalance) {
		errorStatusCode = http.StatusInternalServerError
	}

	c.JSON(errorStatusCode, GetBalanceResponse{
		Error: err.Error(),
	})
}
