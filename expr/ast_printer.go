package expr

import (
	"fmt"
	"strconv"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) Value {
	return expr.Accept(p)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) Value {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) Value {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) Value {
	if expr.Value == nil {
		return "nil"
	}

	switch expr.Value.(type) {
	case float64:
		return strconv.FormatFloat(expr.Value.(float64), 'g', -1, 64)
	case string:
		return "\"" + expr.Value.(string) + "\""
	default:
		return fmt.Sprintf("%v", expr.Value)
	}
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) Value {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) Value {
	var builder strings.Builder

	builder.WriteRune('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteRune(' ')
		s, ok := expr.Accept(p).(string)
		if !ok {
			s = fmt.Sprintf("%v", expr.Accept(p))
		}
		builder.WriteString(s)
	}
	builder.WriteRune(')')
	return builder.String()
}
