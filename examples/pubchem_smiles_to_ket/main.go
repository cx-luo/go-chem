// Package main generates Ketcher KET JSON files from PubChem SMILES records.
//
// Run from the repository root:
//
//	go run ./examples/pubchem_smiles_to_ket -limit 20
//
// The default XML path is pkg/Indigo/Compound_036500001_037000000.xml and
// generated files are written to examples/output.
//
// To additionally fetch reference layouts from a remote Indigo HTTP server,
// pass -spec-api <url>. The reference KET JSON returned by the server is
// written next to the algorithm output as cid_{cid}.spec.ket.json. Example:
//
//	go run ./examples/pubchem_smiles_to_ket \
//	    -xml pkg/Indigo/Compound_036500001_037000000.xml \
//	    -out examples/output \
//	    -limit 200 \
//	    -pretty=false \
//	    -spec-api http://192.168.10.128:8002/v2/indigo/layout
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/render2d"
)

const (
	defaultXMLPath = "pkg/Indigo/Compound_036500001_037000000.xml"
	defaultOutDir  = "examples/output"
)

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

type compoundSMILES struct {
	CID    int
	SMILES string
}

type generationSummary struct {
	Seen      int
	Written   int
	Skipped   int
	ParseFail int
	WriteFail int
	SpecOK    int
	SpecFail  int
}

type generationConfig struct {
	XMLPath       string
	OutDir        string
	Limit         int
	Pretty        bool
	SpecAPI       string
	SpecOnly      bool
	SpecOverwrite bool
	SpecTimeout   time.Duration
}

func main() {
	xmlPath := flag.String("xml", defaultXMLPath, "PubChem compound XML file")
	outDir := flag.String("out", defaultOutDir, "directory for generated KET JSON files")
	limit := flag.Int("limit", 20, "number of KET files to generate; 0 reads the whole file")
	pretty := flag.Bool("pretty", true, "write indented JSON")
	specAPI := flag.String("spec-api", "", "Indigo HTTP layout endpoint (e.g. http://host:8002/v2/indigo/layout); when set, also writes cid_{cid}.spec.ket.json reference layouts")
	specOnly := flag.Bool("spec-only", false, "only generate spec files via -spec-api; skip writing the algorithm KET file")
	specOverwrite := flag.Bool("spec-overwrite", false, "overwrite existing cid_{cid}.spec.ket.json files (default: skip if present)")
	specTimeout := flag.Duration("spec-timeout", 15*time.Second, "HTTP timeout for each -spec-api request")
	flag.Parse()

	cfg := generationConfig{
		XMLPath:       *xmlPath,
		OutDir:        *outDir,
		Limit:         *limit,
		Pretty:        *pretty,
		SpecAPI:       strings.TrimSpace(*specAPI),
		SpecOnly:      *specOnly,
		SpecOverwrite: *specOverwrite,
		SpecTimeout:   *specTimeout,
	}
	if cfg.SpecOnly && cfg.SpecAPI == "" {
		log.Fatalf("-spec-only requires -spec-api")
	}

	summary, err := generateKETFromPubChemXML(cfg)
	if err != nil {
		log.Fatalf("generate KET JSON: %v", err)
	}

	fmt.Printf("PubChem SMILES to KET: %s\n", filepath.Clean(cfg.XMLPath))
	fmt.Printf("  Output directory: %s\n", filepath.Clean(cfg.OutDir))
	fmt.Printf("  Compounds seen:   %d\n", summary.Seen)
	fmt.Printf("  Files written:    %d\n", summary.Written)
	fmt.Printf("  Skipped:          %d\n", summary.Skipped)
	fmt.Printf("  Parse failures:   %d\n", summary.ParseFail)
	fmt.Printf("  Write failures:   %d\n", summary.WriteFail)
	if cfg.SpecAPI != "" {
		fmt.Printf("  Spec API:         %s\n", cfg.SpecAPI)
		fmt.Printf("  Spec written:     %d\n", summary.SpecOK)
		fmt.Printf("  Spec failures:    %d\n", summary.SpecFail)
	}
}

