package solver

import (
	spi "crogo/pkg/solver"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/sat"
	"iter"
)

type logicNgSolver struct {
	*spi.BaseConfigurer
	satSolver         *sat.Solver
	relevantVariables []formula.Variable
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

func (l *logicNgSolver) SetRelevantVariables(variables []spi.Variable) {
	l.relevantVariables = make([]formula.Variable, len(variables))
	formulaFactory := l.satSolver.Factory()
	for i, variable := range variables {
		l.relevantVariables[i] = formulaFactory.Var(variable.String())
	}
}

func (l *logicNgSolver) AddClause(spiLiterals []spi.Literal) {
	literals := l.logicNgLitsFrom(spiLiterals)
	clause := l.satSolver.Factory().Clause(literals...)
	l.satSolver.Add(clause)
}

func (l *logicNgSolver) logicNgLitsFrom(spiLiterals []spi.Literal) []formula.Literal {
	formulaFactory := l.satSolver.Factory()
	literals := make([]formula.Literal, len(spiLiterals))
	for i, spiLiteral := range spiLiterals {
		spiVariable := spi.VariableFrom(spiLiteral)
		literals[i] = formulaFactory.Lit(spiVariable.String(), spiLiteral > 0)
	}
	return literals
}

// TODO override addExactlyOneClause with PB clause

func (l *logicNgSolver) Solutions() iter.Seq[spi.Model] {
	return func(yield func(spi.Model) bool) {
		// TODO actually enumerate all solutions
		result := l.satSolver.Call(
			sat.WithModel(l.relevantVariables),
		)
		if result.OK() && result.Sat() {
			model := result.Model().Literals
			adaptedModel := make([]bool, len(model))
			for i, lit := range model {
				adaptedModel[i] = lit > 0
			}
			yield(adaptedModel)
		}
	}
}
