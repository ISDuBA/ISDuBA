// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import (
	"errors"
	"fmt"
	"strings"
)

// Workflow is a state of an advisory.
type Workflow string

// Different Workflows
const (
	NewWorkflow       Workflow = "new"       // NewWorkflow represents 'new'.
	ReadWorkflow      Workflow = "read"      // ReadWorkflow represents 'read'.
	AssessingWorkflow Workflow = "assessing" // AssessingWorkflow represents 'assessing',
	ReviewWorkflow    Workflow = "review"    // ReviewWorkflow represents 'review'.
	ArchivedWorkflow  Workflow = "archived"  // ArchivedWorkflow represents 'archived'.
	DeleteWorkflow    Workflow = "delete"    // DeleteWorkflow represents 'delete'.
)

// WorkflowRole is a role in the workflow.
type WorkflowRole string

// The different roles
const (
	Admin         WorkflowRole = "admin"          // Admin role
	Importer      WorkflowRole = "importer"       // Importer role
	Editor        WorkflowRole = "editor"         // Editor role
	Reviewer      WorkflowRole = "reviewer"       // Reviewer role
	Auditor       WorkflowRole = "auditor"        // Auditor role
	SourceManager WorkflowRole = "source-manager" // Source Manager role
)

// Transitions is a matrix to tell who is allowed to change between certain states.
// Please call "go generate ./..." in the root dir to update docs/workflow.svg
// if you change this.
var Transitions = map[[2]Workflow][]WorkflowRole{
	{"", NewWorkflow}:                     {Importer}, // Forward
	{NewWorkflow, ReadWorkflow}:           {Editor},
	{ReadWorkflow, AssessingWorkflow}:     {Editor},
	{AssessingWorkflow, ReviewWorkflow}:   {Editor},
	{ReviewWorkflow, ArchivedWorkflow}:    {Reviewer},
	{ReadWorkflow, DeleteWorkflow}:        {Editor, Reviewer},
	{AssessingWorkflow, DeleteWorkflow}:   {Editor, Reviewer},
	{ReviewWorkflow, DeleteWorkflow}:      {Reviewer},
	{ArchivedWorkflow, DeleteWorkflow}:    {Editor, Reviewer},
	{DeleteWorkflow, ArchivedWorkflow}:    {Admin}, // Backward
	{DeleteWorkflow, ReviewWorkflow}:      {Admin},
	{DeleteWorkflow, AssessingWorkflow}:   {Admin},
	{DeleteWorkflow, ReadWorkflow}:        {Admin},
	{ArchivedWorkflow, ReviewWorkflow}:    {Admin, Editor, Reviewer},
	{ArchivedWorkflow, AssessingWorkflow}: {Admin, Editor},
	{ArchivedWorkflow, ReadWorkflow}:      {Admin},
	{ArchivedWorkflow, NewWorkflow}:       {Importer},
	{ReviewWorkflow, AssessingWorkflow}:   {Reviewer, Editor},
	{ReviewWorkflow, ReadWorkflow}:        {Reviewer},
	{ReviewWorkflow, NewWorkflow}:         {Importer},
	{AssessingWorkflow, ReadWorkflow}:     {Editor},
	{AssessingWorkflow, NewWorkflow}:      {Importer},
	{ReadWorkflow, NewWorkflow}:           {Editor},
	{DeleteWorkflow, ""}:                  {Admin},
}

// ParseWorkflowRole parses a workflow role from a string.
func ParseWorkflowRole(s string) (WorkflowRole, error) {
	switch r := WorkflowRole(strings.ToLower(s)); r {
	case Admin, Importer, Editor, Reviewer, Auditor, SourceManager:
		return r, nil
	default:
		return "", fmt.Errorf("unknown workflow role %q", s)
	}
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (wfr *WorkflowRole) UnmarshalText(text []byte) error {
	x, err := ParseWorkflowRole(string(text))
	if err != nil {
		return err
	}
	*wfr = x
	return nil
}

// Scan implements [sql.Scanner].
func (wfr *WorkflowRole) Scan(src any) error {
	if s, ok := src.(string); ok {
		x, err := ParseWorkflowRole(s)
		if err != nil {
			return err
		}
		*wfr = x
		return nil
	}
	return errors.New("unsupported type")
}

// MarshalText implements [encoding.TextMarshaler].
func (wfr WorkflowRole) MarshalText() ([]byte, error) {
	return []byte(wfr), nil
}

// Valid returns true is the workflow represents a valid state.
func (wf Workflow) Valid() bool {
	switch wf {
	case NewWorkflow, ReadWorkflow, AssessingWorkflow, ReviewWorkflow, ArchivedWorkflow, DeleteWorkflow:
		return true
	default:
		return false
	}
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (wf *Workflow) UnmarshalText(text []byte) error {
	x := Workflow(text)
	if !x.Valid() {
		return fmt.Errorf("%q is no a valid workflow state", text)
	}
	*wf = x
	return nil
}

// TransitionsRoles return a list of roles that are allowed to do the requested
// transition.
func (wf Workflow) TransitionsRoles(other Workflow) []WorkflowRole {
	return Transitions[[2]Workflow{wf, other}]
}
