// Package molecule provides InChI (IUPAC International Chemical Identifier) generation
// and InChIKey calculation functionality.
//
// ⚠️  WARNING: This is a SIMPLIFIED implementation that does NOT produce standard-compliant InChI!
//
// LIMITATIONS:
// - Does not use proper canonical atom numbering (lacks graph automorphism algorithm)
// - Connectivity layer will differ from standard InChI
// - InChIKeys will not match official InChIKey for the same molecule
// - DO NOT use for production systems requiring standard InChI/InChIKey
//
// RECOMMENDED: For production use, use the official IUPAC InChI C library via CGO
// or use established frameworks like RDKit, OpenBabel, or CDK.
// See INCHI_LIMITATIONS.md for detailed information and solutions.
//
// InChI Structure:
// InChI is organized in layers, each providing specific information about the molecule:
// 1. Formula Layer (/): Chemical formula (e.g., C6H12O6)
// 2. Connectivity Layer (/c): Atom connections
// 3. Hydrogen Layer (/h): Hydrogen atom distribution
// 4. Double Bond Stereochemistry Layer (/b): Cis/trans configuration
// 5. Tetrahedral Stereochemistry Layer (/t): Chiral centers
// 6. Enantiomer Layer (/m): Enantiomer information
// 7. Stereo Type Layer (/s): Stereochemistry type
//
// References:
// - IUPAC InChI Technical Manual: https://www.inchi-trust.org/technical-faq/
// - InChI API Reference: https://www.inchi-trust.org/downloads/
// - Based on Indigo's molecule_inchi and molecule_inchi_layers implementation

package molecule

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
)

// InChIGenerator generates InChI strings from molecules
type InChIGenerator struct {
	prefix  string // InChI version prefix
	options InChIOptions
}

// InChIOptions contains options for InChI generation
type InChIOptions struct {
	// FixedH indicates if hydrogen layer should be included
	FixedH bool
	// RecMet indicates reconnected metals option
	RecMet bool
	// AuxInfo indicates if auxiliary information should be generated
	AuxInfo bool
	// SNon indicates if stereo information for non-stereogenic centers should be included
	SNon bool
}

// InChIResult contains the generated InChI and additional information
type InChIResult struct {
	InChI    string   // The generated InChI string
	InChIKey string   // The generated InChIKey
	AuxInfo  string   // Auxiliary information
	Warnings []string // Any warnings during generation
	Log      []string // Log messages
}

// NewInChIGenerator creates a new InChI generator with default options
func NewInChIGenerator() *InChIGenerator {
	return &InChIGenerator{
		prefix: "InChI=1S",
		options: InChIOptions{
			FixedH:  false,
			RecMet:  false,
			AuxInfo: false,
			SNon:    false,
		},
	}
}

// SetPrefix sets the InChI version prefix
func (g *InChIGenerator) SetPrefix(prefix string) {
	g.prefix = prefix
}

// SetOptions sets the InChI generation options
func (g *InChIGenerator) SetOptions(options InChIOptions) {
	g.options = options
}

// GenerateInChI generates InChI from a molecule
//
// Algorithm (based on Indigo's molecule_inchi.cpp):
// 1. Normalize and canonicalize the molecule structure
// 2. Decompose into connected components (for multi-component molecules)
// 3. For each component, generate layers:
//   - Formula layer: atom counts in Hill system order (C, H, then alphabetical)
//   - Connectivity layer: canonical numbering and bond connections using DFS
//   - Hydrogen layer: implicit hydrogen distribution
//   - Stereochemistry layers: cis/trans and tetrahedral stereocenters
//
// 4. Sort components by complexity
// 5. Combine layers into final InChI string
func (g *InChIGenerator) GenerateInChI(mol *Molecule) (*InChIResult, error) {
	if mol == nil {
		return nil, fmt.Errorf("molecule is nil")
	}

	// Validate molecule
	if err := g.validateMolecule(mol); err != nil {
		return nil, fmt.Errorf("invalid molecule: %w", err)
	}

	result := &InChIResult{
		Warnings: make([]string, 0),
		Log:      make([]string, 0),
	}

	// Handle empty molecule
	if mol.AtomCount() == 0 {
		result.InChI = g.prefix
		result.InChIKey = ""
		return result, nil
	}

	// Build InChI layers
	layers := g.buildInChILayers(mol)

	// Construct InChI string
	result.InChI = g.constructInChIString(layers)

	// Generate InChIKey
	var err error
	result.InChIKey, err = GenerateInChIKey(result.InChI)
	if err != nil {
		return nil, fmt.Errorf("failed to generate InChIKey: %w", err)
	}

	return result, nil
}

