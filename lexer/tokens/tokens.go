package tokens

import (
	"fmt"
)

type TokenKind uint8

const (
	None TokenKind = iota
	String
	Number
	Ident
	Label
	Assign
	AssignIndex
	Const
	GetValue
	GetIndex
	Call
	ParenOpen
	ParenClose
	BracketClose
)

func (tk TokenKind) PublicString() string {
	return []string{
		"string",
		"number",
		"identifier",
		":",
		"!",
		"![",
		"$",
		"@",
		"@[",
		"&",
		"(",
		")",
		"]",
	}[tk]
}

func (tk TokenKind) String() string {
	return []string{
		"None",
		"String",
		"Number",
		"Ident",
		"Label",
		"Assign",
		"AssignIndex",
		"Const",
		"GetValue",
		"GetIndex",
		"Call",
		"ParenOpen",
		"ParenClose",
		"BracketClose",
	}[tk]
}

type Token struct {
	kind    TokenKind
	lit     string
	col, ln int
}

func New(kind TokenKind, lit string, col, ln int) Token {
	return Token{kind: kind, lit: lit, col: col, ln: ln}
}

func (t Token) Iskind(kind TokenKind) bool {
	return t.kind == kind
}

func (t Token) Islit(lit string) bool {
	return t.lit == lit
}

func (t Token) Kind() TokenKind {
	return t.kind
}

func (t Token) Lit() string {
	return t.lit
}

func (t Token) Errf(format string, a ...any) error {
	return fmt.Errorf("error on line %d, col %d: %s", t.ln, t.col, fmt.Sprintf(format, a...))
}

func (t Token) Err(err error) error {
	if err == nil {
		return nil
	}

	return t.Errf(err.Error())
}

func (t Token) String() string {
	return fmt.Sprintf("{%s `%s` %d %d}", t.kind.String(), t.lit, t.col, t.ln)
}
