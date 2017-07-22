package main

import (
	"fmt"
	"os"

	"github.com/elemental-vm/test-vm/lexer"
	"github.com/elemental-vm/test-vm/vm"
)

func main() {
	filename := os.Args[1]

	lexer, err := lexer.New(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	program, err := lexer.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// fmt.Printf("%#v\n", program)
	// fmt.Println(len(program))

	newvm := vm.New(program)
	os.Exit(int(newvm.Start()))
}
