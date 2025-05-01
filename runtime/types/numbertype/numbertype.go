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
	return ReqNumberType{value: value, ReqBaseType: basetype.New(types.TypeNumber)}
}

func (rnt ReqNumberType) String() string {
	return fmt.Sprint(rnt.value)
}

func (rnt ReqNumberType) Literal() any {
	return rnt.value
}

func (rnt ReqNumberType) Add(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeNumber {
		return rnt.ReqBaseType.Add(other)
	}

	return New(rnt.value + other.Literal().(float32)), nil
}

func (rnt ReqNumberType) Sub(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeNumber {
		return rnt.ReqBaseType.Sub(other)
	}

	return New(rnt.value - other.Literal().(float32)), nil
}

func (rnt ReqNumberType) Mul(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeNumber {
		return rnt.ReqBaseType.Mul(other)
	}

	return New(rnt.value * other.Literal().(float32)), nil
}

func (rnt ReqNumberType) Div(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeNumber {
		return rnt.ReqBaseType.Div(other)
	}

	return New(rnt.value / other.Literal().(float32)), nil
}

func (rnt ReqNumberType) Not() (types.ReqType, error) {
	if rnt.value == 0 {
		return ValueFalse, nil
	}

	return ValueTrue, nil
}

func (rnt ReqNumberType) Cmp(other types.ReqType) (bool, int) {
	if other.Type() != types.TypeNumber {
		return rnt.ReqBaseType.Cmp(other)
	}

	return rnt.value == other.Literal().(float32), cmp.Compare(rnt.value, other.Literal().(float32))
}

func (rnt ReqNumberType) IsFloat() bool {
	return float32(int32(rnt.value)) != rnt.value
}
