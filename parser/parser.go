package parser

import (
	"slices"

	"github.com/ArditZubaku/lox/expr"
	"github.com/ArditZubaku/lox/token"
)

type Parser[T any] struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser[any] {
	return &Parser[any]{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser[T]) expression() expr.Expr[T] {
	return p.equality()
}

func (p *Parser[T]) equality() expr.Expr[T] {
	expression := p.comparison()

	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		// TODO: Rethink the `any`
		expression = expr.Binary[T]{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser[T]) match(types ...token.Type) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser[T]) comparison() expr.Expr[T] {
	panic("unimplemented")
}

func (p *Parser[T]) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

func (p *Parser[T]) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser[T]) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser[T]) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser[T]) previous() token.Token {
	return p.tokens[p.current-1]
}
