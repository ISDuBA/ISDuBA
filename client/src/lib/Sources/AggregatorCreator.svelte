<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { saveAggregator, type Aggregator } from "$lib/Sources/source";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Input, Label, Button } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";
  import validator from "validator";

  let errorMessage: ErrorDetails | null;

  let validUrl: boolean | null = false;
  let urlColor: "red" | "green" | "base" = "base";
  $: if (validUrl !== undefined) {
    if (validUrl === null) {
      urlColor = "base";
    } else if (validUrl) {
      urlColor = "green";
    } else {
      urlColor = "red";
    }
  }

  let aggregator: Aggregator = {
    name: "",
    url: ""
  };

  let formClass = "max-w-[800pt]";

  const checkUrl = async () => {
    if (aggregator.url.startsWith("https://") && aggregator.url.endsWith("aggregator.json")) {
      validUrl = null;
      return;
    }
    if (validator.isFQDN(aggregator.url)) {
      validUrl = null;
      return;
    }
    validUrl = false;
  };

  const saveAll = async () => {
    let result = await saveAggregator(aggregator);
    if (!result.ok) {
      errorMessage = result.error;
    } else {
      push(`/sources`);
    }
  };
</script>

<svelte:head>
  <title>Sources - Add aggregator</title>
</svelte:head>

<div>
  <SectionHeader title="Add new aggregator"></SectionHeader>

  <form on:submit={saveAll} class={formClass}>
    <Label>Name</Label>
    <Input bind:value={aggregator.name}></Input>
    <Label>URL</Label>
    <Input bind:value={aggregator.url} on:input={checkUrl} color={urlColor}></Input>
    <br />
    <Button type="submit" color="light" disabled={validUrl === false || aggregator.name === ""}>
      <i class="bx bx-check me-2"></i>
      <span>Save aggregator</span>
    </Button>
  </form>

  <br />
  <ErrorMessage error={errorMessage}></ErrorMessage>
</div>
