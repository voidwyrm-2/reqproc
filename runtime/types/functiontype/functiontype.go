package functiontype

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/basetype"
)

type ReqFunctionType struct {
	basetype.ReqBaseType
	value reflect.Value
	args  []string
}

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

func parseFuncParams(fn reflect.Type) []string {
	parsed := []string{}

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

func (rft ReqFunctionType) Literal() any {
	return rft.value
}

func (rft ReqFunctionType) String() string {
	return "function{" + strings.Join(rft.args, ", ") + "}"
}
