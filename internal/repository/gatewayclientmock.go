package repository

import (
	"context"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
	"github.com/google/uuid"
)

type GatewayClientMock struct{}

func NewGatewayClient() *GatewayClientMock {
	return &GatewayClientMock{}
}

func (g *GatewayClientMock) CreatePayment(ctx context.Context, paymentRequest internal.PaymentRequest) (string, error) {
	return uuid.New().String(), nil
}
