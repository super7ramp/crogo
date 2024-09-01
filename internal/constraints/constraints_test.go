package constraints

import (
	spi "crogo/pkg/solver"
	"slices"
)

type testingSolverConfigurer struct {
	*spi.BaseConfigurer
	clauses           [][]spi.Literal
	exactlyOneClauses [][]spi.Literal
	andClauses        map[spi.Literal][]spi.Literal
}

func (c *testingSolverConfigurer) AddClause(literals []spi.Literal) {
	c.clauses = append(c.clauses, slices.Clone(literals))
}

func (c *testingSolverConfigurer) AddExactlyOneClause(literals []spi.Literal) {
	c.exactlyOneClauses = append(c.exactlyOneClauses, slices.Clone(literals))
}

func (c *testingSolverConfigurer) AddAndClause(literal spi.Literal, conjunction []spi.Literal) {
	c.andClauses[literal] = slices.Clone(conjunction)
}

func newTestingSolverConfigurer() *testingSolverConfigurer {
	baseConfigurer := spi.BaseConfigurer{}
	solverConfigurer := testingSolverConfigurer{BaseConfigurer: &baseConfigurer}
	baseConfigurer.Configurer = &solverConfigurer
	return &solverConfigurer
}
