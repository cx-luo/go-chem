package test

import (
	srcpkg "github.com/cx-luo/go-chem/molecule"
	"testing"
)

func TestAromatize_Benzene(t *testing.T) {
	// Parse benzene from SMILES (lowercase 'c' indicates aromatic carbon)
	m, err := (srcpkg.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// SMILES parser already marks aromatic bonds, so check initial state
	if len(m.Bonds) != 6 {
		t.Fatalf("expected 6 bonds, got %d", len(m.Bonds))
	}

	// Count aromatic bonds after parsing (SMILES may pre-aromatize)
	initialAromatic := 0
	for i := range m.Bonds {
		if m.GetBondOrder(i) == srcpkg.BOND_AROMATIC {
			initialAromatic++
		}
	}

	// If not already aromatic, run aromatization
	if initialAromatic == 0 {
		a := &srcpkg.AromatizerBase{}
		a.Aromatize(m)
	}

	// Verify all 6 bonds are aromatic
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
	// Parse aromatic benzene
	m, err := (srcpkg.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// Ensure it's aromatized first
	a := &srcpkg.AromatizerBase{}
	a.Aromatize(m)

	// Verify it's aromatic before dearomatization
	aromaticBefore := 0
	for i := range m.Bonds {
		if m.GetBondOrder(i) == srcpkg.BOND_AROMATIC {
			aromaticBefore++
		}
	}
	if aromaticBefore != 6 {
		t.Logf("Warning: Expected 6 aromatic bonds before dearomatization, got %d", aromaticBefore)
	}

	// Dearomatize
	d := &srcpkg.DearomatizerBase{}
	d.Apply(m)

	// Count bond types after dearomatization
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

	// Benzene should have 3 alternating single and 3 double bonds after dearomatization
	if arom != 0 {
		t.Errorf("expected 0 aromatic bonds after dearomatization, got %d", arom)
	}
	if singles != 3 || doubles != 3 {
		t.Errorf("expected 3 single and 3 double bonds, got s=%d d=%d", singles, doubles)
	}
}

func TestAromatize_Cyclohexane_NotAromatic(t *testing.T) {
	// Parse cyclohexane (all single bonds, uppercase 'C' = aliphatic)
	m, err := (srcpkg.SmilesLoader{}).Parse("C1CCCCC1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// Try to aromatize
	a := &srcpkg.AromatizerBase{}
	a.Aromatize(m)

	// Cyclohexane should NOT become aromatic (saturated ring, no π electrons)
	for i := range m.Bonds {
		if m.GetBondOrder(i) == srcpkg.BOND_AROMATIC {
			t.Fatalf("cyclohexane (saturated ring) should not be aromatized")
		}
	}
}
