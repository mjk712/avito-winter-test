package dao

type Inventory struct {
	Id        int `db:"id"`
	UserId    int `db:"user_id"`
	MerchId   int `db:"merch_id"`
	Quantity  int `db:"quantity"`
	MerchName string
}
