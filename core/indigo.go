// Package core provides core functions for Indigo C API library via CGO
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/12
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : indigo.go
// @Software: GoLand
package core

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// Windows: link against import libraries (.lib)
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo

// Linux: use $ORIGIN for runtime library search
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-aarch64

// macOS: use @loader_path (not @executable_path) for shared libraries
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/darwin-x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/darwin-aarch64
#include <stdlib.h>
#include "indigo.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

//// indigoSessionID holds the session ID for Indigo
//var indigoSessionID C.qword

// Indigo represents a session-bound handle to Indigo C library.
type Indigo struct {
	sid uint64
}

// IndigoObject is a lightweight wrapper around Indigo object handle.
type IndigoObject struct {
	id     int
	client *Indigo
}

// IndigoInit creates a new Indigo instance and allocates a session id.
func IndigoInit() (*Indigo, error) {
	sid := C.indigoAllocSessionId()
	if sid == 0 {
		// try to read last error if available
		if errStr := lastErrorString(); errStr != "" {
			return nil, errors.New(errStr)
		}
		return nil, fmt.Errorf("indigo: failed to alloc session id, got %v", sid)
	}
	C.indigoSetSessionId(sid)
	return &Indigo{sid: uint64(sid)}, nil
}

// Close releases session id; call when done with the Indigo instance.
func (in *Indigo) Close() {
	if in == nil {
		return
	}
	if in.sid != 0 {
		in.setSession()
		C.indigoReleaseSessionId(C.ulonglong(in.sid))
		in.sid = 0
	}
}

func (in *Indigo) GetSessionID() uint64 {
	return in.sid
}

// setSession sets the internal session id on the native library for next calls.
func (in *Indigo) setSession() {
	// wrap call to C to set session id for this goroutine call
	C.indigoSetSessionId(C.ulonglong(in.sid))
}

// helper to read last error string from Indigo C API
func lastErrorString() string {
	ptr := C.indigoGetLastError()
	if ptr == nil {
		return ""
	}
	return C.GoString(ptr)
}

// checkResultInt checks integer return values (assume 0 means error in this API).
// Adjust condition depending on actual Indigo C API contract.
func checkResultInt(res C.int) (int, error) {
	if res == 0 {
		return 0, errors.New(lastErrorString())
	}
	return int(res), nil
}

// checkResultInt64 for functions returning long/int handles
func checkResultLong(res C.long) (int, error) {
	if res == 0 {
		return 0, errors.New(lastErrorString())
	}
	return int(res), nil
}

// ---------- Example methods ----------

// Deserialize creates molecule/reaction object from binary serialized CMF format.
func (in *Indigo) Deserialize(arr []byte) (*IndigoObject, error) {
	in.setSession()
	if len(arr) == 0 {
		return nil, errors.New("indigo: empty buffer")
	}
	cbuf := C.CBytes(arr)
	defer C.free(cbuf)
	res := C.indigoUnserialize((*C.uchar)(cbuf), C.int(len(arr)))
	id, err := checkResultLong(C.long(res))
	if err != nil {
		return nil, err
	}
	return &IndigoObject{id: id, client: in}, nil
}

// SetOption sets option value. This mirrors the Python setOption behavior.
func (in *Indigo) SetOption(option string, v1 interface{}, v2 interface{}, v3 interface{}) error {
	in.setSession()
	copt := C.CString(option)
	defer C.free(unsafe.Pointer(copt))

	// Three floats -> color
	if f1, ok1 := v1.(float64); ok1 {
		if f2, ok2 := v2.(float64); ok2 {
			if f3, ok3 := v3.(float64); ok3 {
				if C.indigoSetOptionColor(copt, C.float(f1), C.float(f2), C.float(f3)) == 0 {
					return errors.New(lastErrorString())
				}
				return nil
			}
		}
	}

	// Two ints -> XY
	if i1, ok1 := v1.(int); ok1 {
		if i2, ok2 := v2.(int); ok2 && v3 == nil {
			if C.indigoSetOptionXY(copt, C.int(i1), C.int(i2)) == 0 {
				return errors.New(lastErrorString())
			}
			return nil
		}
	}

	// Single value types
	if v2 == nil && v3 == nil {
		switch val := v1.(type) {
		case string:
			cval := C.CString(val)
			defer C.free(unsafe.Pointer(cval))
			if C.indigoSetOption(copt, cval) == 0 {
				return errors.New(lastErrorString())
			}
			return nil
		case int:
			if C.indigoSetOptionInt(copt, C.int(val)) == 0 {
				return errors.New(lastErrorString())
			}
			return nil
		case float64:
			if C.indigoSetOptionFloat(copt, C.float(val)) == 0 {
				return errors.New(lastErrorString())
			}
			return nil
		case bool:
			var b C.int = 0
			if val {
				b = 1
			}
			if C.indigoSetOptionBool(copt, b) == 0 {
				return errors.New(lastErrorString())
			}
			return nil
		default:
			return errors.New("indigo: bad option value type")
		}
	}

	return errors.New("indigo: bad option parameter combination")
}

