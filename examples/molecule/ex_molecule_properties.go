// Package main demonstrates molecular property calculations
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_properties.go
// @Software: GoLand
package main

import (
	"fmt"
	"github.com/cx-luo/go-indigo/core"
	"log"
)

func main() {
	fmt.Println("=== Molecular Properties Examples ===\n")

	// Test molecules
	molecules := []struct {
		name   string
		smiles string
	}{
		{"Ethanol", "CCO"},
		{"Benzene", "c1ccccc1"},
		{"Acetic Acid", "CC(=O)O"},
		{"Aspirin", "CC(=O)Oc1ccccc1C(=O)O"},
		{"Caffeine", "CN1C=NC2=C1C(=O)N(C(=O)N2C)C"},
	}

	for _, mol := range molecules {
		fmt.Printf("=== %s (%s) ===\n", mol.name, mol.smiles)
		analyzeProperties(mol.smiles)
		fmt.Println()
	}

	// Demonstrate property management
	fmt.Println("=== Custom Properties ===")
	demonstrateProperties()

	fmt.Println("\n=== Examples completed successfully ===")
}

func analyzeProperties(smiles string) {
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	m, err := indigoInit.LoadMoleculeFromString(smiles)
	if err != nil {
		log.Printf("Failed to load molecule: %v\n", err)
		return
	}
	defer m.Close()

	// Basic counts
	atomCount, _ := m.CountAtoms()
	bondCount, _ := m.CountBonds()
	heavyAtoms, _ := m.CountHeavyAtoms()
	fmt.Printf("Atom count: %d\n", atomCount)
	fmt.Printf("Heavy atoms: %d\n", heavyAtoms)
	fmt.Printf("Bond count: %d\n", bondCount)

	// Formula and mass
	grossFormula, err := m.GrossFormula()
	if err == nil {
		fmt.Printf("Gross formula: %s\n", grossFormula)
	}

	molecularFormula, err := m.MolecularFormula()
	if err == nil {
		fmt.Printf("Molecular formula: %s\n", molecularFormula)
	}

	mw, err := m.MolecularWeight()
	if err == nil {
		fmt.Printf("Molecular weight: %.2f\n", mw)
	}

	monoisotopicMass, err := m.MonoisotopicMass()
	if err == nil {
		fmt.Printf("Monoisotopic mass: %.4f\n", monoisotopicMass)
	}

	mostAbundantMass, err := m.MostAbundantMass()
	if err == nil {
		fmt.Printf("Most abundant mass: %.4f\n", mostAbundantMass)
	}

	// Mass composition
	massComp, err := m.MassComposition()
	if err == nil && massComp != "" {
		fmt.Printf("Mass composition: %s\n", massComp)
	}

	// TPSA (Topological Polar Surface Area)
	tpsa, err := m.TPSA(false)
	if err == nil {
		fmt.Printf("TPSA (without S/P): %.2f\n", tpsa)
	}

	tpsaSP, err := m.TPSA(true)
	if err == nil {
		fmt.Printf("TPSA (with S/P): %.2f\n", tpsaSP)
	}

	// Rotatable bonds
	rotBonds, err := m.NumRotatableBonds()
	if err == nil {
		fmt.Printf("Rotatable bonds: %d\n", rotBonds)
	}

	// Rings
	ringCount, err := m.CountSSSR()
	if err == nil {
		fmt.Printf("Rings (SSSR): %d\n", ringCount)
	}

	// Components
	components, err := m.CountComponents()
	if err == nil {
		fmt.Printf("Components: %d\n", components)
	}
}

func demonstrateProperties() {
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	m, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatalf("Failed to load molecule: %v", err)
	}
	defer m.Close()

	// Set molecule name
	fmt.Println("\n1. Setting molecule name:")
	err = m.SetName("Ethanol")
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Println("   Name set successfully")
	}

	name, err := m.Name()
	if err == nil {
		fmt.Printf("   Current name: %s\n", name)
	}

	// Set custom property
	fmt.Println("\n2. Setting custom property:")
	err = m.SetProperty("CAS", "64-17-5")
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Println("   Property 'CAS' set successfully")
	}

	// Check if property exists
	fmt.Println("\n3. Checking property:")
	has, err := m.HasProperty("CAS")
	if err == nil {
		fmt.Printf("   Has 'CAS' property: %v\n", has)
	}

	// Get property value
	cas, err := m.GetProperty("CAS")
	if err == nil {
		fmt.Printf("   CAS number: %s\n", cas)
	}

	// Set another property
	m.SetProperty("Description", "Common alcohol")
	m.SetProperty("Boiling Point", "78.37°C")

	// Display all set properties
	fmt.Println("\n4. All custom properties:")
	properties := []string{"CAS", "Description", "Boiling Point"}
	for _, prop := range properties {
		if has, _ := m.HasProperty(prop); has {
			val, _ := m.GetProperty(prop)
			fmt.Printf("   %s: %s\n", prop, val)
		}
	}

	// Remove a property
	fmt.Println("\n5. Removing property:")
	err = m.RemoveProperty("Description")
	if err == nil {
		fmt.Println("   Property 'Description' removed")
	}

	has, _ = m.HasProperty("Description")
	fmt.Printf("   Has 'Description' property: %v\n", has)
}
