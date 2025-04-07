package buy_item

import (
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

func (uc *Usecase) BuyItem(ctx context.Context, userID int, itemName string) error {
	const op = "MerchShopService.BuyItem"

	item, err := uc.repository.GetMerchByName(ctx, itemName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, errors.New("item not found"))
	}

	err = uc.repository.BuyItem(ctx, userID, item.ID, item.Price)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
