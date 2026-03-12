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
	Id            uint64    `gorm:"column:id;primaryKey" json:"watchlistId"`
	UserId        uint64    `gorm:"column:user_id" json:"userId"`
	WatchlistName string    `gorm:"column:watchlist_name" json:"watchlistName"`
	ScripCount    uint64    `gorm:"column:scrip_count;primaryKey" json:"scripCount"`
	LastUpdatedAt time.Time `gorm:"column:last_updated_at" json:"lastUpdatedAt"`
}

type WatchlistScrip struct {
	Id          uint64 `gorm:"column:id;primaryKey" json:"watchlistScripId"`
	WatchlistId uint64 `gorm:"column:watchlist_id" json:"watchlistId"`
	ScripId     string `gorm:"column:scrip_id" json:"scripId"`
}

type ScripMaster struct {
	Id        uint64 `gorm:"column:id;primaryKey" json:"scriptId"`
	ScripName string `gorm:"column:scrip_name" json:"scripName"`
}

type DatabaseConfiguration struct {
	GormDB *gorm.DB
}
