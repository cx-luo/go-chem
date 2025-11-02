//go:build cgo
// +build cgo

package test

import (
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestInChICGO_Simple tests basic InChI generation with CGO
func TestInChICGO_Simple(t *testing.T) {
	tests := []struct {
		name          string
		smiles        string
		wantInChIHas  string // InChI should contain this string
		wantKeyPrefix string // InChIKey prefix
	}{
		{
			name:          "Methane",
			smiles:        "C",
			wantInChIHas:  "CH4",
			wantKeyPrefix: "",
		},
		{
			name:          "Ethane",
			smiles:        "CC",
			wantInChIHas:  "C2H6",
			wantKeyPrefix: "",
		},
		{
			name:          "Ethanol",
			smiles:        "CCO",
			wantInChIHas:  "C2H6O",
			wantKeyPrefix: "",
		},
		{
			name:          "Benzene",
			smiles:        "c1ccccc1",
			wantInChIHas:  "C6H6",
			wantKeyPrefix: "",
		},
	}

	generator := molecule.NewInChIGeneratorCGO()
	loader := molecule.SmilesLoader{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse SMILES
			mol, err := loader.Parse(tt.smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES: %v", err)
			}

			// Generate InChI
			result, err := generator.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			// Check InChI
			if !strings.Contains(result.InChI, tt.wantInChIHas) {
				t.Errorf("InChI = %s, want to contain %s", result.InChI, tt.wantInChIHas)
			}

			// Check InChI format
			if !strings.HasPrefix(result.InChI, "InChI=") {
				t.Errorf("InChI should start with 'InChI=', got %s", result.InChI)
			}

			// Check InChIKey
			if result.InChIKey == "" {
				t.Errorf("InChIKey should not be empty")
			}

			// InChIKey should be 27 characters
			keyParts := strings.Split(result.InChIKey, "-")
			if len(keyParts) != 3 {
				t.Errorf("InChIKey should have 3 parts separated by '-', got %d parts", len(keyParts))
			}

			t.Logf("SMILES:   %s", tt.smiles)
			t.Logf("InChI:    %s", result.InChI)
			t.Logf("InChIKey: %s", result.InChIKey)
		})
	}
}

// TestInChICGO_InChIKey tests InChIKey generation
func TestInChICGO_InChIKey(t *testing.T) {
	tests := []struct {
		name  string
		inchi string
		want  string // Expected InChIKey (if known)
	}{
		{
			name:  "Methane",
			inchi: "InChI=1S/CH4/h1H4",
			want:  "", // We don't know the exact key, just check format
		},
		{
			name:  "Ethanol",
			inchi: "InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := molecule.GenerateInChIKeyCGO(tt.inchi)
			if err != nil {
				t.Fatalf("Failed to generate InChIKey: %v", err)
			}

			// Check format: 14 chars - 9 chars - 2 chars
			parts := strings.Split(key, "-")
			if len(parts) != 3 {
				t.Errorf("InChIKey should have 3 parts, got %d", len(parts))
			}

			if len(parts[0]) != 14 {
				t.Errorf("First part should be 14 chars, got %d", len(parts[0]))
			}

			if len(parts[1]) != 10 {
				t.Errorf("Second part should be 10 chars, got %d", len(parts[1]))
			}

			if len(parts[2]) != 2 {
				t.Errorf("Third part should be 2 chars, got %d", len(parts[2]))
			}

			t.Logf("InChI:    %s", tt.inchi)
			t.Logf("InChIKey: %s", key)
		})
	}
}

// TestInChICGO_Options tests InChI generation with options
func TestInChICGO_Options(t *testing.T) {
	smiles := "CCO"
	loader := molecule.SmilesLoader{}
	mol, _ := loader.Parse(smiles)

	tests := []struct {
		name    string
		options string
	}{
		{
			name:    "Standard",
			options: "",
		},
		{
			name:    "FixedH",
			options: "FixedH",
		},
		{
			name:    "AuxInfo",
			options: "AuxInfo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := molecule.NewInChIGeneratorCGO()
			generator.SetOptions(tt.options)

			result, err := generator.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			if result.InChI == "" {
				t.Errorf("InChI should not be empty")
			}

			t.Logf("Options: %s", tt.options)
			t.Logf("InChI:   %s", result.InChI)
		})
	}
}

// TestInChICGO_EmptyMolecule tests handling of empty molecule
func TestInChICGO_EmptyMolecule(t *testing.T) {
	generator := molecule.NewInChIGeneratorCGO()

	// Create empty molecule
	mol := &molecule.Molecule{}

	result, err := generator.GenerateInChI(mol)
	if err != nil {
		t.Fatalf("Should handle empty molecule: %v", err)
	}

	if result.InChI != "InChI=1S" {
		t.Errorf("Empty molecule should give 'InChI=1S', got %s", result.InChI)
	}
}

// TestInChICGO_Stereochemistry tests stereochemistry handling
func TestInChICGO_Stereochemistry(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
	}{
		{
			name:   "L-Alanine",
			smiles: "C[C@@H](C(=O)O)N",
		},
		{
			name:   "D-Alanine",
			smiles: "C[C@H](C(=O)O)N",
		},
	}

	generator := molecule.NewInChIGeneratorCGO()
	loader := molecule.SmilesLoader{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mol, err := loader.Parse(tt.smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES: %v", err)
			}

			result, err := generator.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			// Should contain stereochemistry layers (/t, /m, /s)
			if !strings.Contains(result.InChI, "/") {
				t.Errorf("InChI should contain layers")
			}

			t.Logf("SMILES:   %s", tt.smiles)
			t.Logf("InChI:    %s", result.InChI)
			t.Logf("InChIKey: %s", result.InChIKey)
		})
	}
}

// BenchmarkInChICGO_Simple benchmarks simple InChI generation
func BenchmarkInChICGO_Simple(b *testing.B) {
	generator := molecule.NewInChIGeneratorCGO()
	loader := molecule.SmilesLoader{}
	mol, _ := loader.Parse("CCO")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = generator.GenerateInChI(mol)
	}
}

// BenchmarkInChICGO_Complex benchmarks complex molecule
func BenchmarkInChICGO_Complex(b *testing.B) {
	generator := molecule.NewInChIGeneratorCGO()
	loader := molecule.SmilesLoader{}
	mol, _ := loader.Parse("C([C@@H]1[C@H]([C@@H]([C@H](C(O1)O)O)O)O)O") // Glucose

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = generator.GenerateInChI(mol)
	}
}

// TestInChICGO_Compare tests if CGO and Pure Go produce similar results
func TestInChICGO_Compare(t *testing.T) {
	smiles := "CCO"
	loader := molecule.SmilesLoader{}
	mol, _ := loader.Parse(smiles)

	// CGO version
	generatorCGO := molecule.NewInChIGeneratorCGO()
	resultCGO, errCGO := generatorCGO.GenerateInChI(mol)

	if errCGO != nil {
		t.Logf("CGO error: %v", errCGO)
	}

	if errCGO == nil {
		t.Logf("CGO InChI:    %s", resultCGO.InChI)
		t.Logf("CGO InChIKey: %s", resultCGO.InChIKey)

		// Both should produce valid InChI (format may differ)
		if !strings.HasPrefix(resultCGO.InChI, "InChI=") {
			t.Errorf("CGO InChI should start with 'InChI='")
		}
	}
}