// validateMolecule checks if the molecule can be converted to InChI
func (g *InChIGenerator) validateMolecule(mol *Molecule) error {
	// Check for unsupported features
	for i := 0; i < mol.AtomCount(); i++ {
		atom := &mol.Atoms[i]
		// Check for pseudo atoms (not supported in standard InChI)
		if atom.Number == ELEM_PSEUDO {
			return fmt.Errorf("pseudo atoms are not supported in InChI")
		}
		// Check for R-groups (not supported in standard InChI)
		if atom.Number == ELEM_RSITE {
			return fmt.Errorf("R-group atoms are not supported in InChI")
		}
	}
	return nil
}

// inchiLayers holds the different layers of InChI
type inchiLayers struct {
	formula       string   // Formula layer
	connectivity  string   // Connectivity layer (/c)
	hydrogen      string   // Hydrogen layer (/h)
	cistrans      string   // Cis/trans stereochemistry (/b)
	tetrahedral   string   // Tetrahedral stereochemistry (/t)
	enantiomer    string   // Enantiomer layer (/m)
	stereoType    string   // Stereo type (/s)
	components    []string // For multi-component molecules
	hasStereochem bool     // Whether molecule has stereochemistry
}

// buildInChILayers constructs all InChI layers from the molecule
func (g *InChIGenerator) buildInChILayers(mol *Molecule) *inchiLayers {
	layers := &inchiLayers{
		components: make([]string, 0),
	}

	// Generate formula layer (Hill system: C, H, then alphabetical)
	layers.formula = g.generateFormulaLayer(mol)

	// Generate connectivity layer
	layers.connectivity = g.generateConnectivityLayer(mol)

	// Generate hydrogen layer
	layers.hydrogen = g.generateHydrogenLayer(mol)

	// Generate stereochemistry layers
	layers.cistrans = g.generateCisTransLayer(mol)
	layers.tetrahedral = g.generateTetrahedralLayer(mol)

	// Check if there's any stereochemistry
	layers.hasStereochem = layers.cistrans != "" || layers.tetrahedral != ""

	if layers.hasStereochem {
		layers.enantiomer = g.generateEnantiomerLayer(mol)
		layers.stereoType = "1" // Standard stereochemistry
	}

	return layers
}

// generateFormulaLayer generates the chemical formula in Hill system order
// Hill system: C first, then H, then other elements alphabetically
// Reference: Indigo's MainLayerFormula::printFormula (molecule_inchi_layers.cpp, lines 81-101)
func (g *InChIGenerator) generateFormulaLayer(mol *Molecule) string {
	if mol.AtomCount() == 0 {
		return ""
	}

	// Count atoms by element
	elementCount := make(map[int]int)
	for i := 0; i < mol.AtomCount(); i++ {
		atom := &mol.Atoms[i]
		elementCount[atom.Number]++

		// Add implicit hydrogens
		implH := mol.GetImplicitH(i)
		if implH > 0 {
			elementCount[ELEM_H] += implH
		}
	}

	var formula strings.Builder

	// Hill system order: C, H, then alphabetical by symbol
	// Carbon (element 6)
	if count, ok := elementCount[ELEM_C]; ok && count > 0 {
		formula.WriteString("C")
		if count > 1 {
			formula.WriteString(fmt.Sprintf("%d", count))
		}
		delete(elementCount, ELEM_C)
	}

	// Hydrogen (element 1)
	if count, ok := elementCount[ELEM_H]; ok && count > 0 {
		formula.WriteString("H")
		if count > 1 {
			formula.WriteString(fmt.Sprintf("%d", count))
		}
		delete(elementCount, ELEM_H)
	}

	// Remaining elements in alphabetical order
	type elemCount struct {
		element int
		symbol  string
		count   int
	}
	remaining := make([]elemCount, 0, len(elementCount))
	for elem, count := range elementCount {
		if count > 0 {
			symbol := ElementSymbol(elem)
			remaining = append(remaining, elemCount{elem, symbol, count})
		}
	}
	sort.Slice(remaining, func(i, j int) bool {
		return remaining[i].symbol < remaining[j].symbol
	})

	for _, ec := range remaining {
		formula.WriteString(ec.symbol)
		if ec.count > 1 {
			formula.WriteString(fmt.Sprintf("%d", ec.count))
		}
	}

	return formula.String()
}

