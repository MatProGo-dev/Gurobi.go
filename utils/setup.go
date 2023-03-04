/*
setup.go
Description:
	An implementation of the file local_lib_info.go written entirely in go.
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/*
Constants
*/
const GoLibraryFilename string = "gurobi/cgoHelper.go"
const CppHeaderFilename string = "gurobi/gurobi_passthrough.h"

/*
Type Definitions
*/

type SetupFlags struct {
	GurobiHome     string // Directory where Gurobi is installed
	GoFilename     string // Name of the Go File to
	HeaderFilename string // Name of the Headerfile to Create
	PackageName    string // Name of the package
}

type GurobiVersionInfo struct {
	MajorVersion    int
	MinorVersion    int
	TertiaryVersion int
}

/*
CreateGurobiHomeDirectory
Description:

	Creates the Gurobi Home directory where your gurobi program and helper files should be installed when you run
	the Gurobi Installer.
*/
func CreateGurobiHomeDirectory(versionInfo GurobiVersionInfo) (string, error) {
	// Create Base Home Directory
	gurobiHome := fmt.Sprintf(
		"/Library/gurobi%v%v%v",
		versionInfo.MajorVersion,
		versionInfo.MinorVersion,
		versionInfo.TertiaryVersion,
	)

	// Create
	switch runtime.GOOS {
	case "darwin":
		// Decide on if we are on intel processors or apple silicon.
		switch runtime.GOARCH {
		case "amd64":
			gurobiHome := fmt.Sprintf("%v/mac64", gurobiHome)
			return gurobiHome, nil
		case "arm64":
			gurobiHome := fmt.Sprintf("%v/macos_universal2", gurobiHome)
			return gurobiHome, nil
		default:
			return "", fmt.Errorf("There was an error finding the architecture on your mac! What is \"%v\"? Send an issue to the gurobi.go repository.", runtime.GOARCH)
		}
	default:
		return "", fmt.Errorf("The operating system that you are using is not recognized: \"%v\".", runtime.GOOS)
	}

}

func GetDefaultSetupFlags() (SetupFlags, error) {
	// Create Default Struct
	defaultGurobiVersion := GurobiVersionInfo{9, 0, 3}
	defaultHome, err := CreateGurobiHomeDirectory(defaultGurobiVersion)
	if err != nil {
		return SetupFlags{}, err
	}

	mlf := SetupFlags{
		GurobiHome:     defaultHome,
		GoFilename:     GoLibraryFilename,
		HeaderFilename: CppHeaderFilename,
		PackageName:    "gurobi",
	}

	// Search through Mac Library for all instances of Gurobi
	libraryContents, err := os.ReadDir("/Library")
	if err != nil {
		return mlf, err
	}
	gurobiDirectories := []string{}
	for _, content := range libraryContents {
		if content.IsDir() && strings.Contains(content.Name(), "gurobi") {
			fmt.Println(content.Name())
			gurobiDirectories = append(gurobiDirectories, content.Name())
		}
	}

	// Convert Directories into Gurobi Version Info
	gurobiVersionList, err := StringsToGurobiVersionInfoList(gurobiDirectories)
	if err != nil {
		return mlf, err
	}

	highestVersion, err := FindHighestVersion(gurobiVersionList)
	if err != nil {
		return mlf, err
	}

	// Write the highest version's directory into the GurobiHome variable
	mlf.GurobiHome, err = CreateGurobiHomeDirectory(highestVersion)
	if err != nil {
		return mlf, err
	}

	return mlf, nil

}

/*
StringToGurobiVersionInfo
Assumptions:

	Assumes that a valid gurobi name is given.
*/
func StringToGurobiVersionInfo(gurobiDirectoryName string) (GurobiVersionInfo, error) {
	//Locate major and minor version indices in gurobi directory name
	majorVersionAsString := string(gurobiDirectoryName[len("gurobi")])
	minorVersionAsString := string(gurobiDirectoryName[len("gurobi")+1])
	tertiaryVersionAsString := string(gurobiDirectoryName[len("gurobi")+2])

	// Convert using strconv to integers
	majorVersion, err := strconv.Atoi(majorVersionAsString)
	if err != nil {
		return GurobiVersionInfo{}, err
	}

	minorVersion, err := strconv.Atoi(minorVersionAsString)
	if err != nil {
		return GurobiVersionInfo{}, err
	}

	tertiaryVersion, err := strconv.Atoi(tertiaryVersionAsString)
	if err != nil {
		return GurobiVersionInfo{}, err
	}

	return GurobiVersionInfo{
		MajorVersion:    majorVersion,
		MinorVersion:    minorVersion,
		TertiaryVersion: tertiaryVersion,
	}, nil

}

