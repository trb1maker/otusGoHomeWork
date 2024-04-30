package rules

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

var ErrEnum = errors.New("not in enum")

type enum[T string | int] struct {
	enum []T
}

func (e *enum[T]) Validate(v T) error {
	if !slices.Contains(e.enum, v) {
		return ErrEnum
	}
	return nil
}

func newEnumIntValidator(s string) (Rule[int], error) {
	if s == "" {
		return nil, ErrNotValidTag
	}
	elms := strings.Split(s, ",")
	if len(elms) == 0 {
		return nil, ErrNotValidTag
	}
	v := make([]int, 0, len(elms))
	for _, elm := range elms {
		i, err := strconv.Atoi(elm)
		if err != nil {
			return nil, ErrNotValidTag
		}
		v = append(v, i)
	}
	return &enum[int]{enum: v}, nil
}

func newEnumStringValidator(s string) (Rule[string], error) {
	if s == "" {
		return nil, ErrNotValidTag
	}
	elms := strings.Split(s, ",")
	if len(elms) == 0 {
		return nil, ErrNotValidTag
	}
	return &enum[string]{enum: elms}, nil
}
