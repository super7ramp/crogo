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
		"?", "",
		"!", "",
		"À", "A",
		"Å", "A",
		"Ä", "A",
		"Á", "A",
		"É", "E",
		"È", "E",
		"Ê", "E",
		"Ë", "E",
		"Î", "I",
		"Ï", "I",
		"Í", "I",
		"Ô", "O",
		"Ó", "O",
		"Ö", "O",
		"Ù", "U",
		"Û", "U",
		"Ü", "U",
		"Ñ", "N")
	for i := range words {
		words[i] = strings.ToUpper(words[i])
		words[i] = replacer.Replace(words[i])
	}
	return words
}
