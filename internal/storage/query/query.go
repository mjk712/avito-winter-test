package query

import (
	_ "embed"
)

//go:embed scripts/search_user.sql
var SearchUser string

//go:embed scripts/create_new_user.sql
var CreateNewUser string

//go:embed scripts/get_user_by_id.sql
var GetUserByID string

//go:embed scripts/get_user_coin_history.sql
var GetUserCoinHistory string

//go:embed scripts/get_user_inventory.sql
var GetUserInventory string

//go:embed scripts/decrease_user_coins.sql
var DecreaseUserCoins string

//go:embed scripts/increase_user_coins.sql
var IncreaseUserCoins string

//go:embed scripts/get_user_id_by_username.sql
var GetUserIDByUsername string

//go:embed scripts/record_transaction.sql
var RecordTransaction string

//go:embed scripts/get_user_coins.sql
var GetUserCoins string

//go:embed scripts/add_item_to_inventory.sql
var AddItemToInventory string

//go:embed scripts/get_merch_by_name.sql
var GetMerchByName string
