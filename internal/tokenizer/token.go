package tokenizer

type TokenKind struct {
	identifier string
}

type Token struct {
	Position
	kind TokenKind
	text string
}

func (t Token) Kind() TokenKind {
	return t.kind
}

func (t Token) Text() string {
	return t.text
}

var (
	Invalid = TokenKind{"invalid"}
	EOF     = TokenKind{"EOF"}
	Comment = TokenKind{"comment"}

	Identifier = TokenKind{"identifier"}
	String     = TokenKind{"string"}
	Variable   = TokenKind{"variable"}

	// Keywords
	Clone   = TokenKind{"CLONE"}
	Run     = TokenKind{"RUN"}
	In      = TokenKind{"IN"}
	Extract = TokenKind{"EXTRACT"}
	End     = TokenKind{"END"}
	As      = TokenKind{"AS"}
	Set     = TokenKind{"SET"}

	SqrBracketLeft  = TokenKind{"["}
	SqrBracketRight = TokenKind{"["}

	BracketLeft  = TokenKind{"("}
	BracketRight = TokenKind{")"}

	// Separators
	Semicolon = TokenKind{";"}
)
