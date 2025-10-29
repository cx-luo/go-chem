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

import (
	"errors"
)

// BaseReaction represents a base chemical reaction
type BaseReaction struct {
	allMolecules      []*BaseMolecule
	reactionBlocks    []*ReactionBlock
	types             []int
	specialConditions []SpecialCondition
	properties        *PropertiesMap
	meta              *MetaDataStorage
	name              string
	originalFormat    int
	isRetrosynthetic  bool
	reactantCount     int
	productCount      int
	catalystCount     int
	intermediateCount int
	undefinedCount    int
	specialCount      int
}

// NewBaseReaction creates a new BaseReaction
func NewBaseReaction() *BaseReaction {
	return &BaseReaction{
		allMolecules:      make([]*BaseMolecule, 0),
		reactionBlocks:    make([]*ReactionBlock, 0),
		types:             make([]int, 0),
		specialConditions: make([]SpecialCondition, 0),
		properties:        NewPropertiesMap(),
		meta:              NewMetaDataStorage(),
		name:              "",
		originalFormat:    0,
		isRetrosynthetic:  false,
		reactantCount:     0,
		productCount:      0,
		catalystCount:     0,
		intermediateCount: 0,
		undefinedCount:    0,
		specialCount:      0,
	}
}

// Clear clears all data in the reaction
func (br *BaseReaction) Clear() {
	br.reactantCount = 0
	br.productCount = 0
	br.catalystCount = 0
	br.intermediateCount = 0
	br.undefinedCount = 0
	br.specialCount = 0
	br.allMolecules = make([]*BaseMolecule, 0)
	br.reactionBlocks = make([]*ReactionBlock, 0)
	br.specialConditions = make([]SpecialCondition, 0)
	br.types = make([]int, 0)
	br.name = ""
}

// Meta returns the metadata storage
func (br *BaseReaction) Meta() *MetaDataStorage {
	return br.meta
}

// Properties returns the properties map
func (br *BaseReaction) Properties() *PropertiesMap {
	return br.properties
}

// Begin returns the beginning index for iteration
func (br *BaseReaction) Begin() int {
	return br.nextElement(Reactant|Product|Catalyst|Intermediate|Undefined, -1)
}

// End returns the ending index for iteration
func (br *BaseReaction) End() int {
	return len(br.allMolecules)
}

// Next returns the next index for iteration
func (br *BaseReaction) Next(i int) int {
	return br.nextElement(Reactant|Product|Catalyst|Intermediate|Undefined, i)
}

// Count returns the total number of molecules
func (br *BaseReaction) Count() int {
	return len(br.allMolecules)
}

// Remove removes a molecule at the given index
func (br *BaseReaction) Remove(i int) error {
	if i < 0 || i >= len(br.allMolecules) {
		return errors.New("index out of range")
	}

	side := br.types[i]
	switch side {
	case Reactant:
		br.reactantCount--
	case Product:
		br.productCount--
	case Intermediate:
		br.intermediateCount--
	case Undefined:
		br.undefinedCount--
	case Catalyst:
		br.catalystCount--
	}

	// Remove molecule
	br.allMolecules = append(br.allMolecules[:i], br.allMolecules[i+1:]...)
	br.types = append(br.types[:i], br.types[i+1:]...)

	return nil
}

// Molecules returns all molecules
func (br *BaseReaction) Molecules() []*BaseMolecule {
	return br.allMolecules
}

// IntermediateBegin returns the beginning index for intermediates
func (br *BaseReaction) IntermediateBegin() int {
	return br.nextElement(Intermediate, -1)
}

// IntermediateNext returns the next intermediate index
func (br *BaseReaction) IntermediateNext(index int) int {
	return br.nextElement(Intermediate, index)
}

// IntermediateEnd returns the ending index for intermediates
func (br *BaseReaction) IntermediateEnd() int {
	return br.End()
}

// UndefinedBegin returns the beginning index for undefined molecules
func (br *BaseReaction) UndefinedBegin() int {
	return br.nextElement(Undefined, -1)
}

// UndefinedNext returns the next undefined index
func (br *BaseReaction) UndefinedNext(index int) int {
	return br.nextElement(Undefined, index)
}

// UndefinedEnd returns the ending index for undefined molecules
func (br *BaseReaction) UndefinedEnd() int {
	return br.End()
}

