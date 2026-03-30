package models

import "strings"

type BFFAdgToWatchlistRequest struct {
	Action       ActionType `json:"action" validate:"required,checkAction"`
	ScripId      string     `json:"scripId" validate:"required,scripFormat"`
	WatchlistIds []uint64   `json:"watchlistIds" validate:"watchlistRequired"`
}

type WatchlistWithId struct {
	WatchlistId   uint64 `json:"watchlistId"`
	WatchlistName string `json:"watchlistName"`
}

type BFFAdgToWatchlistResponse struct {
	Status          string            `json:"message"`
	Action          ActionType        `json:"action"`
	WatchlistWithId []WatchlistWithId `json:"watchlistNames"`
	Warnings        []string          `json:"warnings,omitempty"`
}

type ActionType string

const (
	ADD ActionType = "ADD"
	GET ActionType = "GET"
	DEL ActionType = "DEL"
)

func (act ActionType) IsValid() bool {
	action := strings.ToUpper(string(act))
	switch action {
	case "ADD", "GET", "DEL":
		return true
	}
	return false
}

type WatchlistValidation struct {
	WatchlistId uint64
	CanInsert   bool
	CanDelete   bool
	Reason      string
}

type WatchlistAndScrip struct {
	WatchlistId uint64
	ScripId     string
}
