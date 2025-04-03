package authenticate

import (
	"avito-winter-test/internal/models/dto"
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type authService interface {
	Authenticate(ctx context.Context, reqData dto.AuthRequest) (string, error)
}
