package test

import (
	srcpkg "go-chem/src/molecule"
	"testing"
)

func TestTPSA_And_Lipinski_OnSimpleMolecules(t *testing.T) {
	cases := []struct {
		smiles  string
		minTPSA float64
		maxTPSA float64
		minHBA  int
		maxHBA  int
		minHBD  int
		maxHBD  int
	}{
		{"CCO", 10, 30, 1, 2, 1, 1},     // ethanol
		{"c1ccccc1", 0, 5, 0, 1, 0, 0},  // benzene
		{"CC(=O)C", 10, 30, 1, 2, 0, 1}, // acetone
	}

	for _, tc := range cases {
		m, err := (srcpkg.SmilesLoader{}).Parse(tc.smiles)
		if err != nil {
			t.Fatalf("parse failed for %s: %v", tc.smiles, err)
		}
		tpsa := m.CalculateTPSA(true)
		if tpsa < tc.minTPSA || tpsa > tc.maxTPSA {
			t.Fatalf("TPSA out of range for %s: %f", tc.smiles, tpsa)
		}
		hba := srcpkg.NumHydrogenBondAcceptors(m)
		if hba < tc.minHBA || hba > tc.maxHBA {
			t.Fatalf("HBA out of range for %s: %d", tc.smiles, hba)
		}
		hbd := srcpkg.NumHydrogenBondDonors(m)
		if hbd < tc.minHBD || hbd > tc.maxHBD {
			t.Fatalf("HBD out of range for %s: %d", tc.smiles, hbd)
		}
		rot := srcpkg.NumRotatableBonds(m)
		if rot < 0 {
			t.Fatalf("rotatable bonds negative for %s", tc.smiles)
		}
	}
}
