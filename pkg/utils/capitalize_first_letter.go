package utils

import "strings"

// CapitalizeFirstLetter is a function to capitalize first letter of a string
func CapitalizeFirstLetter(text string) string {
	if len(text) == 0 {
		return text
	}

	return strings.ToUpper(string(text[0])) + text[1:]
}
