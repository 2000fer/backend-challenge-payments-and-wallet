package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/2000fer/backend-challenge-payments-and-wallet/internal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

// NewPostgresStorage creates a new PostgresStorage instance
func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	// Configure the connection pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return &PostgresStorage{pool: pool}, nil
}

// GetBalance retrieves the current balance for a user
func (s *PostgresStorage) GetBalance(ctx context.Context, userID uint64) (float64, error) {
	var balance float64
	err := s.pool.QueryRow(
		ctx,
		"SELECT balance FROM user_balances WHERE user_id = $1",
		userID,
	).Scan(&balance)

	if err == pgx.ErrNoRows {
		// Return 0 if user has no balance record yet
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("error getting balance: %v", err)
	}

	return balance, nil
}

// GetTransactions retrieves transaction history for a user
func (s *PostgresStorage) GetTransactions(ctx context.Context, userID uint64) ([]internal.Transaction, error) {
	rows, err := s.pool.Query(
		ctx,
		`SELECT id, user_id, amount, transaction_type, status, created_at 
		 FROM transactions 
		 WHERE user_id = $1 
		 ORDER BY created_at DESC 
		 LIMIT 100`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %v", err)
	}
	defer rows.Close()

	var transactions []internal.Transaction
	for rows.Next() {
		var t internal.Transaction
		err := rows.Scan(
			&t.ID,
			&t.Amount,
			&t.Type,
			&t.Status,
			&t.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			continue
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %v", err)
	}

	return transactions, nil
}

// CreatePaymentRequest creates a new payment request
func (s *PostgresStorage) CreatePaymentRequest(ctx context.Context, paymentRequest internal.PaymentRequest) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("error beginning transaction: %v", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Printf("Error in rollback transaction: %v", err)
		}
	}()

	// Generate a new transaction ID
	transactionID := uuid.New().String()

	// Insert the transaction
	_, err = tx.Exec(
		ctx,
		`INSERT INTO transactions 
		 (id, user_id, amount, transaction_type, status, created_at)
		 VALUES ($1, $2, $3, 'payment', 'pending', NOW())`,
		transactionID,
		paymentRequest.UserID,
		paymentRequest.Amount,
	)
	if err != nil {
		return "", fmt.Errorf("error creating transaction: %v", err)
	}

	// Update the user's balance
	_, err = tx.Exec(
		ctx,
		`INSERT INTO user_balances (user_id, balance)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id) 
		 DO UPDATE SET balance = user_balances.balance - $2`,
		paymentRequest.UserID,
		paymentRequest.Amount,
	)
	if err != nil {
		return "", fmt.Errorf("error updating balance: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("error committing transaction: %v", err)
	}

	return transactionID, nil
}

// UpdatePaymentRequest updates the status of a payment request
func (s *PostgresStorage) UpdatePaymentRequest(ctx context.Context, paymentRequest internal.PaymentRequest, transactionID string, status string) error {
	_, err := s.pool.Exec(
		ctx,
		"UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2",
		status,
		transactionID,
	)
	if err != nil {
		return fmt.Errorf("error updating payment request: %v", err)
	}

	if status == internal.PaymentStatusFailed {
		_, err = s.pool.Exec(
			ctx,
			`INSERT INTO user_balances (user_id, balance)
			 VALUES ($1, $2)
			 ON CONFLICT (user_id) 
			 DO UPDATE SET balance = user_balances.balance + $2`,
			paymentRequest.UserID,
			paymentRequest.Amount,
		)
		if err != nil {
			return fmt.Errorf("error updating balance: %v", err)
		}
	}

	return nil
}
