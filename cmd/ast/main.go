package main

import (
	"fmt"

	"github.com/ArditZubaku/lox/expr"
	"github.com/ArditZubaku/lox/token"
)

func main() {
	expression := &expr.Binary{
		Left: &expr.Unary{
			Operator: token.NewToken(token.Minus, "-", nil, 1),
			Right: &expr.Literal{
				Value: 123,
			},
		},
		Operator: token.NewToken(token.Star, "*", nil, 1),
		Right: &expr.Grouping{
			Expression: &expr.Literal{
				Value: 45.67,
			},
		},
	}

	printer := expr.AstPrinter{}
	fmt.Println(printer.Print(expression))
}
