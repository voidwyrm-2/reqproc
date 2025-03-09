package tokens

import (
	"fmt"
)

type TokenKind uint8

const (
	String TokenKind = iota
	Number
	Ident
	Label
	Assign
	Const
	GetValue
	Call
	Plus
	Hyphen
	Asterisk
	ForwardSlash
	GreaterThan
	LessThan
	Caret
	Concat
	Equals
	NotEquals
	ParenOpen
	ParenClose
)

var invalidIdents = []string{
	"and",
	"not",
	"or",
	"band",
	"bor",
	"bnot",
	"xor",
}

type Token struct {
	kind    TokenKind
	lit     string
	col, ln int
}

func New(kind TokenKind, lit string, col, ln int) Token {
	return Token{kind: kind, lit: lit, col: col, ln: ln}
}

func (t Token) String() string {
	return fmt.Sprintf("{%d `%s` %d %d}", t.kind, t.lit, t.col, t.ln)
}
