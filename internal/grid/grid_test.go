package grid

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewGrid_inconsistentLength(t *testing.T) {
	_, err := NewGrid([][]rune{
		{'A', 'B', 'C'},
		{'.', '#'},
	})
	require.EqualError(t, err, "inconsistent number of columns: Row #1 has 2 columns but row #0 has 3")
}

func TestNewGrid_invalidLetter(t *testing.T) {
	_, err := NewGrid([][]rune{
		{'A', 'B', 'C'},
		{'.', '#', '@'},
	})
	assert.EqualError(t, err, "invalid value at row #1, column #2: @")
}

func TestRowCount(t *testing.T) {
	grid, err := NewGrid([][]rune{
		{'A'},
		{'B'},
	})
	require.Nil(t, err)
	assert.Equal(t, 2, grid.RowCount())
}

func TestColumnCount(t *testing.T) {
	grid, err := NewGrid([][]rune{
		{'A'},
		{'B'},
	})
	require.Nil(t, err)
	assert.Equal(t, 1, grid.ColumnCount())
}

func TestSlots_Simple(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
		{'.', '.', '.'},
	})
	actualSlots := grid.Slots()
	expectedSlots := []Slot{
		NewAcrossSlot(0, 3, 0),
		NewAcrossSlot(0, 3, 1),
		NewAcrossSlot(0, 3, 2),
		NewDownSlot(0, 3, 0),
		NewDownSlot(0, 3, 1),
		NewDownSlot(0, 3, 2),
	}
	assert.Equal(t, expectedSlots, actualSlots)
}

func TestSlots_Asymmetrical(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '.', '.'},
		{'.', '.', '.'},
	})
	actualSlots := grid.Slots()
	expectedSlots := []Slot{
		NewAcrossSlot(0, 3, 0),
		NewAcrossSlot(0, 3, 1),
		NewDownSlot(0, 2, 0),
		NewDownSlot(0, 2, 1),
		NewDownSlot(0, 2, 2),
	}
	assert.Equal(t, expectedSlots, actualSlots)
}

func TestSlots_WithBlocks(t *testing.T) {
	grid, _ := NewGrid([][]rune{
		{'.', '#', '.'},
		{'.', '.', '.'},
		{'.', '.', '#'},
	})
	actualSlots := grid.Slots()
	expectedSlots := []Slot{
		NewAcrossSlot(0, 3, 1),
		NewAcrossSlot(0, 2, 2),
		NewDownSlot(0, 3, 0),
		NewDownSlot(1, 3, 1),
		NewDownSlot(0, 2, 2),
	}
	assert.Equal(t, expectedSlots, actualSlots)
}

func TestSlots_Empty(t *testing.T) {
	grid, _ := NewGrid(nil)
	assert.Nil(t, grid.Slots())
}
