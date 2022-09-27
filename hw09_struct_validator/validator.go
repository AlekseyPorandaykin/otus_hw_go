package hw09structvalidator

import (
	"fmt"
	"github.com/fixme_my_friend/hw09_struct_validator/rules"
	"reflect"
	"strconv"
)

var (
	RuleChecker *rules.RuleChecker
	EOF         = "\n"
)

func init() {
	RuleChecker = rules.NewRuleChecker()
}

func Validate(v interface{}) error {
	validErrs, err := validateStruct(reflect.Indirect(reflect.ValueOf(v)))
	if err != nil {
		return err
	}
	return &validErrs
}

type Validator struct {
	ruleChecker *rules.RuleChecker
}

func validateStruct(v reflect.Value) (ValidationErrors, error) {
	validErrs := ValidationErrors{}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		validErrs = append(validErrs, ValidationError{
			Field: t.Name(),
			Err:   ErrVarNotStruct,
		})
		return validErrs, nil
	}
	for i := 0; i < t.NumField(); i++ {
		valueField := v.Field(i)
		structField := t.Field(i)
		if !valueField.CanSet() {
			continue
		}
		errComplexTag, err := validateComplexTagField(valueField, structField)
		if err != nil {
			return validErrs, err
		}
		if len(errComplexTag) > 0 {
			validErrs = append(validErrs, ValidationError{
				Field:       structField.Name,
				Err:         nil,
				NestedError: errComplexTag,
			})
		}
		errSimpleTag, err := validateSimpleField(valueField, structField)
		if err != nil {
			return validErrs, err
		}
		if errSimpleTag != nil {
			validErrs = append(validErrs, *errSimpleTag)
		}
	}
	return validErrs, nil
}

func validateComplexTagField(valueField reflect.Value, structField reflect.StructField) (ValidationErrors, error) {
	validErrs := ValidationErrors{}
	switch valueField.Kind() {
	case reflect.Slice:
		switch str := valueField.Interface().(type) {
		case []string:
			for index, vls := range str {
				validateErr, err := validateElementSlice(vls, structField, index)
				if err != nil {
					return validErrs, err
				}
				if validateErr != nil {
					validErrs = append(validErrs, *validateErr)
				}
			}
		case []int32:
			for index, vls := range str {
				validateErr, err := validateElementSlice(vls, structField, index)
				if err != nil {
					return validErrs, err
				}
				if validateErr != nil {
					validErrs = append(validErrs, *validateErr)
				}
			}
		}
	case reflect.Struct:
		hasNestedTag, err := rules.HasNestedTag([]byte(structField.Tag))
		if err != nil {
			return validErrs, err
		}
		if hasNestedTag {
			validStructErrs, err := validateStruct(valueField)
			if err != nil {
				return validErrs, err
			}
			validErrs = append(validErrs, validStructErrs...)
		}
	}

	return validErrs, nil
}

func validateSimpleField(valueField reflect.Value, structField reflect.StructField) (*ValidationError, error) {
	var value interface{}
	switch valueField.Kind() {
	case reflect.Int:
		value = valueField.Int()
	case reflect.String:
		value = valueField.String()
	default:
		return nil, nil
	}
	validErr := RuleChecker.Valid([]byte(structField.Tag), value)
	if validErr != nil {
		return &ValidationError{Field: structField.Name, Err: validErr}, nil
	}
	return nil, nil
}

func validateElementSlice(val interface{}, structField reflect.StructField, index int) (*ValidationError, error) {
	validErr, err := validateSimpleField(reflect.ValueOf(val), structField)
	if validErr == nil || err != nil {
		return nil, err
	}
	validErr.Field = fmt.Sprintf("[%s]", strconv.Itoa(index))
	return validErr, nil
}
