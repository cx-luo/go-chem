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

// IcrSaver saves reactions to ICR (Indigo Chemical Reaction) format
type IcrSaver struct {
	output *Output
}

// NewIcrSaver creates a new ICR saver
func NewIcrSaver(output *Output) *IcrSaver {
	return &IcrSaver{
		output: output,
	}
}

// SaveReaction saves a reaction to ICR format
func (is *IcrSaver) SaveReaction(rxn *Reaction) error {
	// Implementation would write ICR binary format
	return nil
}

// SaveQueryReaction saves a query reaction to ICR format
func (is *IcrSaver) SaveQueryReaction(qrxn *QueryReaction) error {
	// Implementation would write ICR binary format for query
	return nil
}
