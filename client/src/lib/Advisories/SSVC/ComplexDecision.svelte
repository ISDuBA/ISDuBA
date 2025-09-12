<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Radio } from "flowbite-svelte";
  import { onMount } from "svelte";
  import type { SSVCDecision, SSVCDecisionChild } from "./SSVCCalculator";
  import { getDecision } from "./SSVCCalculator";

  interface Props {
    children: SSVCDecisionChild[] | undefined;
    decisionPoints: SSVCDecision[];
  }

  let { children, decisionPoints }: Props = $props();

  onMount(() => {
    children?.forEach((child) => {
      const firstRadioButton = document.getElementsByName(child.label)?.[0] as HTMLInputElement;
      firstRadioButton.checked = true;
    });
  });
</script>

<div class="complex-decision flex flex-row gap-x-5">
  {#if children && decisionPoints}
    {#each children as child}
      {@const childOptions = getDecision(decisionPoints, child.label)?.options}
      {#if childOptions}
        <div class="flex flex-col">
          <span
            class="text-gary-400 mb-2 text-xs font-bold tracking-tight text-gray-900 dark:text-white"
          >
            {child.label}
          </span>
          <div class="flex flex-row gap-x-3">
            {#each childOptions as option}
              <div title={option.description} class="mb-2 cursor-pointer">
                <Radio
                  name={child.label}
                  value={option.label}
                  class="flex flex-col text-xs tracking-tight"
                  ><span class="mt-2">{option.label}</span></Radio
                >
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/each}
  {/if}
</div>
