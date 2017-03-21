package gravitas

/*
#include "gravity_vm.h"

void log_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
void error_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
const char* precode_trampoline(void *xdata);
const char* loadfile_trampoline(const char *file, size_t *size, uint32_t *fileid, void *xdata);
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"time"
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
		xdata:             vm.delegatePtr,
		log_callback:      C.gravity_log_callback(C.log_trampoline),
		error_callback:    C.gravity_error_callback(C.error_trampoline),
		precode_callback:  C.gravity_precode_callback(C.precode_trampoline),
		loadfile_callback: C.gravity_loadfile_callback(C.loadfile_trampoline),
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

// Time indicates the runtime of the main method of the most recently run closure.
func (vm *VM) Time() time.Duration {
	ms := C.gravity_vm_time(vm.cGravityVM)
	return time.Duration(float64(ms) * float64(time.Millisecond))
}

// Close destroys this VM and frees associated resources.
func (vm *VM) Close() error {
	C.gravity_vm_free(vm.cGravityVM)
	C.free(vm.delegatePtr)
	return nil
}
