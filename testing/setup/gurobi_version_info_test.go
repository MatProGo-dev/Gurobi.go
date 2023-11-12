package test_setup

import (
	"github.com/MatProGo-dev/Gurobi.go/setup"
	"strconv"
	"strings"
	"testing"
)

/*
gurobi_version_info_test.go
Description:

	Testing some of the gurobi_version_info tests.
*/
func TestGurobiVersionInfo_StringToGurobiVersionInfoList1(t *testing.T) {
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
func TestGurobiVersionInfo_StringToGurobiVersionInfoList2(t *testing.T) {
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
TestGurobiVersionInfo_StringsToGurobiVersionInfoList3
Description:

	Tests that the StringToGurobiVersionInfo() function throws an error when a bad minor version
	was given.
*/
func TestGurobiVersionInfo_StringToGurobiVersionInfoList3(t *testing.T) {
	// Constants

	// Collect gurobi version info list

	// Search through Mac Library for all instances of Gurobi
	gurobiDirectory0 := "gurobi11b0"

	// Try to parse each of these directories
	_, err := setup.StringToGurobiVersionInfo(gurobiDirectory0)
	if err == nil {
		t.Errorf("no error was thrown, but we expected one to!")
	} else {
		_, err2 := strconv.Atoi("b")
		if !strings.Contains(
			err.Error(),
			err2.Error(),
		) {
			t.Errorf("unexpected error: %v", err)
		}
	}

}

/*
TestGurobiVersionInfo_StringsToGurobiVersionInfoList4
Description:

	Tests that the StringToGurobiVersionInfo() function throws an error when a bad major version
	was given.
*/
func TestGurobiVersionInfo_StringToGurobiVersionInfoList4(t *testing.T) {
	// Constants

	// Collect gurobi version info list

	// Search through Mac Library for all instances of Gurobi
	gurobiDirectory0 := "gurobic80"

	// Try to parse each of these directories
	_, err := setup.StringToGurobiVersionInfo(gurobiDirectory0)
	if err == nil {
		t.Errorf("no error was thrown, but we expected one to!")
	} else {
		_, err2 := strconv.Atoi("c")
		if !strings.Contains(
			err.Error(),
			err2.Error(),
		) {
			t.Errorf("unexpected error: %v", err)
		}
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

	if !strings.Contains(
		homeDirString,
		"macos_universal2",
	) {
		t.Errorf("expected gurobi home to be in \"macos_universal2\" directory, it was not.")
	}
}

/*
TestGurobiVersionInfo_CreateGurobiHomeDirectory2
Description:

	Testing whether or not the proper string is produced for version
	9.0.0 which was done before the macos_universal directory was used (proper directory should have mac64 in it).
*/
func TestGurobiVersionInfo_CreateGurobiHomeDirectory2(t *testing.T) {
	// Constants
	gv0 := setup.GurobiVersionInfo{MajorVersion: 9}

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
		"gurobi900",
	) {
		t.Errorf("expected gurobi home to be in \"gurobi900\" directory, it was not.")
	}

	if !strings.Contains(
		homeDirString,
		"mac64",
	) {
		t.Errorf("expected gurobi home to be in \"mac64\" directory, it was not.")
	}

}

/*
TestGurobiVersionInfo_StringsToGurobiVersionInfoList1
Description:

	Tests to make sure that the script that converts a correct gurobi directory name into a version info.
*/
func TestGurobiVersionInfo_StringsToGurobiVersionInfoList1(t *testing.T) {
	// Constants

	// Search through Mac Library for all instances of Gurobi
	gurobiDirectories := []string{"gurobi952", "gurobi1002"}

	// Convert multiple gurobi directories to GurobiVersionInfo objects.
	gviSlice0, err := setup.StringsToGurobiVersionInfoList(gurobiDirectories)
	if err != nil {
		t.Errorf("unexpected error with valid directory names: %v", err)
	}

	if gviSlice0[0].MajorVersion != 9 {
		t.Errorf("expected first gvi to have version 9; received %v", gviSlice0[0].MajorVersion)
	}

	if gviSlice0[1].MajorVersion != 10 {
		t.Errorf("expected second gvi to have version 10; received %v", gviSlice0[1].MajorVersion)
	}

}

/*
TestGurobiVersionInfo_StringsToGurobiVersionInfoList1
Description:

	Tests to make sure that the script that raises an error when one of the
	gurobi directories contains a bad tertiary version.
*/
func TestGurobiVersionInfo_StringsToGurobiVersionInfoList2(t *testing.T) {
	// Constants

	// Search through Mac Library for all instances of Gurobi
	gurobiDirectories := []string{"gurobi952", "gurobi100a"}

	// Convert multiple gurobi directories to GurobiVersionInfo objects.
	_, err := setup.StringsToGurobiVersionInfoList(gurobiDirectories)
	if err == nil {
		t.Errorf("no error was thrown, but should have been.")
	} else {
		_, err2 := strconv.Atoi("a")
		if !strings.Contains(
			err.Error(),
			err2.Error(),
		) {
			t.Errorf("unexpected error: %v", err)
		}
	}

}

/*
TestGurobiVersionInfo_FindHighestVersion1
Description:

	Verifies that the FindHighestVersionInfo function should throw an error if an empty slice is given.
*/
func TestGurobiVersionInfo_FindHighestVersion1(t *testing.T) {
	// Constants
	gviSlice0 := []setup.GurobiVersionInfo{}

	// Algorithm
	_, err := setup.FindHighestVersion(gviSlice0)
	if err == nil {
		t.Errorf("an error should have been thrown, but none were generated!")
	} else {
		if !strings.Contains(
			err.Error(),
			"No gurobi versions were provided to FindHighestVersion().",
		) {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

/*
TestGurobiVersionInfo_FindHighestVersion2
Description:

	Verifies that the FindHighestVersionInfo function should return the given GVI if a slice
	of length one is given.
*/
func TestGurobiVersionInfo_FindHighestVersion2(t *testing.T) {
	// Constants
	gviSlice0 := []setup.GurobiVersionInfo{
		{MajorVersion: 11, MinorVersion: 1, TertiaryVersion: 8},
	}

	// Algorithm
	gviOut, err := setup.FindHighestVersion(gviSlice0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if gviOut.MajorVersion != gviSlice0[0].MajorVersion {
		t.Errorf(
			"expected gviOut to be the same, but gviOut.MajorVersion %v doesn't match %v",
			gviOut.MajorVersion,
			gviSlice0[0].MajorVersion,
		)
	}

	if gviOut.MinorVersion != gviSlice0[0].MinorVersion {
		t.Errorf(
			"expected gviOut to be the same, but gviOut.MinorVersion %v doesn't match %v",
			gviOut.MinorVersion,
			gviSlice0[0].MinorVersion,
		)
	}

	if gviOut.TertiaryVersion != gviSlice0[0].TertiaryVersion {
		t.Errorf(
			"expected gviOut to be the same, but gviOut.TertiaryVersion %v doesn't match %v",
			gviOut.TertiaryVersion,
			gviSlice0[0].TertiaryVersion,
		)
	}
}

/*
TestGurobiVersionInfo_FindHighestVersion3
Description:

	Verifies that the FindHighestVersionInfo function should return the correct GVI if a slice
	of length two is given.
*/
func TestGurobiVersionInfo_FindHighestVersion3(t *testing.T) {
	// Constants
	gviSlice0 := []setup.GurobiVersionInfo{
		{MajorVersion: 11, MinorVersion: 1, TertiaryVersion: 8},
		{MajorVersion: 11, MinorVersion: 1, TertiaryVersion: 9},
	}

	// Algorithm
	gviOut, err := setup.FindHighestVersion(gviSlice0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if gviOut.MajorVersion != gviSlice0[1].MajorVersion {
		t.Errorf(
			"expected gviOut to be the same, but gviOut.MajorVersion %v doesn't match %v",
			gviOut.MajorVersion,
			gviSlice0[0].MajorVersion,
		)
	}

	if gviOut.MinorVersion != gviSlice0[1].MinorVersion {
		t.Errorf(
			"expected gviOut to be the same, but gviOut.MinorVersion %v doesn't match %v",
			gviOut.MinorVersion,
			gviSlice0[0].MinorVersion,
		)
	}

	if gviOut.TertiaryVersion != gviSlice0[1].TertiaryVersion {
		t.Errorf(
			"expected gviOut to be the same, but gviOut.TertiaryVersion %v doesn't match %v",
			gviOut.TertiaryVersion,
			gviSlice0[0].TertiaryVersion,
		)
	}
}

/*
TestGurobiVersionInfo_GreaterThan1
Description:

	Added a simple test for applying Greater Than to two GurobiVersionInfo objects
*/
func TestGurobiVersionInfo_GreaterThan1(t *testing.T) {
	// Constants
	gvi1 := setup.GurobiVersionInfo{9, 1, 0}
	gvi2 := setup.GurobiVersionInfo{11, 0, 0}

	// Test
	if !gvi2.GreaterThan(gvi1) {
		t.Errorf("algorithm thinks that gvi2 IS NOT greater than gvi1. Why?")
	}
}
