package main

import (
	"encoding/json"
	"fmt"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/reaction"
)

func Example_moleculeKET() {
	mol := molecule.NewMolecule()
	c := mol.AddAtom(molecule.ELEM_C)
	o := mol.AddAtom(molecule.ELEM_O)
	mol.AddBond(c, o, molecule.BOND_DOUBLE)
	mol.SetAtomXY(c, 0, 0)
	mol.SetAtomXY(o, 1.25, 0)

	ket, err := molecule.SaveMoleculeToKET(mol)
	if err != nil {
		panic(err)
	}

	var doc struct {
		Mol0 struct {
			Atoms []struct {
				Label string `json:"label"`
			} `json:"atoms"`
			Bonds []struct {
				Type  int   `json:"type"`
				Atoms []int `json:"atoms"`
			} `json:"bonds"`
		} `json:"mol0"`
	}
	if err := json.Unmarshal([]byte(ket), &doc); err != nil {
		panic(err)
	}

	fmt.Println(doc.Mol0.Atoms[0].Label, doc.Mol0.Atoms[1].Label)
	fmt.Println(doc.Mol0.Bonds[0].Type, doc.Mol0.Bonds[0].Atoms)

	// Output:
	// C O
	// 2 [0 1]
}

func Example_reactionKET() {
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

	ket, err := reaction.SaveReactionToKETString(rxn, false)
	if err != nil {
		panic(err)
	}

	var doc struct {
		Root struct {
			Nodes []struct {
				Ref  string `json:"$ref"`
				Type string `json:"type"`
			} `json:"nodes"`
		} `json:"root"`
		Mol0 struct {
			Atoms []struct {
				Label string `json:"label"`
			} `json:"atoms"`
		} `json:"mol0"`
		Mol1 struct {
			Atoms []struct {
				Label string `json:"label"`
			} `json:"atoms"`
			Bonds []struct {
				Type int `json:"type"`
			} `json:"bonds"`
		} `json:"mol1"`
	}
	if err := json.Unmarshal([]byte(ket), &doc); err != nil {
		panic(err)
	}

	fmt.Println(len(doc.Root.Nodes))
	fmt.Println(doc.Root.Nodes[0].Ref, doc.Root.Nodes[1].Ref, doc.Root.Nodes[2].Type)
	fmt.Println(doc.Mol0.Atoms[0].Label, len(doc.Mol1.Atoms), doc.Mol1.Bonds[0].Type)

	// Output:
	// 3
	// mol0 mol1 arrow
	// Cl 2 1
}
