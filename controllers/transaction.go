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

	switch data.Status {
	case transaction.COMPLETED:
		return nil, status.Error(codes.FailedPrecondition, "transaction is already completed")
	case transaction.FAILED:
		return nil, status.Error(codes.FailedPrecondition, "transaction is failed")
	case transaction.CANCEL:
		return nil, status.Error(codes.FailedPrecondition, "transaction is already canceled")
	default:
		break
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
