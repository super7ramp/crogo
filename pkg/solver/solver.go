package solver

import "iter"

// Variable is a variable number. 0 is not a valid variable number, variables start at 1.
type Variable = uint

// Literal is an instantiated variable, i.e. it is the variable number if the variable is true, or the negative of the
// variable number if the variable is false. Since variable number starts at 1, 0 is not a valid literal.
type Literal = int

// Model is a slice of the variable states, indexed by the variable number.
type Model = []bool

// Solver defines a SAT solver.
//
// It is an iterator over the models satisfying the problem. A model is a slice indexed by the variables, whose values
// indicate the state of the corresponding variable.
//
// Implementation *may* return only the state of relevant variables defined by Configurer.SetRelevantVariables instead
// of all the variables of the problems.
type Solver interface {
	// Solutions returns an iterator on the solutions.
	Solutions() iter.Seq[Model]
}

// Configurer defines a solver configurer.
type Configurer interface {
	// AllocateVariables gives a hint about the number of variables.
	AllocateVariables(variableCount uint)

	// SetRelevantVariables indicates which variables are relevant for the problem.
	SetRelevantVariables(variables []Variable)

	// AddClause adds the given literals as an *at-least-one* clause, i.e. a disjunction (= or).
	AddClause(literals []Literal)

	// AddExactlyOne adds the given literals as an *exactly-one* clause.
	//
	// An *exactly-one* clause is equivalent to an *at-least-one* and a *at-most-one* clauses.
	AddExactlyOne(literals []Literal)

	// AddAtMostOne adds the given literals as an *at-most-one* clause.
	AddAtMostOne(literals []Literal)

	// AddAnd adds clauses describing the equivalence between the given literal and the given conjunction
	// (= and) of literals, i.e.: *literal ⇔ conjunction\[0\] ∧ conjunction\[1\] ∧ ... ∧ conjunction\[n\]*
	//
	// The corresponding clauses are: *(￢literal ∨ conjunction\[0\]) ∧
	// (￢literal ∨ conjunction\[1\]) ∧ ... ∧ (￢literal ∨ conjunction\[1\]) ∧ (￢conjunction\[0\]
	// ∨ ￢conjunction\[1\] ∨ ... ∨ ￢conjunction\[n\] ∨ literal)*
	AddAnd(literal Literal, literals []Literal)
}

// ConfigurableSolver defines a configurable Solver.
type ConfigurableSolver interface {
	Solver
	Configurer
}

// BaseConfigurer provides default implementations for all the functions of the Configurer interface but for the
// Configurer.AddClause function. These default implementations may be overridden for better performances.
type BaseConfigurer struct {
	Configurer
}

func (c *BaseConfigurer) AllocateVariables(_ uint) {
	// Do nothing.
}

func (c *BaseConfigurer) SetRelevantVariables(_ []Variable) {
	// Do nothing.
}

func (c *BaseConfigurer) AddExactlyOne(literals []Literal) {
	c.AddClause(literals)
	c.AddAtMostOne(literals)
}

func (c *BaseConfigurer) AddAtMostOne(literals []Literal) {
	for i := range literals {
		for j := i + 1; j < len(literals); j++ {
			c.AddClause([]Literal{-literals[i], -literals[j]})
		}
	}
}

func (c *BaseConfigurer) AddAnd(literal Literal, conjunction []Literal) {
	lastClause := make([]Literal, 0, len(conjunction)+1)
	for _, conjunctionLiteral := range conjunction {
		c.AddClause([]Literal{-literal, conjunctionLiteral})
		lastClause = append(lastClause, -conjunctionLiteral)
	}
	lastClause = append(lastClause, literal)
	c.AddClause(lastClause)
}
