package gurobi

// #include <gurobi_passthrough.h>
import "C"
import (
	"errors"
	"fmt"
)

// Model ...
// Gurobi model object
type Model struct {
	AsGRBModel  *C.GRBmodel
	Env         Env
	Variables   []Var
	Constraints []Constr
}

// NewModel ...
// create a new model from the environment.
func NewModel(modelname string, env *Env) (*Model, error) {
	if env == nil {
		return nil, errors.New("This environment is not created yet.")
	}

	var model *C.GRBmodel
	errcode := C.GRBnewmodel(env.env, &model, C.CString(modelname), 0, nil, nil, nil, nil, nil)
	if errcode != 0 {
		return nil, env.MakeError(errcode)
	}

	newenv := C.GRBgetenv(model)
	if newenv == nil {
		return nil, errors.New("Failed retrieve the environment")
	}

	return &Model{AsGRBModel: model, Env: Env{newenv}}, nil
}

// Free ...
// free the model
func (model *Model) Free() {
	if model == nil {
		return
	}
	C.GRBfreemodel(model.AsGRBModel)
}

/*
AddVar
Description:

	Create a variable to the model
	This includes inputs for:
	- obj = Linear coefficient applied to this variable in the objective function (i.e. objective =  ... + obj * newvar + ...)
	- lb = Lower Bound
	- ub = Upper Bound
	- name = Name of the variable
	- constrs = Array of constraints in which the variable participates.
	- columns = Numerical values associated with non-zero values for the new variable?
				(I think these are coefficients appearing in the linear constraint for this variable)

Links:

	Comes from the C API for Gurobi.
	Documentation for 9.0: https://www.gurobi.com/documentation/9.0/refman/c_addvar.html
*/
func (model *Model) AddVar(vtype int8, obj float64, lb float64, ub float64, name string, constrs []*Constr, columns []float64) (*Var, error) {
	if model == nil {
		return nil, errors.New("model is not initialized")
	}

	if len(constrs) != len(columns) {
		return nil, errors.New("either the length of constrs or columns are wrong")
	}

	ind := make([]int32, len(constrs), 0)
	for i, c := range constrs {
		if c.idx < 0 {
			return nil, errors.New("Invalid index in constrs")
		}
		ind[i] = c.idx
	}

	pind := (*C.int)(nil)
	pval := (*C.double)(nil)
	if len(ind) > 0 {
		pind = (*C.int)(&ind[0])
		pval = (*C.double)(&columns[0])
	}

	err := C.GRBaddvar(model.AsGRBModel, C.int(len(constrs)), pind, pval, C.double(obj), C.double(lb), C.double(ub), C.char(vtype), C.CString(name))
	if err != 0 {
		return nil, model.MakeError(err)
	}

	if err := model.Update(); err != nil {
		return nil, err
	}

	model.Variables = append(model.Variables, Var{model, int32(len(model.Variables))})
	return &model.Variables[len(model.Variables)-1], nil
}

