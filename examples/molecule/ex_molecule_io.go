// Package main demonstrates molecule I/O operations
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_io.go
// @Software: GoLand
package main

import (
	"fmt"
	"github.com/cx-luo/go-chem/core"
	"log"
	"os"
)

func main() {
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Molecule I/O Examples ===\n")

	// Example 1: Load from SMILES
	fmt.Println("1. Loading from SMILES:")
	m1, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatalf("Failed to load from SMILES: %v", err)
	}
	defer m1.Close()
	fmt.Println("   Loaded ethanol from SMILES: CCO\n")

	// Example 2: Convert to SMILES
	fmt.Println("2. Converting to SMILES:")
	smiles, err := m1.ToSmiles()
	if err != nil {
		log.Fatalf("Failed to convert to SMILES: %v", err)
	}
	fmt.Printf("   SMILES: %s\n", smiles)

	canonicalSmiles, err := m1.ToCanonicalSmiles()
	if err != nil {
		log.Fatalf("Failed to get canonical SMILES: %v", err)
	}
	fmt.Printf("   Canonical SMILES: %s\n\n", canonicalSmiles)

	// Example 3: Save to MOL file
	fmt.Println("3. Saving to MOL file:")
	err = m1.SaveToFile("ethanol.mol")
	if err != nil {
		log.Fatalf("Failed to save to file: %v", err)
	}
	fmt.Println("   Saved to: ethanol.mol")

	// Clean up
	defer os.Remove("ethanol.mol")

	// Example 4: Load from MOL file
	fmt.Println("4. Loading from MOL file:")
	m2, err := indigoInit.LoadMoleculeFromFile("ethanol.mol")
	if err != nil {
		log.Fatalf("Failed to load from file: %v", err)
	}
	defer m2.Close()
	atomCount, _ := m2.CountAtoms()
	fmt.Printf("   Loaded molecule with %d atoms\n\n", atomCount)

	// Example 5: Convert to MOL format string
	fmt.Println("5. Converting to MOL format string:")
	molString, err := m1.ToMolfile()
	if err != nil {
		log.Fatalf("Failed to convert to MOL: %v", err)
	}
	fmt.Printf("   MOL format (first 100 chars): %s...\n\n", molString[:min(100, len(molString))])

	// Example 6: Save to JSON
	fmt.Println("6. Saving to JSON:")
	err = m1.SaveToJSONFile("ethanol.json")
	if err != nil {
		log.Fatalf("Failed to save JSON: %v", err)
	}
	fmt.Println("   Saved to: ethanol.json")
	defer os.Remove("ethanol.json")

	jsonStr, err := m1.ToJSON()
	if err != nil {
		log.Fatalf("Failed to convert to JSON: %v", err)
	}
	fmt.Printf("   JSON (first 100 chars): %s...\n\n", jsonStr[:min(100, len(jsonStr))])

	// Example 7: Load from buffer
	fmt.Println("7. Loading from buffer:")
	smilesBuffer := []byte("c1ccccc1")
	m3, err := indigoInit.LoadMoleculeFromBuffer(smilesBuffer)
	if err != nil {
		log.Fatalf("Failed to load from buffer: %v", err)
	}
	defer m3.Close()
	fmt.Println("   Loaded benzene from buffer\n")

	// Example 8: Load SMARTS
	fmt.Println("8. Loading SMARTS pattern:")
	pattern, err := indigoInit.LoadSmartsFromString("[OH]")
	if err != nil {
		log.Fatalf("Failed to load SMARTS: %v", err)
	}
	defer pattern.Close()
	fmt.Println("   Loaded SMARTS pattern: [OH]\n")

	// Example 9: Load query molecule
	fmt.Println("9. Loading query molecule:")
	query, err := indigoInit.LoadQueryMoleculeFromString("[#6]CO")
	if err != nil {
		log.Fatalf("Failed to load query: %v", err)
	}
	defer query.Close()
	fmt.Println("   Loaded query molecule: [#6]CO\n")

	// Example 10: Multiple molecules with components
	fmt.Println("10. Loading mixture:")
	mixture, err := indigoInit.LoadMoleculeFromString("CCO.C.O")
	if err != nil {
		log.Fatalf("Failed to load mixture: %v", err)
	}
	defer mixture.Close()
	components, _ := mixture.CountComponents()
	fmt.Printf("   Loaded mixture with %d components\n\n", components)

	fmt.Println("=== I/O Examples completed successfully ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
