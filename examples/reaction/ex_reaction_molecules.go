// Package main demonstrates how to access individual molecules from reactions
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_molecules.go
// @Software: GoLand
package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/core"
)

func main() {
	fmt.Println("=== Accessing Molecules from Reactions ===")

	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	indigoInchi, err := core.InchiInit(indigoInit.GetSessionID())
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatalf("Failed to initialize Indigo: %v", err)
	}
	defer core.DisposeInChI(indigoInchi.GetInchiSessionID())

	// Load a reaction: Fischer esterification
	// ethanol + acetic acid → ethyl acetate + water
	rxn, err := indigoInit.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		log.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	fmt.Println("Reaction: CCO.CC(=O)O>>CC(=O)OCC.O")
	fmt.Println("(Ethanol + Acetic acid → Ethyl acetate + Water)")

	// Example 1: Get individual molecules by index
	fmt.Println("1. Getting Individual Molecules by Index:")

	// Get first reactant (ethanol)
	reactant0Handle, err := rxn.GetReactantMolecule(0)
	if err != nil {
		log.Printf("Failed to get reactant 0: %v", err)
	} else {
		reactant0, err := indigoInit.LoadMoleculeFromHandle(reactant0Handle)
		if err != nil {
			log.Printf("Failed to get reactant 0: %v", err)
		}
		defer reactant0.Close()
		smiles, _ := reactant0.ToSmiles()
		formula, _ := reactant0.GrossFormula()
		mass, _ := reactant0.MolecularWeight()
		fmt.Printf("  Reactant 0: %s (%s, MW=%.2f)\n", smiles, formula, mass)
	}

	// Get second reactant (acetic acid)
	reactant1Handle, err := rxn.GetReactantMolecule(1)
	if err != nil {
		log.Printf("Failed to get reactant 1 handle: %v", err)
	} else {
		reactant1, err := indigoInit.LoadMoleculeFromHandle(reactant1Handle)
		if err != nil {
			log.Printf("Failed to get reactant 1: %v", err)
		}
		defer reactant1.Close()
		smiles, _ := reactant1.ToSmiles()
		formula, _ := reactant1.GrossFormula()
		mass, _ := reactant1.MolecularWeight()
		inChI, err := indigoInchi.GenerateInChI(reactant1)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  Reactant 1: %s (%s, MW=%.2f) inchi: %s \n", smiles, formula, mass, inChI)
	}

	// Get first product (ethyl acetate)
	product0, err := rxn.GetProductMolecule(0)
	if err != nil {
		log.Printf("Failed to get product 0: %v", err)
	} else {
		product0Molecule, err := indigoInit.LoadMoleculeFromHandle(product0)
		if err != nil {
			log.Printf("Failed to get product 0: %v", err)
		}
		defer product0Molecule.Close()
		smiles, _ := product0Molecule.ToSmiles()
		formula, _ := product0Molecule.GrossFormula()
		mass, _ := product0Molecule.MolecularWeight()
		fmt.Printf("  Product 0:  %s (%s, MW=%.2f)\n", smiles, formula, mass)
	}

	// Get second product (water)
	product1, err := rxn.GetProductMolecule(1)
	if err != nil {
		log.Printf("Failed to get product 1: %v", err)
	} else {
		product1Molecule, err := indigoInit.LoadMoleculeFromHandle(product1)
		if err != nil {
			log.Printf("Failed to get product 1: %v", err)
		}
		defer product1Molecule.Close()
		smiles, _ := product1Molecule.ToSmiles()
		formula, _ := product1Molecule.GrossFormula()
		mass, _ := product1Molecule.MolecularWeight()
		fmt.Printf("  Product 1:  %s (%s, MW=%.2f)\n", smiles, formula, mass)
	}

	// Example 2: Get all reactants at once
	fmt.Println("\n2. Getting All Reactants at Once:")

	allReactants, err := rxn.GetAllReactants()
	if err != nil {
		log.Printf("Failed to get all reactants: %v", err)
	} else {
		fmt.Printf("  Found %d reactants:\n", len(allReactants))
		for i, mol := range allReactants {
			molMolecule, err := indigoInit.LoadMoleculeFromHandle(mol)
			if err != nil {
				log.Printf("Failed to get molecule %d: %v", i, err)
			}
			defer molMolecule.Close()
		}
	}

	// Example 3: Get all products at once
	fmt.Println("\n3. Getting All Products at Once:")

	allProducts, err := rxn.GetAllProducts()
	if err != nil {
		log.Printf("Failed to get all products: %v", err)
	} else {
		fmt.Printf("  Found %d products:\n", len(allProducts))
		for i, mol := range allProducts {
			molMolecule, err := indigoInit.LoadMoleculeFromHandle(mol)
			if err != nil {
				log.Printf("Failed to get molecule %d: %v", i, err)
			}
			defer molMolecule.Close()
			smiles, _ := molMolecule.ToSmiles()
			formula, _ := molMolecule.GrossFormula()
			atomCount, _ := molMolecule.CountAtoms()
			fmt.Printf("    [%d] %s (%s, %d atoms)\n", i, smiles, formula, atomCount)
			molMolecule.Close()
		}
	}

	// Example 4: Analyze molecules
	fmt.Println("\n4. Analyzing Reaction Molecules:")

	reactants, _ := rxn.GetAllReactants()
	products, _ := rxn.GetAllProducts()

	var totalReactantMass, totalProductMass float32
	var totalReactantAtoms, totalProductAtoms int

	fmt.Println("  Reactants:")
	for i, mol := range reactants {
		molMolecule, err := indigoInit.LoadMoleculeFromHandle(mol)
		if err != nil {
			log.Printf("Failed to get molecule %d: %v", i, err)
		}
		defer molMolecule.Close()
		mass, _ := molMolecule.MolecularWeight()
		atoms, _ := molMolecule.CountAtoms()
		totalReactantMass += float32(mass)
		totalReactantAtoms += atoms
		smiles, _ := molMolecule.ToSmiles()
		fmt.Printf("    %d. %s (MW=%.2f, atoms=%d)\n", i+1, smiles, mass, atoms)
		molMolecule.Close()
	}

	fmt.Println("  Products:")
	for i, mol := range products {
		molMolecule, err := indigoInit.LoadMoleculeFromHandle(mol)
		if err != nil {
			log.Printf("Failed to get molecule %d: %v", i, err)
		}
		defer molMolecule.Close()
		mass, _ := molMolecule.MolecularWeight()
		atoms, _ := molMolecule.CountAtoms()
		totalProductMass += float32(mass)
		totalProductAtoms += atoms
		smiles, _ := molMolecule.ToSmiles()
		fmt.Printf("    %d. %s (MW=%.2f, atoms=%d)\n", i+1, smiles, mass, atoms)
		molMolecule.Close()
	}

	fmt.Printf("\n  Mass balance: %.2f → %.2f (diff=%.4f)\n",
		totalReactantMass, totalProductMass, totalReactantMass-totalProductMass)
	fmt.Printf("  Atom count: %d → %d (diff=%d)\n",
		totalReactantAtoms, totalProductAtoms, totalReactantAtoms-totalProductAtoms)

	// Example 5: Process molecules from a complex reaction
	fmt.Println("\n5. Processing Complex Reaction:")

	rxn2, err := indigoInit.LoadReactionFromString("COCC(=O)N1CC(C)(N)C1.COC1=CC=C(Cl)C=C1NC(=O)CN>>OC(=O)C(F)(F)F.O=C1N(C2CCC(=O)NC2=O)C(=O)C2=CC3=C(CNC3)C=C12 |f:2.3,c:18,t:13,15,45,47,53,lp:25:2,27:2,29:3,30:3,31:3,32:2,34:1,39:2,40:1,42:2,44:2,50:1|")
	if err != nil {
		log.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn2.Close()

	// Aromatize reactants
	reactants2, _ := rxn2.GetAllReactants()
	fmt.Println("\n  Processing reactants:")
	for i, mol := range reactants2 {
		molMolecule, err := indigoInit.LoadMoleculeFromHandle(mol)
		if err != nil {
			log.Printf("Failed to get molecule %d: %v", i, err)
		}
		defer molMolecule.Close()
		beforeSmiles, _ := molMolecule.ToSmiles()
		molMolecule.Aromatize()
		afterSmiles, _ := molMolecule.ToCanonicalSmiles()
		rings, _ := molMolecule.CountSSSR()

		inchi, err := indigoInchi.GenerateInChI(molMolecule)
		if err != nil {
			panic(err)
		}

		inChIKey, err := indigoInchi.InChIToKey(inchi)
		if err != nil {
			panic(err)
		}

		fmt.Printf("    %d. %s → %s (rings=%d)\n"+
			"	inchi: %s\n"+
			"	inchikey: %s\n", i+1, beforeSmiles, afterSmiles, rings, inchi, inChIKey)
		molMolecule.Close()
	}

	// Process products
	products2, _ := rxn2.GetAllProducts()
	fmt.Println("\n  Processing products:")
	for i, mol := range products2 {
		molMolecule, err := indigoInit.LoadMoleculeFromHandle(mol)
		if err != nil {
			log.Printf("Failed to get molecule %d: %v", i, err)
		}
		defer molMolecule.Close()
		smiles, _ := molMolecule.ToCanonicalSmiles()
		heavy, _ := molMolecule.CountHeavyAtoms()
		bonds, _ := molMolecule.CountBonds()
		fmt.Printf("    %d. %s (heavy atoms=%d, bonds=%d)\n", i+1, smiles, heavy, bonds)
		molMolecule.Close()
	}

	fmt.Println("\n=== Examples completed successfully ===")
}
