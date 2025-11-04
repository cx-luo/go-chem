// Package molecule provides molecule loading functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_loader.go
// @Software: GoLand
package molecule

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// Windows platforms
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo

// Linux platforms
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-aarch64

// macOS platforms
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/darwin-x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/darwin-aarch64

#include <stdlib.h>
#include "indigo.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// LoadMoleculeFromString loads a molecule from a string (SMILES, MOL, etc.)
func LoadMoleculeFromString(data string) (*Molecule, error) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	handle := int(C.indigoLoadMoleculeFromString(cData))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from string: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadMoleculeFromFile loads a molecule from a file
func LoadMoleculeFromFile(filename string) (*Molecule, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadMoleculeFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from file %s: %s", filename, getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadMoleculeFromBuffer loads a molecule from a byte buffer
func LoadMoleculeFromBuffer(buffer []byte) (*Molecule, error) {
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadMoleculeFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from buffer: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadQueryMoleculeFromString loads a query molecule from a string
func LoadQueryMoleculeFromString(data string) (*Molecule, error) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	handle := int(C.indigoLoadQueryMoleculeFromString(cData))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query molecule from string: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadQueryMoleculeFromFile loads a query molecule from a file
func LoadQueryMoleculeFromFile(filename string) (*Molecule, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadQueryMoleculeFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query molecule from file %s: %s", filename, getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadQueryMoleculeFromBuffer loads a query molecule from a byte buffer
func LoadQueryMoleculeFromBuffer(buffer []byte) (*Molecule, error) {
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadQueryMoleculeFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load query molecule from buffer: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadSmartsFromString loads a SMARTS pattern from a string
func LoadSmartsFromString(smarts string) (*Molecule, error) {
	cSmarts := C.CString(smarts)
	defer C.free(unsafe.Pointer(cSmarts))

	handle := int(C.indigoLoadSmartsFromString(cSmarts))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load SMARTS from string: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadSmartsFromFile loads a SMARTS pattern from a file
func LoadSmartsFromFile(filename string) (*Molecule, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	handle := int(C.indigoLoadSmartsFromFile(cFilename))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load SMARTS from file %s: %s", filename, getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadSmartsFromBuffer loads a SMARTS pattern from a byte buffer
func LoadSmartsFromBuffer(buffer []byte) (*Molecule, error) {
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	handle := int(C.indigoLoadSmartsFromBuffer(cBuffer, C.int(len(buffer))))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load SMARTS from buffer: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadStructureFromString loads a structure from a string with parameters
func LoadStructureFromString(data string, params string) (*Molecule, error) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	cParams := C.CString(params)
	defer C.free(unsafe.Pointer(cParams))

	handle := int(C.indigoLoadStructureFromString(cData, cParams))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load structure from string: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadStructureFromFile loads a structure from a file with parameters
func LoadStructureFromFile(filename string, params string) (*Molecule, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cParams := C.CString(params)
	defer C.free(unsafe.Pointer(cParams))

	handle := int(C.indigoLoadStructureFromFile(cFilename, cParams))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load structure from file %s: %s", filename, getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}

// LoadStructureFromBuffer loads a structure from a byte buffer with parameters
func LoadStructureFromBuffer(buffer []byte, params string) (*Molecule, error) {
	if len(buffer) == 0 {
		return nil, fmt.Errorf("empty buffer")
	}

	cBuffer := (*C.byte)(unsafe.Pointer(&buffer[0]))
	cParams := C.CString(params)
	defer C.free(unsafe.Pointer(cParams))

	handle := int(C.indigoLoadStructureFromBuffer(cBuffer, C.int(len(buffer)), cParams))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load structure from buffer: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(m, (*Molecule).Close)
	return m, nil
}
