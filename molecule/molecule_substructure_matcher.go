// Package molecule provides molecular structure manipulation and analysis tools.
// This file implements substructure matching (subgraph isomorphism).
package molecule

import (
	"fmt"
	"sort"
)

// SubstructureMatcher finds substructure matches in molecules
type SubstructureMatcher struct {
	query  *Molecule // Query (smaller) molecule
	target *Molecule // Target (larger) molecule
}

// NewSubstructureMatcher creates a new substructure matcher
func NewSubstructureMatcher(query, target *Molecule) *SubstructureMatcher {
	return &SubstructureMatcher{
		query:  query,
		target: target,
	}
}

// MatchResult represents a substructure match
type MatchResult struct {
	AtomMapping []int // Maps query atom indices to target atom indices
	BondMapping []int // Maps query bond indices to target bond indices
}

// FindAll finds all substructure matches
func (sm *SubstructureMatcher) FindAll() []*MatchResult {
	if sm.query.AtomCount() > sm.target.AtomCount() {
		return nil // Query can't be larger than target
	}
	if sm.query.AtomCount() == 0 {
		return []*MatchResult{{AtomMapping: []int{}, BondMapping: []int{}}}
	}

	var results []*MatchResult
	seen := make(map[string]bool)
	order := sm.matchOrder()
	atomMapping := make([]int, sm.query.AtomCount())
	for i := range atomMapping {
		atomMapping[i] = -1
	}
	used := make([]bool, sm.target.AtomCount())

	sm.enumerateMatches(0, order, atomMapping, used, seen, &results)

	return results
}

// FindFirst finds the first substructure match (faster than FindAll)
func (sm *SubstructureMatcher) FindFirst() *MatchResult {
	if sm.query.AtomCount() > sm.target.AtomCount() {
		return nil
	}
	if sm.query.AtomCount() == 0 {
		return &MatchResult{AtomMapping: []int{}, BondMapping: []int{}}
	}

	order := sm.matchOrder()
	atomMapping := make([]int, sm.query.AtomCount())
	for i := range atomMapping {
		atomMapping[i] = -1
	}
	used := make([]bool, sm.target.AtomCount())

	if sm.findFirstMatch(0, order, atomMapping, used) {
		result := &MatchResult{
			AtomMapping: append([]int(nil), atomMapping...),
			BondMapping: sm.buildBondMapping(atomMapping),
		}
		return result
	}

	return nil
}

// HasMatch checks if there is any substructure match
func (sm *SubstructureMatcher) HasMatch() bool {
	return sm.FindFirst() != nil
}

func (sm *SubstructureMatcher) matchOrder() []int {
	order := make([]int, sm.query.AtomCount())
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		degI := len(sm.query.Vertices[order[i]].Edges)
		degJ := len(sm.query.Vertices[order[j]].Edges)
		if degI != degJ {
			return degI > degJ
		}
		numI := sm.query.Atoms[order[i]].Number
		numJ := sm.query.Atoms[order[j]].Number
		if numI != numJ {
			return numI > numJ
		}
		return order[i] < order[j]
	})
	return order
}

func (sm *SubstructureMatcher) enumerateMatches(depth int, order []int, mapping []int, used []bool, seen map[string]bool, results *[]*MatchResult) {
	if depth == len(order) {
		bondMapping := sm.buildBondMapping(mapping)
		key := matchKey(mapping, bondMapping)
		if seen[key] {
			return
		}
		seen[key] = true
		*results = append(*results, &MatchResult{
			AtomMapping: append([]int(nil), mapping...),
			BondMapping: bondMapping,
		})
		return
	}

	queryIdx := order[depth]
	for targetIdx := 0; targetIdx < sm.target.AtomCount(); targetIdx++ {
		if used[targetIdx] || !sm.atomsMatch(queryIdx, targetIdx) {
			continue
		}
		if !sm.hasRequiredConnections(queryIdx, targetIdx, mapping) {
			continue
		}

		mapping[queryIdx] = targetIdx
		used[targetIdx] = true
		sm.enumerateMatches(depth+1, order, mapping, used, seen, results)
		used[targetIdx] = false
		mapping[queryIdx] = -1
	}
}

