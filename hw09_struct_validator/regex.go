package hw09structvalidator

import (
	"errors"
	"regexp"
)

var ErrRegexp = errors.New("wrong regexp")

type regex struct {
	regexp *regexp.Regexp
}

func (r *regex) Validate(v string) error {
	if !r.regexp.MatchString(v) {
		return ErrRegexp
	}
	return nil
}

func newRegexValidator(s string) (Rule[string], error) {
	if s == "" {
		return nil, ErrNotValidTag
	}
	r, err := regexp.Compile(s)
	if err != nil {
		return nil, ErrNotValidTag
	}
	return &regex{regexp: r}, nil
}
