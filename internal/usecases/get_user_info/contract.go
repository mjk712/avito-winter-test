package get_user_info

import (
	"avito-winter-test/internal/models/dao"
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type repository interface {
	GetUserByID(ctx context.Context, userID int) (dao.User, error)
	GetUserInventory(ctx context.Context, userID int) ([]dao.Inventory, error)
	GetUserCoinHistory(ctx context.Context, userID int) ([]dao.TransactionHistory, error)
}
