package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/voidwyrm-2/reqproc/lexer"
	"github.com/voidwyrm-2/reqproc/runtime/interpreter"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
)

func _main() error {
	if len(os.Args) < 2 {
		return errors.New("expected 'reqproc <file>'")
	}

	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	l := lexer.New(string(content))

	tokens, err := l.Lex()
	if err != nil {
		return err
	}

	interp := interpreter.New(scope.New(nil))

	leftover, err := interp.ExecuteTokens(tokens)

	fmt.Println(leftover)

	return err
}

func main() {
	if err := _main(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
