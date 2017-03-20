package gravitas

/*
#include <stdlib.h>
#include "gravity_compiler.h"

void log_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
void error_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"unsafe"
)

type (
	Compiler struct {
		cGravityCompiler *C.gravity_compiler_t
		delegatePtr      unsafe.Pointer
	}

	Closure struct {
		cGravityClosure *C.gravity_closure_t
	}
)

func NewCompiler(d Delegate) (*Compiler, error) {

	compiler := &Compiler{}

	delegateID := registerDelegate(d)

	// write the delegate ID into a C array; this will be our xdata
	var delID [8]byte
	binary.LittleEndian.PutUint64(delID[:], delegateID)
	compiler.delegatePtr = C.CBytes(delID[:])

	delegate := C.gravity_delegate_t{
		xdata:          compiler.delegatePtr,
		log_callback:   C.gravity_log_callback(C.log_trampoline),
		error_callback: C.gravity_error_callback(C.error_trampoline),
	}

	compiler.cGravityCompiler = C.gravity_compiler_create(&delegate)

	return compiler, nil
}

// Compile converts source code into a Closure.
func (c *Compiler) Compile(source []byte) (*Closure, error) {

	// compile source code into a closure
	c_source := C.CString(string(source))
	main_closure := C.gravity_compiler_run(c.cGravityCompiler, c_source, C.size_t(len(source)), 0, false)
	C.free(unsafe.Pointer(c_source))

	if main_closure == nil {
		return nil, errors.New("compiler error occurred")
	}

	return &Closure{cGravityClosure: main_closure}, nil

}

// Transfer copies the sum of compiled symbols into the VM's memory.
func (c *Compiler) Transfer(vm *VM) error {
	C.gravity_compiler_transfer(c.cGravityCompiler, vm.cGravityVM)
	return nil
}

// Close frees resources associated with Compiler.
func (c *Compiler) Close() error {
	// C.gravity_compiler_free(c.cGravityCompiler)
	C.free(c.delegatePtr)
	return nil
}
