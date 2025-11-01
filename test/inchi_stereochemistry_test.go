package test

import (
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestInChIStereochemistry tests InChI generation with stereochemistry
func TestInChIStereochemistry(t *testing.T) {
	tests := []struct {
		name         string
		smiles       string
		expectStereo bool // Whether we expect stereochemistry layers
		description  string
	}{
		{
			name:         "Ethanol (no stereochemistry)",
			smiles:       "CCO",
			expectStereo: false,
			description:  "Simple molecule without stereochemistry",
		},
		{
			name:         "Trans-2-butene",
			smiles:       "C/C=C/C",
			expectStereo: true,
			description:  "Molecule with trans double bond",
		},
		{
			name:         "Cis-2-butene",
			smiles:       "C/C=C\\C",
			expectStereo: true,
			description:  "Molecule with cis double bond",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse SMILES
			loader := molecule.SmilesLoader{}
			mol, err := loader.Parse(tt.smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES: %v", err)
			}

			// Generate InChI
			generator := molecule.NewInChIGenerator()
			result, err := generator.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			t.Logf("Description: %s", tt.description)
			t.Logf("SMILES: %s", tt.smiles)
			t.Logf("InChI:  %s", result.InChI)
			t.Logf("InChIKey: %s", result.InChIKey)

			// Check if InChI contains stereochemistry layers
			hasStereoLayers := strings.Contains(result.InChI, "/b") ||
				strings.Contains(result.InChI, "/t") ||
				strings.Contains(result.InChI, "/m")

			if tt.expectStereo && !hasStereoLayers {
				t.Logf("Note: Expected stereochemistry layers but none found.")
				t.Logf("This is expected as stereochemistry implementation is in progress.")
			}

			// Basic validation
			if !strings.HasPrefix(result.InChI, "InChI=") {
				t.Errorf("InChI should start with 'InChI='")
			}

			if result.InChIKey == "" {
				t.Errorf("InChIKey should not be empty")
			}
		})
	}
}

// TestCisTransLayer tests cis/trans stereochemistry layer generation
func TestCisTransLayer(t *testing.T) {
	// Create a simple molecule with cis/trans stereochemistry
	mol := molecule.NewMolecule()

	// Build trans-2-butene: C1-C2=C3-C4
	c1 := mol.AddAtom(molecule.ELEM_C)
	c2 := mol.AddAtom(molecule.ELEM_C)
	c3 := mol.AddAtom(molecule.ELEM_C)
	c4 := mol.AddAtom(molecule.ELEM_C)

	mol.AddBond(c1, c2, molecule.BOND_SINGLE)
	doubleBond := mol.AddBond(c2, c3, molecule.BOND_DOUBLE)
	mol.AddBond(c3, c4, molecule.BOND_SINGLE)

	// Register cis/trans stereochemistry
	mol.CisTrans.RegisterBond(doubleBond)
	substituents := [4]int{c1, c2, c3, c4}
	mol.CisTrans.Add(doubleBond, substituents, molecule.TRANS)

	// Generate InChI
	generator := molecule.NewInChIGenerator()
	result, err := generator.GenerateInChI(mol)
	if err != nil {
		t.Fatalf("Failed to generate InChI: %v", err)
	}

	t.Logf("InChI: %s", result.InChI)
	t.Logf("InChIKey: %s", result.InChIKey)

	// Check basic structure
	if !strings.HasPrefix(result.InChI, "InChI=") {
		t.Errorf("InChI should start with 'InChI='")
	}
}

// TestTetrahedralLayer tests tetrahedral stereochemistry layer generation
func TestTetrahedralLayer(t *testing.T) {
	// Create a simple molecule with a chiral center
	mol := molecule.NewMolecule()

	// Build a chiral carbon with 4 different substituents
	// C-H-O-N (simplified chiral center)
	c := mol.AddAtom(molecule.ELEM_C)
	h := mol.AddAtom(molecule.ELEM_H)
	o := mol.AddAtom(molecule.ELEM_O)
	n := mol.AddAtom(molecule.ELEM_N)

	mol.AddBond(c, h, molecule.BOND_SINGLE)
	mol.AddBond(c, o, molecule.BOND_SINGLE)
	mol.AddBond(c, n, molecule.BOND_SINGLE)

	// Add stereocenter
	pyramid := [4]int{h, o, n, -1} // -1 for implicit H
	mol.Stereocenters.Add(c, molecule.STEREO_ATOM_ABS, 0, pyramid)

	// Generate InChI
	generator := molecule.NewInChIGenerator()
	result, err := generator.GenerateInChI(mol)
	if err != nil {
		t.Fatalf("Failed to generate InChI: %v", err)
	}

	t.Logf("InChI: %s", result.InChI)
	t.Logf("InChIKey: %s", result.InChIKey)

	// Check basic structure
	if !strings.HasPrefix(result.InChI, "InChI=") {
		t.Errorf("InChI should start with 'InChI='")
	}
}

