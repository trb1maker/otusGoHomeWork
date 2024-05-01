package rules

import (
	"errors"
	"strconv"
)

var ErrLen = errors.New("wrong length")

func validateLen(v string, rule string) error {
	l, err := strconv.Atoi(rule)
	if err != nil {
		return ErrInvalidRule
	}
	if len(v) != l {
		return ErrLen
	}
	return nil
}
