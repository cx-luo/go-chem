package test

import (
	"fmt"
	srcpkg "go-chem/src"
	"testing"
)

func TestSMILES_Ethene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("C=C")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(m.Atoms) != 2 || len(m.Edges) != 1 {
		t.Fatalf("expected 2 atoms and 1 bond, got %d atoms %d bonds", len(m.Atoms), len(m.Edges))
	}
	if m.GetBondOrder(0) != srcpkg.BOND_DOUBLE {
		t.Fatalf("expected double bond")
	}
}

func TestSMILES_Benzene(t *testing.T) {
	m, err := (srcpkg.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(m.Atoms) != 6 || len(m.Edges) != 6 {
		t.Fatalf("expected 6 atoms and 6 bonds, got %d atoms %d bonds", len(m.Atoms), len(m.Edges))
	}

	err = m.SavePNG("test_smiles.png", 512)
	if err != nil {
		fmt.Println(err)
	}
	aromatic := 0
	for i := range m.Edges {
		if m.GetBondOrder(i) == srcpkg.BOND_AROMATIC {
			aromatic++
		}
	}
	if aromatic == 0 {
		t.Fatalf("expected aromatic bonds in benzene")
	}
}
