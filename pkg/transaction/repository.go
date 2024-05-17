package transaction

import (
	"context"
	"fmt"

	cons "github.com/forum-gamers/glowing-octo-robot/constants"
	"github.com/forum-gamers/glowing-octo-robot/database"
)

func NewTransactionRepo() TransactionRepo {
	return &TransactionRepoImpl{database.DB}
}

func (r *TransactionRepoImpl) Create(ctx context.Context, data *Transaction) error {
	return r.Db.QueryRowContext(
		ctx,
		fmt.Sprintf(`INSERT INTO %s 
		(userId, amount, type, currency, status, transactionDate, description, detail, discount, createdAt, updatedAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id
		`, cons.TRANSACTION),
		data.UserId, data.Amount, data.Type, data.Currency, data.Status, data.TransactionDate,
		data.Description, data.Detail, data.Discount, data.CreatedAt, data.UpdatedAt,
	).Scan(&data.Id)
}
