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

// ReactionCdxmlLoader loads reactions from CDXML format
type ReactionCdxmlLoader struct {
	scanner *Scanner
}

// NewReactionCdxmlLoader creates a new CDXML loader
func NewReactionCdxmlLoader(scanner *Scanner) *ReactionCdxmlLoader {
	return &ReactionCdxmlLoader{
		scanner: scanner,
	}
}

// LoadReaction loads a reaction from CDXML format
func (rcl *ReactionCdxmlLoader) LoadReaction(rxn *Reaction) error {
	// Implementation would parse CDXML format
	return nil
}

// LoadQueryReaction loads a query reaction from CDXML format
func (rcl *ReactionCdxmlLoader) LoadQueryReaction(qrxn *QueryReaction) error {
	// Implementation would parse CDXML format for query
	return nil
}
