// Package molecule provides molecule building functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_builder.go
// @Software: GoLand
package molecule

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

// Bond order constants
const (
	BOND_SINGLE   = 1
	BOND_DOUBLE   = 2
	BOND_TRIPLE   = 3
	BOND_AROMATIC = 4
)

// AddAtom adds an atom to the molecule
// symbol: element symbol (e.g., "C", "N", "O") or pseudoatom label
// Returns the handle of the added atom
func (m *Molecule) AddAtom(symbol string) (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	cSymbol := C.CString(symbol)
	defer C.free(unsafe.Pointer(cSymbol))

	handle := int(C.indigoAddAtom(C.int(m.handle), cSymbol))
	if handle < 0 {
		return 0, fmt.Errorf("failed to add atom %s: %s", symbol, getLastError())
	}

	return handle, nil
}

// ResetAtom resets an atom to a new element
func ResetAtom(atomHandle int, symbol string) error {
	cSymbol := C.CString(symbol)
	defer C.free(unsafe.Pointer(cSymbol))

	ret := int(C.indigoResetAtom(C.int(atomHandle), cSymbol))
	if ret < 0 {
		return fmt.Errorf("failed to reset atom: %s", getLastError())
	}

	return nil
}

// AddRSite adds an R-site to the molecule
// name: R-site name (e.g., "R", "R1", "R2", or list "R1 R3")
// Returns the handle of the added R-site
func (m *Molecule) AddRSite(name string) (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	handle := int(C.indigoAddRSite(C.int(m.handle), cName))
	if handle < 0 {
		return 0, fmt.Errorf("failed to add R-site %s: %s", name, getLastError())
	}

	return handle, nil
}

// SetRSite sets an atom as an R-site
func SetRSite(atomHandle int, name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := int(C.indigoSetRSite(C.int(atomHandle), cName))
	if ret < 0 {
		return fmt.Errorf("failed to set R-site: %s", getLastError())
	}

	return nil
}

// AddBond adds a bond between two atoms
// source: handle of the source atom
// destination: handle of the destination atom
// order: bond order (BOND_SINGLE, BOND_DOUBLE, BOND_TRIPLE, BOND_AROMATIC)
// Returns the handle of the added bond
func (m *Molecule) AddBond(source int, destination int, order int) (int, error) {
	if m.closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	handle := int(C.indigoAddBond(C.int(source), C.int(destination), C.int(order)))
	if handle < 0 {
		return 0, fmt.Errorf("failed to add bond: %s", getLastError())
	}

	return handle, nil
}

// SetBondOrder sets the order of a bond
func SetBondOrder(bondHandle int, order int) error {
	ret := int(C.indigoSetBondOrder(C.int(bondHandle), C.int(order)))
	if ret < 0 {
		return fmt.Errorf("failed to set bond order: %s", getLastError())
	}

	return nil
}

// SetCharge sets the charge of an atom
func SetCharge(atomHandle int, charge int) error {
	ret := int(C.indigoSetCharge(C.int(atomHandle), C.int(charge)))
	if ret < 0 {
		return fmt.Errorf("failed to set charge: %s", getLastError())
	}

	return nil
}

// SetIsotope sets the isotope of an atom
func SetIsotope(atomHandle int, isotope int) error {
	ret := int(C.indigoSetIsotope(C.int(atomHandle), C.int(isotope)))
	if ret < 0 {
		return fmt.Errorf("failed to set isotope: %s", getLastError())
	}

	return nil
}

// SetRadical sets the radical of an atom
func SetRadical(atomHandle int, radical int) error {
	ret := int(C.indigoSetRadical(C.int(atomHandle), C.int(radical)))
	if ret < 0 {
		return fmt.Errorf("failed to set radical: %s", getLastError())
	}

	return nil
}

// ResetRadical resets the radical of an atom
func ResetRadical(atomHandle int) error {
	ret := int(C.indigoResetRadical(C.int(atomHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to reset radical: %s", getLastError())
	}

	return nil
}

// SetImplicitHCount sets the implicit hydrogen count of an atom
func SetImplicitHCount(atomHandle int, count int) error {
	ret := int(C.indigoSetImplicitHCount(C.int(atomHandle), C.int(count)))
	if ret < 0 {
		return fmt.Errorf("failed to set implicit H count: %s", getLastError())
	}

	return nil
}

// Merge merges another molecule into this molecule
func (m *Molecule) Merge(other *Molecule) error {
	if m.closed {
		return fmt.Errorf("molecule is closed")
	}
	if other.closed {
		return fmt.Errorf("other molecule is closed")
	}

	ret := int(C.indigoMerge(C.int(m.handle), C.int(other.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to merge molecules: %s", getLastError())
	}

	return nil
}

// GetRadicalElectrons gets the number of radical electrons of an atom
func GetRadicalElectrons(atomHandle int) (int, bool, error) {
	var electrons C.int
	ret := int(C.indigoGetRadicalElectrons(C.int(atomHandle), &electrons))
	if ret < 0 {
		return 0, false, fmt.Errorf("failed to get radical electrons: %s", getLastError())
	}

	return int(electrons), ret == 1, nil
}

// GetRadical gets the radical of an atom
func GetRadical(atomHandle int) (int, bool, error) {
	var radical C.int
	ret := int(C.indigoGetRadical(C.int(atomHandle), &radical))
	if ret < 0 {
		return 0, false, fmt.Errorf("failed to get radical: %s", getLastError())
	}

	return int(radical), ret == 1, nil
}