// generateConnectivityLayer generates the connectivity layer showing atom connections
// Format: DFS-based connection tree (e.g., "1-2-3(4,5)-6" means 1 connects to 2, 2 to 3, 3 branches to 4 and 5, then continues to 6)
// Reference: Indigo's MainLayerConnections::printConnectionTable (molecule_inchi_layers.cpp, lines 248-422)
func (g *InChIGenerator) generateConnectivityLayer(mol *Molecule) string {
	if mol.AtomCount() <= 1 {
		return ""
	}

	// No need for canonical numbering in simple case - use sequential order
	// Real InChI would use canonical ordering, but for simplicity we use atom order

	// Find starting atom (one with minimum degree, or first heavy atom)
	startIdx := g.findStartAtom(mol)

	// Build connectivity string using DFS
	visited := make([]bool, mol.AtomCount())
	connStr := g.buildConnectivityDFS(mol, startIdx, -1, visited)

	return connStr
}

// findStartAtom finds the best starting atom for DFS traversal
// Prefer atoms with lowest degree (to minimize branches)
func (g *InChIGenerator) findStartAtom(mol *Molecule) int {
	minDegree := mol.AtomCount() + 1
	startIdx := 0

	for i := 0; i < mol.AtomCount(); i++ {
		degree := len(mol.Vertices[i].Edges)
		if degree < minDegree {
			minDegree = degree
			startIdx = i
		}
	}

	return startIdx
}

// buildConnectivityDFS builds connectivity string using depth-first search
// Reference: Indigo's DFS-based connectivity algorithm (molecule_inchi_layers.cpp, lines 318-419)
func (g *InChIGenerator) buildConnectivityDFS(mol *Molecule, current, parent int, visited []bool) string {
	visited[current] = true

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%d", current+1)) // 1-based indexing

	// Get unvisited neighbors
	vertex := &mol.Vertices[current]
	neighbors := make([]int, 0)

	for _, bondIdx := range vertex.Edges {
		bond := &mol.Bonds[bondIdx]
		neighbor := bond.End
		if neighbor == current {
			neighbor = bond.Beg
		}
		if !visited[neighbor] {
			neighbors = append(neighbors, neighbor)
		}
	}

	// Sort neighbors by index for consistency
	sort.Ints(neighbors)

	if len(neighbors) == 0 {
		return result.String()
	}

	if len(neighbors) == 1 {
		// Single branch - continue with dash
		result.WriteString("-")
		result.WriteString(g.buildConnectivityDFS(mol, neighbors[0], current, visited))
	} else {
		// Multiple branches - use parentheses
		result.WriteString("(")
		for i, neighbor := range neighbors {
			if i > 0 {
				result.WriteString(",")
			}
			result.WriteString(g.buildConnectivityDFS(mol, neighbor, current, visited))
		}
		result.WriteString(")")
	}

	return result.String()
}

// generateHydrogenLayer generates the hydrogen atom layer
// Shows which atoms have how many hydrogens
// Reference: Indigo's HydrogensLayer::print (molecule_inchi_layers.cpp, lines 468-528)
func (g *InChIGenerator) generateHydrogenLayer(mol *Molecule) string {
	// Collect atoms with hydrogens
	type hydrogenInfo struct {
		atomIdx int
		count   int
	}

	atomsWithH := make(map[int][]int) // count -> atom indices

	for i := 0; i < mol.AtomCount(); i++ {
		implH := mol.GetImplicitH(i)
		if implH > 0 {
			atomsWithH[implH] = append(atomsWithH[implH], i+1) // 1-based
		}
	}

	if len(atomsWithH) == 0 {
		return ""
	}

	// Sort by H count
	hCounts := make([]int, 0, len(atomsWithH))
	for count := range atomsWithH {
		hCounts = append(hCounts, count)
	}
	sort.Ints(hCounts)

	var result strings.Builder
	first := true

	for _, hCount := range hCounts {
		atoms := atomsWithH[hCount]
		sort.Ints(atoms)

		// Format: 1-3H2 means atoms 1,2,3 each have 2 hydrogens
		// Or: 1,4H means atoms 1 and 4 each have 1 hydrogen

		// Check for ranges
		i := 0
		for i < len(atoms) {
			if !first {
				result.WriteString(",")
			}
			first = false

			start := atoms[i]
			end := start

			// Find consecutive range
			for i+1 < len(atoms) && atoms[i+1] == atoms[i]+1 {
				i++
				end = atoms[i]
			}

			if end > start {
				result.WriteString(fmt.Sprintf("%d-%d", start, end))
			} else {
				result.WriteString(fmt.Sprintf("%d", start))
			}

			i++
		}

		result.WriteString("H")
		if hCount > 1 {
			result.WriteString(fmt.Sprintf("%d", hCount))
		}
	}

	return result.String()
}

