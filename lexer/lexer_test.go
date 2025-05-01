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
				{
					tokens.Number,
					"0",
				},
				{
					tokens.Ident,
					"exit",
				},
			},
		},
		{
			"(|1.1 ref)",
			[]expectedToken{
				{
					tokens.ParenOpen,
					"(",
				},
				{
					tokens.Signature,
					"1.1",
				},
				{
					tokens.Ident,
					"ref",
				},
				{
					tokens.ParenClose,
					")",
				},
			},
		},
		{
			"-1 20 30 40 + -",
			[]expectedToken{
				{
					tokens.Number,
					"-1",
				},
				{
					tokens.Number,
					"20",
				},
				{
					tokens.Number,
					"30",
				},
				{
					tokens.Number,
					"40",
				},
				{
					tokens.Ident,
					"+",
				},
				{
					tokens.Ident,
					"-",
				},
			},
		},
	}

	for _, c := range cases {
		l := New(c.input)
		if actual, err := l.Lex(); err != nil {
			t.Fatal(err.Error())
		} else {
			if len(actual) != len(c.expected) {
				t.Fatalf("expected %d tokens, but found %d instead", len(c.expected), len(actual))
			}

			for i, a := range actual {
				e := c.expected[i]
				if !a.Is(e.kind, e.lit) {
					t.Fatalf("expected (%s, '%s') but found (%s, '%s') instead", e.kind.String(), e.lit, a.Kind().String(), a.Lit())
				}
			}
		}
	}
}
