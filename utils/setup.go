/*
setup.go
Description:
	An implementation of the file local_lib_info.go written entirely in go.
*/

package main

import (
	"fmt"
	"github.com/MatProGo-dev/Gurobi.go/setup"
	"os"
	"strings"
)

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

func main() {
	fmt.Println("Beginning setup for Gurobi.go.")
	sf, err := setup.GetDefaultSetupFlags()
	fmt.Println("Defined the default setup flags.")

	// Next, parse the arguments to make_lib and assign values to the mlf appropriately
	sf, err = setup.ParseMakeLibArguments(sf)
	fmt.Println("Parsed any command line arguments to setup.go")

	fmt.Println(sf)
	fmt.Println(err)

	// Write File
	err = setup.WriteLibGo(sf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Finished writing the lib.go file.")

	err = setup.WriteHeaderFile(sf)
	fmt.Println("Finished writing the header file for Gurobi.go.")

}
