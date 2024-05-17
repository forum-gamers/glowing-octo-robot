package controllers

import (
	"context"
	"time"

	protobuf "github.com/forum-gamers/glowing-octo-robot/generated/transaction"
	h "github.com/forum-gamers/glowing-octo-robot/helpers"
	"github.com/forum-gamers/glowing-octo-robot/pkg/transaction"
	"github.com/forum-gamers/glowing-octo-robot/pkg/user"
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
	}, nil
}