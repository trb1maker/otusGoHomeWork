package hw09structvalidator

import (
	"errors"
	"strings"
)

const (
	tagMin    = "min"
	tagMax    = "max"
	tagEnum   = "in"
	tagLength = "len"
	tagRegexp = "regexp"
)

var ErrNotValidTag = errors.New("not valid tag")

type Rule[T string | int] interface {
	Validate(T) error
}

type pipe[T string | int] []Rule[T]

func (p pipe[T]) Validate(v T) error {
	for _, r := range p {
		if err := r.Validate(v); err != nil {
			return err
		}
	}
	return nil
}

//nolint:dupl
func NewIntRule(s string) (Rule[int], error) {
	elms := strings.Split(s, "|")
	if len(elms) == 0 {
		return nil, ErrNotValidTag
	}
	p := make(pipe[int], 0, len(elms))
	for _, elm := range elms {
		parts := strings.SplitN(elm, ":", 2)
		if len(parts) != 2 {
			return nil, ErrNotValidTag
		}
		switch parts[0] {
		case tagMin:
			v, err := newMinValidator(parts[1])
			if err != nil {
				return nil, err
			}
			p = append(p, v)
		case tagMax:
			v, err := newMaxValidator(parts[1])
			if err != nil {
				return nil, err
			}
			p = append(p, v)
		case tagEnum:
			v, err := newEnumIntValidator(parts[1])
			if err != nil {
				return nil, err
			}
			p = append(p, v)
		default:
			return nil, ErrNotValidTag
		}
	}
	return p, nil
}

//nolint:dupl
func NewStringRule(s string) (Rule[string], error) {
	elms := strings.Split(s, "|")
	if len(elms) == 0 {
		return nil, ErrNotValidTag
	}
	p := make(pipe[string], 0, len(elms))
	for _, elm := range elms {
		parts := strings.SplitN(elm, ":", 2)
		if len(parts) != 2 {
			return nil, ErrNotValidTag
		}
		switch parts[0] {
		case tagLength:
			v, err := newLengthValidator(parts[1])
			if err != nil {
				return nil, err
			}
			p = append(p, v)
		case tagRegexp:
			v, err := newRegexValidator(parts[1])
			if err != nil {
				return nil, err
			}
			p = append(p, v)
		case tagEnum:
			v, err := newEnumStringValidator(parts[1])
			if err != nil {
				return nil, err
			}
			p = append(p, v)
		default:
			return nil, ErrNotValidTag
		}
	}
	return p, nil
}
