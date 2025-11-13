// Package molecule_test provides tests for substructure matching
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_match_test.go
// @Software: GoLand
package molecule_test

import (
	"testing"
)

func TestHasSubstructure(t *testing.T) {
	target, err := indigoInit.LoadMoleculeFromString("c1ccccc1CCO")
	if err != nil {
		t.Fatalf("Failed to load target: %v", err)
	}
	defer target.Close()

	query, err := indigoInit.LoadQueryMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("Failed to load query: %v", err)
	}
	defer query.Close()

	has, err := target.HasSubstructure(query, nil)
	if err != nil {
		t.Fatalf("HasSubstructure failed: %v", err)
	}

	if !has {
		t.Error("Expected to find benzene ring in phenylethanol")
	}
}

func TestCountSubstructureMatches(t *testing.T) {
	// Biphenyl has 2 benzene rings
	target, err := indigoInit.LoadMoleculeFromString("c1ccc(cc1)c2ccccc2")
	if err != nil {
		t.Fatalf("Failed to load target: %v", err)
	}
	defer target.Close()

	query, err := indigoInit.LoadQueryMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("Failed to load query: %v", err)
	}
	defer query.Close()

	// set timeout to 10 seconds to avoid infinite loop
	count, err := target.CountSubstructureMatches(query, nil)
	if err != nil {
		t.Fatalf("CountSubstructureMatches failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 benzene rings in biphenyl, got %d", count)
	}
}

func TestExactMatch(t *testing.T) {
	mol1, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load mol1: %v", err)
	}
	defer mol1.Close()

	mol2, err := indigoInit.LoadQueryMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load mol2: %v", err)
	}
	defer mol2.Close()

	mol3, err := indigoInit.LoadQueryMoleculeFromString("CC")
	if err != nil {
		t.Fatalf("Failed to load mol3: %v", err)
	}
	defer mol3.Close()

	// Test exact match
	isExact, err := mol1.ExactMatch(mol2)
	if err != nil {
		t.Fatalf("ExactMatch failed: %v", err)
	}

	if !isExact {
		t.Error("Expected CCO to exactly match CCO")
	}

	// Test non-match
	isExact2, err := mol1.ExactMatch(mol3)
	if err != nil {
		t.Fatalf("ExactMatch failed: %v", err)
	}

	if isExact2 {
		t.Error("Expected CCO not to match CC")
	}
}

func TestSMARTSMatching(t *testing.T) {
	// Test carboxylic acid pattern in acetic acid
	mol, err := indigoInit.LoadMoleculeFromString("CC(=O)O")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	pattern, err := indigoInit.LoadQueryMoleculeFromString("[CX3](=O)[OX2H1]")
	if err != nil {
		t.Fatalf("Failed to load SMARTS: %v", err)
	}
	defer pattern.Close()

	has, err := mol.HasSubstructure(pattern, nil)
	if err != nil {
		t.Fatalf("HasSubstructure failed: %v", err)
	}

	if !has {
		t.Error("Expected to find carboxylic acid group in acetic acid")
	}

	count, err := mol.CountSubstructureMatches(pattern, nil)
	if err != nil {
		t.Fatalf("CountSubstructureMatches failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 carboxylic acid group, got %d", count)
	}
}

func TestGetSubmolecule(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCCCCC")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Get first 3 atoms (propane from hexane)
	submol, err := mol.GetSubmolecule([]int{0, 1, 2})
	if err != nil {
		t.Fatalf("GetSubmolecule failed: %v", err)
	}
	defer submol.Close()

	atomCount, err := submol.CountAtoms()
	if err != nil {
		t.Fatalf("Failed to count atoms: %v", err)
	}

	if atomCount != 3 {
		t.Errorf("Expected 3 atoms in submolecule, got %d", atomCount)
	}
}

func TestRemoveAtoms(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCCC")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	initialCount, _ := mol.CountAtoms()

	// Remove one atom
	err = mol.RemoveAtoms([]int{0})
	if err != nil {
		t.Fatalf("RemoveAtoms failed: %v", err)
	}

	newCount, _ := mol.CountAtoms()

	if newCount != initialCount-1 {
		t.Errorf("Expected %d atoms after removal, got %d", initialCount-1, newCount)
	}
}

func TestHighlight(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Just test that unhighlight doesn't crash
	err = mol.UnhighlightAll()
	if err != nil {
		t.Errorf("UnhighlightAll failed: %v", err)
	}
}

func TestEmptySubmolecule(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCC")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Try to get submolecule with empty indices
	_, err = mol.GetSubmolecule([]int{})
	if err == nil {
		t.Error("Expected error for empty atom indices")
	}
}

func TestRemoveEmptyAtoms(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCC")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	initialCount, _ := mol.CountAtoms()

	// Remove empty list should not change anything
	err = mol.RemoveAtoms([]int{})
	if err != nil {
		t.Errorf("RemoveAtoms with empty list failed: %v", err)
	}

	newCount, _ := mol.CountAtoms()

	if newCount != initialCount {
		t.Errorf("Expected %d atoms, got %d", initialCount, newCount)
	}
}
