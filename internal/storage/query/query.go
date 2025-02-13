package query

import (
	_ "embed"
)

//go:embed scripts/search_user.sql
var SearchUser string

//go:embed scripts/create_new_user.sql
var CreateNewUser string

//go:embed scripts/get_user_by_id.sql
var GetUserById string

//go:embed scripts/get_user_coin_history.sql
var GetUserCoinHistory string

//go:embed scripts/get_user_inventory.sql
var GetUserInventory string
