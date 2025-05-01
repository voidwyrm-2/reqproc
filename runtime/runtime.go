package runtime

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const REQPROC_VERSION = 3.0

func HandleExitPanic(e error) {
	if e != nil {
		if strings.HasPrefix(e.Error(), "EXIT CODE ") {
			spl := strings.Split(e.Error(), " ")
			ec, err := strconv.Atoi(spl[len(spl)-1])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			os.Exit(ec)
		}
	}
}
