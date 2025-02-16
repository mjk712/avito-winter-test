package dao

type TransactionHistory struct {
	ID              int    `db:"id"`
	FromUserID      *int   `db:"from_user_id"`
	ToUserID        *int   `db:"to_user_id"`
	Amount          int    `db:"amount"`
	TransactionType string `db:"transaction_type"`
	MerchID         *int   `db:"merch_id"`
	Timestamp       string `db:"timestamp"`
	FromUser        string
	ToUser          string
	MerchName       string
}
