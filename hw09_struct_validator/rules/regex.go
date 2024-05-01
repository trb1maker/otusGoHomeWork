package rules

import (
	"errors"
	"regexp"
)

var ErrRegexp = errors.New("regexp error")

func validateRegexp(v string, rule string) error {
	r, err := regexp.Compile(rule)
	if err != nil {
		return ErrInvalidRule
	}
	if !r.MatchString(v) {
		return ErrRegexp
	}
	return nil
}
