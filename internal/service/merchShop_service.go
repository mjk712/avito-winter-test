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
	SendCoin(ctx context.Context, fromUserId int, req dto.SendCoinRequest) error
	BuyItem(ctx context.Context, userId int, itemName string) error
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
			Received: toReceivedTransactionsDTO(transactions, user.Username),
			Sent:     toSentTransactionsDTO(transactions, user.Username),
		},
	}
	return infoResponse, nil
}

func (s *MerchShopServiceImpl) SendCoin(ctx context.Context, fromUserId int, req dto.SendCoinRequest) error {
	const op = "MerchShopService.SendCoin"

	if req.ToUser == "" || req.Amount <= 0 {
		return fmt.Errorf("%s:error validate request data", op)
	}

	toUserId, err := s.storageRepository.GetUserIdByUsername(ctx, req.ToUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if fromUserId == toUserId {
		return fmt.Errorf("%s: %w", op, fmt.Errorf("cannot send to same user: %w", op))
	}

	fromUser, err := s.storageRepository.GetUserById(ctx, fromUserId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if fromUser.Coins < req.Amount {
		return fmt.Errorf("%s: %w", op, errors.New("not enough coins"))
	}

	err = s.storageRepository.TransferCoins(ctx, fromUserId, toUserId, req.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *MerchShopServiceImpl) BuyItem(ctx context.Context, userId int, itemName string) error {
	const op = "MerchShopService.BuyItem"

	item, err := s.storageRepository.GetMerchByName(ctx, itemName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, errors.New("item not found"))
	}

	err = s.storageRepository.BuyItem(ctx, userId, item.Id, item.Price)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

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

func toReceivedTransactionsDTO(transactions []dao.TransactionHistory, username string) []dto.Transaction {
	var receivedTransactions []dto.Transaction
	for _, t := range transactions {
		if t.TransactionType == "transfer" && t.FromUser != username {
			receivedTransactions = append(receivedTransactions, dto.Transaction{
				FromUser: t.FromUser,
				ToUser:   t.ToUser,
				Amount:   t.Amount,
			})
		}
	}
	return receivedTransactions
}

func toSentTransactionsDTO(transactions []dao.TransactionHistory, username string) []dto.Transaction {
	var sentTransactions []dto.Transaction
	for _, t := range transactions {
		if t.TransactionType == "transfer" && t.FromUser == username {

			sentTransactions = append(sentTransactions, dto.Transaction{
				FromUser: t.FromUser,
				ToUser:   t.ToUser,
				Amount:   t.Amount,
			})
		}
	}

	return sentTransactions
}
