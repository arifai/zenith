package utils

import "strings"

// CapitalizeFirstLetter is a function to capitalize first letter of a string.
// This function will split the string by underscore, and join them back with space.
func CapitalizeFirstLetter(text string) string {
	if len(text) == 0 {
		return text
	}

	words := strings.Split(text, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