/*
StringsToGurobiVersionInfoList
Description:

	Receives a set of strings which should be in the format of valid gurobi installation directories
	and returns a list of GurobiVersionInfo objects.

Assumptions:

	Assumes that a valid gurobi name is given.
*/
func StringsToGurobiVersionInfoList(gurobiDirectoryNames []string) ([]GurobiVersionInfo, error) {

	// Convert Directories into Gurobi Version Info
	gurobiVersionList := []GurobiVersionInfo{}
	for _, directory := range gurobiDirectoryNames {
		tempGVI, err := StringToGurobiVersionInfo(directory)
		if err != nil {
			return gurobiVersionList, err
		}
		gurobiVersionList = append(gurobiVersionList, tempGVI)
	}
	// fmt.Println(gurobiVersionList)

	return gurobiVersionList, nil

}

/*
// Iterate through all versions in gurobiVersionList and find the one with the largest major or minor version.
*/
func FindHighestVersion(gurobiVersionList []GurobiVersionInfo) (GurobiVersionInfo, error) {

	// Input Checking
	if len(gurobiVersionList) == 0 {
		return GurobiVersionInfo{}, errors.New("No gurobi versions were provided to FindHighestVersion().")
	}

	// Perform search
	highestVersion := gurobiVersionList[0]
	if len(gurobiVersionList) == 1 {
		return highestVersion, nil
	}

	for _, gvi := range gurobiVersionList {
		// Compare Major version numbers
		if gvi.MajorVersion > highestVersion.MajorVersion {
			highestVersion = gvi
			continue
		}

		// Compare minor version numbers
		if gvi.MinorVersion > highestVersion.MinorVersion {
			highestVersion = gvi
			continue
		}

		// Compare tertiary version numbers
		if gvi.TertiaryVersion > highestVersion.TertiaryVersion {
			highestVersion = gvi
			continue
		}
	}

	return highestVersion, nil

}

func ParseMakeLibArguments(sfIn SetupFlags) (SetupFlags, error) {
	// Iterate through any arguments with mlfIn as the default
	sfOut := sfIn

	// Input Processing
	argIndex := 1 // Skip entry 0
	for argIndex < len(os.Args) {
		// Share parsing data
		fmt.Println("- Parsed input: %v", os.Args[argIndex])

		// Parse Inputs
		switch {
		case os.Args[argIndex] == "--gurobi-home":
			sfOut.GurobiHome = os.Args[argIndex+1]
			argIndex += 2
		case os.Args[argIndex] == "--go-fname":
			sfOut.GoFilename = os.Args[argIndex+1]
			argIndex += 2
		case os.Args[argIndex] == "--pkg":
			sfOut.PackageName = os.Args[argIndex+1]
			argIndex += 2
		default:
			fmt.Printf("Unrecognized input: %v", os.Args[argIndex])
			argIndex++
		}

	}

	return sfOut, nil
}

/*
CreateCXXFlagsDirective
Description:

	Creates the CXX Flags directive in the  file that we will use in lib.go.
*/
func CreateCXXFlagsDirective(sf SetupFlags) (string, error) {
	// Create Statement

	gurobiCXXFlagsString := fmt.Sprintf("// #cgo CXXFLAGS: --std=c++11 -I%v/include \n", sf.GurobiHome)
	//lpSolveCXXFlagsString := "// #cgo CXXFLAGS: -I/usr/local/opt/lp_solve/include\n" // Works as long as lp_solve was installed with Homebrew

	return gurobiCXXFlagsString, nil
}

/*
CreatePackageLine
Description:

	Creates the "package" directive in the  file that we will use in lib.go.
*/
func CreatePackageLine(sf SetupFlags) (string, error) {

	return fmt.Sprintf("package %v\n\n", sf.PackageName), nil
}

