package mpgSolver_test

/*
gurobisolver_test.go
Description:
	Tests all of the different parts of the gurobi solver object.
*/

import (
	"fmt"
	"github.com/MatProGo-dev/Gurobi.go/gurobi"
	"github.com/MatProGo-dev/Gurobi.go/mpgSolver"
	"github.com/MatProGo-dev/MatProInterface.go/optim"
	"gonum.org/v1/gonum/mat"
	"os"
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
	defer os.Remove(gs1.ModelName + ".log")

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
	defer os.Remove(gs.ModelName + ".log")
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

	// Add variables to
	err = gs.AddVariables(vv1.Elements)
	if err != nil {
		t.Errorf("unexpected issue adding variables to gurobi solver's model: %v", err)
	}

	// Test Adding Constraints
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
	defer os.Remove(gs1.ModelName + ".log")
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

/*
TestGurobiSolver_Solve1
Description:

	Tests that the solver's Solve() method work's correctly.
*/
func TestGurobiSolver_Solve1(t *testing.T) {
	// Constants
	exampleName := "testgurobisolver-solve1"

	// Create environment.
	env, err := gurobi.NewEnv(exampleName + ".log")
	if err != nil {
		panic(err.Error())
	}
	defer env.Free()

	// Create an empty model.
	model := optim.NewModel(exampleName + ".model")

	// Add varibles
	x := model.AddVariableVector(2)

	vd1 := mat.NewVecDense(2, []float64{0.0, 1.0})
	vd2 := mat.NewVecDense(2, []float64{2.0, 3.0})

	// Create constraints
	err = model.AddConstraint(x.GreaterEq(*vd1))
	if err != nil {
		t.Errorf(
			"there was an issue creating the proper vector constraint: %v",
			err,
		)
	}

	err = model.AddConstraint(x.LessEq(*vd2))
	if err != nil {
		t.Errorf(
			"there was an issue creating the proper vector constraint: %v",
			err,
		)
	}

	// Set Objective function
	Q1 := *mat.NewDense(2, 2, []float64{1.0, 0.25, 0.25, 0.25})

	prod1, _ := x.Transpose().Multiply(Q1)
	prod2, err := prod1.Multiply(x)
	if err != nil {
		t.Errorf(
			"there was an issue creating product 2: %v",
			err,
		)
	}

	sum, err := prod2.Plus(
		x.Transpose().Multiply(*mat.NewVecDense(2, []float64{0, -0.97})),
	)
	if err != nil {
		t.Errorf(
			"There was an issue computing the final sum: %v",
			err,
		)
	}

	err = model.SetObjective(sum, optim.SenseMinimize)
	if err != nil {
		t.Errorf(
			"There was an issue setting the objective for the model: %v",
			err,
		)
	}

	// Solve the above model with GurobiSolver
	sol, gs, err := mpgSolver.Solve(*model)
	if err != nil {
		t.Errorf("error in solving model: %v", err)
	}

	//// Write model to 'qp.lp'.
	fmt.Printf("%v\n%v\n%v\n", gs.GoopIDToGurobiIndexMap, x.Elements[0].ID, x.Elements[1].ID)
	if err := gs.CurrentModel.Write(exampleName + ".lp"); err != nil {
		panic(err.Error())
	}
	defer os.Remove(exampleName + ".lp")

	// Capture solution information
	optimstatus := sol.Status
	objval := sol.Objective

	fmt.Printf("\nOptimization complete\n")
	if optimstatus != gurobi.OPTIMAL {
		t.Errorf("Optimization status was not optimal! (Received %v)", optimstatus)
	}

	if sol.Values[0] != 0.0 {
		t.Errorf("Expected first variable to be be 0.0; received %v", sol.Values[0])
	}

	if optimstatus == gurobi.OPTIMAL {
		fmt.Printf("Optimal objective: %.4e\n", objval)
		fmt.Printf("  x_1=%.4f, x_2=%.4f\n", sol.Values[0], sol.Values[1])
	}
}
