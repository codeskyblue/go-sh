package sh_test

import (
	"fmt"

	"github.com/shxsun/go-sh"
)

func ExampleCommand() {
	out, err := sh.Command("echo", "hello").Run()
	fmt.Println(string(out), err)
}

func ExampleTest() {
	if sh.Test("dir", "mydir") {
		fmt.Println("mydir exists")
	}
}
