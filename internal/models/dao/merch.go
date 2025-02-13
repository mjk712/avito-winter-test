package dao

type Merch struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}
