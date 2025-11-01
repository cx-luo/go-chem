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
	"io"
)

// Output represents an output writer
type Output struct {
	writer io.Writer
}

// NewOutput creates a new Output
func NewOutput(writer io.Writer) *Output {
	return &Output{writer: writer}
}

// Write writes bytes to the output
func (o *Output) Write(data []byte) (int, error) {
	return o.writer.Write(data)
}

// WriteString writes a string to the output
func (o *Output) WriteString(s string) error {
	_, err := o.writer.Write([]byte(s))
	return err
}

// RSmilesSaver saves reactions to SMILES format
type RSmilesSaver struct {
	output     *Output
	brxn       *BaseReaction
	qrxn       *QueryReaction
	rxn        *Reaction
	SmartsMode bool
	Chemaxon   bool
	comma      bool
}

// NewRSmilesSaver creates a new SMILES reaction saver
func NewRSmilesSaver(output *Output) *RSmilesSaver {
	return &RSmilesSaver{
		output:     output,
		SmartsMode: false,
		Chemaxon:   false,
		comma:      false,
	}
}

// SaveReaction saves a regular reaction to SMILES
func (rss *RSmilesSaver) SaveReaction(reaction *Reaction) error {
	rss.rxn = reaction
	rss.brxn = reaction.BaseReaction
	return rss.saveReaction()
}

// SaveQueryReaction saves a query reaction to SMILES
func (rss *RSmilesSaver) SaveQueryReaction(reaction *QueryReaction) error {
	rss.qrxn = reaction
	rss.brxn = reaction.BaseReaction
	return rss.saveReaction()
}

// saveReaction is the internal implementation for saving reactions
func (rss *RSmilesSaver) saveReaction() error {
	// Write reactants
	for i := rss.brxn.ReactantBegin(); i != rss.brxn.ReactantEnd(); i = rss.brxn.ReactantNext(i) {
		if i != rss.brxn.ReactantBegin() {
			rss.output.WriteString(".")
		}
		rss.writeMolecule(i)
	}

	rss.output.WriteString(">")

	// Write catalysts/agents
	for i := rss.brxn.CatalystBegin(); i != rss.brxn.CatalystEnd(); i = rss.brxn.CatalystNext(i) {
		if i != rss.brxn.CatalystBegin() {
			rss.output.WriteString(".")
		}
		rss.writeMolecule(i)
	}

	rss.output.WriteString(">")

	// Write products
	for i := rss.brxn.ProductBegin(); i != rss.brxn.ProductEnd(); i = rss.brxn.ProductNext(i) {
		if i != rss.brxn.ProductBegin() {
			rss.output.WriteString(".")
		}
		rss.writeMolecule(i)
	}

	return nil
}

// writeMolecule writes a single molecule to SMILES
func (rss *RSmilesSaver) writeMolecule(i int) {
	// Implementation would convert molecule to SMILES
	mol := rss.brxn.GetBaseMolecule(i)
	_ = mol // Placeholder
}
