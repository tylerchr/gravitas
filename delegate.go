package gravitas

/*
#cgo CFLAGS: -I${SRCDIR}/gravity/src/compiler -I${SRCDIR}/gravity/src/runtime -I${SRCDIR}/gravity/src/shared -I${SRCDIR}/gravity/src/utils
#cgo LDFLAGS: -L${SRCDIR}/gravity -lgravity

#include <stdlib.h>
#include "gravity_delegate.h"
#include "gravity_utils.h"

void log_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
void error_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata);
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

type (
	// Delegate aligns with the gravity_delegate_t type in C.
	Delegate interface {
		Log(message string)
		Error(errType int, description string, errDesc ErrorDescription)
	}

	// ErrorDescription aligns with the error_desc_t type in C.
	ErrorDescription struct {
		Code   uint32
		Line   uint32
		Column uint32
		FileID uint32
		Offset uint32
	}
)

var (
	delegatesMu sync.RWMutex
	delegateIdx map[Delegate]uint64
	delegates   map[uint64]Delegate
	delegateID  uint64
)

func registerDelegate(d Delegate) uint64 {

	delegatesMu.Lock()
	defer delegatesMu.Unlock()

	if delegateIdx == nil {
		delegateIdx = make(map[Delegate]uint64)
	}
	if delegates == nil {
		delegates = make(map[uint64]Delegate)
	}

	// reuse existing ID if the delegate is already registered
	if id, ok := delegateIdx[d]; ok {
		return id
	}

	myID := delegateID
	delegateID++
	delegateIdx[d] = myID
	delegates[myID] = d

	return myID
}

func unregisterDelegate(d Delegate) error {

	delegatesMu.Lock()
	defer delegatesMu.Unlock()

	id, ok := delegateIdx[d]
	if !ok {
		return errors.New("delegate not registered")
	}
	delete(delegateIdx, d)
	delete(delegates, id)

	return nil
}

func lookupDelegate(id uint64) (Delegate, error) {
	delegatesMu.RLock()
	defer delegatesMu.RUnlock()
	d, ok := delegates[id]
	if !ok {
		return nil, errors.New("delegate not registered")
	}
	return d, nil
}

//export go_log_callback
func go_log_callback(message *C.char, xdata unsafe.Pointer) {
	// fmt.Printf("LOG: %s (%v)\n", C.GoString(message), xdata)

	// lookup the correct the delegate
	myID := uint64(uintptr(xdata))
	d, err := lookupDelegate(myID)
	if err != nil {
		panic(err)
	}

	// invoke the callback method
	d.Log(C.GoString(message))

}

//export go_err_callback
func go_err_callback(message *C.char, xdata unsafe.Pointer) {
	// fmt.Printf("ERR: %s (%v)\n", C.GoString(message), xdata)

	// lookup the correct the delegate
	myID := uint64(uintptr(xdata))
	d, err := lookupDelegate(myID)
	if err != nil {
		panic(err)
	}

	// invoke the callback method
	d.Error(0, C.GoString(message), ErrorDescription{})
}
