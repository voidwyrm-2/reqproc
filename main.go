package main

import (
	"fmt"

	"github.com/voidwyrm-2/reqproc/runtime/types/functiontype"
	"github.com/voidwyrm-2/reqproc/runtime/types/listtype"
	"github.com/voidwyrm-2/reqproc/runtime/types/numbertype"
	"github.com/voidwyrm-2/reqproc/runtime/types/stringtype"
)

func main() {
	f := functiontype.NewNative(func(str string, f func()) {})
	fmt.Println(f)

	str := stringtype.New("hello")
	fmt.Println(str)

	str2 := stringtype.New(" there")
	fmt.Println(str2)

	str3, _ := str.Add(str2)
	fmt.Println(str3)

	if _, err := str.Sub(str2); err != nil {
		fmt.Println(err.Error())
	}

	ls := listtype.New(stringtype.New("hello"), stringtype.New("there"), stringtype.New("good"), stringtype.New("sir"))
	fmt.Println(ls)

	als, _ := ls.Add(stringtype.New(" wow!"))
	fmt.Println(als)

	ls2 := listtype.New(numbertype.New(10), numbertype.New(20), numbertype.New(30), numbertype.New(40), numbertype.New(50))
	fmt.Println(ls2)

	mls, _ := ls2.Mul(numbertype.New(3))
	fmt.Println(mls)
}
