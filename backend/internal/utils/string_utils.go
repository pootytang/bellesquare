package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

func NormalizeString(s string) string {
	// 1. Unicode NFC Normalization
	normalized := norm.NFC.String(s)
	// 2. Trim whitespace
	normalized = strings.TrimSpace(normalized)
	// 3. Lowercase
	return cases.Title(language.English).String(strings.ToLower(normalized))
}
