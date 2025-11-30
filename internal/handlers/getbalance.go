package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetBalanceResponse struct {
	Balance float64 `json:"balance"`
	Error   string  `json:"error,omitempty"`
}

func GetBalance(c *gin.Context) {
	userID := c.Param("user_id")
	_, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		err = fmt.Errorf("%w: invalid user_id: %w", ErrInvalidRequest, err)
		handleGetBalanceError(c, err)
		return
	}

	c.JSON(http.StatusOK, GetBalanceResponse{
		Balance: 0,
	})
}

func handleGetBalanceError(c *gin.Context, err error) {
	slog.ErrorContext(c.Request.Context(), err.Error())
	errorStatusCode := http.StatusInternalServerError
	// TODO: for each error type send it to telemetry service
	// telemetry.Incr("metric_name", "error_type", errType)

	if errors.Is(err, ErrInvalidRequest) {
		errorStatusCode = http.StatusBadRequest
	}

	c.JSON(errorStatusCode, GetBalanceResponse{
		Error: err.Error(),
	})
}
