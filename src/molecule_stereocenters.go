// Package src provides molecular structure manipulation and analysis tools.
// This file implements stereocenter detection and handling.
package src

import (
	"fmt"
	"math"
)

// Stereocenter type constants
const (
	STEREO_ATOM_ANY = 1 // Any stereoconfiguration
	STEREO_ATOM_AND = 2 // AND (racemic)
	STEREO_ATOM_OR  = 3 // OR (relative)
	STEREO_ATOM_ABS = 4 // Absolute stereochemistry
)

// Stereocenter represents a stereogenic center in a molecule
type Stereocenter struct {
	AtomIdx         int    // Index of the stereocenter atom
	Type            int    // Stereocenter type (ANY, AND, OR, ABS)
	Group           int    // Stereogroup number
	Pyramid         [4]int // Four substituents in specific order (-1 for implicit H)
	IsAtropisomeric bool   // Whether this is an atropisomeric center
	IsTetrahydral   bool   // Whether this is a tetrahedral center
}

// MoleculeStereocenters manages all stereocenters in a molecule
type MoleculeStereocenters struct {
	centers map[int]*Stereocenter // Map from atom index to stereocenter
}

// NewMoleculeStereocenters creates a new empty stereocenter collection
func NewMoleculeStereocenters() *MoleculeStereocenters {
	return &MoleculeStereocenters{
		centers: make(map[int]*Stereocenter),
	}
}

// Clear removes all stereocenters
func (ms *MoleculeStereocenters) Clear() {
	ms.centers = make(map[int]*Stereocenter)
}

// Size returns the number of stereocenters
func (ms *MoleculeStereocenters) Size() int {
	return len(ms.centers)
}

// Exists checks if an atom is a stereocenter
func (ms *MoleculeStereocenters) Exists(atomIdx int) bool {
	_, ok := ms.centers[atomIdx]
	return ok
}

// Add adds a new stereocenter
func (ms *MoleculeStereocenters) Add(atomIdx, stereoType, group int, pyramid [4]int) {
	ms.centers[atomIdx] = &Stereocenter{
		AtomIdx:         atomIdx,
		Type:            stereoType,
		Group:           group,
		Pyramid:         pyramid,
		IsAtropisomeric: false,
		IsTetrahydral:   true,
	}
}

// AddWithInversion adds a stereocenter with optionally inverted pyramid
func (ms *MoleculeStereocenters) AddWithInversion(atomIdx, stereoType, group int, invertPyramid bool, neighbors []int) {
	pyramid := [4]int{-1, -1, -1, -1}

	// Fill pyramid with neighbors
	for i := 0; i < len(neighbors) && i < 4; i++ {
		pyramid[i] = neighbors[i]
	}

	// If we have less than 4 neighbors, add implicit hydrogen
	if len(neighbors) < 4 {
		pyramid[len(neighbors)] = -1
	}

	if invertPyramid {
		// Invert by swapping first two substituents
		pyramid[0], pyramid[1] = pyramid[1], pyramid[0]
	}

	ms.Add(atomIdx, stereoType, group, pyramid)
}

// Remove removes a stereocenter
func (ms *MoleculeStereocenters) Remove(atomIdx int) {
	delete(ms.centers, atomIdx)
}

// Get retrieves a stereocenter by atom index
func (ms *MoleculeStereocenters) Get(atomIdx int) (*Stereocenter, error) {
	if center, ok := ms.centers[atomIdx]; ok {
		return center, nil
	}
	return nil, fmt.Errorf("atom %d is not a stereocenter", atomIdx)
}

// GetType returns the type of a stereocenter
func (ms *MoleculeStereocenters) GetType(atomIdx int) int {
	if center, ok := ms.centers[atomIdx]; ok {
		return center.Type
	}
	return -1
}

// GetGroup returns the stereogroup of a stereocenter
func (ms *MoleculeStereocenters) GetGroup(atomIdx int) int {
	if center, ok := ms.centers[atomIdx]; ok {
		return center.Group
	}
	return -1
}

// SetType sets the type of a stereocenter
func (ms *MoleculeStereocenters) SetType(atomIdx, stereoType int) {
	if center, ok := ms.centers[atomIdx]; ok {
		center.Type = stereoType
	}
}

// SetGroup sets the stereogroup of a stereocenter
func (ms *MoleculeStereocenters) SetGroup(atomIdx, group int) {
	if center, ok := ms.centers[atomIdx]; ok {
		center.Group = group
	}
}

