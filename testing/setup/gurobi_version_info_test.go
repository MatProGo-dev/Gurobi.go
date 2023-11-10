package test_setup

import (
	"github.com/MatProGo-dev/Gurobi.go/setup"
	"strings"
	"testing"
)

/*
gurobi_version_info_test.go
Description:

	Testing some of the gurobi_version_info tests.
*/
func TestGurobiVersionInfo_StringsToGurobiVersionInfoList1(t *testing.T) {
	// Constants

	// Collect gurobi version info list

	// Search through Mac Library for all instances of Gurobi
	gurobiDirectories := []string{"gurobi952", "gurobi1002"}

	// Try to parse each of these directories
	vi1, err := setup.StringToGurobiVersionInfo(gurobiDirectories[0])
	if err != nil {
		t.Errorf("There was an issue parsing the gurobi info from directory name: %v", err)
	}

	if vi1.MajorVersion != 9 {
		t.Errorf("Expected major version for %v to be 9; received %v", gurobiDirectories[0], vi1.MajorVersion)
	}
	if vi1.MinorVersion != 5 {
		t.Errorf("Expected minor version for %v to be 5; received %v", gurobiDirectories[0], vi1.MinorVersion)
	}
	if vi1.TertiaryVersion != 2 {
		t.Errorf("Expected tertiary version for %v to be 5; received %v", gurobiDirectories[0], vi1.TertiaryVersion)
	}

}

/*
TestGurobiVersionInfo_StringsToGurobiVersionInfoList2
Description:

	Testing 10.0.2
*/
func TestGurobiVersionInfo_StringsToGurobiVersionInfoList2(t *testing.T) {
	// Constants

	// Collect gurobi version info list

	// Search through Mac Library for all instances of Gurobi
	gurobiDirectories := []string{"gurobi952", "gurobi1002"}

	// Try to parse each of these directories
	vi1, err := setup.StringToGurobiVersionInfo(gurobiDirectories[1])
	if err != nil {
		t.Errorf("There was an issue parsing the gurobi info from directory name: %v", err)
	}

	if vi1.MajorVersion != 10 {
		t.Errorf("Expected major version for %v to be 10; received %v", gurobiDirectories[0], vi1.MajorVersion)
	}
	if vi1.MinorVersion != 0 {
		t.Errorf("Expected minor version for %v to be 0; received %v", gurobiDirectories[0], vi1.MinorVersion)
	}
	if vi1.TertiaryVersion != 2 {
		t.Errorf("Expected tertiary version for %v to be 2; received %v", gurobiDirectories[0], vi1.TertiaryVersion)
	}

}

/*
TestGurobiVersionInfo_CreateGurobiHomeDirectory1
Description:

	Testing whether or not the proper string is produced for version
	11.0.0 of Gurobi (newest as of Nov. 10, 2023).
*/
func TestGurobiVersionInfo_CreateGurobiHomeDirectory1(t *testing.T) {
	// Constants
	gv0 := setup.GurobiVersionInfo{MajorVersion: 11}

	// Get home directory string for Mac
	homeDirString, err := setup.CreateGurobiHomeDirectory(gv0)
	if err != nil {
		t.Errorf("unexpected error creating home directory: %v", err)
	}

	if !strings.Contains(
		homeDirString,
		"Library",
	) {
		t.Errorf("expected gurobi home to be in Library directory, it was not.")
	}

	if !strings.Contains(
		homeDirString,
		"gurobi1100",
	) {
		t.Errorf("expected gurobi home to be in \"gurobi11000\" directory, it was not.")
	}
}
