package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/storage"
)

type AuthService interface {
	Authenticate(ctx context.Context, reqData dto.AuthRequest) (string, error)
}

type AuthServiceImpl struct {
	storageRepository storage.Storage
}

func NewAuthService(storage storage.Storage) AuthService {
	return &AuthServiceImpl{storageRepository: storage}
}

func (s *AuthServiceImpl) Authenticate(ctx context.Context, reqData dto.AuthRequest) (string, error) {
	const op = "MerchShopService.Authenticate"

	// Валидация входных данных
	if reqData.Username == "" || reqData.Password == "" {
		return "", fmt.Errorf("%s: error validate request data", op)
	}
	// Проверяем существование пользователя
	var user dao.User
	user, err := s.storageRepository.CheckUserAuth(ctx, reqData.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Создаём нового пользователя если он не существует
			user, err = s.storageRepository.CreateNewUser(ctx, reqData.Username, reqData.Password)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}
		} else {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	} else {
		// Нашли пользователя - проверяем пароль
		if reqData.Password != user.Password {
			return "", errors.New("error invalid password")
		}
	}

	// После аутентификации генерируем jwt токен
	token, err := GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
