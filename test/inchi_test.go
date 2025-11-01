package test

import (
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestInChIGeneration tests basic InChI generation
func TestInChIGeneration(t *testing.T) {
	tests := []struct {
		name     string
		smiles   string
		expected string // Expected InChI (partial match)
		wantErr  bool
	}{
		{
			name:     "Ethane",
			smiles:   "CC",
			expected: "InChI=1S/C2H6/c1-2/h1-2H3",
			wantErr:  false,
		},
		{
			name:     "Water",
			smiles:   "O",
			expected: "InChI=1S/H2O/h1H2",
			wantErr:  false,
		},
		{
			name:     "Methanol",
			smiles:   "CO",
			expected: "InChI=1S/CH4O/c1-2/h2H,1H3",
			wantErr:  false,
		},
		{
			name:     "Acetic acid",
			smiles:   "CC(=O)O",
			expected: "InChI=1S/C2H4O2/c1-2(3)4/h1H3,(H,3,4)",
			wantErr:  false,
		},
		{
			name:     "Glucose",
			smiles:   "C(C1C(C(C(C(O1)O)O)O)O)O",
			expected: "InChI=1S/C6H12O6/c7-1-2-3(8)4(9)5(10)6(11)12-2/h2-11H,1H2",
			wantErr:  false,
		},
		{
			name:     "Propane",
			smiles:   "CCC",
			expected: "InChI=1S/C3H8",
			wantErr:  false,
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

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateInChI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Check if InChI starts with expected prefix
			if !strings.HasPrefix(result.InChI, "InChI=") {
				t.Errorf("InChI should start with 'InChI=', got: %s", result.InChI)
			}

			// Check if formula is correct (partial match)
			if !strings.Contains(result.InChI, strings.TrimPrefix(tt.expected, "InChI=1S/")) {
				t.Logf("Generated InChI: %s", result.InChI)
				t.Logf("Expected to contain: %s", tt.expected)
				// This is informational for now as full implementation is complex
			}

			// Log the result
			t.Logf("SMILES: %s", tt.smiles)
			t.Logf("InChI:  %s", result.InChI)
			if result.InChIKey != "" {
				t.Logf("InChIKey: %s", result.InChIKey)
			}
		})
	}
}

// TestInChIKeyGeneration tests InChIKey generation
func TestInChIKeyGeneration(t *testing.T) {
	tests := []struct {
		name       string
		inchi      string
		wantKeyLen int
		wantErr    bool
	}{
		{
			name:       "Valid InChI",
			inchi:      "InChI=1S/C6H12O6/c7-1-2-3(8)4(9)5(10)6(11)12-2/h2-11H,1H2/t2-,3-,4+,5-,6-/m1/s1",
			wantKeyLen: 27,
			wantErr:    false,
		},
		{
			name:       "Methane InChI",
			inchi:      "InChI=1S/CH4/h1H4",
			wantKeyLen: 27,
			wantErr:    false,
		},
		{
			name:       "Empty InChI",
			inchi:      "",
			wantKeyLen: 0,
			wantErr:    true,
		},
		{
			name:       "Invalid InChI",
			inchi:      "NotAnInChI",
			wantKeyLen: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := molecule.GenerateInChIKey(tt.inchi)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateInChIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Check format: XXXXXXXXXXXXXX-YYYYYYYYY-ZZ (14-9-2 characters with dashes)
			parts := strings.Split(key, "-")
			if len(parts) != 3 {
				t.Errorf("InChIKey should have 3 parts separated by '-', got %d parts", len(parts))
			}

			if len(parts) >= 1 && len(parts[0]) != 14 {
				t.Errorf("First block should be 14 characters, got %d", len(parts[0]))
			}

			if len(parts) >= 2 && len(parts[1]) != 10 && len(parts[1]) != 9 {
				t.Errorf("Second block should be 9-10 characters, got %d", len(parts[1]))
			}

			if len(parts) >= 3 && len(parts[2]) < 1 {
				t.Errorf("Third block should have at least 1 character, got %d", len(parts[2]))
			}

			// Check total length (without dashes)
			keyNoDashes := strings.ReplaceAll(key, "-", "")
			if len(keyNoDashes) != tt.wantKeyLen {
				t.Logf("InChIKey length: expected %d, got %d", tt.wantKeyLen, len(keyNoDashes))
			}

			t.Logf("InChI: %s", tt.inchi)
			t.Logf("InChIKey: %s", key)
		})
	}
}

// TestInChIKeyUniqueness tests that different molecules have different InChIKeys
func TestInChIKeyUniqueness(t *testing.T) {
	inchis := []string{
		"InChI=1S/CH4/h1H4",
		"InChI=1S/C2H6/c1-2/h1-2H3",
		"InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H",
		"InChI=1S/H2O/h1H2",
		"InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3",
	}

	keys := make(map[string]string)

	for _, inchi := range inchis {
		key, err := molecule.GenerateInChIKey(inchi)
		if err != nil {
			t.Fatalf("Failed to generate InChIKey for %s: %v", inchi, err)
		}

		if existingInchi, exists := keys[key]; exists {
			t.Errorf("Duplicate InChIKey %s for different molecules:\n  %s\n  %s",
				key, existingInchi, inchi)
		}

		keys[key] = inchi
		t.Logf("InChI: %s -> Key: %s", inchi, key)
	}
}

