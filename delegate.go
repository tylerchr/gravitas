package gravitas

/*
#cgo CFLAGS: -I${SRCDIR}/gravity/src/compiler -I${SRCDIR}/gravity/src/runtime -I${SRCDIR}/gravity/src/shared -I${SRCDIR}/gravity/src/utils
#cgo LDFLAGS: -L${SRCDIR}/gravity -lgravity
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

type (
	// Delegate implements callback methods that may be invoked by various Gravity
	// components to provide state about various lifecycle events.
	//
	// Delegate aligns with the gravity_delegate_t type in C.
	Delegate interface {
		Log(message string)
		Error(errType int, description string, errDesc ErrorDescription)
		Precode() []byte
		Loadfile(file string) (source []byte, fileID uint32)
	}

	// ErrorDescription describes the location in source code of any given error.
	//
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

// registerDelegate allocates a unique ID for the given delegate and stores a
// reference in an internal map. When the unique ID is used as the `xdata` in
// the Gravity library, this enables delegate calls in the library to route to
// the correct Go delegate implementation.
//
// If the delegate is already registered, its preexisting ID is returned.
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

// unregisterDelegate removes the provided delegate from the internal registry.
//
// TODO(tylerc): We need to either implement a reference-counting-like solution
// or stop deduplicating Delegate registrations; otherwise a shared Delegate could
// be prematurely unregistered by a client unaware of other users.
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

// lookupDelegate finds a registered Delegate by its unique ID.
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

	// lookup the correct the delegate
	myID := *(*uint64)(xdata)
	d, err := lookupDelegate(myID)
	if err != nil {
		panic(err)
	}

	// invoke the callback method
	d.Log(C.GoString(message))

}

//export go_err_callback
func go_err_callback(message *C.char, xdata unsafe.Pointer) {

	// lookup the correct the delegate
	myID := *(*uint64)(xdata)
	d, err := lookupDelegate(myID)
	if err != nil {
		panic(err)
	}

	// invoke the callback method
	d.Error(0, C.GoString(message), ErrorDescription{})
}

//export go_loadfile_callback
func go_loadfile_callback(file *C.char, size *uint32, fileid *uint32, xdata unsafe.Pointer) *C.char {

	// lookup the correct the delegate
	myID := *(*uint64)(xdata)
	d, err := lookupDelegate(myID)
	if err != nil {
		panic(err)
	}

	// invoke the callback method
	source, fileID := d.Loadfile(C.GoString(file))

	// populate the size and fileid fields
	*size = uint32(len(source))
	*fileid = fileID

	return C.CString(string(source)) // Gravity will free this for us in gravity_lexer_free

}

//export go_precode_callback
func go_precode_callback(xdata unsafe.Pointer) *C.char {

	// lookup the correct the delegate
	myID := *(*uint64)(xdata)
	d, err := lookupDelegate(myID)
	if err != nil {
		panic(err)
	}

	// invoke the callback method
	source := d.Precode()

	return C.CString(string(source)) // Gravity will free this for us in gravity_lexer_free

}
