// Package molecule coding=utf-8
// @Project : go-chem
// @File    : lipinski.go
package molecule

// NumRotatableBonds returns a naive count of rotatable single bonds (non-terminal, non-ring approximation)
func NumRotatableBonds(m *Molecule) int {
	if m == nil {
		return 0
	}
	count := 0
	for i := range m.Bonds {
		order := m.GetBondOrder(i)
		if order != BOND_SINGLE {
			continue
		}
		b := m.Bonds[i]
		// terminal check
		if len(m.Vertices[b.Beg].Edges) <= 1 {
			continue
		}
		if len(m.Vertices[b.End].Edges) <= 1 {
			continue
		}
		// rough ring exclusion: if both atoms share two or more common neighbors, skip
		if isLikelyRingEdge(m, b.Beg, b.End) {
			continue
		}
		count++
	}
	return count
}

func isLikelyRingEdge(m *Molecule, u, v int) bool {
	// simple heuristic: check neighbor intersection size
	seen := make(map[int]bool)
	for _, eidx := range m.Vertices[u].Edges {
		seen[otherEnd(m, eidx, u)] = true
	}
	inter := 0
	for _, eidx := range m.Vertices[v].Edges {
		if seen[otherEnd(m, eidx, v)] {
			inter++
		}
	}
	return inter >= 2
}

func otherEnd(m *Molecule, eidx int, u int) int {
	e := m.Bonds[eidx]
	if e.Beg == u {
		return e.End
	}
	if e.End == u {
		return e.Beg
	}
	return u
}

// NumHydrogenBondAcceptors naive definition: O or N with available lone pairs proxy via valence
func NumHydrogenBondAcceptors(m *Molecule) int {
	if m == nil {
		return 0
	}
	c := 0
	for i := range m.Atoms {
		n := m.GetAtomNumber(i)
		if n != ELEM_O && n != ELEM_N {
			continue
		}
		if m.GetAtomCharge(i) > 0 {
			continue
		}
		conn := m.getAtomConnectivityNoImplH(i)
		// proxy: oxygen connectivity <=2, nitrogen <=3
		if (n == ELEM_O && conn <= 2) || (n == ELEM_N && conn <= 3) {
			c++
		}
	}
	return c
}

// NumHydrogenBondDonors naive definition: hydrogens implicitly counted on O or N
func NumHydrogenBondDonors(m *Molecule) int {
	if m == nil {
		return 0
	}
	c := 0
	for i := range m.Atoms {
		n := m.GetAtomNumber(i)
		if n != ELEM_O && n != ELEM_N {
			continue
		}
		if m.GetAtomCharge(i) < 0 {
			continue
		}
		h := m.GetImplicitH(i)
		if h > 0 {
			c++
		}
	}
	return c
}
