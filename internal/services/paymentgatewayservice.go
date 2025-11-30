package services

import "github.com/2000fer/backend-challenge-payments-and-wallet/internal"

type PaymentGatewayService struct {
}

func NewPaymentGatewayService() *PaymentGatewayService {
	return &PaymentGatewayService{}
}

func (s *PaymentGatewayService) CreatePayment(paymentRequest internal.PaymentRequest) (string, error) {
	return "", nil
}
