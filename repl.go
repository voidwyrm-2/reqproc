package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/voidwyrm-2/reqproc/runtime/interpreter"
)

func repl() error {
	scn := bufio.NewScanner(os.Stdin)
	acc := []string{}

	fmt.Println("ReqProc REPL")

	for {
		fmt.Print("> ")
		scn.Scan()
		acc = append(acc, scn.Text())

		interp, err := interpreter.New(nil)
		if err != nil {
			fmt.Println(err.Error())
		} else if _, err := interp.Execute(strings.Join(acc, "\n")); err != nil {
			if strings.HasPrefix(err.Error(), "EXIT CODE ") {
				fmt.Printf("exited with code '%s'\n", strings.Split(err.Error(), " ")[2])
				return nil
			}

			acc = acc[:len(acc)-1]
			fmt.Println(err.Error())
		} else if interp.StackLen() == 0 {
			fmt.Println("[]")
		} else {
			fmt.Println(interp.GetStack())
		}
	}
}
