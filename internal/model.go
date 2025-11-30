package internal

import "time"

type PaymentRequest struct {
	UserID uint64  `json:"user_id"`
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
}

type Transaction struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var ValidPaymentMethods = []string{
	"account",
	"card",
}

const (
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
)
