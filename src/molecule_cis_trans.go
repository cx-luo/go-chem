// Package src provides molecular structure manipulation and analysis tools.
// This file implements cis/trans (E/Z) stereochemistry handling.
package src

import (
	"fmt"
)

// Cis/Trans configuration constants
const (
	CIS   = 1 // Cis configuration (Z)
	TRANS = 2 // Trans configuration (E)
)

// CisTrans represents a cis/trans stereochemical bond
type CisTrans struct {
	BondIdx      int    // Index of the double bond
	Parity       int    // CIS or TRANS
	Substituents [4]int // Four substituents around the double bond
	Ignored      bool   // Whether this configuration is explicitly ignored
}

// MoleculeCisTrans manages cis/trans stereochemistry in a molecule
type MoleculeCisTrans struct {
	bonds map[int]*CisTrans // Map from bond index to cis/trans info
}

// NewMoleculeCisTrans creates a new empty cis/trans collection
func NewMoleculeCisTrans() *MoleculeCisTrans {
	return &MoleculeCisTrans{
		bonds: make(map[int]*CisTrans),
	}
}

// Clear removes all cis/trans information
func (mct *MoleculeCisTrans) Clear() {
	mct.bonds = make(map[int]*CisTrans)
}

// Exists checks if a bond has cis/trans information
func (mct *MoleculeCisTrans) Exists() bool {
	return len(mct.bonds) > 0
}

// Count returns the number of cis/trans bonds
func (mct *MoleculeCisTrans) Count() int {
	return len(mct.bonds)
}

// SetParity sets the parity (CIS or TRANS) for a bond
func (mct *MoleculeCisTrans) SetParity(bondIdx, parity int) {
	if ct, ok := mct.bonds[bondIdx]; ok {
		ct.Parity = parity
	}
}

// GetParity returns the parity of a bond
func (mct *MoleculeCisTrans) GetParity(bondIdx int) int {
	if ct, ok := mct.bonds[bondIdx]; ok {
		return ct.Parity
	}
	return 0
}

// IsIgnored checks if a bond's cis/trans configuration is explicitly ignored
func (mct *MoleculeCisTrans) IsIgnored(bondIdx int) bool {
	if ct, ok := mct.bonds[bondIdx]; ok {
		return ct.Ignored
	}
	return false
}

// Ignore marks a bond's cis/trans configuration as ignored
func (mct *MoleculeCisTrans) Ignore(bondIdx int) {
	if ct, ok := mct.bonds[bondIdx]; ok {
		ct.Ignored = true
	} else {
		mct.bonds[bondIdx] = &CisTrans{
			BondIdx:      bondIdx,
			Parity:       0,
			Substituents: [4]int{-1, -1, -1, -1},
			Ignored:      true,
		}
	}
}

// RegisterBond registers a bond as potentially having cis/trans stereochemistry
func (mct *MoleculeCisTrans) RegisterBond(bondIdx int) {
	if _, ok := mct.bonds[bondIdx]; !ok {
		mct.bonds[bondIdx] = &CisTrans{
			BondIdx:      bondIdx,
			Parity:       0,
			Substituents: [4]int{-1, -1, -1, -1},
			Ignored:      false,
		}
	}
}

// Add adds cis/trans information for a bond
func (mct *MoleculeCisTrans) Add(bondIdx int, substituents [4]int, parity int) {
	mct.bonds[bondIdx] = &CisTrans{
		BondIdx:      bondIdx,
		Parity:       parity,
		Substituents: substituents,
		Ignored:      false,
	}
}

// GetSubstituents returns the four substituents around a double bond
func (mct *MoleculeCisTrans) GetSubstituents(bondIdx int) [4]int {
	if ct, ok := mct.bonds[bondIdx]; ok {
		return ct.Substituents
	}
	return [4]int{-1, -1, -1, -1}
}

// GetSubstituentsAll retrieves all substituents around a double bond
func (mct *MoleculeCisTrans) GetSubstituentsAll(mol *Molecule, bondIdx int) [4]int {
	if bondIdx < 0 || bondIdx >= len(mol.Bonds) {
		return [4]int{-1, -1, -1, -1}
	}

	bond := mol.Bonds[bondIdx]
	subst := [4]int{-1, -1, -1, -1}

	// Get neighbors of begin atom
	neighbors1 := mol.GetNeighbors(bond.Beg)
	idx := 0
	for _, n := range neighbors1 {
		if n != bond.End && idx < 2 {
			subst[idx] = n
			idx++
		}
	}

	// Get neighbors of end atom
	neighbors2 := mol.GetNeighbors(bond.End)
	idx = 2
	for _, n := range neighbors2 {
		if n != bond.Beg && idx < 4 {
			subst[idx] = n
			idx++
		}
	}

	return subst
}

