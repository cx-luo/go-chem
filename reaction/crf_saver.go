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

// CrfSaver saves reactions to CRF (Chemistry Resource File) format
type CrfSaver struct {
	output *Output
}

// NewCrfSaver creates a new CRF saver
func NewCrfSaver(output *Output) *CrfSaver {
	return &CrfSaver{
		output: output,
	}
}

// SaveReaction saves a reaction to CRF format
func (cs *CrfSaver) SaveReaction(rxn *Reaction) error {
	// Implementation would write CRF format
	return nil
}
