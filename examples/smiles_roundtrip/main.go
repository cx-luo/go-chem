// Package main demonstrates parsing a SMILES string and writing it back.
//
// Run from the repository root:
//
//	go run ./examples/smiles_roundtrip
package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	input := "CCO"
	mol, err := (molecule.SmilesLoader{}).Parse(input)
	if err != nil {
		log.Fatalf("parse SMILES: %v", err)
	}

	opts := molecule.DefaultSmilesSaverOptions()
	opts.IgnoreHydrogens = true
	output, err := molecule.NewSmilesSaver(opts).SaveSMILES(mol)
	if err != nil {
		log.Fatalf("save SMILES: %v", err)
	}

	formula := molecule.GrossUnitsToStringHill(
		molecule.CollectGross(mol, molecule.GrossFormulaOptions{}),
		false,
	)

	fmt.Printf("Input SMILES:  %s\n", input)
	fmt.Printf("Output SMILES: %s\n", output)
	fmt.Printf("Atoms / Bonds: %d / %d\n", mol.AtomCount(), mol.BondCount())
	fmt.Printf("Formula:       %s\n", formula)
}
