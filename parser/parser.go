package parser

import (
	"slices"

	"github.com/ArditZubaku/lox/expr"
	"github.com/ArditZubaku/lox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

func (p *Parser) equality() expr.Expr {
	expression := p.comparison()

	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser) match(types ...token.Type) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) comparison() expr.Expr {
	expression := p.term()

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		right := p.term()
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser) term() expr.Expr {
	expression := p.factor()

	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		right := p.factor()
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser) factor() expr.Expr {
	expression := p.unary()
	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right := p.unary()
		expression = &expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser) unary() expr.Expr {
	if p.match(token.Bang, token.Minus) {
		operator := p.previous()
		right := p.unary()
		return &expr.Unary{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() expr.Expr {
	if p.match(token.False) {
		return &expr.Literal{
			Value: false,
		}
	}

	if p.match(token.True) {
		return &expr.Literal{
			Value: true,
		}
	}

	if p.match(token.Nil) {
		return &expr.Literal{
			Value: nil,
		}
	}

	if p.match(token.Number, token.String) {
		return &expr.Literal{
			Value: p.previous().Literal,
		}
	}

	if p.match(token.LeftParen) {
		expression := p.expression()
		p.consume(token.RightParen, "Expect ')' after expression.")
		return &expr.Grouping{
			Expression: expression,
		}
	}

	return nil
}

func (p *Parser) consume(paren token.Type, s string) {
	panic("unimplemented")
}

func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}
