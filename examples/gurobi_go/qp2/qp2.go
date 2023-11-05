package main

import (
	"fmt"

	gurobi "github.com/MatProGo-dev/Gurobi.go/gurobi"
)

func main() {
	// Create environment.
	exampleName := "qp2"
	env, err := gurobi.NewEnv(exampleName + ".log")
	if err != nil {
		panic(err.Error())
	}
	defer env.Free()

	// Create an empty model.
	model, err := gurobi.NewModel(exampleName+".model", env)
	if err != nil {
		panic(err.Error())
	}
	defer model.Free()

	// Add varibles
	x, err := model.AddVar(gurobi.CONTINUOUS, 0.0, 0.0, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		panic(err.Error())
	}
	y, err := model.AddVar(gurobi.CONTINUOUS, 0.0, 0.0, gurobi.INFINITY, "y", []*gurobi.Constr{}, []float64{})
	if err != nil {
		panic(err.Error())
	}
	z, err := model.AddVar(gurobi.CONTINUOUS, 0.0, 0.0, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
	if err != nil {
		panic(err.Error())
	}

	// Integrate new variables.
	if err := model.Update(); err != nil {
		panic(err.Error())
	}

	// Set Objective function
	expr := gurobi.QuadExpr{}
	expr.AddTerm(x, 2).AddQTerm(x, x, 1).AddQTerm(x, y, 1).AddQTerm(y, y, 1).AddQTerm(y, z, 1).AddQTerm(z, z, 1)
	if err := model.SetObjective(&expr, gurobi.MINIMIZE); err != nil {
		panic(err.Error())
	}

	// First constraint
	if _, err = model.AddConstr([]*gurobi.Var{x, y, z}, []float64{1, 2, 3}, '>', 4.0, "c0"); err != nil {
		panic(err.Error())
	}

	// Second constraint
	if _, err = model.AddConstr([]*gurobi.Var{x, y, z}, []float64{1, 1, 1}, '>', 1.0, "c1"); err != nil {
		panic(err.Error())
	}

	// Optimize model
	if err := model.Optimize(); err != nil {
		panic(err.Error())
	}

	// Write model to 'qp1.lp'.
	if err := model.Write("qp1.lp"); err != nil {
		panic(err.Error())
	}

	// Capture solution information
	optimstatus, err := model.GetIntAttr(gurobi.INT_ATTR_STATUS)
	if err != nil {
		panic(err.Error())
	}

	objval, err := model.GetDoubleAttr(gurobi.DBL_ATTR_OBJVAL)
	if err != nil {
		panic(err.Error())
	}

	sol, err := model.GetDoubleAttrVars(gurobi.DBL_ATTR_X, []*gurobi.Var{x, y, z})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("\nOptimization complete\n")
	if optimstatus == gurobi.OPTIMAL {
		fmt.Printf("Optimal objective: %.4e\n", objval)
		fmt.Printf("  x=%.4f, y=%.4f, z=%.4f\n", sol[0], sol[1], sol[2])
	} else if optimstatus == gurobi.INF_OR_UNBD {
		fmt.Printf("Model is infeasible or unbounded\n")
	} else {
		fmt.Printf("Optimization was stopped early\n")
	}
}
