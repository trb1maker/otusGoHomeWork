package rules

import (
	"strconv"
)

func validateLen(v string, rule string) error {
	l, err := strconv.Atoi(rule)
	if err != nil {
		return ErrInvalidRule
	}
	if len(v) != l {
		return ErrLength
	}
	return nil
}
