package test

import (
	srcpkg "go-chem/src"
	"testing"
)

func TestGrossFormula_Simple(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("C")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	units := srcpkg.CollectGross(m, srcpkg.GrossFormulaOptions{})
	if len(units) != 1 {
		t.Fatalf("expected 1 unit")
	}
	got := srcpkg.GrossToString(units[0].Isotopes, false)
	// carbon default valence 4 -> CH4
	if got != "C H4" && got != "CH4" { // depending on spacing rules
		t.Fatalf("unexpected gross: %s", got)
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
