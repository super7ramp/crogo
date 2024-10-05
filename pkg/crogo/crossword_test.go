package crogo

import (
	"crogo/pkg/dictionaries"
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
		{'A', 'B', 'C'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	}
	crossword, _ := NewCrossword(cells, words)
	solutions := crossword.Solve()
	for solution := range solutions {
		t.Logf("Unexpected solution %c", solution)
		t.Fail()
	}
}

func TestSolve_Sat_Simple(t *testing.T) {
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

func TestSolve_Sat_Simple_Prefilled(t *testing.T) {
	words := []string{"AAA", "BBB", "CDE", "ABC", "ABD", "ABE"}
	grid := [][]rune{
		{'A', '.', '.'},
		{'B', '.', '.'},
		{'.', '.', '.'},
	}
	crossword, _ := NewCrossword(grid, words)

	actualSolutions := crossword.Solve()

	expectedSolutions := [][][]rune{
		{
			{'A', 'A', 'A'},
			{'B', 'B', 'B'},
			{'C', 'D', 'E'},
		},
	}
	assertSolutionsEqual(t, expectedSolutions, actualSolutions)
}

func TestSolve_Sat_Complex(t *testing.T) {
	words := dictionaries.Ukacd()
	grid := [][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	}
	crossword, _ := NewCrossword(grid, words)

	solutionsIter := crossword.Solve()

	expectedNextSolutions := [][][]rune{
		{
			{'F', 'L', 'U'},
			{'L', 'O', 'G'},
			{'U', 'G', 'S'},
		},
		{
			{'F', 'L', 'U'},
			{'E', 'O', 'N'},
			{'U', 'G', 'S'},
		},
		{
			{'F', 'L', 'O'},
			{'E', 'O', 'N'},
			{'U', 'G', 'S'},
		},
	}
	assertNextSolutionsEqual(t, expectedNextSolutions, solutionsIter)
}

func assertSolutionsEqual(t *testing.T, expected [][][]rune, actual iter.Seq[[][]rune]) {
	expectedRemaining := expected
	for actualSolution := range actual {
		oldLen := len(expectedRemaining)
		expectedRemaining = slices.DeleteFunc(expectedRemaining, func(expectedSolution [][]rune) bool {
			return reflect.DeepEqual(actualSolution, expectedSolution)
		})
		assert.NotEqualf(t, oldLen, len(expectedRemaining), "Unexpected solution %c", actualSolution)
	}
	require.Equalf(t, [][][]rune{}, expectedRemaining, "Missing solutions %c", expectedRemaining)
}

func assertNextSolutionsEqual(t *testing.T, someExpected [][][]rune, actualIter iter.Seq[[][]rune]) {
	getNextActual, stop := iter.Pull(actualIter)
	defer stop()
	for _, expected := range someExpected {
		actual, found := getNextActual()
		require.True(t, found, "Solution not found: %c", expected)
		require.True(t, reflect.DeepEqual(expected, actual), "Expected %c, got %c", expected, actual)
	}
}
