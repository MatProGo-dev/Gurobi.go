package gurobi_test

import (
	"github.com/MatProGo-dev/Gurobi.go/gurobi"
	"testing"
)

/*
TestModel_NewModel1
Description:

	Verifies that the NewModel() function returns an error when
	called with an uninitialized environment.
*/
func TestModel_NewModel1(t *testing.T) {
	// Constants
	var env0 *gurobi.Env

	// Algorithm
	_, err := gurobi.NewModel("testmodel-newmodel1", env0)
	if err == nil {
		t.Errorf("expected error to be thrown, but none were detected!")
	} else {
		if err.Error() != env0.MakeUninitializedError().Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}

}

/*
TestModel_NewModel1
Description:

	Verifies that the NewModel() function does nothing when a nil model
	is provided.
*/
func TestModel_Free1(t *testing.T) {
	// Constants
	var model0 *gurobi.Model

	// Algorithm
	model0.Free()
}

/*
TestModel_AddVar1
Description:

	Verifies that the system properly throws an error when a model is used in
	AddVar that is not initialized.
*/
func TestModel_AddVar1(t *testing.T) {
	// Create nil model pointer
	var model0 *gurobi.Model

	// Test
	_, err := model0.AddVar(
		gurobi.CONTINUOUS, 0.0, 0.0, 1.0, "test",
		[]*gurobi.Constr{},
		[]float64{},
	)
	if err == nil {
		t.Errorf("expected an error, but received none!")
	} else {
		if err.Error() != model0.MakeUninitializedError().Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}

}

/*
TestModel_AddVar2
Description:

	Verifies that the system does not throw an error when calling AddVar
	when the system has been properly initialized.
*/
func TestModel_AddVar2(t *testing.T) {
	// Create nil model pointer
	env0, err := gurobi.NewEnv("testmodel-addvar2.log")
	if err != nil {
		t.Errorf("unexpected error creating new env: %v", err)
	}

	model0, err := gurobi.NewModel("testmodel-addvar2", env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	v1, err := model0.AddVar(
		gurobi.CONTINUOUS, 0.0, 0.0, 1.0, "test",
		[]*gurobi.Constr{},
		[]float64{},
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if lb0, _ := v1.GetDouble(("LB")); lb0 != 0.0 {
		t.Errorf("Retrieved LB (%v) doesn't match expected (%v)", lb0, 0.0)
	}

}
