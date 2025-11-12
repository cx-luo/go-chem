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
func (m *Molecule) CountSubstructureMatches(query *Molecule) (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}
	if query.Closed {
		return 0, fmt.Errorf("query molecule is closed")
	}

	matcherHandle := int(C.indigoSubstructureMatcher(C.int(m.Handle), C.CString("substructure")))
	if matcherHandle < 0 {
		return 0, fmt.Errorf("failed to create substructure matcher: %s", getLastError())
	}
	defer C.indigoFree(C.int(matcherHandle))

	count := int(C.indigoCountMatches(C.int(matcherHandle), C.int(query.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count matches: %s", getLastError())
	}

	return count, nil
}

// HasSubstructure checks if the molecule contains the given substructure
func (m *Molecule) HasSubstructure(query *Molecule) (bool, error) {
	count, err := m.CountSubstructureMatches(query)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExactMatch checks if the molecule exactly matches another molecule
func (m *Molecule) ExactMatch(other *Molecule) (bool, error) {
	if m.Closed {
		return false, fmt.Errorf("molecule is closed")
	}
	if other.Closed {
		return false, fmt.Errorf("other molecule is closed")
	}

	matcherHandle := int(C.indigoSubstructureMatcher(C.int(m.Handle), C.CString("exact")))
	if matcherHandle < 0 {
		return false, fmt.Errorf("failed to create exact matcher: %s", getLastError())
	}
	defer C.indigoFree(C.int(matcherHandle))

	matchHandle := int(C.indigoMatch(C.int(matcherHandle), C.int(other.Handle)))
	return matchHandle >= 0, nil
}

// IterateSubstructureMatches iterates through all substructure matches
func (m *Molecule) IterateSubstructureMatches(query *Molecule) (int, error) {
	if m.Closed {
		return 0, fmt.Errorf("molecule is closed")
	}
	if query.Closed {
		return 0, fmt.Errorf("query molecule is closed")
	}

	matcherHandle := int(C.indigoSubstructureMatcher(C.int(m.Handle), C.CString("substructure")))
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
