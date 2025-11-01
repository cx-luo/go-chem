package test

import (
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

func TestSmilesStereochemistry(t *testing.T) {
	tests := []struct {
		name    string
		smiles  string
		wantErr bool
	}{
		{
			name:    "trans-1,2-dichloroethene",
			smiles:  `Cl/C=C\Cl`,
			wantErr: false,
		},
		{
			name:    "cis-1,2-dichloroethene",
			smiles:  `Cl/C=C/Cl`,
			wantErr: false,
		},
		{
			name:    "trans-2-butene",
			smiles:  `C/C=C\C`,
			wantErr: false,
		},
		{
			name:    "cis-2-butene",
			smiles:  `C/C=C/C`,
			wantErr: false,
		},
		{
			name:    "trans-styrene",
			smiles:  `C/C=C\C1=CC=CC=C1`,
			wantErr: false,
		},
		{
			name:    "backslash first trans",
			smiles:  `ClC\C=C\Cl`,
			wantErr: false,
		},
		{
			name:    "mixed with branches",
			smiles:  `CC(/C)=C\C`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := molecule.SmilesLoader{}
			mol, err := loader.Parse(tt.smiles)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for SMILES %s, but got none", tt.smiles)
				}
				return
			}

			if err != nil {
				t.Errorf("Failed to parse SMILES %s: %v", tt.smiles, err)
				return
			}

			if mol == nil {
				t.Errorf("Got nil molecule for SMILES %s", tt.smiles)
				return
			}

			// Check that molecule has at least one double bond
			hasDoubleBond := false
			for _, bond := range mol.Bonds {
				if bond.Order == molecule.BOND_DOUBLE {
					hasDoubleBond = true
					break
				}
			}

			if !hasDoubleBond {
				t.Errorf("Expected double bond in SMILES %s", tt.smiles)
			}

			t.Logf("Successfully parsed %s: %d atoms, %d bonds", tt.smiles, len(mol.Atoms), len(mol.Bonds))
		})
	}
}

func TestSmilesStereochemistryDirection(t *testing.T) {
	// Test that bond directions are properly set
	loader := molecule.SmilesLoader{}

	// trans: different directions
	mol, err := loader.Parse(`C/C=C\C`)
	if err != nil {
		t.Fatalf("Failed to parse trans-2-butene: %v", err)
	}

	// Find the double bond
	doubleBondIdx := -1
	for i, bond := range mol.Bonds {
		if bond.Order == molecule.BOND_DOUBLE {
			doubleBondIdx = i
			break
		}
	}

	if doubleBondIdx < 0 {
		t.Fatal("No double bond found in trans-2-butene")
	}

	// Check that neighboring single bonds have direction info
	doubleBond := mol.Bonds[doubleBondIdx]
	begNeighbors := mol.GetNeighborBonds(doubleBond.Beg)
	endNeighbors := mol.GetNeighborBonds(doubleBond.End)

	foundDirection := false
	for _, idx := range begNeighbors {
		if idx != doubleBondIdx && mol.Bonds[idx].Direction != 0 {
			foundDirection = true
			t.Logf("Found direction %d on bond %d (begin side)", mol.Bonds[idx].Direction, idx)
		}
	}
	for _, idx := range endNeighbors {
		if idx != doubleBondIdx && mol.Bonds[idx].Direction != 0 {
			foundDirection = true
			t.Logf("Found direction %d on bond %d (end side)", mol.Bonds[idx].Direction, idx)
		}
	}

	if !foundDirection {
		t.Error("Expected to find bond direction information for stereochemistry")
	}
}
