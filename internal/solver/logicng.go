package solver

import (
	spi "crogo/pkg/solver"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/sat"
	"iter"
	"strconv"
)

type logicNgSolver struct {
	*spi.BaseConfigurer
	satSolver *sat.Solver
}

// NewLogicNgSolver creates a new instance of a spi.ConfigurableSolver based on Gophersat.
func NewLogicNgSolver() spi.ConfigurableSolver {
	formulaFactory := formula.NewFactory()
	satSolver := sat.NewSolver(formulaFactory)
	baseConfigurer := spi.BaseConfigurer{}
	solverConfigurer := logicNgSolver{BaseConfigurer: &baseConfigurer, satSolver: satSolver}
	baseConfigurer.Configurer = &solverConfigurer
	return &solverConfigurer
}

func (l logicNgSolver) AddClause(spiLiterals []spi.Literal) {
	formulaFactory := l.satSolver.Factory()
	literals := make([]formula.Literal, len(spiLiterals))
	for i, spiLiteral := range spiLiterals {
		literals[i] = formulaFactory.Lit(strconv.Itoa(spiLiteral), spiLiteral > 0)
	}
	clause := formulaFactory.Clause(literals...)
	l.satSolver.Add(clause)
}

// TODO override addExactlyOneClause with PB clause

func (l logicNgSolver) Solutions() iter.Seq[spi.Model] {
	return func(yield func(spi.Model) bool) {
		// TODO actually enumerate all solutions
		if l.satSolver.Sat() {
			model := l.satSolver.CoreSolver().Model()
			yield(model)
		}
	}
}
