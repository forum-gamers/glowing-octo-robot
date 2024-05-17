package controllers

import (
	"context"
	"errors"
	"time"

	protobuf "github.com/forum-gamers/glowing-octo-robot/generated/wallet"
	h "github.com/forum-gamers/glowing-octo-robot/helpers"
	"github.com/forum-gamers/glowing-octo-robot/pkg/user"
	"github.com/forum-gamers/glowing-octo-robot/pkg/wallet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	protobuf.UnimplementedWalletServiceServer
	GetUser    func(context.Context) user.User
	WalletRepo wallet.WalletRepo
}

func (s *WalletService) CreateWallet(ctx context.Context, in *protobuf.NoArgument) (*protobuf.Wallet, error) {
	id := s.GetUser(ctx).Id

	if exists, err := s.WalletRepo.FindByUserId(ctx, id); err != nil && !errors.Is(err, h.NewAppError(codes.NotFound, "data not found")) {
		return nil, err
	} else if exists.Id != "" {
		return nil, status.Error(codes.AlreadyExists, "wallet is already exists")
	}

	wallet := wallet.Wallet{
		UserId:    id,
		Balance:   0,
		Coin:      0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.WalletRepo.Create(ctx, &wallet); err != nil {
		return nil, err
	}
	return &protobuf.Wallet{
		Id:        wallet.Id,
		UserId:    wallet.UserId,
		Balance:   wallet.Balance,
		Coin:      wallet.Coin,
		CreatedAt: wallet.CreatedAt.String(),
		UpdatedAt: wallet.UpdatedAt.String(),
	}, nil
}
