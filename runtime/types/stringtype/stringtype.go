package stringtype

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
)

type ReqStringType struct {
	basetype.ReqBaseType
	value string
}

func New(value string) ReqStringType {
	return ReqStringType{value: value, ReqBaseType: basetype.New(types.TypeString)}
}

func (rst ReqStringType) String() string {
	return rst.value
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

func (rst ReqStringType) Mul(other types.ReqType) (types.ReqType, error) {
	if other.Type() != types.TypeNumber {
		return rst.ReqBaseType.Add(other)
	}

	value := other.Literal().(float32)
	ivalue := int32(value)
	if float32(ivalue) != value {
		return nil, errors.New("cannot use float value as string multiplier")
	}

	s := ""

	for range ivalue {
		s += rst.value
	}

	return New(s), nil
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

func (rlt ReqStringType) GetIndex(index types.ReqType) (types.ReqType, error) {
	if index.Type() != types.TypeNumber {
		return rlt.ReqBaseType.GetIndex(index)
	}

	nt := index.(numbertype.ReqNumberType)
	if nt.IsFloat() {
		return nil, errors.New("cannot use float value as index")
	}

	n := int(nt.Literal().(float32))
	if n >= len(rlt.value) {
		return nil, fmt.Errorf("index %d out of range for length %d", n, len(rlt.value))
	}

	return New(string(rlt.value[n])), nil
}

/*
func (rlt ReqStringType) SetIndex(index types.ReqType, value types.ReqType) error {
	if index.Type() != types.TypeNumber {
		return rlt.ReqBaseType.SetIndex(index, value)
	}

	nt := index.(numbertype.ReqNumberType)
	if nt.IsFloat() {
		return errors.New("cannot use float value as index")
	}

	n := int(nt.Literal().(float32))
	if n >= len(rlt.value) {
		return fmt.Errorf("index %d out of range for length %d", n, len(rlt.value))
	}

	rlt.value[n] = value

	return nil
}
*/
