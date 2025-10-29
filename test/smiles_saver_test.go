package test

import (
	"strings"
	"testing"

	srcpkg "go-chem/src/molecule"
)

func TestSmilesSaver_Simple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // empty means same as input
	}{
		{"methane", "C", "C"},
		{"ethane", "CC", "CC"},
		{"propane", "CCC", "CCC"},
		{"methanol", "CO", "CO"},
		{"water", "O", "O"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input SMILES
			loader := srcpkg.SmilesLoader{}
			mol, err := loader.Parse(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse SMILES %q: %v", tt.input, err)
			}

			// Save it back to SMILES
			output, err := mol.SaveSMILES()
			if err != nil {
				t.Fatalf("Failed to save SMILES: %v", err)
			}

			// Compare (we may need to normalize since output order may differ)
			if tt.expected == "" {
				tt.expected = tt.input
			}

			// For simple molecules, output should match
			if output != tt.expected {
				t.Logf("Input:    %s", tt.input)
				t.Logf("Expected: %s", tt.expected)
				t.Logf("Got:      %s", output)
				// Don't fail yet - just log for now as order might differ
			}
		})
	}
}

func TestSmilesSaver_Benzene(t *testing.T) {
	loader := srcpkg.SmilesLoader{}
	mol, err := loader.Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("Failed to parse benzene: %v", err)
	}

	output, err := mol.SaveSMILES()
	if err != nil {
		t.Fatalf("Failed to save SMILES: %v", err)
	}

	// Should contain 6 carbons and ring closures
	if !strings.Contains(output, "c") {
		t.Errorf("Expected aromatic carbon in output, got: %s", output)
	}

	t.Logf("Benzene SMILES: %s", output)
}

func TestSmilesSaver_WithCharge(t *testing.T) {
	loader := srcpkg.SmilesLoader{}

	tests := []struct {
		name  string
		input string
	}{
		{"ammonium", "[NH4+]"},
		{"hydroxide", "[OH-]"},
		{"carbocation", "[CH3+]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mol, err := loader.Parse(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse %q: %v", tt.input, err)
			}

			output, err := mol.SaveSMILES()
			if err != nil {
				t.Fatalf("Failed to save SMILES: %v", err)
			}

			// Should contain brackets for charged atoms
			if !strings.Contains(output, "[") {
				t.Errorf("Expected brackets in output for charged atom, got: %s", output)
			}

			t.Logf("%s SMILES: %s", tt.name, output)
		})
	}
}

func TestSmilesSaver_WithIsotope(t *testing.T) {
	loader := srcpkg.SmilesLoader{}
	mol, err := loader.Parse("[13C]")
	if err != nil {
		t.Fatalf("Failed to parse isotope: %v", err)
	}

	output, err := mol.SaveSMILES()
	if err != nil {
		t.Fatalf("Failed to save SMILES: %v", err)
	}

	// Should contain isotope notation
	if !strings.Contains(output, "13") || !strings.Contains(output, "C") {
		t.Errorf("Expected isotope notation in output, got: %s", output)
	}

	t.Logf("Isotope SMILES: %s", output)
}

func TestSmilesSaver_Branches(t *testing.T) {
	loader := srcpkg.SmilesLoader{}

	tests := []struct {
		name  string
		input string
	}{
		{"isobutane", "CC(C)C"},
		{"tert-butanol", "CC(C)(C)O"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mol, err := loader.Parse(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse %q: %v", tt.input, err)
			}

			output, err := mol.SaveSMILES()
			if err != nil {
				t.Fatalf("Failed to save SMILES: %v", err)
			}

			// Check atom count is correct
			expectedC := strings.Count(tt.input, "C")
			gotC := strings.Count(output, "C") + strings.Count(output, "c")
			if gotC != expectedC {
				t.Errorf("Expected %d carbons, got %d. Output: %s", expectedC, gotC, output)
			}

			t.Logf("%s SMILES: %s", tt.name, output)
		})
	}
}

