package reaction_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cx-luo/go-chem/reaction"
)

// TestLoadReactionFromString tests loading a reaction from a SMILES string
func TestLoadReactionFromString(t *testing.T) {
	// Test simple esterification reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction from string: %v", err)
	}
	defer r.Close()

	canonicalSmiles, err := r.ToCanonicalSmiles()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(canonicalSmiles)

	// Verify reactants
	reactantCount, err := r.CountReactants()
	if err != nil {
		t.Errorf("failed to count reactants: %v", err)
	}
	if reactantCount != 2 {
		t.Errorf("expected 2 reactants, got %d", reactantCount)
	}

	// Verify products
	productCount, err := r.CountProducts()
	if err != nil {
		t.Errorf("failed to count products: %v", err)
	}
	if productCount != 2 {
		t.Errorf("expected 2 products, got %d", productCount)
	}
}

// TestLoadReactionFromStringWithCatalyst tests loading a reaction with catalyst
func TestLoadReactionFromStringWithCatalyst(t *testing.T) {
	// Reaction with catalyst notation
	rxn := "C=C.O>>[H]C([H])C([H])([H])O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction with catalyst: %v", err)
	}
	defer r.Close()

	// Verify the reaction loaded successfully
	if r.Handle() < 0 {
		t.Error("invalid reaction handle")
	}
}

// TestLoadReactionFromBuffer tests loading a reaction from a byte buffer
func TestLoadReactionFromBuffer(t *testing.T) {
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
	buffer := []byte(rxn)

	r, err := reaction.LoadReactionFromBuffer(buffer)
	if err != nil {
		t.Fatalf("failed to load reaction from buffer: %v", err)
	}
	defer r.Close()

	reactantCount, _ := r.CountReactants()
	if reactantCount != 2 {
		t.Errorf("expected 2 reactants, got %d", reactantCount)
	}
}

// TestLoadReactionFromFile tests loading a reaction from a file
func TestLoadReactionFromFile(t *testing.T) {
	// Create a temporary RXN file
	tmpFile, err := os.CreateTemp("", "test_reaction_*.rxn")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write a simple RXN file content
	rxnContent := `$RXN

      RDKit

  2  2
$MOL

  RDKit          2D

  4  3  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.2990    0.7500    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    2.5981    0.0000    0.0000 O   0  0  0  0  0  0  0  0  0  0  0  0
    1.2990    2.2500    0.0000 O   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0
  2  3  2  0
  2  4  1  0
M  END
$MOL

  RDKit          2D

  2  1  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.2990    0.7500    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0
M  END
$MOL

  RDKit          2D

  3  2  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    1.2990    0.7500    0.0000 C   0  0  0  0  0  0  0  0  0  0  0  0
    2.5981    0.0000    0.0000 O   0  0  0  0  0  0  0  0  0  0  0  0
  1  2  1  0
  2  3  1  0
M  END
$MOL

  RDKit          2D

  1  0  0  0  0  0  0  0  0  0999 V2000
    0.0000    0.0000    0.0000 O   0  0  0  0  0  0  0  0  0  0  0  0
M  END
`

	_, err = tmpFile.WriteString(rxnContent)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Load the reaction
	r, err := reaction.LoadReactionFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load reaction from file: %v", err)
	}
	defer r.Close()

	// Verify the reaction loaded
	reactantCount, _ := r.CountReactants()
	if reactantCount != 2 {
		t.Errorf("expected 2 reactants, got %d", reactantCount)
	}

	productCount, _ := r.CountProducts()
	if productCount != 2 {
		t.Errorf("expected 2 products, got %d", productCount)
	}
}

// TestLoadQueryReactionFromString tests loading a query reaction
func TestLoadQueryReactionFromString(t *testing.T) {
	// Simple query reaction with wildcards
	rxn := "[#6]C(=O)O.CCO>>[#6]C(=O)OCC.O"

	r, err := reaction.LoadQueryReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load query reaction: %v", err)
	}
	defer r.Close()

	if r.Handle() < 0 {
		t.Error("invalid query reaction handle")
	}
}

// TestLoadQueryReactionFromBuffer tests loading a query reaction from buffer
func TestLoadQueryReactionFromBuffer(t *testing.T) {
	rxn := "[#6]C(=O)O.CCO>>[#6]C(=O)OCC.O"
	buffer := []byte(rxn)

	r, err := reaction.LoadQueryReactionFromBuffer(buffer)
	if err != nil {
		t.Fatalf("failed to load query reaction from buffer: %v", err)
	}
	defer r.Close()

	if r.Handle() < 0 {
		t.Error("invalid query reaction handle")
	}
}

