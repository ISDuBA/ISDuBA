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
	"flag"
	"fmt"
	"log"
	"os/exec"
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
	{{ range $states, $who := . }}
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

func main() {
	output := flag.String("o", "workflow.svg", "SVG file to generate")
	flag.Parse()

	cmd := exec.Command("dot", "-Tsvg", "-o", *output)
	stdin, err := cmd.StdinPipe()
	check(err)

	go func() {
		defer stdin.Close()
		check(tmpl.Execute(stdin, models.Transitions))
	}()
	out, err := cmd.CombinedOutput()
	check(err)
	fmt.Printf("%s\n", out)
}
