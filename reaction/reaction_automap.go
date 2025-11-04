// Package reaction provides reaction atom-to-atom mapping functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_automap.go
// @Software: GoLand
package reaction

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// Windows platforms
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo

// Linux platforms
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-aarch64
#cgo linux LDFLAGS: -L${SRCDIR}/../3rd/linux -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux

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

// Automap modes
const (
	// AutomapModeDiscard discards the existing mapping entirely and considers only the existing reaction centers (the default)
	AutomapModeDiscard = "discard"

	// AutomapModeKeep keeps the existing mapping and maps unmapped atoms
	AutomapModeKeep = "keep"

	// AutomapModeAlter alters the existing mapping, and maps the rest of the reaction but may change the existing mapping
	AutomapModeAlter = "alter"

	// AutomapModeClear removes the mapping from the reaction
	AutomapModeClear = "clear"
)

// Automap options
const (
	// AutomapIgnoreCharges - do not consider atom charges while searching
	AutomapIgnoreCharges = "ignore_charges"

	// AutomapIgnoreIsotopes - do not consider atom isotopes while searching
	AutomapIgnoreIsotopes = "ignore_isotopes"

	// AutomapIgnoreValence - do not consider atom valence while searching
	AutomapIgnoreValence = "ignore_valence"

	// AutomapIgnoreRadicals - do not consider atom radicals while searching
	AutomapIgnoreRadicals = "ignore_radicals"
)

// Automap performs automatic atom-to-atom mapping for the reaction
// mode can be one or more of the following (separated by space):
//   - "discard": discards the existing mapping entirely (default)
//   - "keep": keeps the existing mapping and maps unmapped atoms
//   - "alter": alters the existing mapping
//   - "clear": removes the mapping from the reaction
//   - "ignore_charges": do not consider atom charges while searching
//   - "ignore_isotopes": do not consider atom isotopes while searching
//   - "ignore_valence": do not consider atom valence while searching
//   - "ignore_radicals": do not consider atom radicals while searching
func (r *Reaction) Automap(mode string) error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	cMode := C.CString(mode)
	defer C.free(unsafe.Pointer(cMode))

	ret := int(C.indigoAutomap(C.int(r.handle), cMode))
	if ret < 0 {
		return fmt.Errorf("failed to automap reaction: %s", getLastError())
	}

	return nil
}

// GetAtomMappingNumber returns the atom-to-atom mapping number for a reaction atom
// Returns 0 if no mapping number has been specified
func (r *Reaction) GetAtomMappingNumber(atomHandle int) (int, error) {
	if r.closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	number := int(C.indigoGetAtomMappingNumber(C.int(r.handle), C.int(atomHandle)))
	if number < 0 {
		return 0, fmt.Errorf("failed to get atom mapping number: %s", getLastError())
	}

	return number, nil
}

// SetAtomMappingNumber sets the atom-to-atom mapping number for a reaction atom
func (r *Reaction) SetAtomMappingNumber(atomHandle int, number int) error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoSetAtomMappingNumber(C.int(r.handle), C.int(atomHandle), C.int(number)))
	if ret < 0 {
		return fmt.Errorf("failed to set atom mapping number: %s", getLastError())
	}

	return nil
}

// GetReactingCenter returns the reacting center information for a bond
// Returns the reacting center flags (combination of RC_* constants)
func (r *Reaction) GetReactingCenter(bondHandle int) (int, error) {
	if r.closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	var rc C.int
	ret := int(C.indigoGetReactingCenter(C.int(r.handle), C.int(bondHandle), &rc))
	if ret < 0 {
		return 0, fmt.Errorf("failed to get reacting center: %s", getLastError())
	}

	return int(rc), nil
}

// SetReactingCenter sets the reacting center information for a bond
// rc should be a combination of RC_* constants
func (r *Reaction) SetReactingCenter(bondHandle int, rc int) error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoSetReactingCenter(C.int(r.handle), C.int(bondHandle), C.int(rc)))
	if ret < 0 {
		return fmt.Errorf("failed to set reacting center: %s", getLastError())
	}

	return nil
}

// ClearAAM clears all atom-to-atom mapping information from the reaction
func (r *Reaction) ClearAAM() error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoClearAAM(C.int(r.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to clear AAM: %s", getLastError())
	}

	return nil
}

// CorrectReactingCenters corrects reacting centers according to atom-to-atom mapping
func (r *Reaction) CorrectReactingCenters() error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoCorrectReactingCenters(C.int(r.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to correct reacting centers: %s", getLastError())
	}

	return nil
}
