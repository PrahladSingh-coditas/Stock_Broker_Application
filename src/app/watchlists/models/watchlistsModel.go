package models

import "strings"

type BFFAdgToWatchlistRequest struct {
	Action       ActionType `json:"action" validate:"required"`
	ScripId      string     `json:"scripId" validate:"required"`
	WatchlistIds []uint64   `json:"watchlistIds" validate:"required"`
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
	Watchlist_ID   int64
	Watchlist_Name string
}

type BFFAdgToWatchlistResponse struct {
	Status          string            `json:"message"`
	Action          ActionType        `json:"action"`
	WatchlistWithId []WatchlistWithId `json:"watchlistNames"`
	Warnings        []string          `json:"warnings"`
}
