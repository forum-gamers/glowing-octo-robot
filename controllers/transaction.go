package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	cons "github.com/forum-gamers/glowing-octo-robot/constants"
	protobuf "github.com/forum-gamers/glowing-octo-robot/generated/transaction"
	h "github.com/forum-gamers/glowing-octo-robot/helpers"
	"github.com/forum-gamers/glowing-octo-robot/pkg/transaction"
	"github.com/forum-gamers/glowing-octo-robot/pkg/user"
	"github.com/forum-gamers/glowing-octo-robot/pkg/wallet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	protobuf.UnimplementedTransactionServiceServer
	GetUser         func(context.Context) user.User
	TransactionRepo transaction.TransactionRepo
}

func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	in *protobuf.CreateTransactionInput,
) (*protobuf.Transaction, error) {
	if in.Amount < 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater or equal than 0")
	}

	if in.Discount < 0 {
		return nil, status.Error(codes.InvalidArgument, "discount must be greater or equal than 0")
	}

	if in.Amount < in.Discount {
		return nil, status.Error(codes.InvalidArgument, "discount cannot greater than amount")
	}

	if !transaction.CheckCurrency(in.Currency) {
		return nil, status.Error(codes.InvalidArgument, "invalid or unsupport currency")
	}

	transactionDate := time.Now()
	if in.TransactionDate != "" {
		if parsed, err := h.ParseStrToDate(in.TransactionDate); err != nil {
			return nil, err
		} else {
			transactionDate = parsed
		}

		if h.IsBefore(transactionDate, time.Now()) {
			return nil, status.Error(codes.InvalidArgument, "transaction date cannot before today")
		}
	}

	if !transaction.CheckTransactionType(in.Type) {
		return nil, status.Error(codes.InvalidArgument, "invalid transaction type")
	}

	payload := transaction.Transaction{
		UserId:          s.GetUser(ctx).Id,
		Amount:          in.Amount,
		Type:            in.Type,
		Currency:        in.Currency,
		Status:          transaction.PENDING,
		TransactionDate: transactionDate,
		Description:     in.Description,
		Discount:        in.Discount,
		Detail:          in.Detail,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Signature:       in.Signature,
		ItemId:          in.ItemId,
	}
	if err := s.TransactionRepo.Create(ctx, &payload); err != nil {
		return nil, err
	}

	return &protobuf.Transaction{
		Id:              payload.Id,
		UserId:          payload.UserId,
		Amount:          payload.Amount,
		Type:            payload.Type,
		Currency:        payload.Currency,
		Status:          payload.Status,
		TransactionDate: payload.TransactionDate.String(),
		Description:     payload.Description,
		Discount:        payload.Discount,
		Detail:          payload.Detail,
		CreatedAt:       payload.CreatedAt.String(),
		UpdatedAt:       payload.UpdatedAt.String(),
		Signature:       payload.Signature,
		ItemId:          payload.ItemId,
	}, nil
}

func (s *TransactionService) CancelTransaction(
	ctx context.Context,
	in *protobuf.TransactionIdInput,
) (*protobuf.Transaction, error) {
	if in.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "transaction id is required")
	}

	if !h.IsValidUUID(in.Id) {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	data, err := s.TransactionRepo.FindById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if err := transaction.CheckTransactionStatus(data.Status); err != nil {
		return nil, err
	}

	if err := s.TransactionRepo.UpdateTransactionStatus(ctx, data.Id, transaction.CANCEL); err != nil {
		return nil, err
	}

	return &protobuf.Transaction{
		Id:              data.Id,
		UserId:          data.UserId,
		Amount:          data.Amount,
		Type:            data.Type,
		Currency:        data.Currency,
		Status:          transaction.CANCEL,
		TransactionDate: data.TransactionDate.String(),
		Description:     data.Description,
		Discount:        data.Discount,
		Detail:          data.Detail,
		CreatedAt:       data.CreatedAt.String(),
		UpdatedAt:       data.UpdatedAt.String(),
		Signature:       data.Signature,
		ItemId:          data.ItemId,
	}, nil
}

