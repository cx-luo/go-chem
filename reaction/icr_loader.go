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

// IcrLoader loads reactions from ICR (Indigo Chemical Reaction) format
type IcrLoader struct {
	scanner *Scanner
}

// NewIcrLoader creates a new ICR loader
func NewIcrLoader(scanner *Scanner) *IcrLoader {
	return &IcrLoader{
		scanner: scanner,
	}
}

// LoadReaction loads a reaction from ICR format
func (il *IcrLoader) LoadReaction(rxn *Reaction) error {
	// Implementation would parse ICR binary format
	return nil
}

// LoadQueryReaction loads a query reaction from ICR format
func (il *IcrLoader) LoadQueryReaction(qrxn *QueryReaction) error {
	// Implementation would parse ICR binary format for query
	return nil
}
