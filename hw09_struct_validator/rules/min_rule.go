package rules

var (
	ErrMinValueLess = ErrValidate{"value less min"}
	ErrMinNotInt    = ErrRule{"not int for min"}
)

type MinRule struct {
}

func (rule *MinRule) Validate(ruleStr string, value interface{}) error {
	n, err := getIntFromInterface(value)
	if err != nil {
		return ErrMinNotInt
	}
	minN, err := getIntFromRule(ruleStr, getCompileReg(rule))

	if err != nil {
		return err
	}
	if n < minN {
		return ErrMinValueLess
	}

	return nil
}

func (rule *MinRule) GetTagName() string {
	return "min"
}

func (rule *MinRule) GetRule() string {
	return `(\d)+`
}
