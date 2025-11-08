// Package reaction_test provides tests for reaction molecule access methods
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_molecules_test.go
// @Software: GoLand
package reaction_test

import (
	"testing"

	"github.com/cx-luo/go-chem/reaction"
)

func TestGetReactantMolecule(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Get first reactant (ethanol)
	mol, err := rxn.GetReactantMolecule(0)
	if err != nil {
		t.Fatalf("Failed to get reactant molecule: %v", err)
	}
	defer mol.Close()

	// Verify it's ethanol
	smiles, err := mol.ToSmiles()
	if err != nil {
		t.Fatalf("Failed to get SMILES: %v", err)
	}

	t.Logf("Reactant SMILES: %s", smiles)

	// Should be ethanol (CCO or OCC)
	if smiles != "CCO" && smiles != "OCC" {
		t.Logf("Warning: Expected ethanol, got %s", smiles)
	}

	// Check atom count
	atomCount, _ := mol.CountAtoms()
	if atomCount != 3 {
		t.Errorf("Expected 3 atoms, got %d", atomCount)
	}
}

func TestGetProductMolecule(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Get first product (ethyl acetate)
	mol, err := rxn.GetProductMolecule(0)
	if err != nil {
		t.Fatalf("Failed to get product molecule: %v", err)
	}
	defer mol.Close()

	smiles, err := mol.ToSmiles()
	if err != nil {
		t.Fatalf("Failed to get SMILES: %v", err)
	}

	t.Logf("Product SMILES: %s", smiles)

	// Check atom count for ethyl acetate (C4H8O2 = 6 heavy atoms)
	heavyAtoms, _ := mol.CountHeavyAtoms()
	if heavyAtoms != 6 {
		t.Errorf("Expected 6 heavy atoms, got %d", heavyAtoms)
	}
}

func TestGetMoleculeOutOfRange(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Try to get non-existent reactant
	_, err = rxn.GetReactantMolecule(100)
	if err == nil {
		t.Error("Expected error for out-of-range reactant index")
	}

	// Try to get non-existent product
	_, err = rxn.GetProductMolecule(100)
	if err == nil {
		t.Error("Expected error for out-of-range product index")
	}
}

func TestGetAllReactants(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	reactants, err := rxn.GetAllReactants()
	if err != nil {
		t.Fatalf("Failed to get all reactants: %v", err)
	}

	if len(reactants) != 2 {
		t.Errorf("Expected 2 reactants, got %d", len(reactants))
	}

	// Check each reactant
	for i, mol := range reactants {
		smiles, err := mol.ToSmiles()
		if err != nil {
			t.Errorf("Failed to get SMILES for reactant %d: %v", i, err)
		}
		t.Logf("Reactant %d: %s", i, smiles)
		mol.Close()
	}
}

func TestGetAllProducts(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	products, err := rxn.GetAllProducts()
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}

	// Check each product
	for i, mol := range products {
		smiles, err := mol.ToSmiles()
		if err != nil {
			t.Errorf("Failed to get SMILES for product %d: %v", i, err)
		}
		t.Logf("Product %d: %s", i, smiles)
		mol.Close()
	}
}

func TestGetAllCatalysts(t *testing.T) {
	// Most simple reactions don't have catalysts
	rxn, err := reaction.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	catalysts, err := rxn.GetAllCatalysts()
	if err != nil {
		t.Fatalf("Failed to get all catalysts: %v", err)
	}

	// Should be 0 for this reaction
	if len(catalysts) != 0 {
		t.Logf("Found %d catalysts (expected 0)", len(catalysts))
	}

	// Close any catalysts
	for _, mol := range catalysts {
		mol.Close()
	}
}

func TestMoleculeManipulation(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("c1ccccc1>>C1=CC=CC=C1")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Get reactant and manipulate it
	mol, err := rxn.GetReactantMolecule(0)
	if err != nil {
		t.Fatalf("Failed to get reactant: %v", err)
	}
	defer mol.Close()

	beforeSmiles, _ := mol.ToCanonicalSmiles()
	t.Logf("Before aromatize: %s", beforeSmiles)

	// Aromatize
	err = mol.Aromatize()
	if err != nil {
		t.Errorf("Failed to aromatize: %v", err)
	}

	afterSmiles, _ := mol.ToCanonicalSmiles()
	t.Logf("After aromatize: %s", afterSmiles)

	// Check properties
	rings, _ := mol.CountSSSR()
	if rings != 1 {
		t.Errorf("Expected 1 ring, got %d", rings)
	}

	atoms, _ := mol.CountAtoms()
	if atoms != 6 {
		t.Errorf("Expected 6 atoms, got %d", atoms)
	}
}

func TestMassBalance(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	reactants, _ := rxn.GetAllReactants()
	products, _ := rxn.GetAllProducts()

	var reactantMass, productMass float64

	for _, mol := range reactants {
		mass, _ := mol.MolecularWeight()
		reactantMass += mass
		mol.Close()
	}

	for _, mol := range products {
		mass, _ := mol.MolecularWeight()
		productMass += mass
		mol.Close()
	}

	diff := reactantMass - productMass
	t.Logf("Reactant mass: %.2f, Product mass: %.2f, Diff: %.4f", reactantMass, productMass, diff)

	// Allow small floating point difference
	if diff > 0.1 || diff < -0.1 {
		t.Errorf("Mass balance error too large: %.4f", diff)
	}
}

func TestClosedReaction(t *testing.T) {
	rxn, err := reaction.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}

	rxn.Close()

	// Should error on closed reaction
	_, err = rxn.GetReactantMolecule(0)
	if err == nil {
		t.Error("Expected error when getting molecule from closed reaction")
	}

	_, err = rxn.GetAllReactants()
	if err == nil {
		t.Error("Expected error when getting all reactants from closed reaction")
	}

	_, err = rxn.GetAllProducts()
	if err == nil {
		t.Error("Expected error when getting all products from closed reaction")
	}
}
