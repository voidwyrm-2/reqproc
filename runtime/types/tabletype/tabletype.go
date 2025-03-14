package tabletype

import (
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"

	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

type TableType struct {
	basetype.ReqBaseType
	value map[string]types.ReqType
}

func New(value map[string]types.ReqType) TableType {
	return TableType{value: value, ReqBaseType: basetype.New(types.TypeTable)}
}

func (tbt TableType) String() string {
	return fmt.Sprint(tbt.value)
}

func (tbt TableType) Literal() any {
	return tbt.value
}

func (tbt TableType) Length() (int, error) {
	return len(tbt.value), nil
}
