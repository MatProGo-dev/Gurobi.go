package gurobi_test

import (
	"testing"

	"github.com/MatProGo-dev/Gurobi.go/gurobi"
)

/*
TestVar_GetInt1
Description:

	Tests how we can get the obj value of a given gurobi variable declared using the standard AddVar function.
*/
func TestVar_GetInt1(t *testing.T) {
	// Create environment.
	env, err := gurobi.NewEnv("setobj1.log")
	if err != nil {
		t.Errorf("There was an issue creating the new Env: %v", err)
	}
	defer env.Free()

	// Create an empty model.
	model, err := gurobi.NewModel("setobj1", env)
	if err != nil {
		t.Errorf("There was an issue creating the new model: %v", err)
	}
	defer model.Free()

	// Add varibles
	x, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		t.Errorf("There was an issue adding a variable to the old model: %v", err)
	}

	// Test current value of var (0)
	batchErrorCode, err := x.GetInt("BranchPriority")
	if err != nil {
		t.Errorf("There was an issue getting the obj value: %v", err.Error())
	}

	if batchErrorCode != 0 {
		t.Errorf("The initial BatchErrorCode value was %v; expected %v", batchErrorCode, 0.0)
	}
}

/*
TestVar_GetChar1
Description:

	Tests how we can get the variable type value of a given gurobi
	variable declared using the standard AddVar function.
*/
func TestVar_GetChar1(t *testing.T) {
	// Create environment.
	env, err := gurobi.NewEnv("setobj1.log")
	if err != nil {
		t.Errorf("There was an issue creating the new Env: %v", err)
	}
	defer env.Free()

	// Create an empty model.
	model, err := gurobi.NewModel("setobj1", env)
	if err != nil {
		t.Errorf("There was an issue creating the new model: %v", err)
	}
	defer model.Free()

	// Add varibles
	x, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		t.Errorf("There was an issue adding a variable to the old model: %v", err)
	}

	// Test current value of var (0)
	xVType, err := x.GetChar("VType")
	if err != nil {
		t.Errorf("There was an issue getting the obj value: %v", err.Error())
	}

	if xVType != gurobi.CONTINUOUS {
		t.Errorf("The initial VType value was %v; expected %v", xVType, gurobi.CONTINUOUS)
	}
}

/*
TestVar_GetDouble1
Description:

	Tests how we can get the obj value of a given gurobi variable
	declared using the standard AddVar function.
*/
func TestVar_GetDouble1(t *testing.T) {
	// Create environment.
	env, err := gurobi.NewEnv("setobj1.log")
	if err != nil {
		t.Errorf("There was an issue creating the new Env: %v", err)
	}
	defer env.Free()

	// Create an empty model.
	model, err := gurobi.NewModel("setobj1", env)
	if err != nil {
		t.Errorf("There was an issue creating the new model: %v", err)
	}
	defer model.Free()

	// Add varibles
	x, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		t.Errorf("There was an issue adding a variable to the old model: %v", err)
	}

	// Test current value of var (0)
	initialObjVal, err := x.GetDouble("obj")
	if err != nil {
		t.Errorf("There was an issue getting the obj value: %v", err.Error())
	}

	if initialObjVal != 0.0 {
		t.Errorf("The initial obj value was %v; expected %v", initialObjVal, 0.0)
	}
}

/*
TestVar_GetString1
Description:

	Tests how we can get the obj value of a given gurobi variable declared using the standard AddVar function.
*/
func TestVar_GetString1(t *testing.T) {
	// Create environment.
	env, err := gurobi.NewEnv("setobj1.log")
	if err != nil {
		t.Errorf("There was an issue creating the new Env: %v", err)
	}
	defer env.Free()

	// Create an empty model.
	model, err := gurobi.NewModel("setobj1", env)
	if err != nil {
		t.Errorf("There was an issue creating the new model: %v", err)
	}
	defer model.Free()

	// Add varibles
	x, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		t.Errorf("There was an issue adding a variable to the old model: %v", err)
	}

	// Test current value of var (0)
	xName, err := x.GetString("VarName")
	if err != nil {
		t.Errorf("There was an issue getting the obj value: %v", err.Error())
	}

	if xName != "x" {
		t.Errorf(
			"The initial Name value was %v; expected %v",
			xName,
			gurobi.CONTINUOUS,
		)
	}
}

/*
TestVar_SetObj1
Description:

	Tests how we can set the obj value of a given gurobi variable using the var object.
*/
func TestVar_SetObj1(t *testing.T) {
	// Create environment.
	env, err := gurobi.NewEnv("setobj2.log")
	if err != nil {
		t.Errorf("There was an issue creating the new Env: %v", err)
	}
	defer env.Free()

	// Create an empty model.
	model, err := gurobi.NewModel("setobj1", env)
	if err != nil {
		t.Errorf("There was an issue creating the new model: %v", err)
	}
	defer model.Free()

	// Add varibles
	x, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		t.Errorf("There was an issue adding a variable to the old model: %v", err)
	}

	// Set value of var

	err = x.SetObj(1.0)
	if err != nil {
		t.Errorf("There was an issue setting the obj value of the variable x: %v", err)
	}

	// Retrieve and compare the new obj value of x
	newObjVal, err := x.GetDouble("obj")
	if err != nil {
		t.Errorf("There was an issue getting the obj value: %v", err.Error())
	}

	if newObjVal != 1.0 {
		t.Errorf("The new obj value was %v; expected %v", newObjVal, 1.0)
	}
}
