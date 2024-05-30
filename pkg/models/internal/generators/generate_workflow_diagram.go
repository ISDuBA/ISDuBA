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
	"flag"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"text/template"

	"github.com/ISDuBA/ISDuBA/pkg/models"
)

const tmplTxt = `
digraph workflow_transitions {
	fontname="Helvetica,Arial,sans-serif"
	node [fontname="Helvetica,Arial,sans-serif"]
	edge [fontname="Helvetica,Arial,sans-serif"]
	rankdir=TB;
	node [shape = doublecircle]; start end;
	node [shape = box];
	{{ range $j, $states := $.keys }}
	{{- $who := index $.workflow $states -}}
	{{- $from := index $states 0 -}}
	{{- $to   := index $states 1 -}}
	{{- if eq $from "" }}{{ $from = "start" }}{{ end -}}
	{{- if eq $to "" }}{{ $to = "end" }}{{ end -}}
	{{ $from }} -> {{ $to }} [label = "{{ range $i, $role := $who -}}
	{{- if $i }}, {{ end }}{{ $role -}} 
	{{ end -}}"];
	{{ end }}
}
`

var tmpl = template.Must(template.New("states").Parse(tmplTxt))

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
	output := flag.String("o", "workflow.svg", "SVG file to generate")
	flag.Parse()

	cmd := exec.Command("dot", "-Tsvg", "-o", *output)
	stdin, err := cmd.StdinPipe()
	check(err)

	ks := keys(models.Transitions)
	slices.SortFunc(ks, func(a, b [2]models.Workflow) int {
		if d := cmp.Compare(a[0], b[0]); d != 0 {
			return d
		}
		return cmp.Compare(a[1], b[1])
	})

	go func() {
		defer stdin.Close()
		check(tmpl.Execute(stdin, map[string]any{
			"keys":     ks,
			"workflow": models.Transitions,
		}))
	}()
	out, err := cmd.CombinedOutput()
	check(err)
	fmt.Printf("%s\n", out)
}
