// Package constraints is where crossword problem constraints are built.
//
// The constraints are:
//
//   - Each cell must contain one and only one letter from the alphabet or a block.
//   - Each slot must contain one and only one word from the input word list. This is the tricky
//     part, as there must be a correspondence between cell variables and slot variables. Basically,
//     each slot variable - i.e. a representation of a (slot,word) pair - is equivalent to a
//     conjunction (= and) of cell variables - i.e. (cell,letter) pairs.
//   - Prefilled cells must be kept as is.
//
// Implementation note: Functions here add rules to the solver passed as parameter. Although having
// just a factory of constraints, to be applied separately, would be nice, it does not scale in
// terms of memory: There are too many literals and clauses. Hence, the choice to progressively add
// the clauses to the solver.
package constraints

import (
	"crogo/internal/alphabet"
	. "crogo/internal/grid"
	. "crogo/internal/variables"
	"crogo/pkg/solver"
)

type Constraints struct {
	grid      *Grid
	variables *Variables
	words     []string
}

// cellLiteralsBufferCapacity is the capacity of the buffer used to store cell literals corresponding to a word in a
// slot. Most words/slots should be smaller than this size.
const cellLiteralsBufferCapacity = 20

// NewConstraints constructs a new instance of Constraints.
func NewConstraints(grid *Grid, variables *Variables, words []string) *Constraints {
	return &Constraints{grid, variables, words}
}

// AddOneLetterOrBlockPerCellClausesTo adds the clauses ensuring that each cell must contain exactly one letter from the
// alphabet - or a block - to the given solver configurer.
func (c *Constraints) AddOneLetterOrBlockPerCellClausesTo(solverConfigurer solver.Configurer) {
	literalsBuffer := make([]solver.Literal, 0, CellValueCount())
	for row := 0; row < c.grid.RowCount(); row++ {
		for column := 0; column < c.grid.ColumnCount(); column++ {
			for letterIndex := 0; letterIndex < alphabet.LetterCount(); letterIndex++ {
				letterLiteral := solver.Literal(c.variables.RepresentingCell(row, column, letterIndex))
				literalsBuffer = append(literalsBuffer, letterLiteral)
			}
			blockLiteral := solver.Literal(c.variables.RepresentingCell(row, column, BlockIndex()))
			literalsBuffer = append(literalsBuffer, blockLiteral)
			solverConfigurer.AddExactlyOne(literalsBuffer)
			literalsBuffer = literalsBuffer[:0]
		}
	}
}

// AddOneWordPerSlotClausesTo adds the clauses ensuring that each slot must contain exactly one word from the word list to
// the given solver.
func (c *Constraints) AddOneWordPerSlotClausesTo(solverConfigurer solver.Configurer) {
	slotLiteralsBuffer := make([]solver.Literal, 0, len(c.words))
	cellLiteralsBuffer := make([]solver.Literal, 0, cellLiteralsBufferCapacity)
	for slotIndex, slot := range c.grid.Slots() {
		for wordIndex, word := range c.words {
			if len(word) == slot.Length() {
				slotLiteral := solver.Literal(c.variables.RepresentingSlot(slotIndex, wordIndex))
				slotLiteralsBuffer = append(slotLiteralsBuffer, slotLiteral)
				c.fillCellLiteralsConjunction(&cellLiteralsBuffer, slot, word)
				solverConfigurer.AddAnd(slotLiteral, cellLiteralsBuffer)
				cellLiteralsBuffer = cellLiteralsBuffer[:0]
			} // else skip this word since it obviously doesn't match the slot
		}
		solverConfigurer.AddExactlyOne(slotLiteralsBuffer)
		slotLiteralsBuffer = slotLiteralsBuffer[:0]
	}
}

// fillCellLiteralsConjunction fills the given slice with the cell literals whose conjunction (= and) is equivalent to
// the slot variable of the given slot and word.
//
// Panics if the given word contains a letter which is not in the [alphabet].
func (c *Constraints) fillCellLiteralsConjunction(cellLiterals *[]solver.Literal, slot Slot, word string) {
	slotPositions := slot.Positions()
	wordRunes := []rune(word)
	for i := range len(slotPositions) {
		letterIndex, found := alphabet.IndexOf(wordRunes[i])
		if !found {
			panic("Unsupported character " + string(wordRunes[i]))
		}
		slotPos := slotPositions[i]
		cellLiteral := solver.Literal(c.variables.RepresentingCell(slotPos.Row(), slotPos.Column(), letterIndex))
		*cellLiterals = append(*cellLiterals, cellLiteral)
	}
}

// AddInputGridConstraintsAreSatisfiedClausesTo adds the clauses ensuring that each prefilled letter/block must be
// preserved to the given solver.
func (c *Constraints) AddInputGridConstraintsAreSatisfiedClausesTo(solverConfigurer solver.Configurer) {
	for row := 0; row < c.grid.RowCount(); row++ {
		for column := 0; column < c.grid.ColumnCount(); column++ {
			prefilledLetter := c.grid.LetterAt(row, column)
			var literal solver.Literal
			if prefilledLetter == CellEmpty {
				// Disallow solver to create a block
				literal = solver.Literal(c.variables.RepresentingCell(row, column, BlockIndex())).Negated()
			} else if prefilledLetter == CellBlock {
				literal = solver.Literal(c.variables.RepresentingCell(row, column, BlockIndex()))
			} else {
				letterIndex, _ := alphabet.IndexOf(prefilledLetter)
				literal = solver.Literal(c.variables.RepresentingCell(row, column, letterIndex))
			}
			solverConfigurer.AddClause([]solver.Literal{literal})
		}
	}
}