// GetPyramid returns the pyramid substituents of a stereocenter
func (ms *MoleculeStereocenters) GetPyramid(atomIdx int) [4]int {
	if center, ok := ms.centers[atomIdx]; ok {
		return center.Pyramid
	}
	return [4]int{-1, -1, -1, -1}
}

// InvertPyramid inverts the stereochemistry by swapping two substituents
func (ms *MoleculeStereocenters) InvertPyramid(atomIdx int) {
	if center, ok := ms.centers[atomIdx]; ok {
		center.Pyramid[0], center.Pyramid[1] = center.Pyramid[1], center.Pyramid[0]
	}
}

// IsPossibleStereocenter checks if an atom can be a stereocenter
func (ms *MoleculeStereocenters) IsPossibleStereocenter(mol *Molecule, atomIdx int) (bool, bool, bool) {
	if atomIdx < 0 || atomIdx >= len(mol.Atoms) {
		return false, false, false
	}

	atom := &mol.Atoms[atomIdx]
	neighbors := mol.GetNeighbors(atomIdx)

	// Can't be a stereocenter with < 3 neighbors
	if len(neighbors) < 3 {
		return false, false, false
	}

	// Check for tetrahedral centers (sp3 carbon, nitrogen, etc.)
	possibleImplicitH := false
	possibleLonePair := false

	if atom.Number == ELEM_C {
		// Carbon with 3 or 4 different substituents
		if len(neighbors) == 3 {
			possibleImplicitH = true
		}
		if len(neighbors) >= 3 && len(neighbors) <= 4 {
			if ms.hasDifferentSubstituents(mol, atomIdx, neighbors) {
				return true, possibleImplicitH, possibleLonePair
			}
		}
	} else if atom.Number == ELEM_N {
		// Nitrogen with 3 substituents or 3 + lone pair
		if len(neighbors) == 3 {
			if ms.hasDifferentSubstituents(mol, atomIdx, neighbors) {
				possibleLonePair = true
				return true, false, possibleLonePair
			}
		}
	} else if atom.Number == ELEM_S || atom.Number == ELEM_P {
		// Sulfur or phosphorus with 3-4 different substituents
		if len(neighbors) >= 3 && len(neighbors) <= 4 {
			if ms.hasDifferentSubstituents(mol, atomIdx, neighbors) {
				return true, possibleImplicitH, possibleLonePair
			}
		}
	}

	return false, possibleImplicitH, possibleLonePair
}

// hasDifferentSubstituents checks if all substituents are different (simplified)
func (ms *MoleculeStereocenters) hasDifferentSubstituents(mol *Molecule, atomIdx int, neighbors []int) bool {
	// This is a simplified check - a full implementation would use
	// canonical ranking of substituents
	if len(neighbors) < 3 {
		return false
	}

	// Check if all neighbors have different atomic numbers as a quick test
	seen := make(map[int]bool)
	for _, n := range neighbors {
		num := mol.Atoms[n].Number
		if seen[num] && len(neighbors) == len(seen)+1 {
			// Duplicate and we have exactly one duplicate
			return false
		}
		seen[num] = true
	}

	// If all different atomic numbers, likely different substituents
	return len(seen) >= 3
}

// BuildFromBonds attempts to determine stereocenters from bond directions
func (ms *MoleculeStereocenters) BuildFromBonds(mol *Molecule, sensibleBondsOut *[]int) {
	ms.Clear()

	// Iterate through all atoms looking for potential stereocenters
	for i := range mol.Atoms {
		isPossible, _, _ := ms.IsPossibleStereocenter(mol, i)
		if !isPossible {
			continue
		}

		// Check for bond directions indicating stereochemistry
		neighbors := mol.GetNeighbors(i)
		bonds := mol.GetNeighborBonds(i)

		hasUp := false
		hasDown := false

		for _, bondIdx := range bonds {
			dir := mol.GetBondDirection(bondIdx)
			if dir == BOND_UP {
				hasUp = true
			} else if dir == BOND_DOWN {
				hasDown = true
			}
		}

		// If we have both up and down wedges, this is a defined stereocenter
		if hasUp && hasDown {
			// Determine pyramid from bond directions
			pyramid := ms.buildPyramidFromBonds(mol, i, neighbors, bonds)
			ms.Add(i, STEREO_ATOM_ABS, 1, pyramid)

			if sensibleBondsOut != nil {
				*sensibleBondsOut = append(*sensibleBondsOut, bonds...)
			}
		}
	}
}

