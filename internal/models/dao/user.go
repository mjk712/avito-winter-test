package dao

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Coins    int    `db:"coins"`
}
