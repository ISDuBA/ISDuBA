//go:build ignore

// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/ISDuBA/ISDuBA/pkg/models"
)

const tmplTxt = `/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

// THIS FILE IS MACHINE GENERATED!
// Use "go generate ./..." in the root folder to regenerate it.

export type WorkflowState = string;
export const NEW: WorkflowState = "new";
export const READ: WorkflowState = "read";
export const ASSESSING: WorkflowState = "assessing";
export const REVIEW: WorkflowState = "review";
export const ARCHIVED: WorkflowState = "archived";
export const DELETE: WorkflowState = "delete";

export const WORKFLOW_STATES = [NEW, READ, ASSESSING, REVIEW, ARCHIVED, DELETE];

export type Role = string;
export const ADMIN: Role = "admin";
export const IMPORTER: Role = "importer";
export const EDITOR: Role = "editor";
export const REVIEWER: Role = "reviewer";
export const AUDITOR: Role = "auditor";
export const SOURCE_MANAGER: Role = "source-manager";

export type WorkflowStateTransition = {
  from: WorkflowState;
  to: WorkflowState;
  roles: Role[];
};

{{ $out := false -}}
export const WORKFLOW_TRANSITIONS: WorkflowStateTransition[] = [
  {{ range $j, $key := $.keys }}
  {{- $who := index $.workflow $key }}
  {{- $from := index $key 0 -}}
  {{- $to := index $key 1 -}}
  {{- if eq $to "" }}{{ continue }}{{ end -}}
  {{- if eq $from "" }}{{ continue }}{{ end -}}
  {{- if $out }},
  {{ end }}{{ $out = true -}}
  { from: {{ $from | string | upper }}, to: {{ $to | string | upper }}, roles: [{{ range $i, $role := $who }}
  {{- if $i }}, {{ end }}{{ $role | string | upper }}{{ end }}] }
  {{- end }}
];
`

var tmpl = template.Must(template.New("states").Funcs(
	template.FuncMap{
		"string": fmt.Sprint,
		"upper":  strings.ToUpper,
	}).Parse(tmplTxt))

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func keys[K comparable, V any](m map[K]V) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func main() {
	output := flag.String("o", "workflow.ts", "TS file to generate")
	flag.Parse()
	out, err := os.Create(*output)
	check(err)
	ks := keys(models.Transitions)
	slices.SortFunc(ks, func(a, b [2]models.Workflow) int {
		if d := cmp.Compare(a[0], b[0]); d != 0 {
			return d
		}
		return cmp.Compare(a[1], b[1])
	})
	err1 := tmpl.Execute(out, map[string]any{
		"workflow": models.Transitions,
		"keys":     ks,
	})
	err2 := out.Close()
	check(errors.Join(err1, err2))
}
