package rules

import (
	"errors"
	"strings"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrInvalidRule     = errors.New("invalid rule")
)

type Number interface {
	int64 | uint64 | float64
}

type rule struct {
	name  string
	value string
}

func newRule(s string) (rule, error) {
	elms := strings.SplitN(s, ":", 2)
	if len(elms) != 2 {
		return rule{}, ErrInvalidRule
	}
	return rule{name: elms[0], value: elms[1]}, nil
}

func parseRules(s string) ([]rule, error) {
	elms := strings.Split(s, "|")
	if len(elms) == 0 {
		return nil, ErrInvalidRule
	}
	rules := make([]rule, 0, len(elms))
	for _, elm := range elms {
		r, err := newRule(elm)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func validateString(v string, rr []rule) error {
	var err error
	for _, rule := range rr {
		switch rule.name {
		case "len":
			err = validateLen(v, rule.value)
		case "in":
			err = validateIn(v, rule.value)
		case "regexp":
			err = validateRegexp(v, rule.value)
		default:
			err = ErrInvalidRule
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateString(vv []string, rules string) error {
	rr, err := parseRules(rules)
	if err != nil {
		return err
	}
	for _, v := range vv {
		if err := validateString(v, rr); err != nil {
			return err
		}
	}
	return nil
}

func validateNumber[T Number](v T, rr []rule) error {
	var err error
	for _, rule := range rr {
		switch rule.name {
		case "min":
			err = validateMin(v, rule.value)
		case "max":
			err = validateMax(v, rule.value)
		case "in":
			err = validateIn(v, rule.value)
		default:
			err = ErrInvalidRule
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateNumber[T Number](vv []T, rules string) error {
	rr, err := parseRules(rules)
	if err != nil {
		return err
	}
	for _, v := range vv {
		if err := validateNumber(v, rr); err != nil {
			return err
		}
	}
	return nil
}
