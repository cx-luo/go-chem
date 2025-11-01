package test

import (
	"fmt"
	srcpkg "github.com/cx-luo/go-chem/molecule"
	"strings"
	"testing"
)

func TestSMILES(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("CC(=O)OC1=CC=CC=C1C(=O)O")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// Remove verbose output for CI and correctness testing
	fmt.Println(m.CalcMolecularWeight())
	unit := srcpkg.CollectGross(m, srcpkg.GrossFormulaOptions{})
	fmt.Println(srcpkg.GrossUnitsToStringHill(unit, false))
	fmt.Println(m.CalculateMoleculeHash())

	// Aspirin (acetylsalicylic acid): C9H8O4
	expectedC := 9
	expectedO := 4

	// Check atom counts
	cCount, oCount := 0, 0
	for _, atom := range m.Atoms {
		switch atom.Number {
		case srcpkg.ELEM_C:
			cCount++
		case srcpkg.ELEM_O:
			oCount++
		case srcpkg.ELEM_N:
			t.Fatalf("did not expect nitrogen atom in aspirin")
		}
	}
	if cCount != expectedC {
		t.Fatalf("expected %d C atoms, got %d", expectedC, cCount)
	}
	if oCount != expectedO {
		t.Fatalf("expected %d O atoms, got %d", expectedO, oCount)
	}

	// Verify molecular formula
	grossFormula := srcpkg.GrossUnitsToStringHill(unit, false)
	expectedFormula := "C9 H8 O4"
	if grossFormula != expectedFormula {
		t.Fatalf("expected formula %s, got %s", expectedFormula, grossFormula)
	}

	// Check for at least one C=O double bond
	hasDoubleBondedO := false
	for i, bond := range m.Bonds {
		a1 := m.Atoms[bond.Beg]
		a2 := m.Atoms[bond.End]
		order := m.GetBondOrder(i)
		if ((a1.Number == srcpkg.ELEM_C && a2.Number == srcpkg.ELEM_O) ||
			(a2.Number == srcpkg.ELEM_C && a1.Number == srcpkg.ELEM_O)) && order == srcpkg.BOND_DOUBLE {
			hasDoubleBondedO = true
			break // found, no need to continue
		}
	}
	if !hasDoubleBondedO {
		t.Fatalf("expected at least one C=O double bond in aspirin")
	}
}

func TestSMILES_Ethene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("C=C")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(m.Atoms) != 2 || len(m.Bonds) != 1 {
		t.Fatalf("expected 2 atoms and 1 bond, got %d atoms %d bonds", len(m.Atoms), len(m.Bonds))
	}

	fmt.Println(m.CalcMolecularWeight())
	unit := srcpkg.CollectGross(m, srcpkg.GrossFormulaOptions{})
	fmt.Println(srcpkg.GrossUnitsToStringHill(unit, false))
	fmt.Println(m.CalculateMoleculeHash())

	if m.GetBondOrder(0) != srcpkg.BOND_DOUBLE {
		t.Fatalf("expected double bond")
	}
}

func TestSMILES_Benzene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(m.Atoms) != 6 || len(m.Bonds) != 6 {
		t.Fatalf("expected 6 atoms and 6 bonds, got %d atoms %d bonds", len(m.Atoms), len(m.Bonds))
	}

	fmt.Println(m.CalcMolecularWeight())
	unit := srcpkg.CollectGross(m, srcpkg.GrossFormulaOptions{})
	fmt.Println(srcpkg.GrossUnitsToStringHill(unit, false))

	aromatic := 0
	for i := range m.Bonds {
		if m.GetBondOrder(i) == srcpkg.BOND_AROMATIC {
			aromatic++
		}
	}
	if aromatic == 0 {
		t.Fatalf("expected aromatic bonds in benzene")
	}
}

func TestSMILES_ChargedAtom(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("[NH3+]")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(m.Atoms) != 1 {
		t.Fatalf("expected 1 atom, got %d", len(m.Atoms))
	}
	if m.GetAtomCharge(0) != 1 {
		t.Fatalf("expected charge +1, got %d", m.GetAtomCharge(0))
	}
}

func TestSMILES_Isotope(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("[13C]")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	m.SaveJPEG("./test_mol.jpeg", 512, 50)
	if len(m.Atoms) != 1 {
		t.Fatalf("expected 1 atom, got %d", len(m.Atoms))
	}
	if m.GetAtomIsotope(0) != 13 {
		t.Fatalf("expected isotope 13, got %d", m.GetAtomIsotope(0))
	}
}

func TestSMILES_ComplexBracketed(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("[13C@H]")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(m.Atoms) != 1 {
		t.Fatalf("expected 1 atom, got %d", len(m.Atoms))
	}
	if m.GetAtomIsotope(0) != 13 {
		t.Fatalf("expected isotope 13, got %d", m.GetAtomIsotope(0))
	}
}

func TestSMILES_Thiophene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("c1cscc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(m.Atoms) != 5 || len(m.Bonds) != 5 {
		t.Fatalf("expected 5 atoms and 5 bonds, got %d atoms %d bonds", len(m.Atoms), len(m.Bonds))
	}
	m.SaveJPEG("./test_mol.jpeg", 512, 50)
	d := srcpkg.DearomatizerBase{}
	d.Apply(m)
	fmt.Println(m.SaveSMILES())
	// Check that sulfur is present
	sulfurFound := false
	for i := range m.Atoms {
		if m.GetAtomNumber(i) == srcpkg.ELEM_S {
			sulfurFound = true
			break
		}
	}
	if !sulfurFound {
		t.Fatalf("expected sulfur atom in thiophene")
	}
}

func TestSMILES_RoundTrip(t *testing.T) {
	// Test that parsing and saving gives consistent results
	testCases := []string{
		"C=C",
		"c1ccccc1",
		"[NH3+]",
		"[13C]",
		"c1cscc1",
	}

	for _, input := range testCases {
		m, err := (srcpkg.SmilesLoader{}).Parse(input)
		if err != nil {
			t.Fatalf("parse failed for %s: %v", input, err)
		}

		output, err := m.SaveSMILES()
		if err != nil {
			t.Fatalf("SaveSMILES failed for input %s: %v", input, err)
		}
		if output == "" {
			t.Fatalf("SaveSMILES returned empty string for input: %s", input)
		}

		// Parse the output back to verify it's valid
		m2, err := (srcpkg.SmilesLoader{}).Parse(output)
		if err != nil {
			t.Fatalf("round-trip parse failed for %s -> %s: %v", input, output, err)
		}

		// Basic structure should be the same
		if len(m.Atoms) != len(m2.Atoms) {
			t.Fatalf("atom count mismatch for %s: %d vs %d", input, len(m.Atoms), len(m2.Atoms))
		}
		if len(m.Bonds) != len(m2.Bonds) {
			t.Fatalf("bond count mismatch for %s: %d vs %d", input, len(m.Bonds), len(m2.Bonds))
		}
	}
}

func TestSMILES_Output(t *testing.T) {
	// Test specific output formatting
	m, err := (srcpkg.SmilesLoader{}).Parse("C=C")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	output, err := m.SaveSMILES()
	if err != nil {
		t.Fatalf("SaveSMILES failed: %v", err)
	}
	fmt.Println(m.CalculateMoleculeHash())
	fmt.Println(output)
	if output == "" {
		t.Fatalf("SaveSMILES returned empty string")
	}

	// Should contain the double bond
	if !strings.Contains(output, "=") {
		t.Fatalf("expected double bond in output: %s", output)
	}
}
