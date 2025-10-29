// Package src provides molecular structure manipulation and analysis tools.
// This file implements substructure matching (subgraph isomorphism).
package src

import (
	"fmt"
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

	var results []*MatchResult

	// Try each atom in target as potential match for first query atom
	for targetStart := 0; targetStart < sm.target.AtomCount(); targetStart++ {
		// Initialize mapping
		atomMapping := make([]int, sm.query.AtomCount())
		for i := range atomMapping {
			atomMapping[i] = -1
		}

		used := make([]bool, sm.target.AtomCount())

		// Try to build a complete mapping starting from this atom
		if sm.recursiveMatch(0, targetStart, atomMapping, used) {
			result := &MatchResult{
				AtomMapping: make([]int, len(atomMapping)),
				BondMapping: sm.buildBondMapping(atomMapping),
			}
			copy(result.AtomMapping, atomMapping)
			results = append(results, result)
		}
	}

	return results
}

// FindFirst finds the first substructure match (faster than FindAll)
func (sm *SubstructureMatcher) FindFirst() *MatchResult {
	if sm.query.AtomCount() > sm.target.AtomCount() {
		return nil
	}

	for targetStart := 0; targetStart < sm.target.AtomCount(); targetStart++ {
		atomMapping := make([]int, sm.query.AtomCount())
		for i := range atomMapping {
			atomMapping[i] = -1
		}

		used := make([]bool, sm.target.AtomCount())

		if sm.recursiveMatch(0, targetStart, atomMapping, used) {
			result := &MatchResult{
				AtomMapping: make([]int, len(atomMapping)),
				BondMapping: sm.buildBondMapping(atomMapping),
			}
			copy(result.AtomMapping, atomMapping)
			return result
		}
	}

	return nil
}

// HasMatch checks if there is any substructure match
func (sm *SubstructureMatcher) HasMatch() bool {
	return sm.FindFirst() != nil
}

// recursiveMatch performs recursive backtracking to find matches
func (sm *SubstructureMatcher) recursiveMatch(queryIdx, targetIdx int, mapping []int, used []bool) bool {
	// Check if this query atom can match this target atom
	if !sm.atomsMatch(queryIdx, targetIdx) {
		return false
	}

	// Mark this mapping
	mapping[queryIdx] = targetIdx
	used[targetIdx] = true

	// If we've mapped all query atoms, we have a complete match
	if queryIdx == sm.query.AtomCount()-1 {
		return true
	}

	// Try to match the next query atom
	nextQueryIdx := queryIdx + 1

	// Get neighbors of current query atom that need to be matched
	queryNeighbors := sm.query.GetNeighbors(queryIdx)

	// For each unmapped neighbor of current query atom
	for _, qNeighbor := range queryNeighbors {
		if mapping[qNeighbor] == -1 {
			// This neighbor needs to be matched
			// Try each unused target atom
			for targetNeighbor := 0; targetNeighbor < sm.target.AtomCount(); targetNeighbor++ {
				if !used[targetNeighbor] {
					// Check if there should be a bond between them
					if sm.shouldHaveBond(queryIdx, qNeighbor, targetIdx, targetNeighbor, mapping) {
						if sm.recursiveMatch(qNeighbor, targetNeighbor, mapping, used) {
							return true
						}
					}
				}
			}

			// If we couldn't match this neighbor, backtrack
			mapping[queryIdx] = -1
			used[targetIdx] = false
			return false
		}
	}

	// All neighbors are already mapped, continue with next unmapped atom
	for nextQueryIdx < sm.query.AtomCount() && mapping[nextQueryIdx] != -1 {
		nextQueryIdx++
	}

	if nextQueryIdx >= sm.query.AtomCount() {
		return true // All atoms matched
	}

	// Try to match next unmapped query atom with any unused target atom
	for targetNext := 0; targetNext < sm.target.AtomCount(); targetNext++ {
		if !used[targetNext] {
			// Check connectivity constraints
			if sm.hasRequiredConnections(nextQueryIdx, targetNext, mapping) {
				if sm.recursiveMatch(nextQueryIdx, targetNext, mapping, used) {
					return true
				}
			}
		}
	}

	// Backtrack
	mapping[queryIdx] = -1
	used[targetIdx] = false
	return false
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

	// Bond orders must match
	qOrder := sm.query.GetBondOrder(qBond)
	tOrder := sm.target.GetBondOrder(tBond)

	return qOrder == tOrder
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

			if qBond == -1 && tBond != -1 {
				return false // Query has no bond but target does
			}

			if qBond != -1 {
				if tBond == -1 {
					return false // Query has bond but target doesn't
				}

				// Check bond orders match
				if sm.query.GetBondOrder(qBond) != sm.target.GetBondOrder(tBond) {
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

// Find finds the maximum common substructure (simplified implementation)
func (mcs *MaxCommonSubstructure) Find() *Molecule {
	// This is a simplified placeholder implementation
	// A full MCS algorithm is quite complex (typically using backtracking or clique detection)

	maxSize := 0

	// Try substructures of mol1 in mol2
	for size := min(mcs.mol1.AtomCount(), mcs.mol2.AtomCount()); size > 0; size-- {
		// Generate substructures of mol1 with 'size' atoms
		// Check if any match in mol2
		// This is simplified - a real implementation would be more sophisticated

		if size > maxSize {
			maxSize = size
			break
		}
	}

	// Build MCS molecule from best mapping
	if maxSize == 0 {
		return NewMolecule()
	}

	// Return empty molecule for now (placeholder)
	// A full implementation would build the actual MCS from the mapping
	return NewMolecule()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
