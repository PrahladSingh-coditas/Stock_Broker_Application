package constants

const (
	RunningServerPort = "Running Server on port : %v"
)

const (
	ServiceName       = "watchlist"
	PortDefaultValude = 1000
)

// Database query related
const (
	Successful                    = "successful"
	FieldId                       = "id"
	FielduserId                   = "user_id"
	FieldScripCount               = "scrip_count"
	FieldLastUpdated              = "last_updated"
	UsersTableName                = "users"
	Fieldemail                    = "email"
	UsernameCondition             = "username = ? "
	Password                      = "password"
	Database                      = "database"
	WatchlistTableName            = "watchlists"
	WatchlistScripsTableName      = "watchlist_scrips"
	WatchlistScripMasterTableName = "scrip_masters"
	FieldwatchlistName            = "watchlist_name"
	FieldWatchlistId              = "watchlist_id"
	FieldScripId                  = "scrip_id"
	MaxScripsPerWatchlist         = 10
)

// Scrip service related
const (
	NoScripAddedToWatchlistMsg       = "no scrips were added to any watchlist"
	AlreadyExistsInSomeWatchlistMsg  = "scrip already exists in some watchlists"
	SomeWatchlistsReachedMaxLimitMsg = "some watchlists reached max limit (10)"
	NoWarningsMsg                    = "no warnings!"
	ActionCompletedMsg               = "action completed successfully"
)

// Logger and middleware related
const (
	Username = "username"
	User     = "user"
	UserId   = "userId"
	Latency  = "latency"
	Otp      = "OTP"
	Token    = "token"
	Bearer   = "Bearer "
	Subject  = "sub"
	Purpose  = "purpose"
	Server   = "server"
)

// swagger related fields
const (
	SwaggerTitle = "Stock Broker Application API"
)

// Paths
const (
	BaseConfig = "../../config"
	RootConfig = "./src/config"
)

// redis
const (
	DataFetchFromRedis = "Data fetch from the cache memory"
	UserIdScripIdKey   = "user:%s:scripId:%s"
)
