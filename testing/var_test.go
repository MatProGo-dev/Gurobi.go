package testing

import (
	"testing"

	"github.com/MatProGo-dev/Gurobi.go/gurobi"
)

/*
TestVar_GetDouble1
Description:

	Tests how we can get the obj value of a given gurobi variable declared using the standard AddVar function.
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
