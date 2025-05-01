package nativetype

import (
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

type ReqNativeType struct {
	basetype.ReqBaseType
	value any
}

func New(value any) ReqNativeType {
	return ReqNativeType{value: value, ReqBaseType: basetype.New(types.TypeNative)}
}

func (rnt ReqNativeType) String() string {
	return "`" + fmt.Sprint(rnt.value) + "`"
}

func (rnt ReqNativeType) Literal() any {
	return rnt.value
}
