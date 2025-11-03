// Package main demonstrates InChI generation and usage
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_inchi.go
// @Software: GoLand
package main

import (
	"fmt"
	"log"

	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	fmt.Println("=== InChI Examples ===\n")

	// Example 1: Generate InChI from SMILES
	fmt.Println("1. Generating InChI from SMILES:")
	m1, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatalf("Failed to load molecule: %v", err)
	}
	defer m1.Close()

	inchi, err := m1.ToInChI()
	if err != nil {
		log.Fatalf("Failed to generate InChI: %v", err)
	}
	fmt.Printf("   Ethanol SMILES: CCO\n")
	fmt.Printf("   InChI: %s\n\n", inchi)

	// Example 2: Generate InChI Key
	fmt.Println("2. Generating InChI Key:")
	inchiKey, err := m1.ToInChIKey()
	if err != nil {
		log.Fatalf("Failed to generate InChI Key: %v", err)
	}
	fmt.Printf("   InChI Key: %s\n\n", inchiKey)

	// Example 3: Load molecule from InChI
	fmt.Println("3. Loading molecule from InChI:")
	m2, err := molecule.LoadInChI(inchi)
	if err != nil {
		log.Fatalf("Failed to load from InChI: %v", err)
	}
	defer m2.Close()

	atomCount, _ := m2.CountAtoms()
	fmt.Printf("   Loaded molecule with %d atoms\n", atomCount)

	// Verify by converting back to SMILES
	smiles, _ := m2.ToSmiles()
	fmt.Printf("   SMILES from InChI: %s\n\n", smiles)

	// Example 4: InChI for various molecules
	fmt.Println("4. InChI for various molecules:")
	testMolecules := []struct {
		name   string
		smiles string
	}{
		{"Methanol", "CO"},
		{"Benzene", "c1ccccc1"},
		{"Acetic Acid", "CC(=O)O"},
		{"Aspirin", "CC(=O)Oc1ccccc1C(=O)O"},
		{"Glucose", "C([C@@H]1[C@H]([C@@H]([C@H](C(O1)O)O)O)O)O"},
	}

	for _, mol := range testMolecules {
		fmt.Printf("\n   %s (%s):\n", mol.name, mol.smiles)
		m, err := molecule.LoadMoleculeFromString(mol.smiles)
		if err != nil {
			log.Printf("   Error loading: %v\n", err)
			continue
		}

		inchi, err := m.ToInChI()
		if err != nil {
			log.Printf("   Error generating InChI: %v\n", err)
			m.Close()
			continue
		}

		inchiKey, err := m.ToInChIKey()
		if err != nil {
			log.Printf("   Error generating InChI Key: %v\n", err)
			m.Close()
			continue
		}

		fmt.Printf("   InChI: %s\n", inchi)
		fmt.Printf("   InChI Key: %s\n", inchiKey)

		m.Close()
	}

	// Example 5: InChI roundtrip test
	fmt.Println("\n5. InChI roundtrip test:")
	original, _ := molecule.LoadMoleculeFromString("c1ccccc1")
	defer original.Close()

	inchi1, _ := original.ToInChI()
	fmt.Printf("   Original InChI: %s\n", inchi1)

	reloaded, _ := molecule.LoadInChI(inchi1)
	defer reloaded.Close()

	inchi2, _ := reloaded.ToInChI()
	fmt.Printf("   Reloaded InChI: %s\n", inchi2)

	if inchi1 == inchi2 {
		fmt.Println("   ✓ Roundtrip successful: InChI preserved\n")
	} else {
		fmt.Println("   ✗ Roundtrip failed: InChI changed\n")
	}

	// Example 6: InChI helper functions
	fmt.Println("6. InChI warnings and logs:")
	m3, _ := molecule.LoadMoleculeFromString("CCO")
	defer m3.Close()

	m3.ToInChI()

	warning := molecule.GetInChIWarning()
	log := molecule.GetInChILog()
	auxInfo := molecule.GetInChIAuxInfo()

	if warning != "" {
		fmt.Printf("   Warning: %s\n", warning)
	} else {
		fmt.Println("   No warnings")
	}

	if log != "" {
		fmt.Printf("   Log: %s\n", log)
	}

	if auxInfo != "" {
		fmt.Printf("   Aux Info: %s\n", auxInfo)
	}
	fmt.Println()

	// Example 7: Verify InChI Key uniqueness
	fmt.Println("7. Verifying InChI Key uniqueness:")
	// Same molecule, different SMILES representations
	m4a, _ := molecule.LoadMoleculeFromString("CCO")
	m4b, _ := molecule.LoadMoleculeFromString("OCC")
	defer m4a.Close()
	defer m4b.Close()

	key4a, _ := m4a.ToInChIKey()
	key4b, _ := m4b.ToInChIKey()

	fmt.Printf("   InChI Key for 'CCO': %s\n", key4a)
	fmt.Printf("   InChI Key for 'OCC': %s\n", key4b)

	if key4a == key4b {
		fmt.Println("   ✓ Same InChI Key (as expected)\n")
	} else {
		fmt.Println("   ✗ Different InChI Keys (unexpected)\n")
	}

	fmt.Println("=== InChI Examples completed successfully ===")
}
