package gurobi

// Gurobi linear constraint object
type Constr struct {
	Model *Model
	Index int32
}

/*
VectorConstraintToGurobiSparseFormat
Description:
	Converts a vector constraint of the form L * x <= rhs
	into gurobi's format for AddVars (cbeg,
*/
//func VectorConstraintToGurobiSparseFormat(
//	L [][]float64{}, rhs []float64,
//	) ( ) {
//
//}
