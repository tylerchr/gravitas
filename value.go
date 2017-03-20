package gravitas

// #include <stdlib.h>
// #include "gravity_delegate.h"
// #include "gravity_utils.h"
import "C"
import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

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
