package stdlib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/functiontype"
	"github.com/voidwyrm-2/reqproc/runtime/types/listtype"
	"github.com/voidwyrm-2/reqproc/runtime/types/stringtype"
)

var Stdlib = map[string]map[string]types.ReqType{
	"os": {
		"exit": functiontype.NewNative(1, 0, func(sc *scope.Scope, st *stack.Stack) error {
			if err := st.Expect(types.TypeNumber); err != nil {
				return err
			}

			return fmt.Errorf("%d EXIT", int(st.Pop().Literal().(float32)))
		}),
	},
	"io": {
		"put": functiontype.NewNative(1, 0, func(sc *scope.Scope, st *stack.Stack) error {
			if err := st.Expect(types.TypeAny); err != nil {
				return err
			}

			fmt.Print(st.Pop())
			return nil
		}),
		"putl": functiontype.NewNative(1, 9, func(sc *scope.Scope, st *stack.Stack) error {
			if err := st.Expect(types.TypeAny); err != nil {
				return err
			}

			fmt.Println(st.Pop())
			return nil
		}),
		"readf": functiontype.NewNative(1, 1, func(sc *scope.Scope, st *stack.Stack) error {
			if err := st.Expect(types.TypeString); err != nil {
				return err
			}

			file, err := os.Open(st.Pop().Literal().(string))
			defer file.Close()
			if err != nil {
				return err
			}

			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			st.Push(stringtype.New(string(content)))

			return nil
		}),
		"writef": functiontype.NewNative(2, 0, func(sc *scope.Scope, st *stack.Stack) error {
			if err := st.Expect(types.TypeString, types.TypeString); err != nil {
				return err
			}

			content := st.Pop().Literal().(string)

			file, err := os.Create(st.Pop().Literal().(string))
			defer file.Close()
			if err != nil {
				return err
			}

			_, err = file.WriteString(content)
			return err
		}),
	},
	"web": {
		"download": functiontype.NewNative(1, 1, func(sc *scope.Scope, st *stack.Stack) error {
			if err := st.Expect(types.TypeString); err != nil {
				return err
			}

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
		}),
	},
	"strings": {
		"split": functiontype.NewNative(1, 1, func(sc *scope.Scope, st *stack.Stack) error {
			if err := st.Expect(types.TypeString, types.TypeString); err != nil {
				return err
			}

			delim := st.Pop().Literal().(string)

			spl := strings.Split(st.Pop().Literal().(string), delim)

			v := []types.ReqType{}

			for _, s := range spl {
				v = append(v, stringtype.New(s))
			}

			st.Push(listtype.New(v...))

			return nil
		}),
	},
}
