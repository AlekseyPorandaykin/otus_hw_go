package rules

import (
	"fmt"
	"regexp"
	"strconv"
)

type Rule interface {
	GetTagName() string
	GetRule() string
	Validate(ruleStr string, value interface{}) error
}

func parseStr(compileReg *regexp.Regexp, validateStr []byte) []byte {
	result := []byte{}
	result = compileReg.Expand(result, []byte("$rule"), validateStr, compileReg.FindSubmatchIndex(validateStr))

	return result
}

func getCompileReg(rule Rule) *regexp.Regexp {
	regComp, _ := regexp.Compile(fmt.Sprintf("%s:(?m)(?P<rule>%s)", rule.GetTagName(), rule.GetRule()))

	return regComp
}

func getRuleStr(rule Rule, ruleStr string) string {
	return string(parseStr(getCompileReg(rule), []byte(ruleStr)))
}

func getIntFromInterface(value interface{}) (int, error) {
	if v, ok := value.(int); ok {
		return v, nil
	}
	if v, ok := value.(int64); ok {
		return int(v), nil
	}
	if v, ok := value.(int32); ok {
		return int(v), nil
	}
	return 0, ErrRuleValueNotInt
}

func getIntFromRule(ruleStr string, reg *regexp.Regexp) (int, error) {
	result := parseStr(reg, []byte(ruleStr))
	numRule, err := strconv.Atoi(string(result))
	if err != nil {
		return 0, ErrValueNotSupported
	}
	return numRule, nil
}

func getStringFromInterface(value interface{}) (string, error) {
	if v, ok := value.(string); ok {
		return v, nil
	}
	return "", ErrRuleValueNotString
}
