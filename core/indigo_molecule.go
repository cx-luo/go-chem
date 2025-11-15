// Package core provides core functions for Indigo C API library via CGO
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/12 13:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : indigo_molecule.go
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
	"github.com/cx-luo/go-indigo/molecule"
	"runtime"
	"unsafe"
)

// CreateMolecule creates a new empty molecule
func (in *Indigo) CreateMolecule() (*molecule.Molecule, error) {
	in.setSession()
	handle := int(C.indigoCreateMolecule())
	if handle < 0 {
		return nil, fmt.Errorf("failed to create molecule: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// CreateQueryMolecule creates a new empty query molecule
func (in *Indigo) CreateQueryMolecule() (*molecule.Molecule, error) {
	in.setSession()
	handle := int(C.indigoCreateQueryMolecule())
	if handle < 0 {
		return nil, fmt.Errorf("failed to create query molecule: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadMoleculeFromString loads a molecule from a string and returns IndigoObject.
func (in *Indigo) LoadMoleculeFromString(s string) (*molecule.Molecule, error) {
	in.setSession()
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	handle := int(C.indigoLoadMoleculeFromString(cs))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from string: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadMoleculeFromFile loads a molecule from a file
func (in *Indigo) LoadMoleculeFromFile(filename string) (*molecule.Molecule, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadMoleculeFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from file %s: %s", filename, lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadMoleculeFromBuffer loads a molecule from a byte buffer
func (in *Indigo) LoadMoleculeFromBuffer(buffer []byte) (*molecule.Molecule, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadMoleculeFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from buffer: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadQueryMoleculeFromString loads a query molecule from a string
func (in *Indigo) LoadQueryMoleculeFromString(data string) (*molecule.Molecule, error) {
	in.setSession()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	handle := int(C.indigoLoadQueryMoleculeFromString(cData))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query molecule from string: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadQueryMoleculeFromFile loads a query molecule from a file
func (in *Indigo) LoadQueryMoleculeFromFile(filename string) (*molecule.Molecule, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadQueryMoleculeFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query molecule from file %s: %s", filename, lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadQueryMoleculeFromBuffer loads a query molecule from a byte buffer
func (in *Indigo) LoadQueryMoleculeFromBuffer(buffer []byte) (*molecule.Molecule, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadQueryMoleculeFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query molecule from buffer: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadSmartsFromString loads a SMARTS pattern from a string
func (in *Indigo) LoadSmartsFromString(smarts string) (*molecule.Molecule, error) {
	in.setSession()
	cSmarts := C.CString(smarts)
	defer C.free(unsafe.Pointer(cSmarts))

	handle := int(C.indigoLoadSmartsFromString(cSmarts))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load SMARTS from string: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadSmartsFromFile loads a SMARTS pattern from a file
func (in *Indigo) LoadSmartsFromFile(filename string) (*molecule.Molecule, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadSmartsFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load SMARTS from file %s: %s", filename, lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadSmartsFromBuffer loads a SMARTS pattern from a byte buffer
func (in *Indigo) LoadSmartsFromBuffer(buffer []byte) (*molecule.Molecule, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadSmartsFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load SMARTS from buffer: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadStructureFromString loads a structure from a string with parameters
func (in *Indigo) LoadStructureFromString(data string, params string) (*molecule.Molecule, error) {
	in.setSession()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	cParams := C.CString(params)
	defer C.free(unsafe.Pointer(cParams))

	handle := int(C.indigoLoadStructureFromString(cData, cParams))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load structure from string: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadStructureFromFile loads a structure from a file with parameters
func (in *Indigo) LoadStructureFromFile(filename string, params string) (*molecule.Molecule, error) {
	in.setSession()
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cParams := C.CString(params)
	defer C.free(unsafe.Pointer(cParams))

	handle := int(C.indigoLoadStructureFromFile(cFilename, cParams))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load structure from file %s: %s", filename, lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadStructureFromBuffer loads a structure from a byte buffer with parameters
func (in *Indigo) LoadStructureFromBuffer(buffer []byte, params string) (*molecule.Molecule, error) {
	in.setSession()
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.byte)(unsafe.Pointer(&buffer[0]))
	cParams := C.CString(params)
	defer C.free(unsafe.Pointer(cParams))

	handle := int(C.indigoLoadStructureFromBuffer(cBuffer, C.int(len(buffer)), cParams))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load structure from buffer: %s", lastErrorString())
	}

	return newMolecule(handle), nil
}

// LoadMoleculeFromHandle creates a Molecule object from an existing Indigo handle
// This is useful when getting molecule handles from reactions or other sources
// Note: The molecule will take ownership of the handle and will free it on Close()
func (in *Indigo) LoadMoleculeFromHandle(handle int) (*molecule.Molecule, error) {
	in.setSession()
	if handle < 0 {
		return nil, fmt.Errorf("invalid handle: %d", handle)
	}
	return newMolecule(handle), nil
}

// Similarity Example: similarity between two objects (returns float)
func (in *Indigo) Similarity(item1, item2 *IndigoObject, metrics string) (float64, error) {
	in.setSession()
	cmetrics := C.CString(metrics)
	defer C.free(unsafe.Pointer(cmetrics))
	res := C.indigoSimilarity(C.int(item1.id), C.int(item2.id), cmetrics)
	// Assuming negative or specific value indicates error; adjust as per API
	if res < 0 {
		return 0.0, fmt.Errorf(lastErrorString())
	}
	return float64(res), nil
}

// newMolecule is a helper function to create a Molecule object from a handle
// It sets up the finalizer to ensure proper cleanup
func newMolecule(handle int) *molecule.Molecule {
	m := &molecule.Molecule{
		Handle: handle,
		Closed: false,
	}
	runtime.SetFinalizer(m, (*molecule.Molecule).Close)
	return m
}
