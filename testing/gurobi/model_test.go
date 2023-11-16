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
