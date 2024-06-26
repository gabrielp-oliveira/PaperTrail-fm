package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func FormatRepositoryName(input string) string {
	// Replace spaces with underscores
	replaced := strings.ReplaceAll(input, " ", "_")

	// Remove accents
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, replaced)

	// Replace specific characters with hyphen
	result = strings.Map(func(r rune) rune {
		switch {
		case r == '!':
			return '-'
		// Preserve @ symbol
		case r == '@':
			return r
		// Add more specific cases as needed
		default:
			return r
		}
	}, result)

	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
