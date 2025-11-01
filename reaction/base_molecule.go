/****************************************************************************
 * Copyright (C) from 2009 to Present EPAM Systems.
 *
 * This file is part of Indigo toolkit.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 ***************************************************************************/

package reaction

// BaseMolecule represents a base molecule structure
type BaseMolecule struct {
	ReactionAtomMapping        []int
	ReactionBondReactingCenter []int
	ReactionAtomInversion      []int
	HaveXYZ                    bool
	selected                   []bool
}

// NewBaseMolecule creates a new BaseMolecule
func NewBaseMolecule() *BaseMolecule {
	return &BaseMolecule{
		ReactionAtomMapping:        make([]int, 0),
		ReactionBondReactingCenter: make([]int, 0),
		ReactionAtomInversion:      make([]int, 0),
		HaveXYZ:                    false,
		selected:                   make([]bool, 0),
	}
}

// Clone creates a copy of the molecule
func (m *BaseMolecule) Clone(other *BaseMolecule, mapping, invMapping *[]int) {
	m.ReactionAtomMapping = make([]int, len(other.ReactionAtomMapping))
	copy(m.ReactionAtomMapping, other.ReactionAtomMapping)

	m.ReactionBondReactingCenter = make([]int, len(other.ReactionBondReactingCenter))
	copy(m.ReactionBondReactingCenter, other.ReactionBondReactingCenter)

	m.ReactionAtomInversion = make([]int, len(other.ReactionAtomInversion))
	copy(m.ReactionAtomInversion, other.ReactionAtomInversion)

	m.HaveXYZ = other.HaveXYZ

	m.selected = make([]bool, len(other.selected))
	copy(m.selected, other.selected)
}

// VertexBegin returns the beginning vertex index
func (m *BaseMolecule) VertexBegin() int {
	return 0
}

// VertexEnd returns the ending vertex index
func (m *BaseMolecule) VertexEnd() int {
	return len(m.ReactionAtomMapping)
}

// VertexNext returns the next vertex index
func (m *BaseMolecule) VertexNext(i int) int {
	return i + 1
}

// ClearBondDirections clears bond directions
func (m *BaseMolecule) ClearBondDirections() {
	// Implementation needed
}

// MarkBondsStereocenters marks bonds as stereocenters
func (m *BaseMolecule) MarkBondsStereocenters() {
	// Implementation needed
}

// MarkBondsAlleneStereo marks bonds as allene stereo
func (m *BaseMolecule) MarkBondsAlleneStereo() {
	// Implementation needed
}

// UnfoldHydrogens unfolds hydrogen atoms
func (m *BaseMolecule) UnfoldHydrogens(markers *[]int, param int) {
	// Implementation needed
}

// HasSelection returns true if the molecule has selected atoms
func (m *BaseMolecule) HasSelection() bool {
	for _, sel := range m.selected {
		if sel {
			return true
		}
	}
	return false
}

// Neu creates a new instance of the molecule
func (m *BaseMolecule) Neu() *BaseMolecule {
	return NewBaseMolecule()
}
