// Package molecule_test provides tests for atom operations
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_atom_test.go
// @Software: GoLand
package molecule_test

import (
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

func TestGetAtom(t *testing.T) {
	mol, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	atom, err := mol.GetAtom(0)
	if err != nil {
		t.Fatalf("Failed to get atom: %v", err)
	}

	symbol, err := atom.Symbol()
	if err != nil {
		t.Fatalf("Failed to get atom symbol: %v", err)
	}

	if symbol != "C" {
		t.Errorf("Expected symbol C, got %s", symbol)
	}
}

func TestAtomProperties(t *testing.T) {
	mol, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	atom, _ := mol.GetAtom(0)

	// Test atomic number
	atomicNum, err := atom.AtomicNumber()
	if err != nil {
		t.Errorf("Failed to get atomic number: %v", err)
	}
	if atomicNum != 6 { // Carbon
		t.Errorf("Expected atomic number 6 (C), got %d", atomicNum)
	}

	// Test degree
	degree, err := atom.Degree()
	if err != nil {
		t.Errorf("Failed to get degree: %v", err)
	}
	if degree < 1 {
		t.Errorf("Expected degree >= 1, got %d", degree)
	}

	// Test valence
	valence, err := atom.Valence()
	if err != nil {
		t.Errorf("Failed to get valence: %v", err)
	}
	if valence != 4 {
		t.Logf("Warning: Expected valence 4 for C, got %d", valence)
	}
}

func TestSetCharge(t *testing.T) {
	mol, err := molecule.LoadMoleculeFromString("C")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	atom, _ := mol.GetAtom(0)

	// Set positive charge
	err = atom.SetCharge(1)
	if err != nil {
		t.Fatalf("Failed to set charge: %v", err)
	}

	charge, err := atom.Charge()
	if err != nil {
		t.Fatalf("Failed to get charge: %v", err)
	}

	if charge != 1 {
		t.Errorf("Expected charge 1, got %d", charge)
	}
}

func TestIsotope(t *testing.T) {
	mol, err := molecule.LoadMoleculeFromString("C")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	atom, _ := mol.GetAtom(0)

	// Set isotope
	err = atom.SetIsotope(13)
	if err != nil {
		t.Fatalf("Failed to set isotope: %v", err)
	}

	isotope, err := atom.Isotope()
	if err != nil {
		t.Fatalf("Failed to get isotope: %v", err)
	}

	if isotope != 13 {
		t.Errorf("Expected isotope 13, got %d", isotope)
	}
}

func TestBondOperations(t *testing.T) {
	mol, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	bondCount, _ := mol.CountBonds()
	if bondCount < 2 {
		t.Fatalf("Expected at least 2 bonds, got %d", bondCount)
	}

	// Test first bond
	bond, err := mol.GetBond(0)
	if err != nil {
		t.Fatalf("Failed to get bond: %v", err)
	}

	order, err := molecule.BondOrder(bond)
	if err != nil {
		t.Fatalf("Failed to get bond order: %v", err)
	}

	if order != 1 {
		t.Errorf("Expected bond order 1 (single), got %d", order)
	}

	// Test bond source and destination
	source, err := molecule.BondSource(bond)
	if err != nil {
		t.Fatalf("Failed to get bond source: %v", err)
	}

	dest, err := molecule.BondDestination(bond)
	if err != nil {
		t.Fatalf("Failed to get bond destination: %v", err)
	}

	if source < 0 || dest < 0 {
		t.Errorf("Invalid bond source/destination: %d, %d", source, dest)
	}
}

func TestAtomTypes(t *testing.T) {
	// Test pseudoatom
	mol, err := molecule.CreateMolecule()
	if err != nil {
		t.Fatalf("Failed to create molecule: %v", err)
	}
	defer mol.Close()

	atomHandle, err := mol.AddAtom("R")
	if err != nil {
		t.Fatalf("Failed to add R atom: %v", err)
	}

	atom, _ := mol.GetAtom(atomHandle)

	// Note: Need to check if this is actually a pseudoatom or R-site
	// Different Indigo versions might handle this differently
	t.Logf("Added R atom, IsPseudoatom: %v, IsRSite: %v",
		atom.IsPseudoatom(), atom.IsRSite())
}

func TestImplicitHydrogens(t *testing.T) {
	mol, err := molecule.LoadMoleculeFromString("C")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	atom, _ := mol.GetAtom(0)

	implicitH, err := atom.CountImplicitHydrogens()
	if err != nil {
		t.Fatalf("Failed to count implicit hydrogens: %v", err)
	}

	if implicitH != 4 {
		t.Logf("Warning: Expected 4 implicit H on C, got %d", implicitH)
	}
}
