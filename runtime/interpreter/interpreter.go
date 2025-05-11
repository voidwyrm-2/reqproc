package interpreter

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/voidwyrm-2/reqproc/lexer"
	"github.com/voidwyrm-2/reqproc/lexer/tokens"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/stack"
	"github.com/voidwyrm-2/reqproc/runtime/stdlib"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	"github.com/voidwyrm-2/reqproc/runtime/types/functiontype"
	"github.com/voidwyrm-2/reqproc/runtime/types/listtype"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
	"github.com/voidwyrm-2/reqproc/runtime/types/stringtype"
	"github.com/voidwyrm-2/reqproc/runtime/types/tabletype"
)

func expectKindsForwardInTokens(it int, toks []tokens.Token, kinds ...tokens.TokenKind) error {
	if len(kinds) == 0 {
		return nil
	} else if it+1 >= len(toks) {
		return toks[it].Errf("expected '%s', but found EOF", kinds[0].PublicString())
	}

	for i, k := range kinds {
		if !toks[it+i+1].Iskind(k) {
			return toks[it+i+1].Errf("expected '%s', but found '%s' instead", kinds[0].PublicString(), toks[it+i+1].Lit())
		}
	}

	return nil
}

func CallFunctionType(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack, sameStack bool) error {
	if err := st.Expect(rft.Input()...); err != nil {
		return err
	}

	lit := rft.Literal()

	if fn, ok := lit.(functiontype.NativeFunction); ok {
		return fn(sc, st, func(rft functiontype.ReqFunctionType, sc *scope.Scope, st *stack.Stack) error {
			return CallFunctionType(rft, sc, st, true)
		})
	}

	interp, err := New(sc)
	if err != nil {
		return err
	}

	if sameStack {
		interp.stack = st
	}

	res, err := interp.ExecuteTokens(lit.([]tokens.Token))
	st.Push(res...)

	return err
}

type Interpreter struct {
	scope   *scope.Scope
	stack   *stack.Stack
	modeTry bool
	err     string
}

func New(parentScope *scope.Scope) (Interpreter, error) {
	st := stack.New()
	i := Interpreter{scope: scope.New(parentScope, map[string]types.ReqType{}), stack: &st, err: "", modeTry: false}

	err := i.scope.LoadAllConst(stdlib.Stdlib["__init__"])
	if err != nil {
		return Interpreter{}, err
	}

	return i, nil
}

func (i Interpreter) GetScope() *scope.Scope {
	return i.scope
}

func (i Interpreter) GetStack() stack.Stack {
	return *i.stack
}

func (i *Interpreter) StackPush(values ...types.ReqType) {
	i.stack.Push(values...)
}

func (i *Interpreter) StackPop() (types.ReqType, bool) {
	if i.stack.Len() == 0 {
		return nil, false
	}

	return i.stack.Pop(), true
}

func (i Interpreter) StackLen() int {
	return i.stack.Len()
}

func (i Interpreter) GetErr() string {
	return i.err
}

func (i *Interpreter) SetErr(e string) {
	i.err = e
}

func (i Interpreter) GetModeTry() bool {
	return i.modeTry
}

