package grid

type Pos struct {
	column int
	row    int
}

func NewPos(column, row int) Pos {
	return Pos{column, row}
}

func (p *Pos) Column() int {
	return p.column
}

func (p *Pos) Row() int {
	return p.row
}
