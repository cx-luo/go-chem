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

// Molecule represents a concrete molecule
type Molecule struct {
	*BaseMolecule
}

// NewMolecule creates a new Molecule
func NewMolecule() *Molecule {
	return &Molecule{
		BaseMolecule: NewBaseMolecule(),
	}
}

// AsMolecule returns the molecule as a Molecule type
func (m *BaseMolecule) AsMolecule() *Molecule {
	return &Molecule{BaseMolecule: m}
}

// Reaction represents a chemical reaction
type Reaction struct {
	*BaseReaction
}

// NewReaction creates a new Reaction
func NewReaction() *Reaction {
	return &Reaction{
		BaseReaction: NewBaseReaction(),
	}
}

// Clear clears all data in the reaction
func (r *Reaction) Clear() {
	r.BaseReaction.Clear()
}

// GetMolecule returns the molecule at the given index
func (r *Reaction) GetMolecule(index int) *Molecule {
	return r.GetBaseMolecule(index).AsMolecule()
}

// Aromatize aromatizes bonds in all molecules
func (r *Reaction) Aromatize(options AromaticityOptions) bool {
	aromFound := false
	for i := r.Begin(); i < r.End(); i = r.Next(i) {
		mol := r.allMolecules[i]
		// Aromatize molecule bonds
		// aromFound |= MoleculeAromatizer.AromatizeBonds(mol, options)
		_ = mol // Placeholder
	}
	return aromFound
}

// AsReaction returns this reaction (implements interface)
func (r *Reaction) AsReaction() *Reaction {
	return r
}

// SaveBondOrders saves bond orders from a reaction
func SaveBondOrders(reaction *Reaction, bondTypes *[][]int) {
	for len(*bondTypes) < reaction.End() {
		*bondTypes = append(*bondTypes, make([]int, 0))
	}

	for i := reaction.Begin(); i != reaction.End(); i = reaction.Next(i) {
		mol := reaction.GetMolecule(i)
		// Save molecule bond orders
		_ = mol // Placeholder
	}
}

// LoadBondOrders loads bond orders to a reaction
func LoadBondOrders(reaction *Reaction, bondTypes [][]int) {
	for i := reaction.Begin(); i != reaction.End(); i = reaction.Next(i) {
		mol := reaction.GetMolecule(i)
		// Load molecule bond orders
		_ = mol // Placeholder
	}
}

// CheckForConsistency checks the reaction for consistency
func CheckForConsistency(rxn *Reaction) {
	for i := rxn.Begin(); i != rxn.End(); i = rxn.Next(i) {
		// Check molecule consistency
		_ = rxn.GetMolecule(i)
	}
}
