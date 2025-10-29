// Package src provides molecular structure manipulation and analysis tools.
// This file implements molecular fingerprint generation for similarity searching.
package src

import (
	"hash/fnv"
	"math"
	"math/bits"
)

// FingerprintType represents different types of molecular fingerprints
type FingerprintType int

const (
	// FingerprintPath Path-based fingerprint (similar to Daylight)
	FingerprintPath FingerprintType = iota
	// FingerprintECFP2 Extended Connectivity Fingerprint radius 2
	FingerprintECFP2
	// FingerprintECFP4 Extended Connectivity Fingerprint radius 4
	FingerprintECFP4
	// FingerprintECFP6 Extended Connectivity Fingerprint radius 6
	FingerprintECFP6
)

// FingerprintParameters configures fingerprint generation
type FingerprintParameters struct {
	Type         FingerprintType // Type of fingerprint to generate
	Size         int             // Size in bits (default: 2048)
	MinPath      int             // Minimum path length (for path-based)
	MaxPath      int             // Maximum path length (for path-based)
	UseChirality bool            // Include stereochemistry
}

// DefaultFingerprintParameters returns default parameters
func DefaultFingerprintParameters() FingerprintParameters {
	return FingerprintParameters{
		Type:         FingerprintPath,
		Size:         2048,
		MinPath:      1,
		MaxPath:      7,
		UseChirality: false,
	}
}

// Fingerprint represents a molecular fingerprint as a bit vector
type Fingerprint struct {
	Bits   []uint64 // Bit vector stored as uint64 array
	Size   int      // Total size in bits
	Type   FingerprintType
	params FingerprintParameters
}

// NewFingerprint creates a new fingerprint with specified size
func NewFingerprint(params FingerprintParameters) *Fingerprint {
	if params.Size <= 0 {
		params.Size = 2048
	}
	numWords := (params.Size + 63) / 64
	return &Fingerprint{
		Bits:   make([]uint64, numWords),
		Size:   params.Size,
		Type:   params.Type,
		params: params,
	}
}

// SetBit sets a bit at the given position
func (fp *Fingerprint) SetBit(pos int) {
	if pos < 0 || pos >= fp.Size {
		return
	}
	wordIdx := pos / 64
	bitIdx := uint(pos % 64)
	fp.Bits[wordIdx] |= (1 << bitIdx)
}

// GetBit returns true if the bit at position is set
func (fp *Fingerprint) GetBit(pos int) bool {
	if pos < 0 || pos >= fp.Size {
		return false
	}
	wordIdx := pos / 64
	bitIdx := uint(pos % 64)
	return (fp.Bits[wordIdx] & (1 << bitIdx)) != 0
}

// CountBits returns the number of set bits (population count)
func (fp *Fingerprint) CountBits() int {
	count := 0
	for _, word := range fp.Bits {
		count += bits.OnesCount64(word)
	}
	return count
}

// SetBitsFromHash sets multiple bits based on a hash value
func (fp *Fingerprint) SetBitsFromHash(hash uint32, numBits int) {
	for i := 0; i < numBits; i++ {
		// Use different seeds for each bit to set
		seed := hash + uint32(i*0x9e3779b9)
		pos := int(seed % uint32(fp.Size))
		fp.SetBit(pos)
	}
}

// MoleculeFingerprintBuilder builds fingerprints from molecules
type MoleculeFingerprintBuilder struct {
	mol    *Molecule
	params FingerprintParameters
}

// NewFingerprintBuilder creates a new fingerprint builder
func NewFingerprintBuilder(mol *Molecule, params FingerprintParameters) *MoleculeFingerprintBuilder {
	return &MoleculeFingerprintBuilder{
		mol:    mol,
		params: params,
	}
}

// Build generates the fingerprint
func (fb *MoleculeFingerprintBuilder) Build() *Fingerprint {
	switch fb.params.Type {
	case FingerprintECFP2, FingerprintECFP4, FingerprintECFP6:
		return fb.buildECFP()
	default:
		return fb.buildPathFingerprint()
	}
}

// buildPathFingerprint generates a path-based fingerprint
func (fb *MoleculeFingerprintBuilder) buildPathFingerprint() *Fingerprint {
	fp := NewFingerprint(fb.params)

	// Enumerate all paths up to maxPath length
	for length := fb.params.MinPath; length <= fb.params.MaxPath; length++ {
		fb.enumeratePathsOfLength(fp, length)
	}

	return fp
}

// enumeratePathsOfLength finds all paths of a specific length
func (fb *MoleculeFingerprintBuilder) enumeratePathsOfLength(fp *Fingerprint, length int) {
	if length <= 0 {
		return
	}

	// For each atom as starting point
	for startAtom := 0; startAtom < len(fb.mol.Atoms); startAtom++ {
		visited := make([]bool, len(fb.mol.Atoms))
		path := []int{startAtom}
		fb.dfsPath(fp, startAtom, path, visited, length)
	}
}

