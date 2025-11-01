package form

import (
	"regexp"
	"unicode"
)

func ValidatePassword(s string) bool {
	var hasNumber, hasChar, hasUpperCase bool
	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			hasUpperCase = true
		case unicode.IsNumber(r):
			hasNumber = true
		case unicode.IsLetter(r):
			hasChar = true
		}
	}
	whitespace := regexp.MustCompile(`\s`).MatchString(s)
	return hasNumber && hasUpperCase && hasChar && !whitespace
}
