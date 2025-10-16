package test

import (
	"go-chem/src"
	"testing"
)

// buildBenzeneLike creates a 6-carbon ring with alternating 1-2-1-2-1-2 bonds
func buildBenzeneLike() *src.Molecule {
	m := src.NewMolecule()
	// add 6 carbons
	for i := 0; i < 6; i++ {
		m.AddAtom(src.ELEM_C)
	}
	// ring edges
	edges := [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 0}}
	for i, e := range edges {
		order := src.BOND_SINGLE
		if i%2 == 1 {
			order = src.BOND_DOUBLE
		}
		m.AddBond(e[0], e[1], order)
	}
	return m
}

func TestAromatizerBase_Aromatize_Benzene(t *testing.T) {
	m := buildBenzeneLike()
	a := src.AromatizerBase{}
	a.Aromatize(m)

	// expect all 6 ring bonds to be aromatic
	aromatic := 0
	for i := range m.Bonds {
		if m.GetBondOrder(i) == src.BOND_AROMATIC {
			aromatic++
		}
	}
	if aromatic != 6 {
		t.Fatalf("expected 6 aromatic bonds, got %d", aromatic)
	}
}

func TestDearomatizerBase_Apply_Benzene(t *testing.T) {
	m := buildBenzeneLike()
	// aromatize first
	src.AromatizerBase{}.Aromatize(m)

	// then dearomatize
	d := src.DearomatizerBase{}
	d.Apply(m)

	single := 0
	double := 0
	aromatic := 0
	for i := range m.Bonds {
		switch m.GetBondOrder(i) {
		case src.BOND_SINGLE:
			single++
		case src.BOND_DOUBLE:
			double++
		case src.BOND_AROMATIC:
			aromatic++
		}
	}
	if aromatic != 0 {
		t.Fatalf("expected 0 aromatic bonds after dearomatization, got %d", aromatic)
	}
	if single != 3 || double != 3 {
		t.Fatalf("expected 3 single and 3 double bonds, got single=%d double=%d", single, double)
	}
}

func TestElements_Lookup(t *testing.T) {
	n, err := src.ElementFromString("C")
	if err != nil || n != src.ELEM_C {
		t.Fatalf("ElementFromString(C) failed: n=%d err=%v", n, err)
	}
	if !src.ElementIsHalogen(src.ELEM_F) {
		t.Fatalf("ElementIsHalogen(F) expected true")
	}
}
