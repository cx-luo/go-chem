package test

import (
	"math"
	"strconv"

	srcpkg "github.com/cx-luo/go-chem/molecule"
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

func TestPubChemPropertyRegression(t *testing.T) {
	const smiles = "CCCCN1C(=O)C=CC(=N1)C(=O)NC2=NC(=CS2)C3=CC=C(S3)C"
	m := parseSMILES(t, smiles)

	assertExactInt(t, "HBA", 6, srcpkg.NumHydrogenBondAcceptors(m))
	assertExactInt(t, "HBD", 1, srcpkg.NumHydrogenBondDonors(m))
	assertExactInt(t, "rotatable bonds", 6, srcpkg.NumRotatableBonds(m))
	assertCloseFloat(t, "TPSA", 131.0, m.CalculateTPSA(true), 0.5)
}

func assertExactInt(t *testing.T, field string, want, got int) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %d, want %d", field, got, want)
	}
}

func assertCloseFloat(t *testing.T, field string, want, got, tolerance float64) {
	t.Helper()
	if math.Abs(got-want) > tolerance {
		t.Fatalf("%s: got %s, want %s +/- %s",
			field,
			strconv.FormatFloat(got, 'f', 3, 64),
			strconv.FormatFloat(want, 'f', 3, 64),
			strconv.FormatFloat(tolerance, 'f', 3, 64),
		)
	}
}

func TestNumRotatableBonds(t *testing.T) {
	cases := []struct {
		name   string
		smiles string
		want   int
	}{
		{"ethanol", "CCO", 0},
		{"butane", "CCCC", 1},
		{"cyclohexane", "C1CCCCC1", 0},
		{"biphenyl", "c1ccccc1-c2ccccc2", 1},
		{"explicit_h_butane", "[H]C([H])([H])C([H])([H])C([H])([H])C([H])([H])[H]", 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := parseSMILES(t, tc.smiles)
			got := srcpkg.NumRotatableBonds(m)
			if got != tc.want {
				t.Fatalf("NumRotatableBonds(%s): got %d, want %d", tc.smiles, got, tc.want)
			}
		})
	}
}
