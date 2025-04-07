package authenticate_user

import (
	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/models/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Usecase struct {
	repository repository
}

func New(repository repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (uc *Usecase) AuthenticateUser(ctx context.Context, reqData dto.AuthRequest) (string, error) {
	const op = "MerchShopService.Authenticate"

	// Валидация входных данных
	if reqData.Username == "" || reqData.Password == "" {
		return "", fmt.Errorf("%s: error validate request data", op)
	}
	// Проверяем существование пользователя
	var user dao.User
	user, err := uc.repository.CheckUserAuth(ctx, reqData.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Создаём нового пользователя если он не существует
			user, err = uc.repository.CreateNewUser(ctx, reqData.Username, reqData.Password)
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

func GenerateJWT(userID int) (string, error) {
	const op = "service.GenerateJWT"
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return tokenString, nil
}