// ReactantBegin returns the beginning index for reactants
func (br *BaseReaction) ReactantBegin() int {
	return br.nextElement(Reactant, -1)
}

// ReactantNext returns the next reactant index
func (br *BaseReaction) ReactantNext(index int) int {
	return br.nextElement(Reactant, index)
}

// ReactantEnd returns the ending index for reactants
func (br *BaseReaction) ReactantEnd() int {
	return br.End()
}

// ProductBegin returns the beginning index for products
func (br *BaseReaction) ProductBegin() int {
	return br.nextElement(Product, -1)
}

// ProductNext returns the next product index
func (br *BaseReaction) ProductNext(index int) int {
	return br.nextElement(Product, index)
}

// ProductEnd returns the ending index for products
func (br *BaseReaction) ProductEnd() int {
	return br.End()
}

// CatalystBegin returns the beginning index for catalysts
func (br *BaseReaction) CatalystBegin() int {
	return br.nextElement(Catalyst, -1)
}

// CatalystNext returns the next catalyst index
func (br *BaseReaction) CatalystNext(index int) int {
	return br.nextElement(Catalyst, index)
}

// CatalystEnd returns the ending index for catalysts
func (br *BaseReaction) CatalystEnd() int {
	return br.End()
}

// SideBegin returns the beginning index for a given side
func (br *BaseReaction) SideBegin(side int) int {
	return br.nextElement(side, -1)
}

// SideNext returns the next index for a given side
func (br *BaseReaction) SideNext(side, index int) int {
	return br.nextElement(side, index)
}

// SideEnd returns the ending index for sides
func (br *BaseReaction) SideEnd() int {
	return br.End()
}

// GetSideType returns the side type for a molecule at the given index
func (br *BaseReaction) GetSideType(index int) int {
	if index < 0 || index >= len(br.types) {
		return 0
	}
	return br.types[index]
}

// UndefinedCount returns the count of undefined molecules
func (br *BaseReaction) UndefinedCount() int {
	return br.undefinedCount
}

// IntermediateCount returns the count of intermediate molecules
func (br *BaseReaction) IntermediateCount() int {
	return br.intermediateCount
}

// ReactantsCount returns the count of reactant molecules
func (br *BaseReaction) ReactantsCount() int {
	return br.reactantCount
}

// ProductsCount returns the count of product molecules
func (br *BaseReaction) ProductsCount() int {
	return br.productCount
}

// CatalystCount returns the count of catalyst molecules
func (br *BaseReaction) CatalystCount() int {
	return br.catalystCount
}

// SpecialConditionsCount returns the count of special conditions
func (br *BaseReaction) SpecialConditionsCount() int {
	return len(br.specialConditions)
}

// ReactionBlocksCount returns the count of reaction blocks
func (br *BaseReaction) ReactionBlocksCount() int {
	return len(br.reactionBlocks)
}

// ReactionBlock returns the reaction block at the given index
func (br *BaseReaction) ReactionBlock(index int) *ReactionBlock {
	if index < 0 || index >= len(br.reactionBlocks) {
		return nil
	}
	return br.reactionBlocks[index]
}

// AddReactionBlock adds a new reaction block
func (br *BaseReaction) AddReactionBlock() *ReactionBlock {
	rb := NewReactionBlock()
	br.reactionBlocks = append(br.reactionBlocks, rb)
	return rb
}

// ClearReactionBlocks clears all reaction blocks
func (br *BaseReaction) ClearReactionBlocks() {
	br.reactionBlocks = make([]*ReactionBlock, 0)
}

// ReactionsCount returns the count of reactions
func (br *BaseReaction) ReactionsCount() int {
	return len(br.reactionBlocks)
}

// ReactionBegin returns the beginning index for reactions
func (br *BaseReaction) ReactionBegin() int {
	for i := 0; i < len(br.reactionBlocks); i++ {
		rb := br.reactionBlocks[i]
		if len(rb.Products) > 0 || len(rb.Reactants) > 0 {
			return i
		}
	}
	return len(br.reactionBlocks)
}

// ReactionEnd returns the ending index for reactions
func (br *BaseReaction) ReactionEnd() int {
	if len(br.reactionBlocks) == 0 {
		return 1
	}
	return len(br.reactionBlocks)
}

