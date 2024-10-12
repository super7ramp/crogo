package dictionaries

import (
	"crogo/internal/alphabet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUkacd_OnlyInAlphabet(t *testing.T) {
	for _, word := range Ukacd() {
		runes := []rune(word)
		for _, r := range runes {
			assert.Truef(t, alphabet.Contains(r), "word %s contains non-alphabetical rune %c", word, r)
		}
	}
}

func TestUkacd_NoDupe(t *testing.T) {
	words := Ukacd()
	seen := make(map[string]struct{})
	for _, word := range words {
		_, ok := seen[word]
		assert.Falsef(t, ok, "word %s is duplicated", word)
		seen[word] = struct{}{}
	}
}
