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

type Watchlists struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId        uint64    `gorm:"column:userId;type:bigint;not null" json:"userId"`
	WatchlistName string    `gorm:"column:watchlistName;not null" json:"watchlistName"`
	ScripCount    uint16    `gorm:"column:scripCount;type:smallint;default:0" json:"scripCount"`
	LastUpdatedAt time.Time `gorm:"column:lastUpdatedAt;not null" json:"lastUpdatedAt"`
}

type WatchlistScrip struct {
	ID          uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	WatchlistId uint64 `gorm:"column:watchlistId;not null;uniqueIndex:uq_watchlist_scrip" json:"watchlistId"`
	ScripId     string `gorm:"column:scripId;not null;uniqueIndex:uq_watchlist_scrip" json:"scripId"`

	Watchlists  Watchlists  `gorm:"foreignKey:WatchlistId;references:ID"`
	ScripMaster ScripMaster `gorm:"foreignKey:ScripId;references:ID"`
}

type ScripMaster struct {
	ID        string `gorm:"column:id;primaryKey" json:"id"`
	ScripName string `gorm:"column:scripName;not null" json:"scripName"`
}
