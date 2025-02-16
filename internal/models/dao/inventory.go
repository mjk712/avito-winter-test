package dao

type Inventory struct {
	ID        int `db:"id"`
	UserID    int `db:"user_id"`
	MerchID   int `db:"merch_id"`
	Quantity  int `db:"quantity"`
	MerchName string
}
