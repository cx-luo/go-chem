package test

import (
	"fmt"
	srcpkg "go-chem/src/molecule"
	"testing"
)

func TestGrossFormula_Simple(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("C")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	units := srcpkg.CollectGross(m, srcpkg.GrossFormulaOptions{})
	for _, unit := range units {
		fmt.Println(unit.Isotopes[0])
	}
	fmt.Println(units)
	if len(units) != 1 {
		t.Fatalf("expected 1 unit")
	}
	hill := srcpkg.GrossUnitsToStringHill(units, false)
	// carbon default valence 4 -> CH4
	if hill != "C H4" && hill != "CH4" { // depending on spacing rules
		t.Fatalf("unexpected gross: %s", hill)
	}
}

func TestGrossFormula_Hill_Benzene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	units := srcpkg.CollectGross(m, srcpkg.GrossFormulaOptions{})
	hill := srcpkg.GrossUnitsToStringHill(units, false)
	// Expect C6 H6 in Hill system
	if hill != "C6 H6" && hill != "C6H6" {
		t.Fatalf("unexpected hill: %s", hill)
	}
}
