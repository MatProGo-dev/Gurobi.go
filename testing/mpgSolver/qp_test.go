package mpgSolver_test

import (
	"fmt"
	"github.com/MatProGo-dev/Gurobi.go/mpgSolver"
	"github.com/MatProGo-dev/MatProInterface.go/optim"
	"gonum.org/v1/gonum/mat"
	"testing"
)

/*
TestQP1
Description:

	Creates a simple QP that we wish to minimize using the tools of Gurobi.
*/
func TestQP1(t *testing.T) {
	// Constants
	modelName := "mpg-qp1"
	m := optim.NewModel(modelName)
	x := m.AddVariableVector(2)

	gs := mpgSolver.NewGurobiSolver(
		fmt.Sprintf("solvertest-%v", modelName),
	)

	// Add Variables to Gurobi's Model
	err := gs.AddVariables(x.Elements)
	if err != nil {
		t.Errorf("There was an issue adding x[0] to gs: %v", err)
	}

	// Create Vector Variables
	c1 := optim.KVector(
		*mat.NewVecDense(2, []float64{0.0, 1.0}),
	)

	c2 := optim.KVector(
		*mat.NewVecDense(2, []float64{2.0, 3.0}),
	)

	// Use these to create constraints.
	err = gs.AddConstraint(
		x.LessEq(c2),
	)
	if err != nil {
		t.Errorf("There was an issue creating the proper vector constraint: %v", err)
	}

	err = gs.AddConstraint(
		x.GreaterEq(c1),
	)
	if err != nil {
		t.Errorf("There was an issue creating the proper vector constraint: %v", err)
	}

	// Create objective
	Q1 := optim.Identity(x.Len())
	Q1.Set(0, 1, 0.25)
	Q1.Set(1, 0, 0.25)
	Q1.Set(1, 1, 0.25)

	obj := optim.ScalarQuadraticExpression{
		Q: Q1,
		X: x,
		L: *mat.NewVecDense(x.Len(), []float64{0, -0.97}),
		C: 2.0,
	}

	// Add objective
	err = gs.SetObjective(optim.Objective{obj, optim.SenseMinimize})
	if err != nil {
		t.Errorf("There was an issue setting the objective of the Gurobi solver model: %v", err)
	}

	// Solve!
	sol, err := gs.Optimize()
	if err != nil {
		t.Errorf("There was an issue optimizing the QP: %v", err)
	}

	if len(sol.Values) != 2 {
		t.Errorf("Expected for there to be two variables in solution's values field; received %v", len(sol.Values))
	}

	if sol.Objective > 1.2 {
		t.Errorf("Expected objective to be less than 1.2; received %v", sol.Objective)
	}

	if (sol.Values[x.AtVec(0).(optim.Variable).ID] >= 2.0) && (sol.Values[x.AtVec(0).(optim.Variable).ID] <= 3.0) {
		t.Errorf(
			"Expected 2 <= x^*[0] <= 3. Received %v",
			sol.Values[x.AtVec(0).(optim.Variable).ID],
		)
	}

}
