// Package molecule coding=utf-8
// @Project : go-chem
// @File    : tpsa.go
package molecule

// CalculateTPSA computes the Topological Polar Surface Area following the
// fragment-contribution table from Ertl, Rohde, Selzer (J. Med. Chem. 2000).
//
// When includeSP is true, sulfur and phosphorus atoms are included in the
// summation (this matches PubChem's TPSA which uses the "polar S/P" mode).
//
// Atoms that lie in an aromatic ring whose ring atoms carry an exocyclic
// double bond to oxygen (e.g. pyridinone, pyridazinone, uracil, the pyrimidine
// of purinones) are treated using the Kekulé fragment values, matching the
// behaviour of PubChem/CACTVS where the amide-like resonance preserves sp3 NH
// character instead of being scored as fully aromatic.
func (m *Molecule) CalculateTPSA(includeSP bool) float64 {
	if m == nil {
		return 0
	}

	ensureAromatized(m)
	amideRingAtoms := tpsaAmideLikeRingAtoms(m)
	pureAromaticAtoms := pureAromaticRingAtoms(m)

	var total float64
	for atomIdx := range m.Atoms {
		// Bridgehead atoms shared with a purely aromatic ring keep their
		// aromatic character; only atoms whose sole aromatic membership is
		// in an amide-like ring should fall back to Kekulé fragment values.
		amideOnly := amideRingAtoms[atomIdx] && !pureAromaticAtoms[atomIdx]
		number := m.GetAtomNumber(atomIdx)
		switch number {
		case ELEM_N, ELEM_O:
			total += m.tpsaAtomContribution(atomIdx, amideOnly)
		case ELEM_S, ELEM_P:
			if includeSP {
				total += m.tpsaAtomContribution(atomIdx, amideOnly)
			}
		}
	}
	return total
}

func (m *Molecule) tpsaAtomContribution(atomIdx int, amideOnly bool) float64 {
	env := m.tpsaAtomEnvironment(atomIdx)

	// PubChem/CACTVS scores nitro groups as if they were written in their
	// neutral [N](=O)=O form. Re-shape the environment of the nitro N and the
	// charge-separated [O-] so the lookup lands on the matching fragment.
	if env.number == ELEM_N && isCanonicalNitroNitrogen(m, atomIdx) {
		env.charge = 0
		env.doubleCount = 2
		env.singleCount = env.heavyDegree - 2
		if env.singleCount < 0 {
			env.singleCount = 0
		}
	}
	if env.number == ELEM_O && isNitroOxide(m, atomIdx) {
		env.charge = 0
		env.singleCount = 0
		env.doubleCount = 1
	}

	if amideOnly && env.aromaticCount > 0 {
		env = tpsaKekuleEnvironment(env)
	}
	switch env.number {
	case ELEM_N:
		return tpsaNitrogenContribution(env)
	case ELEM_O:
		return tpsaOxygenContribution(env)
	case ELEM_S:
		return tpsaSulfurContribution(env)
	case ELEM_P:
		return tpsaPhosphorusContribution(env)
	default:
		return 0
	}
}

func isCanonicalNitroNitrogen(m *Molecule, atomIdx int) bool {
	if m.GetAtomNumber(atomIdx) != ELEM_N || m.GetAtomCharge(atomIdx) != 1 {
		return false
	}
	doubleO, negO := 0, 0
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
				negO++
			}
		}
	}
	return doubleO >= 1 && negO >= 1
}

func isNitroOxide(m *Molecule, atomIdx int) bool {
	if m.GetAtomNumber(atomIdx) != ELEM_O || m.GetAtomCharge(atomIdx) >= 0 {
		return false
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		neighbor := otherEnd(m, eidx, atomIdx)
		if isCanonicalNitroNitrogen(m, neighbor) {
			return true
		}
	}
	return false
}

// tpsaKekuleEnvironment converts an aromatic environment to an equivalent
// Kekulé environment for fragment lookup. The atoms touched here belong to a
// ring that PubChem treats as non-aromatic (e.g. pyridinone, pyridazinone,
// chromenone). For nitrogen we reproduce sp3 amide character; for oxygen and
// sulfur we keep ether/thioether character (their Kekulé form has no double
// bond inside the ring).
func tpsaKekuleEnvironment(env tpsaEnvironment) tpsaEnvironment {
	out := env
	out.aromaticCount = 0
	if env.number == ELEM_N {
		if env.hydrogenCount > 0 || env.heavyDegree >= 3 {
			out.singleCount = env.heavyDegree
			out.doubleCount = 0
		} else {
			out.singleCount = env.heavyDegree - 1
			if out.singleCount < 0 {
				out.singleCount = 0
			}
			out.doubleCount = 1
		}
		return out
	}
	out.singleCount = env.heavyDegree
	out.doubleCount = 0
	return out
}

// tpsaAmideLikeRingAtoms returns the set of atoms that belong to at least one
// aromatic ring where any ring atom carries an exocyclic double bond to O.
// Atoms in such rings are scored using Kekulé fragment values to match the
// behaviour of PubChem's TPSA implementation.
func tpsaAmideLikeRingAtoms(m *Molecule) map[int]bool {
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
			if !cycleHasExocyclicOxygen(m, cycle, cycleSet) {
				continue
			}
			for _, atomIdx := range cycle {
				result[atomIdx] = true
			}
		}
	}
	return result
}

