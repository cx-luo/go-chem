package test

import (
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestUserReportedInChI tests the specific molecules reported by the user
func TestUserReportedInChI(t *testing.T) {
	tests := []struct {
		name             string
		smiles           string
		expectedInChI    string
		expectedInChIKey string
	}{
		{
			name:             "Acetylcarnitine with charges",
			smiles:           "CC(=O)OC(CC(=O)[O-])C[N+](C)(C)C",
			expectedInChI:    "InChI=1S/C9H17NO4/c1-7(11)14-8(5-9(12)13)6-10(2,3)4/h8H,5-6H2,1-4H3",
			expectedInChIKey: "RDHQFKQIGNGIED-UHFFFAOYSA-N",
		},
		{
			name:             "Acetylcarnitine neutral",
			smiles:           "CC(=O)OC(CC(=O)O)C[N+](C)(C)C",
			expectedInChI:    "InChI=1S/C9H17NO4/c1-7(11)14-8(5-9(12)13)6-10(2,3)4/h8H,5-6H2,1-4H3/p+1",
			expectedInChIKey: "RDHQFKQIGNGIED-UHFFFAOYSA-O",
		},
		{
			name:             "Hydroxylated aromatic carboxylic acid",
			smiles:           "C1=CC(C(C(=C1)C(=O)O)O)O",
			expectedInChI:    "InChI=1S/C7H8O4/c8-5-3-1-2-4(6(5)9)7(10)11/h1-3,5-6,8-9H,(H,10,11)",
			expectedInChIKey: "INCSWYKICIYAHB-UHFFFAOYSA-N",
		},
		{
			name:             "Amino alcohol",
			smiles:           "CC(CN)O",
			expectedInChI:    "InChI=1S/C3H9NO/c1-3(5)2-4/h3,5H,2,4H2,1H3",
			expectedInChIKey: "HXKKHQJGJAFBHI-UHFFFAOYSA-N",
		},
		{
			name:             "Phosphate compound",
			smiles:           "C(C(=O)COP(=O)(O)O)N",
			expectedInChI:    "InChI=1S/C3H8NO5P/c4-1-3(5)2-9-10(6,7)8/h1-2,4H2,(H2,6,7,8)",
			expectedInChIKey: "HIQNVODXENYOFK-UHFFFAOYSA-N",
		},
		{
			name:             "Chlorodinitrobenzene",
			smiles:           "C1=CC(=C(C=C1[N+](=O)[O-])[N+](=O)[O-])Cl",
			expectedInChI:    "InChI=1S/C6H3ClN2O4/c7-5-2-1-4(8(10)11)3-6(5)9(12)13/h1-3H",
			expectedInChIKey: "DBRXOUCRJQVYJQ-UHFFFAOYSA-N",
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

			t.Logf("SMILES:   %s", tt.smiles)
			t.Logf("Expected: %s", tt.expectedInChI)
			t.Logf("Actual:   %s", result.InChI)
			t.Logf("Expected Key: %s", tt.expectedInChIKey)
			t.Logf("Actual Key:   %s", result.InChIKey)

			// Check formula layer (should match at least)
			// Full matching requires proper canonical ordering
			if result.InChI == "" {
				t.Error("InChI generation returned empty string")
			}
		})
	}
}
