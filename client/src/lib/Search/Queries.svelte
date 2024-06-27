<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { request } from "$lib/utils";
  import { Button, ButtonGroup } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import { push } from "svelte-spa-router";

  let queries: any[] = [];
  $: sortedQueries = queries.sort((a: any, b: any) => {
    if (a.global && !b.global) {
      return -1;
    } else if (!a.global && b.global) {
      return 1;
    }
    return 0;
  });
  export let selectedIndex = -1;
  let errorMessage = "";
  let advancedQueryErrorMessage = "";
  let globalQueryButtonColor = "primary";
  let defaultQueryButtonClass = "flex flex-col p-0 focus:text-black hover:text-black";
  let queryButtonClass = "bg-white hover:bg-gray-100";
  let pressedQueryButtonClass = "bg-gray-200 text-black hover:!bg-gray-100";
  let globalQueryButtonClass = `border-${globalQueryButtonColor}-500 hover:!text-black`;
  let pressedGlobalQueryButtonClass = `border-${globalQueryButtonColor}-500 bg-${globalQueryButtonColor}-600 focus:text-white text-white hover:text-black`;

  const getClass = (isGlobal: boolean, isPressed: boolean) => {
    const addition = isGlobal
      ? isPressed
        ? pressedGlobalQueryButtonClass
        : globalQueryButtonClass
      : isPressed
        ? pressedQueryButtonClass
        : queryButtonClass;
    return `${defaultQueryButtonClass} ${addition}`;
  };

  onMount(async () => {
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      queries = response.content;
    } else if (response.error) {
      errorMessage = `Could not load user defined queries. ${getErrorMessage(response.error)}`;
    }
  });

  const selectQuery = (index: number) => {
    if (selectedIndex == index) {
      selectedIndex = -1;
    } else {
      selectedIndex = index;
    }
  };
</script>

<div class="mb-8 flex flex-col flex-wrap gap-4">
  <div class="flex gap-x-4">
    <ButtonGroup class="flex-wrap">
      {#each sortedQueries as query, index}
        <Button
          size="xs"
          on:click={() => selectQuery(index)}
          class={getClass(query.global, index === selectedIndex)}
        >
          <span title={query.description} class="p-2">{query.name}</span>
        </Button>
      {/each}
      <Button
        title="Configure queries"
        size="xs"
        on:click={() => {
          push("/queries");
        }}
      >
        <i class="bx bx-cog"></i>
      </Button>
    </ButtonGroup>
  </div>
  <ErrorMessage message={advancedQueryErrorMessage}></ErrorMessage>
</div>
<ErrorMessage message={errorMessage}></ErrorMessage>
