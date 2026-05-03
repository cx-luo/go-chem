package test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/render2d"
)

func TestIndigoCoreMoleculePropertyData(t *testing.T) {
	records := readTSV(t, "molecule_properties.tsv")
	for _, row := range records {
		name := row[0]
		t.Run(name, func(t *testing.T) {
			mol := parseSMILES(t, row[1])

			assertEqualInt(t, "atoms", atoi(t, row[2]), mol.AtomCount())
			assertEqualInt(t, "bonds", atoi(t, row[3]), mol.BondCount())
			assertEqualInt(t, "aromatic bonds", atoi(t, row[4]), countBonds(mol, molecule.BOND_AROMATIC))
			assertEqualInt(t, "HBA", atoi(t, row[5]), molecule.NumHydrogenBondAcceptors(mol))
			assertEqualInt(t, "HBD", atoi(t, row[6]), molecule.NumHydrogenBondDonors(mol))
			assertEqualInt(t, "rotatable bonds", atoi(t, row[7]), molecule.NumRotatableBonds(mol))

			tpsa := mol.CalculateTPSA(true)
			minTPSA := atof(t, row[8])
			maxTPSA := atof(t, row[9])
			if tpsa < minTPSA || tpsa > maxTPSA {
				t.Fatalf("TPSA for %s out of range: got %.2f, want %.2f..%.2f", row[1], tpsa, minTPSA, maxTPSA)
			}
		})
	}
}

func TestIndigoCoreMCSData(t *testing.T) {
	records := readTSV(t, "mcs_cases.tsv")
	for _, row := range records {
		name := row[0]
		t.Run(name, func(t *testing.T) {
			query := parseSMILES(t, row[1])
			target := parseSMILES(t, row[2])
			mcs := molecule.NewMaxCommonSubstructure(query, target).Find()

			if mcs.AtomCount() < atoi(t, row[3]) {
				t.Fatalf("MCS atoms: got %d, want at least %s", mcs.AtomCount(), row[3])
			}
			if mcs.BondCount() < atoi(t, row[4]) {
				t.Fatalf("MCS bonds: got %d, want at least %s", mcs.BondCount(), row[4])
			}
		})
	}
}

func TestIndigoCoreSubstructureData(t *testing.T) {
	records := readTSV(t, "substructure_cases.tsv")
	for _, row := range records {
		name := row[0]
		t.Run(name, func(t *testing.T) {
			query := parseSMILES(t, row[1])
			target := parseSMILES(t, row[2])
			matches := molecule.FindSubstructureMatches(query, target)
			assertEqualInt(t, "unique matches", atoi(t, row[3]), len(matches))
		})
	}
}

func TestIndigoCoreRender2DData(t *testing.T) {
	records := readTSV(t, "render2d_cases.tsv")
	for _, row := range records {
		name := row[0]
		t.Run(name, func(t *testing.T) {
			mol := parseSMILES(t, row[1])
			svg, err := render2d.RenderSVG(mol, render2d.Options{
				Width:  atoi(t, row[2]),
				Height: atoi(t, row[3]),
			})
			if err != nil {
				t.Fatalf("RenderSVG failed: %v", err)
			}
			for _, token := range strings.Split(row[4], "|") {
				if !strings.Contains(svg, token) {
					t.Fatalf("SVG output should contain %q, got:\n%s", token, svg)
				}
			}
		})
	}
}

func readTSV(t *testing.T, name string) [][]string {
	t.Helper()

	path := filepath.Join("data", "indigo_core", name)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.Comment = '#'
	r.FieldsPerRecord = -1

	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	for i, row := range rows {
		for j := range row {
			row[j] = strings.TrimSpace(row[j])
		}
		rows[i] = row
	}
	return rows
}

func parseSMILES(t *testing.T, smiles string) *molecule.Molecule {
	t.Helper()
	mol, err := (molecule.SmilesLoader{}).Parse(smiles)
	if err != nil {
		t.Fatalf("parse %s: %v", smiles, err)
	}
	return mol
}

func countBonds(mol *molecule.Molecule, order int) int {
	count := 0
	for i := range mol.Bonds {
		if mol.GetBondOrder(i) == order {
			count++
		}
	}
	return count
}

func atoi(t *testing.T, s string) int {
	t.Helper()
	v, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("parse int %q: %v", s, err)
	}
	return v
}

func atof(t *testing.T, s string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		t.Fatalf("parse float %q: %v", s, err)
	}
	return v
}

func assertEqualInt(t *testing.T, field string, want, got int) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %d, want %d", field, got, want)
	}
}
