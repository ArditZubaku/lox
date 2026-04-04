package main

import "strconv"

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
	case '/':
		// If you encounter two slashes in a row, consume a comment until the end of the line
		if s.match('/') {
			s.consumeLineComment()
			break
		}

		if s.match('*') {
			s.consumeMultiLineComment()
			break
		}

		s.addToken(Slash)
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		if s.isDigit(c) {
			s.scanNumber()
		} else if s.isAlpha(c) {
			s.scanIdentifier()
		} else {
			vm.err(s.line, ErrUnexpectedCharacter)
		}
	}
}

func (s *Scanner) consumeLineComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *Scanner) consumeMultiLineComment() {
	for !s.isAtEnd() {
		switch c := s.peek(); c {
		case '\n':
			s.line++
			s.advance()
		case '*':
			if s.peekNext() == '/' {
				s.current += 2
				return
			}
		default:
			s.advance()
		}
	}

	vm.err(s.line, ErrUnterminatedComment)
}

func (s *Scanner) scanIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	txt := string(s.source[s.start:s.current])
	tokenType, ok := keywords[txt]
	if !ok {
		tokenType = Identifier
	}

	s.addToken(tokenType)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) scanNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the '.'
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	txt := string(s.source[s.start:s.current])
	num, err := strconv.ParseFloat(txt, 64)
	if err != nil {
		vm.err(s.line, ErrInvalidNumberLiteral)
		return
	}

	s.addTokenWithLiteral(Number, num)
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		vm.err(s.line, ErrUnterminatedString)
	}

	// The closing "
	s.advance()

	// Trim the surrounding quotes
	// value := string(s.source[s.start+1 : s.current-1])
	value := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(String, value)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
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