func (s *TransactionService) FindOneBySignature(
	ctx context.Context,
	in *protobuf.SignatureInput,
) (*protobuf.Transaction, error) {
	if in.Signature == "" {
		return nil, status.Error(codes.InvalidArgument, "signature is required")
	}

	data, err := s.TransactionRepo.FindOneBySignature(ctx, in.Signature)
	if err != nil {
		return nil, err
	}

	return &protobuf.Transaction{
		Id:              data.Id,
		UserId:          data.UserId,
		Amount:          data.Amount,
		Type:            data.Type,
		Currency:        data.Currency,
		Status:          data.Status,
		TransactionDate: data.TransactionDate.String(),
		Description:     data.Description,
		Discount:        data.Discount,
		Detail:          data.Detail,
		CreatedAt:       data.CreatedAt.String(),
		UpdatedAt:       data.UpdatedAt.String(),
		Signature:       data.Signature,
		ItemId:          data.ItemId,
	}, nil
}

func (s *TransactionService) SuccessTopup(
	ctx context.Context,
	in *protobuf.SignatureInput,
) (*protobuf.Wallet, error) {
	if in.Signature == "" {
		return nil, status.Error(codes.InvalidArgument, "signature is required")
	}

	tx, err := s.TransactionRepo.StartTransaction(ctx, nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var data transaction.Transaction
	if err := tx.QueryRowContext(
		ctx,
		fmt.Sprintf(`
		SELECT 
		id, userId, amount, type, currency, status, transactionDate, 
		description, discount, detail, signature, itemId, createdAt, updatedAt, fee
		FROM %s
		WHERE signature = $1 FOR UPDATE
		`, cons.TRANSACTION),
		in.Signature,
	).Scan(
		&data.Id, &data.UserId, &data.Amount, &data.Type, &data.Currency, &data.Status, &data.TransactionDate,
		&data.Description, &data.Discount, &data.Detail, &data.Signature, &data.ItemId, &data.CreatedAt, &data.UpdatedAt, &data.Fee,
	); err != nil {
		if err == sql.ErrNoRows {
			err = h.NewAppError(codes.InvalidArgument, "transaction not found")
		}
		tx.Rollback()
		return nil, err
	}

	if err := transaction.CheckTransactionStatus(data.Status); err != nil {
		return nil, err
	}

	var wallet wallet.Wallet
	if err := tx.QueryRowContext(
		ctx,
		fmt.Sprintf(`
		SELECT 
		id, userId, balance, coin, createdAt, updatedAt 
		FROM %s
		WHERE userId = $1 FOR UPDATE
		`, cons.WALLET,
		),
		data.UserId,
	).Scan(
		&wallet.Id, &wallet.UserId, &wallet.Balance, &wallet.Coin, &wallet.CreatedAt, &wallet.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			err = h.NewAppError(codes.InvalidArgument, "wallet not found")
		}
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.ExecContext(
		ctx,
		fmt.Sprintf(`UPDATE %s SET status = $1, updatedAt = NOW() WHERE id = $2`, cons.TRANSACTION),
		transaction.COMPLETED, data.Id,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	updateValue := wallet.Balance + data.Amount
	if _, err := tx.ExecContext(
		ctx,
		fmt.Sprintf(`UPDATE %s SET balance = $1, updatedAt = NOW() WHERE id = $2`, cons.WALLET),
		updateValue, wallet.Id,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &protobuf.Wallet{
		Id:        wallet.Id,
		UserId:    wallet.UserId,
		Balance:   updateValue,
		Coin:      wallet.Coin,
		CreatedAt: wallet.CreatedAt.String(),
		UpdatedAt: time.Now().String(),
	}, nil
}
