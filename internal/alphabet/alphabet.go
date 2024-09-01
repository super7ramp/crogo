package alphabet

import "slices"

// letters is the hardcoded Latin Script
var letters = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
	'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

// LetterAt returns the letter index in the alphabet for the given letter.
func LetterAt(index int) rune {
	return letters[index]
}

// IndexOf returns the index in the alphabet for the given letter, if it exists. The right boolean indicates whether
// it exists.
func IndexOf(letter rune) (int, bool) {
	return slices.BinarySearch(letters, letter)
}

// Contains returns `true` iff the given letter is part of the alphabet.
func Contains(letter rune) bool {
	return slices.Contains(letters, letter)
}

// LetterCount returns the size of the alphabet.
func LetterCount() int {
	return len(letters)
}