// AddVars ...
func (model *Model) AddVars(vtype []int8, obj []float64, lb []float64, ub []float64, name []string, constrs [][]*Constr, columns [][]float64) ([]*Var, error) {
	if model == nil {
		return nil, errors.New("")
	}

	if len(vtype) != len(obj) || len(obj) != len(lb) || len(lb) != len(ub) || len(ub) != len(name) || len(name) != len(constrs) || len(constrs) != len(columns) {
		return nil, errors.New("")
	}

	numnz := 0
	for _, constr := range constrs {
		numnz += len(constr)
	}

	beg := make([]int32, len(constrs))
	ind := make([]int32, numnz)
	val := make([]float64, numnz)
	k := 0
	for i := 0; i < len(constrs); i++ {
		if len(constrs[i]) != len(columns[i]) {
			return nil, errors.New("")
		}

		for j := 0; j < len(constrs[i]); j++ {
			idx := constrs[i][j].idx
			if idx < 0 {
				return nil, errors.New("")
			}
			ind[k+j] = idx
			val[k+j] = columns[i][j]
		}

		beg[i] = int32(k)
		k += len(constrs[i])
	}

	vname := make([](*C.char), len(vtype))
	for i, n := range name {
		vname[i] = C.CString(n)
	}

	pbeg := (*C.int)(nil)
	pind := (*C.int)(nil)
	pval := (*C.double)(nil)
	if len(beg) > 0 {
		pbeg = (*C.int)(&beg[0])
		pind = (*C.int)(&ind[0])
		pval = (*C.double)(&val[0])
	}

	pobj := (*C.double)(nil)
	plb := (*C.double)(nil)
	pub := (*C.double)(nil)
	pvtype := (*C.char)(nil)
	pname := (**C.char)(nil)
	if len(vtype) > 0 {
		pobj = (*C.double)(&obj[0])
		plb = (*C.double)(&lb[0])
		pub = (*C.double)(&ub[0])
		pvtype = (*C.char)(&vtype[0])
		pname = (**C.char)(&vname[0])
	}

	err := C.GRBaddvars(model.AsGRBModel, C.int(len(vtype)), C.int(numnz), pbeg, pind, pval, pobj, plb, pub, pvtype, pname)
	if err != 0 {
		return nil, model.MakeError(err)
	}

	if err := model.Update(); err != nil {
		return nil, err
	}

	vars := make([]*Var, len(vtype), 0)
	xcols := len(model.Variables)
	for i := xcols; i < xcols+len(vtype); i++ {
		model.Variables = append(model.Variables, Var{model, int32(i)})
		vars[i] = &model.Variables[len(model.Variables)-1]
	}
	return vars, nil
}

/*
AddConstr
Description:

	Add a Linear constraint into the model.
	Uses the GRBaddconstr() method from the C api.

Inputs:
  - vars: A slice of variable arrays which provide the indices for the gurobi model's variables.
  - val: A slice of float values which are used as coefficients for the variables in the linear constraint.
  - sense: A flag which determines if this is an equality, less than equal or greater than or equal constraint.
  - rhs: A float value which determines the constant which is on the other side of the constraint.
  - constrname: An optional name for the constraint.

Link:

	https://www.gurobi.com/documentation/9.1/refman/c_addconstr.html
*/
func (model *Model) AddConstr(vars []*Var, val []float64, sense int8, rhs float64, constrname string) (*Constr, error) {
	if model == nil {
		return nil, errors.New("")
	}

	ind := make([]int32, len(vars))
	for i, v := range vars {
		if v.Index < 0 {
			return nil, errors.New("Invalid vars")
		}
		ind[i] = v.Index
	}

	pind := (*C.int)(nil)
	pval := (*C.double)(nil)
	if len(ind) > 0 {
		pind = (*C.int)(&ind[0])
		pval = (*C.double)(&val[0])
	}

	fmt.Printf("pind = %v\n", *pind)
	fmt.Printf("ind = %v\n vars[0] = %v\n", ind, vars[0].Index)

	err := C.GRBaddconstr(
		model.AsGRBModel,
		C.int(len(ind)),
		pind, pval,
		C.char(sense), C.double(rhs), C.CString(constrname))
	if err != 0 {
		return nil, model.MakeError(err)
	}

	if err := model.Update(); err != nil {
		return nil, err
	}

	model.Constraints = append(model.Constraints, Constr{model, int32(len(model.Constraints))})
	return &model.Constraints[len(model.Constraints)-1], nil
}

