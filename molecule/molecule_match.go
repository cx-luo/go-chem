// Package molecule provides substructure matching functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_match.go
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

// SubstructureMatch represents a substructure match
type SubstructureMatch struct {
	handle int
}

// SubstructureMatcher performs substructure matching
func (m *Molecule) SubstructureMatcher(query *Molecule) (*SubstructureMatch, error) {
	if m.Closed {
		return nil, fmt.Errorf("molecule is closed")
	}
	if query.Closed {
		return nil, fmt.Errorf("query molecule is closed")
	}

	handle := int(C.indigoSubstructureMatcher(C.int(m.Handle), C.CString("substructure")))
	if handle < 0 {
		return nil, fmt.Errorf("failed to create substructure matcher: %s", getLastError())
	}

	matchHandle := int(C.indigoMatch(C.int(handle), C.int(query.Handle)))
	if matchHandle < 0 {
		return nil, fmt.Errorf("no match found")
	}

	return &SubstructureMatch{handle: matchHandle}, nil
}

// CountSubstructureMatches counts the number of substructure matches
// modeStr is the matching mode, "" or NORMAL, RES or RESONANCE, TAU..., defaults to ""
func (m *Molecule) CountSubstructureMatches(queryMolecule *Molecule, modeStr *string) (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}
	if queryMolecule.Closed {
		return 0, fmt.Errorf("queryMolecule molecule is closed")
	}

	// Prepare C string only if modeStr is provided and non-empty.
	var cMode *C.char
	if modeStr != nil && *modeStr != "" {
		cMode = C.CString(*modeStr)
		// ensure the C string is freed after the call
		defer C.free(unsafe.Pointer(cMode))
	} else {
		cMode = nil
	}

	matcherHandle := int(C.indigoSubstructureMatcher(C.int(m.Handle), cMode))
	if matcherHandle < 0 {
		return 0, fmt.Errorf("failed to create substructure matcher: %s", getLastError())
	}
	defer C.indigoFree(C.int(matcherHandle))

	count := int(C.indigoCountMatches(C.int(matcherHandle), C.int(queryMolecule.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count matches: %s", getLastError())
	}

	return count, nil
}

// HasSubstructure checks if the molecule contains the given substructure
// modeStr is the matching mode, NORMAL or RESONANCE or TAUTOMER
func (m *Molecule) HasSubstructure(queryMolecule *Molecule, modeStr *string) (bool, error) {
	if m.Closed {
		return false, fmt.Errorf("molecule is closed")
	}
	if queryMolecule.Closed {
		return false, fmt.Errorf("queryMolecule molecule is closed")
	}

	// Prepare C string only if modeStr is provided and non-empty.
	var cMode *C.char
	if modeStr != nil && *modeStr != "" {
		cMode = C.CString(*modeStr)
		// ensure the C string is freed after the call
		defer C.free(unsafe.Pointer(cMode))
	} else {
		cMode = nil
	}
	matcherHandle := int(C.indigoSubstructureMatcher(C.int(m.Handle), cMode))
	if matcherHandle < 0 {
		return false, fmt.Errorf("failed to create substructure matcher: %s", getLastError())
	}
	defer C.indigoFree(C.int(matcherHandle))

	count := int(C.indigoCountMatchesWithLimit(C.int(matcherHandle), C.int(queryMolecule.Handle), 1))
	return count > 0, nil
}

// ExactMatch checks if the molecule exactly matches another molecule.
// Returns (matched, mappingId, error). mappingId is >0 when matched.
func (m *Molecule) ExactMatch(other *Molecule, flags *string) (bool, int, error) {
	if m.Closed {
		return false, 0, fmt.Errorf("molecule is closed")
	}
	if other.Closed {
		return false, 0, fmt.Errorf("other molecule is closed")
	}

	var cFlags *C.char
	if flags != nil && *flags != "" {
		cFlags = C.CString(*flags)
		defer C.free(unsafe.Pointer(cFlags))
	} else {
		cFlags = nil
	}

	match := int(C.indigoExactMatch(C.int(m.Handle), C.int(other.Handle), cFlags))

	if match < 0 {
		return false, 0, fmt.Errorf("indigo error: %s", getLastError())
	}
	if match == 0 {
		return false, 0, nil
	}
	// match > 0
	return true, match, nil
}

// IterateSubstructureMatches iterates through all substructure matches
func (m *Molecule) IterateSubstructureMatches(query *Molecule, modeStr *string) (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}
	if query.Closed {
		return 0, fmt.Errorf("query molecule is closed")
	}

	// Prepare C string only if modeStr is provided and non-empty.
	var cMode *C.char
	if modeStr != nil && *modeStr != "" {
		cMode = C.CString(*modeStr)
		// ensure the C string is freed after the call
		defer C.free(unsafe.Pointer(cMode))
	} else {
		cMode = nil
	}

	matcherHandle := int(C.indigoSubstructureMatcher(C.int(m.Handle), cMode))
	if matcherHandle < 0 {
		return 0, fmt.Errorf("failed to create substructure matcher: %s", getLastError())
	}

	iterHandle := int(C.indigoIterateMatches(C.int(matcherHandle), C.int(query.Handle)))
	if iterHandle < 0 {
		C.indigoFree(C.int(matcherHandle))
		return 0, fmt.Errorf("failed to iterate matches: %s", getLastError())
	}

	return iterHandle, nil
}

