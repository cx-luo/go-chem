package test

import (
	"github.com/cx-luo/go-chem/molecule"
	"testing"
)

// TestSubstructureMatchSimple tests simple substructure matching
func TestSubstructureMatchSimple(t *testing.T) {
	// Create query: C-C
	query := molecule.NewMolecule()
	c1 := query.AddAtom(molecule.ELEM_C)
	c2 := query.AddAtom(molecule.ELEM_C)
	query.AddBond(c1, c2, molecule.BOND_SINGLE)

	// Create target: C-C-C
	target := molecule.NewMolecule()
	t1 := target.AddAtom(molecule.ELEM_C)
	t2 := target.AddAtom(molecule.ELEM_C)
	t3 := target.AddAtom(molecule.ELEM_C)
	target.AddBond(t1, t2, molecule.BOND_SINGLE)
	target.AddBond(t2, t3, molecule.BOND_SINGLE)

	// Query should be a substructure of target
	matcher := molecule.NewSubstructureMatcher(query, target)
	if !matcher.HasMatch() {
		t.Error("C-C should be a substructure of C-C-C")
	}

	// Find matches
	matches := matcher.FindAll()
	if len(matches) < 1 {
		t.Errorf("should find at least 1 match, got %d", len(matches))
	}

	t.Logf("Found %d matches", len(matches))
}

// TestSubstructureMatchNoMatch tests when there's no match
func TestSubstructureMatchNoMatch(t *testing.T) {
	// Create query: C=C (double bond)
	query := molecule.NewMolecule()
	c1 := query.AddAtom(molecule.ELEM_C)
	c2 := query.AddAtom(molecule.ELEM_C)
	query.AddBond(c1, c2, molecule.BOND_DOUBLE)

	// Create target: C-C-C (all single bonds)
	target := molecule.NewMolecule()
	t1 := target.AddAtom(molecule.ELEM_C)
	t2 := target.AddAtom(molecule.ELEM_C)
	t3 := target.AddAtom(molecule.ELEM_C)
	target.AddBond(t1, t2, molecule.BOND_SINGLE)
	target.AddBond(t2, t3, molecule.BOND_SINGLE)

	// Should not match
	matcher := molecule.NewSubstructureMatcher(query, target)
	if matcher.HasMatch() {
		t.Error("C=C should not match in C-C-C")
	}
}

// TestSubstructureMatchBenzene tests matching in benzene
func TestSubstructureMatchBenzene(t *testing.T) {
	// Create query: C=C
	query := molecule.NewMolecule()
	c1 := query.AddAtom(molecule.ELEM_C)
	c2 := query.AddAtom(molecule.ELEM_C)
	query.AddBond(c1, c2, molecule.BOND_DOUBLE)

	// Create target: benzene
	target := buildBenzeneLike()

	matcher := molecule.NewSubstructureMatcher(query, target)
	matches := matcher.FindAll()

	// Benzene has 3 double bonds, so should find matches
	if len(matches) == 0 {
		t.Error("should find C=C in benzene")
	}

	t.Logf("Found %d C=C matches in benzene", len(matches))
}

// TestSubstructureMatchResult tests match result properties
func TestSubstructureMatchResult(t *testing.T) {
	query := molecule.NewMolecule()
	c1 := query.AddAtom(molecule.ELEM_C)
	c2 := query.AddAtom(molecule.ELEM_C)
	query.AddBond(c1, c2, molecule.BOND_SINGLE)

	target := molecule.NewMolecule()
	t1 := target.AddAtom(molecule.ELEM_C)
	t2 := target.AddAtom(molecule.ELEM_C)
	target.AddBond(t1, t2, molecule.BOND_SINGLE)

	matcher := molecule.NewSubstructureMatcher(query, target)
	match := matcher.FindFirst()

	if match == nil {
		t.Fatal("should find a match")
	}

	// Check match result
	if !match.IsComplete() {
		t.Error("match should be complete")
	}

	matchedAtoms := match.GetMatchedAtoms()
	if len(matchedAtoms) != 2 {
		t.Errorf("should match 2 atoms, got %d", len(matchedAtoms))
	}

	matchedBonds := match.GetMatchedBonds()
	if len(matchedBonds) != 1 {
		t.Errorf("should match 1 bond, got %d", len(matchedBonds))
	}
}

