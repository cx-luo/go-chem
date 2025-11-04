// Package molecule provides molecule saving functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_saver.go
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
	"unsafe"
)

// ToSmiles converts the molecule to SMILES format
func (m *Molecule) ToSmiles() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoSmiles(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalSmiles converts the molecule to canonical SMILES format
func (m *Molecule) ToCanonicalSmiles() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoCanonicalSmiles(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to canonical SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToSmarts converts the molecule to SMARTS format
func (m *Molecule) ToSmarts() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoSmarts(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to SMARTS: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalSmarts converts the molecule to canonical SMARTS format
func (m *Molecule) ToCanonicalSmarts() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoCanonicalSmarts(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to canonical SMARTS: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToMolfile returns the molecule as a MOL file string
func (m *Molecule) ToMolfile() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoMolfile(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to MOL: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToFile saves the molecule to a file in MOL format
func (m *Molecule) SaveToFile(filename string) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveMolfileToFile(C.int(m.handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save to file %s: %s", filename, getLastError())
	}

	return nil
}

// SaveMolfile saves the molecule to an output object
// outputHandle is the Indigo handle of an output object
func (m *Molecule) SaveMolfile(outputHandle int) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoSaveMolfile(C.int(m.handle), C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to save to output: %s", getLastError())
	}

	return nil
}

// ToJSON converts the molecule to JSON format
func (m *Molecule) ToJSON() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// Create a string output buffer
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return "", fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Save as JSON
	ret := int(C.indigoSaveJson(C.int(m.handle), C.int(bufferHandle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to save as JSON: %s", getLastError())
	}

	// Get the string
	cStr := C.indigoToString(C.int(bufferHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get JSON string: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToJSONFile saves the molecule to a file in JSON format
func (m *Molecule) SaveToJSONFile(filename string) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveJsonToFile(C.int(m.handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save to JSON file %s: %s", filename, getLastError())
	}

	return nil
}

// ToBase64String converts the molecule to base64 string
func (m *Molecule) ToBase64String() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoToBase64String(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to base64: %s", getLastError())
	}

	return C.GoString(cStr), nil
}
