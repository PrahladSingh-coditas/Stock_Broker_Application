package models

import "strings"

type BFFAdgToWatchlistRequest struct {
	Action       ActionType `json:"action" example:"GET" validate:"required,enum"`
	ScripId      string     `json:"scripId" example:"RELI_12345" validate:"required,scrip_format"`
	WatchlistIds []uint64   `json:"watchlistIds" example:"1,2,3" validate:"watchlist_validation"`
}

type ActionType string

const (
	ADD ActionType = "ADD"
	DEL ActionType = "DEL"
	GET ActionType = "GET"
)

func (act ActionType) IsValid() bool {
	actString := strings.ToUpper(string(act))
	switch actString {
	case "ADD", "DEL", "GET":
		return true
	}
	return false
}

type WatchlistWithId struct {
	Watchlist_ID   int64  `json:"watchlist_id"`
	Watchlist_Name string `json:"watchlist_name"`
}

type BFFAdgToWatchlistResponse struct {
	Status          string            `json:"message"`
	Action          ActionType        `json:"action"`
	WatchlistWithId []WatchlistWithId `json:"watchlistNames"`
	Warnings        []string          `json:"warnings"`
}
