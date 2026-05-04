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
	"bytes"
	"encoding/json"
	"fmt"

	chem "github.com/cx-luo/go-chem/molecule"
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
	if rxn == nil {
		return fmt.Errorf("reaction: cannot save nil reaction to KET")
	}

	data := make(map[string]interface{})
	nodes := make([]interface{}, 0, rxn.Count()+4)
	molSeq := 0

	reactants, err := rjs.addKETSide(data, &nodes, rxn, Reactant, 0, 0, &molSeq)
	if err != nil {
		return err
	}
	catalysts, err := rjs.addKETSide(data, &nodes, rxn, Catalyst, 4.5, 2, &molSeq)
	if err != nil {
		return err
	}
	products, err := rjs.addKETSide(data, &nodes, rxn, Product, 8, 0, &molSeq)
	if err != nil {
		return err
	}
	_, err = rjs.addKETSide(data, &nodes, rxn, Intermediate, 8, -2, &molSeq)
	if err != nil {
		return err
	}
	_, err = rjs.addKETSide(data, &nodes, rxn, Undefined, 0, -2, &molSeq)
	if err != nil {
		return err
	}

	if reactants > 0 || products > 0 || catalysts > 0 || rxn.Count() > 1 {
		nodes = append(nodes, rjs.ketArrowNode(rxn))
	}
	rjs.addKETPlusNodes(&nodes, reactants, 0, 0)
	rjs.addKETPlusNodes(&nodes, products, 8, 0)

	data["root"] = map[string]interface{}{
		"nodes": nodes,
	}

	// Serialize to JSON
	var jsonData []byte
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

// SaveReactionToKETString saves a reaction to a Ketcher KET JSON string.
func SaveReactionToKETString(rxn *Reaction, pretty bool) (string, error) {
	if rxn == nil {
		return "", fmt.Errorf("reaction: cannot save nil reaction to KET")
	}
	var buf bytes.Buffer
	saver := NewReactionJsonSaver(NewOutput(&buf))
	saver.PrettyJSON = pretty
	if err := saver.SaveReaction(rxn.BaseReaction); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (rjs *ReactionJsonSaver) addKETSide(data map[string]interface{}, nodes *[]interface{}, rxn *BaseReaction, side int, x, y float64, molSeq *int) (int, error) {
	count := 0
	for i := rxn.SideBegin(side); i != rxn.SideEnd(); i = rxn.SideNext(side, i) {
		mol := rxn.GetBaseMolecule(i)
		if mol == nil {
			continue
		}

		ref := fmt.Sprintf("mol%d", *molSeq)
		*nodes = append(*nodes, map[string]interface{}{"$ref": ref})

		structure := mol.Structure
		if structure == nil {
			structure = chem.NewMolecule()
		}
		molObj, err := chem.BuildKETMoleculeObjectWithOffset(structure, chem.Vec2f{
			X: x + float64(count)*3,
			Y: y,
		}, chem.KETSaverOptions{BondLength: 1.5})
		if err != nil {
			return count, err
		}
		data[ref] = molObj
		(*molSeq)++
		count++
	}
	return count, nil
}

func (rjs *ReactionJsonSaver) ketArrowNode(rxn *BaseReaction) map[string]interface{} {
	mode := rjs.arrowTypeToString[ArrowBasic]
	if rxn.IsRetrosyntetic() {
		mode = rjs.arrowTypeToString[ArrowRetrosynthetic]
	}
	if mode == "" {
		mode = "open-angle"
	}

	return map[string]interface{}{
		"type": "arrow",
		"data": map[string]interface{}{
			"mode": mode,
			"pos": []map[string]float64{
				{"x": 5, "y": 0, "z": 0},
				{"x": 7, "y": 0, "z": 0},
			},
		},
	}
}

func (rjs *ReactionJsonSaver) addKETPlusNodes(nodes *[]interface{}, count int, x, y float64) {
	if count <= 1 {
		return
	}
	for i := 1; i < count; i++ {
		*nodes = append(*nodes, map[string]interface{}{
			"type":     "plus",
			"location": []float64{x + float64(i)*3 - 0.75, y, 0},
		})
	}
}
