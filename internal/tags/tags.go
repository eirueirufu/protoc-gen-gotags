package tags

import (
	"errors"
	"strings"
	"text/scanner"
)

//go:generate goyacc tags.y
type (
	Tags []*Tag
	Tag  struct {
		Key, Value string
	}

	Lex struct {
		scanner scanner.Scanner
		root    Tags
		err     error
	}
)

func ParseTags(src string) (Tags, error) {
	lex := newLex(src)
	if yyParse(lex) != 0 {
		return nil, lex.err
	}
	return lex.root, nil
}

func newLex(src string) *Lex {
	yyErrorVerbose = true
	lex := new(Lex)
	lex.scanner.Init(strings.NewReader(src))
	lex.scanner.Mode = scanner.ScanIdents |
		scanner.ScanStrings
	return lex
}

func (lex *Lex) Lex(lval *yySymType) int {
	tok := lex.scanner.Scan()
	switch tok {
	case scanner.EOF:
		return 0
	case scanner.Ident:
		identStr := lex.scanner.TokenText()
		lval.tok = identStr
		if identStr == ":" {
			return colon
		}
		return ident
	case scanner.String:
		str := lex.scanner.TokenText()
		str = str[1 : len(str)-1]
		lval.tok = str
		return string_lit
	default:
		return int(tok)
	}
}

func (lex *Lex) Error(s string) {
	lex.err = errors.New(s)
}
