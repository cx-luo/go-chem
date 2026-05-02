package test

import (
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestGetAtomCoordinates tests the GetAtomCoordinates method
func TestGetAtomCoordinates(t *testing.T) {
	mol := molecule.NewMolecule()

	// Add an atom
	idx := mol.AddAtom(6) // Carbon

	// Initially should be [0, 0, 0]
	coords := mol.GetAtomCoordinates(idx)
	if coords[0] != 0 || coords[1] != 0 || coords[2] != 0 {
		t.Errorf("Initial coordinates should be [0, 0, 0], got [%f, %f, %f]", coords[0], coords[1], coords[2])
	}

	// Set coordinates
	mol.SetAtomCoordinates(idx, 1.5, 2.5, 3.5)

	// Get coordinates again
	coords = mol.GetAtomCoordinates(idx)
	if coords[0] != 1.5 || coords[1] != 2.5 || coords[2] != 3.5 {
		t.Errorf("Expected [1.5, 2.5, 3.5], got [%f, %f, %f]", coords[0], coords[1], coords[2])
	}

	// Check HaveXYZ flag
	if !mol.HaveXYZ {
		t.Error("HaveXYZ should be true after setting coordinates")
	}
}

// TestGetAtomCoordinatesOutOfRange tests boundary conditions
func TestGetAtomCoordinatesOutOfRange(t *testing.T) {
	mol := molecule.NewMolecule()
	mol.AddAtom(6) // Carbon

	// Test negative index
	coords := mol.GetAtomCoordinates(-1)
	if coords[0] != 0 || coords[1] != 0 || coords[2] != 0 {
		t.Error("Out of range index should return [0, 0, 0]")
	}

	// Test index beyond range
	coords = mol.GetAtomCoordinates(10)
	if coords[0] != 0 || coords[1] != 0 || coords[2] != 0 {
		t.Error("Out of range index should return [0, 0, 0]")
	}
}

// TestGet2DCoordinates tests the 2D coordinate methods
func TestGet2DCoordinates(t *testing.T) {
	mol := molecule.NewMolecule()

	// Add an atom
	idx := mol.AddAtom(6) // Carbon

	// Initially should be [0, 0]
	coords := mol.GetAtom2DCoordinates(idx)
	if coords[0] != 0 || coords[1] != 0 {
		t.Errorf("Initial 2D coordinates should be [0, 0], got [%f, %f]", coords[0], coords[1])
	}

	// Set 2D coordinates
	mol.SetAtom2DCoordinates(idx, 10.5, 20.5)

	// Get 2D coordinates again
	coords = mol.GetAtom2DCoordinates(idx)
	if coords[0] != 10.5 || coords[1] != 20.5 {
		t.Errorf("Expected [10.5, 20.5], got [%f, %f]", coords[0], coords[1])
	}
}

// TestCoordinatesWithSMILES tests coordinates with a real molecule
func TestCoordinatesWithSMILES(t *testing.T) {
	loader := molecule.SmilesLoader{}
	mol, err := loader.Parse("CCO") // Ethanol
	if err != nil {
		t.Fatalf("Failed to parse SMILES: %v", err)
	}

	// Initially, SMILES parsing doesn't set 3D coordinates
	coords := mol.GetAtomCoordinates(0)
	t.Logf("Initial coordinates of first atom: [%f, %f, %f]", coords[0], coords[1], coords[2])

	// Set coordinates for all atoms
	for i := 0; i < mol.AtomCount(); i++ {
		mol.SetAtomCoordinates(i, float64(i)*1.5, float64(i)*2.0, float64(i)*0.5)
	}

	// Verify coordinates were set
	for i := 0; i < mol.AtomCount(); i++ {
		coords := mol.GetAtomCoordinates(i)
		expectedX := float64(i) * 1.5
		expectedY := float64(i) * 2.0
		expectedZ := float64(i) * 0.5

		if coords[0] != expectedX || coords[1] != expectedY || coords[2] != expectedZ {
			t.Errorf("Atom %d: expected [%f, %f, %f], got [%f, %f, %f]",
				i, expectedX, expectedY, expectedZ, coords[0], coords[1], coords[2])
		}
	}

	// Check that HaveXYZ flag is set
	if !mol.HaveXYZ {
		t.Error("HaveXYZ should be true after setting coordinates")
	}
}
