// Package main demonstrates InChI functionality usage
package main

import (
	"fmt"
	"github.com/cx-luo/go-chem/core"
	"log"
	"strings"
	"sync"

	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	// Initialize Indigo session
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	indigoInchi, err := core.InchiInit(indigoInit.GetSessionID())
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatalf("Failed to initialize Indigo: %v", err)
	}
	defer core.DisposeInChI(indigoInchi.GetInchiSessionID())

	// Print InChI version
	fmt.Println("InChI Version:", indigoInchi.InChIVersion())

	// Example 1: Convert SMILES to InChI
	fmt.Println("\n=== Example 1: Convert SMILES to InChI ===")
	if err := example1(indigoInchi); err != nil {
		log.Printf("Example 1 failed: %v", err)
	}

	// Example 2: Load molecule from InChI
	fmt.Println("\n=== Example 2: Load from InChI ===")
	if err := example2(indigoInchi); err != nil {
		log.Printf("Example 2 failed: %v", err)
	}

	// Example 3: Get InChI with detailed information
	fmt.Println("\n=== Example 3: InChI with Info ===")
	if err := example3(indigoInchi); err != nil {
		log.Printf("Example 3 failed: %v", err)
	}

	// Example 4: Convert InChI to InChIKey
	fmt.Println("\n=== Example 4: InChI to InChIKey ===")
	if err := example4(indigoInchi); err != nil {
		log.Printf("Example 4 failed: %v", err)
	}

	example5(indigoInchi)
}

// example1 demonstrates converting SMILES to InChI
func example1(i *core.IndigoInchi) error {
	// Load molecule from SMILES
	mol, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		return fmt.Errorf("failed to load molecule: %w", err)
	}
	defer mol.Close()

	// Convert to InChI
	inchi, err := i.ToInChI(mol)
	if err != nil {
		return fmt.Errorf("failed to convert to InChI: %w", err)
	}

	fmt.Printf("SMILES: CCO\n")
	fmt.Printf("InChI:  %s\n", inchi)

	// Get InChIKey
	key, err := i.InChIToKey(inchi)
	if err != nil {
		return fmt.Errorf("failed to generate InChIKey: %w", err)
	}

	fmt.Printf("InChIKey: %s\n", key)

	// Check for warnings
	if warning := i.InChIWarning(); warning != "" {
		fmt.Printf("Warning: %s\n", warning)
	}

	return nil
}

// example2 demonstrates loading a molecule from InChI
func example2(i *core.IndigoInchi) error {
	// InChI for ethanol (CCO)
	inchi := "InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3"

	// Load molecule from InChI
	molHandle, err := i.LoadFromInChI(inchi)
	if err != nil {
		return fmt.Errorf("failed to load from InChI: %w", err)
	}

	fmt.Printf("Loaded from InChI: %s\n", inchi)

	mol, err := molecule.LoadMoleculeFromHandle(molHandle)
	if err != nil {
		panic(err)
	}
	defer mol.Close()
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
func example3(i *core.IndigoInchi) error {
	// Load molecule from SMILES (benzene)
	mol, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		return fmt.Errorf("failed to load molecule: %w", err)
	}
	defer mol.Close()

	// Get InChI with detailed information
	result, err := i.ToInChIWithInfo(mol)
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
func example4(i *core.IndigoInchi) error {
	// Various InChI strings
	inchis := []string{
		"InChI=1S/CH4/h1H4",                  // Methane
		"InChI=1S/C2H6/c1-2/h1-2H3",          // Ethane
		"InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H", // Benzene
		"InChI=1S/H2O/h1H2",                  // Water
	}

	for _, inchi := range inchis {
		key, err := i.InChIToKey(inchi)
		if err != nil {
			return fmt.Errorf("failed to generate InChIKey: %w", err)
		}

		fmt.Printf("InChI:    %s\n", inchi)
		fmt.Printf("InChIKey: %s\n\n", key)
	}

	return nil
}

func example5(ii *core.IndigoInchi) {
	// Create test molecules
	smilesList := []string{"CCO", "c1ccccc1", "CC(=O)O", "O"}
	molecules := make([]*molecule.Molecule, 0, len(smilesList))

	for _, smi := range smilesList {
		// Skip if LoadMoleculeFromString fails (maybe due to other issues)
		mol, err := molecule.LoadMoleculeFromString(smi)
		if err != nil {
			fmt.Printf("Skipping %s due to load error: %v\n", smi, err)
			continue
		}
		molecules = append(molecules, mol)
	}

	if len(molecules) == 0 {
		fmt.Println("No molecules loaded successfully")
	}

	// Concurrently convert to InChI
	var wg sync.WaitGroup
	results := make([]string, len(molecules))
	errors := make([]error, len(molecules))

	for i := range molecules {
		wg.Add(1)
		go func(idx int, m *molecule.Molecule) {
			defer wg.Done()
			defer m.Close()

			inchi, err := ii.ToInChI(m)
			results[idx] = inchi
			errors[idx] = err
		}(i, molecules[i])
	}

	wg.Wait()

	// Verify results
	for i := range molecules {
		if errors[i] != nil {
			fmt.Printf("Molecule %d failed: %v", i, errors[i])
			continue
		}

		if results[i] == "" {
			fmt.Printf("Molecule %d: empty InChI", i)
			continue
		}

		if !strings.HasPrefix(results[i], "InChI=") {
			fmt.Printf("Molecule %d: invalid InChI format: %s", i, results[i])
		}
	}

	// Clean up molecules
	for _, mol := range molecules {
		mol.Close()
	}
}
