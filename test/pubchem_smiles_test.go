package test

import (
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

func TestPubChemStyleComplexSMILES(t *testing.T) {
	tests := []struct {
		name        string
		smiles      string
		formula     string
		atoms       int
		bonds       int
		minAromatic int
	}{
		{
			name:        "phenol leading oxygen before aromatic atom",
			smiles:      "Oc1ccccc1",
			formula:     "C6H6O",
			atoms:       7,
			bonds:       7,
			minAromatic: 6,
		},
		{
			name:    "glucose",
			smiles:  "C(C1C(C(C(C(O1)O)O)O)O)O",
			formula: "C6H12O6",
			atoms:   12,
			bonds:   12,
		},
		{
			name:    "L-alanine isomeric",
			smiles:  "C[C@@H](C(=O)O)N",
			formula: "C3H7NO2",
			atoms:   6,
			bonds:   5,
		},
		{
			name:        "caffeine",
			smiles:      "CN1C=NC2=C1C(=O)N(C(=O)N2C)C",
			formula:     "C8H10N4O2",
			atoms:       14,
			bonds:       15,
			minAromatic: 0,
		},
		{
			name:        "ibuprofen",
			smiles:      "CC(C)CC1=CC=C(C=C1)C(C)C(=O)O",
			formula:     "C13H18O2",
			atoms:       15,
			bonds:       15,
			minAromatic: 0,
		},
		{
			name:    "sodium acetate salt",
			smiles:  "[Na+].CC(=O)[O-]",
			formula: "C2H3NaO2",
			atoms:   5,
			bonds:   3,
		},
		{
			name:        "indole aromatic NH",
			smiles:      "c1ccc2[nH]ccc2c1",
			formula:     "C8H7N",
			atoms:       9,
			bonds:       10,
			minAromatic: 10,
		},
		{
			name:        "aromatic selenium two-letter atom",
			smiles:      "c1c[se]cc1",
			formula:     "C4H4Se",
			atoms:       5,
			bonds:       5,
			minAromatic: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mol, err := (molecule.SmilesLoader{}).Parse(tt.smiles)
			if err != nil {
				t.Fatalf("parse %s: %v", tt.smiles, err)
			}
			if mol.AtomCount() != tt.atoms {
				t.Fatalf("atom count for %s: got %d, want %d", tt.smiles, mol.AtomCount(), tt.atoms)
			}
			if mol.BondCount() != tt.bonds {
				t.Fatalf("bond count for %s: got %d, want %d", tt.smiles, mol.BondCount(), tt.bonds)
			}
			if got := compactFormula(mol); got != tt.formula {
				t.Fatalf("formula for %s: got %s, want %s", tt.smiles, got, tt.formula)
			}
			if tt.minAromatic > 0 && countOrder(mol, molecule.BOND_AROMATIC) < tt.minAromatic {
				t.Fatalf("aromatic bond count for %s: got %d, want at least %d", tt.smiles, countOrder(mol, molecule.BOND_AROMATIC), tt.minAromatic)
			}
		})
	}
}

func TestSMILESBracketAtomMapsAndAromaticColonBonds(t *testing.T) {
	tests := []struct {
		name   string
		smiles string
		atoms  int
		bonds  int
	}{
		{name: "atom map", smiles: "[CH3:1][CH2:2][OH:3]", atoms: 3, bonds: 2},
		{name: "explicit aromatic colon", smiles: "c1:c:c:c:c:c:1", atoms: 6, bonds: 6},
		{name: "wildcard pseudo atom", smiles: "*CC", atoms: 3, bonds: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mol, err := (molecule.SmilesLoader{}).Parse(tt.smiles)
			if err != nil {
				t.Fatalf("parse %s: %v", tt.smiles, err)
			}
			if mol.AtomCount() != tt.atoms || mol.BondCount() != tt.bonds {
				t.Fatalf("parsed %s as %d atoms/%d bonds, want %d/%d", tt.smiles, mol.AtomCount(), mol.BondCount(), tt.atoms, tt.bonds)
			}
		})
	}
}

func compactFormula(mol *molecule.Molecule) string {
	unit := molecule.CollectGross(mol, molecule.GrossFormulaOptions{})
	return strings.ReplaceAll(molecule.GrossUnitsToStringHill(unit, false), " ", "")
}

func countOrder(mol *molecule.Molecule, order int) int {
	count := 0
	for i := range mol.Bonds {
		if mol.GetBondOrder(i) == order {
			count++
		}
	}
	return count
}
