package main

import (
	"errors"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(cfg *config) error {
	// Place your code here.
	return nil
}
