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

	"github.com/cx-luo/go-chem/reaction"
)

func main() {
	fmt.Println("=== Accessing Molecules from Reactions ===\n")

	// Load a reaction: Fischer esterification
	// ethanol + acetic acid → ethyl acetate + water
	rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		log.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	fmt.Println("Reaction: CCO.CC(=O)O>>CC(=O)OCC.O")
	fmt.Println("(Ethanol + Acetic acid → Ethyl acetate + Water)\n")

	// Example 1: Get individual molecules by index
	fmt.Println("1. Getting Individual Molecules by Index:")

	// Get first reactant (ethanol)
	reactant0, err := rxn.GetReactantMolecule(0)
	if err != nil {
		log.Printf("Failed to get reactant 0: %v", err)
	} else {
		defer reactant0.Close()
		smiles, _ := reactant0.ToSmiles()
		formula, _ := reactant0.GrossFormula()
		mass, _ := reactant0.MolecularWeight()
		fmt.Printf("  Reactant 0: %s (%s, MW=%.2f)\n", smiles, formula, mass)
	}

	// Get second reactant (acetic acid)
	reactant1, err := rxn.GetReactantMolecule(1)
	if err != nil {
		log.Printf("Failed to get reactant 1: %v", err)
	} else {
		defer reactant1.Close()
		smiles, _ := reactant1.ToSmiles()
		formula, _ := reactant1.GrossFormula()
		mass, _ := reactant1.MolecularWeight()
		fmt.Printf("  Reactant 1: %s (%s, MW=%.2f)\n", smiles, formula, mass)
	}

	// Get first product (ethyl acetate)
	product0, err := rxn.GetProductMolecule(0)
	if err != nil {
		log.Printf("Failed to get product 0: %v", err)
	} else {
		defer product0.Close()
		smiles, _ := product0.ToSmiles()
		formula, _ := product0.GrossFormula()
		mass, _ := product0.MolecularWeight()
		fmt.Printf("  Product 0:  %s (%s, MW=%.2f)\n", smiles, formula, mass)
	}

	// Get second product (water)
	product1, err := rxn.GetProductMolecule(1)
	if err != nil {
		log.Printf("Failed to get product 1: %v", err)
	} else {
		defer product1.Close()
		smiles, _ := product1.ToSmiles()
		formula, _ := product1.GrossFormula()
		mass, _ := product1.MolecularWeight()
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
			smiles, _ := mol.ToSmiles()
			formula, _ := mol.GrossFormula()
			atomCount, _ := mol.CountAtoms()
			fmt.Printf("    [%d] %s (%s, %d atoms)\n", i, smiles, formula, atomCount)
			mol.Close() // Don't forget to close!
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
			smiles, _ := mol.ToSmiles()
			formula, _ := mol.GrossFormula()
			atomCount, _ := mol.CountAtoms()
			fmt.Printf("    [%d] %s (%s, %d atoms)\n", i, smiles, formula, atomCount)
			mol.Close()
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
		mass, _ := mol.MolecularWeight()
		atoms, _ := mol.CountAtoms()
		totalReactantMass += float32(mass)
		totalReactantAtoms += atoms
		smiles, _ := mol.ToSmiles()
		fmt.Printf("    %d. %s (MW=%.2f, atoms=%d)\n", i+1, smiles, mass, atoms)
		mol.Close()
	}

	fmt.Println("  Products:")
	for i, mol := range products {
		mass, _ := mol.MolecularWeight()
		atoms, _ := mol.CountAtoms()
		totalProductMass += float32(mass)
		totalProductAtoms += atoms
		smiles, _ := mol.ToSmiles()
		fmt.Printf("    %d. %s (MW=%.2f, atoms=%d)\n", i+1, smiles, mass, atoms)
		mol.Close()
	}

	fmt.Printf("\n  Mass balance: %.2f → %.2f (diff=%.4f)\n",
		totalReactantMass, totalProductMass, totalReactantMass-totalProductMass)
	fmt.Printf("  Atom count: %d → %d (diff=%d)\n",
		totalReactantAtoms, totalProductAtoms, totalReactantAtoms-totalProductAtoms)

	// Example 5: Process molecules from a complex reaction
	fmt.Println("\n5. Processing Complex Reaction:")

	rxn2, err := reaction.LoadReactionFromString("c1ccccc1.Br2>>c1ccc(Br)cc1.HBr")
	if err != nil {
		log.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn2.Close()

	fmt.Println("  Reaction: Bromination of benzene")
	fmt.Println("  c1ccccc1.Br2>>c1ccc(Br)cc1.HBr")

	// Aromatize reactants
	reactants2, _ := rxn2.GetAllReactants()
	fmt.Println("\n  Processing reactants:")
	for i, mol := range reactants2 {
		beforeSmiles, _ := mol.ToSmiles()
		mol.Aromatize()
		afterSmiles, _ := mol.ToCanonicalSmiles()
		rings, _ := mol.CountSSSR()
		fmt.Printf("    %d. %s → %s (rings=%d)\n", i+1, beforeSmiles, afterSmiles, rings)
		mol.Close()
	}

	// Process products
	products2, _ := rxn2.GetAllProducts()
	fmt.Println("\n  Processing products:")
	for i, mol := range products2 {
		smiles, _ := mol.ToCanonicalSmiles()
		heavy, _ := mol.CountHeavyAtoms()
		bonds, _ := mol.CountBonds()
		fmt.Printf("    %d. %s (heavy atoms=%d, bonds=%d)\n", i+1, smiles, heavy, bonds)
		mol.Close()
	}

	fmt.Println("\n=== Examples completed successfully ===")
}
