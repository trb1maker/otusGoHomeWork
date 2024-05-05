package rules

import (
	"strconv"
)

func max[T Number](v, m T) error {
	if v > m {
		return ErrMax
	}
	return nil
}

func validateMax[T Number](v T, rule string) error {
	switch v := any(v).(type) {
	case int64:
		i, err := strconv.ParseInt(rule, 10, 64)
		if err != nil {
			return ErrInvalidRule
		}
		return max(v, i)
	case uint64:
		i, err := strconv.ParseUint(rule, 10, 64)
		if err != nil {
			return ErrInvalidRule
		}
		return max(v, i)
	case float64:
		i, err := strconv.ParseFloat(rule, 64)
		if err != nil {
			return ErrInvalidRule
		}
		return max(v, i)
	default:
		return ErrUnsupportedType
	}
}
