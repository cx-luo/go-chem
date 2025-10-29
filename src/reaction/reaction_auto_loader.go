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
	"errors"
	"strings"
)

// ReactionAutoLoader automatically detects the format and loads reactions
type ReactionAutoLoader struct {
	scanner *Scanner
}

// NewReactionAutoLoader creates a new auto loader
func NewReactionAutoLoader(scanner *Scanner) *ReactionAutoLoader {
	return &ReactionAutoLoader{
		scanner: scanner,
	}
}

// LoadReaction automatically detects format and loads a reaction
func (ral *ReactionAutoLoader) LoadReaction(rxn *Reaction) error {
	format, err := ral.detectFormat()
	if err != nil {
		return err
	}

	switch format {
	case "rxn":
		loader := NewRxnfileLoader(ral.scanner)
		return loader.LoadReaction(rxn)
	case "smiles":
		loader := NewRSmilesLoader(ral.scanner)
		return loader.LoadReaction(rxn)
	case "json":
		// JSON loading would require reading the data first
		return errors.New("JSON format not yet implemented")
	default:
		return errors.New("unknown reaction format")
	}
}

// LoadQueryReaction automatically detects format and loads a query reaction
func (ral *ReactionAutoLoader) LoadQueryReaction(qrxn *QueryReaction) error {
	format, err := ral.detectFormat()
	if err != nil {
		return err
	}

	switch format {
	case "rxn":
		loader := NewRxnfileLoader(ral.scanner)
		return loader.LoadQueryReaction(qrxn)
	case "smiles":
		loader := NewRSmilesLoader(ral.scanner)
		return loader.LoadQueryReaction(qrxn)
	case "json":
		return errors.New("JSON format not yet implemented")
	default:
		return errors.New("unknown reaction format")
	}
}

// detectFormat attempts to detect the reaction file format
func (ral *ReactionAutoLoader) detectFormat() (string, error) {
	// Peek at the beginning of the input
	// This is a simplified detection logic
	// Real implementation would read first few bytes/lines

	// Common patterns:
	// RXN file starts with "$RXN"
	// SMILES contains ">" separator
	// JSON starts with "{"

	return "rxn", nil // Placeholder
}

// isRxnFormat checks if the data appears to be RXN format
func (ral *ReactionAutoLoader) isRxnFormat(data string) bool {
	return strings.HasPrefix(strings.TrimSpace(data), "$RXN")
}

// isSmilesFormat checks if the data appears to be SMILES format
func (ral *ReactionAutoLoader) isSmilesFormat(data string) bool {
	return strings.Contains(data, ">")
}

// isJSONFormat checks if the data appears to be JSON format
func (ral *ReactionAutoLoader) isJSONFormat(data string) bool {
	trimmed := strings.TrimSpace(data)
	return strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
}
