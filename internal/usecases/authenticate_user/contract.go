package authenticate_user

import (
	"avito-winter-test/internal/models/dao"
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type repository interface {
	CheckUserAuth(ctx context.Context, username string) (dao.User, error)
	CreateNewUser(ctx context.Context, username string, password string) (dao.User, error)
}
