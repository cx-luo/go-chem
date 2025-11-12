// Package main demonstrates atom-level operations
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : atom_operations.go
// @Software: GoLand
package main

import (
	"fmt"
	"github.com/cx-luo/go-chem/core"
	"log"

	"github.com/cx-luo/go-chem/molecule"
)

func main() {
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Atom-Level Operations Examples ===\n")

	// Load a molecule (ethanol)
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Get molecule info
	atomCount, _ := mol.CountAtoms()
	bondCount, _ := mol.CountBonds()
	fmt.Printf("Molecule: CCO (Ethanol)\n")
	fmt.Printf("Atoms: %d, Bonds: %d\n\n", atomCount, bondCount)

	// Example 1: Iterate and query atoms
	fmt.Println("1. Atom Properties:")
	for i := 0; i < atomCount; i++ {
		atom, err := mol.GetAtom(i)
		if err != nil {
			log.Printf("Failed to get atom %d: %v", i, err)
			continue
		}

		symbol, _ := atom.Symbol()
		atomicNum, _ := atom.AtomicNumber()
		degree, _ := atom.Degree()
		valence, _ := atom.Valence()
		charge, _ := atom.Charge()
		implicitH, _ := atom.CountImplicitHydrogens()

		fmt.Printf("  Atom %d: %s (Z=%d)\n", i, symbol, atomicNum)
		fmt.Printf("    Degree: %d, Valence: %d, Charge: %d\n", degree, valence, charge)
		fmt.Printf("    Implicit H: %d\n", implicitH)
	}

	// Example 2: Modify atom properties
	fmt.Println("\n2. Modifying Atom Properties:")
	atom1, _ := mol.GetAtom(1)

	// Set charge
	fmt.Println("  Setting positive charge on carbon...")
	atom1.SetCharge(1)
	charge, _ := atom1.Charge()
	fmt.Printf("  New charge: %d\n", charge)

	// Reset charge
	atom1.SetCharge(0)
	fmt.Println("  Charge reset to 0")

	// Example 3: Isotope operations
	fmt.Println("\n3. Isotope Operations:")
	atom0, _ := mol.GetAtom(0)

	fmt.Println("  Setting carbon-13 isotope...")
	atom0.SetIsotope(13)
	isotope, _ := atom0.Isotope()
	fmt.Printf("  Isotope: %d\n", isotope)

	// Example 4: Build molecule from scratch
	fmt.Println("\n4. Building Molecule from Atoms:")
	newMol, err := molecule.CreateMolecule()
	if err != nil {
		log.Fatalf("Failed to create molecule: %v", err)
	}
	defer newMol.Close()

	// Add atoms
	c1, _ := newMol.AddAtom("C")
	c2, _ := newMol.AddAtom("C")
	o, _ := newMol.AddAtom("O")

	fmt.Println("  Added: C, C, O atoms")

	// Add bonds
	newMol.AddBond(c1, c2, molecule.BOND_SINGLE)
	newMol.AddBond(c2, o, molecule.BOND_SINGLE)

	fmt.Println("  Added bonds: C-C, C-O")

	// Get SMILES
	smiles, _ := newMol.ToSmiles()
	fmt.Printf("  Result SMILES: %s\n", smiles)

	// Example 5: Query bond properties
	fmt.Println("\n5. Bond Properties:")
	for i := 0; i < bondCount; i++ {
		bond, err := mol.GetBond(i)
		if err != nil {
			continue
		}

		order, _ := molecule.BondOrder(bond)
		index, _ := molecule.BondIndex(bond)
		source, _ := molecule.BondSource(bond)
		dest, _ := molecule.BondDestination(bond)

		srcAtom, _ := mol.GetAtom(source)
		dstAtom, _ := mol.GetAtom(dest)
		srcSymbol, _ := srcAtom.Symbol()
		dstSymbol, _ := dstAtom.Symbol()

		fmt.Printf("  Bond %d: %s-%s (order: %d)\n", index, srcSymbol, dstSymbol, order)
	}

	fmt.Println("\n=== Examples completed successfully ===")
}
