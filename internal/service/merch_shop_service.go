package service

import (
	"context"
	"errors"
	"fmt"

	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/storage"
)

type MerchShopService interface {
	GetUserInfo(ctx context.Context, userID int) (dto.InfoResponse, error)
	SendCoin(ctx context.Context, fromUserID int, req dto.SendCoinRequest) error
	BuyItem(ctx context.Context, userID int, itemName string) error
}

type MerchShopServiceImpl struct {
	storageRepository storage.Storage
}

func NewMerchShopService(repo storage.Storage) MerchShopService {
	return &MerchShopServiceImpl{
		storageRepository: repo,
	}
}

func (s *MerchShopServiceImpl) GetUserInfo(ctx context.Context, userID int) (dto.InfoResponse, error) {
	const op = "MerchShopService.GetUserInfo"

	user, err := s.storageRepository.GetUserByID(ctx, userID)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	inventory, err := s.storageRepository.GetUserInventory(ctx, userID)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	transactions, err := s.storageRepository.GetUserCoinHistory(ctx, userID)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	// Проводим маппинг
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

func (s *MerchShopServiceImpl) SendCoin(ctx context.Context, fromUserID int, req dto.SendCoinRequest) error {
	const op = "MerchShopService.SendCoin"

	if req.ToUser == "" || req.Amount <= 0 {
		return fmt.Errorf("%s:error validate request data", op)
	}

	toUserID, err := s.storageRepository.GetUserIDByUsername(ctx, req.ToUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if fromUserID == toUserID {
		return fmt.Errorf("%s: %w", op, errors.New("cannot send to same user"))
	}

	fromUser, err := s.storageRepository.GetUserByID(ctx, fromUserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if fromUser.Coins < req.Amount {
		return fmt.Errorf("%s: %w", op, errors.New("not enough coins"))
	}

	err = s.storageRepository.TransferCoins(ctx, fromUserID, toUserID, req.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *MerchShopServiceImpl) BuyItem(ctx context.Context, userID int, itemName string) error {
	const op = "MerchShopService.BuyItem"

	item, err := s.storageRepository.GetMerchByName(ctx, itemName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, errors.New("item not found"))
	}

	err = s.storageRepository.BuyItem(ctx, userID, item.ID, item.Price)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func toInventoryDTO(inventory []dao.Inventory) []dto.InventoryItem {
	dtoInventory := make([]dto.InventoryItem, 0, 10)
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
