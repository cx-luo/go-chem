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
	"fmt"
)

// RxnfileSaver saves reactions to RXN file format
type RxnfileSaver struct {
	output        *Output
	rxn           *Reaction
	addStereoDesc bool
	saveMode      int
}

// NewRxnfileSaver creates a new RXN file saver
func NewRxnfileSaver(output *Output) *RxnfileSaver {
	return &RxnfileSaver{
		output:        output,
		addStereoDesc: false,
		saveMode:      0,
	}
}

// SaveReaction saves a reaction to RXN format
func (rxs *RxnfileSaver) SaveReaction(rxn *Reaction) error {
	rxs.rxn = rxn
	return rxs.saveReaction()
}

// saveReaction is the internal implementation for saving reactions
func (rxs *RxnfileSaver) saveReaction() error {
	// Write RXN header
	rxs.output.WriteString("$RXN\n")
	rxs.output.WriteString("\n")
	rxs.output.WriteString("\n")
	rxs.output.WriteString("\n")

	// Write counts
	nReactants := rxs.rxn.ReactantsCount()
	nProducts := rxs.rxn.ProductsCount()

	rxs.output.WriteString(fmt.Sprintf("%3d%3d\n", nReactants, nProducts))

	// Write reactants
	for i := rxs.rxn.ReactantBegin(); i != rxs.rxn.ReactantEnd(); i = rxs.rxn.ReactantNext(i) {
		rxs.output.WriteString("$MOL\n")
		rxs.writeMolecule(i)
	}

	// Write products
	for i := rxs.rxn.ProductBegin(); i != rxs.rxn.ProductEnd(); i = rxs.rxn.ProductNext(i) {
		rxs.output.WriteString("$MOL\n")
		rxs.writeMolecule(i)
	}

	return nil
}

// writeMolecule writes a molecule in MOL format
func (rxs *RxnfileSaver) writeMolecule(idx int) {
	mol := rxs.rxn.GetMolecule(idx)
	_ = mol // Placeholder
}
