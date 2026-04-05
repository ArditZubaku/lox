package parser

import (
	"slices"

	"github.com/ArditZubaku/lox/expr"
	"github.com/ArditZubaku/lox/token"
)

type vm interface {
	ReportParseError(err ParseError)
}

type Parser struct {
	tokens  []token.Token
	current int
	vm      vm
}

func New(vm vm, tokens []token.Token) Parser {
	return Parser{
		vm:      vm,
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (expr.Expr, error) {
	return p.expression()
}

// TODO: If this doesn't get extended with more functionality later let's just rename this to Parse()
func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (expr.Expr, error) {
	return p.parseBinary(p.comparison, token.BangEqual, token.EqualEqual)
}

func (p *Parser) comparison() (expr.Expr, error) {
	return p.parseBinary(
		p.term,
		token.Greater,
		token.GreaterEqual,
		token.Less,
		token.LessEqual,
	)
}

func (p *Parser) term() (expr.Expr, error) {
	return p.parseBinary(p.factor, token.Minus, token.Plus)
}

func (p *Parser) factor() (expr.Expr, error) {
	return p.parseBinary(p.unary, token.Slash, token.Star)
}

func (p *Parser) unary() (expr.Expr, error) {
	if p.match(token.Bang, token.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &expr.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (expr.Expr, error) {
	if p.match(token.False) {
		return &expr.Literal{
			Value: false,
		}, nil
	}

	if p.match(token.True) {
		return &expr.Literal{
			Value: true,
		}, nil
	}

	if p.match(token.Nil) {
		return &expr.Literal{
			Value: nil,
		}, nil
	}

	if p.match(token.Number, token.String) {
		return &expr.Literal{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(token.LeftParen) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &expr.Grouping{
			Expression: expression,
		}, nil
	}

	return nil, ParseError{Token: p.peek(), Msg: "Expect expression."}
}

func (p *Parser) parseBinary(
	parseOperand func() (expr.Expr, error),
	operators ...token.Type,
) (expr.Expr, error) {
	left, err := parseOperand()
	if err != nil {
		return nil, err
	}

	for p.match(operators...) {
		operator := p.previous()
		right, err := parseOperand()
		if err != nil {
			return nil, err
		}

		left = &expr.Binary{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) match(types ...token.Type) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) consume(t token.Type, msg string) (token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return token.Token{}, ParseError{Token: p.peek(), Msg: msg}
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

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.Semicolon {
			return
		}

		switch p.peek().Type {
		case
			token.Class,
			token.Fun,
			token.Var,
			token.For,
			token.If,
			token.While,
			token.Print,
			token.Return:
			return
		}

		p.advance()
	}
}
