package molecule_test

import (
	"os"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestLoadMoleculeFromString tests loading a molecule from SMILES
func TestLoadMoleculeFromString(t *testing.T) {
	tests := []struct {
		name    string
		smiles  string
		wantErr bool
	}{
		{"Ethanol", "CCO", false},
		{"Benzene", "c1ccccc1", false},
		{"Acetic acid", "CC(=O)O", false},
		{"Water", "O", false},
		{"Invalid", "invalid smiles", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := molecule.LoadMoleculeFromString(tt.smiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadMoleculeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				defer m.Close()
				if m.Handle() < 0 {
					t.Errorf("invalid molecule handle")
				}
			}
		})
	}
}

// TestLoadMoleculeFromFile tests loading a molecule from a file
func TestLoadMoleculeFromFile(t *testing.T) {
	// Create a temporary MOL file
	tmpFile, err := os.CreateTemp("", "test_molecule_*.mol")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write a simple MOL file (methane)
	molContent := `
  Mrv0541 02231512112D

  1  0  0  0  0  0            999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
M  END
`
	_, err = tmpFile.WriteString(molContent)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Load the molecule
	m, err := molecule.LoadMoleculeFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load molecule from file: %v", err)
	}
	defer m.Close()

	// Verify the molecule loaded
	atomCount, _ := m.CountAtoms()
	if atomCount != 1 {
		t.Errorf("expected 1 atom, got %d", atomCount)
	}
}

// TestLoadMoleculeFromBuffer tests loading a molecule from a buffer
func TestLoadMoleculeFromBuffer(t *testing.T) {
	smiles := "CCO"
	buffer := []byte(smiles)

	m, err := molecule.LoadMoleculeFromBuffer(buffer)
	if err != nil {
		t.Fatalf("failed to load molecule from buffer: %v", err)
	}
	defer m.Close()

	atomCount, _ := m.CountAtoms()
	if atomCount != 3 {
		t.Errorf("expected 3 atoms, got %d", atomCount)
	}
}

// TestLoadQueryMoleculeFromString tests loading a query molecule
func TestLoadQueryMoleculeFromString(t *testing.T) {
	// Query molecule with wildcards
	query := "[#6]CO"

	m, err := molecule.LoadQueryMoleculeFromString(query)
	if err != nil {
		t.Fatalf("failed to load query molecule: %v", err)
	}
	defer m.Close()

	if m.Handle() < 0 {
		t.Error("invalid query molecule handle")
	}
}

// TestLoadSmartsFromString tests loading a SMARTS pattern
func TestLoadSmartsFromString(t *testing.T) {
	tests := []struct {
		name    string
		smarts  string
		wantErr bool
	}{
		{"Alcohol", "[OH]", false},
		{"Carboxylic acid", "C(=O)O", false},
		{"Benzene ring", "c1ccccc1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := molecule.LoadSmartsFromString(tt.smarts)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSmartsFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				defer m.Close()
			}
		})
	}
}

// TestLoadStructureFromString tests loading a structure with parameters
func TestLoadStructureFromString(t *testing.T) {
	smiles := "CCO"

	m, err := molecule.LoadStructureFromString(smiles, "")
	if err != nil {
		t.Fatalf("failed to load structure: %v", err)
	}
	defer m.Close()

	if m.Handle() < 0 {
		t.Error("invalid molecule handle")
	}
}

// TestLoadEmptyBuffer tests loading from empty buffer
func TestLoadEmptyBuffer(t *testing.T) {
	_, err := molecule.LoadMoleculeFromBuffer([]byte{})
	if err == nil {
		t.Error("expected error when loading from empty buffer")
	}
}

// TestLoadMultipleComponents tests loading molecules with multiple components
func TestLoadMultipleComponents(t *testing.T) {
	// Ethanol and methane
	smiles := "CCO.C"

	m, err := molecule.LoadMoleculeFromString(smiles)
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Should have 2 components
	components, err := m.CountComponents()
	if err != nil {
		t.Errorf("failed to count components: %v", err)
	}

	if components != 2 {
		t.Errorf("expected 2 components, got %d", components)
	}
}

// TestLoadAromaticMolecule tests loading aromatic molecules
func TestLoadAromaticMolecule(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
	}{
		{"Benzene", "c1ccccc1"},
		{"Pyridine", "c1ccncc1"},
		{"Furan", "c1ccoc1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := molecule.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load %s: %v", tt.name, err)
			}
			defer m.Close()

			// Should have at least one ring
			rings, err := m.CountSSSR()
			if err != nil {
				t.Errorf("failed to count rings: %v", err)
			}

			if rings < 1 {
				t.Errorf("expected at least 1 ring, got %d", rings)
			}
		})
	}
}

// TestLoadMoleculeWithCharges tests loading molecules with charges
func TestLoadMoleculeWithCharges(t *testing.T) {
	// Sodium acetate
	smiles := "[Na+].CC(=O)[O-]"

	m, err := molecule.LoadMoleculeFromString(smiles)
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Should have 2 components (Na+ and acetate)
	components, err := m.CountComponents()
	if err != nil {
		t.Errorf("failed to count components: %v", err)
	}

	if components != 2 {
		t.Errorf("expected 2 components, got %d", components)
	}
}

// TestLoadMoleculeWithIsotopes tests loading molecules with isotopes
func TestLoadMoleculeWithIsotopes(t *testing.T) {
	// Deuterated methanol
	smiles := "[2H]CO"

	m, err := molecule.LoadMoleculeFromString(smiles)
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	if m.Handle() < 0 {
		t.Error("invalid molecule handle")
	}
}
