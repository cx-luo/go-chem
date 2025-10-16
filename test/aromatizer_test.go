package test

import (
	srcpkg "go-chem/src"
	"testing"
)

func TestAromatize_Benzene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	// Before aromatization, ensure there are bonds
	if len(m.Bonds) != 6 {
		t.Fatalf("expected 6 bonds, got %d", len(m.Bonds))
	}

	(srcpkg.AromatizerBase{}).Aromatize(m)
	aromatic := 0
	for i := range m.Bonds {
		if m.GetBondOrder(i) == srcpkg.BOND_AROMATIC {
			aromatic++
		}
	}
	if aromatic != 6 {
		t.Fatalf("expected 6 aromatic bonds, got %d", aromatic)
	}
}

func TestDearomatize_Benzene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	(srcpkg.AromatizerBase{}).Aromatize(m)
	(srcpkg.DearomatizerBase{}).Apply(m)
	singles, doubles, arom := 0, 0, 0
	for i := range m.Bonds {
		switch m.GetBondOrder(i) {
		case srcpkg.BOND_SINGLE:
			singles++
		case srcpkg.BOND_DOUBLE:
			doubles++
		case srcpkg.BOND_AROMATIC:
			arom++
		}
	}
	if arom != 0 || singles != 3 || doubles != 3 {
		t.Fatalf("expected 3 single and 3 double bonds, got s=%d d=%d a=%d", singles, doubles, arom)
	}
}

func TestAromatize_Cyclohexane_NotAromatic(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("C1CCCCC1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	(srcpkg.AromatizerBase{}).Aromatize(m)
	for i := range m.Bonds {
		if m.GetBondOrder(i) == srcpkg.BOND_AROMATIC {
			t.Fatalf("cyclohexane should not be aromatized")
		}
	}
}
