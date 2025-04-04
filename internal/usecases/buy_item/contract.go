package buy_item

import (
	"avito-winter-test/internal/models/dao"
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type repository interface {
	GetMerchByName(ctx context.Context, name string) (dao.Merch, error)
	BuyItem(ctx context.Context, userID int, itemID int, price int) error
}
