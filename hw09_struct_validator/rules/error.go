package rules

type ErrValidate struct {
	errStr string
}

func (err ErrValidate) Error() string {
	return err.errStr
}

type ErrRule struct {
	errStr string
}

func (err ErrRule) Error() string {
	return err.errStr
}

var (
	ErrRuleValueNotString = ErrRule{"value not string"}
	ErrRuleValueNotInt    = ErrRule{"value not int"}
	ErrRuleNotFound       = ErrRule{"not found rule"}
	ErrIncorrectRule      = ErrRule{"incorrect rule"}

	ErrValueNotSupported = ErrValidate{"not supported value"}
)
