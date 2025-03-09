package scope

import (
	"fmt"

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

func (sc *Scope) Read(name string) (types.ReqType, error) {
	if v, ok := sc.vars[name]; ok {
		return v, nil
	} else if v, ok = sc.consts[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("variable/constant '%s' does not exist", name)
}

func (sc *Scope) Write(name string, value types.ReqType) error {
	if _, ok := sc.vars[name]; ok {
		return fmt.Errorf("variable '%s' already exists", name)
	} else if _, ok := sc.consts[name]; ok {
		return fmt.Errorf("'%s' already exists as a constant", name)
	}

	sc.vars[name] = value

	return nil
}

func (sc *Scope) WriteConst(name string, value types.ReqType) error {
	if _, ok := sc.consts[name]; ok {
		return fmt.Errorf("constant '%s' already exists", name)
	} else if _, ok := sc.vars[name]; ok {
		return fmt.Errorf("'%s' already exists as a variable", name)
	}

	sc.consts[name] = value

	return nil
}

func (sc *Scope) Update(name string, value types.ReqType) error {
	_, err := sc.Read(name)
	if err == nil {
		if _, ok := sc.consts[name]; ok {
			return fmt.Errorf("cannot reassign constant '%s'", name)
		}

		sc.vars[name] = value
	} else if sc.parent != nil {
		return sc.parent.Update(name, value)
	}

	return err
}