func (sm *SubstructureMatcher) findFirstMatch(depth int, order []int, mapping []int, used []bool) bool {
	if depth == len(order) {
		return true
	}

	queryIdx := order[depth]
	for targetIdx := 0; targetIdx < sm.target.AtomCount(); targetIdx++ {
		if used[targetIdx] || !sm.atomsMatch(queryIdx, targetIdx) {
			continue
		}
		if !sm.hasRequiredConnections(queryIdx, targetIdx, mapping) {
			continue
		}

		mapping[queryIdx] = targetIdx
		used[targetIdx] = true
		if sm.findFirstMatch(depth+1, order, mapping, used) {
			return true
		}
		used[targetIdx] = false
		mapping[queryIdx] = -1
	}
	return false
}

func matchKey(atomMapping, bondMapping []int) string {
	atoms := append([]int(nil), atomMapping...)
	bonds := append([]int(nil), bondMapping...)
	sort.Ints(atoms)
	sort.Ints(bonds)
	return fmt.Sprintf("a:%v|b:%v", atoms, bonds)
}

// atomsMatch checks if two atoms are compatible
func (sm *SubstructureMatcher) atomsMatch(queryIdx, targetIdx int) bool {
	qAtom := &sm.query.Atoms[queryIdx]
	tAtom := &sm.target.Atoms[targetIdx]

	// Atom number must match
	if qAtom.Number != tAtom.Number {
		return false
	}

	// Charge must match
	if qAtom.Charge != tAtom.Charge {
		return false
	}

	// Check if target has at least as many connections as query
	qNeighbors := len(sm.query.GetNeighbors(queryIdx))
	tNeighbors := len(sm.target.GetNeighbors(targetIdx))
	if tNeighbors < qNeighbors {
		return false
	}

	return true
}

// shouldHaveBond checks if there should be a bond given the current mapping
func (sm *SubstructureMatcher) shouldHaveBond(qIdx1, qIdx2, tIdx1, tIdx2 int, mapping []int) bool {
	// Check if there's a bond in query
	qBond := sm.query.FindBond(qIdx1, qIdx2)
	if qBond == -1 {
		return true // No bond required in query
	}

	// There must be a corresponding bond in target
	tBond := sm.target.FindBond(tIdx1, tIdx2)
	if tBond == -1 {
		return false
	}

	return bondsMatch(sm.query.GetBondOrder(qBond), sm.target.GetBondOrder(tBond))
}

func bondsMatch(queryOrder, targetOrder int) bool {
	switch queryOrder {
	case BOND_ANY:
		return targetOrder != BOND_ZERO && targetOrder != -1
	case BOND_SINGLE_OR_DOUBLE:
		return targetOrder == BOND_SINGLE || targetOrder == BOND_DOUBLE
	case BOND_SINGLE_OR_AROMATIC:
		return targetOrder == BOND_SINGLE || targetOrder == BOND_AROMATIC
	case BOND_DOUBLE_OR_AROMATIC:
		return targetOrder == BOND_DOUBLE || targetOrder == BOND_AROMATIC
	default:
		return queryOrder == targetOrder
	}
}

// hasRequiredConnections checks if mapped neighbors have correct connections
func (sm *SubstructureMatcher) hasRequiredConnections(queryIdx, targetIdx int, mapping []int) bool {
	// For each already-mapped neighbor of query atom
	queryNeighbors := sm.query.GetNeighbors(queryIdx)

	for _, qNeighbor := range queryNeighbors {
		if mapping[qNeighbor] != -1 {
			// This neighbor is mapped, check if bond exists in target
			tNeighbor := mapping[qNeighbor]

			qBond := sm.query.FindBond(queryIdx, qNeighbor)
			tBond := sm.target.FindBond(targetIdx, tNeighbor)

			if qBond != -1 {
				if tBond == -1 {
					return false // Query has bond but target doesn't
				}

				if !bondsMatch(sm.query.GetBondOrder(qBond), sm.target.GetBondOrder(tBond)) {
					return false
				}
			}
		}
	}

	return true
}