func (model *Model) AddConstrs(vars [][]*Var, val [][]float64, sense []int8, rhs []float64, constrname []string) ([]*Constr, error) {
	if model == nil {
		return nil, errors.New("")
	}

	if len(vars) != len(val) || len(val) != len(sense) || len(sense) != len(rhs) || len(rhs) != len(constrname) {
		return nil, errors.New("")
	}

	numnz := 0
	for _, v := range vars {
		numnz += len(v)
	}

	beg := make([]int32, len(constrname))
	ind := make([]int32, numnz)
	_val := make([]float64, numnz)
	k := 0
	for i := 0; i < len(vars); i++ {
		if len(vars[i]) != len(val[i]) {
			return nil, errors.New("")
		}

		for j := 0; j < len(vars[i]); j++ {
			idx := vars[i][j].Index
			if idx < 0 {
				return nil, errors.New("")
			}
			ind[k+j] = idx
			_val[k+j] = val[i][j]
		}

		beg[i] = int32(k)
		k += len(vars[i])
	}

	name := make([](*C.char), len(constrname))
	for i, n := range constrname {
		name[i] = C.CString(n)
	}

	pbeg := (*C.int)(nil)
	pind := (*C.int)(nil)
	pval := (*C.double)(nil)
	if len(beg) > 0 {
		pbeg = (*C.int)(&beg[0])
		pind = (*C.int)(&ind[0])
		pval = (*C.double)(&_val[0])
	}

	psense := (*C.char)(nil)
	prhs := (*C.double)(nil)
	pname := (**C.char)(nil)
	if len(constrname) > 0 {
		psense = (*C.char)(&sense[0])
		prhs = (*C.double)(&rhs[0])
		pname = (**C.char)(&name[0])
	}

	err := C.GRBaddconstrs(model.AsGRBModel, C.int(len(constrname)), C.int(numnz), pbeg, pind, pval, psense, prhs, pname)
	if err != 0 {
		return nil, model.MakeError(err)
	}

	if err := model.Update(); err != nil {
		return nil, err
	}

	constrs := make([]*Constr, len(constrname))
	xrows := len(model.Constraints)
	for i := xrows; i < xrows+len(constrname); i++ {
		model.Constraints = append(model.Constraints, Constr{model, int32(i)})
		constrs[i] = &model.Constraints[len(model.Constraints)-1]
	}
	return constrs, nil
}

// SetObjective ...
func (model *Model) SetObjective(objectiveExpr interface{}, sense int32) error {

	// Clear Out All Previous Quadratic Objective Terms
	if err := C.GRBdelq(model.AsGRBModel); err != 0 {
		return model.MakeError(err)
	}

	// Detect the Type of Objective We Have
	switch objectiveExpr.(type) {
	case *LinExpr:
		le := objectiveExpr.(*LinExpr)
		if err := model.SetLinearObjective(le, sense); err != nil {
			return err
		}
	case *QuadExpr:
		qe := objectiveExpr.(*QuadExpr)
		if err := model.SetQuadraticObjective(qe, sense); err != nil {
			return err
		}
	default:
		return errors.New("Unexpected objective expression type!")
	}

	return nil
}

/*
SetLinearObjective
Description:

	Adds a linear objective to the model.
*/
func (model *Model) SetLinearObjective(expr *LinExpr, sense int32) error {
	// Constants

	// Algorithm
	for tempIndex, tempVar := range expr.ind {
		// Add Each to the objective, by modifying the obj attribute
		if err := tempVar.SetObj(expr.val[tempIndex]); err != nil {
			return err
		}
	}

	// if err := model.SetDoubleAttrVars(C.GRB_DBL_ATTR_OBJ, expr.lind, expr.lval); err != nil {
	// 	return err
	// }
	if err := model.SetDoubleAttr(C.GRB_DBL_ATTR_OBJCON, expr.offset); err != nil {
		return err
	}
	if err := model.SetIntAttr(C.GRB_INT_ATTR_MODELSENSE, sense); err != nil {
		return err
	}

	// If you successfully complete all steps, then return no errors.
	return nil
}

