package models

type ActionType string

type BFFAdgToWatchlistRequest struct {
	Action       ActionType `json:"action"`
	ScripId      string     `json:"scripId"`
	WatchlistIds []uint64   `json:"watchlistIds"`
}

type WatchlistWithId struct{
	WatchlistId uint64
	WatchlistName string
}

type BFFAdgToWatchlistResponse struct {
	Status          string     `json:"message"`
	Action          ActionType `json:"action"`
	WatchlistWithId []WatchlistWithId  `json:"watchlistNames"`
	Warnings        []string   `json:"warnings"`
}
