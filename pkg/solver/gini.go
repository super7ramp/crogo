package solver

import (
	"github.com/go-air/gini"
	"github.com/go-air/gini/z"
	"iter"
)

type giniSolver struct {
	*BaseConfigurer
	backend           *gini.Gini
	relevantVariables []z.Var
}

// NewGiniSolver creates a new instance of a spi.ConfigurableSolver based on Gini.
func NewGiniSolver() ConfigurableSolver {
	backend := gini.New()
	baseConfigurer := BaseConfigurer{}
	solverConfigurer := giniSolver{BaseConfigurer: &baseConfigurer, backend: backend}
	baseConfigurer.Configurer = &solverConfigurer
	return &solverConfigurer
}

func (g *giniSolver) AddClause(spiLiterals []Literal) {
	for _, spiLiteral := range spiLiterals {
		g.backend.Add(z.Dimacs2Lit(int(spiLiteral)))
	}
	g.backend.Add(0)
}

func (g *giniSolver) SetRelevantVariables(variables []Variable) {
	g.relevantVariables = make([]z.Var, len(variables))
	for i, variable := range variables {
		g.relevantVariables[i] = z.Var(variable)
	}
}

func (g *giniSolver) Solutions() iter.Seq[Model] {
	return func(yield func(Model) bool) {
		for {
			if res := g.backend.Solve(); res != 1 {
				break
			}
			adaptedModel := make([]bool, len(g.relevantVariables))
			for i, variable := range g.relevantVariables {
				adaptedModel[i] = g.backend.Value(variable.Pos())
			}
			if keepGoing := yield(adaptedModel); !keepGoing {
				break
			}
			for i, isPos := range adaptedModel {
				g.backend.Add(boolToGiniLit(i+1, isPos).Not())
			}
			g.backend.Add(0)
		}
	}
}

func boolToGiniLit(variable int, isPos bool) z.Lit {
	var lit z.Lit
	if isPos {
		lit = z.Var(variable).Pos()
	} else {
		lit = z.Var(variable).Neg()
	}
	return lit
}