/*
SetQuadraticObjective
Description:

	Adds a quadratic objective to the model.
*/
func (model *Model) SetQuadraticObjective(expr *QuadExpr, sense int32) error {
	// Constants

	// Algorithm
	if err := model.addQPTerms(expr.qrow, expr.qcol, expr.qval); err != nil {
		return err
	}

	if err := model.SetDoubleAttrVars(C.GRB_DBL_ATTR_OBJ, expr.lind, expr.lval); err != nil {
		return err
	}
	if err := model.SetDoubleAttr(C.GRB_DBL_ATTR_OBJCON, expr.offset); err != nil {
		return err
	}
	if err := model.SetIntAttr(C.GRB_INT_ATTR_MODELSENSE, sense); err != nil {
		return err
	}

	// If you successfully complete all steps, then return no errors.
	return nil
}

func (model *Model) addQPTerms(qrow []*Var, qcol []*Var, qval []float64) error {
	if model == nil {
		return errors.New("")
	}

	if len(qrow) != len(qcol) || len(qcol) != len(qval) {
		return errors.New("")
	}

	_qrow := make([]int32, len(qrow))
	_qcol := make([]int32, len(qcol))
	for i := 0; i < len(qrow); i++ {
		if qrow[i].Index < 0 {
			return errors.New("")
		}
		if qcol[i].Index < 0 {
			return errors.New("")
		}

		_qrow[i] = qrow[i].Index
		_qcol[i] = qcol[i].Index
	}

	pqrow := (*C.int)(nil)
	pqcol := (*C.int)(nil)
	pqval := (*C.double)(nil)
	if len(qrow) > 0 {
		pqrow = (*C.int)(&_qrow[0])
		pqcol = (*C.int)(&_qcol[0])
		pqval = (*C.double)(&qval[0])
	}

	err := C.GRBaddqpterms(model.AsGRBModel, C.int(len(qrow)), pqrow, pqcol, pqval)
	if err != 0 {
		return model.MakeError(err)
	}

	return nil
}

// Update ...
func (model *Model) Update() error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBupdatemodel(model.AsGRBModel)
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

// Optimize ...
func (model *Model) Optimize() error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBoptimize(model.AsGRBModel)
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

// Write ...
func (model *Model) Write(filename string) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBwrite(model.AsGRBModel, C.CString(filename))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

// GetIntAttr ...
func (model *Model) GetIntAttr(attrname string) (int32, error) {
	if model == nil {
		return 0, errors.New("")
	}
	var attr int32
	err := C.GRBgetintattr(model.AsGRBModel, C.CString(attrname), (*C.int)(&attr))
	if err != 0 {
		return 0, model.MakeError(err)
	}
	return attr, nil
}

// GetDoubleAttr ...
func (model *Model) GetDoubleAttr(attrname string) (float64, error) {
	if model == nil {
		return 0, errors.New("")
	}
	var attr float64
	err := C.GRBgetdblattr(model.AsGRBModel, C.CString(attrname), (*C.double)(&attr))
	if err != 0 {
		return 0, model.MakeError(err)
	}
	return attr, nil
}

// GetStringAttr ...
func (model *Model) GetStringAttr(attrname string) (string, error) {
	if model == nil {
		return "", errors.New("")
	}
	var attr *C.char
	err := C.GRBgetstrattr(model.AsGRBModel, C.CString(attrname), (**C.char)(&attr))
	if err != 0 {
		return "", model.MakeError(err)
	}
	return C.GoString(attr), nil
}

// SetIntAttr ...
func (model *Model) SetIntAttr(attrname string, value int32) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBsetintattr(model.AsGRBModel, C.CString(attrname), C.int(value))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

// SetDoubleAttr ...
func (model *Model) SetDoubleAttr(attrname string, value float64) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBsetdblattr(model.AsGRBModel, C.CString(attrname), C.double(value))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

// SetStringAttr ...
func (model *Model) SetStringAttr(attrname string, value string) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBsetstrattr(model.AsGRBModel, C.CString(attrname), C.CString(value))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

