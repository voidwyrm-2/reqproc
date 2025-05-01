package stdlib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/voidwyrm-2/reqproc/runtime"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/functiontype"
	"github.com/voidwyrm-2/reqproc/runtime/types/listtype"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
	"github.com/voidwyrm-2/reqproc/runtime/types/stringtype"
)

/*
Contains all the natively written functions

functions inside __init__ are loaded into the scope automatically at interpreter initialization
*/
var Stdlib = map[string]map[string]types.ReqType{
	"__init__": {
		"exit": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			return fmt.Errorf("EXIT CODE %d", int(st.Pop().Literal().(float32)))
		}, []types.ReqVarType{types.TypeNumber}, []types.ReqVarType{}),

		// meta
		"doc": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			f := st.Pop().(functiontype.ReqFunctionType)

			st.Push(stringtype.New(f.Doc()))

			return nil
		}, []types.ReqVarType{types.TypeFunction}, []types.ReqVarType{types.TypeString}).SetDoc("Returns the docstring of the function it's called on"),

		"sig": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			st.Push(numbertype.New(st.Pop().(functiontype.ReqFunctionType).Signature()))

			return nil
		}, []types.ReqVarType{types.TypeFunction}, []types.ReqVarType{types.TypeNumber}).SetDoc("Returns the I/O signature of the function it's called on"),

		// branches/turing/conditionals
		/*
			"if": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
				fFalse, fTrue, cond := st.Pop().(functiontype.ReqFunctionType), st.Pop().(functiontype.ReqFunctionType), st.Pop().(numbertype.ReqNumberType)

				fmt.Println(st.Slice(), fFalse.Literal())

				if cond.IsFloat() {
					return errors.New("cannot use float value as if condition")
				}

				if int32(cond.Literal().(float32)) != 0 {
					return callf(fTrue, sc, st)
				}

				return callf(fFalse, sc, st)
			}, []types.ReqVarType{types.TypeFunction, types.TypeFunction, types.TypeNumber}, []types.ReqVarType{}),
		*/

		// math
		"+": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			b := st.Pop()

			result, err := st.Pop().Add(b)
			if err != nil {
				return err
			}

			st.Push(result)

			return nil
		}, 2.1).SetDoc("Adds two values together"),

		"-": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			b := st.Pop()

			result, err := st.Pop().Sub(b)
			if err != nil {
				return err
			}

			st.Push(result)

			return nil
		}, 2.1).SetDoc("Subtracts one value from another"),

		"*": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			b := st.Pop()

			result, err := st.Pop().Mul(b)
			if err != nil {
				return err
			}

			st.Push(result)

			return nil
		}, 2.1).SetDoc("Multiplies one value with another"),

		"/": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			b := st.Pop()

			result, err := st.Pop().Div(b)
			if err != nil {
				return err
			}

			st.Push(result)

			return nil
		}, 2.1).SetDoc("Divides one value by another"),

		// stack operations
		"drop": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			st.Pop()

			return nil
		}, 1.0).SetDoc("Takes a value off the stack"),

		"dup": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			value := st.Pop()

			st.Push(value, value)

			return nil
		}, 1.2).SetDoc("Pops a value off the stack then pushes two of that value back onto the stack"),

		"dip": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			f := st.Pop().(functiontype.ReqFunctionType)
			value := st.Pop()

			err := callf(f, sc, st)
			if err != nil {
				return err
			}

			st.Push(value)

			return nil
		}, []types.ReqVarType{types.TypeFunction, types.TypeAny}, []types.ReqVarType{types.TypeAny, types.TypeAny}).SetDoc("Pops a value off the stack then calls a function"),
	},
	"runtime": {
		"version": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			st.Push(numbertype.New(runtime.REQPROC_VERSION))
			return nil
		}, []types.ReqVarType{}, []types.ReqVarType{types.TypeNumber}).SetDoc("Returns the current version of ReqProc"),

		"stacklen": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			st.Push(numbertype.New(float32(st.Len())))

			return nil
		}, []types.ReqVarType{}, []types.ReqVarType{types.TypeNumber}).SetDoc("Returns the length of the stack"),
	},
	//"os": {},
	"io": {
		"put": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			fmt.Print(st.Pop())
			return nil
		}, 1.0),

		"putl": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			fmt.Println(st.Pop())
			return nil
		}, 1.0),

		"dump": functiontype.NewSigNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			values := st.Slice()

			fmt.Println("stack contents:")
			for i := len(values) - 1; i > -1; i-- {
				if i == len(values)-1 {
					fmt.Printf(" [top] %d: %s\n", i, values[i].String())
				} else if i == 0 {
					fmt.Printf(" [bottom] %d: %s\n", i, values[i].String())
				} else {
					fmt.Printf(" %d: %s\n", i, values[i].String())
				}
			}

			return nil
		}, 0.0).SetDoc("Dumps the entire stack"),

		"readf": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			content, err := os.ReadFile(st.Pop().Literal().(string))
			if err != nil {
				return err
			}

			st.Push(stringtype.New(string(content)))

			return nil
		}, []types.ReqVarType{types.TypeString}, []types.ReqVarType{types.TypeString}),

		"writef": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			content := st.Pop().Literal().(string)

			file, err := os.Create(st.Pop().Literal().(string))
			if err != nil {
				return err
			}

			defer file.Close()

			_, err = file.WriteString(content)
			return err
		}, []types.ReqVarType{types.TypeString, types.TypeString}, []types.ReqVarType{}),
	},
	"web": {
		"download": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			resp, err := http.Get(st.Pop().Literal().(string))
			if err != nil {
				return err
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("status code %d, '%s'\n", resp.StatusCode, resp.Status)
			}

			content, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			st.Push(stringtype.New(string(content)))

			return nil
		}, []types.ReqVarType{types.TypeString}, []types.ReqVarType{types.TypeString}),
	},
	"strings": {
		"split": functiontype.NewNative(func(sc *scope.Scope, st *stack.Stack, callf func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error) error {
			delim := st.Pop().Literal().(string)

			spl := strings.Split(st.Pop().Literal().(string), delim)

			v := []types.ReqType{}

			for _, s := range spl {
				v = append(v, stringtype.New(s))
			}

			st.Push(listtype.New(v...))

			return nil
		}, []types.ReqVarType{types.TypeString, types.TypeString}, []types.ReqVarType{types.TypeString}),
	},
	//"ooptools": {},
	//"fptools":  {},
}
