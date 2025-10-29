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
	"sort"
	"strings"
)

// ReactionGrossFormula calculates the gross formula for a reaction
type ReactionGrossFormula struct{}

// Calculate computes the gross formula string for a reaction
func (rgf *ReactionGrossFormula) Calculate(rxn *BaseReaction) string {
	reactantFormula := rgf.calculateSide(rxn, Reactant)
	productFormula := rgf.calculateSide(rxn, Product)

	return fmt.Sprintf("%s >> %s", reactantFormula, productFormula)
}

// calculateSide calculates the formula for one side of the reaction
func (rgf *ReactionGrossFormula) calculateSide(rxn *BaseReaction, side int) string {
	elementCounts := make(map[string]int)

	begin := 0
	next := func(i int) int { return i + 1 }
	end := rxn.End()

	switch side {
	case Reactant:
		begin = rxn.ReactantBegin()
		next = rxn.ReactantNext
		end = rxn.ReactantEnd()
	case Product:
		begin = rxn.ProductBegin()
		next = rxn.ProductNext
		end = rxn.ProductEnd()
	case Catalyst:
		begin = rxn.CatalystBegin()
		next = rxn.CatalystNext
		end = rxn.CatalystEnd()
	}

	for i := begin; i != end; i = next(i) {
		mol := rxn.GetBaseMolecule(i)
		_ = mol // In real implementation, would analyze molecular formula
	}

	return rgf.formatFormula(elementCounts)
}

// formatFormula formats the element counts into a formula string
func (rgf *ReactionGrossFormula) formatFormula(elementCounts map[string]int) string {
	if len(elementCounts) == 0 {
		return ""
	}

	// Sort elements
	elements := make([]string, 0, len(elementCounts))
	for elem := range elementCounts {
		elements = append(elements, elem)
	}
	sort.Strings(elements)

	// Build formula string
	var parts []string
	for _, elem := range elements {
		count := elementCounts[elem]
		if count == 1 {
			parts = append(parts, elem)
		} else {
			parts = append(parts, fmt.Sprintf("%s%d", elem, count))
		}
	}

	return strings.Join(parts, " ")
}
