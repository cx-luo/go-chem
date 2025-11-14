// Package main demonstrates substructure matching operations
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : substructure_matching.go
package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/core"
	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Substructure Matching Examples ===\n")

	// Example 1: Basic substructure search
	fmt.Println("1. Basic Substructure Search:")

	target, err := indigoInit.LoadMoleculeFromString("c1ccccc1CCO")
	if err != nil {
		log.Fatalf("Failed to load target molecule: %v", err)
	}
	defer target.Close()

	query, err := indigoInit.LoadQueryMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatalf("Failed to load query molecule: %v", err)
	}
	defer query.Close()

	fmt.Println("  Target: c1ccccc1CCO (Phenylethanol)")
	fmt.Println("  Query:  c1ccccc1 (Benzene)")

	hasSubstruct, _ := target.HasSubstructure(query, nil)
	fmt.Printf("  Has benzene ring: %v\n", hasSubstruct)

	count, _ := target.CountSubstructureMatches(query, nil)
	fmt.Printf("  Number of matches: %d\n", count)

	// Example 2: Multiple substructure matches
	fmt.Println("\n2. Multiple Substructure Matches:")

	target2, _ := indigoInit.LoadMoleculeFromString("c1ccc(cc1)c2ccccc2")
	defer target2.Close()
	query2, _ := indigoInit.LoadQueryMoleculeFromString("c1ccccc1")
	defer query2.Close()

	fmt.Println("  Target: c1ccc(cc1)c2ccccc2 (Biphenyl)")
	fmt.Println("  Query:  c1ccccc1 (Benzene)")

	count2, _ := target2.CountSubstructureMatches(query2, nil)
	fmt.Printf("  Number of benzene rings: %d\n", count2)

	// Example 3: Exact match
	fmt.Println("\n3. Exact Matching:")

	mol1, _ := indigoInit.LoadMoleculeFromString("CCO")
	defer mol1.Close()
	mol2, _ := indigoInit.LoadMoleculeFromString("CCO")
	defer mol2.Close()
	mol3, _ := indigoInit.LoadMoleculeFromString("OCC")
	defer mol3.Close()

	flags := "ALL"
	isExact1, _, _ := mol1.ExactMatch(mol2, &flags)
	isExact2, mapping2, _ := mol1.ExactMatch(mol3, &flags)

	// fixed: demonstrate atom mapping with correct handles
	if isExact2 && mapping2 != 0 {
		// Try to map the first atom (index 0) of mol1 to mol3
		queryAtom, err := mol1.GetAtom(0)
		if err != nil {
			fmt.Printf("Failed to get first atom from mol1: %v\n", err)
		} else {
			// Use atom handle from queryAtom to map
			mappedAtomHandle := molecule.MapAtom(mapping2, queryAtom.Handle)
			if mappedAtomHandle <= 0 {
				fmt.Printf("  CCO matches OCC: true (atom not mapped)\n")
			} else {
				targetAtomIndex, err := indigoInit.Index(mappedAtomHandle)
				if err != nil {
					fmt.Printf("Failed to get mapped atom index: %v\n", err)
				} else {
					targetAtom, err := mol3.GetAtom(targetAtomIndex)
					if err != nil {
						fmt.Printf("Failed to get atom from OCC: %v\n", err)
					} else {
						symbol, err := targetAtom.Symbol()
						if err != nil {
							fmt.Printf("Failed to get atom symbol: %v\n", err)
						} else {
							fmt.Printf("  CCO matches OCC: %v (first mapped atom: %s)\n", isExact2, symbol)
						}
					}
				}
			}
		}
	}

	// When ExactMatch returns true, mapping1 is a mapping handle that can be used
	// to query atom-to-atom correspondences between the molecules
	fmt.Printf("  CCO matches CCO: %v (mapping handle: %d)\n", isExact1, mapping2)
	fmt.Printf("  CCO matches OCC: %v (same molecule, different notation)\n", isExact2)

	// Example 4: SMARTS pattern matching
	fmt.Println("\n4. SMARTS Pattern Matching:")

	aspirin, _ := indigoInit.LoadMoleculeFromString("CC(=O)Oc1ccccc1C(=O)O")
	defer aspirin.Close()

	// Search for carboxylic acid group
	carboxyl, _ := indigoInit.LoadQueryMoleculeFromString("[CX3](=O)[OX2H1]")
	defer carboxyl.Close()

	fmt.Println("  Molecule: Aspirin")
	fmt.Println("  Pattern: Carboxylic acid [CX3](=O)[OX2H1]")

	hasCarboxyl, _ := aspirin.HasSubstructure(carboxyl, nil)
	carboxylCount, _ := aspirin.CountSubstructureMatches(carboxyl, nil)

	fmt.Printf("  Has carboxylic acid: %v\n", hasCarboxyl)
	fmt.Printf("  Number of groups: %d\n", carboxylCount)

	// Example 5: Functional group detection
	fmt.Println("\n5. Functional Group Detection:")

	testMolecules := []struct {
		smiles string
		name   string
	}{
		{"CCO", "Ethanol"},
		{"CC(=O)C", "Acetone"},
		{"CC(=O)O", "Acetic acid"},
		{"CCN", "Ethylamine"},
	}

	functionalGroups := []struct {
		smarts string
		name   string
	}{
		{"[OX2H]", "Alcohol"},
		{"[CX3](=O)", "Carbonyl"},
		{"[CX3](=O)[OX2H1]", "Carboxylic acid"},
		{"[NX3;H2,H1;!$(NC=O)]", "Primary/Secondary amine"},
	}

	for _, testMol := range testMolecules {
		mol, err := indigoInit.LoadMoleculeFromString(testMol.smiles)
		if err != nil {
			continue
		}

		fmt.Printf("\n  %s (%s):\n", testMol.name, testMol.smiles)

		for _, fg := range functionalGroups {
			pattern, err := indigoInit.LoadQueryMoleculeFromString(fg.smarts)
			if err != nil {
				continue
			}

			has, _ := mol.HasSubstructure(pattern, nil)
			if has {
				count, _ := mol.CountSubstructureMatches(pattern, nil)
				fmt.Printf("    \u2713 %s (%d matches)\n", fg.name, count)
			}

			pattern.Close()
		}

		mol.Close()
	}

	// Example 6: Submolecule extraction
	fmt.Println("\n6. Submolecule Extraction:")

	bigMol, _ := indigoInit.LoadMoleculeFromString("c1ccccc1CCO")
	defer bigMol.Close()

	fmt.Println("  Original: c1ccccc1CCO")
	fmt.Println("  Extracting atoms 0-5 (benzene ring)...")

	submol, err := bigMol.GetSubmolecule([]int{0, 1, 2, 3, 4, 5})
	if err != nil {
		log.Printf("Failed to get submolecule: %v", err)
	} else {
		defer submol.Close()
		subSmiles, _ := submol.ToSmiles()
		fmt.Printf("  Submolecule SMILES: %s\n", subSmiles)
	}

	fmt.Println("\n=== Examples completed successfully ===")
}
