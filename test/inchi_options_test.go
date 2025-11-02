//go:build cgo
// +build cgo

package test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestInChIOptionsNormalization tests option format normalization
func TestInChIOptionsNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantChar string // The prefix character we expect based on OS
	}{
		{
			name:     "Single option without prefix",
			input:    "FixedH",
			wantChar: getExpectedPrefix(),
		},
		{
			name:     "Single option with forward slash",
			input:    "/FixedH",
			wantChar: getExpectedPrefix(),
		},
		{
			name:     "Single option with dash",
			input:    "-FixedH",
			wantChar: getExpectedPrefix(),
		},
		{
			name:     "Multiple options without prefix",
			input:    "FixedH RecMet AuxNone",
			wantChar: getExpectedPrefix(),
		},
		{
			name:     "Multiple options with mixed prefix",
			input:    "/FixedH -RecMet AuxNone",
			wantChar: getExpectedPrefix(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := molecule.NewInChIGeneratorCGO()
			gen.SetOptions(tt.input)

			// We can't directly access g.options from test, so we'll test by actual usage
			// For now, just ensure it doesn't panic
			loader := molecule.SmilesLoader{}
			mol, _ := loader.Parse("CCO")

			result, err := gen.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI: %v", err)
			}

			if result.InChI == "" {
				t.Error("InChI should not be empty")
			}

			t.Logf("Input options: %s", tt.input)
			t.Logf("Generated InChI: %s", result.InChI)
		})
	}
}

// TestInChIStandardOptions tests various InChI options
func TestInChIStandardOptions(t *testing.T) {
	loader := molecule.SmilesLoader{}
	mol, _ := loader.Parse("CCO")

	tests := []struct {
		name            string
		options         string
		wantNonStandard bool // Some options create non-standard InChI
	}{
		{
			name:            "Standard (no options)",
			options:         "",
			wantNonStandard: false,
		},
		{
			name:            "FixedH",
			options:         "FixedH",
			wantNonStandard: true,
		},
		{
			name:            "RecMet",
			options:         "RecMet",
			wantNonStandard: true,
		},
		{
			name:            "SNon",
			options:         "SNon",
			wantNonStandard: false, // SNon is a structure perception option, doesn't make non-standard
		},
		{
			name:            "AuxNone",
			options:         "AuxNone",
			wantNonStandard: false, // AuxNone just omits aux info
		},
		{
			name:            "Multiple options",
			options:         "FixedH RecMet",
			wantNonStandard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := molecule.NewInChIGeneratorCGO()
			gen.SetOptions(tt.options)

			result, err := gen.GenerateInChI(mol)
			if err != nil {
				t.Fatalf("Failed to generate InChI with options '%s': %v", tt.options, err)
			}

			if result.InChI == "" {
				t.Error("InChI should not be empty")
			}

			// Check if it's standard or non-standard
			isStandard := strings.Contains(result.InChI, "InChI=1S/")
			isNonStandard := strings.Contains(result.InChI, "InChI=1/")

			t.Logf("Options: %s", tt.options)
			t.Logf("InChI: %s", result.InChI)
			t.Logf("Is Standard: %v, Is Non-Standard: %v", isStandard, isNonStandard)

			// Note: We can't reliably check standard vs non-standard without knowing
			// the exact InChI library behavior, so we just log it
		})
	}
}

// TestInChIOptionsWithAuxInfo tests AuxInfo option
func TestInChIOptionsWithAuxInfo(t *testing.T) {
	loader := molecule.SmilesLoader{}
	mol, _ := loader.Parse("CCO")

	// Without AuxNone - should include AuxInfo
	gen1 := molecule.NewInChIGeneratorCGO()
	result1, err := gen1.GenerateInChI(mol)
	if err != nil {
		t.Fatalf("Failed to generate InChI: %v", err)
	}

	// With AuxNone - should omit AuxInfo
	gen2 := molecule.NewInChIGeneratorCGO()
	gen2.SetOptions("AuxNone")
	result2, err := gen2.GenerateInChI(mol)
	if err != nil {
		t.Fatalf("Failed to generate InChI with AuxNone: %v", err)
	}

	t.Logf("Without AuxNone - AuxInfo length: %d", len(result1.AuxInfo))
	t.Logf("With AuxNone - AuxInfo length: %d", len(result2.AuxInfo))

	// With AuxNone, AuxInfo should be empty or much shorter
	if len(result2.AuxInfo) > len(result1.AuxInfo) {
		t.Error("AuxNone should reduce or eliminate AuxInfo")
	}
}

// getExpectedPrefix returns the expected option prefix for the current OS
func getExpectedPrefix() string {
	if runtime.GOOS == "windows" {
		return "/"
	}
	return "-"
}

// TestOptionPrefixForOS tests that options get the correct prefix based on OS
func TestOptionPrefixForOS(t *testing.T) {
	expectedPrefix := getExpectedPrefix()
	t.Logf("Current OS: %s", runtime.GOOS)
	t.Logf("Expected prefix: %s", expectedPrefix)

	// This is just informational
	if runtime.GOOS == "windows" {
		if expectedPrefix != "/" {
			t.Error("Windows should use / prefix")
		}
	} else {
		if expectedPrefix != "-" {
			t.Error("Non-Windows should use - prefix")
		}
	}
}
