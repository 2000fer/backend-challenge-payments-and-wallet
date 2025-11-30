package repository

import (
	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
	"github.com/google/uuid"
)

type GatewayClientMock struct{}

func NewGatewayClient() *GatewayClientMock {
	return &GatewayClientMock{}
}

func (g *GatewayClientMock) CreatePayment(paymentRequest internal.PaymentRequest) (string, error) {
	return uuid.New().String(), nil
}
