package molecule_test

import (
	"os"
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestToSmiles tests converting molecule to SMILES
func TestToSmiles(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
	}{
		{"Ethanol", "CCO"},
		{"Benzene", "c1ccccc1"},
		{"Acetic acid", "CC(=O)O"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := molecule.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load molecule: %v", err)
			}
			defer m.Close()

			smiles, err := m.ToSmiles()
			if err != nil {
				t.Errorf("failed to convert to SMILES: %v", err)
			}

			if smiles == "" {
				t.Error("expected non-empty SMILES string")
			}
		})
	}
}

// TestToCanonicalSmiles tests converting molecule to canonical SMILES
func TestToCanonicalSmiles(t *testing.T) {
	// Load the same molecule from different SMILES representations
	smiles1 := "CCO"
	smiles2 := "OCC"

	m1, _ := molecule.LoadMoleculeFromString(smiles1)
	defer m1.Close()

	m2, _ := molecule.LoadMoleculeFromString(smiles2)
	defer m2.Close()

	// Get canonical SMILES
	canon1, err := m1.ToCanonicalSmiles()
	if err != nil {
		t.Errorf("failed to get canonical SMILES: %v", err)
	}

	canon2, err := m2.ToCanonicalSmiles()
	if err != nil {
		t.Errorf("failed to get canonical SMILES: %v", err)
	}

	// Both should have the same canonical SMILES
	if canon1 != canon2 {
		t.Errorf("canonical SMILES differ: %s vs %s", canon1, canon2)
	}
}

// TestToMolfile tests converting molecule to MOL format
func TestToMolfile(t *testing.T) {
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	molfile, err := m.ToMolfile()
	if err != nil {
		t.Errorf("failed to convert to MOL: %v", err)
	}

	if molfile == "" {
		t.Error("expected non-empty MOL file string")
	}

	// Should contain basic MOL file structure
	if !strings.Contains(molfile, "V2000") && !strings.Contains(molfile, "V3000") {
		t.Error("MOL file should contain version marker")
	}

	if !strings.Contains(molfile, "M  END") {
		t.Error("MOL file should contain M  END marker")
	}
}

// TestSaveToFile tests saving molecule to a file
func TestSaveToFile(t *testing.T) {
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_molecule_*.mol")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Save to file
	err = m.SaveToFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to save to file: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("saved file does not exist: %v", err)
	}

	if info.Size() == 0 {
		t.Error("saved file is empty")
	}

	// Try to load the saved file
	m2, err := molecule.LoadMoleculeFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load saved molecule: %v", err)
	}
	defer m2.Close()

	// Verify atom counts match
	count1, _ := m.CountAtoms()
	count2, _ := m2.CountAtoms()
	if count1 != count2 {
		t.Errorf("atom count mismatch: original=%d, loaded=%d", count1, count2)
	}
}

// TestSaveLoadRoundtrip tests saving and loading a molecule
func TestSaveLoadRoundtrip(t *testing.T) {
	// Load a molecule
	m, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Get original properties
	origAtoms, _ := m.CountAtoms()
	origBonds, _ := m.CountBonds()

	// Save to file
	tmpFile, err := os.CreateTemp("", "test_molecule_*.mol")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = m.SaveToFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to save to file: %v", err)
	}

	// Load from file
	m2, err := molecule.LoadMoleculeFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load from file: %v", err)
	}
	defer m2.Close()

	// Verify properties
	loadedAtoms, _ := m2.CountAtoms()
	loadedBonds, _ := m2.CountBonds()

	if origAtoms != loadedAtoms {
		t.Errorf("atom count mismatch: original=%d, loaded=%d", origAtoms, loadedAtoms)
	}

	if origBonds != loadedBonds {
		t.Errorf("bond count mismatch: original=%d, loaded=%d", origBonds, loadedBonds)
	}
}

// TestToJSON tests converting molecule to JSON
func TestToJSON(t *testing.T) {
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	json, err := m.ToJSON()
	if err != nil {
		t.Errorf("failed to convert to JSON: %v", err)
	}

	if json == "" {
		t.Error("expected non-empty JSON string")
	}

	// Should contain JSON structure markers
	if !strings.Contains(json, "{") {
		t.Error("JSON should contain opening brace")
	}
}

// TestSaveToJSONFile tests saving molecule to JSON file
func TestSaveToJSONFile(t *testing.T) {
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_molecule_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Save to JSON file
	err = m.SaveToJSONFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to save to JSON file: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("saved JSON file does not exist: %v", err)
	}

	if info.Size() == 0 {
		t.Error("saved JSON file is empty")
	}
}

// TestToBase64String tests converting molecule to base64
func TestToBase64String(t *testing.T) {
	m, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	base64, err := m.ToBase64String()
	if err != nil {
		t.Errorf("failed to convert to base64: %v", err)
	}

	if base64 == "" {
		t.Error("expected non-empty base64 string")
	}
}

// TestSaveClosedMolecule tests that saving fails on closed molecule
func TestSaveClosedMolecule(t *testing.T) {
	m, _ := molecule.LoadMoleculeFromString("CCO")
	m.Close()

	// All save methods should return errors
	_, err := m.ToSmiles()
	if err == nil {
		t.Error("expected error when converting closed molecule to SMILES")
	}

	_, err = m.ToMolfile()
	if err == nil {
		t.Error("expected error when converting closed molecule to MOL")
	}

	err = m.SaveToFile("test.mol")
	if err == nil {
		t.Error("expected error when saving closed molecule")
	}
}

// TestSmilesRoundtrip tests SMILES conversion roundtrip
func TestSmilesRoundtrip(t *testing.T) {
	// Load a molecule
	m, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Convert to SMILES
	smiles, err := m.ToSmiles()
	if err != nil {
		t.Fatalf("failed to convert to SMILES: %v", err)
	}

	// Load from SMILES
	m2, err := molecule.LoadMoleculeFromString(smiles)
	if err != nil {
		t.Fatalf("failed to load from SMILES: %v", err)
	}
	defer m2.Close()

	// Verify counts
	count1, _ := m.CountAtoms()
	count2, _ := m2.CountAtoms()
	if count1 != count2 {
		t.Errorf("atom count mismatch after SMILES roundtrip: original=%d, reloaded=%d", count1, count2)
	}
}

// TestToSmarts tests converting molecule to SMARTS
func TestToSmarts(t *testing.T) {
	m, err := molecule.LoadSmartsFromString("[OH]")
	if err != nil {
		t.Fatalf("failed to load SMARTS: %v", err)
	}
	defer m.Close()

	smarts, err := m.ToSmarts()
	if err != nil {
		t.Errorf("failed to convert to SMARTS: %v", err)
	}

	if smarts == "" {
		t.Error("expected non-empty SMARTS string")
	}
}

// TestToCanonicalSmarts tests converting molecule to canonical SMARTS
func TestToCanonicalSmarts(t *testing.T) {
	m, err := molecule.LoadSmartsFromString("[OH]")
	if err != nil {
		t.Fatalf("failed to load SMARTS: %v", err)
	}
	defer m.Close()

	smarts, err := m.ToCanonicalSmarts()
	if err != nil {
		t.Errorf("failed to convert to canonical SMARTS: %v", err)
	}

	if smarts == "" {
		t.Error("expected non-empty canonical SMARTS string")
	}
}
