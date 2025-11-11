package core

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// Windows: link against import libraries (.lib)
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo -lindigo-inchi
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo -lindigo-inchi

// Linux: use $ORIGIN for runtime library search
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -lindigo-inchi -Wl,-rpath=${SRCDIR}/../3rd/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -lindigo-inchi -Wl,-rpath=${SRCDIR}/../3rd/linux-aarch64

// macOS: use @loader_path (not @executable_path) for shared libraries
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -lindigo-inchi -Wl,-rpath,${SRCDIR}/../3rd/darwin-x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -lindigo-inchi -Wl,-rpath,${SRCDIR}/../3rd/darwin-aarch64

#include <stdlib.h>
#include "indigo.h"
#include "indigo-inchi.h"
*/
import "C"
import (
	"fmt"
	"github.com/cx-luo/go-chem/molecule"
	"unsafe"
)

// inchiInitialized tracks whether InChI module has been initialized
var inchiInitialized = false

type IndigoInchi struct {
	sid uint64
}

// Molecule represents a chemical molecule in Indigo
type Molecule molecule.Molecule

// getLastError retrieves the last error message from Indigo
func getLastError() string {
	errMsg := C.indigoGetLastError()
	if errMsg == nil {
		return "unknown error"
	}
	return C.GoString(errMsg)
}

// InchiInit initializes the InChI module for the current session
// This should be called before using InChI functions
func InchiInit(sessionID uint64) (*IndigoInchi, error) {
	if inchiInitialized {
		return &IndigoInchi{sid: sessionID}, nil // Already initialized
	}

	ret := int(C.indigoInchiInit(C.ulonglong(sessionID)))
	if ret < 0 {
		return nil, fmt.Errorf("failed to initialize InChI: %s", getLastError())
	}

	inchiInitialized = true
	return &IndigoInchi{sid: sessionID}, nil
}

// DisposeInChI disposes the InChI module
// This should be called when done using InChI functions
func DisposeInChI(sessionID uint64) error {
	if !inchiInitialized {
		return nil // Not initialized
	}

	ret := int(C.indigoInchiDispose(C.qword(sessionID)))
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
func (ii *IndigoInchi) ToInChI(m *Molecule) (string, error) {
	if m == nil {
		return "", fmt.Errorf("molecule is nil")
	}
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}
	if m.Handle < 0 {
		return "", fmt.Errorf("invalid molecule handle")
	}

	cStr := C.indigoInchiGetInchi(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to InChI: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// InChIToKey converts an InChI string to InChIKey
func (ii *IndigoInchi) InChIToKey(inchi string) (string, error) {
	if inchi == "" {
		return "", fmt.Errorf("empty InChI string")
	}

	cInchi := C.CString(inchi)
	defer C.free(unsafe.Pointer(cInchi))

	cKey := C.indigoInchiGetInchiKey(cInchi)
	if cKey == nil {
		return "", fmt.Errorf("failed to generate InChI Key: %s", getLastError())
	}

	return C.GoString(cKey), nil
}

func (ii *IndigoInchi) GetInchiSessionID() uint64 {
	return ii.sid
}

// InChIVersion returns the version of the InChI library
func (ii *IndigoInchi) InChIVersion() string {
	cStr := C.indigoInchiVersion()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// InChIWarning returns any warnings from the last InChI generation
func (ii *IndigoInchi) InChIWarning() string {
	cStr := C.indigoInchiGetWarning()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// InChILog returns log messages from the last InChI generation
func (ii *IndigoInchi) InChILog() string {
	cStr := C.indigoInchiGetLog()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// InChIAuxInfo returns auxiliary information from the last InChI generation
func (ii *IndigoInchi) InChIAuxInfo() string {
	cStr := C.indigoInchiGetAuxInfo()
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// LoadFromInChI loads a molecule from InChI string
func (ii *IndigoInchi) LoadFromInChI(inchi string) (int, error) {
	if inchi == "" {
		return 0, fmt.Errorf("empty InChI string")
	}

	cInchi := C.CString(inchi)
	defer C.free(unsafe.Pointer(cInchi))

	handle := int(C.indigoInchiLoadMolecule(cInchi))
	if handle < 0 {
		return 0, fmt.Errorf("failed to load molecule from InChI: %s", getLastError())
	}

	return handle, nil
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
func (ii *IndigoInchi) ToInChIWithInfo(m *Molecule) (*InChIResult, error) {
	inchi, err := ii.ToInChI(m)
	if err != nil {
		return nil, err
	}

	key, err := ii.InChIToKey(inchi)
	if err != nil {
		return nil, fmt.Errorf("failed to generate InChIKey: %w", err)
	}

	result := &InChIResult{
		InChI:   inchi,
		Key:     key,
		Warning: ii.InChIWarning(),
		Log:     ii.InChILog(),
		AuxInfo: ii.InChIAuxInfo(),
	}

	return result, nil
}
