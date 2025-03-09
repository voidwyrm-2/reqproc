package numbertype

import (
	"cmp"
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

var (
	ValueTrue  = New(1)
	ValueFalse = New(0)
)

type ReqNumberType struct {
	basetype.ReqBaseType
	value float32
}

func New(value float32) ReqNumberType {
	return ReqNumberType{value: value, ReqBaseType: basetype.New(types.TypeString)}
}

func (rst ReqNumberType) String() string {
	return fmt.Sprint(rst.value)
}

func (rst ReqNumberType) Literal() any {
	return rst.value
}

func (rst ReqNumberType) Add(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeString {
		return rst.ReqBaseType.Add(other)
	}

	return New(rst.value + other.Literal().(float32)), nil
}

func (rst ReqNumberType) Sub(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeString {
		return rst.ReqBaseType.Sub(other)
	}

	return New(rst.value - other.Literal().(float32)), nil
}

func (rst ReqNumberType) Mul(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeString {
		return rst.ReqBaseType.Mul(other)
	}

	return New(rst.value * other.Literal().(float32)), nil
}

func (rst ReqNumberType) Div(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeString {
		return rst.ReqBaseType.Div(other)
	}

	return New(rst.value / other.Literal().(float32)), nil
}

func (rst ReqNumberType) Not() (types.ReqType, error) {
	if rst.value == 0 {
		return ValueFalse, nil
	}

	return ValueTrue, nil
}

func (rst ReqNumberType) Cmp(other types.ReqType) (bool, int) {
	if other.Type() != types.TypeString {
		return rst.ReqBaseType.Cmp(other)
	}

	return rst.value == other.Literal().(float32), cmp.Compare(rst.value, other.Literal().(float32))
}
