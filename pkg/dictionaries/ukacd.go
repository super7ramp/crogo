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
	replacer := strings.NewReplacer("-", "", "'", "", ".", "")
	for i := range words {
		words[i] = strings.ToUpper(words[i])
		words[i] = replacer.Replace(words[i])
	}
	return words
}
