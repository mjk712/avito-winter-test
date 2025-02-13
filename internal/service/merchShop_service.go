package service

import (
	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type MerchShopService interface {
	Authenticate(ctx context.Context, reqData dto.AuthRequest) (string, error)
	GetUserInfo(ctx context.Context, userId int) (dto.InfoResponse, error)
}

type MerchShopServiceImpl struct {
	storageRepository storage.Storage
}

func NewMerchShopService(repo storage.Storage) MerchShopService {
	return &MerchShopServiceImpl{
		storageRepository: repo,
	}
}

func (s *MerchShopServiceImpl) Authenticate(ctx context.Context, reqData dto.AuthRequest) (string, error) {
	const op = "MerchShopService.Authenticate"

	//валидируем входные данные
	if reqData.Username == "" || reqData.Password == "" {
		return "", errors.New("error validate request data")
	}
	//проверяем существование пользователя
	var user dao.User
	user, err := s.storageRepository.CheckUserAuth(ctx, reqData.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			//создаём нового пользователя если он не существует
			user, err = s.storageRepository.CreateNewUser(ctx, reqData.Username, reqData.Password)
		} else {
			return "", fmt.Errorf("%s: %w", op, err)
		}

	} else {
		//нашли пользователя - проверяем пароль
		if reqData.Password != user.Password {
			return "", errors.New("error invalid password")
		}
	}

	//после аутентификации генерируем jwt токен
	token, err := GenerateJWT(user.Id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (s *MerchShopServiceImpl) GetUserInfo(ctx context.Context, userId int) (dto.InfoResponse, error) {
	const op = "MerchShopService.GetUserInfo"

	user, err := s.storageRepository.GetUserById(ctx, userId)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	inventory, err := s.storageRepository.GetUserInventory(ctx, userId)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	transactions, err := s.storageRepository.GetUserCoinHistory(ctx, userId)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	//проводим маппинг
	infoResponse := dto.InfoResponse{
		Coins:     user.Coins,
		Inventory: toInventoryDTO(inventory),
		CoinHistory: dto.CoinHistory{
			Received: toReceivedTransactionsDTO(transactions),
			Sent:     toSentTransactionsDTO(transactions),
		},
	}
	return infoResponse, nil
}

func toInventoryDTO(inventory []dao.Inventory) []dto.InventoryItem {
	var dtoInventory []dto.InventoryItem
	for _, item := range inventory {
		dtoInventory = append(dtoInventory, dto.InventoryItem{
			Type:     item.MerchName,
			Quantity: item.Quantity,
		})
	}
	return dtoInventory
}

func toReceivedTransactionsDTO(transactions []dao.TransactionHistory) []dto.Transaction {
	var receivedTransactions []dto.Transaction
	for _, t := range transactions {
		if t.TransactionType == "transfer" && t.FromUser == "" {
			receivedTransactions = append(receivedTransactions, dto.Transaction{
				FromUser: t.FromUser,
				ToUser:   t.ToUser,
				Amount:   t.Amount,
			})
		}
	}
	return receivedTransactions
}

func toSentTransactionsDTO(transactions []dao.TransactionHistory) []dto.Transaction {
	var sentTransactions []dto.Transaction
	for _, t := range transactions {
		if t.TransactionType == "transfer" && t.FromUser == "" {
			sentTransactions = append(sentTransactions, dto.Transaction{
				FromUser: t.FromUser,
				ToUser:   t.ToUser,
				Amount:   t.Amount,
			})
		}
	}
	return sentTransactions
}
