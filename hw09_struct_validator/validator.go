package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/trb1maker/otus_golang_home_work/hw09_struct_validator/rules"
)

const validateTag = "validate"

var (
	valErrs      = []error{rules.ErrMax, rules.ErrMin, rules.ErrIn, rules.ErrLen, rules.ErrRegexp}
	ErrNotStruct = errors.New("not a struct")
)

type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := &strings.Builder{}
	for _, e := range v {
		b.WriteString(e.Err.Error())
	}
	return b.String()
}

func Validate(v interface{}) error {
	s := reflect.ValueOf(v)
	if s.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	if err := validateStruct(s); err != nil {
		return err
	}
	return nil
}

func validateStruct(v reflect.Value) error {
	valErrs := make(ValidationErrors, 0)
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		t := v.Type().Field(i)
		tag, ok := t.Tag.Lookup(validateTag)
		if !ok {
			continue
		}
		if err := validateField(v.Field(i), tag); err != nil {
			if isValidationError(err) {
				valErrs = append(valErrs, ValidationError{
					Field: t.Name,
					Err:   err,
				})
				continue
			}
			return err
		}
	}
	if len(valErrs) != 0 {
		return valErrs
	}
	return nil
}

func validateField(f reflect.Value, rule string) error {
	//nolint:exhaustive
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rules.ValidateNumber([]int64{f.Int()}, rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rules.ValidateNumber([]uint64{f.Uint()}, rule)
	case reflect.Float32, reflect.Float64:
		return rules.ValidateNumber([]float64{f.Float()}, rule)
	case reflect.String:
		return rules.ValidateString([]string{f.String()}, rule)
	case reflect.Slice:
		//nolint:exhaustive
		switch f.Type().Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			l := f.Len()
			v := make([]int64, 0, l)
			for i := 0; i < l; i++ {
				v = append(v, f.Index(i).Int())
			}
			return rules.ValidateNumber(v, rule)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			l := f.Len()
			v := make([]uint64, 0, l)
			for i := 0; i < l; i++ {
				v = append(v, f.Index(i).Uint())
			}
			return rules.ValidateNumber(v, rule)
		case reflect.Float32, reflect.Float64:
			l := f.Len()
			v := make([]float64, 0, l)
			for i := 0; i < l; i++ {
				v = append(v, f.Index(i).Float())
			}
			return rules.ValidateNumber(v, rule)
		case reflect.String:
			l := f.Len()
			v := make([]string, 0, l)
			for i := 0; i < l; i++ {
				v = append(v, f.Index(i).String())
			}
			return rules.ValidateString(v, rule)
		default:
			return rules.ErrUnsupportedType
		}
	default:
		return rules.ErrUnsupportedType
	}
}

func isValidationError(err error) bool {
	for i := range valErrs {
		if errors.Is(err, valErrs[i]) {
			return true
		}
	}
	return false
}