// TestLoadReactionSmartsFromString tests loading a reaction SMARTS
func TestLoadReactionSmartsFromString(t *testing.T) {
	// Simple SMARTS pattern for esterification
	smarts := "[C:1](=[O:2])[OH:3].[C:4][OH:5]>>[C:1](=[O:2])[O:5][C:4].[OH2:3]"

	r, err := reaction.LoadReactionSmartsFromString(smarts)
	if err != nil {
		t.Fatalf("failed to load reaction SMARTS: %v", err)
	}
	defer r.Close()

	if r.Handle() < 0 {
		t.Error("invalid reaction SMARTS handle")
	}
}

// TestLoadReactionSmartsFromBuffer tests loading reaction SMARTS from buffer
func TestLoadReactionSmartsFromBuffer(t *testing.T) {
	smarts := "[C:1](=[O:2])[OH:3].[C:4][OH:5]>>[C:1](=[O:2])[O:5][C:4].[OH2:3]"
	buffer := []byte(smarts)

	r, err := reaction.LoadReactionSmartsFromBuffer(buffer)
	if err != nil {
		t.Fatalf("failed to load reaction SMARTS from buffer: %v", err)
	}
	defer r.Close()

	if r.Handle() < 0 {
		t.Error("invalid reaction SMARTS handle")
	}
}

// TestLoadInvalidReaction tests loading invalid reaction data
func TestLoadInvalidReaction(t *testing.T) {
	// Test with invalid SMILES
	_, err := reaction.LoadReactionFromString("invalid>>reaction>>data")
	if err == nil {
		t.Error("expected error when loading invalid reaction")
	}

	// Test with empty string
	_, err = reaction.LoadReactionFromString("")
	if err == nil {
		t.Error("expected error when loading empty reaction")
	}

	// Test with empty buffer
	_, err = reaction.LoadReactionFromBuffer([]byte{})
	if err == nil {
		t.Error("expected error when loading empty buffer")
	}
}

// TestLoadReactionMultipleComponents tests loading reactions with multiple components
func TestLoadReactionMultipleComponents(t *testing.T) {
	// Reaction with 3 reactants and 2 products
	rxn := "CC.CC.CC>>CCCC.CC"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	reactantCount, _ := r.CountReactants()
	if reactantCount != 3 {
		t.Errorf("expected 3 reactants, got %d", reactantCount)
	}

	productCount, _ := r.CountProducts()
	if productCount != 2 {
		t.Errorf("expected 2 products, got %d", productCount)
	}
}

// TestReactionIterators tests iterating over reaction components
func TestReactionIterators(t *testing.T) {
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Test reactants iterator
	reactIter, err := r.IterateReactants()
	if err != nil {
		t.Fatalf("failed to create reactants iterator: %v", err)
	}
	defer reactIter.Close()

	reactantCount := 0
	for reactIter.HasNext() {
		molHandle, err := reactIter.Next()
		if err != nil {
			t.Errorf("failed to get next reactant: %v", err)
		}
		if molHandle < 0 {
			t.Errorf("invalid molecule handle: %d", molHandle)
		}
		reactantCount++
	}

	if reactantCount != 2 {
		t.Errorf("expected 2 reactants from iterator, got %d", reactantCount)
	}

	// Test products iterator
	prodIter, err := r.IterateProducts()
	if err != nil {
		t.Fatalf("failed to create products iterator: %v", err)
	}
	defer prodIter.Close()

	productCount := 0
	for prodIter.HasNext() {
		molHandle, err := prodIter.Next()
		if err != nil {
			t.Errorf("failed to get next product: %v", err)
		}
		if molHandle < 0 {
			t.Errorf("invalid molecule handle: %d", molHandle)
		}
		productCount++
	}

	if productCount != 2 {
		t.Errorf("expected 2 products from iterator, got %d", productCount)
	}

	// Test all molecules iterator
	molIter, err := r.IterateMolecules()
	if err != nil {
		t.Fatalf("failed to create molecules iterator: %v", err)
	}
	defer molIter.Close()

	moleculeCount := 0
	for molIter.HasNext() {
		molHandle, err := molIter.Next()
		if err != nil {
			t.Errorf("failed to get next molecule: %v", err)
		}
		if molHandle < 0 {
			t.Errorf("invalid molecule handle: %d", molHandle)
		}
		moleculeCount++
	}

	if moleculeCount != 4 {
		t.Errorf("expected 4 molecules from iterator, got %d", moleculeCount)
	}
}
