package main

import (
	"fmt"

	"github.com/ArditZubaku/lox/expr"
	"github.com/ArditZubaku/lox/token"
)

func main() {
	expression := expr.Binary[string]{
		Left: expr.Unary[string]{
			Operator: token.NewToken(token.Minus, "-", nil, 1),
			Right: expr.Literal[string]{
				Value: 123,
			},
		},
		Operator: token.NewToken(token.Star, "*", nil, 1),
		Right: expr.Grouping[string]{
			Expression: expr.Literal[string]{
				Value: 45.67,
			},
		},
	}

	printer := expr.AstPrinter{}
	fmt.Println(printer.Print(expression))
}
