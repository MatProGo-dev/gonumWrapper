package gonumWrapper

import (
	"fmt"
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	"matprogo.dev/wrappers/gonum/gwErrors"
)

/*
solver.go
Description:

	This file defines a solver function that uses the gonum library
	to assemble a linear program and solve it.
*/

/*
SimplexSolver
Description:

		This struct is used to represent a linear program that can be solved
		using the gonum library's simplex solver.
		It represents the problem in the format:
			minimize cᵀ * x
	 		s.t      G * x <= h
	          		 A * x = b
*/
type SimplexSolver struct {
	X           []symbolic.Variable
	G           mat.Matrix
	H           []float64
	A           mat.Matrix
	B           []float64
	C           []float64
	Tolerance   float64
	Env         symbolic.Environment
	Constraints []symbolic.Constraint
	VariableMap map[symbolic.Variable]symbolic.ScalarExpression // Internal variable map used for representing the input variables.
}

///*
//DecomposeVariables
//Description:
//
//	This method decomposes the input variables into positive and negative parts.
//*/
//func (ss *SimplexSolver) DecomposeVariables() {
//	// Setup
//	ss.VariableMap = make(map[symbolic.Variable]symbolic.ScalarExpression)
//
//	// Iterate over the variables
//	var replacement symbolic.ScalarExpression = symbolic.K(0.0)
//	for _, variable := range ss.X {
//		// If there is a positive part, then add it.
//		if variable.Upper > 0 {
//			positivePart := symbolic.Variable{
//				Name:  variable.Name + "_plus",
//				Lower: 0,
//				Upper: variable.Upper,
//				Type:  variable.Type,
//			}
//			ss.Env.Variables = append(ss.Env.Variables, positivePart)
//			replacement = replacement.Plus(positivePart).(symbolic.ScalarExpression)
//		}
//		// If there is a negative part, then add it.
//		if variable.Lower < 0 {
//			negativePart := symbolic.Variable{
//				Name:  variable.Name + "_minus",
//				Lower: 0,
//				Upper: -variable.Lower,
//				Type:  variable.Type,
//			}
//			ss.Env.Variables = append(ss.Env.Variables, negativePart)
//			replacement = replacement.Minus(negativePart).(symbolic.ScalarExpression)
//		}
//		ss.VariableMap[variable] = replacement
//	}
//}

/*
AddConstraint
Description:

	This method adds a constraint to the SimplexSolver object.
	Note that the constraint should take form of either:
	1. A linear equality constraint
	2. A positivity constraint on a variable
*/
func (ss *SimplexSolver) AddConstraint(constraint symbolic.Constraint) {
	// If this is a positivity constraint on a variable, then we can add it directly.
	//if IsPositiveVariableConstraint(constraint) {
	//	// Add the constraint to the list of constraints
	//	ss.Constraints = append(ss.Constraints, constraint)
	//	return
	//}

	// Algorithm
	switch con0 := constraint.(type) {
	case symbolic.ScalarConstraint:
		simplified := con0.Simplify()

		if !simplified.IsLinear() {
			panic(
				gwErrors.NonlinearConstraintError{
					Constraint: simplified,
				},
			)
		}

		// Create Constraint based on what we have so far.
		switch con0.ConstrSense() {
		case symbolic.SenseLessThanEqual:
			// Add slack variable to constraint
			left := con0.Left()
			right := con0.Right()

			sum := left.Minus(right).(symbolic.ScalarExpression)
			//slack := symbolic.Variable{
			//	Name:  "slack_" + fmt.Sprintf("%d", len(ss.Constraints)),
			//	Lower: 0,
			//	Upper: float64(symbolic.Infinity),
			//	Type:  symbolic.Continuous,
			//}

			// Add slack variable to the environment
			// ss.Env.Variables = append(ss.Env.Variables, slack)

			// Add the constraint to the list of constraints
			newConstraint := sum.Minus(sum.Constant()).LessEq(sum.Constant()).(symbolic.Constraint)
			//ss.Constraints = append(ss.Constraints, newConstraint)

			// Update the A matrix
			left_polynomial_like, _ := newConstraint.Left().(symbolic.PolynomialLikeScalar)
			Grow := left_polynomial_like.LinearCoeff(ss.X)
			GrowAsSlice := make([]float64, len(ss.X))
			for ii := range ss.X {
				GrowAsSlice[ii] = Grow.AtVec(ii)
			}
			if ss.G == nil {
				ss.G = mat.NewDense(1, len(GrowAsSlice), GrowAsSlice)
				ss.H = []float64{sum.Constant()}
			} else {
				ss.G.(*mat.Dense).Stack(
					mat.DenseCopyOf(ss.G),
					mat.NewDense(1, len(GrowAsSlice), GrowAsSlice),
				)
				ss.H = append(ss.H, sum.Constant())
			}

		}
	default:
		panic(
			fmt.Errorf(
				"unexpected type of constraint (%T) provided to SimplexSolver.AddConstraint() method.",
				constraint,
			),
		)
	}
}

/*
AddObjective
Description:

	This method adds an objective to the SimplexSolver object.
*/
func (ss *SimplexSolver) AddObjective(objective problem.Objective) {
	// Setup

	// Algorithm
	objAsPolynomialLike, _ := objective.Expression.(symbolic.PolynomialLikeScalar)
	C := objAsPolynomialLike.LinearCoeff(ss.X)

	CAsSlice := make([]float64, len(ss.X))
	for ii := range ss.X {
		switch objective.Sense {
		case problem.SenseMinimize:
			CAsSlice[ii] = C.AtVec(ii)
		case problem.SenseMaximize:
			CAsSlice[ii] = -C.AtVec(ii)
		default:
			panic(
				fmt.Errorf(
					"unexpected sense (%v) provided to SimplexSolver.AddObjective() method.",
					objective.Sense,
				),
			)
		}

	}

	ss.C = CAsSlice
}
