package grid

import (
	"crogo/internal/alphabet"
	"fmt"
	"slices"
)

const (
	CellBlock = '#'
	CellEmpty = '.'
)

type Grid struct {
	cells [][]rune
}

// NewGrid attempts to create a new Grid from given cells. Function returns the grid if given input is valid, otherwise
// it returns an error containing details about the validation failure.
func NewGrid(cells [][]rune) (*Grid, error) {
	err := validate(cells)
	if err != nil {
		return &Grid{}, err
	}
	return &Grid{cells}, nil
}

// validate validates the given cells. Function returns an error iff validation fails.
func validate(cells [][]rune) error {
	if len(cells) == 0 {
		// Trivial case, empty Grid
		return nil
	}
	firstRowColumnCount := len(cells[0])
	for rowIndex, row := range cells {
		columnCount := len(row)
		if firstRowColumnCount != columnCount {
			return fmt.Errorf("inconsistent number of columns: Row #%v has %v columns but row #0 has %v", rowIndex, columnCount, firstRowColumnCount)
		}
		for columnIndex, value := range row {
			if value != CellEmpty && value != CellBlock && !alphabet.Contains(value) {
				return fmt.Errorf("invalid value at row #%v, column #%v: %v", rowIndex, columnIndex, string(value))
			}
		}
	}
	return nil
}

// LetterAt returns the letter at given position.
//
// Special character '#' is returned if the cell contains a block.
// Special character '.' is returned if the cell contains no value.
func (g *Grid) LetterAt(row, column int) rune {
	return g.cells[row][column]
}

// Slots returns the slots of this grid.
func (g *Grid) Slots() []Slot {
	return slices.Concat(g.acrossSlots(), g.downSlots())
}

// acrossSlots computes the across slots.
func (g *Grid) acrossSlots() []Slot {
	var slots []Slot
	columnCount := g.ColumnCount()
	for rowIndex, row := range g.cells {
		columnStartIndex := 0
		for columnIndex, cell := range row {
			if cell == CellBlock {
				if columnIndex-columnStartIndex >= SlotMinLength {
					slots = append(slots, NewAcrossSlot(columnStartIndex, columnIndex, rowIndex))
				}
				columnStartIndex = columnIndex + 1
			}
		}
		if columnCount-columnStartIndex >= SlotMinLength {
			slots = append(slots, NewAcrossSlot(columnStartIndex, columnCount, rowIndex))
		}
	}
	return slots
}

// downSlots computes the down slots.
func (g *Grid) downSlots() []Slot {
	var slots []Slot
	rowCount := g.RowCount()
	columnCount := g.ColumnCount()
	for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
		rowStartIndex := 0
		for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
			if g.LetterAt(rowIndex, columnIndex) == CellBlock {
				if rowIndex-rowStartIndex >= SlotMinLength {
					slots = append(slots, NewDownSlot(rowStartIndex, rowIndex, columnIndex))
				}
				rowStartIndex = rowIndex + 1
			}
		}
		if rowCount-rowStartIndex >= SlotMinLength {
			slots = append(slots, NewDownSlot(rowStartIndex, rowCount, columnIndex))
		}
	}
	return slots
}

// ColumnCount returns the number of columns of the grid.
func (g *Grid) ColumnCount() int {
	if len(g.cells) == 0 {
		return 0
	}
	return len(g.cells[0])
}

// RowCount returns the number of rows of the grid.
func (g *Grid) RowCount() int {
	return len(g.cells)
}

// SlotCount returns the number of slots.
func (g *Grid) SlotCount() int {
	return len(g.Slots())
}
