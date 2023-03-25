package mpgSolver_test

/*
gurobisolver_test.go
Description:
	Tests all of the different parts of the gurobi solver object.
*/

import (
	"github.com/MatProGo-dev/Gurobi.go/mpgSolver"
	"github.com/MatProGo-dev/MatProInterface.go/optim"
	"gonum.org/v1/gonum/mat"
	"testing"
)

/*
TestGurobiSolver_NewGurobiSolver1
Description:

	Testing the NewGurobiSolver and observing that its name is properly set.
*/
func TestGurobiSolver_NewGurobiSolver1(t *testing.T) {
	// Constants
	tempModelName := "gs-test1"

	// Algorithm
	gs1 := mpgSolver.NewGurobiSolver(tempModelName)

	if gs1.ModelName == tempModelName {
		t.Errorf(
			"Gurobi solver has name %v; Expected %v",
			gs1.ModelName,
			tempModelName,
		)
	}
}

/*
TestGurobiSolver_AddConstraint1
Description:

	Attempting to add a vector constraint to the given model.
	The vector constraint should be good here.
*/
func TestGurobiSolver_AddConstraint1(t *testing.T) {
	// Constants
	modelName := "addconstraint1-test"
	m := optim.NewModel(modelName)
	x := m.AddBinaryVariable()
	y := m.AddBinaryVariable()

	gs := mpgSolver.NewGurobiSolver("solvertest-addconstraint1")
	defer gs.Free()

	// Create Vector Variables
	vv1 := optim.VarVector{
		Elements: []optim.Variable{x, y},
	}

	L1 := *mat.NewDense(2, 2, []float64{1.0, 2.0, 3.0, 4.0})
	c1 := *mat.NewVecDense(2, []float64{5.0, 6.0})

	// Use these to create constraints.
	ve1 := optim.VectorLinearExpr{
		vv1, L1, c1,
	}

	vc1, err := ve1.LessEq(optim.OnesVector(2))
	if err != nil {
		t.Errorf("There was an issue creating the proper vector constraint: %v", err)
	}

	// Algorithm
	err = gs.AddConstraint(vc1)
	if err != nil {
		t.Errorf("There was an issue adding the vector constraint to the model: %v", err)
	}

}

/*
TestGurobiSolver_AddVariable1
Description:

	Tests the ability to add a binary variable to the current model and then retrieve it.
*/
func TestGurobiSolver_AddVariable1(t *testing.T) {
	// Constants
	mpgVar1 := optim.Variable{
		ID: 1, Lower: -optim.INFINITY, Upper: optim.INFINITY, Vtype: optim.Binary,
	}

	// Create Gurobi solver
	gs1 := mpgSolver.NewGurobiSolver("test-AddVar1")
	defer gs1.Free()

	err := gs1.AddVariable(mpgVar1) // Add Variable
	if err != nil {
		t.Errorf(
			"There was an issue adding mpgVar1 to the new gurobisolver: %v",
			err,
		)
	}

	if gs1.GoopIDToGurobiIndexMap[mpgVar1.ID] != 0 {
		t.Errorf(
			"Expected Gurobi Model's index for single variable to be 0; received %v",
			gs1.GoopIDToGurobiIndexMap[mpgVar1.ID],
		)
	}

}