// Build determines cis/trans bonds from the molecule structure
func (mct *MoleculeCisTrans) Build(mol *Molecule, excludeBonds []int) {
	mct.Clear()

	exclude := make(map[int]bool)
	for _, b := range excludeBonds {
		exclude[b] = true
	}

	// Look for double bonds that could have cis/trans stereochemistry
	for i, bond := range mol.Bonds {
		if exclude[i] {
			continue
		}

		// Only double bonds can have cis/trans
		if bond.Order != BOND_DOUBLE {
			continue
		}

		// Check if this is a stereogenic double bond
		if mct.IsGeomStereoBond(mol, i) {
			mct.RegisterBondAndSubstituents(mol, i)
		}
	}
}

// IsGeomStereoBond checks if a double bond can have cis/trans stereochemistry
func (mct *MoleculeCisTrans) IsGeomStereoBond(mol *Molecule, bondIdx int) bool {
	if bondIdx < 0 || bondIdx >= len(mol.Bonds) {
		return false
	}

	bond := mol.Bonds[bondIdx]

	// Must be a double bond
	if bond.Order != BOND_DOUBLE {
		return false
	}

	// Get neighbors of both ends
	neighbors1 := mol.GetNeighbors(bond.Beg)
	neighbors2 := mol.GetNeighbors(bond.End)

	// Each end must have exactly 2 or 3 neighbors (including the double bond partner)
	if len(neighbors1) < 2 || len(neighbors1) > 3 || len(neighbors2) < 2 || len(neighbors2) > 3 {
		return false
	}

	// Count non-double-bond neighbors
	nonDbNeighbors1 := 0
	nonDbNeighbors2 := 0

	for _, n := range neighbors1 {
		if n != bond.End {
			nonDbNeighbors1++
		}
	}

	for _, n := range neighbors2 {
		if n != bond.Beg {
			nonDbNeighbors2++
		}
	}

	// Need at least one substituent on each end (not counting the double bond itself)
	// And can't have more than 2 (or it's sp3, not sp2)
	if nonDbNeighbors1 < 1 || nonDbNeighbors1 > 2 || nonDbNeighbors2 < 1 || nonDbNeighbors2 > 2 {
		return false
	}

	// Check that substituents are different (simplified check)
	return mct.hasDistinctSubstituents(mol, bondIdx)
}

// hasDistinctSubstituents checks if the substituents are different enough for stereochemistry
func (mct *MoleculeCisTrans) hasDistinctSubstituents(mol *Molecule, bondIdx int) bool {
	bond := mol.Bonds[bondIdx]

	// Get substituents on each side
	neighbors1 := mol.GetNeighbors(bond.Beg)
	neighbors2 := mol.GetNeighbors(bond.End)

	subs1 := []int{}
	subs2 := []int{}

	for _, n := range neighbors1 {
		if n != bond.End {
			subs1 = append(subs1, n)
		}
	}

	for _, n := range neighbors2 {
		if n != bond.Beg {
			subs2 = append(subs2, n)
		}
	}

	// If we have 2 substituents on each side, check if they're different
	if len(subs1) == 2 {
		if mol.Atoms[subs1[0]].Number == mol.Atoms[subs1[1]].Number {
			// Both substituents on this side are the same element - not stereogenic
			return false
		}
	}

	if len(subs2) == 2 {
		if mol.Atoms[subs2[0]].Number == mol.Atoms[subs2[1]].Number {
			// Both substituents on this side are the same element - not stereogenic
			return false
		}
	}

	return true
}

// RegisterBondAndSubstituents registers a bond and determines its substituents
func (mct *MoleculeCisTrans) RegisterBondAndSubstituents(mol *Molecule, bondIdx int) bool {
	if bondIdx < 0 || bondIdx >= len(mol.Bonds) {
		return false
	}

	subst := mct.GetSubstituentsAll(mol, bondIdx)

	// Determine parity from bond directions if available
	parity := mct.determineParity(mol, bondIdx, subst)

	mct.Add(bondIdx, subst, parity)
	return true
}

