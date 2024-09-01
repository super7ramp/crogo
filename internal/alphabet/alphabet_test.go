package alphabet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	assert.True(t, Contains('A'))
	assert.True(t, Contains('E'))
	assert.True(t, Contains('Z'))
	assert.False(t, Contains('@'))
	assert.False(t, Contains('&'))
	assert.False(t, Contains('À'))
}

func TestLetterAt(t *testing.T) {
	assert.Equal(t, 'A', LetterAt(0))
	assert.Equal(t, 'E', LetterAt(4))
	assert.Equal(t, 'Z', LetterAt(25))
}

func TestLetterAt_Oob(t *testing.T) {
	assert.Panics(t, func() { LetterAt(26) })
}

func TestIndexOf(t *testing.T) {
	index, exists := IndexOf('A')
	assert.Equal(t, 0, index)
	assert.True(t, exists)

	index, exists = IndexOf('E')
	assert.Equal(t, 4, index)
	assert.True(t, exists)

	index, exists = IndexOf('Z')
	assert.Equal(t, 25, index)
	assert.True(t, exists)

	_, exists = IndexOf('@')
	assert.False(t, exists)

	_, exists = IndexOf('&')
	assert.False(t, exists)

	_, exists = IndexOf('À')
	assert.False(t, exists)
}

func TestLetterCount(t *testing.T) {
	assert.Equal(t, 26, LetterCount())
}
