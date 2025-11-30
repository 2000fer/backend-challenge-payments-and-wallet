package internal

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
