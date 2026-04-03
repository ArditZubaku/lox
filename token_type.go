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
	TokenTypeDot
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
