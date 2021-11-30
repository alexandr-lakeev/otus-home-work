package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	LenValidator    = "len"
	InValidator     = "in"
	MinValidator    = "min"
	MaxValidator    = "max"
	RegexpValidator = "regexp"
	NestedValidator = "nested"
)

var (
	ErrNotStruct        = errors.New("expected a struct")
	ErrUnknownValidator = errors.New("unknown validator")
	ErrRegexpNotMatched = errors.New("not matched regexp")
)

type (
	Validators map[string]string

	ValidationError struct {
		Field string
		Err   error
	}

	ValidationErrors []ValidationError
)

func (v ValidationErrors) Error() string {
	errStr := "Validation error:\n"
	for _, err := range v {
		if err.Err == nil {
			continue
		}

		errStr += err.Field + ": " + err.Err.Error() + "\n"
	}

	return errStr
}

func Validate(v interface{}) error {
	return validateStruct("", v)
}

func validateStruct(rootFieldName string, v interface{}) error {
	var resultErrors ValidationErrors

	refVal := reflect.ValueOf(v)
	if refVal.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	fieldsCount := refVal.Type().NumField()

	for i := 0; i < fieldsCount; i++ {
		if rootFieldName != "" {
			rootFieldName += "."
		}

		if err := validateField(rootFieldName, refVal, i); err != nil {
			var fieldValidationErrors ValidationErrors

			if errors.As(err, &fieldValidationErrors) {
				resultErrors = append(resultErrors, fieldValidationErrors...)
			} else {
				return err
			}
		}
	}

	return resultErrors
}

func validateField(rootFieldName string, refVal reflect.Value, num int) error {
	t := refVal.Type()

	field := t.Field(num)
	tag := field.Tag.Get("validate")

	if tag == "" {
		return nil
	}

	kind := refVal.Field(num).Kind()
	value := refVal.Field(num)
	validators := tagToValidators(tag)
	fullFieldName := rootFieldName + field.Name

	// nolint
	switch kind {
	case reflect.Struct:
		if _, ok := validators[NestedValidator]; !ok {
			return nil
		}

		return validateStruct(fullFieldName, value.Interface())
	case reflect.String:
		return validateString(fullFieldName, value.String(), validators)
	case reflect.Int:
		return validateInt(fullFieldName, int(value.Int()), validators)
	case reflect.Int32:
		return validateInt(fullFieldName, int(value.Int()), validators)
	case reflect.Int64:
		return validateInt(fullFieldName, int(value.Int()), validators)
	case reflect.Slice:
		intSlice, ok := value.Interface().([]int)
		if ok {
			return validateIntSlice(fullFieldName, intSlice, validators)
		}

		stringSlice, ok := value.Interface().([]string)
		if ok {
			return validateStringSlice(fullFieldName, stringSlice, validators)
		}
	}

	return nil
}

func validateString(field, value string, validators Validators) error {
	var validationErrors ValidationErrors

	for name, validator := range validators {
		switch name {
		case LenValidator:
			expectedLen, err := strconv.Atoi(validator)
			if err != nil {
				return err
			}

			if len(value) != expectedLen {
				validationErrors = append(validationErrors, ValidationError{
					Field: field,
					Err:   fmt.Errorf("expected len %d, got %d", expectedLen, len(value)),
				})
			}
		case InValidator:
			found := false

			expected := strings.Split(validator, ",")
			for _, exp := range expected {
				if value == exp {
					found = true
					break
				}
			}

			if !found {
				validationErrors = append(validationErrors, ValidationError{
					Field: field,
					Err:   fmt.Errorf("expected %s, got %s", expected, value),
				})
			}
		case RegexpValidator:
			matched, err := regexp.Match(validator, []byte(value))
			if err != nil {
				return err
			}

			if !matched {
				validationErrors = append(validationErrors, ValidationError{
					Field: field,
					Err:   ErrRegexpNotMatched,
				})
			}
		default:
			return ErrUnknownValidator
		}
	}

	return validationErrors
}

func validateInt(field string, value int, validators Validators) error {
	var validationErrors ValidationErrors

	for name, validator := range validators {
		switch name {
		case MinValidator:
			min, err := strconv.Atoi(validator)
			if err != nil {
				return err
			}

			if value < min {
				validationErrors = append(validationErrors, ValidationError{
					Field: field,
					Err:   fmt.Errorf("expected greater or equal %d, got %d", min, value),
				})
			}
		case MaxValidator:
			max, err := strconv.Atoi(validator)
			if err != nil {
				return err
			}

			if value > max {
				validationErrors = append(validationErrors, ValidationError{
					Field: field,
					Err:   fmt.Errorf("expected less or equal %d, got %d", max, value),
				})
			}
		case InValidator:
			found := false

			expected := strings.Split(validator, ",")
			for _, exp := range expected {
				expInt, err := strconv.Atoi(exp)
				if err != nil {
					return err
				}

				if value == expInt {
					found = true
					break
				}
			}

			if !found {
				validationErrors = append(validationErrors, ValidationError{
					Field: field,
					Err:   fmt.Errorf("expected %s, got %d", expected, value),
				})
			}
		default:
			return ErrUnknownValidator
		}
	}

	return validationErrors
}

func validateIntSlice(field string, slice []int, validators Validators) error {
	var validationErrors ValidationErrors

	for key, value := range slice {
		if err := validateInt(fmt.Sprintf("%s[%d]", field, key), value, validators); err != nil {
			var validationError ValidationErrors

			if errors.As(err, &validationError) {
				validationErrors = append(validationErrors, validationError...)
			} else {
				return err
			}
		}
	}

	return validationErrors
}

func validateStringSlice(field string, slice []string, validators Validators) error {
	var validationErrors ValidationErrors

	for key, value := range slice {
		if err := validateString(fmt.Sprintf("%s[%d]", field, key), value, validators); err != nil {
			var validationError ValidationError

			if errors.As(err, &validationErrors) {
				validationErrors = append(validationErrors, validationError)
			} else {
				return err
			}
		}
	}

	return validationErrors
}

func tagToValidators(tag string) Validators {
	tagParts := strings.Split(tag, "|")
	mp := make(Validators, len(tagParts))

	for _, part := range tagParts {
		splited := strings.Split(part, ":")
		if len(splited) > 1 {
			mp[splited[0]] = splited[1]
		} else {
			mp[splited[0]] = ""
		}
	}

	return mp
}
