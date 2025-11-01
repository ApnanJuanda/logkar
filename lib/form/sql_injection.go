package form

import "regexp"

func SQLInjector(input string) string {
	re := regexp.MustCompile(`['\"\n\r\t\;\$\^\*\\]|://`)
	input = re.ReplaceAllLiteralString(input, "")
	return input
}

func SQLInjectorNumber(input string) string {
	re := regexp.MustCompile(`[\Wa-zA-Z_]`)
	input = re.ReplaceAllLiteralString(input, "")
	return input
}
