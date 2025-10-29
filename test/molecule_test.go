package test

import (
	"fmt"
	"go-chem/src/molecule"
	"testing"
)

// TestMoleculeBasics tests basic molecule operations
func TestMoleculeBasics(t *testing.T) {
	m := molecule.NewMolecule()

	// Test adding atoms
	c1 := m.AddAtom(molecule.ELEM_C)
	c2 := m.AddAtom(molecule.ELEM_C)
	c3 := m.AddAtom(molecule.ELEM_C)

	if m.AtomCount() != 3 {
		t.Fatalf("expected 3 atoms, got %d", m.AtomCount())
	}

	// Test adding bonds
	b1 := m.AddBond(c1, c2, molecule.BOND_SINGLE)
	b2 := m.AddBond(c2, c3, molecule.BOND_DOUBLE)

	if m.BondCount() != 2 {
		t.Fatalf("expected 2 bonds, got %d", m.BondCount())
	}

	// Test bond order
	if m.GetBondOrder(b1) != molecule.BOND_SINGLE {
		t.Errorf("bond 1 should be single")
	}
	if m.GetBondOrder(b2) != molecule.BOND_DOUBLE {
		t.Errorf("bond 2 should be double")
	}

	// Test neighbors
	neighbors := m.GetNeighbors(c2)
	if len(neighbors) != 2 {
		t.Fatalf("carbon 2 should have 2 neighbors, got %d", len(neighbors))
	}
}

// TestAtomProperties tests atom property management
func TestAtomProperties(t *testing.T) {
	m := molecule.NewMolecule()

	o := m.AddAtom(molecule.ELEM_O)

	// Test charge
	m.SetAtomCharge(o, -1)
	if m.GetAtomCharge(o) != -1 {
		t.Errorf("expected charge -1, got %d", m.GetAtomCharge(o))
	}

	// Test isotope
	m.SetAtomIsotope(o, 18)
	if m.GetAtomIsotope(o) != 18 {
		t.Errorf("expected isotope 18, got %d", m.GetAtomIsotope(o))
	}

	// Test radical
	m.SetAtomRadical(o, molecule.RADICAL_DOUBLET)
	if m.Atoms[o].Radical != molecule.RADICAL_DOUBLET {
		t.Errorf("expected radical doublet")
	}
}

// TestPseudoAtom tests pseudo atom handling
func TestPseudoAtom(t *testing.T) {
	m := molecule.NewMolecule()

	idx := m.AddAtom(molecule.ELEM_PSEUDO)
	m.SetPseudoAtom(idx, "Ph")

	if !m.IsPseudoAtom(idx) {
		t.Error("atom should be pseudo atom")
	}

	label, err := m.GetPseudoAtom(idx)
	if err != nil {
		t.Errorf("error getting pseudo atom: %v", err)
	}
	if label != "Ph" {
		t.Errorf("expected label 'Ph', got '%s'", label)
	}
}

// TestCoordinates tests 2D and 3D coordinate handling
func TestCoordinates(t *testing.T) {
	m := molecule.NewMolecule()

	c := m.AddAtom(molecule.ELEM_C)

	// Test 3D coordinates
	m.SetAtomXYZ(c, 1.0, 2.0, 3.0)
	pos := m.GetAtomXYZ(c)

	if pos.X != 1.0 || pos.Y != 2.0 || pos.Z != 3.0 {
		t.Errorf("3D coordinates mismatch: got (%f, %f, %f)", pos.X, pos.Y, pos.Z)
	}

	if !m.HaveXYZ {
		t.Error("molecule should have 3D coordinates")
	}

	// Test 2D coordinates
	m.SetAtomXY(c, 4.0, 5.0)
	pos2d := m.GetAtomXY(c)

	if pos2d.X != 4.0 || pos2d.Y != 5.0 {
		t.Errorf("2D coordinates mismatch: got (%f, %f)", pos2d.X, pos2d.Y)
	}
}

// TestMoleculeClone tests molecule cloning
func TestMoleculeClone(t *testing.T) {
	m1 := molecule.NewMolecule()
	m1.Name = "Original"

	c1 := m1.AddAtom(molecule.ELEM_C)
	c2 := m1.AddAtom(molecule.ELEM_C)
	m1.AddBond(c1, c2, molecule.BOND_DOUBLE)

	// Clone the molecule
	m2 := m1.Clone()

	if m2.AtomCount() != 2 {
		t.Errorf("clone should have 2 atoms, got %d", m2.AtomCount())
	}

	if m2.BondCount() != 1 {
		t.Errorf("clone should have 1 bond, got %d", m2.BondCount())
	}

	if m2.Name != "Original" {
		t.Errorf("clone name should be 'Original', got '%s'", m2.Name)
	}

	// Modify clone and ensure original is unchanged
	m2.AddAtom(molecule.ELEM_O)

	if m1.AtomCount() != 2 {
		t.Error("modifying clone should not affect original")
	}

	if m2.AtomCount() != 3 {
		t.Error("clone should have 3 atoms after addition")
	}
}

