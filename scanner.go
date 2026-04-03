package main

type Scanner struct {
	start   int
	current int
	line    int

	source []rune
	tokens []*Token
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  []rune(source),
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
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	default:
		vm.err(s.line, "Unexpected character.")
	}

}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++

	return true
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	txt := string(s.source[s.start:s.current])
	s.tokens = append(
		s.tokens,
		&Token{
			Type:    t,
			Lexeme:  txt,
			Literal: literal,
			Line:    s.line,
		},
	)
}

func (s *Scanner) advance() rune {
	if s.current >= len(s.source) {
		return rune(EOF)
	}
	ch := s.source[s.current]
	s.current++
	return ch
	// s.current++
	// return s.source[s.current]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
