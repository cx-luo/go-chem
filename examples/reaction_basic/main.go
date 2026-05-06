// Package main demonstrates constructing a simple reaction.
//
// Run from the repository root:
//
//	go run ./examples/reaction_basic
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/reaction"
)

func main() {
	reactant := mustParse("CCO")
	product := mustParse("CC=O")

	rxn := reaction.NewReaction()
	reactantIdx := rxn.AddReactant()
	productIdx := rxn.AddProduct()
	rxn.GetMolecule(reactantIdx).SetStructure(reactant)
	rxn.GetMolecule(productIdx).SetStructure(product)

	ket, err := reaction.SaveReactionToKETString(rxn, false)
	if err != nil {
		log.Fatalf("save reaction KET: %v", err)
	}

	fmt.Printf("Reactants / Products: %d / %d\n", rxn.ReactantsCount(), rxn.ProductsCount())
	fmt.Printf("Reactant atoms:       %d\n", reactant.AtomCount())
	fmt.Printf("Product atoms:        %d\n", product.AtomCount())
	fmt.Printf("KET root nodes:       %d\n", ketRootNodeCount(ket))
}

func mustParse(smiles string) *molecule.Molecule {
	mol, err := (molecule.SmilesLoader{}).Parse(smiles)
	if err != nil {
		log.Fatalf("parse %q: %v", smiles, err)
	}
	return mol
}

func ketRootNodeCount(input string) int {
	var doc struct {
		Root struct {
			Nodes []struct{} `json:"nodes"`
		} `json:"root"`
	}
	if err := json.Unmarshal([]byte(input), &doc); err != nil {
		log.Fatalf("decode KET: %v", err)
	}
	return len(doc.Root.Nodes)
}
