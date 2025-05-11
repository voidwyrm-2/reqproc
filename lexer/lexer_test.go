package lexer

import (
	"testing"

	"github.com/voidwyrm-2/reqproc/lexer/tokens"
)

type expectedToken = struct {
	kind tokens.TokenKind
	lit  string
}

type testCase struct {
	input    string
	expected []expectedToken
}

func TestLexer(t *testing.T) {
	cases := []testCase{
		{
			"0 exit",
			[]expectedToken{
				{tokens.Number, "0"},
				{tokens.Ident, "exit"},
			},
		},
		{
			"(|1.1 ref)",
			[]expectedToken{
				{tokens.ParenOpen, "("},
				{tokens.Signature, "1.1"},
				{tokens.Ident, "ref"},
				{tokens.ParenClose, ")"},
			},
		},
		{
			"-1 20 30 40 + -",
			[]expectedToken{
				{tokens.Number, "-1"},
				{tokens.Number, "20"},
				{tokens.Number, "30"},
				{tokens.Number, "40"},
				{tokens.Ident, "+"},
				{tokens.Ident, "-"},
			},
		},
		{
			`[0 2 1 5 3 4] $input
@input (|2.1 @input swap @#) each`,
			[]expectedToken{
				{tokens.BracketOpen, "["},
				{tokens.Number, "0"},
				{tokens.Number, "2"},
				{tokens.Number, "1"},
				{tokens.Number, "5"},
				{tokens.Number, "3"},
				{tokens.Number, "4"},
				{tokens.BracketClose, "]"},
				{tokens.Const, "input"},
				{tokens.GetValue, "input"},
				{tokens.ParenOpen, "("},
				{tokens.Signature, "2.1"},
				{tokens.GetValue, "input"},
				{tokens.Ident, "swap"},
				{tokens.GetIndex, "@#"},
				{tokens.ParenClose, ")"},
				{tokens.Ident, "each"},
			},
		},
		{
			"+++",
			[]expectedToken{
				{tokens.Ident, "+"},
				{tokens.Ident, "+"},
				{tokens.Ident, "+"},
			},
		},
		{
			"cat++dog",
			[]expectedToken{
				{tokens.Ident, "cat"},
				{tokens.Ident, "+"},
				{tokens.Ident, "+"},
				{tokens.Ident, "dog"},
			},
		},
		{
			`"io" import
"Hello there." io.putl`,
			[]expectedToken{
				{tokens.String, "io"},
				{tokens.Ident, "import"},
				{tokens.String, "Hello there."},
				{tokens.Ident, "io.putl"},
			},
		},
	}

	for _, c := range cases {
		l := New(c.input)
		if actual, err := l.Lex(); err != nil {
			t.Fatal(err.Error())
		} else {
			if len(actual) != len(c.expected) {
				t.Fatalf("expected %d tokens, but found %d instead with `%s`\noutput: %v", len(c.expected), len(actual), c.input, actual)
			}

			for i, a := range actual {
				e := c.expected[i]
				if !a.Is(e.kind, e.lit) {
					t.Fatalf("expected (%s, '%s') but found (%s, '%s') instead with `%s`", e.kind.String(), e.lit, a.Kind().String(), a.Lit(), c.input)
				}
			}
		}
	}
}