// GetDoubleAttrVars ...
func (model *Model) GetDoubleAttrVars(attrname string, vars []*Var) ([]float64, error) {
	ind := make([]int32, len(vars))
	for i, v := range vars {
		if v.Index < 0 {
			return []float64{}, errors.New("")
		}
		ind[i] = v.Index
	}
	return model.getDoubleAttrList(attrname, ind)
}

// SetDoubleAttrVars ...
func (model *Model) SetDoubleAttrVars(attrname string, vars []*Var, value []float64) error {
	ind := make([]int32, len(vars))
	for i, v := range vars {
		if v.Index < 0 {
			return errors.New("")
		}
		ind[i] = v.Index
	}
	return model.setDoubleAttrList(attrname, ind, value)
}

func (model *Model) getIntAttrElement(attr string, ind int32) (int32, error) {
	if model == nil {
		return 0.0, model.MakeUninitializedError()
	}
	var value int32
	err := C.GRBgetintattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), (*C.int)(&value))
	if err != 0 {
		return 0, model.MakeError(err)
	}
	return value, nil
}

func (model *Model) getCharAttrElement(attr string, ind int32) (int8, error) {
	if model == nil {
		return 0, errors.New("")
	}
	var value int8
	err := C.GRBgetcharattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), (*C.char)(&value))
	if err != 0 {
		return 0, model.MakeError(err)
	}
	return value, nil
}

func (model *Model) getDoubleAttrElement(attr string, ind int32) (float64, error) {
	if model == nil {
		return 0, errors.New("")
	}
	var value float64
	err := C.GRBgetdblattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), (*C.double)(&value))
	if err != 0 {
		return 0, model.MakeError(err)
	}
	return value, nil
}

func (model *Model) getStringAttrElement(attr string, ind int32) (string, error) {
	if model == nil {
		return "", errors.New("")
	}
	var value *C.char
	err := C.GRBgetstrattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), (**C.char)(&value))
	if err != 0 {
		return "", model.MakeError(err)
	}
	return C.GoString(value), nil
}

func (model *Model) setIntAttrElement(attr string, ind int32, value int32) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBsetintattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), C.int(value))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

func (model *Model) setCharAttrElement(attr string, ind int32, value int8) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBsetcharattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), C.char(value))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

func (model *Model) setDoubleAttrElement(attr string, ind int32, value float64) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBsetdblattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), C.double(value))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

func (model *Model) setStringAttrElement(attr string, ind int32, value string) error {
	if model == nil {
		return errors.New("")
	}
	err := C.GRBsetstrattrelement(model.AsGRBModel, C.CString(attr), C.int(ind), C.CString(value))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

func (model *Model) getDoubleAttrList(attrname string, ind []int32) ([]float64, error) {
	if model == nil {
		return []float64{}, errors.New("")
	}
	if len(ind) == 0 {
		return []float64{}, nil
	}
	value := make([]float64, len(ind))
	err := C.GRBgetdblattrlist(model.AsGRBModel, C.CString(attrname), C.int(len(ind)), (*C.int)(&ind[0]), (*C.double)(&value[0]))
	if err != 0 {
		return []float64{}, model.MakeError(err)
	}
	return value, nil
}

func (model *Model) setDoubleAttrList(attrname string, ind []int32, value []float64) error {
	if model == nil {
		return errors.New("")
	}
	if len(ind) != len(value) {
		return errors.New("")
	}
	if len(ind) == 0 {
		return nil
	}
	err := C.GRBsetdblattrlist(model.AsGRBModel, C.CString(attrname), C.int(len(ind)), (*C.int)(&ind[0]), (*C.double)(&value[0]))
	if err != 0 {
		return model.MakeError(err)
	}
	return nil
}

/*
GetVarByName
Description:
	Collects the GRBVar which has the given gurobi variable by name.
*/
