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

func (s *gophersatSolver) Solutions() iter.Seq[spi.Model] {
	return func(yield func(spi.Model) bool) {
		models := make(chan []bool)
		stop := make(chan struct{}, 1) // FIXME Enumerate does not read from stop channel?
		defer close(stop)
		go s.satSolver.Enumerate(models, stop)
		for model := range models {
			if keepGoing := yield(model); !keepGoing {
				stop <- struct{}{}
				break
			}
		}
	}
}
