package buy_item

import (
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type merchService interface {
	BuyItem(ctx context.Context, userID int, itemName string) error
}
