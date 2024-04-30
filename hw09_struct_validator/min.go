package hw09structvalidator

import (
	"errors"
	"strconv"
)

var ErrMin = errors.New("less then min")

type min struct {
	min int
}

func (m *min) Validate(v int) error {
	if v < m.min {
		return ErrMin
	}
	return nil
}

func newMinValidator(s string) (Rule[int], error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, ErrNotValidTag
	}
	return &min{min: i}, nil
}
