// Package molecule provides InChI (IUPAC International Chemical Identifier) generation
// and InChIKey calculation functionality.
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
// - Based on Indigo's molecule_inchi implementation

package molecule

import (
	"crypto/sha256"
	"encoding/base64"
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
// Algorithm (based on IUPAC InChI specification):
// 1. Normalize and canonicalize the molecule structure
// 2. Decompose into connected components
// 3. For each component, generate layers:
//   - Formula layer: atom counts in Hill system order (C, H, then alphabetical)
//   - Connectivity layer: canonical numbering and bond connections
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
	if count, ok := elementCount[6]; ok && count > 0 {
		formula.WriteString("C")
		if count > 1 {
			formula.WriteString(fmt.Sprintf("%d", count))
		}
		delete(elementCount, 6)
	}

	// Hydrogen (element 1)
	if count, ok := elementCount[1]; ok && count > 0 {
		formula.WriteString("H")
		if count > 1 {
			formula.WriteString(fmt.Sprintf("%d", count))
		}
		delete(elementCount, 1)
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
// Format: 1-2-3,4-5 means atoms connected in a tree structure
func (g *InChIGenerator) generateConnectivityLayer(mol *Molecule) string {
	if mol.AtomCount() <= 1 {
		return ""
	}

	// Create canonical numbering
	canonicalOrder := g.getCanonicalNumbering(mol)

	// Build connectivity string
	visited := make([]bool, mol.AtomCount())
	var connectivityParts []string

	// Start from first atom in canonical order
	for _, startIdx := range canonicalOrder {
		if visited[startIdx] {
			continue
		}

		// BFS/DFS to build connectivity string for this component
		connStr := g.buildConnectivityString(mol, startIdx, visited, canonicalOrder)
		if connStr != "" {
			connectivityParts = append(connectivityParts, connStr)
		}
	}

	if len(connectivityParts) == 0 {
		return ""
	}

	return strings.Join(connectivityParts, ";")
}

// getCanonicalNumbering returns canonical atom numbering
// This is a simplified version - full implementation would use graph automorphism
func (g *InChIGenerator) getCanonicalNumbering(mol *Molecule) []int {
	// For now, return sequential numbering
	// TODO: Implement proper canonical ordering based on:
	// 1. Atomic number
	// 2. Number of connections
	// 3. Bond orders
	// 4. Ring membership
	// 5. Stereochemistry
	order := make([]int, mol.AtomCount())
	for i := range order {
		order[i] = i
	}

	// Simple sorting by atomic number and degree
	type atomInfo struct {
		index   int
		element int
		degree  int
	}
	atoms := make([]atomInfo, mol.AtomCount())
	for i := 0; i < mol.AtomCount(); i++ {
		atom := &mol.Atoms[i]
		vertex := &mol.Vertices[i]
		atoms[i] = atomInfo{
			index:   i,
			element: atom.Number,
			degree:  len(vertex.Edges),
		}
	}

	sort.Slice(atoms, func(i, j int) bool {
		if atoms[i].element != atoms[j].element {
			return atoms[i].element > atoms[j].element
		}
		return atoms[i].degree > atoms[j].degree
	})

	for i, a := range atoms {
		order[i] = a.index
	}

	return order
}

// buildConnectivityString builds connectivity string for a connected component
func (g *InChIGenerator) buildConnectivityString(mol *Molecule, start int, visited []bool, canonicalOrder []int) string {
	var result strings.Builder
	queue := []int{start}
	visited[start] = true

	// Create reverse mapping for canonical numbering
	canonicalIndex := make(map[int]int)
	for i, idx := range canonicalOrder {
		canonicalIndex[idx] = i + 1 // InChI uses 1-based indexing
	}

	// Simple BFS
	isFirst := true
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if !isFirst {
			result.WriteString("-")
		}
		isFirst = false
		result.WriteString(fmt.Sprintf("%d", canonicalIndex[current]))

		neighbors := mol.GetNeighbors(current)
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return result.String()
}

// generateHydrogenLayer generates the hydrogen atom layer
// Shows which atoms have how many hydrogens
func (g *InChIGenerator) generateHydrogenLayer(mol *Molecule) string {
	if !g.options.FixedH {
		// Skip hydrogen layer if not explicitly requested
		return ""
	}

	var parts []string
	for i := 0; i < mol.AtomCount(); i++ {
		implH := mol.GetImplicitH(i)
		if implH > 0 {
			parts = append(parts, fmt.Sprintf("%d,%dH", i+1, implH))
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ",")
}

// generateCisTransLayer generates cis/trans stereochemistry layer for double bonds
//
// Algorithm:
// 1. Iterate through all bonds with cis/trans stereochemistry
// 2. Get substituents and determine configuration
// 3. Encode in InChI format: bond_number+ (trans) or bond_number- (cis)
//
// Reference: Indigo's molecule_inchi.cpp, lines 105-115
// Reference: IUPAC InChI Technical Manual, Section 3.4
func (g *InChIGenerator) generateCisTransLayer(mol *Molecule) string {
	if mol.CisTrans == nil || mol.CisTrans.Count() == 0 {
		return ""
	}

	// Get canonical numbering for proper atom ordering
	canonicalOrder := g.getCanonicalNumbering(mol)
	canonicalIndex := make(map[int]int)
	for i, idx := range canonicalOrder {
		canonicalIndex[idx] = i + 1 // 1-based indexing
	}

	var stereoDescriptors []struct {
		bondCanonical int
		parity        string
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

		// Get canonical indices of the bond atoms
		begCanonical := canonicalIndex[bond.Beg]
		endCanonical := canonicalIndex[bond.End]

		// Encode parity: CIS (-) or TRANS (+)
		parityStr := "-"
		if parity == TRANS {
			parityStr = "+"
		}

		// Store bond descriptor
		// Use smaller canonical number first
		bondCanonical := begCanonical
		if endCanonical < begCanonical {
			bondCanonical = endCanonical
		}

		stereoDescriptors = append(stereoDescriptors, struct {
			bondCanonical int
			parity        string
		}{bondCanonical, parityStr})
	}

	if len(stereoDescriptors) == 0 {
		return ""
	}

	// Sort by canonical bond number
	sort.Slice(stereoDescriptors, func(i, j int) bool {
		return stereoDescriptors[i].bondCanonical < stereoDescriptors[j].bondCanonical
	})

	// Build the layer string
	var result strings.Builder
	for i, desc := range stereoDescriptors {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%d%s", desc.bondCanonical, desc.parity))
	}

	return result.String()
}

// generateTetrahedralLayer generates tetrahedral stereochemistry layer
//
// Algorithm:
// 1. Iterate through all tetrahedral stereocenters
// 2. Compute parity based on pyramid configuration
// 3. Encode in InChI format: atom_number+ or atom_number-
//
// Reference: Indigo's molecule_inchi.cpp, lines 117-131
// Reference: IUPAC InChI Technical Manual, Section 3.5
func (g *InChIGenerator) generateTetrahedralLayer(mol *Molecule) string {
	if mol.Stereocenters == nil || mol.Stereocenters.Size() == 0 {
		return ""
	}

	// Get canonical numbering
	canonicalOrder := g.getCanonicalNumbering(mol)
	canonicalIndex := make(map[int]int)
	for i, idx := range canonicalOrder {
		canonicalIndex[idx] = i + 1 // 1-based indexing
	}

	var stereoDescriptors []struct {
		atomCanonical int
		parity        string
	}

	// Iterate through all atoms to find stereocenters
	for atomIdx := 0; atomIdx < mol.AtomCount(); atomIdx++ {
		if !mol.Stereocenters.Exists(atomIdx) {
			continue
		}

		center, err := mol.Stereocenters.Get(atomIdx)
		if err != nil {
			continue
		}

		// Only process tetrahedral centers
		if !center.IsTetrahydral {
			continue
		}

		// Skip ANY type (undefined stereochemistry)
		if center.Type == STEREO_ATOM_ANY {
			continue
		}

		// Compute parity based on pyramid configuration
		// This is a simplified version - full implementation needs CIP rules
		parity := g.computeTetrahedralParity(mol, center, canonicalIndex)

		atomCanonical := canonicalIndex[atomIdx]
		stereoDescriptors = append(stereoDescriptors, struct {
			atomCanonical int
			parity        string
		}{atomCanonical, parity})
	}

	if len(stereoDescriptors) == 0 {
		return ""
	}

	// Sort by canonical atom number
	sort.Slice(stereoDescriptors, func(i, j int) bool {
		return stereoDescriptors[i].atomCanonical < stereoDescriptors[j].atomCanonical
	})

	// Build the layer string
	var result strings.Builder
	for i, desc := range stereoDescriptors {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%d%s", desc.atomCanonical, desc.parity))
	}

	return result.String()
}

// computeTetrahedralParity computes the parity for a tetrahedral center
// This is a simplified implementation - full version needs Cahn-Ingold-Prelog rules
func (g *InChIGenerator) computeTetrahedralParity(mol *Molecule, center *Stereocenter, canonicalIndex map[int]int) string {
	// Get pyramid configuration
	pyramid := center.Pyramid

	// Convert to canonical indices
	canonicalPyramid := make([]int, 4)
	for i := 0; i < 4; i++ {
		if pyramid[i] == -1 {
			canonicalPyramid[i] = -1 // Implicit hydrogen
		} else {
			canonicalPyramid[i] = canonicalIndex[pyramid[i]]
		}
	}

	// Determine parity by checking ordering
	// This is simplified - should use proper stereochemical determination
	// Count inversions to determine odd/even parity
	inversions := 0
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			if canonicalPyramid[i] != -1 && canonicalPyramid[j] != -1 {
				if canonicalPyramid[i] > canonicalPyramid[j] {
					inversions++
				}
			}
		}
	}

	// Odd number of inversions = '+', even = '-'
	if inversions%2 == 1 {
		return "+"
	}
	return "-"
}

// generateEnantiomerLayer generates enantiomer information layer
//
// This layer indicates the stereochemistry type:
// - "0" = absolute stereochemistry
// - "1" = relative stereochemistry (racemic/relative)
//
// Reference: Indigo's molecule_inchi.cpp, generateEnantiomerLayer method
// Reference: IUPAC InChI Technical Manual, Section 3.6
func (g *InChIGenerator) generateEnantiomerLayer(mol *Molecule) string {
	if mol.Stereocenters == nil || mol.Stereocenters.Size() == 0 {
		return "0" // Default to absolute
	}

	// Check stereocenter types
	hasRelative := false
	hasAnd := false

	for atomIdx := 0; atomIdx < mol.AtomCount(); atomIdx++ {
		if !mol.Stereocenters.Exists(atomIdx) {
			continue
		}

		center, err := mol.Stereocenters.Get(atomIdx)
		if err != nil {
			continue
		}

		switch center.Type {
		case STEREO_ATOM_OR:
			hasRelative = true
		case STEREO_ATOM_AND:
			hasAnd = true
		}
	}

	// Determine enantiomer type
	// If any stereocenter is AND (racemic), use "1"
	if hasAnd {
		return "1"
	}

	// If any stereocenter is relative (OR), use "1"
	if hasRelative {
		return "1"
	}

	// Default to absolute stereochemistry
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
//  1. Split InChI into layers
//  2. For each layer, compute SHA-256 hash
//  3. Take first 65 bits for connectivity block (14 chars in base26)
//  4. Take next 37 bits for stereochemistry block (9 chars in base26)
//  5. Add version and protonation flag characters
//  6. Format: XXXXXXXXXXXXXX-YYYYYYYYY-ZZ
//     X = connectivity layer (14 chars)
//     Y = stereochemistry layer (9 chars)
//     Z = version and flags (2 chars)
//
// Reference: Goodman et al., "InChI version 1, three years on", Journal of Cheminformatics (2012)
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
	// Standard InChI: InChI=1S/formula/c.../h.../b.../t.../m.../s...
	mainPart := inchiBody
	stereoPart := ""

	// Check if there's stereochemistry (indicated by /t, /m, or /s layers)
	if idx := strings.Index(inchiBody, "/t"); idx != -1 {
		mainPart = inchiBody[:idx]
		stereoPart = inchiBody[idx:]
	} else if idx := strings.Index(inchiBody, "/m"); idx != -1 {
		mainPart = inchiBody[:idx]
		stereoPart = inchiBody[idx:]
	} else if idx := strings.Index(inchiBody, "/s"); idx != -1 {
		mainPart = inchiBody[:idx]
		stereoPart = inchiBody[idx:]
	}

	// Hash the main part (connectivity)
	mainHash := sha256.Sum256([]byte(mainPart))
	connectivityBlock := encodeBase26(mainHash[:], 14)

	// Hash the stereochemistry part
	var stereoBlock string
	if stereoPart != "" {
		stereoHash := sha256.Sum256([]byte(stereoPart))
		stereoBlock = encodeBase26(stereoHash[:], 9)
	} else {
		// No stereochemistry - use placeholder
		stereoBlock = "UHFFFAOYSA"
	}

	// Version and protonation flag
	// S = standard InChI, A = no protonation
	versionFlag := "SA"
	if strings.Contains(inchi, "1S") {
		versionFlag = "SA"
	} else if strings.Contains(inchi, "1") {
		versionFlag = "N"
	}

	// Construct InChIKey
	inchiKey := fmt.Sprintf("%s-%s-%s", connectivityBlock, stereoBlock, versionFlag)

	return inchiKey, nil
}

// encodeBase26 encodes byte array to base26 string (A-Z)
// This is used for InChIKey generation
func encodeBase26(data []byte, length int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Convert bytes to big integer
	var num uint64
	for i := 0; i < 8 && i < len(data); i++ {
		num = (num << 8) | uint64(data[i])
	}

	// Convert to base26
	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = alphabet[num%26]
		num /= 26
	}

	return string(result)
}

// ParseInChI parses an InChI string into a molecule
// This is the reverse operation of GenerateInChI
func ParseInChI(inchi string) (*Molecule, error) {
	if !strings.HasPrefix(inchi, "InChI=") {
		return nil, fmt.Errorf("invalid InChI format")
	}

	// TODO: Implement InChI parsing
	// This requires parsing each layer and reconstructing the molecule
	// 1. Parse formula layer
	// 2. Parse connectivity layer
	// 3. Parse hydrogen layer
	// 4. Parse stereochemistry layers
	// 5. Build molecule from parsed data

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
// Returns 0 if equal, -1 if inchi1 < inchi2, 1 if inchi1 > inchi2
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
