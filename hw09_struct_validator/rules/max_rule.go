package rules

var (
	ErrMaxNotInt     = ErrRule{"not int for max"}
	ErrMaxValueLarge = ErrValidate{"value larger max"}
)

type MaxRule struct {
}

func (rule *MaxRule) GetTagName() string {
	return "max"
}

func (rule *MaxRule) GetRule() string {
	return `(\d)+`
}

func (rule *MaxRule) Validate(ruleStr string, value interface{}) error {
	n, err := getIntFromInterface(value)
	if err != nil {
		return ErrMaxNotInt
	}
	maxN, err := getIntFromRule(ruleStr, getCompileReg(rule))
	if err != nil {
		return err
	}
	if n > maxN {
		return ErrMaxValueLarge
	}

	return nil
}