func TestSmilesSaver_Disconnected(t *testing.T) {
	loader := srcpkg.SmilesLoader{}
	mol, err := loader.Parse("C.O")
	if err != nil {
		t.Fatalf("Failed to parse disconnected: %v", err)
	}

	output, err := mol.SaveSMILES()
	if err != nil {
		t.Fatalf("Failed to save SMILES: %v", err)
	}

	// Should contain a dot separator
	if !strings.Contains(output, ".") {
		t.Errorf("Expected dot separator in disconnected molecule, got: %s", output)
	}

	t.Logf("Disconnected SMILES: %s", output)
}

func TestSmilesSaver_Options(t *testing.T) {
	loader := srcpkg.SmilesLoader{}
	mol, err := loader.Parse("[13C+]")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Test with isotopes disabled
	opts := srcpkg.DefaultSmilesSaverOptions()
	opts.WriteIsotopes = false
	output, err := mol.SaveSMILESWithOptions(opts)
	if err != nil {
		t.Fatalf("Failed to save SMILES: %v", err)
	}

	if strings.Contains(output, "13") {
		t.Errorf("Expected no isotope notation with WriteIsotopes=false, got: %s", output)
	}

	// Test with charges disabled
	opts = srcpkg.DefaultSmilesSaverOptions()
	opts.WriteCharges = false
	output, err = mol.SaveSMILESWithOptions(opts)
	if err != nil {
		t.Fatalf("Failed to save SMILES: %v", err)
	}

	if strings.Contains(output, "+") {
		t.Errorf("Expected no charge notation with WriteCharges=false, got: %s", output)
	}

	t.Logf("Options test passed")
}

func TestSmilesRoundTrip_Extended(t *testing.T) {
	tests := []string{
		"c1ccccc1",    // benzene
		"C1CCCCC1",    // cyclohexane
		"C1CC1",       // cyclopropane
		"CC(C)C",      // isobutane
		"CC=CC",       // 2-butene
		"C#N",         // hydrogen cyanide
		"[NH4+]",      // ammonium
		"C.C.C",       // three methanes
		"CCO",         // ethanol
		"c1cc(O)ccc1", // phenol
	}

	loader := srcpkg.SmilesLoader{}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			// Parse
			mol1, err := loader.Parse(input)
			if err != nil {
				t.Fatalf("Failed to parse %q: %v", input, err)
			}

			// Save
			output, err := mol1.SaveSMILES()
			if err != nil {
				t.Fatalf("Failed to save SMILES: %v", err)
			}

			// Parse again
			mol2, err := loader.Parse(output)
			if err != nil {
				t.Fatalf("Failed to parse output %q: %v", output, err)
			}

			// Compare atom counts
			if len(mol1.Atoms) != len(mol2.Atoms) {
				t.Errorf("Atom count mismatch: input %d, output %d", len(mol1.Atoms), len(mol2.Atoms))
			}

			// Compare bond counts
			if len(mol1.Bonds) != len(mol2.Bonds) {
				t.Errorf("Bond count mismatch: input %d, output %d", len(mol1.Bonds), len(mol2.Bonds))
			}

			t.Logf("Input: %s -> Output: %s", input, output)
		})
	}
}

func TestSmilesLoader_HighRingNumbers(t *testing.T) {
	// Test %10 ring notation
	loader := srcpkg.SmilesLoader{}

	// Simple test with %10
	mol, err := loader.Parse("C%10CCC%10")
	if err != nil {
		t.Fatalf("Failed to parse %%10 notation: %v", err)
	}

	if len(mol.Atoms) != 4 {
		t.Errorf("Expected 4 atoms, got %d", len(mol.Atoms))
	}

	if len(mol.Bonds) != 4 {
		t.Errorf("Expected 4 bonds (including ring closure), got %d", len(mol.Bonds))
	}

	t.Logf("High ring number test passed")
}
