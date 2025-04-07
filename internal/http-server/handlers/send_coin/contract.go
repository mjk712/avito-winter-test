package send_coin

import (
	"avito-winter-test/internal/models/dto"
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type sendCoinUsecase interface {
	SendCoins(ctx context.Context, fromUserID int, req dto.SendCoinRequest) error
}
