package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

//vm as in virtual machine
var vm = Lox{}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := vm.runFile(os.Args[1]); err != nil {
			// TODO: Rethink this
			os.Exit(1)
		}
	} else {
		vm.runPrompt()
	}
}

func (l *Lox) runFile(path string) error {
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
	scanner := NewScanner(string(source))
	scanner.scanTokens()

	for _, token := range scanner.tokens {
		fmt.Println(token)
	}
}

func (l *Lox) runPrompt() {
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

func (l *Lox) err(line int, msg string) {
	l.report(line, "", msg)
}

func (l *Lox) report(line int, where string, msg string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, msg)
	l.hadError = true
}
