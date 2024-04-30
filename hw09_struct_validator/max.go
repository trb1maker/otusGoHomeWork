package hw09structvalidator

import (
	"errors"
	"strconv"
)

var ErrMax = errors.New("more then max")

type max struct {
	max int
}

func (m *max) Validate(v int) error {
	if v > m.max {
		return ErrMax
	}
	return nil
}

func newMaxValidator(s string) (Rule[int], error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, ErrNotValidTag
	}
	return &max{max: i}, nil
}
