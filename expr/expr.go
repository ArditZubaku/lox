package expr

import "github.com/ArditZubaku/lox/token"

type Value = any

type Expr interface {
	Accept(Visitor) Value
}

type Visitor interface {
	VisitBinaryExpr(*Binary) Value
	VisitGroupingExpr(*Grouping) Value
	VisitLiteralExpr(*Literal) Value
	VisitUnaryExpr(*Unary) Value
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value Value
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) Value {
	return v.VisitBinaryExpr(b)
}

func (l *Literal) Accept(v Visitor) Value {
	return v.VisitLiteralExpr(l)
}

func (g *Grouping) Accept(v Visitor) Value {
	return v.VisitGroupingExpr(g)
}

func (u *Unary) Accept(v Visitor) Value {
	return v.VisitUnaryExpr(u)
}
