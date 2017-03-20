package gravitas

/*
#include "gravity_vm.h"

void log_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
void error_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"unsafe"
)

// VM interfaces with the Gravity VM.
type VM struct {
	cGravityVM  *C.gravity_vm
	delegatePtr unsafe.Pointer
}

// NewVM instantiates a new Gravity VM using the provided Delegate.
func NewVM(d Delegate) (*VM, error) {

	vm := &VM{}

	delegateID := registerDelegate(d)

	// write the delegate ID into a C array; this will be our xdata
	var delID [8]byte
	binary.LittleEndian.PutUint64(delID[:], delegateID)
	vm.delegatePtr = C.CBytes(delID[:])

	delegate := C.gravity_delegate_t{
		xdata:          vm.delegatePtr,
		log_callback:   C.gravity_log_callback(C.log_trampoline),
		error_callback: C.gravity_error_callback(C.error_trampoline),
	}

	vm.cGravityVM = C.gravity_vm_new(&delegate)

	return vm, nil
}

// RunMain executes the main method of a Closure within the context of vm.
func (vm *VM) RunMain(closure *Closure) (interface{}, error) {

	if ok := C.gravity_vm_runmain(vm.cGravityVM, closure.cGravityClosure); !ok {
		return nil, errors.New("runtime error occurred")
	}

	result := C.gravity_vm_result(vm.cGravityVM)
	return convertGravityValue(result)

}

// Close destroys this VM and frees associated resources.
func (vm *VM) Close() error {
	C.gravity_vm_free(vm.cGravityVM)
	C.free(vm.delegatePtr)
	return nil
}
