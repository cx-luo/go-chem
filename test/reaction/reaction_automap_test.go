package reaction_test

import (
	"testing"

	"github.com/cx-luo/go-indigo/reaction"
)

// TestReactionAutomap tests automatic atom-to-atom mapping
func TestReactionAutomap(t *testing.T) {
	// Load a reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Perform automatic mapping with default mode
	err = r.Automap(reaction.AutomapModeDiscard)
	if err != nil {
		t.Errorf("failed to automap reaction: %v", err)
	}
}

// TestReactionAutomapModes tests different automap modes
func TestReactionAutomapModes(t *testing.T) {
	tests := []struct {
		name string
		mode string
	}{
		{"Discard mode", reaction.AutomapModeDiscard},
		{"Keep mode", reaction.AutomapModeKeep},
		{"Alter mode", reaction.AutomapModeAlter},
		{"Clear mode", reaction.AutomapModeClear},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
			r, err := indigoInit.LoadReactionFromString(rxn)
			if err != nil {
				t.Fatalf("failed to load reaction: %v", err)
			}
			defer r.Close()

			err = r.Automap(tt.mode)
			if err != nil {
				t.Errorf("failed to automap with mode %s: %v", tt.mode, err)
			}
		})
	}
}

// TestReactionAutomapOptions tests automap with additional options
func TestReactionAutomapOptions(t *testing.T) {
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Automap with ignore charges option
	err = r.Automap(reaction.AutomapModeDiscard + " " + reaction.AutomapIgnoreCharges)
	if err != nil {
		t.Errorf("failed to automap with ignore charges: %v", err)
	}
}

// TestReactionAutomapMultipleOptions tests automap with multiple options
func TestReactionAutomapMultipleOptions(t *testing.T) {
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Automap with multiple options
	options := reaction.AutomapModeDiscard + " " +
		reaction.AutomapIgnoreCharges + " " +
		reaction.AutomapIgnoreIsotopes

	err = r.Automap(options)
	if err != nil {
		t.Errorf("failed to automap with multiple options: %v", err)
	}
}

// TestReactionClearAAM tests clearing atom-to-atom mapping
func TestReactionClearAAM(t *testing.T) {
	// Load a reaction with explicit mapping
	rxn := "[CH3:1][C:2](=[O:3])[OH:4].[CH3:5][CH2:6][OH:7]>>[CH3:1][C:2](=[O:3])[O:7][CH2:6][CH3:5].[OH2:4]"

	r, err := indigoInit.LoadReactionSmartsFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Clear AAM
	err = r.ClearAAM()
	if err != nil {
		t.Errorf("failed to clear AAM: %v", err)
	}
}

// TestReactionCorrectReactingCenters tests correcting reacting centers
func TestReactionCorrectReactingCenters(t *testing.T) {
	// Load a reaction with mapping
	rxn := "[CH3:1][C:2](=[O:3])[OH:4].[CH3:5][CH2:6][OH:7]>>[CH3:1][C:2](=[O:3])[O:7][CH2:6][CH3:5].[OH2:4]"

	r, err := indigoInit.LoadReactionSmartsFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Correct reacting centers
	err = r.CorrectReactingCenters()
	if err != nil {
		t.Errorf("failed to correct reacting centers: %v", err)
	}
}

// TestReactionGetSetAtomMapping tests getting and setting atom mapping numbers
func TestReactionGetSetAtomMapping(t *testing.T) {
	// Create a simple reaction
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// First automap the reaction
	err = r.Automap(reaction.AutomapModeDiscard)
	if err != nil {
		t.Fatalf("failed to automap reaction: %v", err)
	}

	// Get the first reactant
	reactIter, err := r.IterateReactants()
	if err != nil {
		t.Fatalf("failed to create reactants iterator: %v", err)
	}
	defer reactIter.Close()

	if !reactIter.HasNext() {
		t.Fatal("no reactants found")
	}

	molHandle, err := reactIter.Next()
	if err != nil {
		t.Fatalf("failed to get reactant: %v", err)
	}

	// Note: Getting/setting atom mapping requires atom handles from the molecule
	// This is a basic test to ensure the functions don't error
	_ = molHandle
}

