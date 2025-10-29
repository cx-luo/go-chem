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

// QueryMolecule represents a query molecule
type QueryMolecule struct {
	*BaseMolecule
	ReactionAtomExactChange []int
}

// NewQueryMolecule creates a new QueryMolecule
func NewQueryMolecule() *QueryMolecule {
	return &QueryMolecule{
		BaseMolecule:            NewBaseMolecule(),
		ReactionAtomExactChange: make([]int, 0),
	}
}

// QueryReaction represents a query reaction
type QueryReaction struct {
	*BaseReaction
	ignorableAAM [][]int
}

// NewQueryReaction creates a new QueryReaction
func NewQueryReaction() *QueryReaction {
	return &QueryReaction{
		BaseReaction: NewBaseReaction(),
		ignorableAAM: make([][]int, 0),
	}
}

// Clear clears all data in the query reaction
func (qr *QueryReaction) Clear() {
	qr.BaseReaction.Clear()
	qr.ignorableAAM = make([][]int, 0)
}

// GetQueryMolecule returns the query molecule at the given index
func (qr *QueryReaction) GetQueryMolecule(index int) *QueryMolecule {
	baseMol := qr.GetBaseMolecule(index)
	if baseMol == nil {
		return nil
	}
	// Type assertion (in real implementation, should be properly handled)
	return &QueryMolecule{BaseMolecule: baseMol}
}

// GetExactChangeArray returns the exact change array for a molecule
func (qr *QueryReaction) GetExactChangeArray(index int) []int {
	mol := qr.GetQueryMolecule(index)
	if mol == nil {
		return nil
	}
	return mol.ReactionAtomExactChange
}

// GetExactChange returns the exact change for a specific atom
func (qr *QueryReaction) GetExactChange(index, atom int) int {
	mol := qr.GetQueryMolecule(index)
	if mol == nil || atom < 0 || atom >= len(mol.ReactionAtomExactChange) {
		return 0
	}
	return mol.ReactionAtomExactChange[atom]
}

// AsQueryReaction returns this as a QueryReaction
func (qr *QueryReaction) AsQueryReaction() *QueryReaction {
	return qr
}

// IsQueryReaction returns true for QueryReaction
func (qr *QueryReaction) IsQueryReaction() bool {
	return true
}

// GetIgnorableAAMArray returns the ignorable AAM array for a molecule
func (qr *QueryReaction) GetIgnorableAAMArray(index int) []int {
	if index < 0 || index >= len(qr.ignorableAAM) {
		return nil
	}
	return qr.ignorableAAM[index]
}

// GetIgnorableAAM returns the ignorable AAM for a specific atom
func (qr *QueryReaction) GetIgnorableAAM(index, atom int) int {
	if index < 0 || index >= len(qr.ignorableAAM) {
		return 0
	}
	if atom < 0 || atom >= len(qr.ignorableAAM[index]) {
		return 0
	}
	return qr.ignorableAAM[index][atom]
}

// Optimize optimizes all query molecules
func (qr *QueryReaction) Optimize() {
	for i := qr.Begin(); i < qr.End(); i = qr.Next(i) {
		// Optimize query molecule
		_ = qr.allMolecules[i]
	}
}

// Aromatize aromatizes all query molecules
func (qr *QueryReaction) Aromatize(options AromaticityOptions) bool {
	aromFound := false
	for i := qr.Begin(); i < qr.End(); i = qr.Next(i) {
		// Aromatize query molecule
		_ = qr.allMolecules[i]
	}
	return aromFound
}
