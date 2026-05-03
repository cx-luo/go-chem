// Package molecule coding=utf-8
// @Project : go-chem
// @File    : lipinski.go
package molecule

// NumRotatableBonds returns the count of rotatable bonds following a
// PubChem/Cactvs-compatible definition: single, non-ring bonds whose endpoints
// are non-terminal heavy atoms, excluding amide C-N bonds, hydrazide N-N
// bonds, bonds touching nitrogen atoms in nitro groups, perhalogenated carbons
// (CF3/CCl3/CBr3/CI3), and atoms participating in any triple bond.
func NumRotatableBonds(m *Molecule) int {
	if m == nil {
		return 0
	}
	ensureAromatized(m)

	count := 0
	for i := range m.Bonds {
		order := m.GetBondOrder(i)
		if order != BOND_SINGLE {
			continue
		}
		b := m.Bonds[i]
		if !isHeavyAtom(m, b.Beg) || !isHeavyAtom(m, b.End) {
			continue
		}
		if !isRotatableEndpoint(m, b.Beg) || !isRotatableEndpoint(m, b.End) {
			continue
		}
		if isAmideBond(m, i) {
			continue
		}
		if isHydrazideBond(m, i) {
			continue
		}
		if isRingBond(m, i) {
			continue
		}
		count++
	}
	return count
}

func isRotatableEndpoint(m *Molecule, atomIdx int) bool {
	if heavyAtomDegree(m, atomIdx) <= 1 {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		if m.GetBondOrder(eidx) == BOND_TRIPLE {
			return false
		}
	}
	if isPerhalogenatedCarbon(m, atomIdx) {
		return false
	}
	if isNitroNitrogen(m, atomIdx) {
		return false
	}
	return true
}

func isPerhalogenatedCarbon(m *Molecule, atomIdx int) bool {
	if m.GetAtomNumber(atomIdx) != ELEM_C {
		return false
	}
	halogens := map[int]int{}
	nonHalogenHeavy := 0
	for _, eidx := range m.Vertices[atomIdx].Edges {
		neighbor := otherEnd(m, eidx, atomIdx)
		switch m.GetAtomNumber(neighbor) {
		case ELEM_F, ELEM_Cl, ELEM_Br, ELEM_I:
			halogens[m.GetAtomNumber(neighbor)]++
		case ELEM_H:
		default:
			if isHeavyAtom(m, neighbor) {
				nonHalogenHeavy++
			}
		}
	}
	if nonHalogenHeavy != 1 {
		return false
	}
	for _, count := range halogens {
		if count >= 3 {
			return true
		}
	}
	return false
}

func isNitroNitrogen(m *Molecule, atomIdx int) bool {
	if m.GetAtomNumber(atomIdx) != ELEM_N {
		return false
	}
	doubleO := 0
	negativeSingleO := 0
	for _, eidx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[eidx]
		neighbor := otherEnd(m, eidx, atomIdx)
		if m.GetAtomNumber(neighbor) != ELEM_O {
			continue
		}
		switch bond.Order {
		case BOND_DOUBLE:
			doubleO++
		case BOND_SINGLE:
			if m.GetAtomCharge(neighbor) < 0 {
				negativeSingleO++
			}
		}
	}
	if doubleO >= 2 {
		return true
	}
	return doubleO >= 1 && negativeSingleO >= 1
}

func isHydrazideBond(m *Molecule, bondIdx int) bool {
	if bondIdx < 0 || bondIdx >= len(m.Bonds) {
		return false
	}
	bond := m.Bonds[bondIdx]
	if m.GetAtomNumber(bond.Beg) != ELEM_N || m.GetAtomNumber(bond.End) != ELEM_N {
		return false
	}
	return nitrogenBondedToCarbonyl(m, bond.Beg) && nitrogenBondedToCarbonyl(m, bond.End)
}

func nitrogenBondedToCarbonyl(m *Molecule, atomIdx int) bool {
	if m.GetAtomNumber(atomIdx) != ELEM_N {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		if isCarbonylCarbon(m, otherEnd(m, eidx, atomIdx)) {
			return true
		}
	}
	return false
}

func isHeavyAtom(m *Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Atoms) {
		return false
	}
	number := m.GetAtomNumber(atomIdx)
	return number > 0 && number != ELEM_H
}

func heavyAtomDegree(m *Molecule, atomIdx int) int {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) {
		return 0
	}
	degree := 0
	for _, eidx := range m.Vertices[atomIdx].Edges {
		neighbor := otherEnd(m, eidx, atomIdx)
		if isHeavyAtom(m, neighbor) {
			degree++
		}
	}
	return degree
}

