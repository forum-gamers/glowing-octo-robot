package transaction

import (
	"context"
	"database/sql"
)

type TransactionRepo interface {
	Create(ctx context.Context, data *Transaction) error
	FindById(ctx context.Context, id string) (Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id string, status TransactionStatus) error
}

type TransactionRepoImpl struct{ Db *sql.DB }

type TransactionStatus = string

const (
	PENDING    TransactionStatus = "pending"
	COMPLETED  TransactionStatus = "completed"
	FAILED     TransactionStatus = "failed"
	CANCEL     TransactionStatus = "cancel"
	REFUND     TransactionStatus = "refund"
	SETTLEMENT TransactionStatus = "settlement"
	DENY       TransactionStatus = "deny"
	EXPIRE     TransactionStatus = "expire"
)

type TransactionCurrency = string

const (
	RUPIAH    TransactionCurrency = "IDR"
	US_DOLLAR TransactionCurrency = "USD"
)

type TransactionType = string

const (
	TOP_UP  TransactionType = "Top up"
	PAYMENT TransactionType = "Payment"
)
