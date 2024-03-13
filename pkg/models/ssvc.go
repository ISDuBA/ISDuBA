// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import (
	_ "embed" // Used for embedding.
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

//go:embed CISA-Coordinator.json
var cisaCoordinator []byte

type ssvcDecision struct {
	Exploitation     string `json:"Exploitation"`
	Automatable      string `json:"Automatable"`
	TechnicalImpact  string `json:"Technical Impact"`
	MissionWellBeing string `json:"Mission & Well-being"`
	Decision         string `json:"Decision"`
}

type ssvcDecisionPointOptionChildCombination struct {
	ChildLabel        string   `json:"child_label"`
	ChildKey          string   `json:"child_key"`
	ChildOptionLabels []string `json:"child_option_labels"`
	ChildOptionKeys   []string `json:"child_option_keys"`
}

type ssvcDecisionPointOption struct {
	Label             string                                      `json:"label"`
	Key               string                                      `json:"key"`
	Description       string                                      `json:"description"`
	ChildCombinations [][]ssvcDecisionPointOptionChildCombination `json:"child_combinations"`
}

type ssvcDecisionPointChild struct {
	Label string `json:"label"`
}

type ssvcDecisionPoint struct {
	Label        string                    `json:"label"`
	DecisionType string                    `json:"decision_type"`
	Key          string                    `json:"key"`
	Options      []ssvcDecisionPointOption `json:"options"`
	Children     []ssvcDecisionPointChild  `json:"children"`
}

type ssvc struct {
	DecisionPoints []ssvcDecisionPoint `json:"decision_points"`
	DecisionTable  []ssvcDecision      `json:"decisions_table"`
}

func (s *ssvc) validateVector(vector string) error {
	parts := strings.Split(vector, "/")
	if len(parts) < 4 {
		return errors.New("vector has invalid length")
	}

	// TODO: Implement me!

	return nil
}

var parsedSSVCv2 = sync.OnceValue(func() *ssvc {
	s := new(ssvc)
	if err := json.Unmarshal(cisaCoordinator, s); err != nil {
		panic(fmt.Sprintf("cannot parse 'CISA-coordinator.json': %v", err))
	}
	return s
})

// ValidateSSVCv2Vector checks if the given SSVCv2 vector is valid.
func ValidateSSVCv2Vector(vector string) error {
	return parsedSSVCv2().validateVector(vector)
}
