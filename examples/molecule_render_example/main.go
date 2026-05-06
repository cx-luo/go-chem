// Package main demonstrates basic molecule parsing and 2D rendering.
//
// Run from the repository root:
//
//	go run ./examples/molecule_render_example
//
// The example writes SVG and PNG files to examples/output/.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/render2d"
)

const defaultPubChemCSV = "test/data/pubchem_smiles_036500001_037000000.csv"

type exampleMolecule struct {
	CID    int
	Name   string
	SMILES string
}

func main() {
	fmt.Println("=== Molecule and Render2D Examples ===")

	dataPath := flag.String("data", defaultPubChemCSV, "CSV file with PubChem cid,smiles columns")
	sampleSize := flag.Int("sample", 20, "number of PubChem SMILES to render")
	seed := flag.Int64("seed", time.Now().UnixNano(), "random seed for reproducible sampling")
	flag.Parse()

	outputDir := filepath.Join("examples", "output", "imgs")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("create output directory: %v", err)
	}

	records, err := loadPubChemSMILES(*dataPath)
	if err != nil {
		log.Fatalf("load PubChem SMILES CSV: %v", err)
	}
	rng := rand.New(rand.NewSource(*seed))
	rng.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})

	fmt.Printf("Data: %s\n", filepath.Clean(*dataPath))
	fmt.Printf("Random seed: %d\n", *seed)
	fmt.Printf("Requested sample size: %d\n", *sampleSize)

	rendered := 0
	skipped := 0
	for _, item := range records {
		if rendered >= *sampleSize {
			break
		}

		mol, err := parseMolecule(item.SMILES)
		if err != nil {
			skipped++
			log.Printf("skip CID %d: parse SMILES: %v", item.CID, err)
			continue
		}
		printMoleculeSummary(item, mol)
		renderMolecule(outputDir, item, mol)
		rendered++
	}

	if rendered < *sampleSize {
		log.Printf("rendered %d of %d requested molecules; skipped %d unparsable SMILES", rendered, *sampleSize, skipped)
	} else if skipped > 0 {
		log.Printf("skipped %d unparsable SMILES while sampling", skipped)
	}

	fmt.Printf("\nRendered files are in %s\n", outputDir)
}

func loadPubChemSMILES(path string) ([]exampleMolecule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(header) < 2 || strings.TrimSpace(header[0]) != "cid" || strings.TrimSpace(header[1]) != "smiles" {
		return nil, fmt.Errorf("expected header cid,smiles, got %q", strings.Join(header, ","))
	}

	var out []exampleMolecule
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(row) < 2 {
			continue
		}
		cid, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("parse CID %q: %w", row[0], err)
		}
		smiles := strings.TrimSpace(row[1])
		if smiles == "" {
			continue
		}
		out = append(out, exampleMolecule{
			CID:    cid,
			Name:   fmt.Sprintf("cid_%d", cid),
			SMILES: smiles,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no PubChem SMILES records found in %s", path)
	}
	return out, nil
}

func parseMolecule(smiles string) (*molecule.Molecule, error) {
	mol, err := (molecule.SmilesLoader{}).Parse(smiles)
	if err != nil {
		return nil, err
	}
	return mol, nil
}

func printMoleculeSummary(item exampleMolecule, mol *molecule.Molecule) {
	formula := molecule.GrossUnitsToStringHill(
		molecule.CollectGross(mol, molecule.GrossFormulaOptions{}),
		false,
	)

	fmt.Printf("\nMolecule: %s\n", item.Name)
	fmt.Printf("  CID:              %d\n", item.CID)
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
