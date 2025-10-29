package test

import (
	"go-chem/src"
	"strings"
	"testing"
)

// TestMolfileLoadBasic tests basic MOL file loading
func TestMolfileLoadBasic(t *testing.T) {
	// Simple ethanol MOL file
	molString := `Ethanol
  Example

  3  2  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.5000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    2.0000    1.0000    0.0000 O   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0  0  0  0
  2  3  1  0  0  0  0
M  END
`

	mol, err := src.LoadMoleculeFromString(molString)
	if err != nil {
		t.Fatalf("error loading MOL: %v", err)
	}

	if mol.AtomCount() != 3 {
		t.Errorf("expected 3 atoms, got %d", mol.AtomCount())
	}

	if mol.BondCount() != 2 {
		t.Errorf("expected 2 bonds, got %d", mol.BondCount())
	}

	// Check atom types
	if mol.GetAtomNumber(0) != src.ELEM_C {
		t.Error("atom 0 should be carbon")
	}
	if mol.GetAtomNumber(1) != src.ELEM_C {
		t.Error("atom 1 should be carbon")
	}
	if mol.GetAtomNumber(2) != src.ELEM_O {
		t.Error("atom 2 should be oxygen")
	}

	// Check bond orders
	if mol.GetBondOrder(0) != src.BOND_SINGLE {
		t.Error("bond 0 should be single")
	}
	if mol.GetBondOrder(1) != src.BOND_SINGLE {
		t.Error("bond 1 should be single")
	}
}

// TestMolfileLoadWithCharge tests loading molecules with charges
func TestMolfileLoadWithCharge(t *testing.T) {
	molString := `Charged
  Example

  2  1  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.5000    0.0000    0.0000 O   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0  0  0  0
M  CHG  1   2  -1
M  END
`

	mol, err := src.LoadMoleculeFromString(molString)
	if err != nil {
		t.Fatalf("error loading MOL: %v", err)
	}

	// Check charge on oxygen
	if mol.GetAtomCharge(1) != -1 {
		t.Errorf("oxygen should have charge -1, got %d", mol.GetAtomCharge(1))
	}
}

// TestMolfileLoadWithIsotope tests loading molecules with isotopes
func TestMolfileLoadWithIsotope(t *testing.T) {
	molString := `Deuterium
  Example

  1  0  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 H   0  0  0  0  0  0  0  0  0  0  0  0
M  ISO  1   1   2
M  END
`

	mol, err := src.LoadMoleculeFromString(molString)
	if err != nil {
		t.Fatalf("error loading MOL: %v", err)
	}

	if mol.GetAtomIsotope(0) != 2 {
		t.Errorf("hydrogen should be deuterium (isotope 2), got %d", mol.GetAtomIsotope(0))
	}
}

// TestMolfileSaveBasic tests basic MOL file saving
func TestMolfileSaveBasic(t *testing.T) {
	// Create a simple molecule (methane: CH4)
	m := src.NewMolecule()
	m.Name = "Methane"

	c := m.AddAtom(src.ELEM_C)
	h1 := m.AddAtom(src.ELEM_H)
	h2 := m.AddAtom(src.ELEM_H)
	h3 := m.AddAtom(src.ELEM_H)
	h4 := m.AddAtom(src.ELEM_H)

	m.AddBond(c, h1, src.BOND_SINGLE)
	m.AddBond(c, h2, src.BOND_SINGLE)
	m.AddBond(c, h3, src.BOND_SINGLE)
	m.AddBond(c, h4, src.BOND_SINGLE)

	// Set some coordinates
	m.SetAtomXYZ(c, 0, 0, 0)
	m.SetAtomXYZ(h1, 1, 0, 0)
	m.SetAtomXYZ(h2, 0, 1, 0)
	m.SetAtomXYZ(h3, 0, 0, 1)
	m.SetAtomXYZ(h4, -1, -1, -1)

	// Save to string
	molString, err := src.SaveMoleculeToString(m)
	if err != nil {
		t.Fatalf("error saving MOL: %v", err)
	}

	// Basic checks on output
	if !strings.Contains(molString, "Methane") {
		t.Error("MOL output should contain molecule name")
	}

	if !strings.Contains(molString, "5  4") {
		t.Error("MOL output should contain counts line with 5 atoms and 4 bonds")
	}

	if !strings.Contains(molString, "M  END") {
		t.Error("MOL output should contain M  END")
	}
}

