package reaction_test

import (
	"testing"

	"github.com/cx-luo/go-chem/core"
)

var indigoInit *core.Indigo
var indigoInchi *core.IndigoInchi

func init() {
	handle, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}
	indigoInit = handle

	indigoInchiHandle, err := indigoInit.InchiInit()
	if err != nil {
		panic(err)
	}
	indigoInchi = indigoInchiHandle
}

// TestCreateReaction tests creating a new reaction
func TestCreateReaction(t *testing.T) {
	r, err := indigoInit.CreateReaction()
	if err != nil {
		t.Fatalf("failed to create reaction: %v", err)
	}
	defer r.Close()

	// Test that the reaction is valid
	if r.Handle < 0 {
		t.Errorf("invalid reaction handle: %d", r.Handle)
	}

	// Test counts
	reactantCount, err := r.CountReactants()
	if err != nil {
		t.Errorf("failed to count reactants: %v", err)
	}
	if reactantCount != 0 {
		t.Errorf("expected 0 reactants, got %d", reactantCount)
	}

	productCount, err := r.CountProducts()
	if err != nil {
		t.Errorf("failed to count products: %v", err)
	}
	if productCount != 0 {
		t.Errorf("expected 0 products, got %d", productCount)
	}

	catalystCount, err := r.CountCatalysts()
	if err != nil {
		t.Errorf("failed to count catalysts: %v", err)
	}
	if catalystCount != 0 {
		t.Errorf("expected 0 catalysts, got %d", catalystCount)
	}
}

// TestCreateQueryReaction tests creating a query reaction
func TestCreateQueryReaction(t *testing.T) {
	r, err := indigoInit.CreateQueryReaction()
	if err != nil {
		t.Fatalf("failed to create query reaction: %v", err)
	}
	defer r.Close()

	if r.Handle < 0 {
		t.Errorf("invalid query reaction handle: %d", r.Handle)
	}
}

// TestReactionClose tests closing a reaction
func TestReactionClose(t *testing.T) {
	r, err := indigoInit.CreateReaction()
	if err != nil {
		t.Fatalf("failed to create reaction: %v", err)
	}

	// Close the reaction
	err = r.Close()
	if err != nil {
		t.Errorf("failed to close reaction: %v", err)
	}

	// Closing again should not error
	err = r.Close()
	if err != nil {
		t.Errorf("second close should not error: %v", err)
	}

	// Operations on closed reaction should error
	_, err = r.CountReactants()
	if err == nil {
		t.Error("expected error when counting reactants on closed reaction")
	}
}

// TestReactionClone tests cloning a reaction
func TestReactionClone(t *testing.T) {
	// Load a simple reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Clone the reaction
	r2, err := r.Clone()
	if err != nil {
		t.Fatalf("failed to clone reaction: %v", err)
	}
	defer r2.Close()

	// Verify the clone has the same structure
	count1, _ := r.CountReactants()
	count2, _ := r2.CountReactants()
	if count1 != count2 {
		t.Errorf("clone has different reactant count: original=%d, clone=%d", count1, count2)
	}

	count1, _ = r.CountProducts()
	count2, _ = r2.CountProducts()
	if count1 != count2 {
		t.Errorf("clone has different product count: original=%d, clone=%d", count1, count2)
	}
}

// TestReactionOptimize tests optimizing a query reaction
func TestReactionOptimize(t *testing.T) {
	r, err := indigoInit.CreateQueryReaction()
	if err != nil {
		t.Fatalf("failed to create query reaction: %v", err)
	}
	defer r.Close()

	// Optimize with default options
	err = r.Optimize("")
	if err != nil {
		t.Errorf("failed to optimize reaction: %v", err)
	}
}

// TestReactionCountMolecules tests counting total molecules in a reaction
func TestReactionCountMolecules(t *testing.T) {
	// Load a reaction: 2 reactants + 2 products = 4 molecules
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	count, err := r.CountMolecules()
	if err != nil {
		t.Errorf("failed to count molecules: %v", err)
	}

	// Should have 4 molecules total
	if count != 4 {
		t.Errorf("expected 4 molecules, got %d", count)
	}
}

// TestReactionGetMolecule tests getting molecules from a reaction
func TestReactionGetMolecule(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	count, _ := r.CountMolecules()

	// Get each molecule
	for i := 0; i < count; i++ {
		molHandle, err := r.GetMolecule(i)
		if err != nil {
			t.Errorf("failed to get molecule %d: %v", i, err)
		}
		if molHandle < 0 {
			t.Errorf("invalid molecule handle at index %d: %d", i, molHandle)
		}
	}

	// Try to get molecule beyond bounds
	_, err = r.GetMolecule(count)
	if err == nil {
		t.Error("expected error when getting molecule beyond bounds")
	}
}
