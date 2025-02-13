package dao

type TransactionHistory struct {
	Id              int    `db:"id"`
	FromUserId      *int   `db:"from_user_id"`
	ToUserId        *int   `db:"to_user_id"`
	Amount          int    `db:"amount"`
	TransactionType string `db:"transaction_type"`
	MerchId         *int   `db:"merch_id"`
	Timestamp       string `db:"timestamp"`
	FromUser        string
	ToUser          string
	MerchName       string
}