// ReactionNext returns the next reaction index
func (br *BaseReaction) ReactionNext(i int) int {
	for i++; i < len(br.reactionBlocks); i++ {
		rb := br.reactionBlocks[i]
		if len(rb.Products) > 0 || len(rb.Reactants) > 0 {
			break
		}
	}
	return i
}

// GetAAM returns the atom-atom mapping for a given molecule and atom
func (br *BaseReaction) GetAAM(index, atom int) int {
	if index < 0 || index >= len(br.allMolecules) {
		return 0
	}
	mol := br.allMolecules[index]
	if atom < 0 || atom >= len(mol.ReactionAtomMapping) {
		return 0
	}
	return mol.ReactionAtomMapping[atom]
}

// GetReactingCenter returns the reacting center for a given molecule and bond
func (br *BaseReaction) GetReactingCenter(index, bond int) int {
	if index < 0 || index >= len(br.allMolecules) {
		return 0
	}
	mol := br.allMolecules[index]
	if bond < 0 || bond >= len(mol.ReactionBondReactingCenter) {
		return 0
	}
	return mol.ReactionBondReactingCenter[bond]
}

// GetInversion returns the inversion for a given molecule and atom
func (br *BaseReaction) GetInversion(index, atom int) int {
	if index < 0 || index >= len(br.allMolecules) {
		return 0
	}
	mol := br.allMolecules[index]
	if atom < 0 || atom >= len(mol.ReactionAtomInversion) {
		return 0
	}
	return mol.ReactionAtomInversion[atom]
}

// GetAAMArray returns the atom-atom mapping array for a molecule
func (br *BaseReaction) GetAAMArray(index int) []int {
	if index < 0 || index >= len(br.allMolecules) {
		return nil
	}
	return br.allMolecules[index].ReactionAtomMapping
}

// GetReactingCenterArray returns the reacting center array for a molecule
func (br *BaseReaction) GetReactingCenterArray(index int) []int {
	if index < 0 || index >= len(br.allMolecules) {
		return nil
	}
	return br.allMolecules[index].ReactionBondReactingCenter
}

// GetInversionArray returns the inversion array for a molecule
func (br *BaseReaction) GetInversionArray(index int) []int {
	if index < 0 || index >= len(br.allMolecules) {
		return nil
	}
	return br.allMolecules[index].ReactionAtomInversion
}

// ClearAAM clears all atom-atom mappings
func (br *BaseReaction) ClearAAM() {
	for i := br.Begin(); i != br.End(); i = br.Next(i) {
		mol := br.allMolecules[i]
		for j := range mol.ReactionAtomMapping {
			mol.ReactionAtomMapping[j] = 0
		}
	}
}

// AddReactant adds a new reactant molecule
func (br *BaseReaction) AddReactant() int {
	return br.addBaseMolecule(Reactant)
}

// AddProduct adds a new product molecule
func (br *BaseReaction) AddProduct() int {
	return br.addBaseMolecule(Product)
}

// AddCatalyst adds a new catalyst molecule
func (br *BaseReaction) AddCatalyst() int {
	return br.addBaseMolecule(Catalyst)
}

// AddIntermediate adds a new intermediate molecule
func (br *BaseReaction) AddIntermediate() int {
	return br.addBaseMolecule(Intermediate)
}

// AddUndefined adds a new undefined molecule
func (br *BaseReaction) AddUndefined() int {
	return br.addBaseMolecule(Undefined)
}

// SpecialConditionCount returns the count of special conditions
func (br *BaseReaction) SpecialConditionCount() int {
	return len(br.specialConditions)
}

// AddSpecialCondition adds a special condition
func (br *BaseReaction) AddSpecialCondition(metaIdx int, bbox Rect2f) int {
	sc := NewSpecialCondition(metaIdx, bbox)
	br.specialConditions = append(br.specialConditions, sc)
	return len(br.specialConditions) - 1
}

// ClearSpecialConditions clears all special conditions
func (br *BaseReaction) ClearSpecialConditions() {
	br.specialConditions = make([]SpecialCondition, 0)
}

