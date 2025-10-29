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

// MoleculeFingerprintParameters represents parameters for fingerprint calculation
type MoleculeFingerprintParameters struct {
	Type       string
	Size       int
	Similarity bool
}

// ReactionFingerprintBuilder builds fingerprints for reactions
type ReactionFingerprintBuilder struct {
	reaction    *BaseReaction
	parameters  MoleculeFingerprintParameters
	Query       bool
	SkipOrd     bool
	SkipSim     bool
	SkipExt     bool
	fingerprint []byte
}

// NewReactionFingerprintBuilder creates a new fingerprint builder
func NewReactionFingerprintBuilder(reaction *BaseReaction, parameters MoleculeFingerprintParameters) *ReactionFingerprintBuilder {
	return &ReactionFingerprintBuilder{
		reaction:    reaction,
		parameters:  parameters,
		Query:       false,
		SkipOrd:     false,
		SkipSim:     false,
		SkipExt:     false,
		fingerprint: make([]byte, 0),
	}
}

// Process processes the reaction to build the fingerprint
func (rfb *ReactionFingerprintBuilder) Process() {
	// Implementation for fingerprint building
	// This would iterate through molecules and build composite fingerprint
}

// Get returns the fingerprint
func (rfb *ReactionFingerprintBuilder) Get() []byte {
	return rfb.fingerprint
}

// GetSim returns the similarity fingerprint
func (rfb *ReactionFingerprintBuilder) GetSim() []byte {
	return rfb.fingerprint
}

// ParseFingerprintType parses the fingerprint type string
func (rfb *ReactionFingerprintBuilder) ParseFingerprintType(fpType string, query bool) {
	rfb.Query = query
	// Parse type and set appropriate flags
	switch fpType {
	case "sim":
		rfb.SkipOrd = true
		rfb.SkipExt = true
	case "ord":
		rfb.SkipSim = true
		rfb.SkipExt = true
	case "ext":
		rfb.SkipOrd = true
		rfb.SkipSim = true
	}
}
