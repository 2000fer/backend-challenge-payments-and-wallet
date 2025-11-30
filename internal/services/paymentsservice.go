package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
)

var (
	ErrNotEnoughBalance       = errors.New("not enough balance")
	ErrGettingBalance         = errors.New("error getting balance")
	ErrPaymentGateway         = errors.New("payment gateway failed")
	ErrCreatingPaymentRequest = errors.New("error creating payment request")
	ErrUpdatingPaymentRequest = errors.New("error updating payment request")
)

type PaymentStorage interface {
	GetBalance(ctx context.Context, userID uint64) (float64, error)
	CreatePaymentRequest(ctx context.Context, paymentRequest internal.PaymentRequest) (string, error)
	UpdatePaymentRequest(ctx context.Context, paymentRequest internal.PaymentRequest, transactionID string, status string) error
}

type GatewayClient interface {
	CreatePayment(ctx context.Context, paymentRequest internal.PaymentRequest) (string, error)
}

type PaymentService struct {
	storage       PaymentStorage
	gatewayClient GatewayClient
}

func NewPaymentService(storage PaymentStorage, gatewayClient GatewayClient) *PaymentService {
	return &PaymentService{
		storage:       storage,
		gatewayClient: gatewayClient,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, paymentRequest internal.PaymentRequest) (string, error) {
	balance, err := s.storage.GetBalance(ctx, paymentRequest.UserID)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrGettingBalance, err.Error())
	}

	if balance < paymentRequest.Amount {
		return "", ErrNotEnoughBalance
	}

	// Create payment request in internal storage
	transactionID, err := s.storage.CreatePaymentRequest(ctx, paymentRequest)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrCreatingPaymentRequest, err.Error())
	}

	// Send payment request to gateway
	if _, err := s.gatewayClient.CreatePayment(ctx, paymentRequest); err != nil {
		// Update transaction failed
		errUpdate := s.storage.UpdatePaymentRequest(ctx, paymentRequest, transactionID, internal.PaymentStatusFailed)
		if errUpdate != nil {
			// TODO: send to contingency plan
			return "", fmt.Errorf("%w: %s", ErrUpdatingPaymentRequest, errUpdate.Error())
		}
		return "", fmt.Errorf("%w: %s", ErrPaymentGateway, err.Error())
	}

	// Update transaction success
	err = s.storage.UpdatePaymentRequest(ctx, paymentRequest, transactionID, internal.PaymentStatusSuccess)
	if err != nil {
		// TODO: send to contingency plan
		return "", fmt.Errorf("%w: %s", ErrUpdatingPaymentRequest, err.Error())
	}

	return transactionID, nil
}
