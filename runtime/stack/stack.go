package stack

import (
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type Stack struct {
	stack []types.ReqType
}

func New(values ...types.ReqType) Stack {
	return Stack{stack: values}
}

func (s *Stack) Push(values ...types.ReqType) {
	s.stack = append(s.stack, values...)
}

func (s Stack) Expect(kinds ...types.ReqVarType) error {
	if len(kinds) == 0 {
		return nil
	}

	gkstr := func(k types.ReqVarType, exp string) error {
		kstr := "a " + k.String()
		if k == types.TypeAny {
			kstr = "any type"
		}

		return fmt.Errorf("expected %s on the stack but %s", kstr, exp)
	}

	if len(s.stack) == 0 {
		return gkstr(kinds[0], "the stack is empty")
	} else if len(s.stack) < len(kinds) {
		return gkstr(kinds[len(s.stack)], "the stack isn't large enough")
	}

	a := 0
	b := len(s.stack) - 1
	for a < len(kinds) && b > -1 {
		if kinds[a]&s.stack[b].Type() != s.stack[b].Type() {
			// fmt.Printf("%s, %s, %v\n", kinds[a], s.stack[b].Type().String(), kinds[a]&s.stack[b].Type() == s.stack[b].Type())
			return gkstr(kinds[a], fmt.Sprintf("found '%s(type %s)' instead", s.stack[b].String(), s.stack[b].Type().String()))
		}

		a++
		b--
	}

	return nil
}

func (s *Stack) Pop() types.ReqType {
	value := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return value
}

func (s Stack) Len() int {
	return len(s.stack)
}

func (s Stack) Slice() []types.ReqType {
	return s.stack
}

func (s Stack) String() string {
	return fmt.Sprint(s.stack)
}
