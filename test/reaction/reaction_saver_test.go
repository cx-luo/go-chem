package reaction_test

import (
	"os"
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/reaction"
)

// TestReactionToRxnfile tests converting a reaction to RXN format
func TestReactionToRxnfile(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Convert to RXN
	rxnString, err := r.ToRxnfile()
	if err != nil {
		t.Fatalf("failed to convert to RXN: %v", err)
	}

	// Verify it's not empty
	if rxnString == "" {
		t.Error("RXN string is empty")
	}

	// Verify it starts with $RXN header
	if !strings.HasPrefix(rxnString, "$RXN") {
		t.Error("RXN string does not start with $RXN header")
	}
}

// TestReactionSaveToFile tests saving a reaction to a file
func TestReactionSaveToFile(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_reaction_*.rxn")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Save to file
	err = r.SaveToFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to save to file: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("saved file does not exist: %v", err)
	}

	if info.Size() == 0 {
		t.Error("saved file is empty")
	}

	// Try to load the saved file
	r2, err := reaction.LoadReactionFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load saved reaction: %v", err)
	}
	defer r2.Close()

	// Verify counts match
	count1, _ := r.CountReactants()
	count2, _ := r2.CountReactants()
	if count1 != count2 {
		t.Errorf("reactant count mismatch: original=%d, loaded=%d", count1, count2)
	}

	count1, _ = r.CountProducts()
	count2, _ = r2.CountProducts()
	if count1 != count2 {
		t.Errorf("product count mismatch: original=%d, loaded=%d", count1, count2)
	}
}

// TestReactionToSmiles tests converting a reaction to SMILES
func TestReactionToSmiles(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Convert to SMILES
	smiles, err := r.ToSmiles()
	if err != nil {
		t.Fatalf("failed to convert to SMILES: %v", err)
	}

	// Verify it's not empty
	if smiles == "" {
		t.Error("SMILES string is empty")
	}

	// Verify it contains reaction arrow
	if !strings.Contains(smiles, ">>") {
		t.Error("SMILES string does not contain reaction arrow >>")
	}
}

// TestReactionToCanonicalSmiles tests converting a reaction to canonical SMILES
func TestReactionToCanonicalSmiles(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Convert to canonical SMILES
	smiles, err := r.ToCanonicalSmiles()
	if err != nil {
		t.Fatalf("failed to convert to canonical SMILES: %v", err)
	}

	// Verify it's not empty
	if smiles == "" {
		t.Error("canonical SMILES string is empty")
	}

	// Verify it contains reaction arrow
	if !strings.Contains(smiles, ">>") {
		t.Error("canonical SMILES string does not contain reaction arrow >>")
	}
}

// TestReactionSaveLoadRoundtrip tests saving and loading a reaction
func TestReactionSaveLoadRoundtrip(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Get original counts
	origReactants, _ := r.CountReactants()
	origProducts, _ := r.CountProducts()

	// Save to file
	tmpFile, err := os.CreateTemp("", "test_reaction_*.rxn")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = r.SaveToFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to save to file: %v", err)
	}

	// Load from file
	r2, err := reaction.LoadReactionFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load from file: %v", err)
	}
	defer r2.Close()

	// Verify counts
	loadedReactants, _ := r2.CountReactants()
	loadedProducts, _ := r2.CountProducts()

	if origReactants != loadedReactants {
		t.Errorf("reactant count mismatch: original=%d, loaded=%d", origReactants, loadedReactants)
	}

	if origProducts != loadedProducts {
		t.Errorf("product count mismatch: original=%d, loaded=%d", origProducts, loadedProducts)
	}
}

// TestReactionSaveComplexReaction tests saving a more complex reaction
func TestReactionSaveComplexReaction(t *testing.T) {
	// Load a more complex reaction with 3 reactants and 2 products
	rxn := "CC(=O)O.CCO.C>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Convert to RXN
	rxnString, err := r.ToRxnfile()
	if err != nil {
		t.Fatalf("failed to convert to RXN: %v", err)
	}

	// Verify structure
	if !strings.Contains(rxnString, "$RXN") {
		t.Error("RXN string missing $RXN header")
	}

	if !strings.Contains(rxnString, "$MOL") {
		t.Error("RXN string missing $MOL sections")
	}
}

// TestReactionSaveClosedReaction tests saving a closed reaction
func TestReactionSaveClosedReaction(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}

	// Close the reaction
	r.Close()

	// Try to save (should error)
	_, err = r.ToRxnfile()
	if err == nil {
		t.Error("expected error when saving closed reaction")
	}

	err = r.SaveToFile("test.rxn")
	if err == nil {
		t.Error("expected error when saving closed reaction to file")
	}
}

// TestReactionSaveWithMapping tests saving a reaction with atom mapping
func TestReactionSaveWithMapping(t *testing.T) {
	// Load a reaction with explicit mapping
	rxn := "[CH3:1][C:2](=[O:3])[OH:4].[CH3:5][CH2:6][OH:7]>>[CH3:1][C:2](=[O:3])[O:7][CH2:6][CH3:5].[OH2:4]"

	r, err := reaction.LoadReactionSmartsFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Save to file
	tmpFile, err := os.CreateTemp("", "test_reaction_*.rxn")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = r.SaveToFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to save to file: %v", err)
	}

	// Verify file exists
	_, err = os.Stat(tmpFile.Name())
	if err != nil {
		t.Errorf("saved file does not exist: %v", err)
	}
}

// TestReactionSmilesRoundtrip tests SMILES conversion roundtrip
func TestReactionSmilesRoundtrip(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Convert to SMILES
	smiles, err := r.ToSmiles()
	if err != nil {
		t.Fatalf("failed to convert to SMILES: %v", err)
	}

	// Load from SMILES
	r2, err := reaction.LoadReactionFromString(smiles)
	if err != nil {
		t.Fatalf("failed to load from SMILES: %v", err)
	}
	defer r2.Close()

	// Verify counts
	count1, _ := r.CountReactants()
	count2, _ := r2.CountReactants()
	if count1 != count2 {
		t.Errorf("reactant count mismatch after SMILES roundtrip: original=%d, reloaded=%d", count1, count2)
	}

	count1, _ = r.CountProducts()
	count2, _ = r2.CountProducts()
	if count1 != count2 {
		t.Errorf("product count mismatch after SMILES roundtrip: original=%d, reloaded=%d", count1, count2)
	}
}
