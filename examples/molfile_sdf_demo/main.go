// Package main demonstrates MOL and SDF loading with a small in-memory molecule.
//
// Run from the repository root:
//
//	go run ./examples/molfile_sdf_demo
package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	source := aceticAcid()

	molfile, err := molecule.SaveMoleculeToString(source)
	if err != nil {
		log.Fatalf("save MOL: %v", err)
	}

	loaded, err := molecule.LoadMoleculeFromString(molfile)
	if err != nil {
		log.Fatalf("load MOL: %v", err)
	}

	sdf := molfile + "\n> <NAME>\nAcetic acid\n\n$$$$\n"
	molecules, err := molecule.LoadSDFFromString(sdf)
	if err != nil {
		log.Fatalf("load SDF: %v", err)
	}

	fmt.Printf("MOL name:      %s\n", loaded.Name)
	fmt.Printf("MOL atoms:     %d\n", loaded.AtomCount())
	fmt.Printf("MOL bonds:     %d\n", loaded.BondCount())
	fmt.Printf("SDF molecules: %d\n", len(molecules))
}

func aceticAcid() *molecule.Molecule {
	mol := molecule.NewMolecule()
	mol.Name = "Acetic acid"

	c1 := mol.AddAtom(molecule.ELEM_C)
	c2 := mol.AddAtom(molecule.ELEM_C)
	o1 := mol.AddAtom(molecule.ELEM_O)
	o2 := mol.AddAtom(molecule.ELEM_O)

	mol.AddBond(c1, c2, molecule.BOND_SINGLE)
	mol.AddBond(c2, o1, molecule.BOND_DOUBLE)
	mol.AddBond(c2, o2, molecule.BOND_SINGLE)

	mol.SetAtomXY(c1, 0, 0)
	mol.SetAtomXY(c2, 1.5, 0)
	mol.SetAtomXY(o1, 2.25, 1.2)
	mol.SetAtomXY(o2, 2.25, -1.2)

	return mol
}
