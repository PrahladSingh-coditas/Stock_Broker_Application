package models

import "strings"

type BFFAdgToWatchlistRequest struct {
	Action       ActionType `json:"action" example:"GET" validate:"required,actionChecker"`
	ScripId      string     `json:"scrip_id" example:"JIO_123" validate:"required,idChecker"`
	WatchlistIds []uint64   `json:"watchlist_id,omitempty" example:"0" validate:"required_if=Action DEL,required_if=Action ADD,excluded_if=Action GET"`
}

type WatchlistWithID struct {
	WatchlistId   uint64 `json:"watchlist_id" gorm:"column:id"`
	WatchlistName string `json:"watchlist_name" gorm:"column:watchlist_name"`
}

type BFFAdgToWatchlistResponse struct {
	Status          string            `json:"status" example:"action completed successfully!"`
	Action          ActionType        `json:"action" example:"GET"`
	WatchlistWithId []WatchlistWithID `json:"watchlist_with_id"`
	Warnings        []string          `json:"warnings" example:"[warnings!]"`
}

type ResultListsForADD struct {
	LimitExceededWatchlistIds []int64 `gorm:"column:limit_exceeded_ids"`
	AddedWatchlistIds         []int64 `gorm:"column:added_ids"`
	SkippedWatchlistIds       []int64 `gorm:"column:skipped_ids"`
	ScripCount                int64   `gorm:"column:scrip_count"`
}

type ActionType string

const (
	ADD ActionType = "ADD"
	DEL ActionType = "DEL"
	GET ActionType = "GET"
)

func (act ActionType) IsValid() bool {
	action := strings.ToUpper(string(act))
	switch action {
	case "ADD", "GET", "DEL":
		return true
	}
	return false
}