// SpecialCondition returns the special condition at the given index
func (br *BaseReaction) SpecialCondition(idx int) (SpecialCondition, error) {
	if idx < 0 || idx >= len(br.specialConditions) {
		return SpecialCondition{}, errors.New("index out of range")
	}
	return br.specialConditions[idx], nil
}

// addBaseMolecule adds a base molecule (internal method)
func (br *BaseReaction) addBaseMolecule(side int) int {
	mol := NewBaseMolecule()
	idx := len(br.allMolecules)
	br.allMolecules = append(br.allMolecules, mol)
	br.addedBaseMolecule(idx, side, mol)
	return idx
}

// addedBaseMolecule is called after a molecule is added
func (br *BaseReaction) addedBaseMolecule(idx, side int, mol *BaseMolecule) {
	switch side {
	case Reactant:
		br.reactantCount++
	case Product:
		br.productCount++
	case Intermediate:
		br.intermediateCount++
	case Undefined:
		br.undefinedCount++
	case Catalyst:
		br.catalystCount++
	}

	// Expand types array if needed
	for len(br.types) <= idx {
		br.types = append(br.types, 0)
	}
	br.types[idx] = side
}

// AddReactantCopy adds a copy of a reactant molecule
func (br *BaseReaction) AddReactantCopy(mol *BaseMolecule, mapping, invMapping *[]int) int {
	idx := len(br.allMolecules)
	newMol := mol.Neu()
	newMol.Clone(mol, mapping, invMapping)
	br.allMolecules = append(br.allMolecules, newMol)
	br.addedBaseMolecule(idx, Reactant, newMol)
	return idx
}

// AddProductCopy adds a copy of a product molecule
func (br *BaseReaction) AddProductCopy(mol *BaseMolecule, mapping, invMapping *[]int) int {
	idx := len(br.allMolecules)
	newMol := mol.Neu()
	newMol.Clone(mol, mapping, invMapping)
	br.allMolecules = append(br.allMolecules, newMol)
	br.addedBaseMolecule(idx, Product, newMol)
	return idx
}

// AddCatalystCopy adds a copy of a catalyst molecule
func (br *BaseReaction) AddCatalystCopy(mol *BaseMolecule, mapping, invMapping *[]int) int {
	idx := len(br.allMolecules)
	newMol := mol.Neu()
	newMol.Clone(mol, mapping, invMapping)
	br.allMolecules = append(br.allMolecules, newMol)
	br.addedBaseMolecule(idx, Catalyst, newMol)
	return idx
}

// AddIntermediateCopy adds a copy of an intermediate molecule
func (br *BaseReaction) AddIntermediateCopy(mol *BaseMolecule, mapping, invMapping *[]int) int {
	idx := len(br.allMolecules)
	newMol := mol.Neu()
	newMol.Clone(mol, mapping, invMapping)
	br.allMolecules = append(br.allMolecules, newMol)
	br.addedBaseMolecule(idx, Intermediate, newMol)
	return idx
}

// AddUndefinedCopy adds a copy of an undefined molecule
func (br *BaseReaction) AddUndefinedCopy(mol *BaseMolecule, mapping, invMapping *[]int) int {
	idx := len(br.allMolecules)
	newMol := mol.Neu()
	newMol.Clone(mol, mapping, invMapping)
	br.allMolecules = append(br.allMolecules, newMol)
	br.addedBaseMolecule(idx, Undefined, newMol)
	return idx
}

// FindAtomByAAM finds an atom by its atom-atom mapping number
func (br *BaseReaction) FindAtomByAAM(molIdx, aam int) int {
	if molIdx < 0 || molIdx >= len(br.allMolecules) {
		return -1
	}

	mol := br.allMolecules[molIdx]
	for i := mol.VertexBegin(); i < mol.VertexEnd(); i = mol.VertexNext(i) {
		if br.GetAAM(molIdx, i) == aam {
			return i
		}
	}
	return -1
}

// FindAamNumber finds the AAM number for a given molecule and atom
func (br *BaseReaction) FindAamNumber(mol *BaseMolecule, atomNumber int) (int, error) {
	for i := br.Begin(); i < br.End(); i = br.Next(i) {
		if br.allMolecules[i] == mol {
			return br.GetAAM(i, atomNumber), nil
		}
	}
	return 0, errors.New("cannot find aam number")
}

