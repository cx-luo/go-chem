// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:50
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : aromatizer.go
// @Software: GoLand
package molecule

// AromatizerBase provides an API to aromatize/dearomatize a molecule.
// Based on the Indigo implementation in molecule_arom.cpp, using Huckel's rule (4n+2 π electrons).
type AromatizerBase struct {
	bondsAromaticCount []int // Count of aromatic cycles containing each bond
}

// Aromatize marks bonds in rings that satisfy Huckel's rule as aromatic.
// This checks for 4n+2 π electrons in planar cyclic systems.
// Based on MoleculeAromatizer::aromatizeBonds in molecule_arom.cpp
func (a *AromatizerBase) Aromatize(m *Molecule) {
	// Initialize bond aromatic counts
	a.bondsAromaticCount = make([]int, len(m.Bonds))

	// Precalculate pi labels for all atoms
	piLabels := make([]int, len(m.Atoms))
	for i := range m.Atoms {
		piLabels[i] = a.getPiLabel(m, i)
	}

	// Find all simple cycles (5-8 membered rings are most common aromatic rings)
	for size := 5; size <= 8; size++ {
		cycles := findSimpleCyclesOfLength(m, size)
		for _, cycle := range cycles {
			if a.isCycleAromatic(m, cycle, piLabels) {
				a.aromatizeCycle(m, cycle)
			}
		}
	}

	// Update bond orders based on aromaticity
	for i := range m.Bonds {
		if a.bondsAromaticCount[i] > 0 {
			m.Bonds[i].Order = BOND_AROMATIC
			if len(m.BondOrders) > i {
				m.BondOrders[i] = BOND_AROMATIC
			}
		}
	}

	// Update aromaticity array based on aromatic bonds
	if len(m.Aromaticity) != len(m.Atoms) {
		m.Aromaticity = make([]int, len(m.Atoms))
	}
	for i := range m.Aromaticity {
		m.Aromaticity[i] = ATOM_ALIPHATIC
	}

	for _, bond := range m.Bonds {
		if bond.Order == BOND_AROMATIC {
			m.Aromaticity[bond.Beg] = ATOM_AROMATIC
			m.Aromaticity[bond.End] = ATOM_AROMATIC
		}
	}

	m.Aromatized = true
}

// Dearomatize converts aromatic bonds back to alternating single/double bonds.
func (a AromatizerBase) Dearomatize(m *Molecule) {
	// Let the dearomatizer handle this properly
	d := DearomatizerBase{}
	d.Apply(m)
	m.Aromatized = false
}

// isCycleAromatic checks if a cycle satisfies Huckel's rule (4n+2 π electrons).
// Based on MoleculeAromatizer::_isCycleAromatic in molecule_arom.cpp
func (a *AromatizerBase) isCycleAromatic(m *Molecule, cycle []int, piLabels []int) bool {
	// Check: all atoms must be potentially aromatic (pi label >= 0)
	for _, atomIdx := range cycle {
		if piLabels[atomIdx] < 0 {
			return false
		}
	}

	// Check for double bonds in the cycle - corresponds to _checkDoubleBonds
	if !a.checkDoubleBonds(m, cycle) {
		return false
	}

	// Calculate π electron count using precalculated pi labels
	piCount := 0
	for _, atomIdx := range cycle {
		piCount += piLabels[atomIdx]
	}

	// Huckel's rule: 4n+2 π electrons (where n = 0, 1, 2, ...)
	// This means: (piCount - 2) % 4 == 0
	if ((piCount - 2) % 4) != 0 {
		return false
	}

	return true
}

// checkDoubleBonds verifies that double bonds in the cycle are acceptable for aromaticity.
// Based on AromatizerBase::_checkDoubleBonds in molecule_arom.cpp
func (a *AromatizerBase) checkDoubleBonds(m *Molecule, cycle []int) bool {
	cycleLen := len(cycle)
	for j := 0; j < cycleLen; j++ {
		vCenter := cycle[j]
		vLeft := cycle[(j-1+cycleLen)%cycleLen]
		vRight := cycle[(j+1)%cycleLen]

		internalDoubleCount := 0
		for _, edgeIdx := range m.Vertices[vCenter].Edges {
			bond := m.Bonds[edgeIdx]
			neiIdx := bond.End
			if neiIdx == vCenter {
				neiIdx = bond.Beg
			}

			if bond.Order == BOND_DOUBLE && a.bondsAromaticCount[edgeIdx] == 0 {
				if neiIdx != vLeft && neiIdx != vRight {
					// External double bond - check if acceptable
					if !a.acceptOutgoingDoubleBond(m, vCenter, neiIdx) {
						return false
					}
				} else {
					// Internal double bond
					internalDoubleCount++
				}
			}
		}
		// Cannot have 2+ consecutive double bonds in the cycle
		if internalDoubleCount >= 2 {
			return false
		}
	}
	return true
}

