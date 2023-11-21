package gurobi_test

import (
	"fmt"
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
TestModel_Check1
Description:

	Tests if Check() will fail if the env is not
	yet defined in the Model object.
*/
func TestModel_Check1(t *testing.T) {
	// Constants
	model0 := gurobi.Model{
		Variables: []gurobi.Var{},
	}

	// Tests
	err := model0.Check()
	if err == nil {
		t.Errorf("expected an error, but none were thrown!")
	} else {
		if err.Error() != model0.Env.MakeUninitializedError().Error() {
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

/*
TestModel_AddVar3
Description:

	Tests that AddVar() successfully handles a call which
	adds multiple variables AND contains pre-defined constraints.
*/
func TestModel_AddVar3(t *testing.T) {
	// Constants
	testIndex := 8
	testName := fmt.Sprintf("testmodel-addvar%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Define empty constraints
	_, err = model0.AddConstr(
		[]*gurobi.Var{},
		[]float64{},
		gurobi.SenseLessThan,
		2.0,
		"max-2",
	)

	emptyConstraint1, err := model0.AddConstr(
		[]*gurobi.Var{},
		[]float64{},
		gurobi.SenseLessThan,
		10.0,
		"max-10",
	)

	// Test
	columns := make([]float64, 1)
	columns[0] = 1.0
	v1, err := model0.AddVar(
		gurobi.CONTINUOUS,
		-1.0,
		-1e3,
		1e3,
		"ellie-goulding",
		[]*gurobi.Constr{emptyConstraint1}, // The constraint to modify
		columns,                            // The new linear coefficient for v1 in constraint emptyConstraint1
	)
	if err != nil {
		t.Errorf("unexpected error while adding variable: %v!", err)
	}

	// Solve the problem to make sure that the constraint was properly modified.

	// Integrate new variables.
	if err := model0.Update(); err != nil {
		panic(err.Error())
	}

	// Optimize model
	if err := model0.Optimize(); err != nil {
		panic(err.Error())
	}

	optimstatus, err := model0.GetIntAttr(gurobi.INT_ATTR_STATUS)
	if err != nil {
		t.Errorf("there was an issue getting the optimization's status: %v", err.Error())
	}

	// Optimization should have been successful was it?
	if optimstatus != gurobi.OPTIMAL {
		t.Errorf("optimization status was not optimal! (%v)", optimstatus)
	}

	// WAs solution 10.0?
	sol, err := model0.GetDoubleAttrVars(
		gurobi.DBL_ATTR_X,
		[]*gurobi.Var{v1},
	)
	if err != nil {
		t.Errorf("there was an issue with getting the solution: %v", err.Error())
	}

	if sol[0] != 10.0 {
		t.Errorf("expected solution to e")
	}

}

/*
TestModel_AddVars1
Description:

	Tests that AddVars() throws an error when an unspecified model has been defined.
*/
func TestModel_AddVars1(t *testing.T) {
	// Constants
	var model0 *gurobi.Model

	// Test
	_, err := model0.AddVars(
		[]int8{gurobi.CONTINUOUS},
		[]float64{0.0},
		[]float64{-1.0},
		[]float64{1.0},
		[]string{"ellie-goulding"},
		[][]*gurobi.Constr{},
		[][]float64{},
	)
	if err == nil {
		t.Errorf("No error was thrown, but there should have been!")
	} else {
		if err.Error() != model0.MakeUninitializedError().Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestModel_AddVars2
Description:

	Tests that AddVars() throws an error when
	there is a different number of vtypes are specified
	than objs
*/
func TestModel_AddVars2(t *testing.T) {
	// Constants
	env0, err := gurobi.NewEnv("testmodel-addvars2.log")
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel("testmodel-addvars2", env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	_, err = model0.AddVars(
		[]int8{gurobi.CONTINUOUS, gurobi.INTEGER},
		[]float64{0.0},
		[]float64{-1.0},
		[]float64{1.0},
		[]string{"ellie-goulding"},
		[][]*gurobi.Constr{},
		[][]float64{},
	)
	if err == nil {
		t.Errorf("No error was thrown, but there should have been!")
	} else {
		if err.Error() != (gurobi.MismatchedLengthError{
			Length1: 2,
			Length2: 1,
			Name1:   "vtypes",
			Name2:   "objs",
		}).Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestModel_AddVars3
Description:

	Tests that AddVars() throws an error when
	there is a different number of objs are specified
	than lbs
*/
func TestModel_AddVars3(t *testing.T) {
	// Constants
	testIndex := 3
	testName := fmt.Sprintf("testmodel-addvars%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	_, err = model0.AddVars(
		[]int8{gurobi.CONTINUOUS},
		[]float64{0.0},
		[]float64{-1.0, 1.0},
		[]float64{1.0},
		[]string{"ellie-goulding"},
		[][]*gurobi.Constr{},
		[][]float64{},
	)
	if err == nil {
		t.Errorf("No error was thrown, but there should have been!")
	} else {
		if err.Error() != (gurobi.MismatchedLengthError{
			Length1: 1,
			Length2: 2,
			Name1:   "objs",
			Name2:   "lbs",
		}).Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestModel_AddVars4
Description:

	Tests that AddVars() throws an error when
	there is a different number of lbs are specified
	than ubs
*/
func TestModel_AddVars4(t *testing.T) {
	// Constants
	testIndex := 4
	testName := fmt.Sprintf("testmodel-addvars%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	_, err = model0.AddVars(
		[]int8{gurobi.CONTINUOUS},
		[]float64{0.0},
		[]float64{-1.0},
		[]float64{1.0, gurobi.INFINITY},
		[]string{"ellie-goulding"},
		[][]*gurobi.Constr{},
		[][]float64{},
	)
	if err == nil {
		t.Errorf("No error was thrown, but there should have been!")
	} else {
		if err.Error() != (gurobi.MismatchedLengthError{
			Length1: 1,
			Length2: 2,
			Name1:   "lbs",
			Name2:   "ubs",
		}).Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestModel_AddVars5
Description:

	Tests that AddVars() throws an error when
	there is a different number of ubs are specified
	than names
*/
func TestModel_AddVars5(t *testing.T) {
	// Constants
	testIndex := 5
	testName := fmt.Sprintf("testmodel-addvars%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	_, err = model0.AddVars(
		[]int8{gurobi.CONTINUOUS},
		[]float64{0.0},
		[]float64{-1.0},
		[]float64{1.0},
		[]string{"ellie-goulding", "chris-brown"},
		[][]*gurobi.Constr{},
		[][]float64{},
	)
	if err == nil {
		t.Errorf("No error was thrown, but there should have been!")
	} else {
		if err.Error() != (gurobi.MismatchedLengthError{
			Length1: 1,
			Length2: 2,
			Name1:   "ubs",
			Name2:   "names",
		}).Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestModel_AddVars6
Description:

	Tests that AddVars() throws an error when
	there is a different number of names are specified
	than constrs
*/
func TestModel_AddVars6(t *testing.T) {
	// Constants
	testIndex := 6
	testName := fmt.Sprintf("testmodel-addvars%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	_, err = model0.AddVars(
		[]int8{gurobi.CONTINUOUS},
		[]float64{0.0},
		[]float64{-1.0},
		[]float64{1.0},
		[]string{"ellie-goulding"},
		[][]*gurobi.Constr{
			[]*gurobi.Constr{},
			[]*gurobi.Constr{},
		},
		[][]float64{},
	)
	if err == nil {
		t.Errorf("No error was thrown, but there should have been!")
	} else {
		if err.Error() != (gurobi.MismatchedLengthError{
			Length1: 1,
			Length2: 2,
			Name1:   "names",
			Name2:   "constrs",
		}).Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestModel_AddVars7
Description:

	Tests that AddVars() throws an error when
	there is a different number of columns are specified
	than columns
*/
func TestModel_AddVars7(t *testing.T) {
	// Constants
	testIndex := 7
	testName := fmt.Sprintf("testmodel-addvars%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	_, err = model0.AddVars(
		[]int8{gurobi.CONTINUOUS},
		[]float64{0.0},
		[]float64{-1.0},
		[]float64{1.0},
		[]string{"ellie-goulding"},
		[][]*gurobi.Constr{
			[]*gurobi.Constr{},
		},
		[][]float64{
			[]float64{1.0},
			[]float64{2.0},
		},
	)
	if err == nil {
		t.Errorf("No error was thrown, but there should have been!")
	} else {
		if err.Error() != (gurobi.MismatchedLengthError{
			Length1: 1,
			Length2: 2,
			Name1:   "constrs",
			Name2:   "columns",
		}).Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestModel_AddVars8
Description:

	Tests that AddVars() successfully handles a call which only adds a single variable.
*/
func TestModel_AddVars8(t *testing.T) {
	// Constants
	testIndex := 8
	testName := fmt.Sprintf("testmodel-addvars%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		t.Errorf("unexpected error creating new environment: %v", err)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		t.Errorf("unexpected error creating new model: %v", err)
	}

	// Test
	vSlice0, err := model0.AddVars(
		[]int8{gurobi.CONTINUOUS},
		[]float64{0.0},
		[]float64{-1.0},
		[]float64{1.0},
		[]string{"ellie-goulding"},
		[][]*gurobi.Constr{},
		[][]float64{},
	)
	if err != nil {
		t.Errorf("unexpected error: %v!", err)
	}

	if len(vSlice0) != 1 {
		t.Errorf("expected output slice to be of length 1; received length %v slice", len(vSlice0))
	}
}