// buildBondMapping builds bond mapping from atom mapping
func (sm *SubstructureMatcher) buildBondMapping(atomMapping []int) []int {
	bondMapping := make([]int, sm.query.BondCount())

	for i, qBond := range sm.query.Bonds {
		qAtom1 := qBond.Beg
		qAtom2 := qBond.End

		tAtom1 := atomMapping[qAtom1]
		tAtom2 := atomMapping[qAtom2]

		if tAtom1 == -1 || tAtom2 == -1 {
			bondMapping[i] = -1
			continue
		}

		tBondIdx := sm.target.FindBond(tAtom1, tAtom2)
		bondMapping[i] = tBondIdx
	}

	return bondMapping
}

// CountMatches counts the number of substructure matches
func (sm *SubstructureMatcher) CountMatches() int {
	matches := sm.FindAll()
	return len(matches)
}

// GetMatchedAtoms returns the set of target atom indices that are matched
func (mr *MatchResult) GetMatchedAtoms() []int {
	matched := make([]int, 0, len(mr.AtomMapping))
	for _, targetIdx := range mr.AtomMapping {
		if targetIdx != -1 {
			matched = append(matched, targetIdx)
		}
	}
	return matched
}

// GetMatchedBonds returns the set of target bond indices that are matched
func (mr *MatchResult) GetMatchedBonds() []int {
	matched := make([]int, 0, len(mr.BondMapping))
	for _, targetIdx := range mr.BondMapping {
		if targetIdx != -1 {
			matched = append(matched, targetIdx)
		}
	}
	return matched
}

// IsComplete checks if the match is complete (all atoms mapped)
func (mr *MatchResult) IsComplete() bool {
	for _, idx := range mr.AtomMapping {
		if idx == -1 {
			return false
		}
	}
	return true
}

// String returns a string representation of the match
func (mr *MatchResult) String() string {
	return fmt.Sprintf("Match: %d atoms, %d bonds", len(mr.AtomMapping), len(mr.BondMapping))
}

// SMARTS-like matching (simplified)

// SMARTSMatcher provides SMARTS-like pattern matching
type SMARTSMatcher struct {
	pattern string
	query   *Molecule
}

// NewSMARTSMatcher creates a new SMARTS matcher (simplified)
func NewSMARTSMatcher(pattern string) (*SMARTSMatcher, error) {
	// This is a placeholder for future SMARTS support
	return &SMARTSMatcher{
		pattern: pattern,
	}, nil
}

// MaxCommonSubstructure finds maximum common substructure between two molecules
type MaxCommonSubstructure struct {
	mol1 *Molecule
	mol2 *Molecule
}

// NewMaxCommonSubstructure creates a new MCS finder
func NewMaxCommonSubstructure(mol1, mol2 *Molecule) *MaxCommonSubstructure {
	return &MaxCommonSubstructure{
		mol1: mol1,
		mol2: mol2,
	}
}

// Find finds a maximum common edge-preserving substructure.
func (mcs *MaxCommonSubstructure) Find() *Molecule {
	if mcs == nil || mcs.mol1 == nil || mcs.mol2 == nil {
		return NewMolecule()
	}
	if mcs.mol1.AtomCount() == 0 || mcs.mol2.AtomCount() == 0 {
		return NewMolecule()
	}

	small := mcs.mol1
	large := mcs.mol2
	if mcs.mol2.AtomCount() < mcs.mol1.AtomCount() {
		small = mcs.mol2
		large = mcs.mol1
	}

	order := make([]int, small.AtomCount())
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		degI := len(small.Vertices[order[i]].Edges)
		degJ := len(small.Vertices[order[j]].Edges)
		if degI != degJ {
			return degI > degJ
		}
		return order[i] < order[j]
	})

	current := make([]int, small.AtomCount())
	best := make([]int, small.AtomCount())
	for i := range current {
		current[i] = -1
		best[i] = -1
	}
	used := make([]bool, large.AtomCount())
	bestAtoms := 0
	bestBonds := 0

	var search func(depth, mappedAtoms int)
	search = func(depth, mappedAtoms int) {
		if mappedAtoms+(len(order)-depth) < bestAtoms {
			return
		}
		if depth == len(order) {
			bonds := countMappedBonds(small, large, current)
			if mappedAtoms > bestAtoms || (mappedAtoms == bestAtoms && bonds > bestBonds) {
				bestAtoms = mappedAtoms
				bestBonds = bonds
				copy(best, current)
			}
			return
		}

		atomIdx := order[depth]
		for targetIdx := 0; targetIdx < large.AtomCount(); targetIdx++ {
			if used[targetIdx] || !mcsAtomsCompatible(small, atomIdx, large, targetIdx) {
				continue
			}
			if !mcsConnectionsCompatible(small, large, atomIdx, targetIdx, current) {
				continue
			}

			current[atomIdx] = targetIdx
			used[targetIdx] = true
			search(depth+1, mappedAtoms+1)
			used[targetIdx] = false
			current[atomIdx] = -1
		}

		// MCS may require skipping atoms from the smaller input.
		search(depth+1, mappedAtoms)
	}

	search(0, 0)
	if bestAtoms == 0 {
		return NewMolecule()
	}

	return buildMappedSubstructure(small, large, best)
}