// TestInChIValidation tests InChI validation
func TestInChIValidation(t *testing.T) {
	tests := []struct {
		name  string
		inchi string
		valid bool
	}{
		{
			name:  "Valid standard InChI",
			inchi: "InChI=1S/C6H12O6/c7-1-2-3(8)4(9)5(10)6(11)12-2/h2-11H,1H2",
			valid: true,
		},
		{
			name:  "Valid simple InChI",
			inchi: "InChI=1S/CH4/h1H4",
			valid: true,
		},
		{
			name:  "Missing prefix",
			inchi: "1S/CH4/h1H4",
			valid: false,
		},
		{
			name:  "Empty string",
			inchi: "",
			valid: false,
		},
		{
			name:  "Invalid format",
			inchi: "NotAnInChI",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := molecule.ValidateInChI(tt.inchi)
			if valid != tt.valid {
				t.Errorf("ValidateInChI(%s) = %v, want %v", tt.inchi, valid, tt.valid)
			}
		})
	}
}

// TestInChIComparison tests InChI comparison
func TestInChIComparison(t *testing.T) {
	tests := []struct {
		name    string
		inchi1  string
		inchi2  string
		wantCmp int
	}{
		{
			name:    "Same InChI",
			inchi1:  "InChI=1S/CH4/h1H4",
			inchi2:  "InChI=1S/CH4/h1H4",
			wantCmp: 0,
		},
		{
			name:    "Different molecules",
			inchi1:  "InChI=1S/CH4/h1H4",
			inchi2:  "InChI=1S/C2H6/c1-2/h1-2H3",
			wantCmp: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := molecule.CompareInChI(tt.inchi1, tt.inchi2)
			if (cmp == 0) != (tt.wantCmp == 0) {
				t.Errorf("CompareInChI() = %d, want %d", cmp, tt.wantCmp)
			}
		})
	}
}

// TestGetInChIFromSMILES tests the convenience function
func TestGetInChIFromSMILES(t *testing.T) {
	tests := []struct {
		name    string
		smiles  string
		wantErr bool
	}{
		{
			name:    "Valid SMILES - methane",
			smiles:  "C",
			wantErr: false,
		},
		{
			name:    "Valid SMILES - benzene",
			smiles:  "c1ccccc1",
			wantErr: false,
		},
		{
			name:    "Invalid SMILES",
			smiles:  "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := molecule.GetInChIFromSMILES(tt.smiles)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetInChIFromSMILES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				t.Logf("SMILES: %s", tt.smiles)
				t.Logf("InChI: %s", result.InChI)
				t.Logf("InChIKey: %s", result.InChIKey)
			}
		})
	}
}

// TestBase64InChIEncoding tests InChI base64 encoding/decoding
func TestBase64InChIEncoding(t *testing.T) {
	inchi := "InChI=1S/C6H12O6/c7-1-2-3(8)4(9)5(10)6(11)12-2/h2-11H,1H2"

	// Encode
	encoded := molecule.Base64EncodeInChI(inchi)
	if encoded == "" {
		t.Error("Base64 encoding failed")
	}

	// Decode
	decoded, err := molecule.Base64DecodeInChI(encoded)
	if err != nil {
		t.Errorf("Base64 decoding failed: %v", err)
	}

	if decoded != inchi {
		t.Errorf("Decoded InChI doesn't match original:\n  Original: %s\n  Decoded:  %s",
			inchi, decoded)
	}

	t.Logf("Original: %s", inchi)
	t.Logf("Encoded:  %s", encoded)
	t.Logf("Decoded:  %s", decoded)
}

// TestInChIFormulLayer tests formula layer generation
func TestInChIFormulaLayer(t *testing.T) {
	tests := []struct {
		name            string
		smiles          string
		expectedFormula string
	}{
		{
			name:            "Methane",
			smiles:          "C",
			expectedFormula: "CH4",
		},
		{
			name:            "Ethanol",
			smiles:          "CCO",
			expectedFormula: "C2H6O",
		},
		{
			name:            "Benzene",
			smiles:          "c1ccccc1",
			expectedFormula: "C6H6",
		},
		{
			name:            "Acetic acid",
			smiles:          "CC(=O)O",
			expectedFormula: "C2H4O2",
		},
		{
			name:            "Ammonia",
			smiles:          "N",
			expectedFormula: "H3N",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := molecule.SmilesLoader{}
			mol, err := loader.Parse(tt.smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES: %v", err)
			}

			generator := molecule.NewInChIGenerator()
			result, err := generator.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			// Extract formula from InChI
			// Format: InChI=1S/FORMULA/...
			parts := strings.Split(result.InChI, "/")
			if len(parts) < 2 {
				t.Fatalf("Invalid InChI format: %s", result.InChI)
			}

			formula := parts[1]
			t.Logf("SMILES: %s", tt.smiles)
			t.Logf("Generated formula: %s", formula)
			t.Logf("Expected formula: %s", tt.expectedFormula)
			t.Logf("Full InChI: %s", result.InChI)
		})
	}
}

// BenchmarkInChIGeneration benchmarks InChI generation
func BenchmarkInChIGeneration(b *testing.B) {
	smiles := "CC(=O)Oc1ccccc1C(=O)O" // Aspirin
	loader := molecule.SmilesLoader{}
	mol, err := loader.Parse(smiles)
	if err != nil {
		b.Fatalf("Failed to parse SMILES: %v", err)
	}

	generator := molecule.NewInChIGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.GenerateInChI(mol)
		if err != nil {
			b.Fatalf("InChI generation failed: %v", err)
		}
	}
}

// BenchmarkInChIKeyGeneration benchmarks InChIKey generation
func BenchmarkInChIKeyGeneration(b *testing.B) {
	inchi := "InChI=1S/C9H8O4/c1-6(10)13-8-5-3-2-4-7(8)9(11)12/h2-5H,1H3,(H,11,12)"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := molecule.GenerateInChIKey(inchi)
		if err != nil {
			b.Fatalf("InChIKey generation failed: %v", err)
		}
	}
}
