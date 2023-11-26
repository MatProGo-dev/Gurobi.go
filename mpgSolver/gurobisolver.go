package mpgSolver

import (
	"fmt"
	"github.com/MatProGo-dev/MatProInterface.go/optim"
	"io"
	"log"
	"os"

	gurobi "github.com/MatProGo-dev/Gurobi.go/gurobi"
)

// Type Definition

type GurobiSolver struct {
	Env                    *gurobi.Env
	CurrentModel           *gurobi.Model
	ModelName              string
	GoopIDToGurobiIndexMap map[uint64]int32 // Maps each Goop ID (uint64) to the idx value used for each Gurobi variable.
}

// Function

/*
NewGurobiSolver
Description:

	Create a new gurobi solver object.
*/
func NewGurobiSolver(modelName string) GurobiSolver {
	// Constants

	// Algorithm
	newGS := GurobiSolver{}
	newGS.CreateModel(modelName)

	return newGS

}

/*
ShowLog
Description:

	Decides whether or not to print logs to the terminal?
*/
func (gs *GurobiSolver) ShowLog(tf bool) error {
	// Constants
	logFileName := gs.ModelName + ".txt"

	// Check to see if logFile exists
	_, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		//Do Nothing. The later lines will create the file.
	} else {
		//Delete the old file.
		err = os.Remove(logFileName)
		if err != nil {
			return fmt.Errorf("There was an issue deleting the old log file: %v", err)
		}
	}

	// Create Logging file
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// log.Fatal(err)
		return fmt.Errorf("There was an issue createing a log file: %v", err)
	}

	// Attach logger to terminal only if tf is true
	if tf {
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		log.SetOutput(file)
	}

	// Create initial file
	log.Println("Log file created.")

	return nil

}

/*
SetTimeLimit
Description:

	Sets the time limit of the current model in gurobi solver gs.

Input:

	limitInS = Value of time limit in seconds (float)
*/
func (gs *GurobiSolver) SetTimeLimit(limitInS float64) error {

	err := gs.Env.SetDBLParam("TimeLimit", limitInS)
	if err != nil {
		return fmt.Errorf("There was an issue using SetDBLParam(): %v", err)
	}

	// If there was no error, return nil
	return nil
}

/*
GetTimeLimit
Description:

	Gets the time limit of the current model in gurobi solver gs.

Input:

	None

Output

	limitInS = Value of time limit in seconds (float)
*/
func (gs *GurobiSolver) GetTimeLimit() (float64, error) {

	limitOut, err := gs.Env.GetDBLParam("TimeLimit")
	if err != nil {
		return -1, fmt.Errorf("There was an error getting the double param TimeLimit: %v", err)
	}

	// If all things succeeded, return good data.
	return limitOut, err
}

/*
CreateModel
Description:
*/
func (gs *GurobiSolver) CreateModel(modelName string) {
	// Constants

	// Algorithm
	env, err := gurobi.NewEnv(modelName + ".log")
	if err != nil {
		panic(err.Error())
	}

	gs.Env = env

	// Create an empty model.
	model, err := gurobi.NewModel(modelName, env)
	if err != nil {
		panic(err.Error())
	}
	gs.CurrentModel = model

	// Create an empty map
	gs.GoopIDToGurobiIndexMap = make(map[uint64]int32)

}

/*
LoadMPIModel()
Description:
	This loads a model, saved as a MatProInterface.Model object.
*/

/*
FreeEnv
Description:

	Frees the Env() method. Useful after the problem is solved.
*/
func (gs *GurobiSolver) FreeEnv() {
	gs.Env.Free()
}

/*
FreeModel
Description

	Frees the Model member. Useful after the problem is solved.
*/
func (gs *GurobiSolver) FreeModel() {
	gs.CurrentModel.Free()
}

/*
Free
Description:

	Frees the Env and Model elements of the system.
*/
func (gs *GurobiSolver) Free() {
	gs.FreeModel()
	gs.FreeEnv()
}

