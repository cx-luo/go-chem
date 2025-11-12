// Package core coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/12 13:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : indigo-reaction.go
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
	"fmt"
	"github.com/cx-luo/go-chem/reaction"
	"runtime"
	"unsafe"
)

// CreateReaction creates a new empty reaction
func (in *Indigo) CreateReaction() (*reaction.Reaction, error) {
	handle := int(C.indigoCreateReaction())
	if handle < 0 {
		return nil, fmt.Errorf("failed to create reaction: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// CreateQueryReaction creates a new empty query reaction
func (in *Indigo) CreateQueryReaction() (*reaction.Reaction, error) {
	handle := int(C.indigoCreateQueryReaction())
	if handle < 0 {
		return nil, fmt.Errorf("failed to create query reaction: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// newReaction is a helper function to create a Reaction object from a handle
// It sets up the finalizer to ensure proper cleanup
func newReaction(handle int) *reaction.Reaction {
	r := &reaction.Reaction{
		Handle: handle,
		Closed: false,
	}
	runtime.SetFinalizer(r, (*reaction.Reaction).Close)
	return r
}

// LoadReactionFromString loads a reaction from a string
// The string should contain reaction data in a supported format (RXN, SMILES, etc.)
func (in *Indigo) LoadReactionFromString(data string) (*reaction.Reaction, error) {
	in.setSession()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	handle := int(C.indigoLoadReactionFromString(cData))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction from string: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionFromFile loads a reaction from a file
func (in *Indigo) LoadReactionFromFile(filename string) (*reaction.Reaction, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadReactionFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction from file %s: %s", filename, lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionFromBuffer loads a reaction from a byte buffer
func (in *Indigo) LoadReactionFromBuffer(buffer []byte) (*reaction.Reaction, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadReactionFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction from buffer: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadQueryReactionFromString loads a query reaction from a string
func (in *Indigo) LoadQueryReactionFromString(data string) (*reaction.Reaction, error) {
	in.setSession()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	handle := int(C.indigoLoadQueryReactionFromString(cData))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query reaction from string: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadQueryReactionFromFile loads a query reaction from a file
func (in *Indigo) LoadQueryReactionFromFile(filename string) (*reaction.Reaction, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadQueryReactionFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query reaction from file %s: %s", filename, lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadQueryReactionFromBuffer loads a query reaction from a byte buffer
func (in *Indigo) LoadQueryReactionFromBuffer(buffer []byte) (*reaction.Reaction, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadQueryReactionFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query reaction from buffer: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionSmartsFromString loads a reaction SMARTS from a string
func (in *Indigo) LoadReactionSmartsFromString(smarts string) (*reaction.Reaction, error) {
	in.setSession()
	cSmarts := C.CString(smarts)
	defer C.free(unsafe.Pointer(cSmarts))

	handle := int(C.indigoLoadReactionSmartsFromString(cSmarts))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction SMARTS from string: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionSmartsFromFile loads a reaction SMARTS from a file
func (in *Indigo) LoadReactionSmartsFromFile(filename string) (*reaction.Reaction, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadReactionSmartsFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction SMARTS from file %s: %s", filename, lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionSmartsFromBuffer loads a reaction SMARTS from a byte buffer
func (in *Indigo) LoadReactionSmartsFromBuffer(buffer []byte) (*reaction.Reaction, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadReactionSmartsFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction SMARTS from buffer: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionWithLibFromString loads a reaction from a string with a monomer library
func (in *Indigo) LoadReactionWithLibFromString(data string, monomerLibraryHandle int) (*reaction.Reaction, error) {
	in.setSession()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	handle := int(C.indigoLoadReactionWithLibFromString(cData, C.int(monomerLibraryHandle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction with library from string: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionWithLibFromFile loads a reaction from a file with a monomer library
func (in *Indigo) LoadReactionWithLibFromFile(filename string, monomerLibraryHandle int) (*reaction.Reaction, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadReactionWithLibFromFile(cFilename, C.int(monomerLibraryHandle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction with library from file %s: %s", filename, lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadReactionWithLibFromBuffer loads a reaction from a byte buffer with a monomer library
func (in *Indigo) LoadReactionWithLibFromBuffer(buffer []byte, monomerLibraryHandle int) (*reaction.Reaction, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadReactionWithLibFromBuffer(cBuffer, C.int(len(buffer)), C.int(monomerLibraryHandle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load reaction with library from buffer: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadQueryReactionWithLibFromString loads a query reaction from a string with a monomer library
func (in *Indigo) LoadQueryReactionWithLibFromString(data string, monomerLibraryHandle int) (*reaction.Reaction, error) {
	in.setSession()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	handle := int(C.indigoLoadQueryReactionWithLibFromString(cData, C.int(monomerLibraryHandle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query reaction with library from string: %s", lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadQueryReactionWithLibFromFile loads a query reaction from a file with a monomer library
func (in *Indigo) LoadQueryReactionWithLibFromFile(filename string, monomerLibraryHandle int) (*reaction.Reaction, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadQueryReactionWithLibFromFile(cFilename, C.int(monomerLibraryHandle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query reaction with library from file %s: %s", filename, lastErrorString())
	}

	return newReaction(handle), nil
}

// LoadQueryReactionWithLibFromBuffer loads a query reaction from a byte buffer with a monomer library
func (in *Indigo) LoadQueryReactionWithLibFromBuffer(buffer []byte, monomerLibraryHandle int) (*reaction.Reaction, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadQueryReactionWithLibFromBuffer(cBuffer, C.int(len(buffer)), C.int(monomerLibraryHandle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query reaction with library from buffer: %s", lastErrorString())
	}

	return newReaction(handle), nil
}
