package basetype

import (
	. "github.com/voidwyrm-2/reqproc/runtime/types"
)

type ReqBaseType struct {
	kind string
}

func New(kind string) ReqBaseType {
	return ReqBaseType{kind: kind}
}

func (rbt ReqBaseType) Type() string {
	return rbt.kind
}

func (rbt ReqBaseType) String() string {
	return "<base>"
}

func (rbt ReqBaseType) Literal() any {
	return nil
}

func (rbt ReqBaseType) Add(other ReqType) (ReqType, error) {
	return nil, InvalidOperation("addition", rbt, other)
}

func (rbt ReqBaseType) Sub(other ReqType) (ReqType, error) {
	return nil, InvalidOperation("subtraction", rbt, other)
}

func (rbt ReqBaseType) Mul(other ReqType) (ReqType, error) {
	return nil, InvalidOperation("multiplication", rbt, other)
}

func (rbt ReqBaseType) Div(other ReqType) (ReqType, error) {
	return nil, InvalidOperation("division", rbt, other)
}

func (rbt ReqBaseType) Not() (ReqType, error) {
	return nil, InvalidSingleOperation("not", rbt)
}

func (rbt ReqBaseType) Cmp(other ReqType) (bool, int) {
	return false, 0
}

func (rbt ReqBaseType) Length() (int, error) {
	return 0, InvalidSingleOperation("length", rbt)
}