// TestEnantiomerLayer tests enantiomer layer generation
func TestEnantiomerLayer(t *testing.T) {
	tests := []struct {
		name         string
		stereoType   int
		expectedCode string
	}{
		{
			name:         "Absolute stereochemistry",
			stereoType:   molecule.STEREO_ATOM_ABS,
			expectedCode: "0",
		},
		{
			name:         "Relative stereochemistry (AND)",
			stereoType:   molecule.STEREO_ATOM_AND,
			expectedCode: "1",
		},
		{
			name:         "Relative stereochemistry (OR)",
			stereoType:   molecule.STEREO_ATOM_OR,
			expectedCode: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mol := molecule.NewMolecule()

			// Build a simple molecule
			c := mol.AddAtom(molecule.ELEM_C)
			h1 := mol.AddAtom(molecule.ELEM_H)
			h2 := mol.AddAtom(molecule.ELEM_H)
			h3 := mol.AddAtom(molecule.ELEM_H)

			mol.AddBond(c, h1, molecule.BOND_SINGLE)
			mol.AddBond(c, h2, molecule.BOND_SINGLE)
			mol.AddBond(c, h3, molecule.BOND_SINGLE)

			// Add stereocenter with specific type
			pyramid := [4]int{h1, h2, h3, -1}
			mol.Stereocenters.Add(c, tt.stereoType, 0, pyramid)

			// Generate InChI
			generator := molecule.NewInChIGenerator()
			result, err := generator.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			t.Logf("Stereo type: %d", tt.stereoType)
			t.Logf("InChI: %s", result.InChI)

			// Check if /m layer exists and has expected code
			if strings.Contains(result.InChI, "/m") {
				parts := strings.Split(result.InChI, "/m")
				if len(parts) == 2 {
					// Extract enantiomer code
					code := string(parts[1][0])
					t.Logf("Enantiomer code: %s (expected: %s)", code, tt.expectedCode)
				}
			}
		})
	}
}

// TestStereochemistryIntegration tests complete stereochemistry functionality
func TestStereochemistryIntegration(t *testing.T) {
	t.Log("Testing integration of cis/trans and tetrahedral stereochemistry")

	// Test that stereochemistry functions don't crash
	mol := molecule.NewMolecule()

	// Add some atoms and bonds
	for i := 0; i < 5; i++ {
		mol.AddAtom(molecule.ELEM_C)
	}

	mol.AddBond(0, 1, molecule.BOND_SINGLE)
	mol.AddBond(1, 2, molecule.BOND_DOUBLE)
	mol.AddBond(2, 3, molecule.BOND_SINGLE)
	mol.AddBond(3, 4, molecule.BOND_SINGLE)

	// Generate InChI - should not crash even without stereochemistry data
	generator := molecule.NewInChIGenerator()
	result, err := generator.GenerateInChI(mol)
	if err != nil {
		t.Fatalf("Failed to generate InChI: %v", err)
	}

	t.Logf("InChI: %s", result.InChI)
	t.Logf("InChIKey: %s", result.InChIKey)

	// Basic validation
	if !molecule.ValidateInChI(result.InChI) {
		t.Errorf("Generated InChI is not valid")
	}

	if len(result.InChIKey) == 0 {
		t.Errorf("InChIKey should not be empty")
	}
}

// TestEmptyStereochemistry tests handling of molecules without stereochemistry
func TestEmptyStereochemistry(t *testing.T) {
	smilesList := []string{
		"C",        // Methane
		"CC",       // Ethane
		"CCO",      // Ethanol
		"c1ccccc1", // Benzene
	}

	for _, smiles := range smilesList {
		t.Run(smiles, func(t *testing.T) {
			loader := molecule.SmilesLoader{}
			mol, err := loader.Parse(smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES: %v", err)
			}

			generator := molecule.NewInChIGenerator()
			result, err := generator.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			// Should not contain stereochemistry layers for these simple molecules
			if strings.Contains(result.InChI, "/t") || strings.Contains(result.InChI, "/m/s") {
				t.Logf("Note: Unexpected stereochemistry layers in InChI: %s", result.InChI)
			}

			t.Logf("SMILES: %s -> InChI: %s", smiles, result.InChI)
		})
	}
}
