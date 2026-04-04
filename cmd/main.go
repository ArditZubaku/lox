package main

import (
	"fmt"
	"os"

	"github.com/ArditZubaku/lox/lox"
)

func main() {
	vm := lox.Lox{}
	if len(os.Args) > 2 {
	fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := vm.RunFile(os.Args[1]); err != nil {
			// TODO: Rethink this
			os.Exit(1)
		}
	} else {
		vm.RunPrompt()
	}
}