// buildPyramidFromBonds constructs the pyramid ordering from bond directions
func (ms *MoleculeStereocenters) buildPyramidFromBonds(mol *Molecule, atomIdx int, neighbors []int, bonds []int) [4]int {
	pyramid := [4]int{-1, -1, -1, -1}

	// This is a simplified implementation
	// A full implementation would consider wedge directions and 3D geometry
	for i := 0; i < len(neighbors) && i < 4; i++ {
		pyramid[i] = neighbors[i]
	}

	return pyramid
}

// BuildFrom3DCoordinates determines stereocenters from 3D coordinates
func (ms *MoleculeStereocenters) BuildFrom3DCoordinates(mol *Molecule) {
	if !mol.HaveXYZ {
		return
	}

	ms.Clear()

	for i := range mol.Atoms {
		isPossible, _, _ := ms.IsPossibleStereocenter(mol, i)
		if !isPossible {
			continue
		}

		neighbors := mol.GetNeighbors(i)
		if len(neighbors) < 3 || len(neighbors) > 4 {
			continue
		}

		// Calculate chirality from 3D coordinates
		if ms.isChiralFrom3D(mol, i, neighbors) {
			pyramid := [4]int{-1, -1, -1, -1}
			for j := 0; j < len(neighbors); j++ {
				pyramid[j] = neighbors[j]
			}
			ms.Add(i, STEREO_ATOM_ABS, 1, pyramid)
		}
	}
}

// isChiralFrom3D determines if an atom is chiral based on 3D geometry
func (ms *MoleculeStereocenters) isChiralFrom3D(mol *Molecule, atomIdx int, neighbors []int) bool {
	if len(neighbors) < 4 {
		return false
	}

	center := mol.Atoms[atomIdx].Pos

	// Get vectors to all four neighbors
	v1 := subVec3f(mol.Atoms[neighbors[0]].Pos, center)
	v2 := subVec3f(mol.Atoms[neighbors[1]].Pos, center)
	v3 := subVec3f(mol.Atoms[neighbors[2]].Pos, center)
	v4 := subVec3f(mol.Atoms[neighbors[3]].Pos, center)

	// Calculate signed volume of tetrahedron
	// If volume is significantly non-zero, the center is chiral
	volume := dotVec3f(v1, crossVec3f(v2, v3))
	volume2 := dotVec3f(v1, crossVec3f(v2, v4))

	// If both volumes have the same sign and are non-zero, it's chiral
	return math.Abs(volume) > 0.01 && math.Abs(volume2) > 0.01
}

// Vector operations for 3D geometry

func subVec3f(a, b Vec3f) Vec3f {
	return Vec3f{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func crossVec3f(a, b Vec3f) Vec3f {
	return Vec3f{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func dotVec3f(a, b Vec3f) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// GetAbsAtoms returns indices of atoms with absolute stereochemistry
func (ms *MoleculeStereocenters) GetAbsAtoms() []int {
	var result []int
	for atomIdx, center := range ms.centers {
		if center.Type == STEREO_ATOM_ABS {
			result = append(result, atomIdx)
		}
	}
	return result
}

// GetOrGroups returns stereogroup numbers for OR groups
func (ms *MoleculeStereocenters) GetOrGroups() []int {
	groups := make(map[int]bool)
	for _, center := range ms.centers {
		if center.Type == STEREO_ATOM_OR {
			groups[center.Group] = true
		}
	}

	result := make([]int, 0, len(groups))
	for g := range groups {
		result = append(result, g)
	}
	return result
}

// GetAndGroups returns stereogroup numbers for AND groups
func (ms *MoleculeStereocenters) GetAndGroups() []int {
	groups := make(map[int]bool)
	for _, center := range ms.centers {
		if center.Type == STEREO_ATOM_AND {
			groups[center.Group] = true
		}
	}

	result := make([]int, 0, len(groups))
	for g := range groups {
		result = append(result, g)
	}
	return result
}

// HaveAbs checks if there are any absolute stereocenters
func (ms *MoleculeStereocenters) HaveAbs() bool {
	for _, center := range ms.centers {
		if center.Type == STEREO_ATOM_ABS {
			return true
		}
	}
	return false
}

// HaveAllAbs checks if all stereocenters are absolute
func (ms *MoleculeStereocenters) HaveAllAbs() bool {
	if len(ms.centers) == 0 {
		return false
	}
	for _, center := range ms.centers {
		if center.Type != STEREO_ATOM_ABS {
			return false
		}
	}
	return true
}

// Iterator methods

// GetAllCenters returns a slice of all stereocenter atom indices
func (ms *MoleculeStereocenters) GetAllCenters() []int {
	result := make([]int, 0, len(ms.centers))
	for atomIdx := range ms.centers {
		result = append(result, atomIdx)
	}
	return result
}
