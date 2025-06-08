package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/voidwyrm-2/reqproc/lexer"
	"github.com/voidwyrm-2/reqproc/runtime"
	"github.com/voidwyrm-2/reqproc/runtime/interpreter"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/types"
	// "github.com/voidwyrm-2/vcheck"
	// "github.com/voidwyrm-2/vcheck/version"
)

//go:embed version.txt
var localVersionString string

const versionURL = "https://github.com/voidwyrm-2/reqproc/blob/main/example/version.txt"

func _main() error {
	fpath := flag.String("f", "", "The file to interpret")
	showVersion := flag.Bool("v", false, "Prints the interpreter version and exits")
	showTokens := flag.Bool("t", false, "Show the generated tokens")
	runREPL := flag.Bool("repl", false, "Run the repl instead of a file")
	// noVersionChecks := flag.Bool("nvc", false, "Do not check for a newer version; useful if internet is not available")

	flag.Parse()

	/*
		if !*noVersionChecks {
			localVersion, err := version.FromString(localVersionString)
			if err != nil {
				return err
			}

			remoteVersion, newerExists, err := vcheck.NewerVersionExists(versionURL)
			if err != nil {
				return err
			}

			if newerExists {
				fmt.Printf("A newer ReqProc version exists (%v -> %v)\nPlease update with 'go get github.com/voidwyrm-2/reqproc@latest' or by ", localVersion, remoteVersion)
				fmt.Printf("downloading the latest release from https://github.com/voidwyrm-2/reqproc/releases/tag/v%v\n", remoteVersion)
				return nil
			}
		}
	*/

	if *showVersion {
		fmt.Println("ReqProc interpreter, version", runtime.REQPROC_VERSION)
	}

	if *runREPL {
		return repl()
	}

	if len(os.Args) < 2 {
		return errors.New("expected 'reqproc <file>'")
	}

	file, err := os.Open(*fpath)
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

	if *showTokens {
		for _, t := range tokens {
			fmt.Println(t)
		}
		fmt.Println()
	}

	interp, err := interpreter.New(scope.New(nil, map[string]types.ReqType{}))
	if err != nil {
		return err
	}

	_, err = interp.ExecuteTokens(tokens)

	return err
}

func main() {
	if err := _main(); err != nil {
		runtime.HandleExitError(err)

		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
