// Package molecule provides molecular structure manipulation using Indigo library via CGO
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule.go
// @Software: GoLand
package molecule

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
	Handle int
	Closed bool
}

// Close frees the Indigo molecule object
func (m *Molecule) Close() error {
	if m.Closed || m.Handle < 0 {
		return nil
	}

	ret := int(C.indigoFree(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to free molecule: %s", getLastError())
	}

	m.Closed = true
	m.Handle = -1
	return nil
}

// Clone creates a deep copy of the molecule
func (m *Molecule) Clone() (*Molecule, error) {
	if m.Closed {
		return nil, fmt.Errorf("molecule is closed")
	}

	newHandle := int(C.indigoClone(C.int(m.Handle)))
	if newHandle < 0 {
		return nil, fmt.Errorf("failed to clone molecule: %s", getLastError())
	}

	newMol := &Molecule{
		Handle: newHandle,
		Closed: false,
	}

	runtime.SetFinalizer(newMol, (*Molecule).Close)
	return newMol, nil
}

// CountAtoms returns the number of atoms in the molecule
func (m *Molecule) CountAtoms() (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountAtoms(C.int(m.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count atoms: %s", getLastError())
	}

	return count, nil
}

// CountBonds returns the number of bonds in the molecule
func (m *Molecule) CountBonds() (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountBonds(C.int(m.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count bonds: %s", getLastError())
	}

	return count, nil
}

// CountHeavyAtoms returns the number of heavy (non-hydrogen) atoms
func (m *Molecule) CountHeavyAtoms() (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountHeavyAtoms(C.int(m.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count heavy atoms: %s", getLastError())
	}

	return count, nil
}

// Aromatize performs aromatization of the molecule
func (m *Molecule) Aromatize() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoAromatize(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to aromatize: %s", getLastError())
	}

	return nil
}

// Dearomatize removes aromaticity from the molecule
func (m *Molecule) Dearomatize() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoDearomatize(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to dearomatize: %s", getLastError())
	}

	return nil
}

// FoldHydrogens folds hydrogens in the molecule
func (m *Molecule) FoldHydrogens() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoFoldHydrogens(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to fold hydrogens: %s", getLastError())
	}

	return nil
}

// UnfoldHydrogens unfolds hydrogens in the molecule
func (m *Molecule) UnfoldHydrogens() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoUnfoldHydrogens(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to unfold hydrogens: %s", getLastError())
	}

	return nil
}

// Layout performs 2D layout of the molecule
func (m *Molecule) Layout() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoLayout(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to layout: %s", getLastError())
	}

	return nil
}

// Clean2D performs 2D cleaning of the molecule
func (m *Molecule) Clean2D() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoClean2d(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to clean2d: %s", getLastError())
	}

	return nil
}

// Normalize normalizes the molecule structure
// It neutralizes charges, resolves 5-valence Nitrogen, removes hydrogens, etc.
func (m *Molecule) Normalize(options string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cOptions := C.CString(options)
	defer C.free(unsafe.Pointer(cOptions))

	ret := int(C.indigoNormalize(C.int(m.Handle), cOptions))
	if ret < 0 {
		return fmt.Errorf("failed to normalize: %s", getLastError())
	}

	return nil
}

// Standardize standardizes the molecule structure
func (m *Molecule) Standardize() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoStandardize(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to standardize: %s", getLastError())
	}

	return nil
}

// Ionize ionizes the molecule at specified pH and pH tolerance
func (m *Molecule) Ionize(pH float32, pHTolerance float32) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoIonize(C.int(m.Handle), C.float(pH), C.float(pHTolerance)))
	if ret < 0 {
		return fmt.Errorf("failed to ionize: %s", getLastError())
	}

	return nil
}

// CountComponents returns the number of connected components
func (m *Molecule) CountComponents() (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountComponents(C.int(m.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count components: %s", getLastError())
	}

	return count, nil
}

// CountSSSR returns the number of smallest set of smallest rings
func (m *Molecule) CountSSSR() (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoCountSSSR(C.int(m.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count SSSR: %s", getLastError())
	}

	return count, nil
}

// newMolecule is a helper function to create a Molecule object from a handle
// It sets up the finalizer to ensure proper cleanup
func newMolecule(handle int) *Molecule {
	m := &Molecule{
		Handle: handle,
		Closed: false,
	}
	runtime.SetFinalizer(m, (*Molecule).Close)
	return m
}

// getLastError retrieves the last error message from Indigo
func getLastError() string {
	errMsg := C.indigoGetLastError()
	if errMsg == nil {
		return "unknown error"
	}
	return C.GoString(errMsg)
}
