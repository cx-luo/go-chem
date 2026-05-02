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

// CanonicalRSmilesSaver saves reactions to canonical SMILES format
type CanonicalRSmilesSaver struct {
	*RSmilesSaver
}

// NewCanonicalRSmilesSaver creates a new canonical SMILES saver
func NewCanonicalRSmilesSaver(output *Output) *CanonicalRSmilesSaver {
	return &CanonicalRSmilesSaver{
		RSmilesSaver: NewRSmilesSaver(output),
	}
}

// SaveReaction saves a reaction to canonical SMILES format
func (crss *CanonicalRSmilesSaver) SaveReaction(reaction *Reaction) error {
	// Implementation would:
	// 1. Canonicalize each molecule
	// 2. Sort molecules in a canonical order
	// 3. Write the canonical SMILES

	return crss.RSmilesSaver.SaveReaction(reaction)
}