func generateKETFromPubChemXML(cfg generationConfig) (generationSummary, error) {
	file, err := os.Open(cfg.XMLPath)
	if err != nil {
		return generationSummary{}, err
	}
	defer file.Close()

	if err := os.MkdirAll(cfg.OutDir, 0o755); err != nil {
		return generationSummary{}, err
	}

	var httpClient *http.Client
	if cfg.SpecAPI != "" {
		httpClient = &http.Client{Timeout: cfg.SpecTimeout}
	}

	decoder := xml.NewDecoder(file)
	result := generationSummary{}
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
		item := extractSMILES(compound)
		if item.SMILES == "" {
			result.Skipped++
			continue
		}

		wroteAlg := false
		if !cfg.SpecOnly {
			mol, err := (molecule.SmilesLoader{}).Parse(item.SMILES)
			if err != nil {
				result.ParseFail++
				continue
			}
			apply2DLayout(mol)

			ket, err := molecule.SaveMoleculeToKETWithOptions(mol, molecule.KETSaverOptions{
				PrettyJSON: cfg.Pretty,
				BondLength: 1.5,
			})
			if err != nil {
				result.WriteFail++
				continue
			}

			path := filepath.Join(cfg.OutDir, fmt.Sprintf("cid_%d.ket.json", item.CID))
			if err := os.WriteFile(path, []byte(ket+"\n"), 0o644); err != nil {
				result.WriteFail++
				continue
			}
			result.Written++
			wroteAlg = true
		}

		if cfg.SpecAPI != "" {
			specPath := filepath.Join(cfg.OutDir, fmt.Sprintf("cid_%d.spec.ket.json", item.CID))
			if cfg.SpecOverwrite || !fileExists(specPath) {
				if err := fetchAndWriteSpec(httpClient, cfg.SpecAPI, item.SMILES, specPath, cfg.Pretty); err != nil {
					result.SpecFail++
					log.Printf("spec API for CID %d: %v", item.CID, err)
				} else {
					result.SpecOK++
				}
			}
		}

		if cfg.SpecOnly && cfg.SpecAPI != "" {
			result.Written = result.SpecOK
		}
		_ = wroteAlg

		stop := false
		switch {
		case cfg.Limit <= 0:
			// no limit
		case cfg.SpecOnly:
			if result.SpecOK+result.SpecFail >= cfg.Limit {
				stop = true
			}
		default:
			if result.Written >= cfg.Limit {
				stop = true
			}
		}
		if stop {
			break
		}
	}

	return result, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// indigoLayoutResponse mirrors the JSON envelope returned by the
// /v2/indigo/layout endpoint. The "struct" field carries the actual KET
// document as a JSON-encoded string.
type indigoLayoutResponse struct {
	Format         string `json:"format"`
	OriginalFormat string `json:"original_format"`
	Struct         string `json:"struct"`
}

func fetchAndWriteSpec(client *http.Client, endpoint, smiles, outPath string, pretty bool) error {
	payload := map[string]string{
		"output_format": "chemical/x-indigo-ket",
		"struct":        smiles,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("status %d: %s", resp.StatusCode, truncate(string(respBytes), 200))
	}

	var envelope indigoLayoutResponse
	if err := json.Unmarshal(respBytes, &envelope); err != nil {
		return fmt.Errorf("decode envelope: %w", err)
	}
	if envelope.Struct == "" {
		return fmt.Errorf("empty struct field in response")
	}

	out := []byte(envelope.Struct)
	if pretty {
		var doc interface{}
		if err := json.Unmarshal(out, &doc); err != nil {
			return fmt.Errorf("decode struct payload: %w", err)
		}
		out, err = json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return fmt.Errorf("re-encode struct payload: %w", err)
		}
	}
	if err := os.WriteFile(outPath, append(out, '\n'), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func extractSMILES(compound pcCompound) compoundSMILES {
	result := compoundSMILES{CID: compound.CID}
	for _, prop := range compound.Props {
		label := strings.TrimSpace(prop.URN.Label)
		name := strings.TrimSpace(prop.URN.Name)
		value := prop.Value.text()
		if value == "" || label != "SMILES" {
			continue
		}

		if name == "Connectivity" {
			result.SMILES = value
			return result
		}
		if name == "Absolute" && result.SMILES == "" {
			result.SMILES = value
		}
	}
	return result
}

func apply2DLayout(mol *molecule.Molecule) {
	const (
		width      = 600
		height     = 400
		coordScale = 0.03
	)

	points := render2d.Layout(mol, render2d.Options{
		Width:  width,
		Height: height,
		Margin: 24,
	})
	for i, point := range points {
		// Convert render-space pixels into compact KET drawing coordinates.
		mol.SetAtomXY(i, (point.X-float64(width)/2)*coordScale, -(point.Y-float64(height)/2)*coordScale)
	}
}
