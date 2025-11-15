// Package reaction_test provides tests for reaction helper functions
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_helpers_test.go
// @Software: GoLand
package reaction_test

import (
	"testing"
)

func TestGetReactant(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Get first reactant (ethanol)
	reactantHandle, err := rxn.GetReactant(0)
	if err != nil {
		t.Fatalf("Failed to get reactant: %v", err)
	}

	if reactantHandle < 0 {
		t.Error("Invalid reactant handle")
	}
}

func TestGetProduct(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Get first product
	productHandle, err := rxn.GetProduct(0)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if productHandle < 0 {
		t.Error("Invalid product handle")
	}
}

func TestGetReactantOutOfRange(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Try to get non-existent reactant
	_, err = rxn.GetReactant(100)
	if err == nil {
		t.Error("Expected error for out-of-range reactant index")
	}
}

func TestIteratorHelpers(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("CCO.CC>>CC=O.C")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Test iterator creation
	reactantCount, _ := rxn.CountReactants()
	if reactantCount != 2 {
		t.Errorf("Expected 2 reactants, got %d", reactantCount)
	}

	productCount, _ := rxn.CountProducts()
	if productCount != 2 {
		t.Errorf("Expected 2 products, got %d", productCount)
	}
}

func TestReactionLayout(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Test layout
	err = rxn.Layout()
	if err != nil {
		t.Errorf("Layout failed: %v", err)
	}
}

func TestReactionClean2D(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Test clean2d
	err = rxn.Clean2D()
	if err != nil {
		t.Errorf("Clean2D failed: %v", err)
	}
}

func TestReactionAromatize(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("C1=CC=CC=C1>>c1ccccc1")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Test aromatize
	err = rxn.Aromatize()
	if err != nil {
		t.Errorf("Aromatize failed: %v", err)
	}
}

func TestReactionDearomatize(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("c1ccccc1>>C1=CC=CC=C1")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Test dearomatize
	err = rxn.Dearomatize()
	if err != nil {
		t.Errorf("Dearomatize failed: %v", err)
	}
}

func TestGetCatalyst(t *testing.T) {
	// Most reactions don't have catalysts, so this is a basic test
	rxn, err := indigoInit.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	catalystCount, _ := rxn.CountCatalysts()
	if catalystCount > 0 {
		_, err := rxn.GetCatalyst(0)
		if err != nil {
			t.Errorf("Failed to get catalyst: %v", err)
		}
	}
}

func TestReactionCountMethods(t *testing.T) {
	rxn, err := indigoInit.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
	if err != nil {
		t.Fatalf("Failed to load reaction: %v", err)
	}
	defer rxn.Close()

	reactants, err := rxn.CountReactants()
	if err != nil {
		t.Errorf("CountReactants failed: %v", err)
	}
	if reactants != 2 {
		t.Errorf("Expected 2 reactants, got %d", reactants)
	}

	products, err := rxn.CountProducts()
	if err != nil {
		t.Errorf("CountProducts failed: %v", err)
	}
	if products != 2 {
		t.Errorf("Expected 2 products, got %d", products)
	}

	totalMolecules, err := rxn.CountMolecules()
	if err != nil {
		t.Errorf("CountMolecules failed: %v", err)
	}
	if totalMolecules != reactants+products {
		t.Errorf("Expected %d total molecules, got %d", reactants+products, totalMolecules)
	}
}