func mcsAtomsCompatible(mol1 *Molecule, atom1 int, mol2 *Molecule, atom2 int) bool {
	a1 := mol1.Atoms[atom1]
	a2 := mol2.Atoms[atom2]
	return a1.Number == a2.Number &&
		a1.Charge == a2.Charge &&
		a1.Isotope == a2.Isotope &&
		a1.Radical == a2.Radical
}

func mcsConnectionsCompatible(small, large *Molecule, atomIdx, targetIdx int, mapping []int) bool {
	for _, neighbor := range small.GetNeighbors(atomIdx) {
		mappedNeighbor := mapping[neighbor]
		if mappedNeighbor == -1 {
			continue
		}

		smallBond := small.FindBond(atomIdx, neighbor)
		largeBond := large.FindBond(targetIdx, mappedNeighbor)
		if smallBond == -1 {
			continue
		}
		if largeBond == -1 {
			return false
		}
		if !bondsMatch(small.GetBondOrder(smallBond), large.GetBondOrder(largeBond)) {
			return false
		}
	}
	return true
}

func countMappedBonds(small, large *Molecule, mapping []int) int {
	count := 0
	for _, bond := range small.Bonds {
		tBeg := mapping[bond.Beg]
		tEnd := mapping[bond.End]
		if tBeg == -1 || tEnd == -1 {
			continue
		}
		targetBond := large.FindBond(tBeg, tEnd)
		if targetBond != -1 && bondsMatch(bond.Order, large.GetBondOrder(targetBond)) {
			count++
		}
	}
	return count
}

func buildMappedSubstructure(small, large *Molecule, mapping []int) *Molecule {
	result := NewMolecule()
	atomMap := make(map[int]int)

	for smallIdx, largeIdx := range mapping {
		if largeIdx == -1 {
			continue
		}
		newIdx := result.AddAtom(small.Atoms[smallIdx].Number)
		result.Atoms[newIdx] = small.Atoms[smallIdx]
		atomMap[smallIdx] = newIdx
	}

	for _, bond := range small.Bonds {
		newBeg, okBeg := atomMap[bond.Beg]
		newEnd, okEnd := atomMap[bond.End]
		if !okBeg || !okEnd {
			continue
		}
		targetBond := large.FindBond(mapping[bond.Beg], mapping[bond.End])
		if targetBond == -1 || !bondsMatch(bond.Order, large.GetBondOrder(targetBond)) {
			continue
		}
		newBond := result.AddBond(newBeg, newEnd, bond.Order)
		result.Bonds[newBond].Direction = bond.Direction
	}

	return result
}

// Convenience functions

// IsSubstructureOf checks if query is a substructure of target
func IsSubstructureOf(query, target *Molecule) bool {
	matcher := NewSubstructureMatcher(query, target)
	return matcher.HasMatch()
}

// FindSubstructureMatches finds all matches of query in target
func FindSubstructureMatches(query, target *Molecule) []*MatchResult {
	matcher := NewSubstructureMatcher(query, target)
	return matcher.FindAll()
}

// CountSubstructureMatches counts matches of query in target
func CountSubstructureMatches(query, target *Molecule) int {
	matcher := NewSubstructureMatcher(query, target)
	return matcher.CountMatches()
}