/*
AddVariable
Description:

	Adds a variable to the Gurobi Model.
*/
func (gs *GurobiSolver) AddVariable(varIn optim.Variable) error {
	// Constants

	// Convert Variable Type
	vType, err := VarTypeToGRBVType(varIn.Vtype)
	if err != nil {
		return fmt.Errorf("There was an error defining gurobi type: %v", err)
	}

	// Add Variable to Current Model
	tempVar, err := gs.CurrentModel.AddVar(int8(vType), 0.0, varIn.Lower, varIn.Upper, fmt.Sprintf("x%v", varIn.ID), []*gurobi.Constr{}, []float64{})

	fmt.Printf("%v: L=%v, U=%v, name=%v\n", int8(vType), varIn.Lower, varIn.Upper, fmt.Sprintf("x%v", varIn.ID))

	// Update Map from GoopID to Gurobi Idx
	gs.GoopIDToGurobiIndexMap[varIn.ID] = int32(len(gs.CurrentModel.Variables) - 1)
	if int32(len(gs.CurrentModel.Variables)-1) != tempVar.Index {
		fmt.Println("Hehe; This is Kwesi's secret message to tell you that mpgSolver.AddVariable() is not working as expected!")
	}

	return err
}

func VarTypeToGRBVType(vtype optim.VarType) (int8, error) {
	switch vtype {
	case optim.Continuous:
		return gurobi.CONTINUOUS, nil
	case optim.Binary:
		return gurobi.BINARY, nil
	}

	return gurobi.BINARY, fmt.Errorf("Unexpected mpg variable type for conversion: %v", vtype)
}

/*
AddVariables
Description:

	Adds a set of variables to the Gurobi Model.
*/
func (gs *GurobiSolver) AddVariables(varSliceIn []optim.Variable) error {
	// Constants

	// Iterate through ALL variable address in varSliceIn
	for _, tempVar := range varSliceIn {
		err := gs.AddVariable(tempVar)
		if err != nil {
			// Terminate early.
			return fmt.Errorf("Error in AddVariable(): %v", err)
		}
	}

	// If we successfully made it through all Variable objects, then return no errors.
	return nil
}

