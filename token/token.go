package token

import "fmt"

type Token struct {
	Type    Type
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(
	Type Type,
	lexeme string,
	literal any,
	line int,
) Token {
	return Token{
		Type:    Type,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
}
