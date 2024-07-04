package main

import "fmt"

/*
lp2.go
Description:
	This file creates a linear program that finds a point on the edges of a polytope.
*/

import (
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
	"math"
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

func main() {
	// Create Constants
	AData := []float64{
		1, 0,
		0, 1,
		-1, 0,
		0, -1,
	}
	A := mat.NewDense(4, 2, AData)

	bData := []float64{2, 1, 0, 0}
	//b := mat.NewVecDense(4, bData)

	cData := []float64{-1, -1}
	//c := mat.NewVecDense(2, cData)

	// Define Equality Constraints
	// GData := []float64{0, 0}
	// G := mat.NewDense(1, 2, GData)

	// hbuf := []float64{0}
	// h := mat.NewVecDense(0, hbuf)

	// Solve!
	optValGonum, xOptGonum, err := SolveWithGonum(*A, bData, cData)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Optimal Value: %v\n", optValGonum)
	fmt.Printf("Optimal Point: %v\n", xOptGonum)

}
