// Package main demonstrates basic molecule operations using go-indigo
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : basic_usage.go
// @Software: GoLand
package main

import (
	"fmt"
	"github.com/cx-luo/go-indigo/core"
	"log"
)

func main() {
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	fmt.Println("=== go-indigo Molecule Basic Usage Examples ===\n")

	// Example 1: Create an empty molecule
	fmt.Println("1. Creating an empty molecule:")
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		log.Fatalf("Failed to create molecule: %v", err)
	}
	defer m.Close()
	fmt.Printf("   Created molecule with handle: %d\n\n", m.Handle)

	// Example 2: Load molecule from SMILES
	fmt.Println("2. Loading molecule from SMILES:")
	ethanol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatalf("Failed to load ethanol: %v", err)
	}
	defer ethanol.Close()

	atomCount, _ := ethanol.CountAtoms()
	bondCount, _ := ethanol.CountBonds()
	heavyAtoms, _ := ethanol.CountHeavyAtoms()
	fmt.Printf("   Ethanol (CCO):\n")
	fmt.Printf("   - Total atoms: %d\n", atomCount)
	fmt.Printf("   - Heavy atoms: %d\n", heavyAtoms)
	fmt.Printf("   - Bonds: %d\n\n", bondCount)

	// Example 3: Clone a molecule
	fmt.Println("3. Cloning a molecule:")
	ethanolClone, err := ethanol.Clone()
	if err != nil {
		log.Fatalf("Failed to clone molecule: %v", err)
	}
	defer ethanolClone.Close()
	fmt.Println("   Successfully cloned ethanol molecule\n")

	// Example 4: Aromatize a molecule
	fmt.Println("4. Aromatization:")
	benzene, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatalf("Failed to load benzene: %v", err)
	}
	defer benzene.Close()

	err = benzene.Aromatize()
	if err != nil {
		log.Printf("   Aromatization warning: %v\n", err)
	} else {
		fmt.Println("   Benzene aromatized successfully")
	}

	// Count rings
	ringCount, _ := benzene.CountSSSR()
	fmt.Printf("   Benzene rings (SSSR): %d\n\n", ringCount)

	// Example 5: Fold/Unfold hydrogens
	fmt.Println("5. Hydrogen management:")
	molecule2, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatalf("Failed to load molecule: %v", err)
	}
	defer molecule2.Close()

	fmt.Println("   Folding hydrogens...")
	err = molecule2.FoldHydrogens()
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Println("   Hydrogens folded successfully")
	}

	fmt.Println("   Unfolding hydrogens...")
	err = molecule2.UnfoldHydrogens()
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Println("   Hydrogens unfolded successfully")
	}
	fmt.Println()

	// Example 6: Layout molecule (2D coordinates)
	fmt.Println("6. 2D Layout:")
	molecule3, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatalf("Failed to load molecule: %v", err)
	}
	defer molecule3.Close()

	err = molecule3.Layout()
	if err != nil {
		log.Printf("   Layout error: %v\n", err)
	} else {
		fmt.Println("   2D layout generated successfully")
	}
	fmt.Println()

	// Example 7: Normalize and Standardize
	fmt.Println("7. Normalization and Standardization:")
	molecule4, err := indigoInit.LoadMoleculeFromString("CC(=O)O")
	if err != nil {
		log.Fatalf("Failed to load molecule: %v", err)
	}
	defer molecule4.Close()

	err = molecule4.Normalize("")
	if err != nil {
		log.Printf("   Normalization error: %v\n", err)
	} else {
		fmt.Println("   Molecule normalized successfully")
	}

	err = molecule4.Standardize()
	if err != nil {
		log.Printf("   Standardization error: %v\n", err)
	} else {
		fmt.Println("   Molecule standardized successfully")
	}
	fmt.Println()

	// Example 8: Count components
	fmt.Println("8. Counting connected components:")
	mixture, err := indigoInit.LoadMoleculeFromString("CCO.C.O")
	if err != nil {
		log.Fatalf("Failed to load mixture: %v", err)
	}
	defer mixture.Close()

	components, _ := mixture.CountComponents()
	fmt.Printf("   Mixture 'CCO.C.O' has %d components\n\n", components)

	fmt.Println("=== Examples completed successfully ===")
}
