package stack

import (
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type Stack struct {
	stack []types.ReqType
}

func New() Stack {
	return Stack{stack: []types.ReqType{}}
}

func (s *Stack) Push(values ...types.ReqType) {
	s.stack = append(s.stack, values...)
}

func (s *Stack) Expect(types ...types.ReqVarType) error {
	if len(s.stack) == 0 {
		return fmt.Errorf("expected a %s on the stack but the stack is empty", types[0])
	} else if len(s.stack) < len(types) {
		return fmt.Errorf("expected a %s on the stack but the is not large enough", types[len(s.stack)])
	}

	a := 0
	b := len(s.stack) - 1
	for a < len(types) && b > -1 {
		if types[a]&s.stack[b].Type() != s.stack[b].Type() {
			return fmt.Errorf("expected a %s on the stack but found a %s instead", types[a], s.stack[b])
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
