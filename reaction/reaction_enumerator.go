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

// ReactionProductEnumerator enumerates reaction products
type ReactionProductEnumerator struct {
	reaction *QueryReaction
	state    *ReactionEnumeratorState
}

// ReactionEnumeratorState maintains state for reaction enumeration
type ReactionEnumeratorState struct {
	currentIndex int
	totalCount   int
}

// NewReactionProductEnumerator creates a new product enumerator
func NewReactionProductEnumerator(reaction *QueryReaction) *ReactionProductEnumerator {
	return &ReactionProductEnumerator{
		reaction: reaction,
		state: &ReactionEnumeratorState{
			currentIndex: 0,
			totalCount:   0,
		},
	}
}

// Next generates the next product
func (rpe *ReactionProductEnumerator) Next() (*Molecule, bool) {
	if rpe.state.currentIndex >= rpe.state.totalCount {
		return nil, false
	}

	// Generate next product
	product := NewMolecule()
	rpe.state.currentIndex++

	return product, true
}

// HasNext returns true if there are more products to enumerate
func (rpe *ReactionProductEnumerator) HasNext() bool {
	return rpe.state.currentIndex < rpe.state.totalCount
}

// Reset resets the enumerator to the beginning
func (rpe *ReactionProductEnumerator) Reset() {
	rpe.state.currentIndex = 0
}
