package wallet

import (
	"context"
	"fmt"

	cons "github.com/forum-gamers/glowing-octo-robot/constants"
	"github.com/forum-gamers/glowing-octo-robot/database"
	h "github.com/forum-gamers/glowing-octo-robot/helpers"
	"google.golang.org/grpc/codes"
)

func NewWalletRepo() WalletRepo {
	return &WalletRepoImpl{database.DB}
}

func (r *WalletRepoImpl) Create(ctx context.Context, data *Wallet) error {
	return r.Db.QueryRowContext(
		ctx,
		fmt.Sprintf(`INSERT INTO %s 
		(userId, balance, coin, createdAt, updatedAt)
		VALUES
		($1, $2, $3, $4, $5) RETURNING id
		`, cons.WALLET),
		data.UserId, data.Balance, data.Coin, data.CreatedAt, data.UpdatedAt,
	).Scan(&data.Id)
}

func (r *WalletRepoImpl) FindByUserId(ctx context.Context, userId string) (data Wallet, err error) {
	rows, err := r.Db.QueryContext(ctx, fmt.Sprintf(`SELECT * FROM %s WHERE userId = $1`, cons.WALLET), userId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&data.Id, &data.UserId, &data.Balance, &data.Coin, &data.CreatedAt, &data.UpdatedAt,
		); err != nil {
			return
		}
	}

	if data.Id == "" {
		err = h.NewAppError(codes.NotFound, "data not found")
		return
	}
	return
}
