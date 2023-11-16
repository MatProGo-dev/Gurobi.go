package gurobi

// #include <gurobi_passthrough.h>
import "C"
import (
	"errors"
	"fmt"
)

type Error struct {
	ErrorCode int32
	Message   string
}

/*
Error Methods
*/

func (err Error) Error() string {
	return err.Message
}

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

/*
MakeUninitializedError
Description:

	This function simply returns a fixed error for when the model is not initialized.
*/
func (model *Model) MakeUninitializedError() error {
	return fmt.Errorf("The gurobi model was not yet initialized!")
}
