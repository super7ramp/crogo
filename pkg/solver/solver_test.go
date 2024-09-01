package solver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testingSolverConfigurer struct {
	*BaseConfigurer
	clauses [][]Literal
}

func (c *testingSolverConfigurer) AddClause(literals []Literal) {
	c.clauses = append(c.clauses, literals)
}

func newTestingSolverConfigurer() *testingSolverConfigurer {
	baseConfigurer := BaseConfigurer{}
	solverConfigurer := testingSolverConfigurer{BaseConfigurer: &baseConfigurer}
	baseConfigurer.Configurer = &solverConfigurer
	return &solverConfigurer
}

func TestAddExactlyOne(t *testing.T) {
	solverConfigurer := newTestingSolverConfigurer()
	solverConfigurer.AddExactlyOne([]Literal{1, 2, 3})
	assert.Equal(t, [][]Literal{{1, 2, 3}, {-1, -2}, {-1, -3}, {-2, -3}}, solverConfigurer.clauses)
}

func TestAddAtMostOne(t *testing.T) {
	solverConfigurer := newTestingSolverConfigurer()
	solverConfigurer.AddAtMostOne([]Literal{1, 2, 3})
	assert.Equal(t, [][]Literal{{-1, -2}, {-1, -3}, {-2, -3}}, solverConfigurer.clauses)
}

func TestAddAnd(t *testing.T) {
	solverConfigurer := newTestingSolverConfigurer()
	solverConfigurer.AddAnd(42, []Literal{-1, 6, -7})
	assert.Equal(t, [][]Literal{{-42, -1}, {-42, 6}, {-42, -7}, {1, -6, 7, 42}}, solverConfigurer.clauses)
}
