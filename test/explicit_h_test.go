package test

import (
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

func TestExplicitHydrogenInBrackets(t *testing.T) {
	tests := []struct {
		smiles   string
		expected string
		name     string
	}{
		{
			smiles:   "CCCC[SnH](CCCC)CCCC",
			expected: "C12H28Sn",
			name:     "Tin with explicit hydrogen",
		},
		{
			smiles:   "[CH-]1C=CC=C1.[CH-]1C=CC=C1.[Fe+2]",
			expected: "C10H10Fe",
			name:     "Ferrocene with explicit H on charged carbons",
		},
		{
			smiles:   "[NH3+]",
			expected: "H3N+",
			name:     "Ammonium ion with explicit H count",
		},
		{
			smiles:   "[CH3-]",
			expected: "CH3-",
			name:     "Methyl anion with explicit H count",
		},
		{
			smiles:   "C[SnH2]C",
			expected: "C2H8Sn",
			name:     "Tin with 2 explicit hydrogens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := molecule.SmilesLoader{}
			mol, err := loader.Parse(tt.smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES %s: %v", tt.smiles, err)
			}

			// Get molecular formula
			opts := molecule.GrossFormulaOptions{AddIsotopes: false, AddRSites: false}
			units := molecule.CollectGross(mol, opts)
			formula := molecule.GrossUnitsToStringHill(units, false)

			// Remove spaces from formula for comparison
			formula = stripSpaces(formula)
			expected := stripSpaces(tt.expected)

			if formula != expected {
				t.Errorf("SMILES: %s\nExpected formula: %s\nActual formula: %s", tt.smiles, expected, formula)
			}
		})
	}
}

func stripSpaces(s string) string {
	result := ""
	for _, ch := range s {
		if ch != ' ' {
			result += string(ch)
		}
	}
	return result
}
