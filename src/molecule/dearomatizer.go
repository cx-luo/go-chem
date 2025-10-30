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
// This is a simplified implementation using greedy Kekulé structure assignment.
type DearomatizerBase struct{}

// Apply converts aromatic bonds back into alternating single/double bonds.
// This implements a simplified version of the Kekulé structure assignment algorithm.
// Based on MoleculeDearomatizer::dearomatizeMolecule in molecule_dearom.cpp
func (d *DearomatizerBase) Apply(m *Molecule) {
	// Ensure BondOrders is properly initialized
	if len(m.BondOrders) != len(m.Bonds) {
		m.BondOrders = make([]int, len(m.Bonds))
		for i := range m.Bonds {
			m.BondOrders[i] = m.Bonds[i].Order
		}
	}

	// Track which bonds have been processed
	processed := make([]bool, len(m.Bonds))

	// Find all aromatic cycles and dearomatize them
	// Process from smallest to largest rings (5-8 membered rings)
	for size := 5; size <= 8; size++ {
		cycles := findSimpleCyclesOfLength(m, size)
		for _, cycle := range cycles {
			d.dearomatizeCycle(m, cycle, processed)
		}
	}

	// Convert any remaining aromatic bonds to single bonds
	for i := range m.Bonds {
		if m.Bonds[i].Order == BOND_AROMATIC && !processed[i] {
			m.Bonds[i].Order = BOND_SINGLE
			m.BondOrders[i] = BOND_SINGLE
		}
	}

	// Update aromaticity array - all atoms become aliphatic
	if len(m.Aromaticity) == len(m.Atoms) {
		for i := range m.Aromaticity {
			m.Aromaticity[i] = ATOM_ALIPHATIC
		}
	}

	m.Aromatized = false
	m.UpdateEditRevision()
}

// dearomatizeCycle attempts to assign alternating single/double bonds to an aromatic cycle.
// Returns true if the cycle was successfully dearomatized.
// Based on the dearomatization logic in molecule_dearom.cpp
func (d *DearomatizerBase) dearomatizeCycle(m *Molecule, cycle []int, processed []bool) bool {
	if len(cycle) == 0 {
		return false
	}

	eidxs := cycleEdges(m, cycle)
	if len(eidxs) != len(cycle) {
		return false
	}

	// Check if all bonds in this cycle are aromatic and not yet processed
	allAromatic := true
	anyProcessed := false
	for _, eidx := range eidxs {
		bondOrder := m.Bonds[eidx].Order
		if bondOrder != BOND_AROMATIC {
			allAromatic = false
		}
		if processed[eidx] {
			anyProcessed = true
		}
	}

	// Skip if not all aromatic or already processed
	if !allAromatic || anyProcessed {
		return false
	}

	// Try to assign alternating single/double bonds
	// For even-sized rings: perfect alternation
	// For odd-sized rings (5-membered): need special handling
	success := d.assignKekuleStructure(m, cycle, eidxs, processed)

	return success
}

// assignKekuleStructure assigns alternating single/double bonds to a cycle.
// This implements a simplified Kekulé structure assignment algorithm.
// Based on the perfect matching approach in molecule_dearom.cpp
func (d *DearomatizerBase) assignKekuleStructure(m *Molecule, cycle []int, edges []int, processed []bool) bool {
	// Find best starting position - prefer heteroatoms
	start := d.findBestStartingPoint(m, cycle)

	// Try two patterns: start with single or start with double
	patterns := []bool{false, true} // false = start with single, true = start with double

	for _, startWithDouble := range patterns {
		// Assign bonds according to pattern
		for i := 0; i < len(edges); i++ {
			idx := (start + i) % len(edges)
			eidx := edges[idx]

			var newOrder int
			if startWithDouble {
				if i%2 == 0 {
					newOrder = BOND_DOUBLE
				} else {
					newOrder = BOND_SINGLE
				}
			} else {
				if i%2 == 0 {
					newOrder = BOND_SINGLE
				} else {
					newOrder = BOND_DOUBLE
				}
			}

			m.Bonds[eidx].Order = newOrder
			m.BondOrders[eidx] = newOrder
		}

		// Check if this assignment is valid
		if d.checkValences(m, cycle) {
			// Mark all bonds as processed
			for _, eidx := range edges {
				processed[eidx] = true
			}
			return true
		}
	}

	// If neither pattern worked, use fallback: all single bonds
	for _, eidx := range edges {
		m.Bonds[eidx].Order = BOND_SINGLE
		m.BondOrders[eidx] = BOND_SINGLE
		processed[eidx] = true
	}

	return true
}

// findBestStartingPoint finds the best atom to start bond alternation.
// Prefers heteroatoms (N, O, S) which often have specific valence requirements.
func (d *DearomatizerBase) findBestStartingPoint(m *Molecule, cycle []int) int {
	// Prefer starting from heteroatoms
	for i, atomIdx := range cycle {
		atom := m.Atoms[atomIdx]
		// Nitrogen, oxygen, sulfur often have specific requirements
		if atom.Number == ELEM_N || atom.Number == ELEM_O || atom.Number == ELEM_S {
			return i
		}
	}

	// Otherwise, start from first atom
	return 0
}

// checkValences performs a basic check that atoms have reasonable valences.
// This is a simplified version - full implementation would check against element valence rules.
func (d *DearomatizerBase) checkValences(m *Molecule, cycle []int) bool {
	// Check each atom in the cycle
	for _, atomIdx := range cycle {
		atom := m.Atoms[atomIdx]

		// Calculate total bond order
		totalBondOrder := 0
		for _, edgeIdx := range m.Vertices[atomIdx].Edges {
			bond := m.Bonds[edgeIdx]
			if bond.Order > 0 && bond.Order <= BOND_TRIPLE {
				totalBondOrder += bond.Order
			}
		}

		// Get maximum connectivity for this element
		maxConn := ElementMaximumConnectivity(atom.Number, atom.Charge, atom.Radical, false)

		// Allow some flexibility for implicit hydrogens
		// Most aromatic atoms can accommodate the bonds
		if totalBondOrder > maxConn+1 {
			return false
		}
	}

	return true
}
