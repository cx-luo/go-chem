// Package molecule provides molecular structure manipulation using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule.go
// @Software: GoLand
package molecule

/*
#cgo CFLAGS: -I${SRCDIR}/../include

// Windows: link against import libraries (.lib)
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../lib/windows-x86_64 -lindigo -lindigo-inchi
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../lib/windows-i386 -lindigo -lindigo-inchi

// Linux: use $ORIGIN for runtime library search
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../lib/linux-x86_64 -lindigo -lindigo-inchi -Wl,-rpath='$ORIGIN/../lib/linux-x86_64'
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../lib/linux-aarch64 -lindigo -lindigo-inchi -Wl,-rpath='$ORIGIN/../lib/linux-aarch64'

// macOS: use @loader_path (not @executable_path) for shared libraries
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../lib/darwin-x86_64 -lindigo -lindigo-inchi -Wl,-rpath,'@loader_path/../lib/darwin-x86_64'
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../lib/darwin-aarch64 -lindigo -lindigo-inchi -Wl,-rpath,'@loader_path/../lib/darwin-aarch64'
#include <stdlib.h>
#include "indigo.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Radical constants
const (
	RADICAL_SINGLET = 2
	RADICAL_DOUBLET = 3
	RADICAL_TRIPLET = 4
)

// Special element constants
const (
	ELEM_PSEUDO   = -1
	ELEM_RSITE    = -2
	ELEM_TEMPLATE = -3
)

// Molecule represents a chemical molecule in Indigo
type Molecule struct {
	handle int
	closed bool
}

// indigoSessionID holds the session ID for Indigo
var indigoSessionID C.qword

func init() {
	// Initialize Indigo session
	indigoSessionID = C.indigoAllocSessionId()
	C.indigoSetSessionId(indigoSessionID)
}

// CreateMolecule creates a new empty molecule
func CreateMolecule() (*Molecule, error) {
	handle := int(C.indigoCreateMolecule())
	if handle < 0 {
		return nil, fmt.Errorf("failed to create molecule: %s", getLastError())
	}

	return newMolecule(handle), nil
}

// CreateQueryMolecule creates a new empty query molecule
func CreateQueryMolecule() (*Molecule, error) {
	handle := int(C.indigoCreateQueryMolecule())
	if handle < 0 {
		return nil, fmt.Errorf("failed to create query molecule: %s", getLastError())
	}

	return newMolecule(handle), nil
}

// Close frees the Indigo molecule object
func (m *Molecule) Close() error {
	if m.closed || m.handle < 0 {
		return nil
	}

	ret := int(C.indigoFree(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to free molecule: %s", getLastError())
	}

	m.closed = true
	m.handle = -1
	return nil
}

// Handle returns the internal Indigo handle
func (m *Molecule) Handle() int {
	return m.handle
}

// Clone creates a deep copy of the molecule
func (m *Molecule) Clone() (*Molecule, error) {
	if m.closed {
		return nil, fmt.Errorf("molecule is closed")
	}

	newHandle := int(C.indigoClone(C.int(m.handle)))
	if newHandle < 0 {
		return nil, fmt.Errorf("failed to clone molecule: %s", getLastError())
	}

	newMol := &Molecule{
		handle: newHandle,
		closed: false,
	}

	runtime.SetFinalizer(newMol, (*Molecule).Close)
	return newMol, nil
}

// CountAtoms returns the number of atoms in the molecule
func (m *Molecule) CountAtoms() (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountAtoms(C.int(m.handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count atoms: %s", getLastError())
	}

	return count, nil
}

// CountBonds returns the number of bonds in the molecule
func (m *Molecule) CountBonds() (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountBonds(C.int(m.handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count bonds: %s", getLastError())
	}

	return count, nil
}

// CountHeavyAtoms returns the number of heavy (non-hydrogen) atoms
func (m *Molecule) CountHeavyAtoms() (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountHeavyAtoms(C.int(m.handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count heavy atoms: %s", getLastError())
	}

	return count, nil
}

// Aromatize performs aromatization of the molecule
func (m *Molecule) Aromatize() error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoAromatize(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to aromatize: %s", getLastError())
	}

	return nil
}

// Dearomatize removes aromaticity from the molecule
func (m *Molecule) Dearomatize() error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoDearomatize(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to dearomatize: %s", getLastError())
	}

	return nil
}

// FoldHydrogens folds hydrogens in the molecule
func (m *Molecule) FoldHydrogens() error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoFoldHydrogens(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to fold hydrogens: %s", getLastError())
	}

	return nil
}

// UnfoldHydrogens unfolds hydrogens in the molecule
func (m *Molecule) UnfoldHydrogens() error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoUnfoldHydrogens(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to unfold hydrogens: %s", getLastError())
	}

	return nil
}

// Layout performs 2D layout of the molecule
func (m *Molecule) Layout() error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoLayout(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to layout: %s", getLastError())
	}

	return nil
}

// Clean2D performs 2D cleaning of the molecule
func (m *Molecule) Clean2D() error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoClean2d(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to clean2d: %s", getLastError())
	}

	return nil
}

// Normalize normalizes the molecule structure
// It neutralizes charges, resolves 5-valence Nitrogen, removes hydrogens, etc.
func (m *Molecule) Normalize(options string) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	cOptions := C.CString(options)
	defer C.free(unsafe.Pointer(cOptions))

	ret := int(C.indigoNormalize(C.int(m.handle), cOptions))
	if ret < 0 {
		return fmt.Errorf("failed to normalize: %s", getLastError())
	}

	return nil
}

// Standardize standardizes the molecule structure
func (m *Molecule) Standardize() error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoStandardize(C.int(m.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to standardize: %s", getLastError())
	}

	return nil
}

// Ionize ionizes the molecule at specified pH and pH tolerance
func (m *Molecule) Ionize(pH float32, pHTolerance float32) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoIonize(C.int(m.handle), C.float(pH), C.float(pHTolerance)))
	if ret < 0 {
		return fmt.Errorf("failed to ionize: %s", getLastError())
	}

	return nil
}

// CountComponents returns the number of connected components
func (m *Molecule) CountComponents() (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountComponents(C.int(m.handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count components: %s", getLastError())
	}

	return count, nil
}

// CountSSSR returns the number of smallest set of smallest rings
func (m *Molecule) CountSSSR() (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountSSSR(C.int(m.handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count SSSR: %s", getLastError())
	}

	return count, nil
}

// getLastError retrieves the last error message from Indigo
func getLastError() string {
	errMsg := C.indigoGetLastError()
	if errMsg == nil {
		return "unknown error"
	}
	return C.GoString(errMsg)
}

// newMolecule is a helper function to create a Molecule object from a handle
// It sets up the finalizer to ensure proper cleanup
func newMolecule(handle int) *Molecule {
	m := &Molecule{
		handle: handle,
		closed: false,
	}
	runtime.SetFinalizer(m, (*Molecule).Close)
	return m
}
