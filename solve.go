package gonumWrapper

import (
	"fmt"
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"gonum.org/v1/gonum/optimize/convex/lp"
)

/*
solve.go
Description:
	This file contains all of the functions needed
	to support the solving of an optimization problem given
	in the optimization problem.
*/

/*
Solve()
Description:

	This function attempts to solve the model given in model
	using an available solver in gonum.
*/
func Solve(p problem.OptimizationProblem) (Solution, error) {
	// Input Processing + Setup
	//TODO: Introduce problem.Check() method.

	// Create the solver object
	solver := SimplexSolver{
		X:         p.Variables,
		Tolerance: 0.1,
		A:         nil,
		B:         nil,
	}

	// Add Constraints to solver
	for _, con := range p.Constraints {
		solver.AddConstraint(con)
	}

	// Add Objective to solver
	solver.AddObjective(p.Objective)

	// Solve the problem using gonum
	fmt.Println("Solving with gonum...")
	cNew, aNew, bNew := lp.Convert(solver.C, solver.G, solver.H, solver.A, solver.B)
	opt, xNew, lp_err := lp.Simplex(cNew, aNew, bNew, solver.Tolerance, nil)
	if lp_err != nil {
		return Solution{}, lp_err
	}

	// Convert constraints to matrices
	return Solution{
		ObjectiveValue: opt,
		Variables:      xNew,
	}, nil
}
