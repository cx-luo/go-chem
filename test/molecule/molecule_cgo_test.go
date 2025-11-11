package molecule_test

import (
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestCreateMolecule tests creating a new molecule
func TestCreateMolecule(t *testing.T) {
	m, err := molecule.CreateMolecule()
	if err != nil {
		t.Fatalf("failed to create molecule: %v", err)
	}
	defer m.Close()

	if m.Handle < 0 {
		t.Errorf("invalid molecule handle: %d", m.Handle)
	}

	// Test counts
	atomCount, err := m.CountAtoms()
	if err != nil {
		t.Errorf("failed to count atoms: %v", err)
	}
	if atomCount != 0 {
		t.Errorf("expected 0 atoms, got %d", atomCount)
	}

	bondCount, err := m.CountBonds()
	if err != nil {
		t.Errorf("failed to count bonds: %v", err)
	}
	if bondCount != 0 {
		t.Errorf("expected 0 bonds, got %d", bondCount)
	}
}

// TestCreateQueryMolecule tests creating a query molecule
func TestCreateQueryMolecule(t *testing.T) {
	m, err := molecule.CreateQueryMolecule()
	if err != nil {
		t.Fatalf("failed to create query molecule: %v", err)
	}
	defer m.Close()

	if m.Handle < 0 {
		t.Errorf("invalid query molecule handle: %d", m.Handle)
	}
}

// TestMoleculeClose tests closing a molecule
func TestMoleculeClose(t *testing.T) {
	m, err := molecule.CreateMolecule()
	if err != nil {
		t.Fatalf("failed to create molecule: %v", err)
	}

	// Close the molecule
	err = m.Close()
	if err != nil {
		t.Errorf("failed to close molecule: %v", err)
	}

	// Closing again should not error
	err = m.Close()
	if err != nil {
		t.Errorf("second close should not error: %v", err)
	}

	// Operations on closed molecule should error
	_, err = m.CountAtoms()
	if err == nil {
		t.Error("expected error when counting atoms on closed molecule")
	}
}

// TestMoleculeClone tests cloning a molecule
func TestMoleculeClone(t *testing.T) {
	// Load a molecule
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Clone the molecule
	m2, err := m.Clone()
	if err != nil {
		t.Fatalf("failed to clone molecule: %v", err)
	}
	defer m2.Close()

	// Verify the clone has the same structure
	count1, _ := m.CountAtoms()
	count2, _ := m2.CountAtoms()
	if count1 != count2 {
		t.Errorf("clone has different atom count: original=%d, clone=%d", count1, count2)
	}
}

// TestMoleculeAromatize tests aromatization
func TestMoleculeAromatize(t *testing.T) {
	// Load benzene
	m, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Aromatize
	err = m.Aromatize()
	if err != nil {
		t.Errorf("failed to aromatize: %v", err)
	}
}

// TestMoleculeDearomatize tests dearomatization
func TestMoleculeDearomatize(t *testing.T) {
	// Load benzene
	m, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Dearomatize
	err = m.Dearomatize()
	if err != nil {
		t.Errorf("failed to dearomatize: %v", err)
	}
}

// TestMoleculeFoldUnfoldHydrogens tests hydrogen folding and unfolding
func TestMoleculeFoldUnfoldHydrogens(t *testing.T) {
	// Load ethanol
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Fold hydrogens
	err = m.FoldHydrogens()
	if err != nil {
		t.Errorf("failed to fold hydrogens: %v", err)
	}

	// Unfold hydrogens
	err = m.UnfoldHydrogens()
	if err != nil {
		t.Errorf("failed to unfold hydrogens: %v", err)
	}
}

// TestMoleculeLayout tests 2D layout
func TestMoleculeLayout(t *testing.T) {
	// Load a molecule
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Layout
	err = m.Layout()
	if err != nil {
		t.Errorf("failed to layout molecule: %v", err)
	}
}

// TestMoleculeClean2D tests 2D cleaning
func TestMoleculeClean2D(t *testing.T) {
	// Load a molecule
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Clean 2D
	err = m.Clean2D()
	if err != nil {
		t.Errorf("failed to clean 2D: %v", err)
	}
}

// TestMoleculeNormalize tests normalization
func TestMoleculeNormalize(t *testing.T) {
	// Load a molecule
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Normalize with default options
	err = m.Normalize("")
	if err != nil {
		t.Errorf("failed to normalize: %v", err)
	}
}

// TestMoleculeStandardize tests standardization
func TestMoleculeStandardize(t *testing.T) {
	// Load a molecule
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Standardize
	err = m.Standardize()
	if err != nil {
		t.Errorf("failed to standardize: %v", err)
	}
}

// TestMoleculeIonize tests ionization
func TestMoleculeIonize(t *testing.T) {
	// Load acetic acid
	m, err := molecule.LoadMoleculeFromString("CC(=O)O")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Ionize at pH 7.0
	err = m.Ionize(7.0, 0.5)
	if err != nil {
		t.Errorf("failed to ionize: %v", err)
	}
}

// TestMoleculeCountComponents tests counting connected components
func TestMoleculeCountComponents(t *testing.T) {
	// Load a molecule with two components
	m, err := molecule.LoadMoleculeFromString("CCO.C")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	count, err := m.CountComponents()
	if err != nil {
		t.Errorf("failed to count components: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 components, got %d", count)
	}
}

// TestMoleculeCountSSSR tests counting smallest set of smallest rings
func TestMoleculeCountSSSR(t *testing.T) {
	// Load benzene (1 ring)
	m, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	count, err := m.CountSSSR()
	if err != nil {
		t.Errorf("failed to count SSSR: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 ring, got %d", count)
	}
}

// TestMoleculeCountHeavyAtoms tests counting heavy atoms
func TestMoleculeCountHeavyAtoms(t *testing.T) {
	// Load ethanol (C2H6O - 3 heavy atoms)
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	count, err := m.CountHeavyAtoms()
	if err != nil {
		t.Errorf("failed to count heavy atoms: %v", err)
	}

	if count != 3 {
		t.Errorf("expected 3 heavy atoms, got %d", count)
	}
}
