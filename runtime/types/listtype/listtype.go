package listtype

import (
	"errors"
	"fmt"
	"strings"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
)

type ReqListType struct {
	basetype.ReqBaseType
	value []types.ReqType
}

func New(value ...types.ReqType) ReqListType {
	return ReqListType{value: value, ReqBaseType: basetype.New(types.TypeList)}
}

func (rlt ReqListType) String() string {
	formatted := []string{}

	for _, v := range rlt.value {
		formatted = append(formatted, v.String())
	}

	return "[" + strings.Join(formatted, ", ") + "]"
}

func (rlt ReqListType) Literal() any {
	return rlt.value
}

func (rlt ReqListType) Add(other types.ReqType) (types.ReqType, error) {
	if other.Type() == types.TypeList {
		v := other.Literal().([]types.ReqType)
		return New(append(rlt.value, v...)...), nil
	}

	for i := range rlt.value {
		v, err := rlt.value[i].Add(other)
		if err != nil {
			return nil, err
		}

		rlt.value[i] = v
	}

	return rlt, nil
}

func (rlt ReqListType) Sub(other types.ReqType) (types.ReqType, error) {
	if other.Type() == types.TypeList {
		return rlt.ReqBaseType.Sub(other)
	}

	for i := range rlt.value {
		v, err := rlt.value[i].Sub(other)
		if err != nil {
			return nil, err
		}

		rlt.value[i] = v
	}

	return rlt, nil
}

func (rlt ReqListType) Mul(other types.ReqType) (types.ReqType, error) {
	if other.Type() == types.TypeList {
		return rlt.ReqBaseType.Mul(other)
	}

	for i := range rlt.value {
		v, err := rlt.value[i].Mul(other)
		if err != nil {
			return nil, err
		}

		rlt.value[i] = v
	}

	return rlt, nil
}

func (rlt ReqListType) Div(other types.ReqType) (types.ReqType, error) {
	if other.Type() == types.TypeList {
		return rlt.ReqBaseType.Div(other)
	}

	for i := range rlt.value {
		v, err := rlt.value[i].Div(other)
		if err != nil {
			return nil, err
		}

		rlt.value[i] = v
	}

	return rlt, nil
}

func (rlt ReqListType) Length() (int, error) {
	return len(rlt.value), nil
}

func (rlt ReqListType) GetIndex(index types.ReqType) (types.ReqType, error) {
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

	return rlt.value[n], nil
}

func (rlt ReqListType) SetIndex(index types.ReqType, value types.ReqType) error {
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
