package gwErrors

import (
	"fmt"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
)

/*
nonlinear_constraint.go
Description:

	Defines an error that is thrown when a nonlinear constraint
	is used in a context where only linear constraints are allowed.
*/

// Type
// ====

type NonlinearConstraintError struct {
	Constraint symbolic.Constraint
}

// Methods
// =======

func (nce NonlinearConstraintError) Error() string {
	// Setup

	// Compute degree
	lhs, err := symbolic.ToPolynomialLike(nce.Constraint.Left())
	if err != nil {
		return fmt.Sprintf(
			"nonlinear constraint (type %T) used in context where only linear constraints are allowed.",
			nce.Constraint,
		)
	}
	rhs, err := symbolic.ToPolynomialLike(nce.Constraint.Right())
	if err != nil {
		return fmt.Sprintf(
			"nonlinear constraint (type %T) used in context where only linear constraints are allowed.",
			nce.Constraint,
		)
	}

	degree := max(lhs.Degree(), rhs.Degree())

	return fmt.Sprintf(
		"nonlinear constraint (type %T) with degree %d used in context where only linear constraints are allowed.",
		nce.Constraint,
		degree,
	)
}
