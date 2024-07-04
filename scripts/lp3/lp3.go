package main

/*
lp3.go
Description:
	This file creates a linear program that finds a point on the edges of a polytope.
*/

import (
	"fmt"
	problem "github.com/MatProGo-dev/MatProInterface.go/problem"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
	"math"
	gonumWrapper "matprogo.dev/wrappers/gonum"
)

func SolveWithGonum(A mat.Dense, bData []float64, cData []float64) (float64, []float64, error) {
	// Constants
	tol := 0.1
	var_cnt := len(cData)

	// Solve!
	cNew, aNew, bNew := lp.Convert(cData, &A, bData, nil, nil)
	opt, xNew, lp_err := lp.Simplex(cNew, aNew, bNew, tol, nil)

	// Print Results
	if lp_err != nil {
		fmt.Printf("LP error: %v\n", lp_err)
		// return
		return 0, nil, lp_err
	}

	// Collect Results
	fmt.Printf("[gonum] LP Success.\n")
	fmt.Printf("[gonum] Optimal Objective Value: %v\n", opt)

	x := make([]float64, var_cnt)
	for i := 0; i < var_cnt; i++ {
		x[i] = math.RoundToEven(xNew[i] - xNew[i+var_cnt])
	}
	fmt.Printf("%v", x)

	return opt, x, nil
}

func SolveWithMatProGo(A mat.Dense, bData []float64, cData []float64) (float64, []float64, error) {
	// Create variables
	prob := problem.NewProblem("lp1-mpg")
	x := prob.AddVariableVector(2)

	// Set up constraints
	con1 := symbolic.DenseToKMatrix(A).Multiply(x).LessEq(bData[0])

	prob.Constraints = append(prob.Constraints, con1)

	// Set up objective
	c := getKVector.From(cData)
	newObjective := problem.NewObjective(
		x.Transpose().Multiply(c),
		problem.SenseMinimize,
	)
	prob.Objective = *newObjective

	// Solve problem
	solution, err := gonumWrapper.Solve(*prob)
	if err != nil {
		return 0, nil, err
	}

	// Return
	return solution.ObjectiveValue, solution.Variables, nil
}

func main() {
	// Create Constants
	AData := []float64{
		1, 1,
	}
	A := mat.NewDense(1, 2, AData)

	bData := []float64{2}
	//b := mat.NewVecDense(4, bData)

	cData := []float64{-1, -1}
	//c := mat.NewVecDense(2, cData)

	// Solve!
	optValGonum, xOptGonum, err := SolveWithGonum(*A, bData, cData)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Optimal Value: %v\n", optValGonum)
	fmt.Printf("Optimal Point: %v\n", xOptGonum)

	// Solve using the symbolic math library
	optValMatProGo, xOptMatProGo, err := SolveWithMatProGo(*A, bData, cData)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Optimal Value: %v\n", optValMatProGo)
	fmt.Printf("Optimal Point: %v\n", xOptMatProGo)

}
