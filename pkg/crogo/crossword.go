package crogo

import (
	. "crogo/internal/constraints"
	. "crogo/internal/grid"
	. "crogo/internal/solver"
	. "crogo/internal/variables"
	"crogo/pkg/solver"
	"iter"
)

// Crossword is the crossword structure, holding variables and constraints information.
type Crossword struct {
	variables   *Variables
	constraints *Constraints
}

// Solutions is an iterator over crossword solutions.
type Solutions = iter.Seq[[][]rune]

// NewCrossword constructs a new instance of Crossword.
func NewCrossword(cells [][]rune, words []string) (*Crossword, error) {
	grid, err := NewGrid(cells)
	if err != nil {
		return nil, err
	}
	variables := NewVariables(grid, len(words))
	constraints := NewConstraints(grid, variables, words)
	return &Crossword{variables, constraints}, nil
}

// Solve solves this crossword using builtin solver.
func (c *Crossword) Solve() Solutions {
	//defaultSolver := NewGophersatSolver()
	defaultSolver := NewLogicNgSolver()
	return c.SolveWith(defaultSolver)
}

// SolveWith solves this crossword using the given solver.
func (c *Crossword) SolveWith(configurableSolver solver.ConfigurableSolver) Solutions {
	c.addClausesTo(configurableSolver)
	return c.solutions(configurableSolver)
}

// addClausesTo adds clauses to the given solver configurer.
func (c *Crossword) addClausesTo(solverConfigurer solver.Configurer) {
	solverConfigurer.AllocateVariables(uint(c.variables.Count()))
	solverConfigurer.SetRelevantVariables(c.variables.RepresentingCells())
	c.constraints.AddOneLetterOrBlockPerCellClausesTo(solverConfigurer)
	c.constraints.AddOneWordPerSlotClausesTo(solverConfigurer)
	c.constraints.AddInputGridConstraintsAreSatisfiedClausesTo(solverConfigurer)
}

func (c *Crossword) solutions(s solver.Solver) Solutions {
	return func(yield func([][]rune) bool) {
		adaptedYield := func(model solver.Model) bool {
			return yield(c.variables.BackToDomain(model))
		}
		s.Solutions()(adaptedYield)
	}
}
