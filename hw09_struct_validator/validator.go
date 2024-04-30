package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

const tagName = "validate"

var ErrUnsupportedType = errors.New("unsupported type")

type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Err.Error()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	s := make([]string, 0, len(v))
	for _, e := range v {
		s = append(s, e.Error())
	}
	return strings.Join(s, ", ")
}

func Validate(v interface{}) error {
	vi := reflect.ValueOf(v)
	if vi.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}
	return validateStruct(vi)
}

func validateStruct(v reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)
	l := v.NumField()
	for i := 0; i < l; i++ {
		field := v.Field(i)
		name := v.Type().Field(i).Name
		tag, ok := v.Type().Field(i).Tag.Lookup(tagName)
		if !ok {
			continue
		}
		err := validateField(field, name, tag)
		if err == nil {
			continue
		}
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			validationErrors = append(validationErrors, *validationErr)
		} else {
			return err
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

//nolint:gocognit
func validateField(field reflect.Value, fieldName, fieldTag string) error {
	//nolint:exhaustive
	switch field.Kind() {
	case reflect.String:
		rule, err := NewStringRule(fieldTag)
		if err != nil {
			return err
		}
		if err := rule.Validate(field.String()); err != nil {
			return &ValidationError{fieldName, err}
		}
	case reflect.Int:
		rule, err := NewIntRule(fieldTag)
		if err != nil {
			return err
		}
		if err := rule.Validate(int(field.Int())); err != nil {
			return &ValidationError{fieldName, err}
		}
	case reflect.Slice:
		if field.Len() == 0 {
			return nil
		}
		//nolint:exhaustive
		switch field.Index(0).Kind() {
		case reflect.String:
			rule, err := NewStringRule(fieldTag)
			if err != nil {
				return err
			}
			for i := 0; i < field.Len(); i++ {
				if err := rule.Validate(field.Index(i).String()); err != nil {
					return &ValidationError{fieldName, err}
				}
			}
		case reflect.Int:
			rule, err := NewIntRule(fieldTag)
			if err != nil {
				return err
			}
			for i := 0; i < field.Len(); i++ {
				if err := rule.Validate(int(field.Index(i).Int())); err != nil {
					return &ValidationError{fieldName, err}
				}
			}
		default:
			return &ValidationError{fieldName, ErrUnsupportedType}
		}
	default:
		return &ValidationError{fieldName, ErrUnsupportedType}
	}
	return nil
}
