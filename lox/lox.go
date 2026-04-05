package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ArditZubaku/lox/parser"
	"github.com/ArditZubaku/lox/scanner"
	"github.com/ArditZubaku/lox/token"
)

type Lox struct {
	hadError bool
}

func (l *Lox) RunPrompt() {
	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if ok := input.Scan(); !ok {
			break
		}

		line := input.Bytes()
		if string(line) == ".exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		l.run(line)
		l.hadError = false
	}
}

func (l *Lox) RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	l.run(bytes)

	if l.hadError {
		os.Exit(65)
	}

	return nil
}

func (l *Lox) run(source []byte) {
	sc := scanner.NewScanner(l, string(source))
	sc.ScanTokens()

	for _, token := range sc.GetTokens() {
		fmt.Println(token)
	}
}

func (l *Lox) ReportErr(line int, err error) {
	l.report(line, "", err)
}

func (l *Lox) ReportParseError(err parser.ParseError) {
	if err.Token.Type == token.EOF {
		l.report(err.Token.Line, " at end", err)
	} else {
		l.report(err.Token.Line, " at '"+err.Token.Lexeme+"'", err)
	}
}

func (l *Lox) report(line int, where string, err error) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	l.hadError = true
}
