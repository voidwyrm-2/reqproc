package interpreter

import (
	"io"
	"os"
	"path"
	"strconv"

	"github.com/voidwyrm-2/reqproc/lexer"
	"github.com/voidwyrm-2/reqproc/lexer/tokens"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/stdlib"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/functiontype"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
	"github.com/voidwyrm-2/reqproc/runtime/types/stringtype"
	"github.com/voidwyrm-2/reqproc/runtime/types/tabletype"
)

func CallFunctionType(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error {
	lit := rft.Literal()

	if fn, ok := lit.(func(sc *scope.Scope, st *stack.Stack) error); ok {
		return fn(sc, st)
	}

	interp := New(sc)

	res, err := interp.ExecuteTokens(lit.([]tokens.Token))
	st.Push(res...)

	return err
}

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

	labels := map[string]int{}

	for i, t := range toks {
		if t.Iskind(tokens.Label) {
			if _, ok := labels[t.Lit()]; ok {
				return []types.ReqType{}, t.Errf("cannot redefine existing label '%s'", t.Lit())
			}

			labels[t.Lit()] = i
		}
	}

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
		case tokens.Label:
			it++
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
				} else if err = i.scope.Write(next.Lit(), nil); err != nil {
					return []types.ReqType{}, next.Err(err)
				}
				it += 2
			case "err":
				if err := expectKindsForward(tokens.Ident); err != nil {
					return []types.ReqType{}, cur.Err(err)
				} else if l, ok := labels[next.Lit()]; !ok {
					return []types.ReqType{}, next.Errf("label '%s' is not defined", next.Lit())
				} else if i.err != "" {
					it = l
				} else {
					it++
				}
			case "geterr":
				i.stack.Push(stringtype.New(i.err))
				it++
			case "errcl":
				i.err = ""
				it++
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

					if path.Ext(modname) == ".req" {
						err := func() error {
							mod, err := os.Open(modname)
							defer mod.Close()
							if err != nil {
								return err
							}

							content, err := io.ReadAll(mod)
							if err != nil {
								return err
							}

							interp := New(nil)

							_, err = interp.Execute(string(content))
							if err != nil {
								return err
							}

							return i.scope.WriteConst(modname, tabletype.New(interp.scope.Consts()))
						}()
						if err != nil {
							i.err = cur.Err(err).Error()
						}
					} else {
						if mod, ok := stdlib.Stdlib[modname]; ok {
							if err := i.scope.WriteConst(modname, tabletype.New(mod)); err != nil {
								return []types.ReqType{}, cur.Err(err)
							}
						} else {
							return []types.ReqType{}, cur.Errf("module '%s' does not exist in the standard library", modname)
						}
					}
				}
				it++
			default:
				return []types.ReqType{}, toks[it].Errf("unexpected identifier '%s'", toks[it].Lit())
			}
		case tokens.Call:
			if v, err := i.scope.Read(cur.Lit()); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else if v.Type() != types.TypeFunction {
				return []types.ReqType{}, cur.Errf("'%s' is not callable", v.Type().String())
			} else {
				if err := CallFunctionType(v.(functiontype.ReqFunctionType), i.scope, &i.stack); err != nil {
					i.err = cur.Err(err).Error()
				}
			}
			it++
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
