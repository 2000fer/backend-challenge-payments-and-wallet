package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
	"github.com/2000fer/backend-challenge-payments-and-wallet/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for dependencies
type mockPaymentStorage struct {
	mock.Mock
}

func (m *mockPaymentStorage) GetBalance(ctx context.Context, userID uint64) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockPaymentStorage) CreatePaymentRequest(ctx context.Context, paymentRequest internal.PaymentRequest) (string, error) {
	args := m.Called(ctx, paymentRequest)
	return args.String(0), args.Error(1)
}

func (m *mockPaymentStorage) UpdatePaymentRequest(ctx context.Context, paymentRequest internal.PaymentRequest, transactionID string, status string) error {
	args := m.Called(ctx, paymentRequest, transactionID, status)
	return args.Error(0)
}

type mockGatewayClient struct {
	mock.Mock
}

func (m *mockGatewayClient) CreatePayment(ctx context.Context, paymentRequest internal.PaymentRequest) (string, error) {
	args := m.Called(ctx, paymentRequest)
	return args.String(0), args.Error(1)
}

// ... (previous imports and mock implementations remain the same)

func TestPaymentService_CreatePayment(t *testing.T) {
	tests := []struct {
		name          string
		request       internal.PaymentRequest
		setupMocks    func(*mockPaymentStorage, *mockGatewayClient)
		expectedID    string
		expectedError error
	}{
		{
			name: "successful payment creation",
			request: internal.PaymentRequest{
				UserID: 1234,
				Amount: 100.50,
				Method: "card",
			},
			setupMocks: func(ps *mockPaymentStorage, gc *mockGatewayClient) {
				// Mock getting balance
				ps.On("GetBalance", mock.Anything, uint64(1234)).Return(1000.0, nil)

				// Mock creating payment request
				ps.On("CreatePaymentRequest", mock.Anything, mock.MatchedBy(func(pr internal.PaymentRequest) bool {
					return pr.UserID == 1234 && pr.Amount == 100.50
				})).Return("payment-123", nil)

				// Mock gateway call
				gc.On("CreatePayment", mock.Anything, mock.MatchedBy(func(pr internal.PaymentRequest) bool {
					return pr.UserID == 1234 && pr.Amount == 100.50
				})).Return("gateway-tx-123", nil)

				// Mock updating payment request
				ps.On("UpdatePaymentRequest",
					mock.Anything,
					mock.MatchedBy(func(pr internal.PaymentRequest) bool {
						return pr.UserID == 1234
					}),
					"payment-123", // Should be the same as CreatePaymentRequest return value
					internal.PaymentStatusSuccess,
				).Return(nil)
			},
			expectedID: "payment-123",
		},
		{
			name: "insufficient balance",
			request: internal.PaymentRequest{
				UserID: 1234,
				Amount: 1000.50,
				Method: "card",
			},
			setupMocks: func(ps *mockPaymentStorage, gc *mockGatewayClient) {
				ps.On("GetBalance", mock.Anything, uint64(1234)).Return(500.0, nil)
			},
			expectedError: services.ErrNotEnoughBalance,
		},
		{
			name: "gateway error",
			request: internal.PaymentRequest{
				UserID: 1234,
				Amount: 100.50,
				Method: "card",
			},
			setupMocks: func(ps *mockPaymentStorage, gc *mockGatewayClient) {
				ps.On("GetBalance", mock.Anything, uint64(1234)).Return(1000.0, nil)
				ps.On("CreatePaymentRequest", mock.Anything, mock.Anything).Return("payment-123", nil)
				gc.On("CreatePayment", mock.Anything, mock.Anything).Return("", errors.New("gateway error"))
				ps.On("UpdatePaymentRequest",
					mock.Anything,
					mock.Anything,
					"payment-123",
					internal.PaymentStatusFailed,
				).Return(nil)
			},
			expectedError: services.ErrPaymentGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mockPaymentStorage)
			mockGateway := new(mockGatewayClient)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockGateway)
			}

			service := services.NewPaymentService(mockStorage, mockGateway)
			paymentID, err := service.CreatePayment(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Empty(t, paymentID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, paymentID)
			}

			mockStorage.AssertExpectations(t)
			mockGateway.AssertExpectations(t)
		})
	}
}
