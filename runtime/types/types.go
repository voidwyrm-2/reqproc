package types

import (
	"fmt"
)

type ReqVarType int8

const (
	TypeBase              = -1
	TypeString ReqVarType = 1 << (iota - 1)
	TypeNumber
	TypeList
	TypeTable
	TypeFunction
	TypeNative
)

const TypeAny = TypeString | TypeNumber | TypeList | TypeTable | TypeFunction | TypeNative

var typeNameMapFrom, typeNameMapInto = func() (map[ReqVarType]string, map[string]ReqVarType) {
	a := map[ReqVarType]string{
		TypeAny:      "any",
		TypeString:   "string",
		TypeNumber:   "number",
		TypeList:     "list",
		TypeTable:    "table",
		TypeFunction: "function",
		TypeNative:   "native",
	}

	b := map[string]ReqVarType{}
	for k, v := range a {
		b[v] = k
	}

	return a, b
}()

func TypeFromString(s string) (ReqVarType, error) {
	if v, ok := typeNameMapInto[s]; ok {
		return v, nil
	}

	return TypeBase, fmt.Errorf("'%s' is not a valid ReqVarType", s)
}

func (rvt ReqVarType) String() string {
	if v, ok := typeNameMapFrom[rvt]; ok {
		return v
	}

	panic(fmt.Sprintf("invalid ReqVarType %d", rvt))
}

var IllegalVariableNames = map[string]struct{}{
	"true":   {},
	"false":  {},
	"import": {},
	"def":    {},
}

type ReqType interface {
	Type() ReqVarType
	String() string
	Literal() any
	Add(other ReqType) (ReqType, error)
	Sub(other ReqType) (ReqType, error)
	Mul(other ReqType) (ReqType, error)
	Div(other ReqType) (ReqType, error)
	Not() (ReqType, error)
	Cmp(other ReqType) (bool, int) // equality, less/greater
	Length() (int, error)
	GetIndex(index ReqType) (ReqType, error)
	SetIndex(index ReqType, value ReqType) error
}

func ExpectType(expected, actual ReqVarType) error {
	if expected != actual {
		return fmt.Errorf("expected type '%s' but found '%s' instead", expected.String(), actual.String())
	}

	return nil
}

func InvalidOperation(operation string, typeA, typeB ReqType) error {
	return fmt.Errorf("invalid operation '%s' for types '%s' and '%s'", operation, typeA.Type(), typeB.Type())
}

func InvalidSingleOperation(operation string, typeA ReqType) error {
	return fmt.Errorf("invalid operation '%s' for types '%s'", operation, typeA.Type())
}
