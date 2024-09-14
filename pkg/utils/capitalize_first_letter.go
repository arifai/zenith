package utils

import "strings"

// CapitalizeFirstLetter capitalizes the first letter of each word in the given text,
// where words are separated by underscores.
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