/*
CreateLDFlagsDirective
Description:

	Creates the LD_FLAGS directive in the file that we will use in lib.go.
*/
func CreateLDFlagsDirective(sf SetupFlags) (string, error) {
	// Constants
	AsGVI, err := sf.ToGurobiVersionInfo()
	if err != nil {
		return "", err
	}

	// Locate the desired files for mac in the directory.
	// libContent, err := os.ReadDir(mlfIn.GurobiHome)
	// if err != nil {
	// 	return "", err
	// }

	ldFlagsDirective := fmt.Sprintf("// #cgo LDFLAGS: -L%v/lib", sf.GurobiHome)

	targetedFilenames := []string{"gurobi_c++", fmt.Sprintf("gurobi%v%v", AsGVI.MajorVersion, AsGVI.MinorVersion)}

	for _, target := range targetedFilenames {
		ldFlagsDirective = fmt.Sprintf("%v -l%v", ldFlagsDirective, target)
	}
	ldFlagsDirective = fmt.Sprintf("%v \n", ldFlagsDirective)

	return ldFlagsDirective, nil
}

func (mlf *SetupFlags) ToGurobiVersionInfo() (GurobiVersionInfo, error) {
	// Split the GurobiHome variable by the name gurobi
	GurobiWordIndexStart := strings.Index(mlf.GurobiHome, "gurobi")
	GurobiDirNameIndexEnd := len(mlf.GurobiHome) - len("/mac64") - 1

	return StringToGurobiVersionInfo(string(mlf.GurobiHome[GurobiWordIndexStart : GurobiDirNameIndexEnd+1]))

}

func GetAHeaderFilenameFrom(dirName string) (string, error) {
	// Constants

	// Algorithm

	// Search through dirName directory for all instances of .a files
	libraryContents, err := os.ReadDir(dirName)
	if err != nil {
		return "", err
	}
	headerNames := []string{}
	for _, content := range libraryContents {
		if content.Type().IsRegular() && strings.Contains(content.Name(), ".a") {
			fmt.Println(content.Name())
			headerNames = append(headerNames, content.Name())
		}
	}

	return headerNames[0], nil

}

/*
WriteLibGo
Description:

	Creates the library file which imports the proper libraries for cgo.
	By default this is named according to GoLibraryFilename.
*/
func WriteLibGo(sf SetupFlags) error {
	// Constants

	// Algorithm

	// First Create all Strings that we would like to write to lib.go
	// 1. Create package definition
	packageDirective, err := CreatePackageLine(sf)
	if err != nil {
		return err
	}

	// 2. Create CXX_FLAGS argument
	cxxDirective, err := CreateCXXFlagsDirective(sf)
	if err != nil {
		return err
	}

	// 3. Create LDFLAGS Argument
	ldflagsDirective, err := CreateLDFlagsDirective(sf)
	if err != nil {
		return err
	}

	// Now Write to File
	f, err := os.Create(sf.GoFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write all directives to file
	_, err = f.WriteString(fmt.Sprintf("%v%v%v import \"C\"\n", packageDirective, cxxDirective, ldflagsDirective))
	if err != nil {
		return err
	}

	return nil

}

/*
WriteHeaderFile
Description:

	This script writes the C++ header file which goes in gurobi.go but references
	the true gurobi_c.h file.
*/
func WriteHeaderFile(sf SetupFlags) error {
	// Constants

	// Algorithm

	// Now Write to File
	f, err := os.Create(sf.HeaderFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write a small comment + import
	simpleComment := fmt.Sprintf("// This header file was created by setup.go \n// It simply connects gurobi.go to the local distribution (along with the cgo directives in %v\n\n", GoLibraryFilename)
	simpleImport := fmt.Sprintf("#include <%v/include/gurobi_c.h>\n", sf.GurobiHome)

	// Write all directives to file
	_, err = f.WriteString(fmt.Sprintf("%v%v", simpleComment, simpleImport))
	if err != nil {
		return err
	}

	// Return nil if everything went well.
	return nil
}

func main() {

	sf, err := GetDefaultSetupFlags()

	// Next, parse the arguments to make_lib and assign values to the mlf appropriately
	sf, err = ParseMakeLibArguments(sf)

	fmt.Println(sf)
	fmt.Println(err)

	// Write File
	err = WriteLibGo(sf)
	if err != nil {
		fmt.Println(err)
	}

	err = WriteHeaderFile(sf)

}
