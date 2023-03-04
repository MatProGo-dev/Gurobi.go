package gurobi

// #include <gurobi_passthrough.h>
import "C"
import "errors"

type Error struct {
	errcode int32
	message string
}

func (err Error) Error() string {
	return err.message
}

// make an error object from error code.
func (env *Env) makeError(errcode C.int) error {
	if env == nil {
		return errors.New("This environment has not initialized yet.")
	}

	if errcode != 0 {
		return Error{int32(errcode), C.GoString(C.GRBgeterrormsg(env.env))}
	}

	return nil
}

func (model *Model) makeError(errcode C.int) error {
	return model.Env.makeError(errcode)
}
