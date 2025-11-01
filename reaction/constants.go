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

// Stereo changes during reaction
const (
	StereoUnmarked = 0
	StereoInverts  = 1
	StereoRetains  = 2
)

// Reacting centers
const (
	RCNotCenter    = -1
	RCUnmarked     = 0
	RCCenter       = 1
	RCUnchanged    = 2
	RCMadeOrBroken = 4
	RCOrderChanged = 8
	RCTotal        = 16
)

// Reaction side types
const (
	Reactant     = 1
	Product      = 2
	Intermediate = 4
	Undefined    = 8
	Catalyst     = 16
)
