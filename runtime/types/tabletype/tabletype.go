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
	return fmt.Sprint(tbt.value)
}

func (tbt ReqTableType) Literal() any {
	return tbt.value
}

func (tbt ReqTableType) Length() (int, error) {
	return len(tbt.value), nil
}