// generateCisTransLayer generates cis/trans stereochemistry layer for double bonds
// Reference: Indigo's CisTransStereochemistryLayer::print (molecule_inchi_layers.cpp, lines 567-610)
func (g *InChIGenerator) generateCisTransLayer(mol *Molecule) string {
	if mol.CisTrans == nil || mol.CisTrans.Count() == 0 {
		return ""
	}

	var stereoDescriptors []struct {
		beg    int
		end    int
		parity string
	}

	// Process each bond with cis/trans stereochemistry
	for bondIdx := 0; bondIdx < mol.BondCount(); bondIdx++ {
		bond := &mol.Bonds[bondIdx]

		// Only process double bonds
		if bond.Order != BOND_DOUBLE {
			continue
		}

		// Check if this bond has stereochemistry
		parity := mol.CisTrans.GetParity(bondIdx)
		if parity == 0 {
			continue
		}

		// Skip if explicitly ignored
		if mol.CisTrans.IsIgnored(bondIdx) {
			continue
		}

		// Get bond atoms (1-based)
		beg := bond.Beg + 1
		end := bond.End + 1

		// Encode parity: + for TRANS, - for CIS
		parityStr := "-"
		if parity == TRANS {
			parityStr = "+"
		}

		// Always put smaller index first
		if beg > end {
			beg, end = end, beg
		}

		stereoDescriptors = append(stereoDescriptors, struct {
			beg    int
			end    int
			parity string
		}{beg, end, parityStr})
	}

	if len(stereoDescriptors) == 0 {
		return ""
	}

	// Sort by bond indices
	sort.Slice(stereoDescriptors, func(i, j int) bool {
		if stereoDescriptors[i].end != stereoDescriptors[j].end {
			return stereoDescriptors[i].end < stereoDescriptors[j].end
		}
		return stereoDescriptors[i].beg < stereoDescriptors[j].beg
	})

	// Build the layer string
	var result strings.Builder
	for i, desc := range stereoDescriptors {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%d-%d%s", desc.end, desc.beg, desc.parity))
	}

	return result.String()
}

// generateTetrahedralLayer generates tetrahedral stereochemistry layer
// Reference: Indigo's TetrahedralStereochemistryLayer::print (molecule_inchi_layers.cpp, lines 711-730)
func (g *InChIGenerator) generateTetrahedralLayer(mol *Molecule) string {
	if mol.Stereocenters == nil || mol.Stereocenters.Size() == 0 {
		return ""
	}

	var stereoDescriptors []struct {
		atomIdx int
		parity  string
	}

	// Find first stereocenter to determine reference
	firstSign := 0
	for atomIdx := 0; atomIdx < mol.AtomCount(); atomIdx++ {
		if !mol.Stereocenters.Exists(atomIdx) {
			continue
		}

		center, err := mol.Stereocenters.Get(atomIdx)
		if err != nil || !center.IsTetrahydral {
			continue
		}

		if center.Type == STEREO_ATOM_ANY {
			continue
		}

		sign := g.computeTetrahedralSign(center)
		if firstSign == 0 {
			firstSign = -sign
		}

		// Compute parity relative to first center
		parity := "+"
		if sign*firstSign == -1 {
			parity = "-"
		}

		stereoDescriptors = append(stereoDescriptors, struct {
			atomIdx int
			parity  string
		}{atomIdx + 1, parity}) // 1-based
	}

	if len(stereoDescriptors) == 0 {
		return ""
	}

	// Sort by atom index
	sort.Slice(stereoDescriptors, func(i, j int) bool {
		return stereoDescriptors[i].atomIdx < stereoDescriptors[j].atomIdx
	})

	// Build the layer string
	var result strings.Builder
	for i, desc := range stereoDescriptors {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%d%s", desc.atomIdx, desc.parity))
	}

	return result.String()
}

