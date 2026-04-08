package constants

import (
	"errors"

	"gorm.io/gorm"
)

// error messages
const (
	ErrInvalidActiontype           = "invalid action type"
	ErrScripDoesnotExist           = "scrip does not exist"
	ErrUserNotFound                = "user does not exist"
	ErrDatabaseQuery               = "database query error"
	ErrInvalidWatchists            = "all watchlists are invalid"
	ErrAllWatchlistsFull           = "no watchlists of the user have capacity"
	ErrScripNotAddedToAnyWatchlist = "scrip not added to any watchlists"
)

// Request Validation Errors
const (
	ErrInvalidPayload               = "invalid required payload"
	ErrUnexpectedValue              = "unexpected value for the field."
	ErrFailedToAddScripToWatchlists = "failed to add scrip to watchlists"
	ErrScripNotInWatchlists         = "scrip does not exist in any watchlists of the user"
	ErrScripNotAddedToAnyWatchlists = "scrip was not added to any watchlist - maybe full or duplicate"
)

// errors
var (
	InternalServerError               = "internal server error"
	DatabaseQueryError                = errors.New(ErrDatabaseQuery)
	ScripDoesNotExistError            = errors.New(ErrScripDoesnotExist)
	ScripNotInWatchlistsError         = errors.New(ErrScripNotInWatchlists)
	UserNotFoundError                 = gorm.ErrRecordNotFound
	InvalidActionTypeError            = errors.New(ErrInvalidActiontype)
	InvalidWatchlistsError            = errors.New(ErrInvalidWatchists)
	AllWatchlistsFullError            = errors.New(ErrAllWatchlistsFull)
	ScripNotAddedToAnyWatchlistError  = errors.New(ErrScripNotAddedToAnyWatchlist)
	ScripNotAddedToAnyWatchlistsError = errors.New(ErrFailedToAddScripToWatchlists)
)

const (
	ErrRequestFailed = "%s request failed"
)
