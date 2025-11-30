package services

import "github.com/2000fer/backend-challenge-payments-and-wallet/internal"

type WalletService struct {
}

func NewWalletService() *WalletService {
	return &WalletService{}
}

func (s *WalletService) GetBalance(userID string) (float64, error) {
	return 0, nil
}

func (s *WalletService) GetTransactions(userID string) ([]internal.Transaction, error) {
	return []internal.Transaction{}, nil
}
