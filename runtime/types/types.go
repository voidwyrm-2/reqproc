package types

import "fmt"

const (
	TypeString   = "string"
	TypeNumber   = "number"
	TypeList     = "list"
	TypeFunction = "function"
)

type ReqType interface {
	Type() string
	String() string
	Literal() any
	Add(other ReqType) (ReqType, error)
	Sub(other ReqType) (ReqType, error)
	Mul(other ReqType) (ReqType, error)
	Div(other ReqType) (ReqType, error)
	Not() (ReqType, error)
	Cmp(other ReqType) (bool, int) // equality, less/greater
	Length() (int, error)
}

func InvalidOperation(operation string, typeA, typeB ReqType) error {
	return fmt.Errorf("invalid operation '%s' for types '%s' and '%s'", operation, typeA.Type(), typeB.Type())
}

func InvalidSingleOperation(operation string, typeA ReqType) error {
	return fmt.Errorf("invalid operation '%s' for types '%s'", operation, typeA.Type())
}
