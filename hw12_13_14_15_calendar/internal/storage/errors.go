package storage

import "errors"

var (
	ErrNotFound = errors.New("identifier not found")
	ErrNoData   = errors.New("no data")
)
