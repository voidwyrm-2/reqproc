package test

import (
	"testing"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
)

func TestUiuaFunctions(t *testing.T) {
	cases := []stackTestCase{
		{
			`4 1 2 @+ dip`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeNumber, float32(5)},
				{types.TypeNumber, float32(2)},
			},
			true,
		},
		{
			`4 10 @range dip`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(0), numbertype.New(1), numbertype.New(2), numbertype.New(3)}},
				{types.TypeNumber, float32(10)},
			},
			true,
		},
		{
			`5 range`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(0), numbertype.New(1), numbertype.New(2), numbertype.New(3), numbertype.New(4)}},
			},
			true,
		},
		{
			`10 range`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, []types.ReqType{numbertype.New(0), numbertype.New(1), numbertype.New(2), numbertype.New(3), numbertype.New(4), numbertype.New(5), numbertype.New(6), numbertype.New(7), numbertype.New(8), numbertype.New(9)}},
			},
			true,
		},
	}

	testStack(t, cases)
}

func TestOsFs(t *testing.T) {
	cases := []stackTestCase{
		{
			`"os" import @os.fs`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeTable, nil},
			},
			true,
		},
		{
			`"os" import "." os.fs.items`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeList, nil},
			},
			true,
		},
	}

	testStack(t, cases)
}

func TestFFI(t *testing.T) {
	cases := []stackTestCase{
		{
			`"ffi" import 20 ffi.toNative`,
			[]struct {
				t types.ReqVarType
				v any
			}{
				{types.TypeNative, nil},
			},
			true,
		},
	}

	testStack(t, cases)
}
