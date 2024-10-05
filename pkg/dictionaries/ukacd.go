package dictionaries

import (
	_ "embed"
	"strings"
)

//go:embed UKACD18plus.txt
var ukacd string

// Ukacd returns the UKACD dictionary as a slice of strings.
func Ukacd() []string {
	words := strings.Split(ukacd, "\n")
	replacer := strings.NewReplacer(
		"-", "",
		"'", "",
		".", "",
		" ", "",
		"À", "A",
		"É", "E",
		"È", "E",
		"Ê", "E",
		"Ë", "E",
		"Î", "I",
		"Ï", "I",
		"Ô", "O",
		"Ù", "U",
		"Û", "U",
		"Ü", "U")
	for i := range words {
		words[i] = strings.ToUpper(words[i])
		words[i] = replacer.Replace(words[i])
	}
	return words
}