// dfsPath performs depth-first search to find paths
func (fb *MoleculeFingerprintBuilder) dfsPath(fp *Fingerprint, current int, path []int, visited []bool, remaining int) {
	if remaining == 0 {
		// Hash this path and set bits
		hash := fb.hashPath(path)
		fp.SetBitsFromHash(hash, 2) // Set 2 bits per path
		return
	}

	visited[current] = true

	neighbors := fb.mol.GetNeighbors(current)
	for _, neighbor := range neighbors {
		if !visited[neighbor] {
			newPath := make([]int, len(path))
			copy(newPath, path)
			newPath = append(newPath, neighbor)
			fb.dfsPath(fp, neighbor, newPath, visited, remaining-1)
		}
	}

	visited[current] = false
}

// hashPath creates a hash for a path of atoms
func (fb *MoleculeFingerprintBuilder) hashPath(path []int) uint32 {
	h := fnv.New32a()

	for i, atomIdx := range path {
		atom := &fb.mol.Atoms[atomIdx]

		// Include atom number
		h.Write([]byte{byte(atom.Number)})

		// Include bond order if not first atom
		if i > 0 {
			bondIdx := fb.mol.FindBond(path[i-1], atomIdx)
			if bondIdx >= 0 {
				h.Write([]byte{byte(fb.mol.GetBondOrder(bondIdx))})
			}
		}
	}

	return h.Sum32()
}

// buildECFP generates an ECFP (Extended Connectivity Fingerprint)
func (fb *MoleculeFingerprintBuilder) buildECFP() *Fingerprint {
	fp := NewFingerprint(fb.params)

	radius := fb.getECFPRadius()

	// Initialize atom identifiers
	atomIdentifiers := make([]uint32, len(fb.mol.Atoms))
	for i := range fb.mol.Atoms {
		atomIdentifiers[i] = fb.getInitialAtomIdentifier(i)
	}

	// Iterate for each radius
	for r := 0; r <= radius; r++ {
		// Set bits for current identifiers
		for _, identifier := range atomIdentifiers {
			fp.SetBitsFromHash(identifier, 2)
		}

		// Update identifiers for next iteration
		if r < radius {
			newIdentifiers := make([]uint32, len(fb.mol.Atoms))
			for i := range fb.mol.Atoms {
				newIdentifiers[i] = fb.updateECFPIdentifier(i, atomIdentifiers)
			}
			atomIdentifiers = newIdentifiers
		}
	}

	return fp
}

// getECFPRadius returns the radius for ECFP fingerprints
func (fb *MoleculeFingerprintBuilder) getECFPRadius() int {
	switch fb.params.Type {
	case FingerprintECFP2:
		return 1
	case FingerprintECFP4:
		return 2
	case FingerprintECFP6:
		return 3
	default:
		return 2
	}
}

// getInitialAtomIdentifier returns initial identifier for an atom
func (fb *MoleculeFingerprintBuilder) getInitialAtomIdentifier(atomIdx int) uint32 {
	atom := &fb.mol.Atoms[atomIdx]

	h := fnv.New32a()

	// Atom number
	h.Write([]byte{byte(atom.Number & 0xFF)})
	h.Write([]byte{byte((atom.Number >> 8) & 0xFF)})

	// Number of heavy atom neighbors
	neighbors := fb.mol.GetNeighbors(atomIdx)
	heavyCount := 0
	for _, n := range neighbors {
		if fb.mol.Atoms[n].Number != ELEM_H {
			heavyCount++
		}
	}
	h.Write([]byte{byte(heavyCount)})

	// Total connectivity
	totalH := fb.mol.GetImplicitH(atomIdx)
	for _, n := range neighbors {
		if fb.mol.Atoms[n].Number == ELEM_H {
			totalH++
		}
	}
	h.Write([]byte{byte(len(neighbors) + totalH)})

	// Charge
	h.Write([]byte{byte(atom.Charge + 128)}) // Offset to handle negative

	return h.Sum32()
}

