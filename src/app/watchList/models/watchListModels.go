package models

type ActionType string

const (
	ADD ActionType = "ADD"
	DEL ActionType = "DEL"
	GET ActionType = "GET"
)

type BFFAdgToWatchListRequest struct {
	Action       ActionType `json:"action" validate:"required"`
	ScripId      string     `json:"scripId" validate:"required,scripFormat"`
	WatchlistIds []uint64   `json:"watchlistIds" validate:"watchListValidation"`
}

type WatchListWithId struct {
	ID            uint64 `json:"id"`
	WatchListName string `json:"watchlistName"`
}

type BFFAdgToWatchlistResponse struct {
	Status          string            `json:"message"`
	Action          ActionType        `json:"action"`
	WatchlistWithId []WatchListWithId `json:"watchlistNames"`
	Warning         []string          `json:"warnings"`
}
