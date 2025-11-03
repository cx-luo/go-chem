// Package main provides examples of using InChI with CGO bindings
//
// This example demonstrates how to use the CGO-based InChI generator
// which calls the official InChI library (libinchi.dll/libinchi.so)
//
// Build instructions:
//   Windows: go build -tags=cgo examples/inchi_cgo_example.go
//   Linux:   go build -tags=cgo examples/inchi_cgo_example.go
//
// Make sure the InChI library is in the 3rd/ directory:
//   - Windows: 3rd/libinchi.dll
//   - Linux:   3rd/libinchi.so

package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	fmt.Println("=== InChI CGO Examples ===\n")

	// Example 1: Generate InChI from SMILES using CGO
	example1()

	// Example 2: Compare CGO vs Pure Go implementation
	example2()

	// Example 3: Generate InChIKey using CGO
	example3()

	// Example 4: Advanced molecules
	example4()
}

// Example 1: Simple InChI generation using CGO
func example1() {
	fmt.Println("Example 1: Generate InChI using CGO")
	fmt.Println("------------------------------------")

	molecules := []struct {
		name   string
		smiles string
	}{
		{"Methane", "C"},
		{"Ethanol", "CCO"},
		{"Benzene", "c1ccccc1"},
		{"Acetic acid", "CC(=O)O"},
	}

	// Create CGO-based generator
	generator := molecule.NewInChIGeneratorCGO()

	for _, mol := range molecules {
		// Parse SMILES
		moleculeFromString, err := molecule.LoadMoleculeFromString(mol.smiles)
		if err != nil {
			return
		}
		if err != nil {
			log.Printf("Error parsing SMILES %s: %v\n", mol.name, err)
			continue
		}

		// Generate InChI using CGO
		result, err := generator.GenerateInChI(moleculeFromString)
		if err != nil {
			log.Printf("Error generating InChI for %s: %v\n", mol.name, err)
			continue
		}

		fmt.Printf("Molecule: %s\n", mol.name)
		fmt.Printf("  SMILES:   %s\n", mol.smiles)
		fmt.Printf("  InChI:    %s\n", result.InChI)
		fmt.Printf("  InChIKey: %s\n", result.InChIKey)

		if len(result.Warnings) > 0 {
			fmt.Printf("  Warnings: %v\n", result.Warnings)
		}
		fmt.Println()
	}
}

// Example 2: Compare CGO vs Pure Go
func example2() {
	fmt.Println("\nExample 2: Compare CGO vs Pure Go")
	fmt.Println("----------------------------------")

	smiles := "CCO" // Ethanol

	// Parse SMILES
	mol, err := molecule.LoadMoleculeFromString(smiles)
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	// Generate using CGO
	generatorCGO := molecule.NewInChIGeneratorCGO()
	resultCGO, err := generatorCGO.GenerateInChI(mol)
	if err != nil {
		log.Printf("CGO error: %v\n", err)
	} else {
		fmt.Printf("CGO InChI:    %s\n", resultCGO.InChI)
		fmt.Printf("CGO InChIKey: %s\n", resultCGO.InChIKey)
	}

	fmt.Println()
}

// Example 3: Generate InChIKey from InChI
func example3() {
	fmt.Println("\nExample 3: Generate InChIKey using CGO")
	fmt.Println("---------------------------------------")

	inchis := []string{
		"InChI=1S/CH4/h1H4",
		"InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3",
		"InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H",
	}

	for _, inchi := range inchis {
		key, err := molecule.GenerateInChIKeyCGO(inchi)
		if err != nil {
			log.Printf("Error generating InChIKey: %v\n", err)
			continue
		}

		fmt.Printf("InChI:    %s\n", inchi)
		fmt.Printf("InChIKey: %s\n\n", key)
	}
}

// Example 4: Advanced molecules
func example4() {
	fmt.Println("\nExample 4: Advanced Molecules")
	fmt.Println("------------------------------")

	molecules := []struct {
		name   string
		smiles string
	}{
		{"L-Alanine", "C[C@@H](C(=O)O)N"},
		{"Glucose", "C([C@@H]1[C@H]([C@@H]([C@H](C(O1)O)O)O)O)O"},
		{"Caffeine", "CN1C=NC2=C1C(=O)N(C(=O)N2C)C"},
	}

	generator := molecule.NewInChIGeneratorCGO()

	for _, mol := range molecules {
		m, err := molecule.LoadMoleculeFromString(mol.smiles)
		if err != nil {
			log.Printf("Error parsing SMILES %s: %v\n", mol.name, err)
			continue
		}

		result, err := generator.GenerateInChI(m)
		if err != nil {
			log.Printf("Error generating InChI for %s: %v\n", mol.name, err)
			continue
		}

		fmt.Printf("Molecule: %s\n", mol.name)
		fmt.Printf("  SMILES:   %s\n", mol.smiles)
		fmt.Printf("  InChI:    %s\n", result.InChI)
		fmt.Printf("  InChIKey: %s\n", result.InChIKey)
		fmt.Println()
	}
}

// Example: Batch processing with error handling
func exampleBatchProcessing() {
	fmt.Println("\nExample: Batch Processing with CGO")
	fmt.Println("-----------------------------------")

	smilesList := []string{
		"C",        // Methane
		"CC",       // Ethane
		"CCC",      // Propane
		"c1ccccc1", // Benzene
		"CCO",      // Ethanol
	}

	generator := molecule.NewInChIGeneratorCGO()

	successCount := 0
	for i, smiles := range smilesList {
		mol, err := molecule.LoadMoleculeFromString(smiles)
		if err != nil {
			fmt.Printf("%d. Error parsing: %v\n", i+1, err)
			continue
		}

		result, err := generator.GenerateInChI(mol)
		if err != nil {
			fmt.Printf("%d. Error generating: %v\n", i+1, err)
			continue
		}

		fmt.Printf("%d. %-10s -> %s\n", i+1, smiles, result.InChIKey)
		successCount++
	}

	fmt.Printf("\nProcessed: %d/%d molecules\n", successCount, len(smilesList))
}

// Example: Options
func exampleWithOptions() {
	fmt.Println("\nExample: InChI with Custom Options")
	fmt.Println("-----------------------------------")

	smiles := "CCO"

	mol, _ := molecule.LoadMoleculeFromString(smiles)
	// Standard InChI
	gen1 := molecule.NewInChIGeneratorCGO()
	result1, _ := gen1.GenerateInChI(mol)
	fmt.Printf("Standard:  %s\n", result1.InChI)

	// With FixedH option (will be auto-prefixed with / or - based on OS)
	gen2 := molecule.NewInChIGeneratorCGO()
	gen2.SetOptions("FixedH")
	result2, _ := gen2.GenerateInChI(mol)
	fmt.Printf("FixedH:    %s\n", result2.InChI)

	// With RecMet option
	gen3 := molecule.NewInChIGeneratorCGO()
	gen3.SetOptions("RecMet")
	result3, _ := gen3.GenerateInChI(mol)
	fmt.Printf("RecMet:    %s\n", result3.InChI)

	// With multiple options
	gen4 := molecule.NewInChIGeneratorCGO()
	gen4.SetOptions("FixedH RecMet AuxNone")
	result4, _ := gen4.GenerateInChI(mol)
	fmt.Printf("Multiple:  %s\n", result4.InChI)
}
