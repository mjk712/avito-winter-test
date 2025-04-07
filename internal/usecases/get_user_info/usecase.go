package get_user_info

import (
	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/models/dto"
	"context"
	"fmt"
)

type Usecase struct {
	repository repository
}

func New(repository repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (uc *Usecase) GetUserInfo(ctx context.Context, userID int) (dto.InfoResponse, error) {
	const op = "MerchShopService.GetUserInfo"

	user, err := uc.repository.GetUserByID(ctx, userID)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	inventory, err := uc.repository.GetUserInventory(ctx, userID)
	if err != nil {
		return dto.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	transactions, err := uc.repository.GetUserCoinHistory(ctx, userID)
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
