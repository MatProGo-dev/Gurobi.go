package main

import (
	"fmt"
	"github.com/MatProGo-dev/Gurobi.go/mpgSolver"
	"gonum.org/v1/gonum/mat"

	gurobi "github.com/MatProGo-dev/Gurobi.go/gurobi"
	"github.com/MatProGo-dev/MatProInterface.go/optim"
)

func main() {
	// Constants
	exampleName := "matprointerface-qp1"
	// Create environment.
	env, err := gurobi.NewEnv(exampleName + ".log")
	if err != nil {
		panic(err.Error())
	}
	defer env.Free()

	// Create an empty model.
	model := optim.NewModel(exampleName + ".model")

	// Add varibles
	x := model.AddVariableVector(2)

	vd1 := mat.NewVecDense(2, []float64{0.0, 1.0})
	vd2 := mat.NewVecDense(2, []float64{2.0, 3.0})

	// Create constraints
	err = model.AddConstraint(x.GreaterEq(*vd1))
	if err != nil {
		panic(fmt.Sprintf("there was an issue creating the proper vector constraint: %v", err))
	}

	err = model.AddConstraint(x.LessEq(*vd2))
	if err != nil {
		panic(fmt.Sprintf("there was an issue creating the proper vector constraint: %v", err))
	}

	// Set Objective function
	Q1 := *mat.NewDense(2, 2, []float64{1.0, 0.25, 0.25, 0.25})

	prod1, _ := x.Transpose().Multiply(Q1)
	prod2, err := prod1.Multiply(x)
	if err != nil {
		panic(fmt.Sprintf("there was an issue creating product 2: %v", err))
	}

	sum, err := prod2.Plus(
		x.Transpose().Multiply(*mat.NewVecDense(2, []float64{0, -0.97})),
	)
	if err != nil {
		panic(fmt.Sprintf("There was an issue computing the final sum: %v", err))
	}

	err = model.SetObjective(sum, optim.SenseMinimize)
	if err != nil {
		panic(fmt.Sprintf("There was an issue setting the objective for the model: %v", err))
	}

	// Solve the above model with GurobiSolver
	sol, gs, err := mpgSolver.Solve(*model)
	if err != nil {
		panic(
			fmt.Errorf("error in solving model: %v", err),
		)
	}

	//// Write model to 'qp.lp'.
	fmt.Printf("%v\n%v\n%v\n", gs.GoopIDToGurobiIndexMap, x.Elements[0].ID, x.Elements[1].ID)
	if err := gs.CurrentModel.Write(exampleName + ".lp"); err != nil {
		panic(err.Error())
	}

	// Capture solution information
	optimstatus := sol.Status
	objval := sol.Objective

	fmt.Printf("\nOptimization complete\n")
	if optimstatus == gurobi.OPTIMAL {
		fmt.Printf("Optimal objective: %.4e\n", objval)
		fmt.Printf("  x_1=%.4f, x_2=%.4f\n", sol.Values[0], sol.Values[1])
	} else if optimstatus == gurobi.INF_OR_UNBD {
		fmt.Printf("Model is infeasible or unbounded\n")
	} else {
		fmt.Printf("Optimization was stopped early\n")
	}
}
