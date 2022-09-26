package rules

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrRegexpNotValidRule     = ErrValidate{"not valid regular"}
	ErrRegexpIncorrectExt     = ErrRule{"incorrect regular ext"}
	ErrRegexpNotSupportedType = ErrRule{"type value not supported for regexp"}
)

type RegexpRule struct {
}

func (rule *RegexpRule) GetTagName() string {
	return "regexp"
}

func (rule *RegexpRule) GetRule() string {
	return `(.)+`
}

func (rule *RegexpRule) getRule(ruleStr string) string {
	return strings.Replace(getRuleStr(rule, ruleStr), `\\`, `\`, -1)
}

func (rule *RegexpRule) Validate(ruleStr string, value interface{}) error {
	v, err := getValueRegexp(value)
	if err != nil {
		return err
	}
	result, err := regexp.Match(rule.getRule(ruleStr), []byte(v))
	if err != nil {
		return ErrRegexpIncorrectExt
	}
	if !result {
		return ErrRegexpNotValidRule
	}
	return nil
}

func getValueRegexp(value interface{}) (string, error) {
	switch value.(type) {
	case string:
		sValue, _ := getStringFromInterface(value)
		return sValue, nil
	case int, int32, int64:
		iValue, _ := getIntFromInterface(value)
		return strconv.Itoa(iValue), nil
	default:
		return "", ErrRegexpNotSupportedType
	}
}