// computeTetrahedralSign computes the stereochemical sign for a tetrahedral center
// Reference: Indigo's _getMappingSign (molecule_inchi_layers.cpp, lines 797-825)
func (g *InChIGenerator) computeTetrahedralSign(center *Stereocenter) int {
	pyramid := center.Pyramid

	// Move minimal element to end
	minIdx := 0
	minVal := pyramid[0]
	for i := 1; i < 4; i++ {
		if pyramid[i] < minVal {
			minVal = pyramid[i]
			minIdx = i
		}
	}

	// Swap to put min at position 3
	if minIdx != 3 {
		pyramid[minIdx], pyramid[3] = pyramid[3], pyramid[minIdx]
	}

	// Count inversions in first 3 elements
	cnt := 0
	for i := 0; i < 2; i++ {
		if pyramid[i] > pyramid[i+1] {
			cnt++
		}
	}
	if pyramid[0] > pyramid[2] {
		cnt++
	}

	if cnt%2 == 0 {
		return 1
	}
	return -1
}

// generateEnantiomerLayer generates enantiomer information layer
// Reference: Indigo's TetrahedralStereochemistryLayer::printEnantiomers (molecule_inchi_layers.cpp, lines 745-755)
func (g *InChIGenerator) generateEnantiomerLayer(mol *Molecule) string {
	if mol.Stereocenters == nil || mol.Stereocenters.Size() == 0 {
		return "0" // Default to absolute
	}

	// Find first stereocenter
	for atomIdx := 0; atomIdx < mol.AtomCount(); atomIdx++ {
		if !mol.Stereocenters.Exists(atomIdx) {
			continue
		}

		center, err := mol.Stereocenters.Get(atomIdx)
		if err != nil || !center.IsTetrahydral {
			continue
		}

		if center.Type == STEREO_ATOM_ANY {
			continue
		}

		sign := g.computeTetrahedralSign(center)
		if sign == 1 {
			return "1"
		} else if sign == -1 {
			return "0"
		}
	}

	return "0"
}

// constructInChIString combines all layers into final InChI string
func (g *InChIGenerator) constructInChIString(layers *inchiLayers) string {
	var result strings.Builder

	// Start with prefix
	result.WriteString(g.prefix)

	// Add formula layer
	if layers.formula != "" {
		result.WriteString("/")
		result.WriteString(layers.formula)
	}

	// Add connectivity layer
	if layers.connectivity != "" {
		result.WriteString("/c")
		result.WriteString(layers.connectivity)
	}

	// Add hydrogen layer
	if layers.hydrogen != "" {
		result.WriteString("/h")
		result.WriteString(layers.hydrogen)
	}

	// Add cis/trans stereochemistry
	if layers.cistrans != "" {
		result.WriteString("/b")
		result.WriteString(layers.cistrans)
	}

	// Add tetrahedral stereochemistry
	if layers.tetrahedral != "" {
		result.WriteString("/t")
		result.WriteString(layers.tetrahedral)
	}

	// Add enantiomer layer if stereochemistry exists
	if layers.hasStereochem {
		result.WriteString("/m")
		result.WriteString(layers.enantiomer)
		result.WriteString("/s")
		result.WriteString(layers.stereoType)
	}

	return result.String()
}

// GenerateInChIKey generates InChIKey from InChI string
//
// InChIKey Algorithm (IUPAC specification):
//  1. Split InChI into main (connectivity) and stereochemistry parts
//  2. Hash each part using SHA-256
//  3. Encode hash bits as base-26 characters (A-Z)
//  4. Format: XXXXXXXXXXXXXX-YYYYYYYYY-ZZ
//     X = connectivity hash (14 chars)
//     Y = stereochemistry hash (9-10 chars)
//     Z = version and flags (1-2 chars)
//
// Reference: InChI Technical FAQ, InChIKey specification
func GenerateInChIKey(inchi string) (string, error) {
	if inchi == "" {
		return "", fmt.Errorf("empty InChI string")
	}

	// Validate InChI prefix
	if !strings.HasPrefix(inchi, "InChI=") {
		return "", fmt.Errorf("invalid InChI format: missing 'InChI=' prefix")
	}

	// Remove prefix for hashing
	inchiBody := strings.TrimPrefix(inchi, "InChI=")

	// Split into main and stereochemistry parts
	// Stereochemistry starts at /t, /m, or /b layer
	mainPart := inchiBody
	stereoPart := ""

	// Find stereochemistry layers
	stereoIdx := -1
	for _, layer := range []string{"/b", "/t", "/m", "/s"} {
		idx := strings.Index(inchiBody, layer)
		if idx != -1 && (stereoIdx == -1 || idx < stereoIdx) {
			stereoIdx = idx
		}
	}

	if stereoIdx != -1 {
		mainPart = inchiBody[:stereoIdx]
		stereoPart = inchiBody[stereoIdx:]
	}

	// Hash the main part (connectivity)
	mainHash := sha256.Sum256([]byte(mainPart))
	connectivityBlock := encodeBase26(mainHash[:], 14)

	// Hash the stereochemistry part
	var stereoBlock string
	if stereoPart != "" {
		stereoHash := sha256.Sum256([]byte(stereoPart))
		stereoBlock = encodeBase26(stereoHash[:], 10)
	} else {
		// No stereochemistry - use standard placeholder
		stereoBlock = "UHFFFAOYSA"
	}

	// Version and protonation flag
	// Format: XY where X is version, Y is protonation state
	// S = standard InChI version 1
	// A = no protonation
	// N = non-standard
	// O = +1 protonation
	versionFlag := "N" // Non-standard by default
	if strings.Contains(inchi, "/p") {
		// Protonation layer present
		if strings.Contains(inchi, "/p+1") {
			versionFlag = "O"
		} else {
			versionFlag = "M"
		}
	} else {
		versionFlag = "N"
	}

	// Construct InChIKey
	inchiKey := fmt.Sprintf("%s-%s-%s", connectivityBlock, stereoBlock, versionFlag)

	return inchiKey, nil
}

