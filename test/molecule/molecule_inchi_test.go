package molecule_test

import (
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestToInChI tests converting molecule to InChI
func TestToInChI(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
	}{
		{"Ethanol", "CCO"},
		{"Benzene", "c1ccccc1"},
		{"Acetic acid", "CC(=O)O"},
		{"Water", "O"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := molecule.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load molecule: %v", err)
			}
			defer m.Close()

			inchi, err := m.ToInChI()
			if err != nil {
				t.Errorf("failed to convert to InChI: %v", err)
			}

			if inchi == "" {
				t.Error("expected non-empty InChI string")
			}

			// InChI should start with "InChI="
			if !strings.HasPrefix(inchi, "InChI=") {
				t.Errorf("InChI should start with 'InChI=', got: %s", inchi)
			}
		})
	}
}

// TestToInChIKey tests converting molecule to InChI Key
func TestToInChIKey(t *testing.T) {
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	inchiKey, err := m.ToInChIKey()
	if err != nil {
		t.Errorf("failed to convert to InChI Key: %v", err)
	}

	if inchiKey == "" {
		t.Error("expected non-empty InChI Key string")
	}

	// InChI Key should be 27 characters (14-10-1 format)
	if len(inchiKey) < 25 { // At least 25 chars
		t.Errorf("InChI Key seems too short: %s (len=%d)", inchiKey, len(inchiKey))
	}
}

// TestLoadInChI tests loading a molecule from InChI
func TestLoadInChI(t *testing.T) {
	// Load ethanol
	m1, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m1.Close()

	// Get InChI
	inchi, err := m1.ToInChI()
	if err != nil {
		t.Fatalf("failed to convert to InChI: %v", err)
	}

	// Load from InChI
	m2, err := molecule.LoadFromInChI(inchi)
	if err != nil {
		t.Fatalf("failed to load from InChI: %v", err)
	}
	defer m2.Close()

	// Verify atom counts match
	count1, _ := m1.CountAtoms()
	count2, _ := m2.CountAtoms()

	if count1 != count2 {
		t.Errorf("atom count mismatch: original=%d, from InChI=%d", count1, count2)
	}
}

// TestInChIRoundtrip tests InChI conversion roundtrip
func TestInChIRoundtrip(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
	}{
		{"Methanol", "CO"},
		{"Ethanol", "CCO"},
		{"Propanol", "CCCO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load molecule
			m1, err := molecule.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load molecule: %v", err)
			}
			defer m1.Close()

			// Convert to InChI
			inchi, err := m1.ToInChI()
			if err != nil {
				t.Fatalf("failed to convert to InChI: %v", err)
			}

			// Load from InChI
			m2, err := molecule.LoadFromInChI(inchi)
			if err != nil {
				t.Fatalf("failed to load from InChI: %v", err)
			}
			defer m2.Close()

			// Verify structure preserved
			count1, _ := m1.CountAtoms()
			count2, _ := m2.CountAtoms()

			if count1 != count2 {
				t.Errorf("atom count mismatch after roundtrip: %d vs %d", count1, count2)
			}
		})
	}
}

// TestInChIKeyConsistency tests that same molecule produces same InChI Key
func TestInChIKeyConsistency(t *testing.T) {
	// Load the same molecule twice
	m1, _ := molecule.LoadMoleculeFromString("CCO")
	defer m1.Close()

	m2, _ := molecule.LoadMoleculeFromString("OCC")
	defer m2.Close()

	// Get InChI Keys
	key1, err := m1.ToInChIKey()
	if err != nil {
		t.Fatalf("failed to get InChI Key 1: %v", err)
	}

	key2, err := m2.ToInChIKey()
	if err != nil {
		t.Fatalf("failed to get InChI Key 2: %v", err)
	}

	// Keys should be the same (same molecule, different SMILES)
	if key1 != key2 {
		t.Errorf("InChI Keys should match for same molecule: %s vs %s", key1, key2)
	}
}

// TestInChIHelperFunctions tests warning, log, and auxinfo functions
func TestInChIHelperFunctions(t *testing.T) {
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Generate InChI (this may produce warnings/logs)
	_, err = m.ToInChI()
	if err != nil {
		t.Fatalf("failed to generate InChI: %v", err)
	}

	// These functions should not panic
	warning := molecule.InChIWarning()
	log := molecule.InChILog()
	auxInfo := molecule.InChIAuxInfo()

	// Just verify they return strings (may be empty)
	_ = warning
	_ = log
	_ = auxInfo
}

// TestInChIOnClosedMolecule tests that InChI fails on closed molecule
func TestInChIOnClosedMolecule(t *testing.T) {
	m, _ := molecule.LoadMoleculeFromString("CCO")
	m.Close()

	_, err := m.ToInChI()
	if err == nil {
		t.Error("expected error when generating InChI on closed molecule")
	}

	_, err = m.ToInChIKey()
	if err == nil {
		t.Error("expected error when generating InChI Key on closed molecule")
	}
}

// TestLoadInvalidInChI tests loading from invalid InChI
func TestLoadInvalidInChI(t *testing.T) {
	_, err := molecule.LoadFromInChI("invalid inchi string")
	if err == nil {
		t.Error("expected error when loading from invalid InChI")
	}
}

// TestInChIForComplexMolecules tests InChI for complex structures
func TestInChIForComplexMolecules(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
	}{
		{"Aspirin", "CC(=O)Oc1ccccc1C(=O)O"},
		{"Caffeine", "CN1C=NC2=C1C(=O)N(C(=O)N2C)C"},
		{"Glucose", "C([C@@H]1[C@H]([C@@H]([C@H](C(O1)O)O)O)O)O"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := molecule.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load %s: %v", tt.name, err)
			}
			defer m.Close()

			inchi, err := m.ToInChI()
			if err != nil {
				t.Errorf("failed to generate InChI for %s: %v", tt.name, err)
			}

			if !strings.HasPrefix(inchi, "InChI=") {
				t.Errorf("invalid InChI format for %s: %s", tt.name, inchi)
			}

			// Also test InChI Key
			key, err := m.ToInChIKey()
			if err != nil {
				t.Errorf("failed to generate InChI Key for %s: %v", tt.name, err)
			}

			if key == "" {
				t.Errorf("empty InChI Key for %s", tt.name)
			}
		})
	}
}
