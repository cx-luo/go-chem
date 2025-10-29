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

// ReactionBlock represents a block of reactants, products, and catalysts
type ReactionBlock struct {
	Reactants []int
	Products  []int
	Catalysts []int
}

// NewReactionBlock creates a new ReactionBlock
func NewReactionBlock() *ReactionBlock {
	return &ReactionBlock{
		Reactants: make([]int, 0),
		Products:  make([]int, 0),
		Catalysts: make([]int, 0),
	}
}

// Copy copies data from another reaction block
func (rb *ReactionBlock) Copy(other *ReactionBlock) {
	rb.Reactants = make([]int, len(other.Reactants))
	copy(rb.Reactants, other.Reactants)

	rb.Products = make([]int, len(other.Products))
	copy(rb.Products, other.Products)

	rb.Catalysts = make([]int, len(other.Catalysts))
	copy(rb.Catalysts, other.Catalysts)
}
