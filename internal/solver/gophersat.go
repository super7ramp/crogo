package solver

import (
	spi "crogo/pkg/solver"
	sat "github.com/crillab/gophersat/solver"
	"iter"
)

type gophersatSolver struct {
	*spi.BaseConfigurer
	satSolver *sat.Solver
}

// NewGophersatSolver creates a new instance of a spi.ConfigurableSolver based on Gophersat.
func NewGophersatSolver() spi.ConfigurableSolver {
	problem := sat.Problem{}
	satSolver := sat.New(&problem)
	baseConfigurer := spi.BaseConfigurer{}
	solverConfigurer := gophersatSolver{BaseConfigurer: &baseConfigurer, satSolver: satSolver}
	baseConfigurer.Configurer = &solverConfigurer
	return &solverConfigurer
}

// AddClause adds the given literals as an *at-least-one* clause, i.e. a disjunction (= or).
func (s *gophersatSolver) AddClause(spiLiterals []spi.Literal) {
	literals := gophersatLitsFrom(spiLiterals...)
	clause := sat.NewClause(literals)
	s.satSolver.AppendClause(clause)
}

// TODO override addExactlyOneClause with PB clause

func gophersatLitsFrom(vals ...spi.Literal) []sat.Lit {
	// That's a useless copy...
	res := make([]sat.Lit, len(vals))
	for i, val := range vals {
		res[i] = sat.IntToLit(int32(val))
	}
	return res
}

// Solutions returns an iterator on the solutions.
func (s *gophersatSolver) Solutions() iter.Seq[spi.Model] {
	return func(yield func(spi.Model) bool) {
		// TODO actually enumerate all solutions
		status := s.satSolver.Solve()
		if status == sat.Sat {
			yield(s.satSolver.Model())
		}
	}
}
