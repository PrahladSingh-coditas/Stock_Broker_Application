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

type Watchlists struct {
	Id            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId        uint64    `gorm:"column:user_id;not null" json:"user_id"`
	WatchlistName string    `gorm:"column:watchlist_name;not null" json:"watchlist_name"`
	ScripCount    uint64    `gorm:"column:scrip_count;default:0" json:"scrip_count"`
	LastUpdated   time.Time `gorm:"column:last_updated;not null" json:"last_updated"`
}

type WatchlistScrips struct {
	Id          uint64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	WatchlistId uint64      `gorm:"column:watchlist_id;not null;uniqueIndex:idx_watchlist_scrip" json:"watchlist_id"`
	ScripId     string      `gorm:"column:scrip_id;not null;uniqueIndex:idx_watchlist_scrip" json:"scrip_id"`
	Scrip       ScripMaster `gorm:"foreignKey:ScripId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Watchlist   Watchlists  `gorm:"foreignKey:WatchlistId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ScripMaster struct {
	Id        string `gorm:"column:id;primaryKey" json:"id"`
	ScripName string `gorm:"column:scrip_name;not null" json:"scrip_name"`
}

type DatabaseConfiguration struct {
	GormDB *gorm.DB
}
