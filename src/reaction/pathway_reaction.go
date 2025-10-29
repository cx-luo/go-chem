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

// PathwayReaction represents a multi-step reaction pathway
type PathwayReaction struct {
	*BaseReaction
	steps []*BaseReaction
}

// NewPathwayReaction creates a new pathway reaction
func NewPathwayReaction() *PathwayReaction {
	return &PathwayReaction{
		BaseReaction: NewBaseReaction(),
		steps:        make([]*BaseReaction, 0),
	}
}

// AddStep adds a reaction step to the pathway
func (pr *PathwayReaction) AddStep(reaction *BaseReaction) {
	pr.steps = append(pr.steps, reaction)
}

// GetStep returns a reaction step by index
func (pr *PathwayReaction) GetStep(index int) *BaseReaction {
	if index < 0 || index >= len(pr.steps) {
		return nil
	}
	return pr.steps[index]
}

// StepsCount returns the number of steps in the pathway
func (pr *PathwayReaction) StepsCount() int {
	return len(pr.steps)
}

// AsPathwayReaction returns this as a PathwayReaction
func (pr *PathwayReaction) AsPathwayReaction() *PathwayReaction {
	return pr
}

// IsPathwayReaction returns true for PathwayReaction
func (pr *PathwayReaction) IsPathwayReaction() bool {
	return true
}
