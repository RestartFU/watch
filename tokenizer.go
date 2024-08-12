package main

import (
	"errors"
	"unicode/utf8"
)

const (
	RuneBom   = 0xfeff
	RuneEof   = 0
	RuneError = 65533
)

type pos struct {
	offset int
	line   int
	column int
}

type tokenKind struct {
	identifier string
}

type token struct {
	pos
	kind tokenKind
	text string
}

type tokenizer struct {
	pos
	data              string
	r                 rune
	w                 int
	currentLineOffset int
	insertSemicolon   bool
}

func newTokenizer(data string) *tokenizer {
	t := &tokenizer{
		pos:  pos{line: 1},
		data: data,
	}
	t.Next()
	if t.r == RuneBom {
		t.Next()
	}
	return t
}

func (t *tokenizer) Next() rune {
	if t.offset >= len(t.data) {
		return RuneEof
	} else {
		t.offset += t.w
		t.r, t.w = utf8.DecodeRuneInString(t.data[t.offset:])
		t.column = t.offset - t.currentLineOffset
		if t.offset >= len(t.data) {
			return RuneEof
		}
	}
	return t.r
}

func (t *tokenizer) skipWhiteSpace(newLine bool) {
loop:
	for t.offset < len(t.data) {
		switch t.r {
		case ' ', '\t', '\v', '\f', '\r':
			t.Next()
		case '\n':
			if newLine {
				break loop
			}
			t.line++
			t.currentLineOffset = t.offset
			t.column = 1
			t.Next()
		default:
			switch t.r {
			case 0x2028, 0x2029, 0xFEFF:
				t.Next()
				continue loop
			}
			break loop
		}
	}
}

func isLetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}
func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAny(r rune, runes ...rune) bool {
	for _, v := range runes {
		if r == v {
			return true
		}
	}
	return false
}

func (t *tokenizer) Token() (token token, err error) {
	t.skipWhiteSpace(t.insertSemicolon)

	token.pos = t.pos
	token.kind = Invalid

	var currRune rune = t.r
	t.Next()
	switch {
	case currRune == RuneError:
		// TODO: actually make this work.
		token.kind = EOF

	case currRune == RuneEof:
		token.kind = EOF
		err = errors.New("EOF")

	case isAny(currRune, '\n', ';', '\\'):
		t.insertSemicolon = false
		token.text = "\n"
		token.kind = Semicolon
		t.line++
		t.currentLineOffset = t.offset
		t.column = 1
		return
	case isLetter(currRune):
		token.kind = Identifier
		for t.offset < len(t.data) {
			switch {
			case isLetter(t.r) || isDigit(t.r) || t.r == '_':
				t.Next()
				continue
			}
			break
		}
		for _, v := range []tokenKind{
			Clone,
			Begin,
			Extract,
			Then,
			End,
			With,
		} {
			if t.data[token.offset:t.offset] == v.identifier {
				token.kind = v
				break
			}
		}
	case currRune == '#':
		token.kind = Comment
		for t.offset < len(t.data) {
			if t.r == '\n' {
				break
			}
			t.Next()
		}
	case currRune == ':':
		token.kind = String
		for t.offset < len(t.data) {
			switch {
			case t.r == '\n':
				t.Next()
				break
			default:
				t.Next()
				continue
			}
			break
		}
	default:
		err = errors.New("invalid character")
	}
	text := t.data[token.offset:t.offset]
	switch token.kind {
	case String:
		text = text[2:]
	case EOF, Semicolon:
		t.insertSemicolon = false
	case Identifier:
		t.insertSemicolon = true
	default:
		t.insertSemicolon = false
	}
	token.text = text
	return
}

var (
	Invalid = tokenKind{"invalid"}
	EOF     = tokenKind{"EOF"}
	Comment = tokenKind{"comment"}

	Identifier = tokenKind{"identifier"}
	String     = tokenKind{"string"}

	// Keywords
	Clone   = tokenKind{"CLONE"}
	Begin   = tokenKind{"BEGIN"}
	Extract = tokenKind{"EXTRACT"}
	Then    = tokenKind{"THEN"}
	End     = tokenKind{"END"}
	With    = tokenKind{"WITH"}

	// Separators
	Semicolon = tokenKind{";"}
)
