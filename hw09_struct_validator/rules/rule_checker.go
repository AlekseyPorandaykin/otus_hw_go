package rules

import (
	"fmt"
	"regexp"
	"strings"
)

type RuleChecker struct {
	rules []Rule
}

func NewRuleChecker() *RuleChecker {
	return (new(RuleChecker)).
		Add(new(InRule)).
		Add(new(LenRule)).
		Add(new(MaxRule)).
		Add(new(MinRule)).
		Add(new(RegexpRule))
}

func (c *RuleChecker) Add(rule Rule) *RuleChecker {
	c.rules = append(c.rules, rule)
	return c
}

func (c *RuleChecker) hasRule(ruleStr []byte) (bool, error) {
	hasTag, err := c.hasRuleTag(ruleStr)
	if err != nil || !hasTag {
		return hasTag, err
	}
	for _, r := range c.rules {
		hasValid, errRule := regexp.Match(fmt.Sprintf(`validate:"%s:(.+)"`, r.GetTagName()), ruleStr)
		if errRule != nil {
			return false, ErrIncorrectRule
		}
		if hasValid {
			return true, nil
		}
	}

	return false, nil
}

func (c *RuleChecker) hasRuleTag(ruleStr []byte) (bool, error) {
	isMatch, err := regexp.Match(`(.)*validate:(.)*`, ruleStr)
	if err != nil {
		return false, ErrIncorrectRule
	}
	return isMatch, nil
}

func (c *RuleChecker) Valid(ruleStr []byte, value interface{}) error {
	if hasRule, err := c.hasRule(ruleStr); !hasRule || err != nil {
		return err
	}
	rules, err := getRules(ruleStr)
	if err != nil {
		return err
	}
	for _, r := range rules {
		if errV := c.checkRule(r, value); errV != nil {
			return errV
		}
	}

	return nil
}

func (c *RuleChecker) checkRule(ruleStr string, value interface{}) error {
	for _, rule := range c.rules {
		ok, errR := regexp.Match(fmt.Sprintf("%s:%s", rule.GetTagName(), rule.GetRule()), []byte(ruleStr))
		if errR != nil {
			return ErrIncorrectRule
		}
		if ok {
			return rule.Validate(ruleStr, value)
		}
	}
	return ErrRuleNotFound
}

func HasNestedTag(ruleStr []byte) (bool, error) {
	hasTag, err := regexp.Match(`(.)*validate:"nested"(.)*`, ruleStr)
	if !hasTag {
		return false, nil
	}
	if err != nil {
		return false, ErrIncorrectRule
	}
	return true, nil
}

func getRules(validateStr []byte) ([]string, error) {
	compileReg, err := regexp.Compile(`validate:"(?m)(?P<rule>.+)"`)
	if err != nil {
		return []string{}, ErrIncorrectRule
	}
	ruleStr := parseStr(compileReg, validateStr)

	return strings.Split(string(ruleStr), "|"), nil
}
