package interpreter

import (
	"github.com/voidwyrm-2/reqproc/lexer/tokens"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type Interpreter struct {
	scope *scope.Scope
	stack stack.Stack
	err   string
}

func New(parentScope *scope.Scope) Interpreter {
	return Interpreter{scope: parentScope, stack: stack.New(), err: ""}
}

func (i Interpreter) GetScope() *scope.Scope {
	return i.scope
}

func (i *Interpreter) ExecuteTokens(toks tokens.Token) ([]types.ReqType, error) {
	return []types.ReqType{}, nil
}
