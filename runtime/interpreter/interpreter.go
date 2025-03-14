package interpreter

import (
	"io"
	"os"
	"strconv"

	"github.com/voidwyrm-2/reqproc/lexer"
	"github.com/voidwyrm-2/reqproc/lexer/tokens"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
	"github.com/voidwyrm-2/reqproc/runtime/types/stringtype"
	"github.com/voidwyrm-2/reqproc/runtime/types/tabletype"
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

func (i *Interpreter) ExecuteTokens(toks []tokens.Token) ([]types.ReqType, error) {
	it := 0

	for it < len(toks) {
		cur := toks[it]
		next := tokens.Token{}
		if it+1 < len(toks) {
			next = toks[it+1]
		}

		expectKindsForward := func(kinds ...tokens.TokenKind) error {
			if len(kinds) == 0 {
				return nil
			} else if it+1 >= len(toks) {
				return toks[it].Errf("expected '%s', but found EOF", kinds[0].PublicString())
			}

			for i, k := range kinds {
				if !toks[it+i+1].Iskind(k) {
					return toks[it].Errf("expected '%s', but found '%s' instead", kinds[0].PublicString(), toks[it+i+1].Lit())
				}
			}

			return nil
		}

		switch cur.Kind() {
		case tokens.String:
			i.stack.Push(stringtype.New(cur.Lit()))
			it++
		case tokens.Number:
			if n, err := strconv.ParseFloat(toks[it].Lit(), 32); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else {
				i.stack.Push(numbertype.New(float32(n)))
			}
			it++
		case tokens.Ident:
			switch toks[it].Lit() {
			case "def":
				if err := expectKindsForward(tokens.Ident); err != nil {
					return []types.ReqType{}, cur.Err(err)
				}

				i.scope.Write(next.Lit(), nil)
				it += 2
			case "true":
				i.stack.Push(numbertype.New(1))
				it++
			case "false":
				i.stack.Push(numbertype.New(0))
				it++
			case "import":
				if err := i.stack.Expect(types.TypeString); err != nil {
					return []types.ReqType{}, cur.Err(err)
				} else {
					modname := i.stack.Pop().Literal().(string)

					mod, err := os.Open(modname)
					if err != nil {
						return []types.ReqType{}, cur.Err(err)
					}

					content, err := io.ReadAll(mod)
					if err != nil {
						return []types.ReqType{}, cur.Err(err)
					}

					interp := New(nil)

					_, err = interp.Execute(string(content))
					if err != nil {
						return []types.ReqType{}, cur.Err(err)
					}

					i.scope.WriteConst(modname, tabletype.New(interp.scope.Consts()))
				}
				it++
			default:
				return []types.ReqType{}, toks[it].Errf("unexpected identifier '%s'", toks[it].Lit())
			}
		case tokens.GetValue:
			if v, err := i.scope.Read(cur.Lit()); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else {
				i.stack.Push(v)
			}
			it++
		case tokens.Assign:
			if err := i.stack.Expect(types.TypeAny); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else if err = i.scope.Update(cur.Lit(), i.stack.Pop()); err != nil {
				return []types.ReqType{}, cur.Err(err)
			}
			it++
		case tokens.Const:
			if err := i.stack.Expect(types.TypeAny); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else if err = i.scope.WriteConst(cur.Lit(), i.stack.Pop()); err != nil {
				return []types.ReqType{}, cur.Err(err)
			}
			it++
		default:
			return []types.ReqType{}, toks[it].Errf("unexpected token '%s'", toks[it].Lit())
		}
	}

	return i.stack.Slice(), nil
}

func (i *Interpreter) Execute(text string) ([]types.ReqType, error) {
	l := lexer.New(text)

	tokens, err := l.Lex()
	if err != nil {
		return []types.ReqType{}, err
	}

	return i.ExecuteTokens(tokens)
}
