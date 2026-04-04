//go:generate stringer -type=TokenType

package main

type TokenType byte

const (
	// Single-character tokens.
	LeftParen TokenType = iota // LeftParen
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals.
	Identifier
	String
	Number

	// Keywords.
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

// TODO: Try replacing this with the https://github.com/lemire/constmap just for fun
// keywords is a map of the reserved keywords
var keywords = map[string]TokenType{
	"and":   And,
	"class": Class,
	"else":  Else,
	"false": False,
	"for":   For,
	// TODO: I might replace this with `fn` at some point, I like fn better
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}
