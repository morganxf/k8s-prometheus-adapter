package provider

import (
	"unicode"
)

func SplitDigitLetter(str string) (string, string) {
	var digits []rune
	var letter string
	for i, c := range str {
		if unicode.IsDigit(c) || c == '.' {
			digits = append(digits, c)
		} else {
			letter = str[i:]
			break
		}
	}
	return string(digits), letter
}