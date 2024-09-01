package grid

const SlotMinLength = 2

type Slot struct {
	start  int
	end    int
	offset int
	isDown bool
}

func newSlot(start, end, offset int, isDown bool) Slot {
	return Slot{
		start,
		end,
		offset,
		isDown,
	}
}

// NewAcrossSlot creates a new across Slot.
func NewAcrossSlot(startColumn, endColumn, row int) Slot {
	return newSlot(startColumn, endColumn, row, false)
}

// NewDownSlot creates a new down Slot.
func NewDownSlot(startRow, endRow, column int) Slot {
	return newSlot(startRow, endRow, column, true)
}

// Positions returns the positions of the cells of this slot.
func (s *Slot) Positions() []Pos {
	var positions []Pos
	for i := s.start; i < s.end; i++ {
		if s.isDown {
			positions = append(positions, NewPos(s.offset, i))
		} else {
			positions = append(positions, NewPos(i, s.offset))
		}
	}
	return positions
}

// Length returns the length of this slot.
func (s *Slot) Length() int {
	return s.end - s.start
}
