package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/storage"
)

func TestGetUserInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	service := NewMerchShopService(mockRepo)

	user := dao.User{
		ID:       1,
		Username: "testUser",
		Coins:    1000,
	}

	inventory := []dao.Inventory{
		{MerchName: "t-shirt", Quantity: 2},
	}

	transactions := []dao.TransactionHistory{
		{FromUser: "user1", ToUser: "testUser", Amount: 100, TransactionType: "transfer"},
	}

	mockRepo.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)
	mockRepo.EXPECT().GetUserInventory(gomock.Any(), 1).Return(inventory, nil)
	mockRepo.EXPECT().GetUserCoinHistory(gomock.Any(), 1).Return(transactions, nil)

	info, err := service.GetUserInfo(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1000, info.Coins)
	assert.Equal(t, 1, len(info.Inventory))
	assert.Equal(t, 1, len(info.CoinHistory.Received))
}

func TestSendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	service := NewMerchShopService(mockRepo)

	req := dto.SendCoinRequest{
		ToUser: "receiver",
		Amount: 100,
	}

	fromUser := dao.User{
		ID:    1,
		Coins: 1000,
	}

	mockRepo.EXPECT().GetUserIDByUsername(gomock.Any(), "receiver").Return(2, nil)
	mockRepo.EXPECT().GetUserByID(gomock.Any(), 1).Return(fromUser, nil)
	mockRepo.EXPECT().TransferCoins(gomock.Any(), 1, 2, 100).Return(nil)

	err := service.SendCoin(context.Background(), 1, req)
	assert.NoError(t, err)
}

func TestSendCoin_NotEnoughCoins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	service := NewMerchShopService(mockRepo)

	req := dto.SendCoinRequest{
		ToUser: "receiver",
		Amount: 100,
	}

	fromUser := dao.User{
		ID:    1,
		Coins: 50,
	}

	mockRepo.EXPECT().GetUserIDByUsername(gomock.Any(), "receiver").Return(2, nil)
	mockRepo.EXPECT().GetUserByID(gomock.Any(), 1).Return(fromUser, nil)

	err := service.SendCoin(context.Background(), 1, req)
	assert.Error(t, err)
	assert.Equal(t, "MerchShopService.SendCoin: not enough coins", err.Error())
}

func TestBuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	service := NewMerchShopService(mockRepo)

	item := dao.Merch{
		ID:    1,
		Name:  "t-shirt",
		Price: 80,
	}

	mockRepo.EXPECT().GetMerchByName(gomock.Any(), "t-shirt").Return(item, nil)
	mockRepo.EXPECT().BuyItem(gomock.Any(), 1, 1, 80).Return(nil)

	err := service.BuyItem(context.Background(), 1, "t-shirt")
	assert.NoError(t, err)
}

func TestBuyItem_ItemNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := storage.NewMockStorage(ctrl)
	service := NewMerchShopService(mockRepo)

	mockRepo.EXPECT().GetMerchByName(gomock.Any(), "t-shirt").
		Return(dao.Merch{}, errors.New("item not found"))

	err := service.BuyItem(context.Background(), 1, "t-shirt")
	assert.Error(t, err)
	assert.Equal(t, "MerchShopService.BuyItem: item not found", err.Error())
}
