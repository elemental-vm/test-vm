package main

import (
	"fmt"
	"os"

	"flag"

	"github.com/elemental-vm/test-vm/lexer"
	"github.com/elemental-vm/test-vm/vm"
)

var (
	debug   bool
	compile bool
	outFile string
)

func init() {
	flag.BoolVar(&debug, "d", false, "Enable debug output")
	flag.BoolVar(&compile, "c", false, "Compile to byte file")
	flag.StringVar(&outFile, "o", "", "Output file")
}

func main() {
	flag.Parse()

	filename := flag.Arg(0)

	theLexer, err := lexer.New(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	program, err := theLexer.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if compile {
		file, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		file.Write(lexer.FileHeader)
		file.Write(program)
		file.Close()
		return
	}

	if debug {
		fmt.Printf("%#v\n", program)
		fmt.Println(len(program))
	}

	newvm := vm.New(program)
	os.Exit(int(newvm.Start(debug)))
}
