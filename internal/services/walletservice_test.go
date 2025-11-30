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

type mockWalletRepo struct {
	mock.Mock
}

func (m *mockWalletRepo) GetBalance(ctx context.Context, userID uint64) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockWalletRepo) GetTransactions(ctx context.Context, userID uint64) ([]internal.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]internal.Transaction), args.Error(1)
}

func TestWalletService_GetBalance(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint64
		setupMock     func(*mockWalletRepo)
		expected      float64
		expectedError error
	}{
		{
			name:   "successful balance retrieval",
			userID: 1234,
			setupMock: func(m *mockWalletRepo) {
				m.On("GetBalance", mock.Anything, uint64(1234)).
					Return(1000.50, nil)
			},
			expected: 1000.50,
		},
		{
			name:   "error getting balance",
			userID: 1234,
			setupMock: func(m *mockWalletRepo) {
				m.On("GetBalance", mock.Anything, uint64(1234)).
					Return(0.0, errors.New("database error"))
			},
			expected:      0.0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockWalletRepo)
			tt.setupMock(mockRepo)

			service := services.NewWalletService(mockRepo)
			balance, err := service.GetBalance(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, balance)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestWalletService_GetTransactions(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint64
		setupMock     func(*mockWalletRepo)
		expected      []internal.Transaction
		expectedError error
	}{
		{
			name:   "successful transactions retrieval",
			userID: 1234,
			setupMock: func(m *mockWalletRepo) {
				transactions := []internal.Transaction{
					{
						ID:     "1",
						Amount: 100.50,
						Type:   "debit",
					},
				}
				m.On("GetTransactions", mock.Anything, uint64(1234)).
					Return(transactions, nil)
			},
			expected: []internal.Transaction{
				{
					ID:     "1",
					Amount: 100.50,
					Type:   "debit",
				},
			},
		},
		{
			name:   "error getting transactions",
			userID: 1234,
			setupMock: func(m *mockWalletRepo) {
				m.On("GetTransactions", mock.Anything, uint64(1234)).
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockWalletRepo)
			tt.setupMock(mockRepo)

			service := services.NewWalletService(mockRepo)
			transactions, err := service.GetTransactions(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, transactions)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
