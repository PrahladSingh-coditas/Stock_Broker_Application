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
	OtpSent      uint64 `gorm:"column:otpSent;default:null" json:"otpSent"`
	OtpExpiresAt uint64 `gorm:"column:otpExpiresAt;default:null" json:"otpExpiresAt"`
}

type DatabaseConfiguration struct {
	GormDB *gorm.DB
}

type Watchlist struct {
	Id            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId        uint64    `gorm:"column:user_id;not null" json:"userId"`
	WatchlistName string    `gorm:"column:watchlist_name;not null" json:"watchlistName"`
	ScripCount    uint16    `gorm:"column:scrip_count;default:0" json:"scripCount"`
	LastUpdatedAt time.Time `gorm:"column:last_updated_at;not null" json:"lastUpdatedAt"`
}

type ScripMaster struct {
	Id        string `gorm:"column:id;primaryKey" json:"id"`
	ScripName string `gorm:"column:scrip_name;not null" json:"scripName"`
}

type WatchlistScrip struct {
	Id          uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	WatchlistId uint64 `gorm:"column:watchlist_id;not null;uniqueIndex:uq_watchlist_scrip" json:"watchlistId"`
	ScripId     string `gorm:"column:scrip_id;not null;uniqueIndex:uq_watchlist_scrip" json:"scripId"`

	Watchlist Watchlist   `gorm:"foreignKey:WatchlistId;references:Id"`
	Scrip     ScripMaster `gorm:"foreignKey:ScripId;references:Id"`
}
