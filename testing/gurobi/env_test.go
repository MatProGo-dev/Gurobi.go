package gurobi_test

import (
	"github.com/MatProGo-dev/Gurobi.go/gurobi"
	"testing"
)

/*
TestEnv_NewEnv1
Description:

	Verifies that NewEnv() correctly creates a new environment with log file given.
*/
func TestEnv_NewEnv1(t *testing.T) {
	// Constants
	logfilename1 := "thomTide.log"

	// Algorithm
	env, err := gurobi.NewEnv(logfilename1)
	if err != nil {
		t.Errorf("There was an issue creating the new Env variable: %v", err)
	}
	defer env.Free()

}

/*
TestEnv_SetTimeLimit1
Description:

	Verifies that NewEnv() correctly creates a new environment with log file given.
*/
func TestEnv_SetTimeLimit1(t *testing.T) {
	// Constants
	logfilename1 := "thomTide.log"
	var newTimeLimit float64 = 132

	// Algorithm
	env, err := gurobi.NewEnv(logfilename1)
	if err != nil {
		t.Errorf("There was an issue creating the new Env variable: %v", err)
	}
	defer env.Free()

	err = env.SetTimeLimit(newTimeLimit)
	if err != nil {
		t.Errorf("There was an error setting the time limit of the environment! %v", err)
	}

	detectedTimeLimit, err := env.GetTimeLimit()
	if err != nil {
		t.Errorf("There was an error getting the time limit of the environment! %v", err)
	}

	if detectedTimeLimit != newTimeLimit {
		t.Errorf("The detected time limit (%v) was not equal to the expected time limit (%v s).", detectedTimeLimit, newTimeLimit)
	}

}

/*
TestEnv_SetTimeLimit2
Description:

	Tests that the function throws an error, when it receives a
	uninitialized.
*/
func TestEnv_SetTimeLimit2(t *testing.T) {
	// Constants
	var env0 *gurobi.Env

	// Algorithm
	err := env0.SetTimeLimit(1e2)
	if err == nil {
		t.Errorf("expected an error to be thrown, but received none!")
	} else {
		if err.Error() != env0.MakeUninitializedError().Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestEnv_SetDBLParam1
Description:

	Verifies that we can set the value of 'TimeLimit' in current model.
*/
func TestEnv_SetDBLParam1(t *testing.T) {
	// Constants
	logfilename1 := "thomTide.log"
	var newTimeLimit float64 = 132

	// Algorithm
	env, err := gurobi.NewEnv(logfilename1)
	if err != nil {
		t.Errorf("There was an issue creating the new Env variable: %v", err)
	}
	defer env.Free()

	err = env.SetDBLParam("TimeLimit", newTimeLimit)
	if err != nil {
		t.Errorf("There was an error setting the time limit of the environment! %v", err)
	}

	detectedTimeLimit, err := env.GetDBLParam("TimeLimit")
	if err != nil {
		t.Errorf("There was an error getting the time limit of the environment! %v", err)
	}

	if detectedTimeLimit != newTimeLimit {
		t.Errorf("The detected time limit (%v) was not equal to the expected time limit (%v s).", detectedTimeLimit, newTimeLimit)
	}

}

/*
TestEnv_SetDBLParam2
Description:

	Verifies that we can set the value of 'BestObjStop' in current model.
*/
func TestEnv_SetDBLParam2(t *testing.T) {
	// Constants
	logfilename1 := "thomTide.log"
	var newVal float64 = 132
	var paramToModify string = "BestObjStop"

	// Algorithm
	env, err := gurobi.NewEnv(logfilename1)
	if err != nil {
		t.Errorf("There was an issue creating the new Env variable: %v", err)
	}
	defer env.Free()

	err = env.SetDBLParam(paramToModify, newVal)
	if err != nil {
		t.Errorf("There was an error setting the time limit of the environment! %v", err)
	}

	detectedVal, err := env.GetDBLParam(paramToModify)
	if err != nil {
		t.Errorf("There was an error getting the time limit of the environment! %v", err)
	}

	if detectedVal != newVal {
		t.Errorf("The detected %v (%v) was not equal to the expected %v (%v).", paramToModify, detectedVal, paramToModify, newVal)
	}

}

/*
TestEnv_Check1
Description:

	Verifies that the method returns an error when the env is uninitialized.
*/
func TestEnv_Check1(t *testing.T) {
	// Constants
	var env0 *gurobi.Env

	// Check
	err := env0.Check()
	if err == nil {
		t.Errorf("an error should have been thrown, but none were detected!")
	} else {
		if err.Error() != env0.MakeUninitializedError().Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}

}

/*
TestEnv_Check2
Description:

	Verifies that the method returns no error when the env is
	properly initialized.
*/
func TestEnv_Check2(t *testing.T) {
	// Constants
	env0, err := gurobi.NewEnv("testenv-check2.log")
	if err != nil {
		t.Errorf("error creating test environment: %v", err)
	}

	// Check
	err = env0.Check()
	if err != nil {
		t.Errorf("unexpected error during Check(): %v", err)
	}

}

/*
TestEnv_GetTimeLimit1
Description:

	Tests that the function throws an error, when it receives a
	uninitialized.
*/
func TestEnv_GetTimeLimit1(t *testing.T) {
	// Constants
	var env0 *gurobi.Env

	// Algorithm
	_, err := env0.GetTimeLimit()
	if err == nil {
		t.Errorf("expected an error to be thrown, but received none!")
	} else {
		if err.Error() != env0.MakeUninitializedError().Error() {
			t.Errorf("unexpected error: %v", err)
		}
	}
}
