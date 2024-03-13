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
	"time"
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

func (dp *ssvcDecisionPoint) findOption(option string) *ssvcDecisionPointOption {
	for i := range dp.Options {
		opt := &dp.Options[i]
		if opt.Key == option {
			return opt
		}
	}
	return nil
}

func (s ssvc) findDecisionPointByKey(key string) *ssvcDecisionPoint {
	for i := range s.DecisionPoints {
		dp := &s.DecisionPoints[i]
		if dp.Key == key {
			return dp
		}
	}
	return nil
}

func (s *ssvc) validateVector(vector string) error {
	parts := strings.Split(vector, "/")
	if len(parts) < 4 {
		return errors.New("vector has invalid length")
	}
	if parts[0] != "SSVCv2" {
		return errors.New("vector does not start with 'SSVCv2'")
	}
	if parts[len(parts)-1] != "" {
		return errors.New("vector is not terminated with '/'")
	}

	const timestampFormat = "2006-01-02T15:04:05Z"
	ts := parts[len(parts)-2]
	if _, err := time.Parse(timestampFormat, ts); err != nil {
		return fmt.Errorf("vector timestamp is invalid: %v", err)
	}
	decisions := parts[1 : len(parts)-2]
	for _, decision := range decisions {
		key, option, ok := strings.Cut(decision, ":")
		if !ok {
			return fmt.Errorf("decision %q has no ':'", decision)
		}
		dp := s.findDecisionPointByKey(key)
		if dp == nil {
			return fmt.Errorf("no decision point with key %q found", key)
		}
		if dp.findOption(option) == nil {
			return fmt.Errorf("decision point %q has no option %q", dp.Key, option)
		}
	}
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
