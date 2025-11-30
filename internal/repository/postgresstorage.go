package repository

import (
	"math/rand/v2"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
	"github.com/google/uuid"
)

type PostgresStorage struct{}

func NewPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

func (s *PostgresStorage) GetBalance(userID uint64) (float64, error) {
	return rand.Float64() * 100, nil
}

func (s *PostgresStorage) GetTransactions(userID uint64) ([]internal.Transaction, error) {
	return nil, nil
}

func (s *PostgresStorage) CreatePaymentRequest(paymentRequest internal.PaymentRequest) (string, error) {
	return uuid.New().String(), nil
}

func (s *PostgresStorage) UpdatePaymentRequest(transactionID string, status string) error {
	return nil
}
