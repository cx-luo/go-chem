//go:build ignore
// +build ignore

// Example program demonstrating reaction package usage
// Build and run: go run example_reaction.go
package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/reaction"
)

func main() {
	fmt.Println("=== Reaction Package Example ===\n")

	// Example 1: Load a reaction from SMILES
	fmt.Println("1. Loading reaction from SMILES:")
	rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
	fmt.Printf("   Input: %s\n", rxn)

	r, err := reaction.LoadReactionFromString(rxn)
	if err != nil {
		log.Fatalf("Failed to load reaction: %v", err)
	}
	defer r.Close()

	// Get counts
	reactantCount, _ := r.CountReactants()
	productCount, _ := r.CountProducts()
	catalystCount, _ := r.CountCatalysts()

	fmt.Printf("   Reactants: %d\n", reactantCount)
	fmt.Printf("   Products: %d\n", productCount)
	fmt.Printf("   Catalysts: %d\n\n", catalystCount)

	// Example 2: Convert to canonical SMILES
	fmt.Println("2. Converting to canonical SMILES:")
	canonicalSmiles, err := r.ToCanonicalSmiles()
	if err != nil {
		log.Fatalf("Failed to convert to canonical SMILES: %v", err)
	}
	fmt.Printf("   Canonical: %s\n\n", canonicalSmiles)

	// Example 3: Automatic atom mapping
	fmt.Println("3. Performing automatic atom mapping:")
	err = r.Automap(reaction.AutomapModeDiscard)
	if err != nil {
		log.Fatalf("Failed to automap: %v", err)
	}
	fmt.Println("   Mapping completed successfully\n")

	// Example 4: Convert to RXN format
	fmt.Println("4. Converting to RXN format:")
	rxnString, err := r.ToRxnfile()
	if err != nil {
		log.Fatalf("Failed to convert to RXN: %v", err)
	}
	fmt.Printf("   RXN length: %d bytes\n", len(rxnString))
	fmt.Println("   First line:", rxnString[:50], "...\n")

	// Example 5: Iterate over reactants
	fmt.Println("5. Iterating over reactants:")
	reactIter, err := r.IterateReactants()
	if err != nil {
		log.Fatalf("Failed to create iterator: %v", err)
	}
	defer reactIter.Close()

	idx := 0
	for reactIter.HasNext() {
		molHandle, err := reactIter.Next()
		if err != nil {
			log.Printf("Failed to get next molecule: %v", err)
			continue
		}
		fmt.Printf("   Reactant %d: handle=%d\n", idx, molHandle)
		idx++
	}
	fmt.Println()

	// Example 6: Load a reaction SMARTS
	fmt.Println("6. Loading reaction SMARTS:")
	smarts := "[C:1](=[O:2])[OH:3].[C:4][OH:5]>>[C:1](=[O:2])[O:5][C:4].[OH2:3]"
	fmt.Printf("   Input: %s\n", smarts)

	r2, err := reaction.LoadReactionSmartsFromString(smarts)
	if err != nil {
		log.Fatalf("Failed to load reaction SMARTS: %v", err)
	}
	defer r2.Close()

	r2ReactantCount, _ := r2.CountReactants()
	r2ProductCount, _ := r2.CountProducts()
	fmt.Printf("   Reactants: %d\n", r2ReactantCount)
	fmt.Printf("   Products: %d\n\n", r2ProductCount)

	// Example 7: Clone a reaction
	fmt.Println("7. Cloning a reaction:")
	r3, err := r.Clone()
	if err != nil {
		log.Fatalf("Failed to clone reaction: %v", err)
	}
	defer r3.Close()

	clonedReactants, _ := r3.CountReactants()
	clonedProducts, _ := r3.CountProducts()
	fmt.Printf("   Cloned reaction has %d reactants and %d products\n\n", clonedReactants, clonedProducts)

	// Example 8: Create a new empty reaction
	fmt.Println("8. Creating a new empty reaction:")
	r4, err := reaction.CreateReaction()
	if err != nil {
		log.Fatalf("Failed to create reaction: %v", err)
	}
	defer r4.Close()
	fmt.Println("   Empty reaction created successfully\n")

	fmt.Println("=== All examples completed successfully! ===")
}
