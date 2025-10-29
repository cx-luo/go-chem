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

// ReactionMultistepDetector detects multi-step reactions
type ReactionMultistepDetector struct {
	reaction *BaseReaction
}

// NewReactionMultistepDetector creates a new multistep detector
func NewReactionMultistepDetector(reaction *BaseReaction) *ReactionMultistepDetector {
	return &ReactionMultistepDetector{
		reaction: reaction,
	}
}

// Detect detects if the reaction is multi-step
func (rmd *ReactionMultistepDetector) Detect() bool {
	// Implementation would analyze the reaction to detect multiple steps
	return len(rmd.reaction.reactionBlocks) > 1
}

// SplitIntoSteps splits a multi-step reaction into individual steps
func (rmd *ReactionMultistepDetector) SplitIntoSteps() []*BaseReaction {
	steps := make([]*BaseReaction, 0)

	for i := 0; i < len(rmd.reaction.reactionBlocks); i++ {
		// Create a new reaction for each block
		step := NewBaseReaction()
		// Copy the reaction block data to the step
		// Implementation details...
		steps = append(steps, step)
	}

	return steps
}