// determineParity determines CIS or TRANS parity from bond directions
func (mct *MoleculeCisTrans) determineParity(mol *Molecule, bondIdx int, subst [4]int) int {
	// This is a simplified implementation
	// A full implementation would analyze bond directions and 3D coordinates

	// Check for explicit bond directions
	bond := mol.Bonds[bondIdx]

	// Check if we have 2D or 3D coordinates to determine parity
	if mol.HaveXYZ {
		return mct.determineParityFrom3D(mol, bondIdx, subst)
	}

	// Check bond directions
	dir1 := -1
	dir2 := -1

	for _, edgeIdx := range mol.GetNeighborBonds(bond.Beg) {
		if edgeIdx != bondIdx {
			dir := mol.GetBondDirection(edgeIdx)
			if dir != 0 {
				dir1 = dir
				break
			}
		}
	}

	for _, edgeIdx := range mol.GetNeighborBonds(bond.End) {
		if edgeIdx != bondIdx {
			dir := mol.GetBondDirection(edgeIdx)
			if dir != 0 {
				dir2 = dir
				break
			}
		}
	}

	// If we have directions on both sides, determine parity
	if dir1 > 0 && dir2 > 0 {
		if dir1 == dir2 {
			return CIS
		}
		return TRANS
	}

	// Can't determine parity
	return 0
}

// determineParityFrom3D determines parity from 3D coordinates
func (mct *MoleculeCisTrans) determineParityFrom3D(mol *Molecule, bondIdx int, subst [4]int) int {
	if !mol.HaveXYZ {
		return 0
	}

	// Get the double bond
	bond := mol.Bonds[bondIdx]

	// Need all 4 substituents to be valid
	if subst[0] < 0 || subst[1] < 0 || subst[2] < 0 || subst[3] < 0 {
		return 0
	}

	// Get positions
	pos1 := mol.Atoms[bond.Beg].Pos
	pos2 := mol.Atoms[bond.End].Pos
	subPos1 := mol.Atoms[subst[0]].Pos
	subPos2 := mol.Atoms[subst[2]].Pos

	// Calculate if substituents are on same side
	bondVec := subVec3f(pos2, pos1)
	sub1Vec := subVec3f(subPos1, pos1)
	sub2Vec := subVec3f(subPos2, pos2)

	// Cross products to determine sides
	cross1 := crossVec3f(bondVec, sub1Vec)
	dot := dotVec3f(cross1, sub2Vec)

	// If dot product is positive, they're on opposite sides (trans)
	// If negative, same side (cis)
	if dot > 0.01 {
		return TRANS
	} else if dot < -0.01 {
		return CIS
	}

	return 0
}

// FlipBond updates cis/trans info when a bond is flipped
func (mct *MoleculeCisTrans) FlipBond(mol *Molecule, atomParent, atomFrom, atomTo int) {
	// Find bonds involving these atoms that might have cis/trans info
	for bondIdx, ct := range mct.bonds {
		bond := mol.Bonds[bondIdx]

		// Check if this bond involves the flipped connection
		if bond.Beg == atomParent || bond.End == atomParent {
			// Update substituents
			for i := range ct.Substituents {
				if ct.Substituents[i] == atomFrom {
					ct.Substituents[i] = atomTo
				}
			}
		}
	}
}

// Validate checks and removes invalid cis/trans configurations
func (mct *MoleculeCisTrans) Validate(mol *Molecule) {
	toRemove := []int{}

	for bondIdx := range mct.bonds {
		if !mct.IsGeomStereoBond(mol, bondIdx) {
			toRemove = append(toRemove, bondIdx)
		}
	}

	for _, bondIdx := range toRemove {
		delete(mct.bonds, bondIdx)
	}
}

// GetAllBonds returns indices of all cis/trans bonds
func (mct *MoleculeCisTrans) GetAllBonds() []int {
	result := make([]int, 0, len(mct.bonds))
	for bondIdx := range mct.bonds {
		result = append(result, bondIdx)
	}
	return result
}

// String returns a string representation of the cis/trans configuration
func (mct *MoleculeCisTrans) String(bondIdx int) string {
	if ct, ok := mct.bonds[bondIdx]; ok {
		if ct.Ignored {
			return "ignored"
		}
		switch ct.Parity {
		case CIS:
			return "cis (Z)"
		case TRANS:
			return "trans (E)"
		default:
			return "undefined"
		}
	}
	return "none"
}

// CisTransError Error type for cis/trans operations
type CisTransError struct {
	Message string
}

func (e *CisTransError) Error() string {
	return fmt.Sprintf("cis/trans error: %s", e.Message)
}
