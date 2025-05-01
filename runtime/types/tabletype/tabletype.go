package tabletype

import (
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"

	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

type ReqTableType struct {
	basetype.ReqBaseType
	value map[string]types.ReqType
}

func New(value map[string]types.ReqType) ReqTableType {
	return ReqTableType{value: value, ReqBaseType: basetype.New(types.TypeTable)}
}

func (tbt ReqTableType) String() string {
	/*
		if m, ok := tbt.value["__string"]; ok {
			if m.Type() == types.TypeFunction {
			} else if m.Type() == types.TypeString {
				return m.Literal().(string)
			}
		}
	*/

	return fmt.Sprint(tbt.value)
}

func (tbt ReqTableType) Literal() any {
	return tbt.value
}

func (tbt ReqTableType) Add(other types.ReqType) (types.ReqType, error) {
	return tbt.ReqBaseType.Add(other)
}

func (tbt ReqTableType) Sub(other types.ReqType) (types.ReqType, error) {
	return tbt.ReqBaseType.Sub(other)
}

func (tbt ReqTableType) Mul(other types.ReqType) (types.ReqType, error) {
	return tbt.ReqBaseType.Mul(other)
}

func (tbt ReqTableType) Div(other types.ReqType) (types.ReqType, error) {
	return tbt.ReqBaseType.Div(other)
}

func (tbt ReqTableType) Cmp(other types.ReqType) (bool, int) {
	return tbt.ReqBaseType.Cmp(other)
}

func (tbt ReqTableType) Length() (int, error) {
	return len(tbt.value), nil
}
