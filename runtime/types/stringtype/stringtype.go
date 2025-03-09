package stringtype

import (
	"cmp"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

type ReqStringType struct {
	basetype.ReqBaseType
	value string
}

func New(value string) ReqStringType {
	return ReqStringType{value: value, ReqBaseType: basetype.New(types.TypeString)}
}

func (rst ReqStringType) String() string {
	return "`" + rst.value + "`"
}

func (rst ReqStringType) Literal() any {
	return rst.value
}

func (rst ReqStringType) Add(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeString {
		return rst.ReqBaseType.Add(other)
	}

	return New(rst.value + other.Literal().(string)), nil
}

func (rst ReqStringType) Cmp(other types.ReqType) (bool, int) {
	if other.Type() != types.TypeString {
		return rst.ReqBaseType.Cmp(other)
	}

	return rst.value == other.Literal().(string), cmp.Compare(rst.value, other.Literal().(string))
}

func (rst ReqStringType) Length() (int, error) {
	return len(rst.value), nil
}
