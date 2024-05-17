package transaction

import (
	"context"
	"database/sql"
)

type TransactionRepo interface {
	Create(ctx context.Context, data *Transaction) error
}

type TransactionRepoImpl struct{ Db *sql.DB }

type TransactionStatus = string

const (
	PENDING   TransactionStatus = "pending"
	COMPLETED TransactionStatus = "completed"
	FAILED    TransactionStatus = "failed"
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
