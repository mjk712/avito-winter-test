package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"avito-winter-test/internal/models/dao"
	"avito-winter-test/internal/storage/query"

	"github.com/golang-migrate/migrate/v4"
	// Подключаем драйвер PostgreSQL для работы с миграциями
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Подключаем поддержку миграций из файловой системы
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	CheckUserAuth(ctx context.Context, username string) (dao.User, error)
	CreateNewUser(ctx context.Context, username, password string) (dao.User, error)
	GetUserByID(ctx context.Context, userID int) (dao.User, error)
	GetUserInventory(ctx context.Context, userID int) ([]dao.Inventory, error)
	GetUserCoinHistory(ctx context.Context, userID int) ([]dao.TransactionHistory, error)
	GetUserIDByUsername(ctx context.Context, username string) (int, error)
	TransferCoins(ctx context.Context, fromUserID, toUserID, amount int) error
	GetMerchByName(ctx context.Context, name string) (dao.Merch, error)
	BuyItem(ctx context.Context, userID, itemID, price int) error
}

type Repository struct {
	DB *sqlx.DB
}

func New(connectionString string) (Storage, error) {
	const op = "storage.postgres.new"
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	m, err := migrate.New("file:///app/internal/storage/migrations", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Repository{DB: db}, nil
}

func (r *Repository) CheckUserAuth(ctx context.Context, username string) (dao.User, error) {
	const op = "storage.postgres.check_user_auth"
	var user dao.User
	err := r.DB.QueryRowxContext(ctx, query.SearchUser, username).StructScan(&user)
	if err != nil {
		return dao.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *Repository) CreateNewUser(ctx context.Context, username, password string) (dao.User, error) {
	const op = "storage.postgres.create_new_user"
	var user dao.User
	err := r.DB.QueryRowxContext(ctx, query.CreateNewUser, username, password).StructScan(&user)
	if err != nil {
		return dao.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID int) (dao.User, error) {
	const op = "storage.postgres.get_user_by_id"
	var user dao.User
	err := r.DB.QueryRowxContext(ctx, query.GetUserByID, userID).Scan(&user.ID, &user.Username, &user.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dao.User{}, fmt.Errorf("%s: user not found", op)
		}
		return dao.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *Repository) GetUserInventory(ctx context.Context, userID int) ([]dao.Inventory, error) {
	const op = "storage.postgres.get_user_inventory"
	rows, err := r.DB.QueryxContext(ctx, query.GetUserInventory, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var inventory []dao.Inventory
	for rows.Next() {
		var item dao.Inventory
		var merchName string
		if err := rows.Scan(&merchName, &item.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		item.MerchName = merchName
		inventory = append(inventory, item)
	}

	return inventory, nil
}

func (r *Repository) GetUserCoinHistory(ctx context.Context, userID int) ([]dao.TransactionHistory, error) {
	const op = "storage.postgres.get_user_coin_history"
	rows, err := r.DB.QueryxContext(ctx, query.GetUserCoinHistory, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var history []dao.TransactionHistory
	for rows.Next() {
		var t dao.TransactionHistory
		var fromUser, toUser, merchName *string
		if err := rows.Scan(&t.TransactionType, &t.Amount, &t.Timestamp, &fromUser, &toUser, &merchName); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if fromUser != nil {
			t.FromUser = *fromUser
		}
		if toUser != nil {
			t.ToUser = *toUser
		}
		if merchName != nil {
			t.MerchName = *merchName
		}
		history = append(history, t)
	}

	return history, nil
}

func (r *Repository) GetUserIDByUsername(ctx context.Context, username string) (int, error) {
	const op = "storage.postgres.get_user_by_name"
	var userID int
	err := r.DB.QueryRowxContext(ctx, query.GetUserIDByUsername, username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: user not found", op)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return userID, nil
}

func (r *Repository) TransferCoins(ctx context.Context, fromUserID, toUserID, amount int) error {
	const op = "storage.postgres.transfer_coins"
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(sql.ErrTxDone, err) {
			log.Printf("transaction rollback failed: %v", err)
		}
	}()

	_, err = tx.ExecContext(ctx, query.DecreaseUserCoins, amount, fromUserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query.IncreaseUserCoins, amount, toUserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query.RecordTransaction, fromUserID, toUserID, amount, "transfer")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit()
}

func (r *Repository) GetMerchByName(ctx context.Context, name string) (dao.Merch, error) {
	const op = "storage.postgres.get_merch_name"
	var merch dao.Merch
	err := r.DB.QueryRowContext(ctx, query.GetMerchByName, name).Scan(&merch.ID, &merch.Name, &merch.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dao.Merch{}, fmt.Errorf("%s: merch not found", op)
		}
		return dao.Merch{}, fmt.Errorf("%s: %w", op, err)
	}
	return merch, nil
}

func (r *Repository) BuyItem(ctx context.Context, userID, itemID, price int) error {
	const op = "storage.postgres.buy_item"
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("transaction rollback failed: %v", err)
		}
	}()

	var userCoins int
	err = tx.QueryRowContext(ctx, query.GetUserCoins, userID).Scan(&userCoins)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if userCoins < price {
		return fmt.Errorf("%s: not enough coins", op)
	}

	_, err = tx.ExecContext(ctx, query.DecreaseUserCoins, price, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query.AddItemToInventory, userID, itemID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query.RecordTransaction, sql.NullInt64{Valid: false}, userID, price, "purchase")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit()
}
