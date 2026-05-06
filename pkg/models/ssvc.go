// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
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

// SSVCHistoryEntry represents a singular ssvc change event
type SSVCHistoryEntry struct {
	SSVC             *string   `json:"ssvc"`
	ChangeDate       time.Time `json:"changedate"`
	ChangeNumber     int64     `json:"change_number"`
	Actor            *string   `json:"actor,omitempty"`
	DocumentsID      int64     `json:"documents_id"`
	DocumentsVersion *string   `json:"documents_version"`
}

// SSVCChange bundles a SSVCHistoryEntry with the previous SSVC
type SSVCChange struct {
	SSVCHistoryEntry
	SSVCPrev *string `json:"ssvc_prev"`
}

// SSVCResponse represents a singular SSVC
type SSVCResponse struct {
	SSVC *string `json:"ssvc,omitempty"`
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
	// Parts:
	// 1: SSVCv2
	// 2: Exploitation
	// 3: Automatable
	// 4: Technical Impact
	// The order between Mission & Well Being and the optional
	// Decision Points is not set. Having the mandatory in the middle
	// makes little sense, but wouldn't necessarily be incorrect.
	// ToDo: Evaluate whether that should be allowed. No for now.
	// 5/6/7: Mission Prevalence (Optional, must appear alongside Public Well-Being Impact)
	// 5/6/7: Public Well-Being Impact (Optional, must appear alongside Mission Prevalence)
	// 5/7: Mission & Well Being (Optional in theory as long as the other optional decions are set. We require it.)
	// 6/8: Decision
	// 7/9: Date
	// 8/10: End after last /
	if len(parts) < 8 || len(parts) > 10 {
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
	uniqueKeys := map[string]struct{}{}

	for part, decision := range decisions {
		key, option, ok := strings.Cut(decision, ":")
		if !ok {
			return fmt.Errorf("decision %q has no ':'", decision)
		}
		if _, ok = uniqueKeys[key]; ok {
			return fmt.Errorf("decision about %q was defined multiple times", key)
		}
		uniqueKeys[key] = struct{}{}
		dp := s.findDecisionPointByKey(key)
		if dp == nil {
			return fmt.Errorf("no decision point with key %q found", key)
		}
		if !dp.checkValidOrder(part) {
			return fmt.Errorf("invalid order of decision points. %q at point %d", dp.Label, part)
		}
		if dp.findOption(option) == nil {
			return fmt.Errorf("decision point %q has no option %q", dp.Key, option)
		}
	}
	return nil
}

func (dp *ssvcDecisionPoint) checkValidOrder(part int) bool {
	switch dp.Label {
	case "Exploitation":
		return part == 0
	case "Automatable":
		return part == 1
	case "Technical Impact":
		return part == 2
	case "Mission & Well-being":
		return part == 3 || part == 5
	case "Mission Prevalence":
		return 2 < part && 6 > part
	case "Public Well-being Impact":
		return 2 < part && 6 > part
	case "Decision":
		return part == 4 || part == 6
	}
	return true
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