func isRingBond(m *Molecule, bondIdx int) bool {
	if bondIdx < 0 || bondIdx >= len(m.Bonds) {
		return false
	}
	bond := m.Bonds[bondIdx]
	if bond.Beg < 0 || bond.Beg >= len(m.Vertices) || bond.End < 0 || bond.End >= len(m.Vertices) {
		return false
	}

	visited := make([]bool, len(m.Atoms))
	queue := []int{bond.Beg}
	visited[bond.Beg] = true

	for len(queue) > 0 {
		atomIdx := queue[0]
		queue = queue[1:]
		for _, eidx := range m.Vertices[atomIdx].Edges {
			if eidx == bondIdx {
				continue
			}
			neighbor := otherEnd(m, eidx, atomIdx)
			if neighbor == bond.End {
				return true
			}
			if neighbor < 0 || neighbor >= len(visited) || visited[neighbor] {
				continue
			}
			visited[neighbor] = true
			queue = append(queue, neighbor)
		}
	}
	return false
}

func isAmideBond(m *Molecule, bondIdx int) bool {
	if bondIdx < 0 || bondIdx >= len(m.Bonds) {
		return false
	}
	bond := m.Bonds[bondIdx]
	begNumber := m.GetAtomNumber(bond.Beg)
	endNumber := m.GetAtomNumber(bond.End)
	if begNumber == ELEM_N && isCarbonylCarbon(m, bond.End) {
		return true
	}
	return endNumber == ELEM_N && isCarbonylCarbon(m, bond.Beg)
}

func isCarbonylCarbon(m *Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) || m.GetAtomNumber(atomIdx) != ELEM_C {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[eidx]
		if bond.Order != BOND_DOUBLE {
			continue
		}
		neighbor := otherEnd(m, eidx, atomIdx)
		neighborNumber := m.GetAtomNumber(neighbor)
		if neighborNumber == ELEM_O || neighborNumber == ELEM_S {
			return true
		}
	}
	return false
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

// NumHydrogenBondAcceptors returns the PubChem/Cactvs-style hydrogen bond
// acceptor count. It includes N, O, S, and F atoms with non-positive charge,
// excluding nitrogen atoms involved in amide/sulfonamide/phosphoramide bonds
// and those that contribute their lone pair to an aromatic ring
// (pyrrole-like). Heavier halogens (Cl, Br, I) are not counted.
func NumHydrogenBondAcceptors(m *Molecule) int {
	if m == nil {
		return 0
	}
	ensureAromatized(m)
	amideRingAtoms := tpsaAmideLikeRingAtoms(m)
	pureAromaticAtoms := pureAromaticRingAtoms(m)

	c := 0
	for i := range m.Atoms {
		n := m.GetAtomNumber(i)
		if m.GetAtomCharge(i) > 0 {
			continue
		}
		switch n {
		case ELEM_O, ELEM_F:
			c++
		case ELEM_S:
			if !isSulfonylLikeSulfur(m, i) {
				c++
			}
		case ELEM_N:
			if isHydrogenBondAcceptorNitrogen(m, i, amideRingAtoms, pureAromaticAtoms) {
				c++
			}
		}
	}
	return c
}

func isHydrogenBondAcceptorNitrogen(m *Molecule, atomIdx int, amideRingAtoms, pureAromaticAtoms map[int]bool) bool {
	if isAmideLikeNitrogen(m, atomIdx) {
		return false
	}
	// Atoms whose only aromatic membership is a ring that PubChem treats as
	// non-aromatic (one of its atoms carries an exocyclic C=O / C=S, e.g.
	// acridone, pyridinone) keep their nitrogen lone pair available, so the
	// pyrrole-like exclusion that would otherwise apply is suppressed. If the
	// same atom also belongs to a purely aromatic ring (e.g. the triazole
	// bridgehead in a triazolo-quinazolinone), the lone pair remains delocalised
	// and the pyrrole exclusion still applies.
	if amideRingAtoms[atomIdx] && !pureAromaticAtoms[atomIdx] {
		return true
	}
	if isPyrroleLikeAromaticNitrogen(m, atomIdx) {
		return false
	}
	return true
}

// pureAromaticRingAtoms returns the atoms that belong to at least one aromatic
// ring whose ring atoms carry no exocyclic double bond to O / S. These atoms
// retain the standard aromatic interpretation regardless of any other
// (amide-like) ring they may also be part of.
func pureAromaticRingAtoms(m *Molecule) map[int]bool {
	result := make(map[int]bool)
	for size := 5; size <= 7; size++ {
		for _, cycle := range findSimpleCyclesOfLength(m, size) {
			if !cycleAllAromaticBonds(m, cycle) {
				continue
			}
			cycleSet := make(map[int]bool, len(cycle))
			for _, atomIdx := range cycle {
				cycleSet[atomIdx] = true
			}
			if cycleHasExocyclicOxygen(m, cycle, cycleSet) {
				continue
			}
			for _, atomIdx := range cycle {
				result[atomIdx] = true
			}
		}
	}
	return result
}

