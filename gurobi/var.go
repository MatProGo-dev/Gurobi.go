package gurobi

/*
var.go
Description:
	A set of functions for manipulating the gurobi Var object.
Notes:
	The available attributes for variables is listed on Gurobi's website at links like:
	https://www.gurobi.com/documentation/current/refman/attributes.html
*/
// Gurobi variable object
type Var struct {
	Model *Model
	Index int32
}

func (v *Var) GetInt(attr string) (int32, error) {
	return v.Model.getIntAttrElement(attr, v.Index)
}

func (v *Var) GetChar(attr string) (int8, error) {
	return v.Model.getCharAttrElement(attr, v.Index)
}

func (v *Var) GetDouble(attr string) (float64, error) {
	return v.Model.getDoubleAttrElement(attr, v.Index)
}

func (v *Var) GetString(attr string) (string, error) {
	return v.Model.getStringAttrElement(attr, v.Index)
}

func (v *Var) SetInt(attr string, value int32) error {
	return v.Model.setIntAttrElement(attr, v.Index, value)
}

func (v *Var) SetChar(attr string, value int8) error {
	return v.Model.setCharAttrElement(attr, v.Index, value)
}

func (v *Var) SetDouble(attr string, value float64) error {
	return v.Model.setDoubleAttrElement(attr, v.Index, value)
}

func (v *Var) SetString(attr string, value string) error {
	return v.Model.setStringAttrElement(attr, v.Index, value)
}

func (v *Var) SetObj(value float64) error {
	err := v.Model.setDoubleAttrElement("Obj", v.Index, value)
	if err != nil {
		return err
	}

	// Update model and return
	return v.Model.Update()
}