// CreateArray creates an Indigo array for rendering multiple objects
func (in *Indigo) CreateArray() (int, error) {
	handle := int(C.indigoCreateArray())
	if handle < 0 {
		return 0, fmt.Errorf("failed to create array: %s", lastErrorString())
	}
	return handle, nil
}

// ArrayAdd adds an object to an array
func (in *Indigo) ArrayAdd(arrayHandle int, objectHandle int) error {
	if arrayHandle < 0 {
		return fmt.Errorf("invalid array handle")
	}
	if objectHandle < 0 {
		return fmt.Errorf("invalid object handle")
	}

	ret := int(C.indigoArrayAdd(C.int(arrayHandle), C.int(objectHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to add object to array: %s", lastErrorString())
	}

	return nil
}

// FreeObject frees an Indigo object (array, buffer, etc.)
func (in *Indigo) FreeObject(handle int) error {
	if handle < 0 {
		return nil // Already invalid
	}

	ret := int(C.indigoFree(C.int(handle)))
	if ret < 0 {
		return fmt.Errorf("failed to free object: %s", lastErrorString())
	}

	return nil
}

// CreateWriteBuffer creates an output buffer for rendering
func (in *Indigo) CreateWriteBuffer() (int, error) {
	handle := int(C.indigoWriteBuffer())
	if handle < 0 {
		return 0, fmt.Errorf("failed to create write buffer: %s", lastErrorString())
	}

	runtime.SetFinalizer(&handle, func(h *int) {
		if *h >= 0 {
			C.indigoFree(C.int(*h))
		}
	})

	return handle, nil
}

// GetBufferData retrieves data from a write buffer
func (in *Indigo) GetBufferData(bufferHandle int) ([]byte, error) {
	if bufferHandle < 0 {
		return nil, fmt.Errorf("invalid buffer handle")
	}

	var size C.int
	var dataPtr *C.char
	ret := C.indigoToBuffer(C.int(bufferHandle), &dataPtr, &size)
	if ret < 0 || dataPtr == nil {
		return nil, fmt.Errorf("failed to get buffer data: %s", lastErrorString())
	}

	// Copy C data to Go slice
	// Note: dataPtr is managed by Indigo internally, don't free it
	data := C.GoBytes(unsafe.Pointer(dataPtr), size)
	return data, nil
}

// GetOption returns option value as string.
func (in *Indigo) GetOption(option string) (string, error) {
	in.setSession()
	copt := C.CString(option)
	defer C.free(unsafe.Pointer(copt))
	ptr := C.indigoGetOption(copt)
	if ptr == nil {
		return "", errors.New(lastErrorString())
	}
	return C.GoString(ptr), nil
}

// GetOptionInt returns integer option value.
func (in *Indigo) GetOptionInt(option string) (int, error) {
	in.setSession()
	copt := C.CString(option)
	defer C.free(unsafe.Pointer(copt))
	var out C.int
	if C.indigoGetOptionInt(copt, &out) == 0 {
		return 0, errors.New(lastErrorString())
	}
	return int(out), nil
}

// Next obtains the next element from an iterator, returns 0 if there is no next element
func (in *Indigo) Next(iterHandle int) (int, error) {
	in.setSession()
	result := int(C.indigoNext(C.int(iterHandle)))
	if result < 0 {
		return 0, fmt.Errorf("failed to get next element: %s", lastErrorString())
	}
	return result, nil
}

// HasNext checks if there is a next element without obtaining it
func (in *Indigo) HasNext(iterHandle int) (bool, error) {
	in.setSession()
	result := int(C.indigoHasNext(C.int(iterHandle)))
	if result < 0 {
		return false, fmt.Errorf("failed to check next element: %s", lastErrorString())
	}
	return result != 0, nil
}

// Index returns the index of an element
func (in *Indigo) Index(itemHandle int) (int, error) {
	in.setSession()
	result := int(C.indigoIndex(C.int(itemHandle)))
	if result < 0 {
		return 0, fmt.Errorf("failed to get index: %s", lastErrorString())
	}
	return result, nil
}

// Remove removes an item from its container (usually a molecule)
func (in *Indigo) Remove(itemHandle int) error {
	in.setSession()
	result := int(C.indigoRemove(C.int(itemHandle)))
	if result < 0 {
		return fmt.Errorf("failed to remove item: %s", lastErrorString())
	}
	return nil
}

// GetOriginalFormat returns the original format of an item
func (in *Indigo) GetOriginalFormat(itemHandle int) (string, error) {
	in.setSession()
	cStr := C.indigoGetOriginalFormat(C.int(itemHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get original format: %s", lastErrorString())
	}
	return C.GoString(cStr), nil
}
