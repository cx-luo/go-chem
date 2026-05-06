// Package main demonstrates saving molecules and reactions to KET JSON.
//
// Run from the repository root:
//
//	go run ./examples/ket_saver_demo
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/reaction"
)

func main() {
	mol := carbonylMolecule()
	ket, err := molecule.SaveMoleculeToKET(mol)
	if err != nil {
		log.Fatalf("save molecule KET: %v", err)
	}
	printMoleculeKETSummary(ket)

	rxn := saltReaction()
	reactionKET, err := reaction.SaveReactionToKETString(rxn, false)
	if err != nil {
		log.Fatalf("save reaction KET: %v", err)
	}
	printReactionKETSummary(reactionKET)
}

func carbonylMolecule() *molecule.Molecule {
	mol := molecule.NewMolecule()
	c := mol.AddAtom(molecule.ELEM_C)
	o := mol.AddAtom(molecule.ELEM_O)
	mol.AddBond(c, o, molecule.BOND_DOUBLE)
	mol.SetAtomXY(c, 0, 0)
	mol.SetAtomXY(o, 1.25, 0)
	return mol
}

func saltReaction() *reaction.Reaction {
	reactant := molecule.NewMolecule()
	cl := reactant.AddAtom(molecule.ELEM_Cl)
	reactant.SetAtomXY(cl, 0, 0)

	product := molecule.NewMolecule()
	na := product.AddAtom(molecule.ELEM_Na)
	prodCl := product.AddAtom(molecule.ELEM_Cl)
	product.AddBond(na, prodCl, molecule.BOND_SINGLE)
	product.SetAtomXY(na, 0, 0)
	product.SetAtomXY(prodCl, 1.5, 0)

	rxn := reaction.NewReaction()
	reactantIdx := rxn.AddReactant()
	productIdx := rxn.AddProduct()
	rxn.GetMolecule(reactantIdx).SetStructure(reactant)
	rxn.GetMolecule(productIdx).SetStructure(product)
	return rxn
}

func printMoleculeKETSummary(input string) {
	var doc struct {
		Mol0 struct {
			Atoms []struct {
				Label string `json:"label"`
			} `json:"atoms"`
			Bonds []struct {
				Type int `json:"type"`
			} `json:"bonds"`
		} `json:"mol0"`
	}
	if err := json.Unmarshal([]byte(input), &doc); err != nil {
		log.Fatalf("decode molecule KET: %v", err)
	}
	fmt.Printf("Molecule KET: %d atoms, %d bonds, first atom %s\n",
		len(doc.Mol0.Atoms), len(doc.Mol0.Bonds), doc.Mol0.Atoms[0].Label)
}

func printReactionKETSummary(input string) {
	var doc struct {
		Root struct {
			Nodes []struct {
				Ref  string `json:"$ref"`
				Type string `json:"type"`
			} `json:"nodes"`
		} `json:"root"`
	}
	if err := json.Unmarshal([]byte(input), &doc); err != nil {
		log.Fatalf("decode reaction KET: %v", err)
	}
	fmt.Printf("Reaction KET: %d root nodes (%s, %s, %s)\n",
		len(doc.Root.Nodes), doc.Root.Nodes[0].Ref, doc.Root.Nodes[1].Ref, doc.Root.Nodes[2].Type)
}
