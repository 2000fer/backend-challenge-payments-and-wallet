package internal

type PaymentRequest struct {
	UserID uint64  `json:"user_id"`
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
}

type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

var ValidPaymentMethods = []string{
	"account",
	"card",
}

const (
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
)
