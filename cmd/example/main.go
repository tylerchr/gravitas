package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tylerchr/gravitas"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("USAGE: %s path/to/program.gravity")
		os.Exit(1)
	}

	// load a program from disk
	source, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// create a Go delegate
	var delegate SampleDelegate

	// create a VM
	vm, _ := gravitas.NewVM(delegate)
	defer vm.Close()
	// fmt.Printf("%#v\n", vm)

	// set up a compiler
	compiler, _ := gravitas.NewCompiler(delegate)
	defer compiler.Close()
	// fmt.Printf("%#v\n", compiler)

	// compile the program
	closure, _ := compiler.Compile(source)
	// fmt.Printf("%#v\n", closure)

	// copy all the symbols from the compiler into the VM
	compiler.Transfer(vm)

	result, err := vm.Run(closure)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T: %#v\n", result, result)

}

type SampleDelegate struct{}

func (md SampleDelegate) Log(message string) {
	log.Println(message)
}

func (md SampleDelegate) Error(errType int, message string, errDesc gravitas.ErrorDescription) {
	log.Printf("error [%d]: %s (%#v)\n", errType, message, errDesc)
}
