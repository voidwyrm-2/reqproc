package stack

import (
	"errors"
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type Stack struct {
	stack []types.ReqType
}

func New() Stack {
	return Stack{stack: []types.ReqType{}}
}

func (s *Stack) Push(value types.ReqType) {
	s.stack = append(s.stack, value)
}

func (s *Stack) Pop() (types.ReqType, error) {
	if len(s.stack) == 0 {
		return nil, errors.New("stack underflow")
	}

	value := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return value, nil
}

func (s Stack) Len() int {
	return len(s.stack)
}

func (s Stack) String() string {
	return fmt.Sprint(s.stack)
}
