package dictionaries

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUkacd_NoAccent(t *testing.T) {
	for _, word := range Ukacd() {
		if strings.ContainsAny(word, "ÀÅÄÁÉÈÊËÎÏÍÔÓÖÙÛÜÑ") {
			assert.Failf(t, "word %s contains accent", word)
		}
	}
}

func TestUkacd_NoPunctuation(t *testing.T) {
	for _, word := range Ukacd() {
		if strings.ContainsAny(word, "!-'.?") {
			assert.Failf(t, "word %s contains accent", word)
		}
	}
}

func TestUkacd_NoSpace(t *testing.T) {
	for _, word := range Ukacd() {
		if strings.Contains(word, " ") {
			assert.Failf(t, "word %s contains space", word)
		}
	}
}
