package solver

import (
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/sat"
	"iter"
	"slices"
)

type logicNgSolver struct {
	*BaseConfigurer
	satSolver         *sat.Solver
	relevantVariables []formula.Variable
}

// NewLogicNgSolver creates a new instance of a spi.ConfigurableSolver based on LogicNg.
func NewLogicNgSolver() ConfigurableSolver {
	formulaFactory := formula.NewFactory()
	satSolver := sat.NewSolver(formulaFactory)
	baseConfigurer := BaseConfigurer{}
	solverConfigurer := logicNgSolver{BaseConfigurer: &baseConfigurer, satSolver: satSolver}
	baseConfigurer.Configurer = &solverConfigurer
	return &solverConfigurer
}

func (l *logicNgSolver) SetRelevantVariables(variables []Variable) {
	l.relevantVariables = make([]formula.Variable, len(variables))
	formulaFactory := l.satSolver.Factory()
	for i, variable := range variables {
		l.relevantVariables[i] = formulaFactory.Var(variable.String())
	}
}

func (l *logicNgSolver) AddClause(spiLiterals []Literal) {
	literals := l.logicNgLitsFrom(spiLiterals)
	clause := l.satSolver.Factory().Clause(literals...)
	l.satSolver.Add(clause)
}

func (l *logicNgSolver) logicNgLitsFrom(spiLiterals []Literal) []formula.Literal {
	formulaFactory := l.satSolver.Factory()
	literals := make([]formula.Literal, len(spiLiterals))
	for i, spiLiteral := range spiLiterals {
		spiVariable := VariableFrom(spiLiteral)
		literals[i] = formulaFactory.Lit(spiVariable.String(), spiLiteral > 0)
	}
	return literals
}

func (l *logicNgSolver) AddExactlyOne(spiLiterals []Literal) {
	literals := l.logicNgLitsFrom(spiLiterals)
	clause := l.satSolver.Factory().PBC(formula.EQ, 1, literals, slices.Repeat([]int{1}, len(literals)))
	l.satSolver.Add(clause)
}

func (l *logicNgSolver) Solutions() iter.Seq[Model] {
	return func(yield func(Model) bool) {
		for {
			result := l.satSolver.Call(sat.WithModel(l.relevantVariables))
			if !result.OK() || !result.Sat() {
				break
			}

			model := result.Model().Literals
			adaptedModel := make([]bool, len(model))
			for i, lit := range model {
				adaptedModel[i] = lit.IsPos()
			}
			if keepGoing := yield(adaptedModel); !keepGoing {
				break
			}

			differentModel := make([]formula.Literal, len(model))
			factory := l.satSolver.Factory()
			for i, lit := range model {
				differentModel[i] = lit.Negate(factory)
			}
			differentModelClause := factory.Clause(differentModel...)
			l.satSolver.Add(differentModelClause)
		}
	}
}
