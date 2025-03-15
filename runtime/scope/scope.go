package scope

import (
	"fmt"
	"slices"
	"strings"

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
	if strings.Contains(name, ".") {
		path := strings.Split(name, ".")

		v, err := sc.Read(path[0])
		if err != nil {
			return nil, err
		}

		return sc.nestedRead(path[1:], v)
	}

	if v, ok := sc.vars[name]; ok {
		if v == nil {
			return nil, fmt.Errorf("variable '%s' has not had a value assigned to it yet", name)
		}

		return v, nil
	} else if v, ok = sc.consts[name]; ok {
		return v, nil
	}

	if sc.parent != nil {
		return sc.parent.Read(name)
	}

	return nil, fmt.Errorf("variable/constant '%s' does not exist", name)
}

func (sc Scope) nestedRead(path []string, tbl types.ReqType) (types.ReqType, error) {
	if path[0] == "" {
		return nil, fmt.Errorf("the dot indexed path cannot be empty")
	} else if tbl.Type() != types.TypeTable {
		return nil, fmt.Errorf("'%s' is not a dot indexable type", tbl.Type().String())
	}

	m := tbl.Literal().(map[string]types.ReqType)

	v, ok := m[path[0]]
	if !ok {
		return nil, fmt.Errorf("key '%s' does not exist", path[0])
	}

	if len(path) == 1 {
		return v, nil
	}

	return sc.nestedRead(path[1:], v)
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
