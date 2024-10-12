package dictionaries

import (
	_ "embed"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

//go:embed UKACD18plus.txt
var ukacd string

// cleaner is a string transformer that transform or remove any character that the crossword solver doesn't support.
var cleaner = transform.Chain(norm.NFD,
	runes.Remove(runes.In(unicode.Mn)),
	runes.Remove(runes.In(unicode.Punct)),
	runes.Remove(runes.In(unicode.Space)),
	runes.Map(func(r rune) rune { return unicode.ToUpper(r) }),
	norm.NFC)

// Ukacd returns the UKACD dictionary as a slice of strings.
func Ukacd() []string {
	cleanUkacd, _, _ := transform.String(cleaner, ukacd)
	cleanUkacd = strings.ReplaceAll(cleanUkacd, "Ø", "OE")
	return strings.Split(cleanUkacd, "\n")
}