func (i *Interpreter) SetModeTry(m bool) {
	i.modeTry = m
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
			return expectKindsForwardInTokens(it, toks, kinds...)
		}

		collectInPairForward := func(anchor tokens.Token, a, b tokens.TokenKind, filter func(tokens.Token) bool) ([]tokens.Token, error) {
			contained := []tokens.Token{}

			nest := 0
			for _, t := range toks[it+1:] {
				if t.Iskind(a) {
					nest++
				} else if t.Iskind(b) {
					if nest > 0 {
						nest--
					} else {
						return contained, nil
					}
				} else if !filter(t) {
					return []tokens.Token{}, t.Errf("token '%s' is not valid inside %s %s", t.Kind().PublicString(), a.PublicString(), b.PublicString())
				}

				contained = append(contained, t)
			}

			if len(contained) == 0 {
				contained = append(contained, anchor)
			}

			return []tokens.Token{}, contained[len(contained)-1].Errf("no '%s' to match '%s'", b.PublicString(), a.PublicString())
		}

		switch cur.Kind() {
		case tokens.Label:
			it++
		case tokens.String:
			i.stack.Push(stringtype.New(cur.Lit()))
			it++
		case tokens.Number:
			if val, err := numbertype.FromString(toks[it].Lit()); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else {
				i.stack.Push(val)
			}
			it++
		case tokens.Ident:
			switch cur.Lit() {
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
					it += 2
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

							interp, err := New(nil)
							if err != nil {
								return err
							}

							_, err = interp.Execute(string(content))
							if err != nil {
								return err
							}

							return i.scope.WriteConst(modname, tabletype.New(interp.scope.Consts()))
						}()
						if err != nil {
							if i.modeTry {
								i.err = cur.Err(err).Error()
							} else {
								return []types.ReqType{}, err
							}
						}
					} else {
						if mod, ok := stdlib.Stdlib[modname]; ok && !strings.HasPrefix(modname, "__") {
							if err := i.scope.WriteConst(modname, tabletype.New(mod)); err != nil {
								return []types.ReqType{}, cur.Err(err)
							}
						} else {
							return []types.ReqType{}, cur.Errf("module '%s' does not exist in the standard library", modname)
						}
					}
				}
				it++
			case "try":
				i.modeTry = true
				it++
			case "notry":
				i.modeTry = false
				it++
			default:
				if v, err := i.scope.Read(cur.Lit()); err != nil { // does that variable or const exist?
					return []types.ReqType{}, cur.Err(err)
				} else if v.Type() != types.TypeFunction { // can we call it?
					return []types.ReqType{}, cur.Errf("'%s' is not callable", v.Type().String())
				} else { // all good, let's call it
					if err := CallFunctionType(v.(functiontype.ReqFunctionType), i.scope, i.stack, false); err != nil {
						if strings.HasPrefix(err.Error(), "EXIT CODE ") { // special handling for the exit function
							return []types.ReqType{}, err
						}

						if i.modeTry {
							i.err = cur.Err(err).Error()
						} else {
							return []types.ReqType{}, cur.Err(err)
						}
					}
				}
				it++
			}
		case tokens.GetValue:
			if f, ok := stdlib.Stdlib["__keyword__"][cur.Lit()]; ok {
				i.stack.Push(f)
			} else if v, err := i.scope.Read(cur.Lit()); err != nil {
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
		case tokens.GetIndex:
			// we don't check for any specific type because 1. we don't have to change it for any future indexables, and 2. the GetIndex functions handles that
			// index, indexable
			if err := i.stack.Expect(types.TypeAny, types.TypeAny); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else {
				index, indexable := i.stack.Pop(), i.stack.Pop()

				result, err := indexable.GetIndex(index)
				if err != nil {
					return []types.ReqType{}, cur.Err(err)
				}

				i.stack.Push(result)
			}
			it++
		case tokens.AssignIndex:
			// see the commend under `case tokens.GetIndex` for why we aren't checking the types
			// item, index, indexable
			if err := i.stack.Expect(types.TypeAny, types.TypeAny, types.TypeAny); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else {
				item, index, indexable := i.stack.Pop(), i.stack.Pop(), i.stack.Pop()

				err := indexable.SetIndex(index, item)
				if err != nil {
					return []types.ReqType{}, cur.Err(err)
				}

				i.stack.Push(indexable)
			}
			it++
		case tokens.Const:
			if err := i.stack.Expect(types.TypeAny); err != nil {
				return []types.ReqType{}, cur.Err(err)
			} else if err = i.scope.WriteConst(cur.Lit(), i.stack.Pop()); err != nil {
				return []types.ReqType{}, cur.Err(err)
			}
			it++
		case tokens.ParenOpen:
			// collect the tokens until the matching ')'
			// the function keeps track of how many pairs deep we are
			if fcontent, err := collectInPairForward(toks[it], tokens.ParenOpen, tokens.ParenClose, func(t tokens.Token) bool {
				return true
			}); err != nil {
				if i.modeTry {
					i.err = err.Error()
				} else {
					return []types.ReqType{}, err
				}
			} else {
				if err = expectKindsForwardInTokens(-1, fcontent, tokens.Signature); err != nil {
					return []types.ReqType{}, err
				}

				sig, err := strconv.ParseFloat(fcontent[0].Lit(), 32)
				if err != nil {
					return []types.ReqType{}, cur.Err(err)
				}

				i.stack.Push(functiontype.New(fcontent[1:], float32(sig)))

				it += len(fcontent) + 2
			}
		case tokens.BracketOpen:
			if lcontent, err := collectInPairForward(toks[it], tokens.BracketOpen, tokens.BracketClose, func(t tokens.Token) bool {
				return t.Iskind(tokens.GetValue) || t.Iskind(tokens.String) || t.Iskind(tokens.Number)
			}); err != nil {
				if i.modeTry {
					i.err = err.Error()
				} else {
					return []types.ReqType{}, err
				}
			} else {
				list := []types.ReqType{}

				for _, t := range lcontent {
					if t.Iskind(tokens.GetValue) {
						if v, err := i.scope.Read(t.Lit()); err != nil {
							return []types.ReqType{}, err
						} else {
							list = append(list, v)
						}
					} else if t.Iskind(tokens.String) {
						list = append(list, stringtype.New(t.Lit()))
					} else if t.Iskind(tokens.Number) {
						if val, err := numbertype.FromString(t.Lit()); err != nil {
							return []types.ReqType{}, cur.Err(err)
						} else {
							list = append(list, val)
						}
					} else {
						panic(fmt.Sprintf("somehow found token %s (or, `%s`) during list parsing", t.String(), t.Kind().PublicString()))
					}
				}

				i.stack.Push(listtype.New(list...))
				it += len(lcontent) + 2
			}
		default:
			return []types.ReqType{}, cur.Errf("unexpected token '%s'", cur.Lit())
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
