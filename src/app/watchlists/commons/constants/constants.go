package constants

//Authentications API URL Keys
const (
	ServiceName       = "watchlists"
	PortDefaultValude = 1000
)

//Swagger Titile
const SwaggerTitle = "Stock Broker Application API"

//Success Messages
const (
	GetUserIdSuccessMsg       = "User Exists"
	WatchlistSuccessMsg       = "Used ActionType Successfully"
	ScripIdExistsSuccessMsg   = "ScrpId Exists"
	GetWatchlistIdsSuccessMsg = " WatchlistIds Fetched Successfully"
)

//Database Field Names
const (
	//Users
	UsersTableName = "users"
	FieldUsername  = "username"
	FieldUserId    = "id"

	//wtchlists
	WatclistsTableName = "watchlists"
	FieldWatchlistId   = "id"
	FieldWatchlistName = "watchlist_name"
	FieldWatchUserId   = "user_id"
	FieldScripCount    = "scrip_count"

	//scrip_master
	ScripMasterTableName = "scrip_masters"
	FieldScripId         = "id"
	FieldScripName       = "scrip_name"

	//watchlist_scrip
	WatchScripTableName = "watchlist_scrips"
	FieldWatchScripId   = "id"
	FieldWatchId        = "watchlist_id"
	FieldWScripId       = "scrip_id"
)

const (
	Username = "username = ?"
)
