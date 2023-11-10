[![codecov](https://codecov.io/gh/MatProGo-dev/Gurobi.go/graph/badge.svg?token=PLWY38PGZA)](https://codecov.io/gh/MatProGo-dev/Gurobi.go)
# Gurobi.go
This is a simple wrapper that goes around the C API for Gurobi.
When this module is installed properly (with `go generate`),
you can use it to set up Mathematical Programs with or
without using the [MatProInterface](https://github.com/MatProGo-dev/MatProInterface.go).

| ![](images/qp1/scalar-range-optimization1.png) | ![](images/qp1/geogebra-export1-yz-slice.png) |
|:----------------------------------------------:|:---------------------------------------------:|
|    Effectively Solve Mathematical Programs     |  Using the Tools of Mathematical Programming  |

The above example will be discussed in more detail in a Wiki or example for this library.

## Example Usage

To get a sense of how to use Gurobi.go, feel free to use the examples in the `examples` directory as a guide. You can solve
optimization problems using the `gurobi` package in this library or using `MatProInterface.go`, a modelling library developed
by the MatProGo-dev team.

To illustrate how to use each, we solve the QP above using each option.

<details>
  <summary>gurobi</summary>
  
  ```
  package main

  import (
      "fmt"
  
      gurobi "github.com/MatProGo-dev/Gurobi.go/gurobi"
  )
  
  func main() {
      // Constants
      exampleName := "gurobi_go-qp1"
      // Create environment.
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
      x1, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "x", []*gurobi.Constr{}, []float64{})
      if err != nil {
          panic(err.Error())
      }
      x2, err := model.AddVar(gurobi.CONTINUOUS, 0.0, -gurobi.INFINITY, gurobi.INFINITY, "y", []*gurobi.Constr{}, []float64{})
      if err != nil {
          panic(err.Error())
      }
  
      // Integrate new variables.
      if err := model.Update(); err != nil {
          panic(err.Error())
      }
  
      // Set Objective function
      expr := gurobi.QuadExpr{}
      expr.AddTerm(x2, -0.97).AddQTerm(x1, x1, 1).AddQTerm(x1, x2, 0.5).AddQTerm(x2, x2, 0.25)
      if err := model.SetObjective(&expr, gurobi.MINIMIZE); err != nil {
          panic(err.Error())
      }
  
      // First constraint
      if _, err = model.AddConstr([]*gurobi.Var{x1, x2}, []float64{1.0, 0.0}, '>', 0.0, "c0"); err != nil {
          panic(err.Error())
      }
  
      // Second constraint
      if _, err = model.AddConstr([]*gurobi.Var{x1, x2}, []float64{0.0, 1}, '>', 1.0, "c1"); err != nil {
          panic(err.Error())
      }
  
      // Third constraint
      if _, err = model.AddConstr([]*gurobi.Var{x1, x2}, []float64{1.0, 0.0}, '<', 2.0, "c2"); err != nil {
          panic(err.Error())
      }
  
      // Fourth constraint
      if _, err = model.AddConstr([]*gurobi.Var{x1, x2}, []float64{0.0, 1}, '<', 3.0, "c3"); err != nil {
          panic(err.Error())
      }
  
      // Optimize model
      if err := model.Optimize(); err != nil {
          panic(err.Error())
      }
  
      // Write model to 'qp1.lp'.
      if err := model.Write(exampleName + ".lp"); err != nil {
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
  
      sol, err := model.GetDoubleAttrVars(gurobi.DBL_ATTR_X, []*gurobi.Var{x1, x2})
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
  }
  ```

</details>

<details>
  <summary>MatProInterface.go</summary>

  ```
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
  ```

</details>

### Notes

Note: As of November 2023, only QP1 has been tested.

## Installation

Warning: The setup script is designed to only work on Mac OS X. If you are interested in using this on a Windows machine, then there are no guarantees that it will work.

### Installation in Module

If you are using this as part of another Go project, then you will need to install this as a module ready for import.
To do this, follow the instructions below:

1. Use a "-d" `go get -d github.com/MatProGo-dev/Gurobi.go/gurobi`. Pay attention to which version appears in your terminal output.
2. Enter Go's internal installation of gurobi.go. For example, run `cd ~/go/pkg/mod/github.com/MatProGo-dev/Gurobi.go@v0.0.0-20221111000100-e629c3f29605` where the suffix is the version number from the previous output.
3. Run go generate with sudo privileges from this installation. `sudo go generate`.

### Development Installation

If you wish to improve upon Gurobi.go, then you can simply clone the repository into your local file system and then run `go generate`.

1. Clone the library using `git clone github.com/MatProGo-dev/Gurobi.go `
2. Enter the repository: `cd Gurobi.go`.
3. Run the setup script from inside the cloned repository: `go generate`.

## LICENSE
See [LICENSE](LICENSE).

## To-Dos

* [ ] Create Documentation for how to implement QP1 with
  * [ ] MatProInterface.go
  * [ ] Gurobi.go
* [ ] Set up CodeCoverage pipeline.