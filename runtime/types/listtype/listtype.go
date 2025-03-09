package listtype

import (
	"strings"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
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
