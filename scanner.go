package main

type Scanner struct {
	start   int
	current int
	line    int

	source string
	tokens []*Token
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []*Token{},
		line:    1,
		start:   0,
		current: 0,
	}
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(
		s.tokens,
		&Token{
			Type:    EOF,
			Lexeme:  "",
			Literal: nil,
			Line:    s.line,
		},
	)
}

func (s *Scanner) scanToken() {
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
