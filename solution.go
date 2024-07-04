package gonumWrapper

/*
Solution
Description:

	This struct is used to store the results of a solution
	to an optimization problem.
*/
type Solution struct {
	ObjectiveValue float64
	Variables      []float64
}
