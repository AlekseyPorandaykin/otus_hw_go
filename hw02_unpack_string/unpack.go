package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	EscapeSymbol = rune(92)
	EmptySymbol  = rune(0)
)

type characterSequence struct {
	prevValue    rune
	currentValue rune
}

func Unpack(str string) (string, error) {
	var result strings.Builder
	var valueEscaped bool
	var characterSeq characterSequence

	for _, value := range str {
		characterSeq.currentValue = value
		if isValueEscaped(characterSeq) && !valueEscaped {
			valueEscaped = true
		} else {
			valueEscaped = false
		}
		// Add symbol to result string
		if isSymbol(characterSeq) || valueEscaped {
			if characterSeq.prevValue != EmptySymbol && !valueEscaped {
				result.WriteRune(characterSeq.prevValue)
			}
			characterSeq.prevValue = characterSeq.currentValue
			continue
		}
		// Repeat symbols
		if characterSeq.prevValue == EmptySymbol {
			return "", ErrInvalidString
		}
		num, errAtoi := strconv.Atoi(string(characterSeq.currentValue))
		if errAtoi != nil {
			return "", errAtoi
		}
		if num > 0 {
			result.WriteString(strings.Repeat(string(characterSeq.prevValue), num))
		}
		characterSeq.prevValue = EmptySymbol
	}
	// Add last symbol to result string
	if characterSeq.prevValue != EmptySymbol {
		result.WriteRune(characterSeq.prevValue)
	}
	return result.String(), nil
}

func isSymbol(characterSeq characterSequence) bool {
	return !unicode.IsDigit(characterSeq.currentValue)
}

func isValueEscaped(characterSeq characterSequence) bool {
	return (!isSymbol(characterSeq) || characterSeq.currentValue == EscapeSymbol) && characterSeq.prevValue == EscapeSymbol
}
