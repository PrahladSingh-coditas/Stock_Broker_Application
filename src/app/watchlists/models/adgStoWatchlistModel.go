package models

import "strings"

type BffAdgStoWatchlistRequest struct {
	Action       Actiontype `json:"action" example:"GET" validate:"required,actionChecker"` //ADD or  DEL or GET
	ScripId      string     `json:"scripId" example:"RELI" validate:"required"`
	WatchlistIds []uint64   `json:"watchlistIds,omitempty" example:"1,2,3" validate:"required_if=Action ADD required_if=Action DEL excluded_if=Action GET"`
}

type WatchlistWithId struct {
	Id   uint64
	Name string
}

type BffAdgStoWatchlistResponse struct {
	Status          string            `json:"status"`         //success or failure, if succeded - descriptive message, if failed - failed adding/ deleting/ getting scrip to wathclist...
	Action          Actiontype        `json:"action"`         //the one from the request
	WatchlistWithId []WatchlistWithId `json:"watchlistNames"` //
	Warnings        []string          `json:"warnings"`       //only while adding if a scrip already exists in a watchlist
}

type Actiontype string

const (
	AddAction Actiontype = "ADD"
	DelAction Actiontype = "DEL"
	GetAction Actiontype = "GET"
)

func (action Actiontype) IsValid() bool {
	actionString := strings.ToUpper(string(action))
	switch actionString {
	case string(AddAction), string(DelAction), string(GetAction):
		return true
	}
	return false
}
