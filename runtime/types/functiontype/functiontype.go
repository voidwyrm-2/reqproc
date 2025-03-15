package functiontype

import (
	"fmt"

	"github.com/voidwyrm-2/reqproc/lexer/tokens"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

type ReqFunctionType struct {
	basetype.ReqBaseType
	// value reflect.Value
	// args  []types.ReqVarType
	native    func(sc *scope.Scope, st *stack.Stack) error
	tokens    []tokens.Token
	args, ret int
}

/*

func tokenizeSigniture(sig string) []string {
	tokens := []string{""}

	t := 0
	nest := 0
	for _, ch := range sig {
		if ch == '(' {
			nest++
		} else if ch == ')' {
			nest--
		} else if nest == 1 {
			if ch == ' ' || ch == ',' {
				if tokens[t] != "" {
					t++
					tokens = append(tokens, "")
				}
			} else {
				tokens[t] += string(ch)
			}
		}
	}

	return tokens
}

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

func NewNative(fn any) ReqFunctionType {
	value := reflect.ValueOf(fn)
	if value.Kind().String() != "func" {
		panic("value given for new ReqFunctionType is not a function")
	}

	return ReqFunctionType{value: value, args: parseFuncParams(value.Type()), ReqBaseType: basetype.New(types.TypeFunction)}
}
*/

func NewNative(args, ret int, fn func(sc *scope.Scope, st *stack.Stack) error) ReqFunctionType {
	return ReqFunctionType{native: fn, args: args, ret: ret, ReqBaseType: basetype.New(types.TypeFunction)}
}

func New(tokens []tokens.Token) ReqFunctionType {
	f := ReqFunctionType{tokens: tokens, ReqBaseType: basetype.New(types.TypeFunction)}
	f.parseForArgs()
	return f
}

func (rft ReqFunctionType) Args() int {
	return rft.args
}

func (rft ReqFunctionType) Ret() int {
	return rft.ret
}

func (rft *ReqFunctionType) parseForArgs() {
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
	return fmt.Sprintf("function{%d}", rft.args)
}
