package gurobi_test

import (
	"fmt"
	"os"
	"testing"

	gurobi "github.com/MatProGo-dev/Gurobi.go/gurobi"
)

func Test_LP1(t *testing.T) {
	// Create environment.
	env, err := gurobi.NewEnv("lp.log")
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove("lp.log")
	defer env.Free()

	// Create an empty model.
	model, err := gurobi.NewModel("lp-hypercube1", env)
	if err != nil {
		panic(err.Error())
	}
	defer model.Free()

	// Add varibles
	x, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		panic(err.Error())
	}
	y, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "y", []*gurobi.Constr{}, []float64{})
	if err != nil {
		panic(err.Error())
	}

	// Integrate new variables.
	if err := model.Update(); err != nil {
		panic(err.Error())
	}

	// Set Objective function
	expr := gurobi.LinExpr{}
	expr.AddTerm(x, 1).AddTerm(y, 1)
	if err := model.SetObjective(&expr, gurobi.MINIMIZE); err != nil {
		panic(err.Error())
	}

	// Add box constraints (four of them)
	if _, err = model.AddConstr([]*gurobi.Var{x}, []float64{1}, '<', 1.0, "x_upper_bound"); err != nil {
		panic(err.Error())
	}
	if _, err = model.AddConstr([]*gurobi.Var{x}, []float64{-1}, '<', 1.0, "x_lower_bound"); err != nil {
		panic(err.Error())
	}
	if _, err = model.AddConstr([]*gurobi.Var{y}, []float64{1}, '<', 1.0, "y_upper_bound"); err != nil {
		panic(err.Error())
	}
	if _, err = model.AddConstr([]*gurobi.Var{y}, []float64{-1}, '<', 1.0, "y_lower_bound"); err != nil {
		panic(err.Error())
	}

	// Optimize model
	if err := model.Optimize(); err != nil {
		panic(err.Error())
	}

	// Write model to 'lp-hyperbox1.lp'.
	if err := model.Write("lp-hyperbox1.lp"); err != nil {
		panic(err.Error())
	}
	defer os.Remove("lp-hyperbox1.lp")

	// Capture solution information
	optimstatus, err := model.GetIntAttr(gurobi.INT_ATTR_STATUS)
	if err != nil {
		panic(err.Error())
	}

	objval, err := model.GetDoubleAttr(gurobi.DBL_ATTR_OBJVAL)
	if err != nil {
		panic(err.Error())
	}

	sol, err := model.GetDoubleAttrVars(gurobi.DBL_ATTR_X, []*gurobi.Var{x, y})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("\nOptimization complete\n")
	if optimstatus == gurobi.OPTIMAL {
		fmt.Printf("Optimal objective: %.4e\n", objval)
		fmt.Printf("  x=%.4f, y=%.4f\n", sol[0], sol[1])
	} else if optimstatus == gurobi.INF_OR_UNBD {
		fmt.Printf("Model is infeasible or unbounded\n")
	} else {
		fmt.Printf("Optimization was stopped early\n")
	}

	// Checks
	if sol[0] != -1.0 {
		t.Errorf("The value of x is %v, not -1.0.", sol[0])
	}

	if sol[1] != -1.0 {
		t.Errorf("The value of y is %v, not -1.0.", sol[1])
	}

}
