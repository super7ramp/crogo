// Package variables is where translation of problem data from/to integer variables occurs.
//
// There are two kinds of variables:
//
//   - Variables representing cells: For each pair (cell,letter) is associated a variable. See
//     RepresentingCell for the translation.
//   - Variables representing slots: For each pair (slot,word) is associated a variable. They are placed
//     "after" the variables representing cells in the model. See RepresentingSlot for the translation.
package variables

import (
	"crogo/internal/alphabet"
	"crogo/internal/grid"
)

type Variables struct {
	grid      *grid.Grid
	wordCount int
}

// NewVariables constructs a new instance of Variables.
func NewVariables(grid *grid.Grid, wordCount int) *Variables {
	return &Variables{grid, wordCount}
}

// CellValueCount returns the number of values that a cell of a solved grid can take.
func CellValueCount() int {
	return alphabet.LetterCount() + 1
}

// BlockIndex returns the numerical representation of a block (the value of a shaded cell).
func BlockIndex() int {
	return alphabet.LetterCount()
}

// RepresentingCell returns the variable associated to the given cell and value.
func (v *Variables) RepresentingCell(row, column, value int) int {
	return row*v.grid.ColumnCount()*CellValueCount() +
		column*CellValueCount() +
		value +
		1 // variable must be strictly positive
}

// RepresentingCells returns all the variables associated to the cells of the grid.
func (v *Variables) RepresentingCells() []uint {
	variables := make([]uint, v.RepresentingCellCount())
	for rowIndex := 0; rowIndex < v.grid.RowCount(); rowIndex++ {
		for columnIndex := 0; columnIndex < v.grid.ColumnCount(); columnIndex++ {
			for value := 0; value < CellValueCount(); value++ {
				variables = append(variables, uint(v.RepresentingCell(rowIndex, columnIndex, value)))
			}
		}
	}
	return variables
}

// RepresentingSlot returns the variable associated to the given word at the given slot.
//
// RepresentingSlot variables are put after cell variables, so first slot variable corresponds to the number of cell variables
// plus 1 (because variables start at 1).
func (v *Variables) RepresentingSlot(slotIndex, wordIndex int) int {
	return v.RepresentingCellCount() + // last cell variable
		slotIndex*v.wordCount +
		wordIndex +
		1
}

// BackToDomain translates the variables states back to a crossword grid.
func (v *Variables) BackToDomain(model []bool) [][]rune {
	columnCount := v.grid.ColumnCount()
	rowCount := v.grid.RowCount()
	outputGrid := make([][]rune, rowCount)
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		outputGrid[rowIndex] = make([]rune, columnCount)
		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			for value := 0; value < CellValueCount(); value++ {
				variable := v.RepresentingCell(rowIndex, columnIndex, value) - 1
				if model[variable] {
					if value == BlockIndex() {
						outputGrid[rowIndex][columnIndex] = grid.CellBlock
					} else {
						outputGrid[rowIndex][columnIndex] = alphabet.LetterAt(value)
					}
					break
				}
			}
		}
	}
	return outputGrid
}

// RepresentingCellCount returns the number of variables representing cells.
func (v *Variables) RepresentingCellCount() int {
	return v.grid.ColumnCount() * v.grid.RowCount() * CellValueCount()
}

// RepresentingSlotCount returns the number of variables representing slots.
func (v *Variables) RepresentingSlotCount() int {
	return v.grid.SlotCount() * v.wordCount
}

// Count returns the number of variables.
func (v *Variables) Count() int {
	return v.RepresentingCellCount() + v.RepresentingSlotCount()
}