func isAmideLikeNitrogen(m *Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) || m.GetAtomNumber(atomIdx) != ELEM_N {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[eidx]
		neighbor := otherEnd(m, eidx, atomIdx)
		if isCarbonylCarbon(m, neighbor) || isPhosphorylPhosphorus(m, neighbor) {
			return true
		}
		// Single bond to an amidine / guanidine / isothiourea sp2 carbon: the
		// nitrogen's lone pair is delocalised into the C=N, so it does not
		// behave as a free hydrogen-bond acceptor. The C=N imine partner
		// itself is reached through a double bond and is therefore not
		// excluded by this branch.
		if bond.Order == BOND_SINGLE && isAmidineLikeCarbon(m, neighbor, atomIdx) {
			return true
		}
	}
	return false
}

// isAmidineLikeCarbon reports whether atomIdx is an sp2 carbon that has a
// double bond to a nitrogen that is not the excluded atom. This identifies
// the central carbon of amidines (R-N=C-NR'), guanidines (HN=C(NR)NR') and
// isothiouronium (S-C(=N)-N) groups.
func isAmidineLikeCarbon(m *Molecule, atomIdx, excludeAtomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) || m.GetAtomNumber(atomIdx) != ELEM_C {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[eidx]
		if bond.Order != BOND_DOUBLE {
			continue
		}
		neighbor := otherEnd(m, eidx, atomIdx)
		if neighbor == excludeAtomIdx {
			continue
		}
		if m.GetAtomNumber(neighbor) == ELEM_N {
			return true
		}
	}
	return false
}

func isPyrroleLikeAromaticNitrogen(m *Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) || m.GetAtomNumber(atomIdx) != ELEM_N {
		return false
	}
	aromatic := 0
	doubleInRing := 0
	for _, eidx := range m.Vertices[atomIdx].Edges {
		switch m.GetBondOrder(eidx) {
		case BOND_AROMATIC:
			aromatic++
		case BOND_DOUBLE:
			doubleInRing++
		}
	}
	if aromatic == 0 {
		return false
	}
	if doubleInRing > 0 {
		return false
	}
	if aromatic >= 3 {
		return true
	}
	if atomHydrogenCount(m, atomIdx) > 0 {
		return true
	}
	return heavyAtomDegree(m, atomIdx) >= 3
}

// isSulfonylLikeSulfur identifies sulfur atoms with two double-bonded oxygens
// (sulfone / sulfonamide / sulfonate / sulfonyl halide / sulfate). PubChem
// does not count those sulfurs as hydrogen-bond acceptors regardless of the
// remaining substituents because their lone pairs are heavily depleted.
func isSulfonylLikeSulfur(m *Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) || m.GetAtomNumber(atomIdx) != ELEM_S {
		return false
	}
	doubleO := 0
	for _, eidx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[eidx]
		if bond.Order != BOND_DOUBLE {
			continue
		}
		if m.GetAtomNumber(otherEnd(m, eidx, atomIdx)) == ELEM_O {
			doubleO++
		}
	}
	return doubleO >= 2
}

func isSulfonylSulfur(m *Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) || m.GetAtomNumber(atomIdx) != ELEM_S {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[eidx]
		if bond.Order == BOND_DOUBLE && m.GetAtomNumber(otherEnd(m, eidx, atomIdx)) == ELEM_O {
			return true
		}
	}
	return false
}

func isPhosphorylPhosphorus(m *Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) || m.GetAtomNumber(atomIdx) != ELEM_P {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[eidx]
		if bond.Order == BOND_DOUBLE && m.GetAtomNumber(otherEnd(m, eidx, atomIdx)) == ELEM_O {
			return true
		}
	}
	return false
}

func atomHydrogenCount(m *Molecule, atomIdx int) int {
	count := 0
	for _, eidx := range m.Vertices[atomIdx].Edges {
		if m.GetAtomNumber(otherEnd(m, eidx, atomIdx)) == ELEM_H {
			count++
		}
	}
	if !m.IsPseudoAtom(atomIdx) && !m.IsTemplateAtom(atomIdx) && m.Atoms[atomIdx].Number > 0 {
		count += m.GetImplicitH(atomIdx)
	}
	return count
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
		// Skip special atom types that don't support implicit H calculation
		if m.IsPseudoAtom(i) || m.IsTemplateAtom(i) || n == ELEM_RSITE {
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
