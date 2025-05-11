package basetype

import (
	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type ReqBaseType struct {
	kind types.ReqVarType
}

func New(kind types.ReqVarType) ReqBaseType {
	return ReqBaseType{kind: kind}
}

func (rbt ReqBaseType) Type() types.ReqVarType {
	return rbt.kind
}

func (rbt ReqBaseType) String() string {
	return "<base>"
}

func (rbt ReqBaseType) Literal() any {
	return nil
}

func (rbt ReqBaseType) Add(other types.ReqType) (types.ReqType, error) {
	return nil, types.InvalidOperation("addition", rbt, other)
}

func (rbt ReqBaseType) Sub(other types.ReqType) (types.ReqType, error) {
	return nil, types.InvalidOperation("subtraction", rbt, other)
}

func (rbt ReqBaseType) Mul(other types.ReqType) (types.ReqType, error) {
	return nil, types.InvalidOperation("multiplication", rbt, other)
}

func (rbt ReqBaseType) Div(other types.ReqType) (types.ReqType, error) {
	return nil, types.InvalidOperation("division", rbt, other)
}

func (rbt ReqBaseType) Not() (types.ReqType, error) {
	return nil, types.InvalidSingleOperation("not", rbt)
}

func (rbt ReqBaseType) Cmp(other types.ReqType) (bool, int) {
	return false, 0
}

func (rbt ReqBaseType) Length() (int, error) {
	return 0, types.InvalidSingleOperation("length", rbt)
}

func (rbt ReqBaseType) GetIndex(index types.ReqType) (types.ReqType, error) {
	return nil, types.InvalidSingleOperation("index get", rbt)
}

func (rbt ReqBaseType) SetIndex(index types.ReqType, value types.ReqType) error {
	return types.InvalidSingleOperation("index set", rbt)
}