// TestMolfileRoundTrip tests loading and saving
func TestMolfileRoundTrip(t *testing.T) {
	// Original MOL
	originalMol := `Propane
  Example

  3  2  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.5000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    3.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0  0  0  0
  2  3  1  0  0  0  0
M  END
`

	// Load
	mol, err := src.LoadMoleculeFromString(originalMol)
	if err != nil {
		t.Fatalf("error loading: %v", err)
	}

	// Save
	savedMol, err := src.SaveMoleculeToString(mol)
	if err != nil {
		t.Fatalf("error saving: %v", err)
	}

	// Load again
	mol2, err := src.LoadMoleculeFromString(savedMol)
	if err != nil {
		t.Fatalf("error reloading: %v", err)
	}

	// Compare
	if mol.AtomCount() != mol2.AtomCount() {
		t.Errorf("atom counts differ: %d vs %d", mol.AtomCount(), mol2.AtomCount())
	}

	if mol.BondCount() != mol2.BondCount() {
		t.Errorf("bond counts differ: %d vs %d", mol.BondCount(), mol2.BondCount())
	}

	// Check atom types preserved
	for i := 0; i < mol.AtomCount(); i++ {
		if mol.GetAtomNumber(i) != mol2.GetAtomNumber(i) {
			t.Errorf("atom %d type differs", i)
		}
	}

	// Check bond orders preserved
	for i := 0; i < mol.BondCount(); i++ {
		if mol.GetBondOrder(i) != mol2.GetBondOrder(i) {
			t.Errorf("bond %d order differs", i)
		}
	}
}

// TestMolfileBondTypes tests different bond types
func TestMolfileBondTypes(t *testing.T) {
	molString := `BondTypes
  Example

  4  3  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.5000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    3.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    4.5000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0  0  0  0
  2  3  2  0  0  0  0
  3  4  3  0  0  0  0
M  END
`

	mol, err := src.LoadMoleculeFromString(molString)
	if err != nil {
		t.Fatalf("error loading: %v", err)
	}

	// Check bond types: single, double, triple
	if mol.GetBondOrder(0) != src.BOND_SINGLE {
		t.Error("bond 0 should be single")
	}
	if mol.GetBondOrder(1) != src.BOND_DOUBLE {
		t.Error("bond 1 should be double")
	}
	if mol.GetBondOrder(2) != src.BOND_TRIPLE {
		t.Error("bond 2 should be triple")
	}
}

// TestMolfilePseudoAtom tests pseudo atom handling
func TestMolfilePseudoAtom(t *testing.T) {
	molString := `PseudoAtom
  Example

  2  1  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.5000    0.0000    0.0000 Ph  0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0  0  0  0
M  END
`

	mol, err := src.LoadMoleculeFromString(molString)
	if err != nil {
		t.Fatalf("error loading: %v", err)
	}

	// Check that second atom is pseudo
	if !mol.IsPseudoAtom(1) {
		t.Error("atom 1 should be pseudo atom")
	}

	label, err := mol.GetPseudoAtom(1)
	if err != nil {
		t.Errorf("error getting pseudo atom: %v", err)
	}
	if label != "Ph" {
		t.Errorf("pseudo atom label should be 'Ph', got '%s'", label)
	}
}

// TestMolfileStereo tests stereochemical bonds
func TestMolfileStereo(t *testing.T) {
	molString := `Stereo
  Example

  4  3  0  0  1  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.5000    1.0000    0.0000 H   0  0  0  0  0  0  0  0  0  0  0  0
    1.5000   -1.0000    0.0000 H   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0  0  0  0
  2  3  1  1  0  0  0
  2  4  1  6  0  0  0
M  END
`

	mol, err := src.LoadMoleculeFromString(molString)
	if err != nil {
		t.Fatalf("error loading: %v", err)
	}

	// Check chiral flag
	if mol.ChiralFlag != 1 {
		t.Error("molecule should have chiral flag set")
	}

	// Check bond directions
	if mol.GetBondDirection(1) != src.BOND_UP {
		t.Error("bond 1 should be wedge (up)")
	}
	if mol.GetBondDirection(2) != src.BOND_DOWN {
		t.Error("bond 2 should be hash (down)")
	}
}
