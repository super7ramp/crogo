package solver

import (
	spi "crogo/pkg/solver"
	"github.com/go-air/gini"
	"github.com/go-air/gini/z"
	"iter"
)

type giniSolver struct {
	*spi.BaseConfigurer
	backend           *gini.Gini
	relevantVariables []z.Var
}

// NewGiniSolver creates a new instance of a spi.ConfigurableSolver based on Gini.
func NewGiniSolver() spi.ConfigurableSolver {
	backend := gini.New()
	baseConfigurer := spi.BaseConfigurer{}
	solverConfigurer := giniSolver{BaseConfigurer: &baseConfigurer, backend: backend}
	baseConfigurer.Configurer = &solverConfigurer
	return &solverConfigurer
}

func (g *giniSolver) AddClause(spiLiterals []spi.Literal) {
	for _, spiLiteral := range spiLiterals {
		g.backend.Add(z.Dimacs2Lit(int(spiLiteral)))
	}
	g.backend.Add(0)
}

func (g *giniSolver) SetRelevantVariables(variables []spi.Variable) {
	g.relevantVariables = make([]z.Var, len(variables))
	for i, variable := range variables {
		g.relevantVariables[i] = z.Var(variable)
	}
}

func (g *giniSolver) Solutions() iter.Seq[spi.Model] {
	return func(yield func(spi.Model) bool) {
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