// updateECFPIdentifier updates atom identifier based on neighbors
func (fb *MoleculeFingerprintBuilder) updateECFPIdentifier(atomIdx int, currentIdentifiers []uint32) uint32 {
	h := fnv.New32a()

	// Include current identifier
	id := currentIdentifiers[atomIdx]
	h.Write([]byte{
		byte(id & 0xFF),
		byte((id >> 8) & 0xFF),
		byte((id >> 16) & 0xFF),
		byte((id >> 24) & 0xFF),
	})

	// Get neighbor identifiers and sort them
	neighbors := fb.mol.GetNeighbors(atomIdx)
	neighborIds := make([]uint32, len(neighbors))
	for i, n := range neighbors {
		neighborIds[i] = currentIdentifiers[n]
	}

	// Simple sort (bubble sort for small arrays)
	for i := 0; i < len(neighborIds); i++ {
		for j := i + 1; j < len(neighborIds); j++ {
			if neighborIds[i] > neighborIds[j] {
				neighborIds[i], neighborIds[j] = neighborIds[j], neighborIds[i]
			}
		}
	}

	// Include sorted neighbor identifiers
	for _, nid := range neighborIds {
		h.Write([]byte{
			byte(nid & 0xFF),
			byte((nid >> 8) & 0xFF),
			byte((nid >> 16) & 0xFF),
			byte((nid >> 24) & 0xFF),
		})
	}

	return h.Sum32()
}

// Similarity calculations

// TanimotoSimilarity calculates Tanimoto coefficient between two fingerprints
func TanimotoSimilarity(fp1, fp2 *Fingerprint) float64 {
	if len(fp1.Bits) != len(fp2.Bits) {
		return 0.0
	}

	intersection := 0
	union := 0

	for i := range fp1.Bits {
		// Count bits in intersection (AND)
		intersection += bits.OnesCount64(fp1.Bits[i] & fp2.Bits[i])
		// Count bits in union (OR)
		union += bits.OnesCount64(fp1.Bits[i] | fp2.Bits[i])
	}

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// DiceSimilarity calculates Dice coefficient between two fingerprints
func DiceSimilarity(fp1, fp2 *Fingerprint) float64 {
	if len(fp1.Bits) != len(fp2.Bits) {
		return 0.0
	}

	intersection := 0
	count1 := 0
	count2 := 0

	for i := range fp1.Bits {
		intersection += bits.OnesCount64(fp1.Bits[i] & fp2.Bits[i])
		count1 += bits.OnesCount64(fp1.Bits[i])
		count2 += bits.OnesCount64(fp2.Bits[i])
	}

	if count1+count2 == 0 {
		return 0.0
	}

	return 2.0 * float64(intersection) / float64(count1+count2)
}

// CosineSimilarity calculates cosine similarity between two fingerprints
func CosineSimilarity(fp1, fp2 *Fingerprint) float64 {
	if len(fp1.Bits) != len(fp2.Bits) {
		return 0.0
	}

	intersection := 0
	count1 := 0
	count2 := 0

	for i := range fp1.Bits {
		intersection += bits.OnesCount64(fp1.Bits[i] & fp2.Bits[i])
		count1 += bits.OnesCount64(fp1.Bits[i])
		count2 += bits.OnesCount64(fp2.Bits[i])
	}

	if count1 == 0 || count2 == 0 {
		return 0.0
	}

	return float64(intersection) / math.Sqrt(float64(count1)*float64(count2))
}

// GenerateFingerprint is a convenience function to generate fingerprint with default params
func GenerateFingerprint(mol *Molecule) *Fingerprint {
	params := DefaultFingerprintParameters()
	builder := NewFingerprintBuilder(mol, params)
	return builder.Build()
}

// GenerateFingerprintECFP4 generates an ECFP4 fingerprint
func GenerateFingerprintECFP4(mol *Molecule) *Fingerprint {
	params := FingerprintParameters{
		Type: FingerprintECFP4,
		Size: 2048,
	}
	builder := NewFingerprintBuilder(mol, params)
	return builder.Build()
}

// ToHexString converts fingerprint to hex string for storage/display
func (fp *Fingerprint) ToHexString() string {
	result := make([]byte, 0, len(fp.Bits)*16)
	for _, word := range fp.Bits {
		for i := 0; i < 8; i++ {
			b := byte((word >> (i * 8)) & 0xFF)
			result = append(result, hexChar(b>>4), hexChar(b&0x0F))
		}
	}
	return string(result)
}

func hexChar(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	return 'a' + (n - 10)
}

// FingerprintDistance calculates various distance metrics
type FingerprintDistance struct{}

// HammingDistance calculates Hamming distance (number of differing bits)
func (fd *FingerprintDistance) HammingDistance(fp1, fp2 *Fingerprint) int {
	if len(fp1.Bits) != len(fp2.Bits) {
		return -1
	}

	distance := 0
	for i := range fp1.Bits {
		distance += bits.OnesCount64(fp1.Bits[i] ^ fp2.Bits[i])
	}

	return distance
}

// EuclideanDistance calculates Euclidean distance
func (fd *FingerprintDistance) EuclideanDistance(fp1, fp2 *Fingerprint) float64 {
	hamming := float64(fd.HammingDistance(fp1, fp2))
	return math.Sqrt(hamming)
}
