package variables

import (
	. "crogo/internal/grid"
	. "crogo/pkg/solver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCell(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	})
	variables := NewVariables(grid, 100_000 /* does not matter here */)

	assert.Equal(t, Variable(1), variables.RepresentingCell(0, 0, 0))
	assert.Equal(t, Variable(2), variables.RepresentingCell(0, 0, 1))
	assert.Equal(t, Variable(27), variables.RepresentingCell(0, 0, 26))

	assert.Equal(t, Variable(28), variables.RepresentingCell(0, 1, 0))
	assert.Equal(t, Variable(29), variables.RepresentingCell(0, 1, 1))
	assert.Equal(t, Variable(54), variables.RepresentingCell(0, 1, 26))

	assert.Equal(t, Variable(243), variables.RepresentingCell(2, 2, 26))
}

func TestRepresentingSlot(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	})
	variables := NewVariables(grid, 100_000)

	assert.Equal(t, Variable(244), variables.RepresentingSlot(0, 0))
	assert.Equal(t, Variable(245), variables.RepresentingSlot(0, 1))
	assert.Equal(t, Variable(100_243), variables.RepresentingSlot(0, 99_999))

	assert.Equal(t, Variable(100_244), variables.RepresentingSlot(1, 0))
	assert.Equal(t, Variable(100_245), variables.RepresentingSlot(1, 1))

	assert.Equal(t, Variable(600_243), variables.RepresentingSlot(5, 99_999))
}

func TestRepresentingCellCount(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	})
	variables := NewVariables(grid, 100_000 /* does not matter here */)
	assert.Equal(t, 243, variables.RepresentingCellCount())
}

func TestRepresentingSlotCount(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	})
	variables := NewVariables(grid, 100_000)
	assert.Equal(t, 600_000, variables.RepresentingSlotCount())
}

func TestCount(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	})
	variables := NewVariables(grid, 100_000)
	assert.Equal(t, 600_243, variables.Count())
}

func TestBackToDomain(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '#', '.'},
		{'.', '.', '.'},
	})
	variables := NewVariables(grid, 1)
	var model []bool
	for cell := 0; cell < 3; cell++ {
		model = append(model, true) // state of variable 'A' for the current cell
		for variable := 1; variable < CellValueCount(); variable++ {
			model = append(model, false) // states of variable 'B' to '#' for the current cell
		}
	}
	model = append(model, false) // state of variable 'A' for the cell 4
	model = append(model, true)  // state of variable 'A' for the cell 4
	for variable := 2; variable < CellValueCount(); variable++ {
		model = append(model, false) // states of variable 'C' to '#' for the cell 4
	}
	for variable := 0; variable < CellValueCount()-1; variable++ {
		model = append(model, false) // states of variable 'A' to 'Z' for the cell 5
	}
	model = append(model, true)  // state of variable '#' for the cell 5
	model = append(model, false) // state of variable 'A' for the cell 6
	model = append(model, true)  // state of variable 'B' for the cell 6
	for variable := 2; variable < CellValueCount(); variable++ {
		model = append(model, false) // states of variable 'C' to '#' for the cell 6
	}
	for cell := 5; cell < 9; cell++ {
		model = append(model, false) // state of variable 'A' for the current cell
		model = append(model, false) // state of variable 'B' for the current cell
		model = append(model, true)  // state of variable 'C' for the current cell
		for variable := 3; variable < CellValueCount(); variable++ {
			model = append(model, false) // states of variable 'D' to '#' for the current cell
		}
	}

	solvedGrid := variables.BackToDomain(model)

	assert.Equal(t, [][]rune{
		{'A', 'A', 'A'},
		{'B', '#', 'B'},
		{'C', 'C', 'C'},
	}, solvedGrid)
}
