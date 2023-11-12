package setup

import (
	"fmt"
	"os"
	"strings"
)

/*
setup_flags.go
Description:
	Helper functions and methods for setting up Gurobi.go
*/

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

	fmt.Printf("- Highest Version detected was %v \n", highestVersion)

	// Write the highest version's directory into the GurobiHome variable
	mlf.GurobiHome, err = CreateGurobiHomeDirectory(highestVersion)
	if err != nil {
		return mlf, err
	}

	return mlf, nil

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
			fmt.Printf("Unrecognized input: %v\n", os.Args[argIndex])
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

	GurobiHomeNameWithoutStart := mlf.GurobiHome[GurobiWordIndexStart:]
	GurobiDirNameIndexEnd := strings.Index(GurobiHomeNameWithoutStart, "/")

	return StringToGurobiVersionInfo(
		string(GurobiHomeNameWithoutStart[:GurobiDirNameIndexEnd]),
	)

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
	fmt.Println("- Created package line for lib.go file.")

	// 2. Create CXX_FLAGS argument
	cxxDirective, err := CreateCXXFlagsDirective(sf)
	if err != nil {
		return err
	}
	fmt.Println("- Created CXX Flags Directive line for lib.go file.")

	// 3. Create LDFLAGS Argument
	ldflagsDirective, err := CreateLDFlagsDirective(sf)
	if err != nil {
		return err
	}
	fmt.Println("- Created LDFlags directive line for lib.go file.")

	// Now Write to File
	f, err := os.Create(sf.GoFilename)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Println("- Created empty go file.")

	// Write all directives to file
	_, err = f.WriteString(fmt.Sprintf("%v%v%v import \"C\"\n", packageDirective, cxxDirective, ldflagsDirective))
	if err != nil {
		return err
	}
	fmt.Println("- Wrote all package lines to lib.go file.")

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
