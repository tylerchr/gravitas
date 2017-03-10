package main

/*
#cgo CFLAGS: -I${SRCDIR}/gravity/src/compiler -I${SRCDIR}/gravity/src/runtime -I${SRCDIR}/gravity/src/shared -I${SRCDIR}/gravity/src/utils
#cgo LDFLAGS: -L ${SRCDIR}/gravity/src -lgravity

#include <stdlib.h>
#include "gravity_compiler.h"
#include "gravity_delegate.h"
#include "gravity_utils.h"
#include "gravity_vm.h"

void log_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
void error_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"unsafe"
)

func main() {

	file := flag.String("program", "", "gravity program to execute")

	flag.Parse()

	if *file == "" {
		fmt.Println("must provide a -file argument")
		os.Exit(1)
	}

	// here we load a program from disk
	// var size C.size_t
	// c_source := C.file_read(C.CString(*file), &size)
	// source := C.GoString(c_source)
	// fmt.Printf("%T %s %d %d\n", source, source, len(source), size)

	source, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	// create the Gravity delegate
	// we'll use this all over
	delegate := C.gravity_delegate_t{
		log_callback:   C.gravity_log_callback(C.log_trampoline),
		error_callback: C.gravity_error_callback(C.error_trampoline),
	}
	// fmt.Printf("Delegate: %T %#v\n", delegate, delegate)

	// create the Gravity VM
	vm := C.gravity_vm_new(&delegate)
	// fmt.Printf("VM: %T %#v\n", vm, vm)

	// create a Gravity compiler
	compiler := C.gravity_compiler_create(&delegate)
	// fmt.Printf("Compiler: %T %#v\n", compiler, compiler)

	// compile source code into a closure
	c_source := C.CString(string(source))
	main_closure := C.gravity_compiler_run(compiler, c_source, C.size_t(len(source)), 0, false)
	C.free(unsafe.Pointer(c_source))
	if main_closure == nil {
		os.Exit(1)
	}

	// fmt.Printf("null: %t\n", main_closure == nil)
	// fmt.Printf("Closure: %T %#v\n", main_closure, main_closure)

	C.gravity_compiler_transfer(compiler, vm)

	if C.gravity_vm_run(vm, main_closure) {

		result := C.gravity_vm_result(vm)

		gores, err := convertGravityValue(result)
		fmt.Printf("%T %#v %v\n", gores, gores, err)

		// var buffer [512]byte
		// buf := (*C.char)(unsafe.Pointer(&buffer))
		// C.gravity_value_dump(result, buf, C.uint16_t(len(buffer)))
		// fmt.Printf("RESULT: %s\n", C.GoString(buf))

	}

}

// convertGravityValue returns a copy of the gravity_value_t as a Go object.
func convertGravityValue(v C.gravity_value_t) (interface{}, error) {

	// fmt.Printf("value: %#v\n", v)

	vv := [8]byte(v.anon0)
	vvv := binary.LittleEndian.Uint64(vv[:])

	switch v.isa {

	case C.gravity_class_function:
		fmt.Println("it's a function")
		return nil, errors.New("cannot convert Gravity type: function")

	case C.gravity_class_closure:
		fmt.Println("it's a closure")
		return nil, errors.New("cannot convert Gravity type: closure")

	case C.gravity_class_null:
		if int(vvv) == 0 {
			// it's "null"
		} else if int(vvv) == 1 {
			// it's "undefined"
		}
		return nil, nil

	case C.gravity_class_string:
		s := (*C.gravity_string_t)(unsafe.Pointer(uintptr(vvv)))
		return C.GoString(s.s), nil

	case C.gravity_class_int:
		return int(vvv), nil

	case C.gravity_class_float:
		return math.Float64frombits(vvv), nil

	case C.gravity_class_bool:
		return int(vvv) == 1, nil

	}

	return nil, fmt.Errorf("unknown type: %T\n", v.isa)

}

//export go_log_callback
func go_log_callback(message *C.char, xdata unsafe.Pointer) {
	fmt.Printf("LOG: %s (%v)\n", C.GoString(message), xdata)
}

//export go_err_callback
func go_err_callback(message *C.char, xdata unsafe.Pointer) {
	fmt.Printf("ERR: %s (%v)\n", C.GoString(message), xdata)
}
