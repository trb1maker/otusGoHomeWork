package rules

import (
	"errors"
	"strconv"
	"strings"
)

var ErrIn = errors.New("not in list")

// in.go:11:1: cognitive complexity 31 of func `validateIn` is high (> 30) (gocognit)
//
//nolint:gocognit
func validateIn[T Number | string](v T, rule string) error {
	list := strings.Split(rule, ",")
	if len(list) == 0 {
		return ErrInvalidRule
	}
	switch v := any(v).(type) {
	case int64:
		for _, l := range list {
			i, err := strconv.ParseInt(l, 10, 64)
			if err != nil {
				return ErrInvalidRule
			}
			if v == i {
				return nil
			}
		}
		return ErrIn
	case uint64:
		for _, l := range list {
			i, err := strconv.ParseUint(l, 10, 64)
			if err != nil {
				return ErrInvalidRule
			}
			if v == i {
				return nil
			}
		}
		return ErrIn
	case float64:
		for _, l := range list {
			i, err := strconv.ParseFloat(l, 64)
			if err != nil {
				return ErrInvalidRule
			}
			if v == i {
				return nil
			}
		}
		return ErrIn
	case string:
		for _, l := range list {
			if l == v {
				return nil
			}
		}
		return ErrIn
	}
	return ErrUnsupportedType
}
