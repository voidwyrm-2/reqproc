package types

import "fmt"

type ReqVarType int8

const (
	TypeBase              = -1
	TypeString ReqVarType = (iota + 1) << 1
	TypeNumber
	TypeList
	TypeTable
	TypeFunction
)

const TypeAny = TypeString | TypeNumber | TypeList | TypeTable | TypeFunction

func (rvt ReqVarType) String() string {
	switch rvt {
	case TypeAny:
		return "any"
	case TypeString:
		return "string"
	case TypeNumber:
		return "number"
	case TypeList:
		return "list"
	case TypeTable:
		return "table"
	case TypeFunction:
		return "function"
	default:
		panic(fmt.Sprintf("invalid ReqVarType %d", rvt))
	}
}

var IllegalVariableNames = []string{
	"true",
	"false",
	"import",
	"def",
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
}

func InvalidOperation(operation string, typeA, typeB ReqType) error {
	return fmt.Errorf("invalid operation '%s' for types '%s' and '%s'", operation, typeA.Type(), typeB.Type())
}

func InvalidSingleOperation(operation string, typeA ReqType) error {
	return fmt.Errorf("invalid operation '%s' for types '%s'", operation, typeA.Type())
}
