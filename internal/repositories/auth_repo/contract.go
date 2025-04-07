package auth_repo

import (
	"avito-winter-test/internal/models/dao"
	"context"
)

type storage interface {
	CheckUserAuth(ctx context.Context, username string) (dao.User, error)
	CreateNewUser(ctx context.Context, username, password string) (dao.User, error)
}
