package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/storage"
)

func TestAuthenticate_NewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	authService := NewAuthService(mockRepo)

	req := dto.AuthRequest{
		Username: "testUser",
		Password: "testPassword",
	}

	mockRepo.EXPECT().CheckUserAuth(gomock.Any(), "testUser").Return(dao.User{}, sql.ErrNoRows)

	mockRepo.EXPECT().CreateNewUser(gomock.Any(), "testUser", "testPassword").
		Return(dao.User{ID: 1, Username: "testUser", Password: "testPassword"}, nil)

	token, err := authService.Authenticate(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthenticate_ExistingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	authService := NewAuthService(mockRepo)

	req := dto.AuthRequest{
		Username: "existUser",
		Password: "password",
	}

	mockRepo.EXPECT().CheckUserAuth(gomock.Any(), "existUser").
		Return(dao.User{ID: 1, Username: "existUser", Password: "password"}, nil)

	token, err := authService.Authenticate(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	authService := NewAuthService(mockRepo)

	req := dto.AuthRequest{
		Username: "existUser",
		Password: "wrongPassword",
	}

	mockRepo.EXPECT().CheckUserAuth(gomock.Any(), "existUser").
		Return(dao.User{ID: 1, Username: "existUser", Password: "password"}, nil)

	_, err := authService.Authenticate(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, "error invalid password", err.Error())
}
