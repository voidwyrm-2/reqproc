package scope

import (
	"fmt"
	"slices"

	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type Scope struct {
	vars   map[string]types.ReqType
	consts map[string]types.ReqType
	parent *Scope
}

func New(parent *Scope) *Scope {
	return &Scope{vars: map[string]types.ReqType{}, consts: map[string]types.ReqType{}, parent: parent}
}

func (sc Scope) Vars() map[string]types.ReqType {
	return sc.vars
}

func (sc Scope) Consts() map[string]types.ReqType {
	return sc.consts
}

func (sc Scope) Read(name string) (types.ReqType, error) {
	if v, ok := sc.vars[name]; ok {
		if v == nil {
			return nil, fmt.Errorf("variable '%s' has not had a value assigned to it yet", name)
		}

		return v, nil
	} else if v, ok = sc.consts[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("variable/constant '%s' does not exist", name)
}

func (sc *Scope) Write(name string, value types.ReqType) error {
	if slices.Contains(types.IllegalVariableNames, name) {
		return fmt.Errorf("'%s' is not a valid variable name", name)
	} else if _, ok := sc.vars[name]; ok {
		return fmt.Errorf("variable '%s' already exists", name)
	} else if _, ok := sc.consts[name]; ok {
		return fmt.Errorf("'%s' already exists as a constant", name)
	}

	sc.vars[name] = value

	return nil
}

func (sc *Scope) WriteConst(name string, value types.ReqType) error {
	if slices.Contains(types.IllegalVariableNames, name) {
		return fmt.Errorf("'%s' is not a valid variable name", name)
	} else if _, ok := sc.consts[name]; ok {
		return fmt.Errorf("constant '%s' already exists", name)
	} else if _, ok := sc.vars[name]; ok {
		return fmt.Errorf("'%s' already exists as a variable", name)
	}

	sc.consts[name] = value

	return nil
}

func (sc *Scope) Update(name string, value types.ReqType) error {
	if _, ok := sc.consts[name]; ok {
		return fmt.Errorf("cannot reassign constant '%s'", name)
	}

	_, ok := sc.vars[name]

	if ok {
		sc.vars[name] = value
	} else if sc.parent != nil {
		return sc.parent.Update(name, value)
	} else {
		return fmt.Errorf("variable/constant '%s' does not exist", name)
	}

	return nil
}
