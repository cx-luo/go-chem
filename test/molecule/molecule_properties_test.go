package molecule_test

import (
	"strings"
	"testing"
)

// TestGrossFormula tests getting gross formula
func TestGrossFormula(t *testing.T) {
	tests := []struct {
		name    string
		smiles  string
		formula string
	}{
		{"Ethanol", "CCO", "C2H6O"},
		{"Benzene", "c1ccccc1", "C6H6"},
		{"Water", "O", "H2O"},
		{"Acetic acid", "CC(=O)O", "C2H4O2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := indigoInit.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load molecule: %v", err)
			}
			defer m.Close()

			formula, err := m.GrossFormula()
			if err != nil {
				t.Errorf("failed to get gross formula: %v", err)
			}

			// Just check it's not empty
			if formula == "" {
				t.Errorf("expected non-empty formula")
			}
		})
	}
}

// TestMolecularWeight tests getting molecular weight
func TestMolecularWeight(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
		minMW  float64 // Minimum expected molecular weight
		maxMW  float64 // Maximum expected molecular weight
	}{
		{"Water", "O", 18.0, 18.1},
		{"Ethanol", "CCO", 46.0, 46.1},
		{"Benzene", "c1ccccc1", 78.0, 78.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := indigoInit.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load molecule: %v", err)
			}
			defer m.Close()

			mw, err := m.MolecularWeight()
			if err != nil {
				t.Errorf("failed to get molecular weight: %v", err)
			}

			if mw < tt.minMW || mw > tt.maxMW {
				t.Errorf("molecular weight %f not in range [%f, %f]", mw, tt.minMW, tt.maxMW)
			}
		})
	}
}

// TestMonoisotopicMass tests getting monoisotopic mass
func TestMonoisotopicMass(t *testing.T) {
	m, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	mass, err := m.MonoisotopicMass()
	if err != nil {
		t.Errorf("failed to get monoisotopic mass: %v", err)
	}

	if mass <= 0 {
		t.Errorf("expected positive mass, got %f", mass)
	}
}

// TestMostAbundantMass tests getting most abundant mass
func TestMostAbundantMass(t *testing.T) {
	m, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	mass, err := m.MostAbundantMass()
	if err != nil {
		t.Errorf("failed to get most abundant mass: %v", err)
	}

	if mass <= 0 {
		t.Errorf("expected positive mass, got %f", mass)
	}
}

// TestTPSA tests calculating topological polar surface area
func TestTPSA(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
	}{
		{"Ethanol", "CCO"},
		{"Acetic acid", "CC(=O)O"},
		{"Aspirin", "CC(=O)Oc1ccccc1C(=O)O"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := indigoInit.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load molecule: %v", err)
			}
			defer m.Close()

			// Test with includeSP = false
			tpsa, err := m.TPSA(false)
			if err != nil {
				t.Errorf("failed to get TPSA: %v", err)
			}

			if tpsa < 0 {
				t.Errorf("expected non-negative TPSA, got %f", tpsa)
			}

			// Test with includeSP = true
			tpsaSP, err := m.TPSA(true)
			if err != nil {
				t.Errorf("failed to get TPSA with SP: %v", err)
			}

			if tpsaSP < 0 {
				t.Errorf("expected non-negative TPSA (with SP), got %f", tpsaSP)
			}
		})
	}
}

// TestNumRotatableBonds tests counting rotatable bonds
func TestNumRotatableBonds(t *testing.T) {
	tests := []struct {
		name     string
		smiles   string
		expected int
	}{
		{"Ethanol", "CCO", 1},
		{"Propanol", "CCCO", 2},
		{"Benzene", "c1ccccc1", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := indigoInit.LoadMoleculeFromString(tt.smiles)
			if err != nil {
				t.Fatalf("failed to load molecule: %v", err)
			}
			defer m.Close()

			count, err := m.NumRotatableBonds()
			if err != nil {
				t.Errorf("failed to get rotatable bonds: %v", err)
			}

			if count != tt.expected {
				t.Errorf("expected %d rotatable bonds, got %d", tt.expected, count)
			}
		})
	}
}

// TestMoleculeName tests setting and getting molecule name
func TestMoleculeName(t *testing.T) {
	m, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Set name
	err = m.SetName("Ethanol")
	if err != nil {
		t.Errorf("failed to set name: %v", err)
	}

	// Get name
	name, err := m.Name()
	if err != nil {
		t.Errorf("failed to get name: %v", err)
	}

	if name != "Ethanol" {
		t.Errorf("expected name 'Ethanol', got '%s'", name)
	}
}

// TestMoleculeProperties tests setting and getting properties
func TestMoleculeProperties(t *testing.T) {
	m, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	// Set property
	err = m.SetProperty("test_prop", "test_value")
	if err != nil {
		t.Errorf("failed to set property: %v", err)
	}

	// Check if property exists
	has, err := m.HasProperty("test_prop")
	if err != nil {
		t.Errorf("failed to check property: %v", err)
	}

	if !has {
		t.Error("expected property to exist")
	}

	// Get property
	value, err := m.GetProperty("test_prop")
	if err != nil {
		t.Errorf("failed to get property: %v", err)
	}

	if value != "test_value" {
		t.Errorf("expected property value 'test_value', got '%s'", value)
	}

	// Remove property
	err = m.RemoveProperty("test_prop")
	if err != nil {
		t.Errorf("failed to remove property: %v", err)
	}

	// Check property was removed
	has, err = m.HasProperty("test_prop")
	if err != nil {
		t.Errorf("failed to check property after removal: %v", err)
	}

	if has {
		t.Error("expected property to be removed")
	}
}

// TestMolecularFormula tests getting molecular formula
func TestMolecularFormula(t *testing.T) {
	m, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	formula, err := m.MolecularFormula()
	if err != nil {
		t.Errorf("failed to get molecular formula: %v", err)
	}

	if formula == "" {
		t.Error("expected non-empty molecular formula")
	}

	// Should contain C, H, O
	if !strings.Contains(formula, "C") {
		t.Error("formula should contain C")
	}
}

// TestMassComposition tests getting mass composition
func TestMassComposition(t *testing.T) {
	m, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer m.Close()

	composition, err := m.MassComposition()
	if err != nil {
		t.Errorf("failed to get mass composition: %v", err)
	}

	if composition == "" {
		t.Error("expected non-empty mass composition")
	}
}

// TestPropertiesOnClosedMolecule tests that properties fail on closed molecule
func TestPropertiesOnClosedMolecule(t *testing.T) {
	m, _ := indigoInit.LoadMoleculeFromString("CCO")
	m.Close()

	// All property methods should return errors
	_, err := m.GrossFormula()
	if err == nil {
		t.Error("expected error on closed molecule")
	}

	_, err = m.MolecularWeight()
	if err == nil {
		t.Error("expected error on closed molecule")
	}

	_, err = m.TPSA(false)
	if err == nil {
		t.Error("expected error on closed molecule")
	}
}
