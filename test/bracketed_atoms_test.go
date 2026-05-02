package test

import (
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestBracketedAtomsWithoutExplicitH tests that bracketed atoms without
// explicit H specification (like [O], [N+]) correctly have 0 implicit hydrogens,
// not calculated by valence rules.
func TestBracketedAtomsWithoutExplicitH(t *testing.T) {
	tests := []struct {
		name     string
		smiles   string
		expected string
	}{
		{
			name:     "Hydroperoxyl radical",
			smiles:   "O[O]",
			expected: "HO2",
		},
		{
			name:     "TEMPO derivative with ketone",
			smiles:   "CC1(CC(=O)CC(N1[O])(C)C)C",
			expected: "C9H16NO2",
		},
		{
			name:     "Glycol peroxide",
			smiles:   "C(CO[O])O",
			expected: "C2H5O3",
		},
		{
			name:     "TEMPO with amine",
			smiles:   "CC1(CC(CC(N1[O])(C)C)N)C",
			expected: "C9H19N2O",
		},
		{
			name:     "TEMPO with chloroacetamide",
			smiles:   "CC1(CC(CC(N1[O])(C)C)NC(=O)CCl)C",
			expected: "C11H20ClN2O2",
		},
		{
			name:     "Simple TEMPO",
			smiles:   "CC1(CCCC(N1[O])(C)C)C",
			expected: "C9H18NO",
		},
		{
			name:     "TEMPO with maleimide",
			smiles:   "CC1(CC(CC(N1[O])(C)C)N2C(=O)C=CC2=O)C",
			expected: "C13H19N2O3",
		},
		{
			name:     "TEMPO with iodoacetamide",
			smiles:   "CC1(CC(CC(N1[O])(C)C)NC(=O)CI)C",
			expected: "C11H20IN2O2",
		},
		{
			name:     "Complex salt with potassium",
			smiles:   "CC1(C([N+](=C(N1[O])C2=CC=C(C=C2)C(=O)[O-])[O-])(C)C)C.[K+]",
			expected: "C14H16KN2O4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := molecule.SmilesLoader{}
			mol, err := loader.Parse(tt.smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES %s: %v", tt.smiles, err)
			}

			opts := molecule.GrossFormulaOptions{AddIsotopes: false, AddRSites: false}
			units := molecule.CollectGross(mol, opts)
			formula := molecule.GrossUnitsToStringHill(units, false)

			// Remove spaces for comparison
			formula = strings.ReplaceAll(formula, " ", "")
			expected := strings.ReplaceAll(tt.expected, " ", "")

			if formula != expected {
				t.Errorf("SMILES: %s\nExpected formula: %s\nActual formula: %s", tt.smiles, expected, formula)
			}
		})
	}
}

// TestBracketedVsUnbracketedAtoms verifies the difference in behavior between
// bracketed and unbracketed atoms for implicit hydrogen calculation.
func TestBracketedVsUnbracketedAtoms(t *testing.T) {
	tests := []struct {
		name            string
		smiles          string
		expectedFormula string
		description     string
	}{
		{
			name:            "Unbracketed O gets implicit H",
			smiles:          "O",
			expectedFormula: "H2O",
			description:     "Plain O calculates 2 implicit H by valence",
		},
		{
			name:            "Bracketed [O] has no H",
			smiles:          "[O]",
			expectedFormula: "O",
			description:     "[O] explicitly specifies 0 H",
		},
		{
			name:            "Bracketed [OH] has 1 H",
			smiles:          "[OH]",
			expectedFormula: "HO",
			description:     "[OH] explicitly specifies 1 H",
		},
		{
			name:            "Unbracketed N gets implicit H",
			smiles:          "N",
			expectedFormula: "H3N",
			description:     "Plain N calculates 3 implicit H by valence",
		},
		{
			name:            "Bracketed [N] has no H",
			smiles:          "[N]",
			expectedFormula: "N",
			description:     "[N] explicitly specifies 0 H",
		},
		{
			name:            "Bracketed [NH2] has 2 H",
			smiles:          "[NH2]",
			expectedFormula: "H2N",
			description:     "[NH2] explicitly specifies 2 H",
		},
		{
			name:            "Charged [N+] has no H",
			smiles:          "[N+]",
			expectedFormula: "N+",
			description:     "[N+] explicitly specifies 0 H",
		},
		{
			name:            "Charged [NH4+] has 4 H",
			smiles:          "[NH4+]",
			expectedFormula: "H4N+",
			description:     "[NH4+] explicitly specifies 4 H",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := molecule.SmilesLoader{}
			mol, err := loader.Parse(tt.smiles)
			if err != nil {
				t.Fatalf("Failed to parse SMILES %s: %v", tt.smiles, err)
			}

			opts := molecule.GrossFormulaOptions{AddIsotopes: false, AddRSites: false}
			units := molecule.CollectGross(mol, opts)
			formula := molecule.GrossUnitsToStringHill(units, false)

			// Remove spaces for comparison
			formula = strings.ReplaceAll(formula, " ", "")
			expected := strings.ReplaceAll(tt.expectedFormula, " ", "")

			if formula != expected {
				t.Errorf("%s\nSMILES: %s\nExpected: %s\nActual: %s",
					tt.description, tt.smiles, expected, formula)
			}
		})
	}
}
