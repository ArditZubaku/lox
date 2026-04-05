package parser

import (
	"github.com/ArditZubaku/lox/token"
)

type ParseError struct {
	Token token.Token
	Msg   string
}

func (e ParseError) Error() string {
	return e.Msg
}