// TestFindBond tests bond finding functionality
func TestFindBond(t *testing.T) {
	m := molecule.NewMolecule()

	c1 := m.AddAtom(molecule.ELEM_C)
	c2 := m.AddAtom(molecule.ELEM_C)
	c3 := m.AddAtom(molecule.ELEM_C)

	b1 := m.AddBond(c1, c2, molecule.BOND_SINGLE)
	m.AddBond(c2, c3, molecule.BOND_DOUBLE)

	// Find existing bond
	found := m.FindBond(c1, c2)
	if found != b1 {
		t.Errorf("should find bond between c1 and c2, got %d", found)
	}

	// Find in reverse order
	found = m.FindBond(c2, c1)
	if found != b1 {
		t.Errorf("should find bond between c2 and c1, got %d", found)
	}

	// Try to find non-existent bond
	found = m.FindBond(c1, c3)
	if found != -1 {
		t.Error("should not find bond between c1 and c3")
	}
}

// TestImplicitHydrogens tests implicit hydrogen calculation
func TestImplicitHydrogens(t *testing.T) {
	m := molecule.NewMolecule()

	// Methyl carbon (CH3)
	c1 := m.AddAtom(molecule.ELEM_C)
	c2 := m.AddAtom(molecule.ELEM_C)
	m.AddBond(c1, c2, molecule.BOND_SINGLE)

	// c1 has 1 bond, should have 3 implicit H
	implH := m.GetImplicitH(c1)
	if implH != 3 {
		t.Errorf("methyl carbon should have 3 implicit H, got %d", implH)
	}
}

// TestGetOtherBondEnd tests getting the other end of a bond
func TestGetOtherBondEnd(t *testing.T) {
	m := molecule.NewMolecule()

	c1 := m.AddAtom(molecule.ELEM_C)
	c2 := m.AddAtom(molecule.ELEM_C)
	b := m.AddBond(c1, c2, molecule.BOND_TRIPLE)

	fmt.Println(m.CalcMolecularWeight())
	unit := molecule.CollectGross(m, molecule.GrossFormulaOptions{})
	fmt.Println(molecule.GrossUnitsToStringHill(unit, false))

	other := m.GetOtherBondEnd(b, c1)
	if other != c2 {
		t.Errorf("other end should be c2, got %d", other)
	}

	other = m.GetOtherBondEnd(b, c2)
	if other != c1 {
		t.Errorf("other end should be c1, got %d", other)
	}
}

// TestMolecularWeight tests molecular weight calculation
func TestMolecularWeight(t *testing.T) {
	m := molecule.NewMolecule()

	// Water: H2O
	o := m.AddAtom(molecule.ELEM_O)
	h1 := m.AddAtom(molecule.ELEM_H)
	h2 := m.AddAtom(molecule.ELEM_H)
	m.AddBond(o, h1, molecule.BOND_SINGLE)
	m.AddBond(o, h2, molecule.BOND_SINGLE)

	mw := m.CalcMolecularWeight()

	// Water molecular weight should be approximately 18 (O(16.00) + H(1.008)*2)
	// Allow some tolerance
	if mw < 17.0 || mw > 19.0 {
		t.Errorf("water molecular weight should be ~18, got %f", mw)
	}
}

// TestBondDirection tests stereochemical bond directions
func TestBondDirection(t *testing.T) {
	m := molecule.NewMolecule()

	c1 := m.AddAtom(molecule.ELEM_C)
	c2 := m.AddAtom(molecule.ELEM_C)
	b := m.AddBond(c1, c2, molecule.BOND_SINGLE)

	// Set up wedge
	m.SetBondDirection(b, molecule.BOND_UP)

	if m.GetBondDirection(b) != molecule.BOND_UP {
		t.Error("bond direction should be UP")
	}

	// Change to down
	m.SetBondDirection(b, molecule.BOND_DOWN)

	if m.GetBondDirection(b) != molecule.BOND_DOWN {
		t.Error("bond direction should be DOWN")
	}
}

// TestEditRevision tests edit revision tracking
func TestEditRevision(t *testing.T) {
	m := molecule.NewMolecule()

	rev1 := m.GetEditRevision()

	// Add atom should increment revision
	m.AddAtom(molecule.ELEM_C)
	rev2 := m.GetEditRevision()

	if rev2 <= rev1 {
		t.Error("edit revision should increase after adding atom")
	}

	// Clear should increment revision
	m.Clear()
	rev3 := m.GetEditRevision()

	if rev3 <= rev2 {
		t.Error("edit revision should increase after clear")
	}
}
