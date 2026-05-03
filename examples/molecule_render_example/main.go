// Package main demonstrates basic molecule parsing and 2D rendering.
//
// Run from the repository root:
//
//	go run ./examples/molecule_render_example
//
// The example writes SVG and PNG files to examples/output/.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/render2d"
)

type exampleMolecule struct {
	Name   string
	SMILES string
}

func main() {
	fmt.Println("=== Molecule and Render2D Examples ===")

	outputDir := filepath.Join("examples", "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("create output directory: %v", err)
	}

	examples := []exampleMolecule{
		{Name: "ethanol", SMILES: "CCO"},
		{Name: "benzene", SMILES: "c1ccccc1"},
		{Name: "phenol", SMILES: "c1ccccc1O"},
		{Name: "aspirin", SMILES: "CC(=O)OC1=CC=CC=C1C(=O)O"},
		{Name: "glucose", SMILES: "C1C(C(C(C(C(O1)O)O)O)O)O"},
		{Name: "2,3,6-triiodobenzaldehyde", SMILES: "C1=CC(=C(C(=C1I)C=O)I)I"},
		{Name: "benzaldehyde", SMILES: "C1=CC=C(C=C1)C=O"},
		{Name: "imidazole", SMILES: "CCCCN1C(=O)C=CC(=N1)C(=O)NC2=NC(=CS2)C3=CC=C(S3)C"},
	}

	for _, item := range examples {
		mol := mustParse(item.SMILES)
		printMoleculeSummary(item, mol)
		renderMolecule(outputDir, item, mol)
	}

	fmt.Printf("\nRendered files are in %s\n", outputDir)
}

func mustParse(smiles string) *molecule.Molecule {
	mol, err := (molecule.SmilesLoader{}).Parse(smiles)
	if err != nil {
		log.Fatalf("parse %q: %v", smiles, err)
	}
	return mol
}

func printMoleculeSummary(item exampleMolecule, mol *molecule.Molecule) {
	formula := molecule.GrossUnitsToStringHill(
		molecule.CollectGross(mol, molecule.GrossFormulaOptions{}),
		false,
	)

	fmt.Printf("\nMolecule: %s\n", item.Name)
	fmt.Printf("  SMILES:           %s\n", item.SMILES)
	fmt.Printf("  Atoms / Bonds:    %d / %d\n", mol.AtomCount(), mol.BondCount())
	fmt.Printf("  Formula:          %s\n", formula)
	fmt.Printf("  Molecular Weight: %.3f\n", mol.CalcMolecularWeight())
	fmt.Printf("  TPSA:             %.2f\n", mol.CalculateTPSA(true))
	fmt.Printf("  HBA / HBD:        %d / %d\n",
		molecule.NumHydrogenBondAcceptors(mol),
		molecule.NumHydrogenBondDonors(mol),
	)
	fmt.Printf("  Rotatable Bonds:  %d\n", molecule.NumRotatableBonds(mol))
}

func renderMolecule(outputDir string, item exampleMolecule, mol *molecule.Molecule) {
	options := render2d.Options{
		Width:            360,
		Height:           260,
		Margin:           32,
		ShowCarbonLabels: false,
		UseAtomColors:    true,
	}

	svgPath := filepath.Join(outputDir, item.Name+".svg")
	if err := render2d.SaveSVG(svgPath, mol, options); err != nil {
		log.Fatalf("render SVG for %s: %v", item.Name, err)
	}

	pngPath := filepath.Join(outputDir, item.Name+".png")
	if err := render2d.SavePNG(pngPath, mol, options); err != nil {
		log.Fatalf("render PNG for %s: %v", item.Name, err)
	}

	fmt.Printf("  SVG:              %s\n", svgPath)
	fmt.Printf("  PNG:              %s\n", pngPath)
}
