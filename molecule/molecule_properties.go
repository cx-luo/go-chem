// Package molecule provides molecule properties functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_properties.go
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

// GrossFormula returns the gross formula of the molecule
func (m *Molecule) GrossFormula() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	handle := int(C.indigoGrossFormula(C.int(m.handle)))
	if handle < 0 {
		return "", fmt.Errorf("failed to get gross formula: %s", getLastError())
	}

	cStr := C.indigoToString(C.int(handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert gross formula to string: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// MolecularFormula returns the molecular formula of the molecule
func (m *Molecule) MolecularFormula() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	handle := int(C.indigoMolecularFormula(C.int(m.handle)))
	if handle < 0 {
		return "", fmt.Errorf("failed to get molecular formula: %s", getLastError())
	}

	cStr := C.indigoToString(C.int(handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert molecular formula to string: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// MolecularWeight returns the molecular weight of the molecule
func (m *Molecule) MolecularWeight() (float64, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	weight := float64(C.indigoMolecularWeight(C.int(m.handle)))
	if weight < 0 {
		return 0, fmt.Errorf("failed to get molecular weight: %s", getLastError())
	}

	return weight, nil
}

// MostAbundantMass returns the most abundant mass of the molecule
func (m *Molecule) MostAbundantMass() (float64, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	mass := float64(C.indigoMostAbundantMass(C.int(m.handle)))
	if mass < 0 {
		return 0, fmt.Errorf("failed to get most abundant mass: %s", getLastError())
	}

	return mass, nil
}

// MonoisotopicMass returns the monoisotopic mass of the molecule
func (m *Molecule) MonoisotopicMass() (float64, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	mass := float64(C.indigoMonoisotopicMass(C.int(m.handle)))
	if mass < 0 {
		return 0, fmt.Errorf("failed to get monoisotopic mass: %s", getLastError())
	}

	return mass, nil
}

// MassComposition returns the mass composition of the molecule
func (m *Molecule) MassComposition() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoMassComposition(C.int(m.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get mass composition: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// TPSA returns the topological polar surface area
// includeSP: if true, includes sulfur and phosphorus
func (m *Molecule) TPSA(includeSP bool) (float64, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	var sp C.int
	if includeSP {
		sp = 1
	} else {
		sp = 0
	}

	tpsa := float64(C.indigoTPSA(C.int(m.handle), sp))
	if tpsa < 0 {
		return 0, fmt.Errorf("failed to get TPSA: %s", getLastError())
	}

	return tpsa, nil
}

// NumRotatableBonds returns the number of rotatable bonds
func (m *Molecule) NumRotatableBonds() (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	count := int(C.indigoNumRotatableBonds(C.int(m.handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to get rotatable bonds count: %s", getLastError())
	}

	return count, nil
}

// Name returns the name of the molecule
func (m *Molecule) Name() (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoName(C.int(m.handle))
	if cStr == nil {
		return "", nil // No name set
	}

	return C.GoString(cStr), nil
}

// SetName sets the name of the molecule
func (m *Molecule) SetName(name string) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := int(C.indigoSetName(C.int(m.handle), cName))
	if ret < 0 {
		return fmt.Errorf("failed to set name: %s", getLastError())
	}

	return nil
}

// HasProperty checks if the molecule has a property
func (m *Molecule) HasProperty(prop string) (bool, error) {
	if m.closed {
		return false, fmt.Errorf("molecule is closed")
	}

	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))

	ret := int(C.indigoHasProperty(C.int(m.handle), cProp))
	if ret < 0 {
		return false, fmt.Errorf("failed to check property: %s", getLastError())
	}

	return ret > 0, nil
}

// GetProperty gets a property value from the molecule
func (m *Molecule) GetProperty(prop string) (string, error) {
	if m.closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))

	cStr := C.indigoGetProperty(C.int(m.handle), cProp)
	if cStr == nil {
		return "", fmt.Errorf("failed to get property: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SetProperty sets a property value on the molecule
func (m *Molecule) SetProperty(prop string, value string) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	ret := int(C.indigoSetProperty(C.int(m.handle), cProp, cValue))
	if ret < 0 {
		return fmt.Errorf("failed to set property: %s", getLastError())
	}

	return nil
}

// RemoveProperty removes a property from the molecule
func (m *Molecule) RemoveProperty(prop string) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}

	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))

	ret := int(C.indigoRemoveProperty(C.int(m.handle), cProp))
	if ret < 0 {
		return fmt.Errorf("failed to remove property: %s", getLastError())
	}

	return nil
}