// encodeBase26 encodes byte array to base26 string (A-Z)
func encodeBase26(data []byte, length int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Convert first 8 bytes to uint64
	var num uint64
	for i := 0; i < 8 && i < len(data); i++ {
		num = (num << 8) | uint64(data[i])
	}

	// Handle additional bytes if we need more entropy
	if length > 12 {
		// Use more bytes from hash
		var num2 uint64
		for i := 8; i < 16 && i < len(data); i++ {
			num2 = (num2 << 8) | uint64(data[i])
		}
		// Combine both numbers
		num = num ^ (num2 >> 32)
	}

	// Convert to base26
	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = alphabet[num%26]
		num /= 26
	}

	return string(result)
}

// encodeBase26Better uses better hash distribution
func encodeBase26Better(data []byte, length int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)

	// Use bytes directly with better distribution
	for i := 0; i < length; i++ {
		if i < len(data) {
			// Mix multiple bytes for better distribution
			idx := i
			val := uint(data[idx])
			if idx+1 < len(data) {
				val = (val*31 + uint(data[idx+1])) % 26
			} else {
				val = val % 26
			}
			result[i] = alphabet[val]
		} else {
			result[i] = 'A'
		}
	}

	return string(result)
}

// ParseInChI parses an InChI string into a molecule
func ParseInChI(inchi string) (*Molecule, error) {
	if !strings.HasPrefix(inchi, "InChI=") {
		return nil, fmt.Errorf("invalid InChI format")
	}

	// TODO: Implement InChI parsing
	return nil, fmt.Errorf("InChI parsing not yet implemented")
}

// ValidateInChI checks if an InChI string is valid
func ValidateInChI(inchi string) bool {
	if !strings.HasPrefix(inchi, "InChI=") {
		return false
	}

	// Check for required layers
	parts := strings.Split(inchi, "/")
	if len(parts) < 2 {
		return false
	}

	return true
}

// CompareInChI compares two InChI strings for equivalence
func CompareInChI(inchi1, inchi2 string) int {
	// Normalize by removing version differences
	norm1 := strings.TrimPrefix(inchi1, "InChI=1S/")
	norm1 = strings.TrimPrefix(norm1, "InChI=1/")

	norm2 := strings.TrimPrefix(inchi2, "InChI=1S/")
	norm2 = strings.TrimPrefix(norm2, "InChI=1/")

	if norm1 == norm2 {
		return 0
	}
	if norm1 < norm2 {
		return -1
	}
	return 1
}

// GetInChIFromSMILES is a convenience function that converts SMILES to InChI
func GetInChIFromSMILES(smiles string) (*InChIResult, error) {
	// Load molecule from SMILES
	loader := SmilesLoader{}
	mol, err := loader.Parse(smiles)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMILES: %w", err)
	}

	// Generate InChI
	generator := NewInChIGenerator()
	return generator.GenerateInChI(mol)
}

// Base64EncodeInChI encodes InChI to base64 for compact storage
func Base64EncodeInChI(inchi string) string {
	return base64.StdEncoding.EncodeToString([]byte(inchi))
}

// Base64DecodeInChI decodes base64-encoded InChI
func Base64DecodeInChI(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// Helper function to avoid unused import error
var _ = binary.Size
