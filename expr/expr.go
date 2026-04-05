package expr

import "github.com/ArditZubaku/lox/token"

// TODO: This should be constrained better and not be an `any`
// NOTE: I'd rather name this `Do` but I see the convention is to name it `Accept`
type Expr[T any] interface {
	Accept(Visitor[T]) T
}

type Visitor[T any] interface {
	VisitBinaryExpr(Binary[T]) T
	VisitGroupingExpr(Grouping[T]) T
	VisitLiteralExpr(Literal[T]) T
	VisitUnaryExpr(Unary[T]) T
}

type Binary[T any] struct {
	Left     Expr[T]
	Operator token.Token
	Right    Expr[T]
}

type Grouping[T any] struct {
	Expression Expr[T]
}

type Literal[T any] struct {
	Value any
}

type Unary[T any] struct {
	Operator token.Token
	Right    Expr[T]
}

func (b Binary[T]) Accept(v Visitor[T]) T {
	return v.VisitBinaryExpr(b)
}

func (l Literal[T]) Accept(v Visitor[T]) T {
	return v.VisitLiteralExpr(l)
}

func (g Grouping[T]) Accept(v Visitor[T]) T {
	return v.VisitGroupingExpr(g)
}

func (u Unary[T]) Accept(v Visitor[T]) T {
	return v.VisitUnaryExpr(u)
}
