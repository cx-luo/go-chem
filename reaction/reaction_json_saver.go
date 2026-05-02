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
	"encoding/json"
)

// KETVersion represents the KET format version
type KETVersion int

const (
	KETVersionAuto KETVersion = 0
	KETVersion1    KETVersion = 1
	KETVersion2    KETVersion = 2
)

// ReactionJsonSaver saves reactions to JSON format
type ReactionJsonSaver struct {
	output             *Output
	AddStereoDesc      bool
	AddReactionData    bool
	PrettyJSON         bool
	KetVersion         KETVersion
	UseNativePrecision bool
	LayoutOptions      LayoutOptions
	arrowTypeToString  map[int]string
}

// NewReactionJsonSaver creates a new JSON reaction saver
func NewReactionJsonSaver(output *Output) *ReactionJsonSaver {
	return &ReactionJsonSaver{
		output:             output,
		AddStereoDesc:      false,
		AddReactionData:    true,
		PrettyJSON:         false,
		KetVersion:         KETVersionAuto,
		UseNativePrecision: false,
		LayoutOptions:      LayoutOptions{},
		arrowTypeToString: map[int]string{
			ArrowBasic:                                   "open-angle",
			ArrowFilledTriangle:                          "filled-triangle",
			ArrowFilledBow:                               "filled-bow",
			ArrowDashed:                                  "dashed-open-angle",
			ArrowFailed:                                  "failed",
			ArrowBothEndsFilledTriangle:                  "both-ends-filled-triangle",
			ArrowEquilibriumFilledHalfBow:                "equilibrium-filled-half-bow",
			ArrowEquilibriumFilledTriangle:               "equilibrium-filled-triangle",
			ArrowEquilibriumOpenAngle:                    "equilibrium-open-angle",
			ArrowUnbalancedEquilibriumFilledHalfBow:      "unbalanced-equilibrium-filled-half-bow",
			ArrowUnbalancedEquilibriumLargeFilledHalfBow: "unbalanced-equilibrium-large-filled-half-bow",
			ArrowRetrosynthetic:                          "retrosynthetic",
		},
	}
}

// SaveReaction saves a reaction to JSON format
func (rjs *ReactionJsonSaver) SaveReaction(rxn *BaseReaction) error {
	data := make(map[string]interface{})

	// Build JSON structure
	data["root"] = map[string]interface{}{
		"nodes": []interface{}{},
	}

	// Add reactants
	reactants := make([]interface{}, 0)
	for i := rxn.ReactantBegin(); i != rxn.ReactantEnd(); i = rxn.ReactantNext(i) {
		mol := rxn.GetBaseMolecule(i)
		_ = mol // Placeholder
		reactants = append(reactants, map[string]interface{}{})
	}

	// Add products
	products := make([]interface{}, 0)
	for i := rxn.ProductBegin(); i != rxn.ProductEnd(); i = rxn.ProductNext(i) {
		mol := rxn.GetBaseMolecule(i)
		_ = mol // Placeholder
		products = append(products, map[string]interface{}{})
	}

	// Serialize to JSON
	var jsonData []byte
	var err error
	if rjs.PrettyJSON {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	_, err = rjs.output.Write(jsonData)
	return err
}
