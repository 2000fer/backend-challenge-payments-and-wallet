package services

import (
	"context"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
)

type Storage interface {
	GetBalance(ctx context.Context, userID uint64) (float64, error)
	GetTransactions(ctx context.Context, userID uint64) ([]internal.Transaction, error)
}

type WalletService struct {
	storage Storage
}

func NewWalletService(storage Storage) *WalletService {
	return &WalletService{
		storage: storage,
	}
}

func (s *WalletService) GetBalance(ctx context.Context, userID uint64) (float64, error) {
	return s.storage.GetBalance(ctx, userID)
}

func (s *WalletService) GetTransactions(ctx context.Context, userID uint64) ([]internal.Transaction, error) {
	return s.storage.GetTransactions(ctx, userID)
}
