// Package src coding=utf-8
// @Project : go-chem
// @File    : tpsa.go
package src

// CalculateTPSA computes a simplified Topological Polar Surface Area.
// This is a placeholder approximation: counts contributions from O and N atoms
// based on element type and simple local environment. It is not a full Indigo TPSA.
func (m *Molecule) CalculateTPSA(includeSP bool) float64 {
	if m == nil {
		return 0
	}
	var total float64
	for atomIdx := range m.Atoms {
		number := m.GetAtomNumber(atomIdx)
		if number != ELEM_O && number != ELEM_N {
			continue
		}
		charge := m.Atoms[atomIdx].Charge
		// simple environment features
		deg := len(m.Vertices[atomIdx].Edges)
		maxOrder := 0
		singleCount, doubleCount, tripleCount, aromaticCount := 0, 0, 0, 0
		for _, eidx := range m.Vertices[atomIdx].Edges {
			o := m.GetBondOrder(eidx)
			if o > maxOrder {
				maxOrder = o
			}
			switch o {
			case BOND_SINGLE:
				singleCount++
			case BOND_DOUBLE:
				doubleCount++
			case BOND_TRIPLE:
				tripleCount++
			case BOND_AROMATIC:
				aromaticCount++
			}
		}

		contrib := 0.0
		if number == ELEM_O {
			// base oxygen contribution
			contrib = 12.0
			if charge > 0 {
				contrib -= 2.0
			}
			if doubleCount > 0 {
				contrib += 2.0
			}
			if aromaticCount > 0 {
				contrib -= 1.0
			}
		} else if number == ELEM_N {
			// base nitrogen contribution
			contrib = 3.0
			if charge > 0 {
				contrib -= 1.0
			}
			if doubleCount > 0 {
				contrib += 0.5
			}
			if aromaticCount > 0 {
				contrib += 0.5
			}
		}
		if includeSP {
			// slight correction for sp hybridization proxy: deg<=2 and maxOrder>=2
			if deg <= 2 && maxOrder >= BOND_DOUBLE {
				contrib *= 0.9
			}
		}
		if contrib < 0 {
			contrib = 0
		}
		total += contrib
	}
	return total
}
