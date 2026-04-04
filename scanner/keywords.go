package scanner

import "github.com/ArditZubaku/lox/token"

// TODO: Try replacing this with the https://github.com/lemire/constmap just for fun
// keywords is a map of the reserved keywords
var keywords = map[string]token.TokenType{
	"and":    token.And,
	"class":  token.Class,
	"else":   token.Else,
	"false":  token.False,
	"for":    token.For,
	"fun":    token.Fun, // TODO: I might replace this with `fn` at some point, I like fn better
	"if":     token.If,
	"nil":    token.Nil,
	"or":     token.Or,
	"print":  token.Print,
	"return": token.Return,
	"super":  token.Super,
	"this":   token.This,
	"true":   token.True,
	"var":    token.Var,
	"while":  token.While,
}
