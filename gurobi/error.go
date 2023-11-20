package gurobi

// #include <gurobi_passthrough.h>
import "C"
import (
	"errors"
	"fmt"
)

/*
Error Objects
*/

// Gurobi error object (not the same as a regular error)
type Error struct {
	ErrorCode int32
	Message   string
}

type MismatchedLengthError struct {
	Length1 int
	Length2 int
	Name1   string // Name of the object with Length1
	Name2   string // Name of the object with Length2
}

/*
Error Methods
*/

func (err Error) Error() string {
	return err.Message
}

func (err MismatchedLengthError) Error() string {
	// Assemble string
	return fmt.Sprintf(
		"the length of %v (%v) must match that of %v (%v)!",
		err.Name1,
		err.Length1,
		err.Name2,
		err.Length2,
	)
}

/*
Other Error-Related methods
*/

// make an error object from error code.
func (env *Env) MakeError(errcode C.int) error {
	if env == nil {
		return errors.New("This environment has not initialized yet.")
	}

	if errcode != 0 {
		return Error{int32(errcode), C.GoString(C.GRBgeterrormsg(env.env))}
	}

	return nil
}

func (env *Env) MakeUninitializedError() error {
	return fmt.Errorf("The gurobi environment was not yet initialized!")
}

func (model *Model) MakeError(errcode C.int) error {
	return model.Env.MakeError(errcode)
}
