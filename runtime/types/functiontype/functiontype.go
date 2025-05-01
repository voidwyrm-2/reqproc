package functiontype

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/voidwyrm-2/reqproc/lexer/tokens"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

// callf is needed so native functions can call functions popped off the stack
type NativeFunction = func(sc *scope.Scope, st *stack.Stack, callf func(rft ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error

// A runnable function value; can be either written natively in Go, or written in ReqProc
type ReqFunctionType struct {
	basetype.ReqBaseType
	doc           string
	native        NativeFunction
	tokens        []tokens.Token
	signature     float32
	input, output []types.ReqVarType
}

/*
func parseFuncParams(fn reflect.Type) []types.ReqVarType {
	parsed := []types.ReqVarType{}

	tokens := tokenizeSigniture(fn.String())

	for _, t := range tokens {
		switch t {
		case "string":
			parsed = append(parsed, types.TypeString)
		case "func":
			parsed = append(parsed, types.TypeFunction)
		default:
			if strings.HasPrefix(t, "int") || strings.HasPrefix(t, "uint") || strings.HasPrefix(t, "float") {
				parsed = append(parsed, types.TypeNumber)
			} else if strings.HasPrefix(t, "[]") {
				parsed = append(parsed, types.TypeList)
			} else {
				panic(fmt.Sprintf("invalid parameter type '%s'", t))
			}
		}
	}

	return parsed
}
*/

func parseSignature(signature float32) (uint32, uint32) {
	s := fmt.Sprint(signature)
	if !strings.Contains(s, ".") {
		s += ".0"
	}

	if n, err := strconv.ParseUint(strings.Split(s, ".")[1], 10, 32); err != nil {
		panic(err.Error())
	} else {
		return uint32(signature), uint32(n)
	}
}

func makeTypeSlice(length uint32, t types.ReqVarType) []types.ReqVarType {
	sl := make([]types.ReqVarType, 0, length)

	for range length {
		sl = append(sl, t)
	}

	return sl
}

func makeSignature(i, o int) float32 {
	sig := float32(o)
	for sig > 1 {
		sig *= 0.1
	}

	return float32(i) + sig
}

func NewNative(fn func(sc *scope.Scope, st *stack.Stack, callf func(rft ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error, input, output []types.ReqVarType) ReqFunctionType {
	return ReqFunctionType{
		native:      fn,
		signature:   makeSignature(len(input), len(output)),
		input:       input,
		output:      output,
		ReqBaseType: basetype.New(types.TypeFunction),
	}
}

func NewSigNative(fn func(sc *scope.Scope, st *stack.Stack, callf func(rft ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error, signature float32) ReqFunctionType {
	input, output := parseSignature(signature)

	return NewNative(fn, makeTypeSlice(input, types.TypeAny), makeTypeSlice(output, types.TypeAny))
}

func New(tokens []tokens.Token, signature float32) ReqFunctionType {
	input, output := parseSignature(signature)

	return ReqFunctionType{
		tokens:      tokens,
		signature:   signature,
		input:       makeTypeSlice(input, types.TypeAny),
		output:      makeTypeSlice(output, types.TypeAny),
		ReqBaseType: basetype.New(types.TypeFunction),
	}
}

func (rft ReqFunctionType) SetDoc(doc string) types.ReqType {
	rft.doc = doc
	return rft
}

func (rft ReqFunctionType) Input() []types.ReqVarType {
	return rft.input
}

func (rft ReqFunctionType) Output() []types.ReqVarType {
	return rft.output
}

func (rft ReqFunctionType) Signature() float32 {
	return rft.signature
}

func (rft ReqFunctionType) ExpectSignature(expected float32) error {
	if rft.signature != expected {
		return fmt.Errorf("expected signature %f but found %f instead", expected, rft.signature)
	}

	return nil
}

func (rft ReqFunctionType) Doc() string {
	return rft.doc
}

func (rft ReqFunctionType) Literal() any {
	if rft.native == nil {
		return rft.tokens
	}

	return rft.native
}

/*
func (rft ReqFunctionType) String() string {
	formatted := []string{}

	for _, a := range rft.args {
		formatted = append(formatted, a.String())
	}

	return "function{" + strings.Join(formatted, ", ") + "}"
}
*/

func (rft ReqFunctionType) String() string {
	input, output := make([]string, 0, len(rft.input)), make([]string, 0, len(rft.output))

	for _, t := range rft.input {
		input = append(input, t.String())
	}

	for _, t := range rft.output {
		output = append(output, t.String())
	}

	return fmt.Sprintf("function{%s -> %s}", strings.Join(input, ", "), strings.Join(output, ", "))
}
