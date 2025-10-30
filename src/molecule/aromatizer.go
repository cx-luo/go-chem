// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:50
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : aromatizer.go
// @Software: GoLand
package molecule

// AromatizerBase provides an API to aromatize/dearomatize a molecule.
// Based on the Indigo implementation, using Huckel's rule (4n+2 π electrons).
type AromatizerBase struct{}

// Aromatize marks bonds in rings that satisfy Huckel's rule as aromatic.
// This checks for 4n+2 π electrons in planar cyclic systems.
func (a AromatizerBase) Aromatize(m *Molecule) {
	// Find all simple cycles (5-8 membered rings are most common aromatic rings)
	for size := 5; size <= 8; size++ {
		cycles := findSimpleCyclesOfLength(m, size)
		for _, cycle := range cycles {
			if a.isCycleAromatic(m, cycle) {
				a.aromatizeCycle(m, cycle)
			}
		}
	}

	// Update aromaticity array based on aromatic bonds
	if len(m.Aromaticity) != len(m.Atoms) {
		m.Aromaticity = make([]int, len(m.Atoms))
		for i := range m.Aromaticity {
			m.Aromaticity[i] = ATOM_ALIPHATIC
		}
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
// Based on C++ implementation in molecule_arom.cpp::_isCycleAromatic
func (a AromatizerBase) isCycleAromatic(m *Molecule, cycle []int) bool {
	// First check: all atoms must be potentially aromatic
	for _, atomIdx := range cycle {
		if !a.canBeAromatic(m, atomIdx) {
			return false
		}
	}

	// Calculate π electron count using pi labels
	piCount := 0
	for _, atomIdx := range cycle {
		piLabel := a.getPiLabel(m, atomIdx, cycle)
		if piLabel < 0 {
			return false // Atom cannot contribute to aromaticity
		}
		piCount += piLabel
	}

	// Huckel's rule: 4n+2 π electrons (where n = 0, 1, 2, ...)
	// This means: piCount - 2 must be divisible by 4
	if ((piCount - 2) % 4) != 0 {
		return false
	}

	// Additional check: must have alternating or all aromatic bonds
	return a.checkBondPattern(m, cycle)
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
// Based on MoleculeAromatizer::_getPiLabel in C++.
func (a AromatizerBase) getPiLabel(m *Molecule, atomIdx int, cycle []int) int {
	atom := m.Atoms[atomIdx]

	// Count bonds in and out of the ring
	doubleBondsInRing := 0
	doubleBondsOutRing := 0
	neighbors := 0

	for _, bondIdx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[bondIdx]
		otherAtom := bond.End
		if otherAtom == atomIdx {
			otherAtom = bond.Beg
		}

		neighbors++
		if bond.Order == BOND_DOUBLE {
			if isInCycle(otherAtom, cycle) {
				doubleBondsInRing++
			} else {
				doubleBondsOutRing++
			}
		} else if bond.Order == BOND_TRIPLE {
			return -1 // Triple bonds cannot be aromatic
		}
	}

	// If there's a double bond in the ring, atom contributes 1 π electron
	if doubleBondsInRing > 0 {
		return 1
	}

	// If there are multiple external double bonds, cannot be aromatic
	if doubleBondsOutRing > 1 {
		return -1
	}

	// Special case for external double bonds (e.g., C=O, C=S in aromatic rings)
	if doubleBondsOutRing == 1 {
		// Check if it's an acceptable exocyclic double bond (e.g., C=O, C=S, C=N)
		if atom.Number == ELEM_C || atom.Number == ELEM_S {
			return 0 // Vacant p-orbital
		}
		return -1
	}

	// Determine π contribution based on atom type and connectivity
	return a.getPiLabelByConnectivity(m, atomIdx, neighbors)
}

// getPiLabelByConnectivity determines π electron contribution based on atom type and connectivity.
// Based on MoleculeAromatizer::_getPiLabelByConn in C++.
func (a AromatizerBase) getPiLabelByConnectivity(m *Molecule, atomIdx int, connectivity int) int {
	atom := m.Atoms[atomIdx]

	// Simplified pi label assignment based on common aromatic atoms
	switch atom.Number {
	case ELEM_C:
		// Carbon with 3 connections: sp2 hybridized, contributes 1 π electron
		if connectivity == 3 {
			return 1
		}
	case ELEM_N:
		// Nitrogen with 2 connections: pyridine-like, contributes 1 π electron
		if connectivity == 2 {
			return 1
		}
		// Nitrogen with 3 connections: pyrrole-like, contributes 2 π electrons (lone pair)
		if connectivity == 3 {
			return 2
		}
	case ELEM_O:
		// Oxygen with 2 connections: furan-like, contributes 2 π electrons (lone pair)
		if connectivity == 2 {
			return 2
		}
	case ELEM_S:
		// Sulfur with 2 connections: thiophene-like, contributes 2 π electrons (lone pair)
		if connectivity == 2 {
			return 2
		}
	case ELEM_P:
		// Phosphorus similar to nitrogen
		if connectivity == 2 {
			return 1
		}
		if connectivity == 3 {
			return 2
		}
	}

	return -1 // Cannot determine or not aromatic
}

// checkBondPattern verifies that bonds in the cycle can be aromatic.
func (a AromatizerBase) checkBondPattern(m *Molecule, cycle []int) bool {
	// Check that all bonds are either single, double, or already aromatic
	for i := 0; i < len(cycle); i++ {
		u := cycle[i]
		v := cycle[(i+1)%len(cycle)]

		bondIdx := -1
		for _, eidx := range m.Vertices[u].Edges {
			e := m.Bonds[eidx]
			if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
				bondIdx = eidx
				break
			}
		}

		if bondIdx == -1 {
			return false
		}

		order := m.GetBondOrder(bondIdx)
		if order != BOND_SINGLE && order != BOND_DOUBLE && order != BOND_AROMATIC {
			return false
		}
	}
	return true
}

// aromatizeCycle marks all bonds in a cycle as aromatic.
// Based on C++ implementation in molecule_arom.cpp::_aromatizeCycle
func (a AromatizerBase) aromatizeCycle(m *Molecule, cycle []int) {
	// Mark all bonds in the cycle as aromatic
	for i := 0; i < len(cycle); i++ {
		u := cycle[i]
		v := cycle[(i+1)%len(cycle)]

		for _, eidx := range m.Vertices[u].Edges {
			e := m.Bonds[eidx]
			if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
				m.setBondOrderInternal(eidx, BOND_AROMATIC)
				break
			}
		}
	}
}

// isInCycle checks if an atom is in the given cycle.
func isInCycle(atomIdx int, cycle []int) bool {
	for _, a := range cycle {
		if a == atomIdx {
			return true
		}
	}
	return false
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