// FindReactingCenter finds the reacting center for a given molecule and bond
func (br *BaseReaction) FindReactingCenter(mol *BaseMolecule, bondNumber int) (int, error) {
	for i := br.Begin(); i < br.End(); i = br.Next(i) {
		if br.allMolecules[i] == mol {
			return br.GetReactingCenter(i, bondNumber), nil
		}
	}
	return 0, errors.New("cannot find reacting center")
}

// FindMolecule finds the index of a molecule
func (br *BaseReaction) FindMolecule(mol *BaseMolecule) int {
	for i := br.Begin(); i != br.End(); i = br.Next(i) {
		if br.GetBaseMolecule(i) == mol {
			return i
		}
	}
	return -1
}

// MarkStereocenterBonds marks stereocenter bonds
func (br *BaseReaction) MarkStereocenterBonds() {
	for i := br.Begin(); i < br.End(); i = br.Next(i) {
		br.allMolecules[i].ClearBondDirections()
		br.allMolecules[i].MarkBondsStereocenters()
		br.allMolecules[i].MarkBondsAlleneStereo()
	}
}

// HaveCoord checks if all molecules have coordinates
func HaveCoord(reaction *BaseReaction) bool {
	for i := reaction.Begin(); i < reaction.End(); i = reaction.Next(i) {
		if !reaction.GetBaseMolecule(i).HaveXYZ {
			return false
		}
	}
	return true
}

// HasSelection checks if any molecule has selection
func (br *BaseReaction) HasSelection() bool {
	for i := br.Begin(); i < br.End(); i = br.Next(i) {
		if br.GetBaseMolecule(i).HasSelection() {
			return true
		}
	}
	return false
}

// GetBaseMolecule returns the base molecule at the given index
func (br *BaseReaction) GetBaseMolecule(index int) *BaseMolecule {
	if index < 0 || index >= len(br.allMolecules) {
		return nil
	}
	return br.allMolecules[index]
}

// nextElement finds the next element of a given type
func (br *BaseReaction) nextElement(typ, index int) int {
	if index == -1 {
		index = 0
	} else {
		index++
	}

	for ; index < len(br.allMolecules); index++ {
		if index < len(br.types) && (br.types[index]&typ) != 0 {
			break
		}
	}
	return index
}

// IsRetrosyntetic checks if the reaction is retrosynthetic
func (br *BaseReaction) IsRetrosyntetic() bool {
	return br.isRetrosynthetic
}

// SetIsRetrosyntetic sets the reaction as retrosynthetic
func (br *BaseReaction) SetIsRetrosyntetic() {
	br.isRetrosynthetic = true
}

// MultitaleCount returns the count of multitale objects
func (br *BaseReaction) MultitaleCount() int {
	return br.meta.GetMetaCount("ReactionMultitailArrowObject")
}

// UnfoldHydrogens unfolds hydrogens in all molecules
func (br *BaseReaction) UnfoldHydrogens() {
	markers := make([]int, 0)
	for i := br.Begin(); i != br.End(); i = br.Next(i) {
		mol := br.GetBaseMolecule(i)
		mol.UnfoldHydrogens(&markers, -1)
	}
}

// Clone creates a deep copy of the reaction
func (br *BaseReaction) Clone(other *BaseReaction, molMapping *[]int) {
	br.Clear()

	if molMapping != nil {
		*molMapping = make([]int, other.End())
		for i := range *molMapping {
			(*molMapping)[i] = -1
		}
	}

	index := 0
	for i := 0; i < len(other.allMolecules); i++ {
		rmol := other.allMolecules[i]

		switch other.types[i] {
		case Reactant:
			index = br.AddReactantCopy(rmol, nil, nil)
		case Product:
			index = br.AddProductCopy(rmol, nil, nil)
		case Catalyst:
			index = br.AddCatalystCopy(rmol, nil, nil)
		case Intermediate:
			index = br.AddIntermediateCopy(rmol, nil, nil)
		case Undefined:
			index = br.AddUndefinedCopy(rmol, nil, nil)
		}

		if molMapping != nil {
			(*molMapping)[i] = index
		}
	}

	br.name = other.name
	br.meta.Clone(other.meta)
	br.properties.Copy(other.properties)
	br.isRetrosynthetic = other.isRetrosynthetic
}
