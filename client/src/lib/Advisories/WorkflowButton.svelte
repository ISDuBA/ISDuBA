<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import CBadge from "$lib/Components/CBadge.svelte";
  import { getAllowedWorkflowChanges } from "$lib/permissions";
  import { getContext } from "svelte";
  import WorkflowStateIcon from "./WorkflowStateIcon.svelte";

  interface Props {
    state: string;
    tooltip: string;
    onClick?: (event: any) => void;
  }
  let { state, tooltip, onClick = undefined }: Props = $props();

  let currentState: () => string = getContext("currentState");
  let updateStateFn: () => (newState: string) => void = getContext("updateStateFn");

  let allowedTargetStates = $derived(
    getAllowedWorkflowChanges([currentState()]).map((transition) => transition.to)
  );

  let disabled = $derived.by(() => {
    return currentState() === undefined || !allowedTargetStates.includes(state);
  });

  let badgeColor: any = $derived.by(() => {
    if (state === currentState()) {
      return "green";
    } else if (!disabled) {
      return "dark";
    } else {
      return "none";
    }
  });

  const buttonClass =
    "h-fit w-fit rounded-xs border-0 p-0 hover:bg-transparent cursor-pointer disabled:cursor-default";
</script>

<button
  {disabled}
  class={buttonClass}
  onclick={(event: any) => {
    if (onClick) {
      onClick(event);
    } else if (updateStateFn) {
      const callback = updateStateFn();
      callback(state);
    }
  }}
>
  <CBadge showHoverEffect={false} title={tooltip} class="flex w-fit gap-1" color={badgeColor}>
    <WorkflowStateIcon advisoryState={state} />
    <span>{state}</span>
  </CBadge>
</button>
