package crogo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"iter"
	"reflect"
	"slices"
	"testing"
)

func TestNewCrossword(t *testing.T) {
	words := []string{"ABC", "DEF", "AA", "BB", "CC"}
	cells := [][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
	}
	_, err := NewCrossword(cells, words)
	assert.Nil(t, err)
}

func TestNewCrossword_Error(t *testing.T) {
	words := []string{"ABC", "DEF", "AA", "BB", "CC"}
	cells := [][]rune{{'_', '_', '_'}}
	_, err := NewCrossword(cells, words)
	assert.EqualError(t, err, "invalid value at row #0, column #0: _")
}

func TestSolve_Unsat(t *testing.T) {
	words := []string{"AAA", "BBB", "CDF" /* should be CDE */, "ABC", "ABD", "ABE"}
	cells := [][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	}
	crossword, _ := NewCrossword(cells, words)
	solutions := crossword.Solve()
	for solution := range solutions {
		t.Logf("Unexpected solution %v", solution)
		t.Fail()
	}
}

func TestSolve_Sat(t *testing.T) {
	words := []string{"AAA", "BBB", "CDE", "ABC", "ABD", "ABE"}
	grid := [][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	}
	crossword, _ := NewCrossword(grid, words)

	actualSolutions := crossword.Solve()

	expectedSolutions := [][][]rune{
		{
			{'B', 'B', 'B'},
			{'B', 'B', 'B'},
			{'B', 'B', 'B'},
		},
		{
			{'A', 'B', 'C'},
			{'A', 'B', 'D'},
			{'A', 'B', 'E'},
		},
		{
			{'A', 'A', 'A'},
			{'B', 'B', 'B'},
			{'C', 'D', 'E'},
		},
		{
			{'A', 'A', 'A'},
			{'A', 'A', 'A'},
			{'A', 'A', 'A'},
		},
	}
	assertSolutionsEqual(t, expectedSolutions, actualSolutions)
}

func assertSolutionsEqual(t *testing.T, expected [][][]rune, actual iter.Seq[[][]rune]) {
	for actualSolution := range actual {
		oldLen := len(expected)
		slices.DeleteFunc(expected, func(expectedSolution [][]rune) bool {
			return reflect.DeepEqual(actualSolution, expectedSolution)
		})
		require.NotEqualf(t, oldLen, len(expected), "Unexpected solution %v", actualSolution)
	}
	require.Equalf(t, [][][]rune{}, expected, "Missing solutions %v", expected)
}
