package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"column:username;uniqueIndex" json:"username"`
	Password     string `gorm:"column:password" json:"password"`
	PanCard      string `gorm:"column:panCard;uniqueIndex" json:"panCard"`
	PhoneNumber  uint64 `gorm:"column:phoneNumber" json:"phoneNumber"`
	Email        string `gorm:"column:email;uniqueIndex" json:"email"`
	OtpSent      uint64 `json:"otpSent" gorm:"column:otpSent;default:null"`
	OtpExpiresAt uint64 `json:"otpExpiresAt" gorm:"column:otpExpiresAt;default:null"`
}

type Watchlists struct {
	Id            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"watchlistId"`
	UserId        uint64    `gorm:"column:user_id;not null" json:"userId"`
	WatchlistName string    `gorm:"column:watchlist_name;not null" json:"watchlistName"`
	ScripCount    uint64    `gorm:"column:scrip_count;default:0" json:"scripCount"`
	LastUpdatedAt time.Time `gorm:"column:last_updated_at;not null" json:"lastUpdatedAt"`
}

type WatchlistScrips struct {
	Id          uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"watchlistScripId"`
	WatchlistId uint64 `gorm:"column:watchlist_id;not null;uniqueIndex:uq_watchlist_scrip" json:"watchlist_Id"`
	ScripId     string `gorm:"column:scrip_id;not null;uniqueIndex:uq_watchlist_scrip" json:"scrip_Id"`

	Watchlists  Watchlists   `gorm:"foreignKey:WatchlistId;references:Id"`
	ScripMaster ScripMasters `gorm:"foreignKey:ScripId;references:Id"`
}

type ScripMasters struct {
	Id        string `gorm:"column:id;primaryKey" json:"scripId"`
	ScripName string `gorm:"column:scrip_name;not null" json:"scripName"`
}

type DatabaseConfiguration struct {
	GormDB *gorm.DB
}
