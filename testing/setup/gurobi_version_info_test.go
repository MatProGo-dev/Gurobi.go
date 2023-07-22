package test_setup

import (
	"github.com/MatProGo-dev/Gurobi.go/setup"
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
