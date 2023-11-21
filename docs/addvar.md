# addVar

One of the fundamental methods in Gurobi's API
is the method `addVar` which creates new
gurobi variables.

It has some inputs that I was not familiar with:
(`vind` and `vval`) which were clarified by the example included below.

What it shows is that, if there are constraints that exist before you create a variable,
you can modify those constraints using this special method in the C API.

Below, I use it to modify an empty constraint to become an upper bound constraint
on the single variable in this simple LP.


## Example Script

```Go

package main

// #include <gurobi_passthrough.h>
import "C"
import (
	"fmt"
	"github.com/MatProGo-dev/Gurobi.go/gurobi"
	"os"
)

func main() {
	testIndex := 8
	testName := fmt.Sprintf("testmodel-addvars%v", testIndex)

	env0, err := gurobi.NewEnv(testName + `.log`)
	if err != nil {
		fmt.Printf("unexpected error creating new environment: %v", err)
		os.Exit(1)
	}

	model0, err := gurobi.NewModel(testName+`-model`, env0)
	if err != nil {
		fmt.Printf("unexpected error creating new model: %v", err)
		os.Exit(2)
	}
	defer model0.Free()

	// Define empty constraints
	_, err = model0.AddConstr(
		[]*gurobi.Var{},
		[]float64{},
		gurobi.SenseLessThan,
		2.0,
		"max-2",
	)
	if err != nil {
		fmt.Printf("error occurred while creating constraint #0: %v", err)
		os.Exit(3)
	}

	_, err = model0.AddConstr(
		[]*gurobi.Var{},
		[]float64{},
		gurobi.SenseLessThan,
		10.0,
		"max-10",
	)
	if err != nil {
		fmt.Printf("error occurred while creating constraint #1: %v", err)
		os.Exit(3)
	}

	// Create sparse matrix definition
	ind := make([]int32, 1)
	ind[0] = 1

	columns := make([]float64, 1)
	columns[0] = 0.5

	pind := (*C.int)(&ind[0])
	pcolumns := (*C.double)(&columns[0])

	// Call GRBaddvar
	errCode := C.GRBaddvar(
		(*C.GRBmodel)(model0.AsGRBModel),
		C.int(1),
		pind, pcolumns,
		C.double(-1.0),
		C.double(-1e2), C.double(1e2),
		C.char(gurobi.CONTINUOUS),
		C.CString("ellie-goulding"),
	)
	if errCode != 0 {
		fmt.Printf("unexpected error: %v!", errCode)
		os.Exit(4)
	}

	// Integrate new variables.
	if err := model0.Update(); err != nil {
		panic(err.Error())
	}

	// Optimize model
	if err := model0.Optimize(); err != nil {
		panic(err.Error())
	}

	// Write model to 'qp1.lp'.
	if err := model0.Write(testName + ".lp"); err != nil {
		panic(err.Error())
	}

	// Capture solution information
	optimstatus, err := model0.GetIntAttr(gurobi.INT_ATTR_STATUS)
	if err != nil {
		panic(err.Error())
	}

	objval, err := model0.GetDoubleAttr(gurobi.DBL_ATTR_OBJVAL)
	if err != nil {
		panic(err.Error())
	}

	sol, err := model0.GetDoubleAttrVars(
		gurobi.DBL_ATTR_X,
		[]*gurobi.Var{
			&gurobi.Var{Index: 0},
		},
	)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("\nOptimization complete\n")
	if optimstatus == gurobi.OPTIMAL {
		fmt.Printf("Optimal objective: %.4e\n", objval)
		fmt.Printf("  x=%.4f\n", sol[0])
	} else if optimstatus == gurobi.INF_OR_UNBD {
		fmt.Printf("Model is infeasible or unbounded\n")
	} else {
		fmt.Printf("Optimization was stopped early\n")
	}

}

```