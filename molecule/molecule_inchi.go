// Package molecule provides InChI functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_inchi.go
// @Software: GoLand
package molecule

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd
#cgo windows LDFLAGS: -L${SRCDIR}/../3rd/win -lindigo -lindigo-inchi
#cgo linux LDFLAGS: -L${SRCDIR}/../3rd/linux -lindigo -lindigo-inchi -Wl,-rpath,${SRCDIR}/../3rd/linux

#include <stdlib.h>
#include "indigo.h"
#include "indigo-inchi.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ToInChI converts the molecule to InChI format
// This uses Indigo's InChI plugin
func (m *Molecule) ToInChI() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoInchiGetInchi(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to InChI: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToInChIKey converts the molecule to InChI Key format
func (m *Molecule) ToInChIKey() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// First get the InChI
	inchi, err := m.ToInChI()
	if err != nil {
		return "", err
	}

	// Generate InChI Key from InChI
	cInchi := C.CString(inchi)
	defer C.free(unsafe.Pointer(cInchi))

	cKey := C.indigoInchiGetInchiKey(cInchi)
	if cKey == nil {
		return "", fmt.Errorf("failed to generate InChI Key: %s", getLastError())
	}

	return C.GoString(cKey), nil
}

// GetInChIWarning returns any warnings from the last InChI generation
func GetInChIWarning() string {
	cStr := C.indigoInchiGetWarning()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// GetInChILog returns log messages from the last InChI generation
func GetInChILog() string {
	cStr := C.indigoInchiGetLog()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// GetInChIAuxInfo returns auxiliary information from the last InChI generation
func GetInChIAuxInfo() string {
	cStr := C.indigoInchiGetAuxInfo()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// LoadInChI loads a molecule from InChI string
func LoadInChI(inchi string) (*Molecule, error) {
	cInchi := C.CString(inchi)
	defer C.free(unsafe.Pointer(cInchi))

	handle := int(C.indigoInchiLoadMolecule(cInchi))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from InChI: %s", getLastError())
	}

	m := &Molecule{
		handle: handle,
		closed: false,
	}

	return m, nil
}
