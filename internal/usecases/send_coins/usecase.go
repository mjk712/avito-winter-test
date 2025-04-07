package send_coins

import (
	"avito-winter-test/internal/models/dto"
	"context"
	"errors"
	"fmt"
)

type Usecase struct {
	repository repository
}

func New(repository repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (uc *Usecase) SendCoins(ctx context.Context, fromUserID int, req dto.SendCoinRequest) error {
	const op = "MerchShopService.SendCoin"

	if req.ToUser == "" || req.Amount <= 0 {
		return fmt.Errorf("%s:error validate request data", op)
	}

	toUserID, err := uc.repository.GetUserIDByUsername(ctx, req.ToUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if fromUserID == toUserID {
		return fmt.Errorf("%s: %w", op, errors.New("cannot send to same user"))
	}

	fromUser, err := uc.repository.GetUserByID(ctx, fromUserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if fromUser.Coins < req.Amount {
		return fmt.Errorf("%s: %w", op, errors.New("not enough coins"))
	}

	err = uc.repository.TransferCoins(ctx, fromUserID, toUserID, req.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
