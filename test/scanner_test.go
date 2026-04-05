package scanner

import (
	"testing"
)

type Token struct {
	Type   int
	Lexeme string
	Line   int
}

// Value version
func buildValueTokens(n int) []Token {
	tokens := make([]Token, 0, n)
	for i := range n {
		tokens = append(tokens, Token{
			Type:   i,
			Lexeme: "abc",
			Line:   1,
		})
	}
	return tokens
}

// Pointer version
func buildPointerTokens(n int) []*Token {
	tokens := make([]*Token, 0, n)
	for i := range n {
		tokens = append(tokens, &Token{
			Type:   i,
			Lexeme: "abc",
			Line:   1,
		})
	}
	return tokens
}

func BenchmarkValueTokens(b *testing.B) {
	for b.Loop() {
		_ = buildValueTokens(1_000_000)
	}
}

func BenchmarkPointerTokens(b *testing.B) {
	for b.Loop() {
		_ = buildPointerTokens(1_000_000)
	}
}
