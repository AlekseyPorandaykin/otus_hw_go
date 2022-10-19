package rules

import (
	"strconv"
)

var (
	ErrLenStringValue  = ErrValidate{"error len string"}
	ErrLenNotFoundRule = ErrRule{"not found len"}
)

type LenRule struct {
}

func (rule *LenRule) GetTagName() string {
	return "len"
}

func (rule *LenRule) GetRule() string {
	return `(\d)+`
}

func (rule *LenRule) Validate(ruleStr string, value interface{}) error {
	v, err := getStringFromInterface(value)
	if err != nil {
		return err
	}
	lenStr, err := rule.getLen(ruleStr)

	if err != nil {
		return err
	}
	if len([]rune(v)) != lenStr {
		return ErrLenStringValue
	}
	return nil
}

func (rule *LenRule) getLen(ruleStr string) (int, error) {
	if result := getRuleStr(rule, ruleStr); result != "" {
		iLen, err := strconv.Atoi(result)
		if err != nil {
			return 0, ErrValueNotSupported
		}
		return iLen, nil
	}

	return 0, ErrLenNotFoundRule
}
