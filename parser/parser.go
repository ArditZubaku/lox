package parser

import (
	"slices"

	"github.com/ArditZubaku/lox/expr"
	"github.com/ArditZubaku/lox/token"
)

// TODO: Rethink the `any` or I could probably get rid of generics EVERYWHERE
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
	expression := p.term()

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		right := p.term()
		expression = expr.Binary[T]{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser[T]) term() expr.Expr[T] {
	expression := p.factor()

	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		right := p.factor()
		expression = expr.Binary[T]{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser[T]) factor() expr.Expr[T] {
	expression := p.unary()
	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right := p.unary()
		expression = expr.Binary[T]{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser[T]) unary() expr.Expr[T] {
	if p.match(token.Bang, token.Minus) {
		operator := p.previous()
		right := p.unary()
		return expr.Unary[T]{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser[T]) primary() expr.Expr[T] {
	if p.match(token.False) {
		return expr.Literal[T]{
			Value: false,
		}
	}

	if p.match(token.True) {
		return expr.Literal[T]{
			Value: true,
		}
	}

	if p.match(token.Nil) {
		return expr.Literal[T]{
			Value: nil,
		}
	}

	if p.match(token.Number, token.String) {
		return expr.Literal[T]{
			Value: p.previous().Literal,
		}
	}

	if p.match(token.LeftParen) {
		expression := p.expression()
		p.consume(token.RightParen, "Expect ')' after expression.")
		return expr.Grouping[T]{
			Expression: expression,
		}
	}

	return nil
}

func (p *Parser[T]) consume(paren token.Type, s string) {
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
