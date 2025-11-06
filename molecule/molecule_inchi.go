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

// Windows: link against import libraries (.lib)
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo -lindigo-inchi
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo -lindigo-inchi

// Linux: use $ORIGIN for runtime library search
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -lindigo-inchi -Wl,-rpath='$ORIGIN/../3rd/linux-x86_64'
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -lindigo-inchi -Wl,-rpath='$ORIGIN/../3rd/linux-aarch64'

// macOS: use @loader_path (not @executable_path) for shared libraries
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -lindigo-inchi -Wl,-rpath,'@loader_path/../3rd/darwin-x86_64'
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -lindigo-inchi -Wl,-rpath,'@loader_path/../3rd/darwin-aarch64'

#include <stdlib.h>
#include "indigo.h"
#include "indigo-inchi.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// inchiInitialized tracks whether InChI module has been initialized
var inchiInitialized = false

// InChIVersion returns the version of the InChI library
func InChIVersion() string {
	cStr := C.indigoInchiVersion()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// InitInChI initializes the InChI module for the current session
// This should be called before using InChI functions
func InitInChI() error {
	if inchiInitialized {
		return nil // Already initialized
	}

	ret := int(C.indigoInchiInit(indigoSessionID))
	if ret < 0 {
		return fmt.Errorf("failed to initialize InChI: %s", getLastError())
	}

	inchiInitialized = true
	return nil
}

// DisposeInChI disposes the InChI module
// This should be called when done using InChI functions
func DisposeInChI() error {
	if !inchiInitialized {
		return nil // Not initialized
	}

	ret := int(C.indigoInchiDispose(indigoSessionID))
	if ret < 0 {
		return fmt.Errorf("failed to dispose InChI: %s", getLastError())
	}

	inchiInitialized = false
	return nil
}

// ResetInChIOptions resets InChI options to default
func ResetInChIOptions() error {
	ret := int(C.indigoInchiResetOptions())
	if ret < 0 {
		return fmt.Errorf("failed to reset InChI options: %s", getLastError())
	}
	return nil
}

// ToInChI converts the molecule to InChI format
// This uses Indigo's InChI plugin
func (m *Molecule) ToInChI() (string, error) {
	if m == nil {
		return "", fmt.Errorf("molecule is nil")
	}
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}
	if m.handle < 0 {
		return "", fmt.Errorf("invalid molecule handle")
	}

	// Ensure InChI module is initialized
	if !inchiInitialized {
		if err := InitInChI(); err != nil {
			return "", fmt.Errorf("failed to initialize InChI: %w", err)
		}
	}

	cStr := C.indigoInchiGetInchi(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to InChI: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToInChIKey converts the molecule to InChI Key format
func (m *Molecule) ToInChIKey() (string, error) {
	if m == nil {
		return "", fmt.Errorf("molecule is nil")
	}
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}
	if m.handle < 0 {
		return "", fmt.Errorf("invalid molecule handle")
	}

	// First get the InChI
	inchi, err := m.ToInChI()
	if err != nil {
		return "", fmt.Errorf("failed to get InChI: %w", err)
	}

	// Generate InChI Key from InChI
	return InChIToKey(inchi)
}

// InChIToKey converts an InChI string to InChIKey
func InChIToKey(inchi string) (string, error) {
	if inchi == "" {
		return "", fmt.Errorf("empty InChI string")
	}

	// Ensure InChI module is initialized
	if !inchiInitialized {
		if err := InitInChI(); err != nil {
			return "", fmt.Errorf("failed to initialize InChI: %w", err)
		}
	}

	cInchi := C.CString(inchi)
	defer C.free(unsafe.Pointer(cInchi))

	cKey := C.indigoInchiGetInchiKey(cInchi)
	if cKey == nil {
		return "", fmt.Errorf("failed to generate InChI Key: %s", getLastError())
	}

	return C.GoString(cKey), nil
}

// InChIWarning returns any warnings from the last InChI generation
func InChIWarning() string {
	cStr := C.indigoInchiGetWarning()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// InChILog returns log messages from the last InChI generation
func InChILog() string {
	cStr := C.indigoInchiGetLog()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// InChIAuxInfo returns auxiliary information from the last InChI generation
func InChIAuxInfo() string {
	cStr := C.indigoInchiGetAuxInfo()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// LoadFromInChI loads a molecule from InChI string
func LoadFromInChI(inchi string) (*Molecule, error) {
	if inchi == "" {
		return nil, fmt.Errorf("empty InChI string")
	}

	// Ensure InChI module is initialized
	if !inchiInitialized {
		if err := InitInChI(); err != nil {
			return nil, fmt.Errorf("failed to initialize InChI: %w", err)
		}
	}

	cInchi := C.CString(inchi)
	defer C.free(unsafe.Pointer(cInchi))

	handle := int(C.indigoInchiLoadMolecule(cInchi))
	if handle < 0 {
		return nil, fmt.Errorf("failed to load molecule from InChI: %s", getLastError())
	}

	return newMolecule(handle), nil
}

// InChIResult contains the result of InChI generation
type InChIResult struct {
	InChI   string // The InChI string
	Key     string // The InChIKey
	Warning string // Warning messages
	Log     string // Log messages
	AuxInfo string // Auxiliary information
}

// ToInChIWithInfo converts the molecule to InChI format and returns detailed information
func (m *Molecule) ToInChIWithInfo() (*InChIResult, error) {
	inchi, err := m.ToInChI()
	if err != nil {
		return nil, err
	}

	key, err := InChIToKey(inchi)
	if err != nil {
		return nil, fmt.Errorf("failed to generate InChIKey: %w", err)
	}

	result := &InChIResult{
		InChI:   inchi,
		Key:     key,
		Warning: InChIWarning(),
		Log:     InChILog(),
		AuxInfo: InChIAuxInfo(),
	}

	return result, nil
}