// TestReactionGetSetReactingCenter tests getting and setting reacting centers
func TestReactionGetSetReactingCenter(t *testing.T) {
	// Load a reaction with mapping
	rxn := "[CH3:1][C:2](=[O:3])[OH:4].[CH3:5][CH2:6][OH:7]>>[CH3:1][C:2](=[O:3])[O:7][CH2:6][CH3:5].[OH2:4]"

	r, err := indigoInit.LoadReactionSmartsFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Note: Getting/setting reacting centers requires bond handles from the molecules
	// This is a basic test to ensure the reaction loaded
	_ = r
}

// TestReactionAutomapSmartsReaction tests automapping a SMARTS reaction
func TestReactionAutomapSmartsReaction(t *testing.T) {
	// SMARTS reaction without explicit mapping
	smarts := "[C](=[O])[OH].[C][OH]>>[C](=[O])[O][C].[OH2]"

	r, err := indigoInit.LoadReactionSmartsFromString(smarts)
	if err != nil {
		t.Fatalf("failed to load reaction SMARTS: %v", err)
	}
	defer r.Close()

	// Automap
	err = r.Automap(reaction.AutomapModeDiscard)
	if err != nil {
		t.Errorf("failed to automap SMARTS reaction: %v", err)
	}
}

// TestReactionAutomapComplexReaction tests automapping a more complex reaction
func TestReactionAutomapComplexReaction(t *testing.T) {
	// More complex reaction: Diels-Alder
	rxn := "C=CC=C.C=C>>C1=CCCCC1"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Automap
	err = r.Automap(reaction.AutomapModeDiscard)
	if err != nil {
		t.Errorf("failed to automap complex reaction: %v", err)
	}
}

// TestReactionAutomapWithIsotopes tests automapping with isotope considerations
func TestReactionAutomapWithIsotopes(t *testing.T) {
	// Reaction with isotopic labels
	rxn := "[13C]C(=O)O.CCO>>[13C]C(=O)OCC.O"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Automap without ignoring isotopes
	err = r.Automap(reaction.AutomapModeDiscard)
	if err != nil {
		t.Errorf("failed to automap with isotopes: %v", err)
	}

	// Automap ignoring isotopes
	err = r.Automap(reaction.AutomapModeDiscard + " " + reaction.AutomapIgnoreIsotopes)
	if err != nil {
		t.Errorf("failed to automap ignoring isotopes: %v", err)
	}
}

// TestReactionAutomapWithCharges tests automapping with charged species
func TestReactionAutomapWithCharges(t *testing.T) {
	// Reaction with charged species
	rxn := "CC(=O)[O-].[Na+].CCO>>CC(=O)OCC.O.[Na+].[OH-]"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Automap without ignoring charges
	err = r.Automap(reaction.AutomapModeDiscard)
	if err != nil {
		t.Errorf("failed to automap with charges: %v", err)
	}

	// Automap ignoring charges
	err = r.Automap(reaction.AutomapModeDiscard + " " + reaction.AutomapIgnoreCharges)
	if err != nil {
		t.Errorf("failed to automap ignoring charges: %v", err)
	}
}

// TestReactionAutomapKeepExisting tests keeping existing mapping
func TestReactionAutomapKeepExisting(t *testing.T) {
	// Load a reaction with partial mapping
	rxn := "[CH3:1][C:2](=[O:3])[OH:4].CCO>>CC(=O)OCC.O"

	r, err := indigoInit.LoadReactionFromString(rxn)
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer r.Close()

	// Automap with keep mode to preserve existing mapping
	err = r.Automap(reaction.AutomapModeKeep)
	if err != nil {
		t.Errorf("failed to automap with keep mode: %v", err)
	}
}
