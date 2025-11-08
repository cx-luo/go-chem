// Package main demonstrates InChI functionality usage
package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	// Initialize Indigo session
	if err := molecule.InitInChI(); err != nil {
		log.Fatalf("Failed to initialize Indigo: %v", err)
	}
	defer molecule.DisposeInChI()

	// Print InChI version
	fmt.Println("InChI Version:", molecule.InChIVersion())

	// Example 1: Convert SMILES to InChI
	fmt.Println("\n=== Example 1: Convert SMILES to InChI ===")
	if err := example1(); err != nil {
		log.Printf("Example 1 failed: %v", err)
	}

	// Example 2: Load molecule from InChI
	fmt.Println("\n=== Example 2: Load from InChI ===")
	if err := example2(); err != nil {
		log.Printf("Example 2 failed: %v", err)
	}

	// Example 3: Get InChI with detailed information
	fmt.Println("\n=== Example 3: InChI with Info ===")
	if err := example3(); err != nil {
		log.Printf("Example 3 failed: %v", err)
	}

	// Example 4: Convert InChI to InChIKey
	fmt.Println("\n=== Example 4: InChI to InChIKey ===")
	if err := example4(); err != nil {
		log.Printf("Example 4 failed: %v", err)
	}
}

// example1 demonstrates converting SMILES to InChI
func example1() error {
	// Load molecule from SMILES
	mol, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		return fmt.Errorf("failed to load molecule: %w", err)
	}
	defer mol.Close()

	// Convert to InChI
	inchi, err := mol.ToInChI()
	if err != nil {
		return fmt.Errorf("failed to convert to InChI: %w", err)
	}

	fmt.Printf("SMILES: CCO\n")
	fmt.Printf("InChI:  %s\n", inchi)

	// Get InChIKey
	key, err := mol.ToInChIKey()
	if err != nil {
		return fmt.Errorf("failed to generate InChIKey: %w", err)
	}

	fmt.Printf("InChIKey: %s\n", key)

	// Check for warnings
	if warning := molecule.InChIWarning(); warning != "" {
		fmt.Printf("Warning: %s\n", warning)
	}

	return nil
}

// example2 demonstrates loading a molecule from InChI
func example2() error {
	// InChI for ethanol (CCO)
	inchi := "InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3"

	// Load molecule from InChI
	mol, err := molecule.LoadFromInChI(inchi)
	if err != nil {
		return fmt.Errorf("failed to load from InChI: %w", err)
	}
	defer mol.Close()

	fmt.Printf("Loaded from InChI: %s\n", inchi)

	// Convert back to SMILES to verify
	smiles, err := mol.ToSmiles()
	if err != nil {
		return fmt.Errorf("failed to convert to SMILES: %w", err)
	}

	fmt.Printf("Converted to SMILES: %s\n", smiles)

	// Get atom count
	atomCount, err := mol.CountAtoms()
	if err != nil {
		return fmt.Errorf("failed to count atoms: %w", err)
	}

	fmt.Printf("Number of atoms: %d\n", atomCount)

	return nil
}

// example3 demonstrates getting InChI with detailed information
func example3() error {
	// Load molecule from SMILES (benzene)
	mol, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		return fmt.Errorf("failed to load molecule: %w", err)
	}
	defer mol.Close()

	// Get InChI with detailed information
	result, err := mol.ToInChIWithInfo()
	if err != nil {
		return fmt.Errorf("failed to generate InChI: %w", err)
	}

	fmt.Printf("SMILES: c1ccccc1 (benzene)\n")
	fmt.Printf("InChI:  %s\n", result.InChI)
	fmt.Printf("InChIKey: %s\n", result.Key)

	if result.Warning != "" {
		fmt.Printf("Warning: %s\n", result.Warning)
	}

	if result.Log != "" {
		fmt.Printf("Log: %s\n", result.Log)
	}

	if result.AuxInfo != "" {
		fmt.Printf("AuxInfo: %s\n", result.AuxInfo)
	}

	return nil
}

// example4 demonstrates converting InChI to InChIKey directly
func example4() error {
	// Various InChI strings
	inchis := []string{
		"InChI=1S/CH4/h1H4",                  // Methane
		"InChI=1S/C2H6/c1-2/h1-2H3",          // Ethane
		"InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H", // Benzene
		"InChI=1S/H2O/h1H2",                  // Water
	}

	for _, inchi := range inchis {
		key, err := molecule.InChIToKey(inchi)
		if err != nil {
			return fmt.Errorf("failed to generate InChIKey: %w", err)
		}

		fmt.Printf("InChI:    %s\n", inchi)
		fmt.Printf("InChIKey: %s\n\n", key)
	}

	return nil
}
