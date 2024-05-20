package transaction

import (
	"context"
	"fmt"

	cons "github.com/forum-gamers/glowing-octo-robot/constants"
	"github.com/forum-gamers/glowing-octo-robot/database"
	h "github.com/forum-gamers/glowing-octo-robot/helpers"
	"google.golang.org/grpc/codes"
)

func NewTransactionRepo() TransactionRepo {
	return &TransactionRepoImpl{database.DB}
}

func (r *TransactionRepoImpl) Create(ctx context.Context, data *Transaction) error {
	return r.Db.QueryRowContext(
		ctx,
		fmt.Sprintf(`INSERT INTO %s 
		(userId, amount, type, currency, status, transactionDate, description, detail, discount, signature, itemId, createdAt, updatedAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id
		`, cons.TRANSACTION),
		data.UserId, data.Amount, data.Type, data.Currency, data.Status, data.TransactionDate,
		data.Description, data.Detail, data.Discount, data.Signature, data.ItemId, data.CreatedAt, data.UpdatedAt,
	).Scan(&data.Id)
}

func (r *TransactionRepoImpl) FindById(ctx context.Context, id string) (result Transaction, err error) {
	rows, err := r.Db.QueryContext(
		ctx,
		fmt.Sprintf(`
		SELECT 
		id, userId, amount, type, currency, status, transactionDate, 
		description, discount, detail, signature, itemId, createdAt, updatedAt 
		FROM %s 
		WHERE id = $1`, cons.TRANSACTION),
		id,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&result.Id, &result.UserId, &result.Amount, &result.Type, &result.Currency, &result.Status, &result.TransactionDate,
			&result.Description, &result.Discount, &result.Detail, &result.Signature, &result.ItemId, &result.CreatedAt, &result.UpdatedAt,
		); err != nil {
			return
		}
	}

	if result.Id == "" {
		err = h.NewAppError(codes.InvalidArgument, "data not found")
		return
	}
	return
}

func (r *TransactionRepoImpl) UpdateTransactionStatus(ctx context.Context, id string, status TransactionStatus) error {
	_, err := r.Db.ExecContext(
		ctx,
		fmt.Sprintf(`UPDATE %s SET status = $1, updatedAt = NOW() WHERE id = $2`, cons.TRANSACTION),
		status, id,
	)
	return err
}
