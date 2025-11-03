// Package reaction provides reaction saving functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_saver.go
// @Software: GoLand
package reaction

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd
#cgo windows LDFLAGS: -L${SRCDIR}/../3rd/win -lindigo
#cgo linux LDFLAGS: -L${SRCDIR}/../3rd/linux -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux

#include <stdlib.h>
#include "indigo.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ToRxnfile returns the reaction as an RXN file string
func (r *Reaction) ToRxnfile() (string, error) {
	if r.closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoRxnfile(C.int(r.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to RXN: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToFile saves the reaction to a file in RXN format
func (r *Reaction) SaveToFile(filename string) error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveRxnfileToFile(C.int(r.handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to file %s: %s", filename, getLastError())
	}

	return nil
}

// SaveRxnfile saves the reaction to an output object
// outputHandle is the Indigo handle of an output object
func (r *Reaction) SaveRxnfile(outputHandle int) error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoSaveRxnfile(C.int(r.handle), C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to output: %s", getLastError())
	}

	return nil
}

// ToSmiles converts the reaction to SMILES format
func (r *Reaction) ToSmiles() (string, error) {
	if r.closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoSmiles(C.int(r.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalSmiles converts the reaction to canonical SMILES format
func (r *Reaction) ToCanonicalSmiles() (string, error) {
	if r.closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoCanonicalSmiles(C.int(r.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to canonical SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// CreateStringOutput creates a string output buffer
// Returns the Indigo handle for the output buffer
func CreateStringOutput() (int, error) {
	handle := int(C.indigoWriteBuffer())
	if handle < 0 {
		return 0, fmt.Errorf("failed to create string output: %s", getLastError())
	}
	return handle, nil
}

// GetStringOutput retrieves the string from an output buffer
func GetStringOutput(outputHandle int) (string, error) {
	cStr := C.indigoToString(C.int(outputHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get string from output: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// FreeOutput frees an output buffer
func FreeOutput(outputHandle int) error {
	ret := int(C.indigoFree(C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to free output: %s", getLastError())
	}
	return nil
}
