package grid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPositions_Across(t *testing.T) {
	slot := NewAcrossSlot(1, 4, 1)
	positions := slot.Positions()
	expectedPositions := []Pos{NewPos(1, 1), NewPos(2, 1), NewPos(3, 1)}
	assert.Equal(t, expectedPositions, positions)
}

func TestPositions_Down(t *testing.T) {
	slot := NewDownSlot(1, 4, 1)
	positions := slot.Positions()
	expectedPositions := []Pos{NewPos(1, 1), NewPos(1, 2), NewPos(1, 3)}
	assert.Equal(t, expectedPositions, positions)
}
