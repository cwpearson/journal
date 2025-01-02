package tags

import (
	"strings"
	"unicode"
)

func Clean(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace all punctuation with space
	var result strings.Builder
	for _, ch := range s {
		if !unicode.IsPunct(ch) {
			result.WriteRune(ch)
		}
	}

	// Replace repeated whitespace with single dash
	words := strings.Fields(result.String())
	return strings.Join(words, "-")
}

func CleanAll(tags []string) []string {
	res := []string{}
	for _, t := range tags {
		res = append(res, Clean(t))
	}
	return res
}
