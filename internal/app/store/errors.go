package store

import "errors"

var (
	ErrRecordNotFound = errors.New("user not found")
	ErrInvalidData    = errors.New("invalid data")
	ErrDataConflict   = errors.New("data conflict")
)
