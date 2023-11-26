package test_setup

import (
	"github.com/MatProGo-dev/Gurobi.go/setup"
	"os"
	"strings"
	"testing"
)

/*
setup_flags_test.go
Description
	Tests the methods designed for the CreateLDFlagsDirective object.
*/

/*
TestSetupFlags_GetDefaultSetupFlags1
Description:
	Simply runs the default setup flags test on the current installation for gurobi available on the testing
	computer.
*/

func TestSetupFlags_GetDefaultSetupFlags1(t *testing.T) {
	// Constants

	// Test
	sf, err := setup.GetDefaultSetupFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gv, err := sf.ToGurobiVersionInfo()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if gv.MajorVersion != 10 {
		t.Errorf("unexpected gurobi version found: %v", gv.MajorVersion)
	}
}

/*
TestSetupFlags_ParseMakeLibArguments1
Description:

	Tests that the parsing method for parsing command line arguments into the SetupFlags script
	correctly parses things/leaves others alone.
*/
func TestSetupFlags_ParseMakeLibArguments1(t *testing.T) {
	// Constants
	sf0, _ := setup.GetDefaultSetupFlags()

	// Test
	sfOut, err := setup.ParseMakeLibArguments(sf0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !strings.Contains(sfOut.PackageName, sf0.PackageName) {
		t.Errorf(
			"Package name should be preserved after ParseMakeLibArguments, but PackageName was changed from %v to %v",
			sf0.PackageName,
			sfOut.PackageName,
		)
	}
}

/*
TestSetupFlags_WriteHeaderFile1
Description:

	Tests that the writing of the header file works without any errors when given a good enough file name.
*/
func TestSetupFlags_WriteHeaderFile1(t *testing.T) {
	// Constants
	sf0, _ := setup.GetDefaultSetupFlags()
	sf0.HeaderFilename = "./testHeader.txt"

	// Test
	err := setup.WriteHeaderFile(sf0)
	if err != nil {
		t.Errorf("unexpected issue when writing header file: %v", err)
	}
	defer os.Remove(sf0.HeaderFilename)
}

/*
TestSetupFlags_CreatePackageLine1
Description:

	Tests that the package line is properly created in a go file.
*/
func TestSetupFlags_CreatePackageLine1(t *testing.T) {
	// Constants
	sf0, _ := setup.GetDefaultSetupFlags()

	// Algorithm
	line0, err := setup.CreatePackageLine(sf0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !strings.Contains(
		line0,
		"package",
	) {
		t.Errorf("line is expected to contain the word \"package\", but it did not.")
	}
}

/*
TestSetupFlags_CreateLDFlagsDirective1
Description:

	Tests the LD_Flags string is properly created by the CreateLDFlagsDirective().
*/
func TestSetupFlags_CreateLDFlagsDirective1(t *testing.T) {
	// Constants
	sf0, _ := setup.GetDefaultSetupFlags()

	// Algorithm
	directive0, err := setup.CreateLDFlagsDirective(sf0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !strings.Contains(
		directive0,
		"LDFLAGS",
	) {
		t.Errorf("line is expected to contain the word \"LDFLAGS\", but it did not.")
	}

	if !strings.Contains(
		directive0,
		"gurobi_c++",
	) {
		t.Errorf("line is expected to contain the word \"gurobi_c++\", but it did not.")
	}
}

/*
TestSetupFlags_CreateLDFlagsDirective1
Description:

	Tests the LD_Flags string is properly created by the CreateLDFlagsDirective().
*/
func TestSetupFlags_CreateLDFlagsDirective2(t *testing.T) {
	// Constants
	sf0, _ := setup.GetDefaultSetupFlags()
	sf0.GurobiHome = "tempt"

	// Algorithm
	_, err := setup.CreateLDFlagsDirective(sf0)
	if err == nil {
		t.Errorf("no error was thrown, but they should have been!")
	} else {
		if !strings.Contains(
			err.Error(),
			"\"gurobi\" is not found in the setup flag's GurobiHome!",
		) {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestSetupFlags_WriteLibGo1
Description:

	Testing that the WriteLibGo1 file
*/
func TestSetupFlags_WriteLibGo1(t *testing.T) {
	// Constants
	sf0, _ := setup.GetDefaultSetupFlags()
	sf0.GoFilename = "./test_lib.go"
	sf0.PackageName = "test_setup"

	// Algorithm
	err := setup.WriteLibGo(sf0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(sf0.GoFilename)
}
