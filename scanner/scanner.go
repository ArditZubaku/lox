package scanner

import (
	"strconv"

	"github.com/ArditZubaku/lox/token"
)

type vm interface {
	ReportErr(line int, err error)
}

type Scanner struct {
	start   int
	current int
	line    int

	source []rune
	tokens []*token.Token

	vm vm
}

func NewScanner(vm vm, source string) *Scanner {
	return &Scanner{
		vm:      vm,
		source:  []rune(source),
		tokens:  []*token.Token{},
		line:    1,
		start:   0,
		current: 0,
	}
}

func (s *Scanner) ScanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(
		s.tokens,
		&token.Token{
			Type:    token.EOF,
			Lexeme:  "",
			Literal: nil,
			Line:    s.line,
		},
	)
}

func (s *Scanner) GetTokens() []*token.Token {
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LeftParen)
	case ')':
		s.addToken(token.RightParen)
	case '{':
		s.addToken(token.LeftBrace)
	case '}':
		s.addToken(token.RightBrace)
	case ',':
		s.addToken(token.Comma)
	case '.':
		s.addToken(token.Dot)
	case '-':
		s.addToken(token.Minus)
	case '+':
		s.addToken(token.Plus)
	case ';':
		s.addToken(token.Semicolon)
	case '*':
		s.addToken(token.Star)
	case '!':
		if s.match('=') {
			s.addToken(token.BangEqual)
		} else {
			s.addToken(token.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EqualEqual)
		} else {
			s.addToken(token.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LessEqual)
		} else {
			s.addToken(token.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GreaterEqual)
		} else {
			s.addToken(token.Greater)
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

		s.addToken(token.Slash)
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
			s.vm.ReportErr(s.line, ErrUnexpectedCharacter)
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

	s.vm.ReportErr(s.line, ErrUnterminatedComment)
}

func (s *Scanner) scanIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	txt := string(s.source[s.start:s.current])
	Type, ok := keywords[txt]
	if !ok {
		Type = token.Identifier
	}

	s.addToken(Type)
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
		s.vm.ReportErr(s.line, ErrInvalidNumberLiteral)
		return
	}

	s.addTokenWithLiteral(token.Number, num)
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
		s.vm.ReportErr(s.line, ErrUnterminatedString)
	}

	// The closing "
	s.advance()

	// Trim the surrounding quotes
	// value := string(s.source[s.start+1 : s.current-1])
	value := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(token.String, value)
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

func (s *Scanner) addToken(t token.Type) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t token.Type, literal any) {
	txt := string(s.source[s.start:s.current])
	s.tokens = append(
		s.tokens,
		&token.Token{
			Type:    t,
			Lexeme:  txt,
			Literal: literal,
			Line:    s.line,
		},
	)
}

func (s *Scanner) advance() rune {
	if s.current >= len(s.source) {
		return rune(token.EOF)
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
