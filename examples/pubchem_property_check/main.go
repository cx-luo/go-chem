// Package main checks calculated molecule properties against PubChem XML data.
//
// Run from the repository root:
//
//	go run ./examples/pubchem_property_check -limit 100
//
// The default XML path is pkg/Indigo/Compound_036500001_037000000.xml.
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cx-luo/go-chem/molecule"
)

const defaultXMLPath = "pkg/Indigo/Compound_036500001_037000000.xml"

type pcCompound struct {
	CID   int          `xml:"PC-Compound_id>PC-CompoundType>PC-CompoundType_id>PC-CompoundType_id_cid"`
	Props []pcInfoData `xml:"PC-Compound_props>PC-InfoData"`
}

type pcInfoData struct {
	URN struct {
		Label string `xml:"PC-Urn_label"`
		Name  string `xml:"PC-Urn_name"`
	} `xml:"PC-InfoData_urn>PC-Urn"`
	Value pcInfoValue `xml:"PC-InfoData_value"`
}

type pcInfoValue struct {
	String string `xml:"PC-InfoData_value_sval"`
	Float  string `xml:"PC-InfoData_value_fval"`
	Int    string `xml:"PC-InfoData_value_ival"`
}

func (v pcInfoValue) text() string {
	switch {
	case v.String != "":
		return strings.TrimSpace(v.String)
	case v.Float != "":
		return strings.TrimSpace(v.Float)
	default:
		return strings.TrimSpace(v.Int)
	}
}

type expectedProperties struct {
	CID       int
	SMILES    string
	Formula   string
	Weight    *float64
	TPSA      *float64
	HBA       *int
	HBD       *int
	Rotatable *int
}

type mismatch struct {
	CID      int
	SMILES   string
	Property string
	Got      string
	Want     string
}

type propertyComparison struct {
	Property string
	Checked  bool
	Matched  bool
	Got      string
	Want     string
}

type propertyStats struct {
	Checked    int
	Matched    int
	Mismatched int
	Skipped    int
}

type summary struct {
	Seen          int
	Checked       int
	Matched       int
	Mismatched    int
	Skipped       int
	PropertyStats map[string]propertyStats
	Mismatches    []mismatch
}

var propertyOrder = []string{
	"Molecular Formula",
	"Molecular Weight",
	"TPSA",
	"HBA",
	"HBD",
	"Rotatable Bond Count",
}

func main() {
	xmlPath := flag.String("xml", defaultXMLPath, "PubChem compound XML file")
	limit := flag.Int("limit", 100, "number of compounds to check; 0 checks the whole file")
	maxMismatches := flag.Int("max-mismatches", 20, "maximum mismatches to print")
	weightTolerance := flag.Float64("weight-tolerance", 0.2, "absolute tolerance for molecular weight")
	tpsaTolerance := flag.Float64("tpsa-tolerance", 0.5, "absolute tolerance for TPSA")
	strict := flag.Bool("strict", true, "exit with a non-zero status when mismatches are found")
	flag.Parse()

	result, err := checkProperties(*xmlPath, *limit, *maxMismatches, *weightTolerance, *tpsaTolerance)
	if err != nil {
		log.Fatalf("check PubChem properties: %v", err)
	}

	fmt.Printf("PubChem property check: %s\n", filepath.Clean(*xmlPath))
	fmt.Printf("  Compounds seen:     %d\n", result.Seen)
	fmt.Printf("  Compounds checked:  %d\n", result.Checked)
	fmt.Printf("  Fully matched:      %d\n", result.Matched)
	fmt.Printf("  With mismatches:    %d\n", result.Mismatched)
	fmt.Printf("  Skipped:            %d\n", result.Skipped)

	fmt.Println("\nProperty comparisons:")
	for _, property := range propertyOrder {
		stats := result.PropertyStats[property]
		fmt.Printf("  %-22s checked: %d, matched: %d, mismatched: %d, skipped: %d\n",
			property, stats.Checked, stats.Matched, stats.Mismatched, stats.Skipped)
	}

	if len(result.Mismatches) > 0 {
		fmt.Printf("\nFirst %d mismatches:\n", len(result.Mismatches))
		for _, item := range result.Mismatches {
			fmt.Printf("  CID %d %s\n", item.CID, item.SMILES)
			fmt.Printf("    %s: got %s, want %s\n", item.Property, item.Got, item.Want)
		}
	}

	if *strict && result.Mismatched > 0 {
		os.Exit(1)
	}
}

func checkProperties(xmlPath string, limit, maxMismatches int, weightTolerance, tpsaTolerance float64) (summary, error) {
	file, err := os.Open(xmlPath)
	if err != nil {
		return summary{}, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	result := newSummary()
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}

		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "PC-Compound" {
			continue
		}

		var compound pcCompound
		if err := decoder.DecodeElement(&compound, &start); err != nil {
			return result, err
		}

		result.Seen++
		expected := extractExpectedProperties(compound)
		if expected.SMILES == "" {
			result.Skipped++
			continue
		}

		mol, err := (molecule.SmilesLoader{}).Parse(expected.SMILES)
		if err != nil {
			result.Skipped++
			if len(result.Mismatches) < maxMismatches {
				result.Mismatches = append(result.Mismatches, mismatch{
					CID:      expected.CID,
					SMILES:   expected.SMILES,
					Property: "SMILES parse",
					Got:      err.Error(),
					Want:     "parse without error",
				})
			}
			continue
		}

		result.Checked++
		comparisons := compareProperties(expected, mol, weightTolerance, tpsaTolerance)
		hasMismatch := false
		for _, comparison := range comparisons {
			if result.recordPropertyComparison(expected, comparison, maxMismatches) {
				hasMismatch = true
			}
		}
		if !hasMismatch {
			result.Matched++
		} else {
			result.Mismatched++
		}

		if limit > 0 && result.Seen >= limit {
			break
		}
	}

	return result, nil
}

