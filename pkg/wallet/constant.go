package wallet

import (
	"context"
	"database/sql"
)

type WalletRepo interface {
	Create(ctx context.Context, data *Wallet) error
	FindByUserId(ctx context.Context, userId string) (Wallet, error)
}

type WalletRepoImpl struct{ Db *sql.DB }