// TestSubstructureFindFirst tests finding first match
func TestSubstructureFindFirst(t *testing.T) {
	query := molecule.NewMolecule()
	c1 := query.AddAtom(molecule.ELEM_C)
	c2 := query.AddAtom(molecule.ELEM_C)
	query.AddBond(c1, c2, molecule.BOND_SINGLE)

	target := molecule.NewMolecule()
	for i := 0; i < 10; i++ {
		target.AddAtom(molecule.ELEM_C)
	}
	for i := 0; i < 9; i++ {
		target.AddBond(i, i+1, molecule.BOND_SINGLE)
	}

	matcher := molecule.NewSubstructureMatcher(query, target)
	match := matcher.FindFirst()

	if match == nil {
		t.Fatal("should find a match")
	}

	// FindFirst should be faster than FindAll but still work
	t.Log("FindFirst successfully found a match")
}

// TestSubstructureConvenienceFunctions tests convenience functions
func TestSubstructureConvenienceFunctions(t *testing.T) {
	query := molecule.NewMolecule()
	c1 := query.AddAtom(molecule.ELEM_C)
	c2 := query.AddAtom(molecule.ELEM_C)
	query.AddBond(c1, c2, molecule.BOND_SINGLE)

	target := molecule.NewMolecule()
	t1 := target.AddAtom(molecule.ELEM_C)
	t2 := target.AddAtom(molecule.ELEM_C)
	t3 := target.AddAtom(molecule.ELEM_C)
	target.AddBond(t1, t2, molecule.BOND_SINGLE)
	target.AddBond(t2, t3, molecule.BOND_SINGLE)

	// Test IsSubstructureOf
	if !molecule.IsSubstructureOf(query, target) {
		t.Error("query should be substructure of target")
	}

	// Test CountSubstructureMatches
	count := molecule.CountSubstructureMatches(query, target)
	if count == 0 {
		t.Error("should find at least 1 match")
	}

	// Test FindSubstructureMatches
	matches := molecule.FindSubstructureMatches(query, target)
	if len(matches) == 0 {
		t.Error("should find matches")
	}

	t.Logf("Found %d matches using convenience function", len(matches))
}

// TestSubstructureDifferentAtomTypes tests matching with different atom types
func TestSubstructureDifferentAtomTypes(t *testing.T) {
	// Query: C-O
	query := molecule.NewMolecule()
	c := query.AddAtom(molecule.ELEM_C)
	o := query.AddAtom(molecule.ELEM_O)
	query.AddBond(c, o, molecule.BOND_SINGLE)

	// Target: C-C-O-C
	target := molecule.NewMolecule()
	c1 := target.AddAtom(molecule.ELEM_C)
	c2 := target.AddAtom(molecule.ELEM_C)
	o2 := target.AddAtom(molecule.ELEM_O)
	c3 := target.AddAtom(molecule.ELEM_C)
	target.AddBond(c1, c2, molecule.BOND_SINGLE)
	target.AddBond(c2, o2, molecule.BOND_SINGLE)
	target.AddBond(o2, c3, molecule.BOND_SINGLE)

	// Should find matches
	matcher := molecule.NewSubstructureMatcher(query, target)
	matches := matcher.FindAll()

	if len(matches) == 0 {
		t.Error("should find C-O pattern in target")
	}

	t.Logf("Found %d C-O matches", len(matches))
}

// TestSubstructureSelfMatch tests matching a molecule with itself
func TestSubstructureSelfMatch(t *testing.T) {
	mol := buildBenzeneLike()

	matcher := molecule.NewSubstructureMatcher(mol, mol)
	match := matcher.FindFirst()

	if match == nil {
		t.Error("molecule should match itself")
	}

	if !match.IsComplete() {
		t.Error("self-match should be complete")
	}
}