// Highlight highlights atoms and bonds from a match
func (m *Molecule) Highlight(match *SubstructureMatch) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoHighlightedTarget(C.int(match.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to highlight match: %s", getLastError())
	}

	return nil
}

// UnhighlightAll removes all highlights from the molecule
func (m *Molecule) UnhighlightAll() error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoUnhighlight(C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to unhighlight: %s", getLastError())
	}

	return nil
}

// RemoveAtoms removes specified atoms from the molecule
func (m *Molecule) RemoveAtoms(atomIndices []int) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	if len(atomIndices) == 0 {
		return nil
	}

	// Convert to C array
	cIndices := make([]C.int, len(atomIndices))
	for i, idx := range atomIndices {
		cIndices[i] = C.int(idx)
	}

	ret := int(C.indigoRemoveAtoms(C.int(m.Handle), C.int(len(atomIndices)), &cIndices[0]))
	if ret < 0 {
		return fmt.Errorf("failed to remove atoms: %s", getLastError())
	}

	return nil
}

// RemoveBonds removes specified bonds from the molecule
func (m *Molecule) RemoveBonds(bondIndices []int) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	if len(bondIndices) == 0 {
		return nil
	}

	// Convert to C array
	cIndices := make([]C.int, len(bondIndices))
	for i, idx := range bondIndices {
		cIndices[i] = C.int(idx)
	}

	ret := int(C.indigoRemoveBonds(C.int(m.Handle), C.int(len(bondIndices)), &cIndices[0]))
	if ret < 0 {
		return fmt.Errorf("failed to remove bonds: %s", getLastError())
	}

	return nil
}

// GetSubmolecule returns a submolecule containing only specified atoms
func (m *Molecule) GetSubmolecule(atomIndices []int) (*Molecule, error) {
	if m.Closed {
		return nil, fmt.Errorf("molecule is closed")
	}

	if len(atomIndices) == 0 {
		return nil, fmt.Errorf("no atoms specified")
	}

	// Convert to C array
	cIndices := make([]C.int, len(atomIndices))
	for i, idx := range atomIndices {
		cIndices[i] = C.int(idx)
	}

	handle := int(C.indigoGetSubmolecule(C.int(m.Handle), C.int(len(atomIndices)), &cIndices[0]))
	if handle < 0 {
		return nil, fmt.Errorf("failed to get submolecule: %s", getLastError())
	}

	return newMolecule(handle), nil
}
