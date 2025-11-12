// Package main demonstrates building molecules from scratch
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_builder.go
// @Software: GoLand
package main

import (
	"fmt"
	"github.com/cx-luo/go-chem/core"
	"log"

	"github.com/cx-luo/go-chem/molecule"
)

var indigoInit *core.Indigo

func init() {
	handle, err := core.IndigoInit()
	indigoInit = handle
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("=== Molecule Builder Examples ===\n")

	// Example 1: Build ethanol (C-C-O)
	fmt.Println("1. Building ethanol from scratch:")
	ethanol := buildEthanol()
	defer ethanol.Close()

	atomCount, _ := ethanol.CountAtoms()
	bondCount, _ := ethanol.CountBonds()
	fmt.Printf("   Built ethanol: %d atoms, %d bonds\n", atomCount, bondCount)

	smiles, _ := ethanol.ToSmiles()
	fmt.Printf("   SMILES: %s\n\n", smiles)

	// Example 2: Build benzene ring
	fmt.Println("2. Building benzene ring:")
	benzene := buildBenzene()
	defer benzene.Close()

	atomCount, _ = benzene.CountAtoms()
	bondCount, _ = benzene.CountBonds()
	ringCount, _ := benzene.CountSSSR()
	fmt.Printf("   Built benzene: %d atoms, %d bonds, %d rings\n", atomCount, bondCount, ringCount)

	smiles, _ = benzene.ToSmiles()
	fmt.Printf("   SMILES: %s\n\n", smiles)

	// Example 3: Build water molecule
	fmt.Println("3. Building water molecule:")
	water := buildWater()
	defer water.Close()

	atomCount, _ = water.CountAtoms()
	bondCount, _ = water.CountBonds()
	fmt.Printf("   Built water: %d atoms, %d bonds\n", atomCount, bondCount)

	smiles, _ = water.ToSmiles()
	fmt.Printf("   SMILES: %s\n\n", smiles)

	// Example 4: Set atom charge
	fmt.Println("4. Building charged molecule (acetate ion):")
	acetate := buildAcetate()
	defer acetate.Close()

	smiles, _ = acetate.ToSmiles()
	fmt.Printf("   Built acetate ion\n")
	fmt.Printf("   SMILES: %s\n\n", smiles)

	// Example 5: Merge molecules
	fmt.Println("5. Merging molecules:")
	mol1, _ := indigoInit.CreateMolecule()
	defer mol1.Close()
	c1, _ := mol1.AddAtom("C")
	c2, _ := mol1.AddAtom("C")
	mol1.AddBond(c1, c2, molecule.BOND_SINGLE)

	mol2, _ := indigoInit.CreateMolecule()
	defer mol2.Close()
	mol2.AddAtom("O")

	err := mol1.Merge(mol2)
	if err != nil {
		log.Printf("   Merge error: %v\n", err)
	} else {
		atomCount, _ := mol1.CountAtoms()
		components, _ := mol1.CountComponents()
		fmt.Printf("   Merged: %d atoms, %d components\n\n", atomCount, components)
	}

	fmt.Println("=== Builder Examples completed successfully ===")
}

// buildEthanol builds ethanol (C-C-O) from scratch
func buildEthanol() *molecule.Molecule {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		log.Fatalf("Failed to create molecule: %v", err)
	}

	// Add atoms
	c1, _ := m.AddAtom("C")
	c2, _ := m.AddAtom("C")
	o, _ := m.AddAtom("O")

	// Add bonds
	m.AddBond(c1, c2, molecule.BOND_SINGLE)
	m.AddBond(c2, o, molecule.BOND_SINGLE)

	return m
}

// buildBenzene builds a benzene ring from scratch
func buildBenzene() *molecule.Molecule {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		log.Fatalf("Failed to create molecule: %v", err)
	}

	// Add 6 carbon atoms
	var atoms [6]int
	for i := 0; i < 6; i++ {
		atoms[i], _ = m.AddAtom("C")
	}

	// Add aromatic bonds in a ring
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		m.AddBond(atoms[i], atoms[next], molecule.BOND_AROMATIC)
	}

	return m
}

// buildWater builds a water molecule
func buildWater() *molecule.Molecule {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		log.Fatalf("Failed to create molecule: %v", err)
	}

	// Add atoms
	o, _ := m.AddAtom("O")
	h1, _ := m.AddAtom("H")
	h2, _ := m.AddAtom("H")

	// Add bonds
	m.AddBond(o, h1, molecule.BOND_SINGLE)
	m.AddBond(o, h2, molecule.BOND_SINGLE)

	return m
}

// buildAcetate builds acetate ion (CH3COO-)
func buildAcetate() *molecule.Molecule {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		log.Fatalf("Failed to create molecule: %v", err)
	}

	// Add atoms
	c1, _ := m.AddAtom("C") // CH3
	c2, _ := m.AddAtom("C") // COO-
	o1, _ := m.AddAtom("O")
	o2, _ := m.AddAtom("O")

	// Add bonds
	m.AddBond(c1, c2, molecule.BOND_SINGLE)
	m.AddBond(c2, o1, molecule.BOND_DOUBLE)
	m.AddBond(c2, o2, molecule.BOND_SINGLE)

	// Set charge on O2
	molecule.SetCharge(o2, -1)

	return m
}