func newSummary() summary {
	stats := make(map[string]propertyStats, len(propertyOrder))
	for _, property := range propertyOrder {
		stats[property] = propertyStats{}
	}
	return summary{PropertyStats: stats}
}

func (s *summary) recordPropertyComparison(expected expectedProperties, comparison propertyComparison, maxMismatches int) bool {
	stats := s.PropertyStats[comparison.Property]
	if comparison.Checked {
		stats.Checked++
		if comparison.Matched {
			stats.Matched++
		} else {
			stats.Mismatched++
		}
	} else {
		stats.Skipped++
	}
	s.PropertyStats[comparison.Property] = stats

	if !comparison.Checked || comparison.Matched {
		return false
	}
	if len(s.Mismatches) < maxMismatches {
		s.Mismatches = append(s.Mismatches, mismatch{
			CID:      expected.CID,
			SMILES:   expected.SMILES,
			Property: comparison.Property,
			Got:      comparison.Got,
			Want:     comparison.Want,
		})
	}
	return true
}

func extractExpectedProperties(compound pcCompound) expectedProperties {
	expected := expectedProperties{CID: compound.CID}
	for _, prop := range compound.Props {
		label := strings.TrimSpace(prop.URN.Label)
		name := strings.TrimSpace(prop.URN.Name)
		value := prop.Value.text()
		if value == "" {
			continue
		}

		switch {
		case label == "SMILES" && name == "Connectivity":
			expected.SMILES = value
		case label == "SMILES" && name == "Absolute" && expected.SMILES == "":
			expected.SMILES = value
		case label == "Molecular Formula":
			expected.Formula = value
		case label == "Molecular Weight":
			expected.Weight = parseFloat(value)
		case label == "Topological" && name == "Polar Surface Area":
			expected.TPSA = parseFloat(value)
		case label == "Count" && name == "Hydrogen Bond Acceptor":
			expected.HBA = parseInt(value)
		case label == "Count" && name == "Hydrogen Bond Donor":
			expected.HBD = parseInt(value)
		case label == "Count" && name == "Rotatable Bond":
			expected.Rotatable = parseInt(value)
		}
	}
	return expected
}

func compareProperties(expected expectedProperties, mol *molecule.Molecule, weightTolerance, tpsaTolerance float64) []propertyComparison {
	var comparisons []propertyComparison
	skip := func(property string) {
		comparisons = append(comparisons, propertyComparison{
			Property: property,
		})
	}
	add := func(property, got, want string, matched bool) {
		comparisons = append(comparisons, propertyComparison{
			Property: property,
			Checked:  true,
			Matched:  matched,
			Got:      got,
			Want:     want,
		})
	}

	if expected.Formula != "" {
		got := molecule.GrossUnitsToStringHill(
			molecule.CollectGross(mol, molecule.GrossFormulaOptions{}),
			false,
		)
		got = compactFormula(got)
		want := compactFormula(expected.Formula)
		add("Molecular Formula", got, want, got == want)
	} else {
		skip("Molecular Formula")
	}

	if expected.Weight != nil {
		got := mol.CalcMolecularWeight()
		matched := math.Abs(got-*expected.Weight) <= weightTolerance
		add("Molecular Weight", formatFloat(got), formatFloat(*expected.Weight), matched)
	} else {
		skip("Molecular Weight")
	}

	if expected.TPSA != nil {
		got := mol.CalculateTPSA(true)
		matched := math.Abs(got-*expected.TPSA) <= tpsaTolerance
		add("TPSA", formatFloat(got), formatFloat(*expected.TPSA), matched)
	} else {
		skip("TPSA")
	}

	if expected.HBA != nil {
		got := molecule.NumHydrogenBondAcceptors(mol)
		add("HBA", strconv.Itoa(got), strconv.Itoa(*expected.HBA), got == *expected.HBA)
	} else {
		skip("HBA")
	}

	if expected.HBD != nil {
		got := molecule.NumHydrogenBondDonors(mol)
		add("HBD", strconv.Itoa(got), strconv.Itoa(*expected.HBD), got == *expected.HBD)
	} else {
		skip("HBD")
	}

	if expected.Rotatable != nil {
		got := molecule.NumRotatableBonds(mol)
		add("Rotatable Bond Count", strconv.Itoa(got), strconv.Itoa(*expected.Rotatable), got == *expected.Rotatable)
	} else {
		skip("Rotatable Bond Count")
	}

	return comparisons
}

func parseFloat(value string) *float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseInt(value string) *int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &parsed
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
}

func compactFormula(value string) string {
	return strings.Join(strings.Fields(value), "")
}
