package hw09structvalidator

import (
	"errors"
	"strconv"
)

var ErrLen = errors.New("wrong length")

type length struct {
	len int
}

func (l *length) Validate(v string) error {
	if len(v) != l.len {
		return ErrLen
	}
	return nil
}

func newLengthValidator(s string) (Rule[string], error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, ErrNotValidTag
	}
	return &length{len: i}, nil
}
