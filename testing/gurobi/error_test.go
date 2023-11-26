package gurobi_test

import (
	"github.com/MatProGo-dev/Gurobi.go/gurobi"
	"testing"
)

/*
error_test.go
Description:
	This file tests the Error class in the gurobi package.
*/

/*
TestError_Error1
Description:

	Tests the Error() method.
*/
func TestError_Error1(t *testing.T) {
	// Constants
	err1 := gurobi.Error{
		ErrorCode: int32(2),
		Message:   "heart",
	}

	// Test
	if err1.Error() != "heart" {
		t.Errorf("unexpected error result: %v", err1)
	}

}

/*
TestError_Error2
Description:

	Tests the Error() method with empty message
*/
func TestError_Error2(t *testing.T) {
	// Constants
	err1 := gurobi.Error{
		ErrorCode: int32(2),
	}

	// Test
	if err1.Error() != "" {
		t.Errorf("unexpected error result: %v", err1)
	}

}
