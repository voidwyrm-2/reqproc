package lexer

import (
	"unicode"

	"github.com/voidwyrm-2/reqproc/lexer/tokens"
)

func isIdent(ch rune) bool {
	return !unicode.IsSpace(ch) && ch < 256 && ch != '!' && ch != '$' && ch != '@' && ch != '&' && ch != '(' && ch != ')' && ch != '[' && ch != ']'
}

type Lexer struct {
	text         string
	idx, col, ln int
	ch           rune
}

func New(text string) Lexer {
	l := Lexer{text: text, idx: -1, col: 0, ln: 1, ch: -1}
	l.advance()
	return l
}

func (l *Lexer) advance() {
	l.idx++
	l.col++

	if l.idx < len(l.text) {
		l.ch = rune(l.text[l.idx])
	} else {
		l.ch = -1
	}

	if l.ch == '\n' {
		l.ln++
		l.col = 1
	}
}

func (l *Lexer) peek() rune {
	if l.idx+1 < len(l.text) {
		return rune(l.text[l.idx+1])
	}

	return -1
}

func (l Lexer) errfp(col, ln int, format string, a ...any) error {
	return tokens.New(tokens.TokenKind(0), "", col, ln).Errf(format, a...)
}

func (l Lexer) errf(format string, a ...any) error {
	return l.errfp(l.col, l.ln, format, a...)
}

func (l Lexer) illch() error {
	return l.errf("illegal character '%c'", l.ch)
}

func (l Lexer) isNumber() bool {
	return l.ch >= '0' && l.ch <= '9'
}

func (l Lexer) isIdent() bool {
	return isIdent(l.ch)
}

func (l *Lexer) collectString(raw bool) (tokens.Token, error) {
	start := l.col
	startln := l.ln
	lit := ""
	escaped := false

	l.advance()

	qch := '"'
	if raw {
		qch = '`'
	}

	for l.ch != -1 {
		if escaped {
			switch l.ch {
			case '\\', '\'', '"':
				lit += string(l.ch)
			case 'n':
				lit += "\n"
			case 'f':
				lit += "\f"
			case 't':
				lit += "\t"
			case '0':
				lit += string(rune(0))
			case 'a':
				lit += string(rune(7))
			default:
				return tokens.Token{}, l.errf("invalid escape sequence character '%c'", l.ch)
			}
			escaped = false
		} else if l.ch == '\\' && !raw {
			escaped = true
		} else if l.ch == qch {
			break
		} else {
			lit += string(l.ch)
		}

		l.advance()
	}

	if l.ch == -1 {
		return tokens.Token{}, l.errfp(start, startln, "unterminated string literal")
	}

	l.advance()

	return tokens.New(tokens.String, lit, start, startln), nil
}

func (l *Lexer) collectNumber() tokens.Token {
	start := l.col
	startln := l.ln
	lit := ""
	dot := false

	for l.ch != -1 && (l.isNumber() || l.ch == '.') {
		if l.ch == '.' {
			if dot {
				break
			}

			dot = true
		}

		lit += string(l.ch)
		l.advance()
	}

	return tokens.New(tokens.Number, lit, start, startln)
}

func (l *Lexer) collectIdent(kind tokens.TokenKind, adv bool) tokens.Token {
	start := l.col
	startln := l.ln
	lit := ""

	if adv {
		l.advance()
	}

	for l.ch != -1 && l.isIdent() {
		lit += string(l.ch)
		l.advance()
	}

	return tokens.New(kind, lit, start, startln)
}

func (l *Lexer) Lex() ([]tokens.Token, error) {
	toks := []tokens.Token{}

	for l.ch != -1 {
		switch l.ch {
		case ';':
			for l.ch != -1 && l.ch != '\n' {
				l.advance()
			}
		case '(':
			toks = append(toks, tokens.New(tokens.ParenOpen, string(l.ch), l.col, l.ln))
			l.advance()
		case ')':
			toks = append(toks, tokens.New(tokens.ParenClose, string(l.ch), l.col, l.ln))
			l.advance()
		case ']':
			toks = append(toks, tokens.New(tokens.BracketClose, string(l.ch), l.col, l.ln))
			l.advance()
		case '&':
			if !isIdent(l.peek()) {
				return []tokens.Token{}, l.illch()
			}
			toks = append(toks, l.collectIdent(tokens.Call, true))
		case '!':
			if l.peek() == '[' {
				toks = append(toks, tokens.New(tokens.AssignIndex, string(l.ch)+"[", l.col, l.ln))
				l.advance()
				l.advance()
			} else if !isIdent(l.peek()) {
				return []tokens.Token{}, l.illch()
			} else {
				toks = append(toks, l.collectIdent(tokens.Assign, true))
			}
		case '$':
			if !isIdent(l.peek()) {
				return []tokens.Token{}, l.illch()
			}
			toks = append(toks, l.collectIdent(tokens.Const, true))
		case '@':
			if l.peek() == '[' {
				toks = append(toks, tokens.New(tokens.GetIndex, string(l.ch)+"[", l.col, l.ln))
				l.advance()
				l.advance()
			} else if !isIdent(l.peek()) {
				return []tokens.Token{}, l.illch()
			} else {
				toks = append(toks, l.collectIdent(tokens.GetValue, true))
			}
		case ':':
			toks = append(toks, l.collectIdent(tokens.Label, true))
		case '"':
			{
				t, err := l.collectString(false)
				if err != nil {
					return []tokens.Token{}, err
				}

				toks = append(toks, t)
			}
		case '`':
			{
				t, err := l.collectString(true)
				if err != nil {
					return []tokens.Token{}, err
				}

				toks = append(toks, t)
			}
		default:
			if l.isNumber() {
				toks = append(toks, l.collectNumber())
			} else if l.isIdent() {
				toks = append(toks, l.collectIdent(tokens.Ident, false))
			} else if unicode.IsSpace(l.ch) {
				l.advance()
			} else {
				return []tokens.Token{}, l.illch()
			}
		}
	}

	return toks, nil
}
