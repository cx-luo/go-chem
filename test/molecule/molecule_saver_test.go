// Package molecule_test provides tests for enhanced molecule saving functionality
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_saver_enhanced_test.go
// @Software: GoLand
package molecule_test

import (
	"os"
	"strings"
	"testing"
)

// Test ChemAxon CXSMILES format
func TestToCXSmiles(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	cxsmiles, err := mol.ToCXSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to CXSmiles: %v", err)
	}

	if cxsmiles == "" {
		t.Error("CXSmiles is empty")
	}

	t.Logf("CXSmiles: %s", cxsmiles)
}

// Test Canonical ChemAxon CXSMILES format
func TestToCanonicalCXSmiles(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	cxsmiles, err := mol.ToCanonicalCXSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to canonical CXSmiles: %v", err)
	}

	if cxsmiles == "" {
		t.Error("Canonical CXSmiles is empty")
	}

	t.Logf("Canonical CXSmiles: %s", cxsmiles)
}

// Test Daylight SMILES format
func TestToDaylightSmiles(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	smiles, err := mol.ToDaylightSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to Daylight SMILES: %v", err)
	}

	if smiles == "" {
		t.Error("Daylight SMILES is empty")
	}

	t.Logf("Daylight SMILES: %s", smiles)
}

// Test CML file saving
func TestSaveToCMLFile(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	filename := "test_ethanol.cml"
	defer os.Remove(filename)

	err = mol.SaveToCMLFile(filename)
	if err != nil {
		t.Fatalf("Failed to save to CML file: %v", err)
	}

	// Check if file exists and has content
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("CML file is empty")
	}

	t.Logf("CML file size: %d bytes", info.Size())
}

// Test CDXML file saving
func TestSaveToCDXMLFile(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	filename := "test_ethanol.cdxml"
	defer os.Remove(filename)

	err = mol.SaveToCDXMLFile(filename)
	if err != nil {
		t.Fatalf("Failed to save to CDXML file: %v", err)
	}

	// Check if file exists and has content
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("CDXML file is empty")
	}

	t.Logf("CDXML file size: %d bytes", info.Size())
}

// Test CDX file saving
func TestSaveToCDXFile(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	filename := "test_ethanol.cdx"
	defer os.Remove(filename)

	err = mol.SaveToCDXFile(filename)
	if err != nil {
		t.Fatalf("Failed to save to CDX file: %v", err)
	}

	// Check if file exists and has content
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("CDX file is empty")
	}

	t.Logf("CDX file size: %d bytes", info.Size())
}

// Test RDF format conversion
func TestToRDF(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	rdf, err := mol.ToRDF()
	if err != nil {
		t.Fatalf("Failed to convert to RDF: %v", err)
	}

	if rdf == "" {
		t.Error("RDF is empty")
	}

	if !strings.Contains(rdf, "$RDFILE") && !strings.Contains(rdf, "$RXN") {
		t.Logf("Warning: RDF format may not be standard (first 100 chars): %s", rdf[:min(100, len(rdf))])
	}

	t.Logf("RDF length: %d bytes", len(rdf))
}

// Test RDF file saving
func TestSaveToRDFFile(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	filename := "test_ethanol.rdf"
	defer os.Remove(filename)

	err = mol.SaveToRDFFile(filename)
	if err != nil {
		t.Fatalf("Failed to save to RDF file: %v", err)
	}

	// Check if file exists and has content
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("RDF file is empty")
	}

	t.Logf("RDF file size: %d bytes", info.Size())
}

// Test buffer conversion
func TestToBuffer(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	buffer, err := mol.ToBuffer()
	if err != nil {
		t.Fatalf("Failed to convert to buffer: %v", err)
	}

	if len(buffer) == 0 {
		t.Error("Buffer is empty")
	}

	// Buffer should contain MOL format data
	bufferStr := string(buffer)
	if !strings.Contains(bufferStr, "V2000") && !strings.Contains(bufferStr, "V3000") {
		t.Errorf("Buffer doesn't seem to contain valid MOL format")
	}

	t.Logf("Buffer size: %d bytes", len(buffer))
}

// Test KET format conversion
func TestToKet(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	ket, err := mol.ToKet()
	if err != nil {
		t.Fatalf("Failed to convert to KET: %v", err)
	}

	if ket == "" {
		t.Error("KET is empty")
	}

	// KET should contain JSON structure
	if !strings.Contains(ket, "{") || !strings.Contains(ket, "}") {
		t.Error("KET doesn't contain valid JSON structure")
	}

	t.Logf("KET length: %d bytes", len(ket))
}

// Test KET file saving
func TestSaveToKetFile(t *testing.T) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	filename := "test_ethanol.ket"
	defer os.Remove(filename)

	err = mol.SaveToKetFile(filename)
	if err != nil {
		t.Fatalf("Failed to save to KET file: %v", err)
	}

	// Check if file exists and has content
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("KET file is empty")
	}

	t.Logf("KET file size: %d bytes", info.Size())
}

// Test CXSmiles with complex molecule
func TestCXSmilesWithComplexMolecule(t *testing.T) {
	// Aspirin with stereochemistry
	smilesInput := "CC(=O)Oc1ccccc1C(=O)O"

	mol, err := indigoInit.LoadMoleculeFromString(smilesInput)
	if err != nil {
		t.Fatalf("Failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Test regular SMILES
	smiles, err := mol.ToSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to SMILES: %v", err)
	}
	t.Logf("SMILES: %s", smiles)

	// Test canonical SMILES
	canonicalSmiles, err := mol.ToCanonicalSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to canonical SMILES: %v", err)
	}
	t.Logf("Canonical SMILES: %s", canonicalSmiles)

	// Test CXSmiles
	cxsmiles, err := mol.ToCXSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to CXSmiles: %v", err)
	}
	t.Logf("CXSmiles: %s", cxsmiles)

	// Test canonical CXSmiles
	canonicalCXSmiles, err := mol.ToCanonicalCXSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to canonical CXSmiles: %v", err)
	}
	t.Logf("Canonical CXSmiles: %s", canonicalCXSmiles)

	// Test Daylight SMILES
	daylightSmiles, err := mol.ToDaylightSmiles()
	if err != nil {
		t.Fatalf("Failed to convert to Daylight SMILES: %v", err)
	}
	t.Logf("Daylight SMILES: %s", daylightSmiles)
}

// Test format consistency
func TestFormatConsistency(t *testing.T) {
	testCases := []string{
		"CCO",      // Ethanol
		"c1ccccc1", // Benzene
		"CC(C)C",   // Isobutane
		"CC(=O)O",  // Acetic acid
		"C1CCCCC1", // Cyclohexane
	}

	for _, smiles := range testCases {
		t.Run(smiles, func(t *testing.T) {
			mol, err := indigoInit.LoadMoleculeFromString(smiles)
			if err != nil {
				t.Fatalf("Failed to load molecule %s: %v", smiles, err)
			}
			defer mol.Close()

			// Test all SMILES variants
			formats := map[string]func() (string, error){
				"SMILES":       mol.ToSmiles,
				"Canonical":    mol.ToCanonicalSmiles,
				"CXSmiles":     mol.ToCXSmiles,
				"Canonical CX": mol.ToCanonicalCXSmiles,
				"Daylight":     mol.ToDaylightSmiles,
			}

			for name, fn := range formats {
				result, err := fn()
				if err != nil {
					t.Errorf("Failed to convert to %s: %v", name, err)
					continue
				}
				if result == "" {
					t.Errorf("%s is empty", name)
				}
			}
		})
	}
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