// acceptOutgoingDoubleBond checks if an external double bond is acceptable.
// Based on MoleculeAromatizer::_acceptOutgoingDoubleBond in molecule_arom.cpp
func (a *AromatizerBase) acceptOutgoingDoubleBond(m *Molecule, atom int, neighbor int) bool {
	atomNum := m.Atoms[atom].Number
	neighNum := m.Atoms[neighbor].Number

	// C=O, C=N, C=S are acceptable (atom contributes 0 pi electrons)
	if atomNum == ELEM_C {
		if neighNum == ELEM_N || neighNum == ELEM_O || neighNum == ELEM_S {
			return true
		}
	}

	// S=O is acceptable
	if atomNum == ELEM_S && neighNum == ELEM_O {
		return true
	}

	return false
}

// canBeAromatic checks if an atom can participate in aromatic systems.
// Based on Element::canBeAromatic in C++.
func (a AromatizerBase) canBeAromatic(m *Molecule, atomIdx int) bool {
	atomNum := m.Atoms[atomIdx].Number
	// C, N, O, S, P, As, Se can be aromatic
	switch atomNum {
	case ELEM_C, ELEM_N, ELEM_O, ELEM_S, ELEM_P, ELEM_As, ELEM_Se:
		return true
	}
	return false
}

// getPiLabel calculates the π electron contribution of an atom.
// Returns: 0 (vacant p-orbital), 1 (radical/double bond), 2 (lone pair), -1 (cannot be aromatic)
// Based on MoleculeAromatizer::_getPiLabel in molecule_arom.cpp
func (a *AromatizerBase) getPiLabel(m *Molecule, atomIdx int) int {
	atom := m.Atoms[atomIdx]

	// Pseudo atoms, R-sites, templates cannot be aromatic
	if m.IsPseudoAtom(atomIdx) || m.IsTemplateAtom(atomIdx) {
		return -1
	}

	// Check if element can be aromatic
	if !ElementCanBeAromatic(atom.Number) {
		return -1
	}

	// Count connectivity and identify bond types
	nonAromConn := 0
	aromBonds := 0
	nDoubleTotal := 0
	nDoubleExt := 0

	for _, edgeIdx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[edgeIdx]
		bondOrder := bond.Order

		if bondOrder == BOND_TRIPLE {
			return -1 // Triple bonds cannot participate in aromaticity
		}

		if bondOrder == BOND_AROMATIC {
			aromBonds++
		} else if bondOrder == BOND_DOUBLE {
			nDoubleTotal++
			// Check if this might be an external double bond
			// External double bonds typically go to O, N, S heteroatoms
			neighbor := bond.End
			if neighbor == atomIdx {
				neighbor = bond.Beg
			}
			neighAtom := m.Atoms[neighbor].Number
			// C=O, C=N, C=S, S=O are typically external
			if (atom.Number == ELEM_C && (neighAtom == ELEM_O || neighAtom == ELEM_N || neighAtom == ELEM_S)) ||
				(atom.Number == ELEM_S && neighAtom == ELEM_O) {
				if !a.acceptOutgoingDoubleBond(m, atomIdx, neighbor) {
					return -1
				}
				nDoubleExt++
			}
		} else {
			nonAromConn++
		}
	}

	// If there are double bonds that are not external, they contribute 1 π electron
	nDoubleRing := nDoubleTotal - nDoubleExt
	if nDoubleRing > 0 {
		return 1
	}

	// Multiple external double bonds not allowed
	if nDoubleExt > 1 {
		return -1
	}

	// Single external double bond: atom contributes specific π electrons
	if nDoubleExt == 1 {
		if atom.Number == ELEM_S {
			return 2 // S with external =O can contribute lone pair
		}
		return 0 // Vacant p-orbital (no pi contribution)
	}

	// No double bonds - calculate based on connectivity
	// For proper aromaticity detection, we need to account for implicit hydrogens
	connectivity := nonAromConn + aromBonds

	// If no aromatic bonds and only single bonds, check if atom is saturated
	if aromBonds == 0 && nDoubleTotal == 0 {
		// For carbon, if connectivity is 2 and we could add 2 more H to reach 4,
		// it would be sp3 (saturated), not aromatic
		if atom.Number == ELEM_C && connectivity == 2 {
			// Could be aromatic (sp2) OR saturated (sp3)
			// Need to check if implicit H count suggests sp3
			// In a proper implementation, this would check getImplicitH
			// For now, if no double bonds in the molecule and pure single bonds,
			// assume it's saturated
			hasAnyDoubleBonds := false
			for i := range m.Bonds {
				if m.Bonds[i].Order == BOND_DOUBLE || m.Bonds[i].Order == BOND_AROMATIC {
					hasAnyDoubleBonds = true
					break
				}
			}
			if !hasAnyDoubleBonds {
				return -1 // Saturated ring, not aromatic
			}
		}
	}

	return a.getPiLabelByConnectivity(m, atomIdx, connectivity)
}

