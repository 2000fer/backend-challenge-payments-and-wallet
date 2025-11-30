package services

import "github.com/2000fer/backend-challenge-payments-and-wallet/internal"

type Storage interface {
	GetBalance(userID uint64) (float64, error)
	GetTransactions(userID uint64) ([]internal.Transaction, error)
}

type WalletService struct {
	storage Storage
}

func NewWalletService(storage Storage) *WalletService {
	return &WalletService{
		storage: storage,
	}
}

func (s *WalletService) GetBalance(userID uint64) (float64, error) {
	return s.storage.GetBalance(userID)
}

func (s *WalletService) GetTransactions(userID uint64) ([]internal.Transaction, error) {
	return s.storage.GetTransactions(userID)
}
