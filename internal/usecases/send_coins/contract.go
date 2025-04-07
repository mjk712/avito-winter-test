package send_coins

import (
	"avito-winter-test/internal/models/dao"
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type repository interface {
	GetUserIDByUsername(ctx context.Context, username string) (int, error)
	GetUserByID(ctx context.Context, userID int) (dao.User, error)
	TransferCoins(ctx context.Context, fromUserID int, toUserID int, amount int) error
}
