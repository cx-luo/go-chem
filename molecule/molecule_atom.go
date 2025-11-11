// Package molecule provides atom manipulation functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_atom.go
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
)

// Atom represents an atom in a molecule
type Atom struct {
	handle int
}

// GetAtom returns an atom by its index
func (m *Molecule) GetAtom(index int) (*Atom, error) {
	if m.Closed {
		return nil, fmt.Errorf("molecule is closed")
	}

	handle := int(C.indigoGetAtom(C.int(m.Handle), C.int(index)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to get atom at index %d: %s", index, getLastError())
	}

	return &Atom{handle: handle}, nil
}

// GetBond returns a bond by its index
func (m *Molecule) GetBond(index int) (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}

	handle := int(C.indigoGetBond(C.int(m.Handle), C.int(index)))
	if handle < 0 {
		return 0, fmt.Errorf("failed to get bond at index %d: %s", index, getLastError())
	}

	return handle, nil
}

// Symbol returns the element symbol of an atom
func (a *Atom) Symbol() (string, error) {
	cStr := C.indigoSymbol(C.int(a.handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get atom symbol: %s", getLastError())
	}
	return C.GoString(cStr), nil
}

// Degree returns the number of neighbors of an atom
func (a *Atom) Degree() (int, error) {
	degree := int(C.indigoDegree(C.int(a.handle)))
	if degree < 0 {
		return 0, fmt.Errorf("failed to get atom degree: %s", getLastError())
	}
	return degree, nil
}

// Index returns the index of an atom in its molecule
func (a *Atom) Index() (int, error) {
	index := int(C.indigoAtomIndex(C.int(a.handle)))
	if index < 0 {
		return 0, fmt.Errorf("failed to get atom index: %s", getLastError())
	}
	return index, nil
}

// AtomicNumber returns the atomic number of an atom
func (a *Atom) AtomicNumber() (int, error) {
	number := int(C.indigoAtomicNumber(C.int(a.handle)))
	if number < 0 {
		return 0, fmt.Errorf("failed to get atomic number: %s", getLastError())
	}
	return number, nil
}

// Charge returns the charge of an atom
func (a *Atom) Charge() (int, error) {
	var charge C.int
	ret := int(C.indigoGetCharge(C.int(a.handle), &charge))
	if ret < 0 {
		return 0, fmt.Errorf("failed to get atom charge: %s", getLastError())
	}
	return int(charge), nil
}

// SetCharge sets the charge of an atom
func (a *Atom) SetCharge(charge int) error {
	ret := int(C.indigoSetCharge(C.int(a.handle), C.int(charge)))
	if ret < 0 {
		return fmt.Errorf("failed to set atom charge: %s", getLastError())
	}
	return nil
}

// Isotope returns the isotope of an atom (0 if not set)
func (a *Atom) Isotope() (int, error) {
	isotope := int(C.indigoIsotope(C.int(a.handle)))
	if isotope < 0 {
		return 0, fmt.Errorf("failed to get isotope: %s", getLastError())
	}
	return isotope, nil
}

// SetIsotope sets the isotope of an atom
func (a *Atom) SetIsotope(isotope int) error {
	ret := int(C.indigoSetIsotope(C.int(a.handle), C.int(isotope)))
	if ret < 0 {
		return fmt.Errorf("failed to set isotope: %s", getLastError())
	}
	return nil
}

// Valence returns the valence of an atom
func (a *Atom) Valence() (int, error) {
	valence := int(C.indigoValence(C.int(a.handle)))
	if valence < 0 {
		return 0, fmt.Errorf("failed to get valence: %s", getLastError())
	}
	return valence, nil
}

// ExplicitValence returns the explicit valence of an atom
func (a *Atom) ExplicitValence() (int, error) {
	var valence C.int
	ret := int(C.indigoGetExplicitValence(C.int(a.handle), &valence))
	if ret < 0 {
		return 0, fmt.Errorf("failed to get explicit valence: %s", getLastError())
	}
	return int(valence), nil
}

// SetExplicitValence sets the explicit valence of an atom
func (a *Atom) SetExplicitValence(valence int) error {
	ret := int(C.indigoSetExplicitValence(C.int(a.handle), C.int(valence)))
	if ret < 0 {
		return fmt.Errorf("failed to set explicit valence: %s", getLastError())
	}
	return nil
}

// Radical returns the radical type of an atom
func (a *Atom) Radical() (int, error) {
	var radical C.int
	ret := int(C.indigoGetRadical(C.int(a.handle), &radical))
	if ret < 0 {
		return 0, fmt.Errorf("failed to get radical: %s", getLastError())
	}
	return int(radical), nil
}

// SetRadical sets the radical type of an atom
func (a *Atom) SetRadical(radical int) error {
	ret := int(C.indigoSetRadical(C.int(a.handle), C.int(radical)))
	if ret < 0 {
		return fmt.Errorf("failed to set radical: %s", getLastError())
	}
	return nil
}

// CountImplicitHydrogens returns the number of implicit hydrogens
func (a *Atom) CountImplicitHydrogens() (int, error) {
	count := int(C.indigoCountImplicitHydrogens(C.int(a.handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count implicit hydrogens: %s", getLastError())
	}
	return count, nil
}

// SetImplicitHCount sets the implicit hydrogen count
func (a *Atom) SetImplicitHCount(count int) error {
	ret := int(C.indigoSetImplicitHCount(C.int(a.handle), C.int(count)))
	if ret < 0 {
		return fmt.Errorf("failed to set implicit H count: %s", getLastError())
	}
	return nil
}

// IsPseudoatom checks if an atom is a pseudoatom
func (a *Atom) IsPseudoatom() bool {
	return int(C.indigoIsPseudoatom(C.int(a.handle))) > 0
}

// IsRSite checks if an atom is an R-site
func (a *Atom) IsRSite() bool {
	return int(C.indigoIsRSite(C.int(a.handle))) > 0
}

// BondOrder returns the order of a bond
func BondOrder(bondHandle int) (int, error) {
	order := int(C.indigoBondOrder(C.int(bondHandle)))
	if order < 0 {
		return 0, fmt.Errorf("failed to get bond order: %s", getLastError())
	}
	return order, nil
}

// BondIndex returns the index of a bond
func BondIndex(bondHandle int) (int, error) {
	index := int(C.indigoBondIndex(C.int(bondHandle)))
	if index < 0 {
		return 0, fmt.Errorf("failed to get bond index: %s", getLastError())
	}
	return index, nil
}

// BondSource returns the source atom of a bond
func BondSource(bondHandle int) (int, error) {
	handle := int(C.indigoSource(C.int(bondHandle)))
	if handle < 0 {
		return 0, fmt.Errorf("failed to get bond source: %s", getLastError())
	}
	return handle, nil
}

// BondDestination returns the destination atom of a bond
func BondDestination(bondHandle int) (int, error) {
	handle := int(C.indigoDestination(C.int(bondHandle)))
	if handle < 0 {
		return 0, fmt.Errorf("failed to get bond destination: %s", getLastError())
	}
	return handle, nil
}
