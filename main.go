package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

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

	_, err = interp.ExecuteTokens(tokens)

	return err
}

func main() {
	exm, err := regexp.Compile(`error on line [0-9]+, col [0-9]+: [0-9]+ EXIT`)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := _main(); err != nil {
		if exm.MatchString(err.Error()) {
			spl := strings.Split(err.Error(), " ")
			ec, err := strconv.Atoi(spl[len(spl)-2])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			os.Exit(ec)
		}

		fmt.Println(err.Error())
		os.Exit(1)
	}
}