func cycleAllAromaticBonds(m *Molecule, cycle []int) bool {
	cycleLen := len(cycle)
	for i := 0; i < cycleLen; i++ {
		u := cycle[i]
		v := cycle[(i+1)%cycleLen]
		found := false
		for _, eidx := range m.Vertices[u].Edges {
			e := m.Bonds[eidx]
			if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
				found = true
				if e.Order != BOND_AROMATIC {
					return false
				}
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func cycleHasExocyclicOxygen(m *Molecule, cycle []int, cycleSet map[int]bool) bool {
	for _, atomIdx := range cycle {
		for _, eidx := range m.Vertices[atomIdx].Edges {
			bond := m.Bonds[eidx]
			if bond.Order != BOND_DOUBLE {
				continue
			}
			neighbor := bond.End
			if neighbor == atomIdx {
				neighbor = bond.Beg
			}
			if cycleSet[neighbor] {
				continue
			}
			switch m.GetAtomNumber(neighbor) {
			case ELEM_O, ELEM_S:
				return true
			}
		}
	}
	return false
}

func tpsaNitrogenContribution(env tpsaEnvironment) float64 {
	if env.aromaticCount > 0 {
		return tpsaAromaticNitrogen(env)
	}
	switch env.charge {
	case 0:
		return tpsaNeutralNitrogen(env)
	case 1:
		return tpsaCationicNitrogen(env)
	}
	return 0
}

func tpsaAromaticNitrogen(env tpsaEnvironment) float64 {
	if env.charge >= 1 {
		switch {
		case env.hydrogenCount > 0:
			return 14.14
		case env.aromaticCount >= 3:
			return 4.10
		case env.singleCount > 0:
			return 3.88
		default:
			return 4.10
		}
	}
	switch {
	case env.hydrogenCount > 0:
		return 15.79
	case env.aromaticCount >= 3:
		return 4.41
	case env.singleCount > 0:
		return 4.93
	case env.doubleCount > 0:
		return 8.39
	default:
		return 12.89
	}
}

func tpsaNeutralNitrogen(env tpsaEnvironment) float64 {
	switch {
	case env.tripleCount > 0:
		if env.doubleCount > 0 {
			return 13.60
		}
		return 23.79
	case env.doubleCount >= 2:
		return 11.68
	case env.doubleCount == 1:
		switch env.hydrogenCount {
		case 0:
			return 12.36
		default:
			return 23.85
		}
	}
	switch env.hydrogenCount {
	case 0:
		return 3.24
	case 1:
		return 12.03
	default:
		return 26.02
	}
}

func tpsaCationicNitrogen(env tpsaEnvironment) float64 {
	switch {
	case env.tripleCount > 0:
		return 4.36
	case env.doubleCount >= 1:
		switch env.hydrogenCount {
		case 0:
			return 3.01
		case 1:
			return 13.97
		default:
			return 25.59
		}
	}
	switch env.hydrogenCount {
	case 0:
		return 0
	case 1:
		return 4.44
	case 2:
		return 16.61
	default:
		return 27.64
	}
}

func tpsaOxygenContribution(env tpsaEnvironment) float64 {
	if env.aromaticCount > 0 {
		return 13.14
	}
	switch {
	case env.charge < 0:
		return 23.06
	case env.doubleCount > 0:
		return 17.07
	case env.hydrogenCount > 0:
		return 20.23
	default:
		return 9.23
	}
}

func tpsaSulfurContribution(env tpsaEnvironment) float64 {
	if env.aromaticCount > 0 {
		if env.doubleCount > 0 {
			return 21.70
		}
		return 28.24
	}
	switch env.doubleCount {
	case 0:
		if env.hydrogenCount > 0 {
			return 38.80
		}
		return 25.30
	case 1:
		if env.heavyDegree <= 1 {
			return 32.09
		}
		return 19.21
	default:
		return 8.38
	}
}

func tpsaPhosphorusContribution(env tpsaEnvironment) float64 {
	switch env.doubleCount {
	case 0:
		if env.hydrogenCount > 0 {
			return 23.47
		}
		return 13.59
	case 1:
		if env.heavyDegree <= 2 {
			return 34.14
		}
		return 9.81
	}
	return 9.81
}

type tpsaEnvironment struct {
	number        int
	charge        int
	hydrogenCount int
	heavyDegree   int
	doubleOCount  int
	singleCount   int
	doubleCount   int
	tripleCount   int
	aromaticCount int
}

func (m *Molecule) tpsaAtomEnvironment(atomIdx int) tpsaEnvironment {
	env := tpsaEnvironment{
		number:        m.GetAtomNumber(atomIdx),
		charge:        m.GetAtomCharge(atomIdx),
		hydrogenCount: m.atomHydrogenCount(atomIdx),
	}
	for _, eidx := range m.Vertices[atomIdx].Edges {
		order := m.GetBondOrder(eidx)
		neighbor := otherEnd(m, eidx, atomIdx)
		if isHeavyAtom(m, neighbor) {
			env.heavyDegree++
		}
		switch order {
		case BOND_SINGLE:
			env.singleCount++
		case BOND_DOUBLE:
			env.doubleCount++
			if m.GetAtomNumber(neighbor) == ELEM_O {
				env.doubleOCount++
			}
		case BOND_TRIPLE:
			env.tripleCount++
		case BOND_AROMATIC:
			env.aromaticCount++
		}
	}
	return env
}

func (m *Molecule) atomHydrogenCount(atomIdx int) int {
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

func (m *Molecule) atomIsAttachedToCarbonyl(atomIdx int) bool {
	for _, eidx := range m.Vertices[atomIdx].Edges {
		if isCarbonylCarbon(m, otherEnd(m, eidx, atomIdx)) {
			return true
		}
	}
	return false
}

func ensureAromatized(m *Molecule) {
	if m == nil || m.Aromatized {
		return
	}
	(&AromatizerBase{}).Aromatize(m)
}