// getPiLabelByConnectivity determines π electron contribution based on atom type and connectivity.
// Based on MoleculeAromatizer::_getPiLabelByConn in molecule_arom.cpp
func (a *AromatizerBase) getPiLabelByConnectivity(m *Molecule, atomIdx int, connectivity int) int {
	atom := m.Atoms[atomIdx]

	// Check for radicals - always contributes 1 π electron
	if atom.Radical == RADICAL_DOUBLET || atom.Radical == RADICAL_TRIPLET {
		return 1
	}

	// Get element group for vacant pi orbital / lone pair calculation
	group := ElementGroup(atom.Number)

	switch atom.Number {
	case ELEM_C:
		// Carbon: 4 valence electrons
		// conn=3 → 1 unpaired electron → contributes 1
		if connectivity == 3 {
			return 1
		}
		if connectivity == 2 {
			// Carbene or similar: 2 lone electrons
			return 1
		}
	case ELEM_N:
		// Nitrogen: 5 valence electrons, charge affects behavior
		charge := atom.Charge
		if charge == 0 {
			if connectivity == 2 {
				// =N- (pyridine): 1 π electron
				return 1
			}
			if connectivity == 3 {
				// -N< (pyrrole): lone pair contributes 2
				return 2
			}
		} else if charge == 1 {
			// N+ with 3 bonds: pyridinium (vacant orbital)
			if connectivity == 3 {
				return 0
			}
		}
	case ELEM_O:
		// Oxygen: 6 valence electrons
		if connectivity == 2 {
			// -O- (furan): two lone pairs, contributes 2
			return 2
		}
	case ELEM_S:
		// Sulfur: 6 valence electrons (like oxygen)
		if connectivity == 2 {
			// -S- (thiophene): lone pair contributes 2
			return 2
		}
	case ELEM_P:
		// Phosphorus: 5 valence electrons (like nitrogen)
		if connectivity == 2 {
			return 1
		}
		if connectivity == 3 {
			return 2
		}
	case ELEM_B:
		// Boron: 3 valence electrons, typically vacant p-orbital
		if connectivity == 3 {
			return 0 // Vacant orbital
		}
	case ELEM_As, ELEM_Se:
		// Similar to P and S respectively
		if connectivity == 2 {
			if atom.Number == ELEM_As {
				return 1
			}
			return 2
		}
	}

	// For other cases, use simplified electron counting
	// Group valence electrons minus connectivity
	electrons := group - connectivity - atom.Charge
	if electrons == 1 {
		return 1 // Unpaired electron
	}
	if electrons >= 2 {
		return 2 // Lone pair
	}
	if electrons == 0 {
		return 0 // Vacant orbital
	}

	return -1 // Cannot be aromatic
}

// aromatizeCycle marks all bonds in a cycle as aromatic.
// Based on AromatizerBase::_aromatizeCycle in molecule_arom.cpp
func (a *AromatizerBase) aromatizeCycle(m *Molecule, cycle []int) {
	// Mark all bonds in the cycle by incrementing their aromatic count
	for i := 0; i < len(cycle); i++ {
		u := cycle[i]
		v := cycle[(i+1)%len(cycle)]

		// Find edge between u and v
		for _, eidx := range m.Vertices[u].Edges {
			e := m.Bonds[eidx]
			if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
				a.bondsAromaticCount[eidx]++
				break
			}
		}
	}

	// Also mark single bonds inside the aromatic cycle as aromatic
	// For example, in fused rings like naphthalene
	cycleAtomSet := make(map[int]bool, len(cycle))
	for _, atomIdx := range cycle {
		cycleAtomSet[atomIdx] = true
	}

	for _, atomIdx := range cycle {
		for _, edgeIdx := range m.Vertices[atomIdx].Edges {
			bond := m.Bonds[edgeIdx]
			otherAtom := bond.End
			if otherAtom == atomIdx {
				otherAtom = bond.Beg
			}

			// Check if both ends are in the cycle
			if !cycleAtomSet[otherAtom] {
				continue
			}

			// If already marked as aromatic, skip
			if a.bondsAromaticCount[edgeIdx] > 0 {
				continue
			}

			// Mark single bonds inside the cycle as aromatic
			if bond.Order == BOND_SINGLE {
				a.bondsAromaticCount[edgeIdx]++
			}
		}
	}
}

