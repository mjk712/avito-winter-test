package get_user_info

import (
	"context"

	"avito-winter-test/internal/models/dto"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type getUserInfoUsecase interface {
	GetUserInfo(ctx context.Context, userID int) (dto.InfoResponse, error)
}
