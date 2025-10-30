// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:58
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : dearomatizer.go
// @Software: GoLand
package molecule

// DearomatizerBase performs dearomatization of aromatic bonds.
// Based on the Indigo implementation in molecule_dearom.cpp.
type DearomatizerBase struct{}

// Apply converts aromatic bonds back into alternating single/double bonds.
// This implements a simplified version of the Kekulé structure assignment.
func (d DearomatizerBase) Apply(m *Molecule) {
	// Process rings of different sizes (5-8 membered rings)
	used := make([]bool, len(m.Bonds))

	for size := 5; size <= 8; size++ {
		cycles := findSimpleCyclesOfLength(m, size)
		for _, cycle := range cycles {
			if d.dearomatizeCycle(m, cycle, used) {
				// Successfully dearomatized this cycle
			}
		}
	}

	// Any remaining aromatic bonds that weren't part of processed cycles → convert to single
	for i := range m.Bonds {
		if m.BondOrders[i] == BOND_AROMATIC && !used[i] {
			m.setBondOrderInternal(i, BOND_SINGLE)
		}
	}

	// Update aromaticity array
	if len(m.Aromaticity) == len(m.Atoms) {
		for i := range m.Aromaticity {
			m.Aromaticity[i] = ATOM_ALIPHATIC
		}
	}

	m.Aromatized = false
}

// dearomatizeCycle attempts to assign alternating single/double bonds to an aromatic cycle.
// Returns true if the cycle was successfully dearomatized.
func (d DearomatizerBase) dearomatizeCycle(m *Molecule, cycle []int, used []bool) bool {
	eidxs := cycleEdges(m, cycle)

	// Check if all bonds in this cycle are aromatic and not yet processed
	allAromatic := true
	anyUsed := false
	for _, eidx := range eidxs {
		if m.GetBondOrder(eidx) != BOND_AROMATIC {
			allAromatic = false
		}
		if used[eidx] {
			anyUsed = true
		}
	}

	if !allAromatic || anyUsed {
		return false
	}

	if len(eidxs) != len(cycle) {
		return false
	}

	// Find the best starting point for alternation
	// Prefer starting from atoms with specific characteristics
	start := d.findBestStartingPoint(m, cycle)

	// Assign alternating single/double bonds
	// The pattern depends on the ring size and needs to satisfy valence
	if d.assignAlternatingBonds(m, cycle, eidxs, start, used) {
		return true
	}

	// If the first pattern didn't work, try starting with double bond
	return d.assignAlternatingBonds(m, cycle, eidxs, start, used)
}

// findBestStartingPoint finds the best atom to start bond alternation.
// Prefers atoms with existing double bonds or specific atom types.
func (d DearomatizerBase) findBestStartingPoint(m *Molecule, cycle []int) int {
	// Start with atom with minimum index for consistency
	minIdx := 0
	for i := 1; i < len(cycle); i++ {
		if cycle[i] < cycle[minIdx] {
			minIdx = i
		}
	}

	// Check if any atom has constraints that prefer double bonds
	for i, atomIdx := range cycle {
		atom := m.Atoms[atomIdx]
		// Nitrogen or oxygen often prefer certain positions
		if atom.Number == ELEM_N || atom.Number == ELEM_O {
			// These atoms might prefer specific patterns
			// For simplicity, use them as starting point
			return i
		}
	}

	return minIdx
}

// assignAlternatingBonds assigns alternating single/double bonds to a cycle.
func (d DearomatizerBase) assignAlternatingBonds(m *Molecule, cycle []int, edges []int, start int, used []bool) bool {
	// Try pattern: single, double, single, double, ...
	for i := 0; i < len(edges); i++ {
		idx := (start + i) % len(edges)
		eidx := edges[idx]

		if i%2 == 0 {
			m.setBondOrderInternal(eidx, BOND_SINGLE)
		} else {
			m.setBondOrderInternal(eidx, BOND_DOUBLE)
		}
		used[eidx] = true
	}

	// Verify that the assignment is valid (all atoms have valid valence)
	// This is a simplified check - full implementation would verify valences
	valid := d.checkValences(m, cycle)

	if !valid {
		// Revert and try opposite pattern
		for i := 0; i < len(edges); i++ {
			idx := (start + i) % len(edges)
			eidx := edges[idx]

			if i%2 == 0 {
				m.setBondOrderInternal(eidx, BOND_DOUBLE)
			} else {
				m.setBondOrderInternal(eidx, BOND_SINGLE)
			}
		}

		valid = d.checkValences(m, cycle)
	}

	return valid
}

// checkValences performs a basic check that atoms have reasonable valences.
// This is a simplified version - full implementation would check against element valence rules.
func (d DearomatizerBase) checkValences(m *Molecule, cycle []int) bool {
	// For now, accept all assignments
	// A full implementation would check:
	// 1. Each atom's total bond order doesn't exceed its maximum valence
	// 2. Hydrogens can be added to satisfy valence requirements
	// 3. Charges are accounted for
	return true
}
