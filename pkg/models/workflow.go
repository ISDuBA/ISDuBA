// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import "fmt"

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

// The different roles
const (
	Admin    = "admin"    // Admin role
	Importer = "importer" // Importer role
	Editor   = "editor"   // Editor role
	Reviewer = "reviewer" // Reviewer role
	Auditor  = "auditor"  // Auditor role
)

// Transitions is a matrix to tell who is allowed to change between certain states.
// Please call "go generate ./..." in the root dir to update docs/workflow.svg
// if you change this.
var Transitions = map[[2]Workflow][]string{
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
	{ArchivedWorkflow, ReviewWorkflow}:    {Admin},
	{ArchivedWorkflow, AssessingWorkflow}: {Admin},
	{ArchivedWorkflow, ReadWorkflow}:      {Admin},
	{ArchivedWorkflow, NewWorkflow}:       {Importer},
	{ReviewWorkflow, AssessingWorkflow}:   {Reviewer},
	{ReviewWorkflow, ReadWorkflow}:        {Reviewer},
	{ReviewWorkflow, NewWorkflow}:         {Importer},
	{AssessingWorkflow, ReadWorkflow}:     {Editor},
	{AssessingWorkflow, NewWorkflow}:      {Importer},
	{ReadWorkflow, NewWorkflow}:           {Editor},
	{DeleteWorkflow, ""}:                  {Admin},
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
func (wf Workflow) TransitionsRoles(other Workflow) []string {
	return Transitions[[2]Workflow{wf, other}]
}

// CommentingAllowed returns true if commenting is allowed.
func (wf Workflow) CommentingAllowed() bool {
	return wf == ReadWorkflow || wf == AssessingWorkflow
}
