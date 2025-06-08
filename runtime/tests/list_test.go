package test

import (
	"testing"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
	"github.com/voidwyrm-2/reqproc/runtime/types/stringtype"
)

func TestLists(t *testing.T) {
	cases := []stackTestCase{
		{
			`[0 1 2 3 4]`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(0), numbertype.New(1), numbertype.New(2), numbertype.New(3), numbertype.New(4)}},
			},
			true,
		},
		{
			`[30 "hello" -5 "wow" 1]`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(30), stringtype.New("hello"), numbertype.New(-5), stringtype.New("wow"), numbertype.New(1)}},
			},
			true,
		},
		{
			"[0 1 2 3 4] 1 +",
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(1), numbertype.New(2), numbertype.New(3), numbertype.New(4), numbertype.New(5)}},
			},
			true,
		},
		{
			"[0 1 2 3 4] 20 +",
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(20), numbertype.New(21), numbertype.New(22), numbertype.New(23), numbertype.New(24)}},
			},
			true,
		},
		{
			"[0 1 2 3 4] 3 -",
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(-3), numbertype.New(-2), numbertype.New(-1), numbertype.New(0), numbertype.New(1)}},
			},
			true,
		},
		{
			"[2 3 1 0 4] 4 *",
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(8), numbertype.New(12), numbertype.New(4), numbertype.New(0), numbertype.New(16)}},
			},
			true,
		},
	}

	testStack(t, cases)
}