/*
AddConstraint
Description:

	Adds a single constraint to the gurobi model object inside of the current GurobiSolver object.
*/
func (gs *GurobiSolver) AddConstraint(constrIn optim.Constraint, errors ...error) error {
	// Input Checking
	err := optim.CheckErrors(errors)
	if err != nil {
		return err
	}

	if !optim.IsConstraint(constrIn) {
		return fmt.Errorf("The input to AddConstr is not recognized as a constraint!")
	}

	// Constants
	switch constrIn.(type) {
	case optim.ScalarConstraint:
		// Cast and simplify
		constrAsSC, _ := constrIn.(optim.ScalarConstraint)

		simplifiedConstr, err := constrAsSC.Simplify()
		if err != nil {
			fmt.Println("1")
			return err
		}

		if tf, _ := simplifiedConstr.IsLinear(); !tf {
			return fmt.Errorf("cannot handle quadratic constraints yet in Gurobi.go; create an issue if you want this feature!")
		}

		gurobiVarSlice, L, senseOut, C, err := gs.ToGurobiLinearConstraint(simplifiedConstr)
		if err != nil {
			fmt.Println("2")
			return err
		}

		for ii, v := range gurobiVarSlice {
			fmt.Printf("var slice post TGLC: slice[%v] = %v\n", ii, v)
		}

		// Call Gurobi library's AddConstr() function
		_, err = gs.CurrentModel.AddConstr(
			gurobiVarSlice, L, senseOut, C,
			fmt.Sprintf("goop Constraint #%v", len(gs.CurrentModel.Constraints)),
		)
		if err != nil {
			return fmt.Errorf("There was an issue with adding the constraint to the gurobi model: %v", err)
		}
	case optim.VectorConstraint:
		// Cast
		constrAsVC, _ := constrIn.(optim.VectorConstraint)

		if constrAsVC.Check() != nil {
			return constrAsVC.Check()
		}

		// Extract the Scalar Constraint from Each element of this vector constraint
		// TODO: Finish logic here so that we can extract scalar constraints from an arbitrary vector constraint.
		for vecIdx := 0; vecIdx < constrAsVC.LeftHandSide.Len(); vecIdx++ {
			tempConstr, err := constrAsVC.AtVec(vecIdx)
			if err != nil {
				return err
			}

			err = gs.AddConstraint(tempConstr)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Unexpected type of constraint input: %T (%v)", constrIn, constrIn)
	}

	// Create no errors if there were no errors!
	return nil
}

/*
SetObjective
Description:

	This algorithm should set the objective based on the value of the expression provided as input to this function.
*/
func (gs *GurobiSolver) SetObjective(objIn optim.Objective) error {

	objExpression := objIn.ScalarExpression

	// Handle this differently for different types of expression inputs
	switch objExpression.(type) {
	case optim.ScalarLinearExpr:
		gurobiLE := &gurobi.LinExpr{}
		for varIndex, goopIndex := range objExpression.IDs() {
			gurobiIndex := gs.GoopIDToGurobiIndexMap[goopIndex]

			// Add each linear term to the expression.
			tempGurobiVar := gurobi.Var{
				Model: gs.CurrentModel,
				Index: gurobiIndex,
			}
			gurobiLE = gurobiLE.AddTerm(&tempGurobiVar, objExpression.Coeffs()[varIndex])
		}

		// Add a constant term to the expression
		gurobiLE = gurobiLE.AddConstant(objExpression.Constant())

		fmt.Println(gurobiLE)

		// Add linear expression to the objective.
		err := gs.CurrentModel.SetLinearObjective(gurobiLE, int32(objIn.Sense))
		if err != nil {
			return fmt.Errorf("There was an issue setting the linear objective with SetLinearObjective(): %v", err)
		}

		return nil

	case optim.ScalarQuadraticExpression:
		objExpressionAsQE := objExpression.(optim.ScalarQuadraticExpression)
		gurobiQE := &gurobi.QuadExpr{}

		// Create quadratic part of quadratic expression
		for varIndex1, goopIndex1 := range objExpression.IDs() {
			gurobiIndex1 := gs.GoopIDToGurobiIndexMap[goopIndex1]

			for varIndex2, goopIndex2 := range objExpression.IDs() {
				gurobiIndex2 := gs.GoopIDToGurobiIndexMap[goopIndex2]

				// Add each linear term to the expression.
				tempGurobiVar1 := gurobi.Var{
					Model: gs.CurrentModel,
					Index: gurobiIndex1,
				}
				tempGurobiVar2 := gurobi.Var{
					Model: gs.CurrentModel,
					Index: gurobiIndex2,
				}

				gurobiQE = gurobiQE.AddQTerm(&tempGurobiVar1, &tempGurobiVar2, objExpressionAsQE.Q.At(varIndex1, varIndex2))
			}
		}

		// Create linear part of quadratic expression
		for varIndex, goopIndex := range objExpression.IDs() {
			gurobiIndex := gs.GoopIDToGurobiIndexMap[goopIndex]

			// Add each linear term to the expression.
			tempGurobiVar := gurobi.Var{
				Model: gs.CurrentModel,
				Index: gurobiIndex,
			}
			gurobiQE = gurobiQE.AddTerm(&tempGurobiVar, objExpressionAsQE.L.AtVec(varIndex))
		}

		// Create offset
		gurobiQE = gurobiQE.AddConstant(objExpressionAsQE.C)

		// Return
		fmt.Println(*gurobiQE)

		err := gs.CurrentModel.SetQuadraticObjective(gurobiQE, int32(objIn.Sense))
		if err != nil {
			return fmt.Errorf("There was an issue setting the quadratic objective with SetQuadraticObjective(): %v", err)
		}

		return nil

	default:
		return fmt.Errorf("Unexpected objective type given to gurobisolver's SetObjective(): %T", objExpression)
	}
}

/*
Optimize
Description:
*/
func (gs *GurobiSolver) Optimize() (optim.Solution, error) {
	// Make sure that all changes are applied to the given model.
	err := gs.CurrentModel.Update()
	if err != nil {
		return optim.Solution{}, fmt.Errorf("There was an issue updating the current gurobi model: %v", err)
	}

	// Optimize
	err = gs.CurrentModel.Optimize()
	if err != nil {
		return optim.Solution{}, fmt.Errorf("There was an issue optimizing the current model: %v", err)
	}

	// Construct solution:
	// - Status
	tempSolution := optim.Solution{}
	tempStatus, err := gs.CurrentModel.GetIntAttr("Status")
	if err != nil {
		return tempSolution, fmt.Errorf("There was an issue collecting the model's status: %v", err)
	}
	tempSolution.Status = optim.OptimizationStatus(tempStatus)

	// - Values
	tempValues := make(map[uint64]float64)
	for _, tempGurobiVar := range gs.CurrentModel.Variables {
		val, err := tempGurobiVar.GetDouble("X")
		if err != nil {
			return tempSolution, fmt.Errorf("Error while retrieving the optimal values of the problem: %v", err)
		}
		// identify goop index that has this gurobi variables data
		for goopIndex, gurobiIndex := range gs.GoopIDToGurobiIndexMap {
			if gurobiIndex == tempGurobiVar.Index {
				tempValues[goopIndex] = val
				break // When you find it, save the value and return the value to the map.
			}
		}
	}
	tempSolution.Values = tempValues

	// - Objective
	tempObjective, err := gs.CurrentModel.GetDoubleAttr("ObjVal")
	if err != nil {
		return tempSolution, fmt.Errorf("There was an issue getting the objective value of the current model.")
	}
	tempSolution.Objective = tempObjective

	// All steps were successful, return solution!
	return tempSolution, nil
}

/*
DeleteSolver
Description:

	Attempts to delete all info about the current solver.
*/
func (gs *GurobiSolver) DeleteSolver() error {
	// Free model and environment
	gs.CurrentModel.Free()

	gs.Env.Free()

	return nil
}

/*
OptimizeModel
Description:

	Getting
*/
func Solve(model optim.Model) (optim.Solution, GurobiSolver, error) {
	// Create GurobiSolver
	solver := NewGurobiSolver(model.Name + "_GurobiSolver")

	// Add Variables
	err := solver.AddVariables(model.Variables)
	if err != nil {
		return optim.Solution{},
			solver,
			fmt.Errorf("error adding MPG variables to gurobi model: %v", err)
	}

	// Add Constraints
	for i, constraint := range model.Constraints {
		err = solver.AddConstraint(constraint)
		if err != nil {
			return optim.Solution{},
				solver,
				fmt.Errorf(
					"there was an issue adding %v-th constraint: %v",
					i, constraint,
				)
		}
	}

	// Add Objective
	err = solver.SetObjective(*model.Obj)
	if err != nil {
		return optim.Solution{},
			solver,
			fmt.Errorf(
				"there was an issue adding the model's objective: %v", err,
			)
	}

	// Call Solver
	sol, err := solver.Optimize()
	if err != nil {
		return optim.Solution{},
			solver,
			fmt.Errorf(
				"there was an issue with optimizing the model: %v",
				err,
			)
	}

	// Return final solution
	return sol, solver, err

}

func (gs *GurobiSolver) ToGurobiLinearConstraint(constr optim.ScalarConstraint) (
	[]*gurobi.Var, []float64, int8, float64, error,
) {
	// constr
	constrSimplified, err := constr.Simplify()
	if err != nil {
		return nil, nil, int8(-1), -1, err
	}

	// RightHandSide now contains just a constant

	// Algorithm
	switch left := constrSimplified.LeftHandSide.(type) {
	case optim.Variable:
		copiedConstr := constrSimplified
		copiedConstr.LeftHandSide = left.ToScalarLinearExpression()
		return gs.ToGurobiLinearConstraint(copiedConstr)
	case optim.ScalarLinearExpr:
		// Create slice of gurobi.Var objects that matches whats in expr
		tempVarSlice := make([]*gurobi.Var, left.X.Len())
		newL := make([]float64, left.L.Len())
		for GoopIdx, tempGoopID := range left.IDs() {
			tempGurobiIdx := gs.GoopIDToGurobiIndexMap[tempGoopID]

			fmt.Printf("Gurobi Index: %v, MPG Index: %v\n", tempGurobiIdx, tempGoopID)

			// Locate the gurobi variable in the current model that has matching ID
			for jj, tempGurobiVar := range gs.CurrentModel.Variables {
				if tempGurobiIdx == tempGurobiVar.Index {
					tempVarSlice[GoopIdx] = &(gs.CurrentModel.Variables[jj])
					newL[GoopIdx] = left.L.AtVec(GoopIdx)
				}
			}
		}

		fmt.Printf("tempVarSlice = %v\n", tempVarSlice)
		fmt.Printf("L: %v\n", left.L)
		fmt.Printf(" POST tempVarSlice[0].ID: %v\n", tempVarSlice[0].Index)

		// Return
		return tempVarSlice, newL, int8(constrSimplified.Sense), float64(constrSimplified.Right().(optim.K)) - left.C, nil

	default:
		return nil, nil, int8(-1), -1, fmt.Errorf("unexpected left hand side input of type %T", left)

	}

}
