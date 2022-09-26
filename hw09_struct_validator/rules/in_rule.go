package rules

import (
	"strconv"
	"strings"
)

var (
	ErrInNotInRange = ErrValidate{"value not in range"}
)

type InRule struct {
}

func (rule *InRule) GetTagName() string {
	return "in"
}

func (rule *InRule) GetRule() string {
	return `(\w+,)+\w+`
}

func (rule *InRule) Validate(ruleStr string, value interface{}) error {
	switch value.(type) {
	case string:
		for _, r := range rule.getRange(ruleStr) {
			if value == r {
				return nil
			}
		}
		return ErrInNotInRange
	case int, int32, int64:
		v, err := getIntFromInterface(value)
		if err != nil {
			return err
		}
		for _, r := range rule.getRange(ruleStr) {
			rer, errR := strconv.Atoi(r)
			if errR != nil {
				return ErrValueNotSupported
			}
			if v == rer {
				return nil
			}
		}
		return ErrInNotInRange
	default:
		return ErrValueNotSupported
	}
}

func (rule *InRule) getRange(ruleStr string) []string {
	return strings.Split(getRuleStr(rule, ruleStr), ",")
}