// --- helpers (graph cycles) ---

func findSimpleCyclesOfLength(m *Molecule, k int) [][]int {
	var cycles [][]int
	visited := make([]bool, len(m.Vertices))
	var path []int
	var dfs func(start, current, depth int)
	dfs = func(start, current, depth int) {
		if depth == k {
			// edge back to start?
			if areNeighbors(m, current, start) {
				cycle := append([]int(nil), path...)
				cycles = append(cycles, cycle)
			}
			return
		}
		visited[current] = true
		for _, eidx := range m.Vertices[current].Edges {
			e := m.Bonds[eidx]
			next := e.Beg
			if next == current {
				next = e.End
			} else if e.End != current {
				continue
			}
			if next < start { // avoid duplicates by ordering
				continue
			}
			if !visited[next] {
				path = append(path, next)
				dfs(start, next, depth+1)
				path = path[:len(path)-1]
			}
		}
		visited[current] = false
	}
	for i := range m.Vertices {
		path = []int{i}
		dfs(i, i, 1)
	}
	return dedupCycles(cycles)
}

func areNeighbors(m *Molecule, u, v int) bool {
	for _, eidx := range m.Vertices[u].Edges {
		e := m.Bonds[eidx]
		if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
			return true
		}
	}
	return false
}

func cycleEdges(m *Molecule, cycle []int) []int {
	var res []int
	for i := 0; i < len(cycle); i++ {
		u := cycle[i]
		v := cycle[(i+1)%len(cycle)]
		for _, eidx := range m.Vertices[u].Edges {
			e := m.Bonds[eidx]
			if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
				res = append(res, eidx)
				break
			}
		}
	}
	return res
}

// dedupCycles removes duplicate cycles irrespective of rotation or direction.
func dedupCycles(cycles [][]int) [][]int {
	seen := make(map[string]struct{})
	var res [][]int
	for _, c := range cycles {
		if len(c) == 0 {
			continue
		}
		key := normalizeCycleKey(c)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		res = append(res, c)
	}
	return res
}

func normalizeCycleKey(cycle []int) string {
	n := len(cycle)
	if n == 0 {
		return ""
	}
	// build two sequences: forward and reversed
	fwd := append([]int(nil), cycle...)
	rev := make([]int, n)
	for i := 0; i < n; i++ {
		rev[i] = cycle[n-1-i]
	}
	// rotate to put minimal vertex first
	fwd = rotateToMinFirst(fwd)
	rev = rotateToMinFirst(rev)
	// choose lexicographically smaller
	if lessSeq(rev, fwd) {
		return seqKey(rev)
	}
	return seqKey(fwd)
}

func rotateToMinFirst(seq []int) []int {
	n := len(seq)
	if n == 0 {
		return seq
	}
	minIdx := 0
	for i := 1; i < n; i++ {
		if seq[i] < seq[minIdx] {
			minIdx = i
		}
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = seq[(minIdx+i)%n]
	}
	return res
}

func lessSeq(a, b []int) bool {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return len(a) < len(b)
}

func seqKey(a []int) string {
	// fast integer join
	if len(a) == 0 {
		return ""
	}
	// convert to string with separators
	// avoid importing strings.Builder it is already used elsewhere
	key := ""
	for i, v := range a {
		if i > 0 {
			key += ","
		}
		// small numbers; fmt is ok
		key += fmtInt(v)
	}
	return key
}

func fmtInt(v int) string {
	// minimal int to string without strconv import
	if v == 0 {
		return "0"
	}
	neg := false
	if v < 0 {
		neg = true
		v = -v
	}
	buf := make([]byte, 0, 12)
	for v > 0 {
		d := byte(v % 10)
		buf = append(buf, '0'+d)
		v /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
