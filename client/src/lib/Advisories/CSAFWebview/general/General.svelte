<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { Status } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
  import { getReadableDateString } from "../helpers";
  import Cvss from "./CVSS.svelte";

  let trackingVersion = $derived(appStore.state.webview.doc?.trackingVersion);
  let generator = $derived(appStore.state.webview.doc?.generator);
  let publisherName = $derived(appStore.state.webview.doc?.publisher.name);
  let publisherCategory = $derived(appStore.state.webview.doc?.publisher.category);
  let publisherNamespace = $derived(appStore.state.webview.doc?.publisher.namespace);
  let publisherIssuingAuthority = $derived(appStore.state.webview.doc?.publisher.issuing_authority);
  let publisherContactDetails = $derived(appStore.state.webview.doc?.publisher.contact_details);
  let category = $derived(appStore.state.webview.doc?.category);
  let title = $derived(appStore.state.webview.doc?.title);
  let lang = $derived(appStore.state.webview.doc?.lang);
  let sourceLang = $derived(appStore.state.webview.doc?.sourceLang);
  let csafVersion = $derived(appStore.state.webview.doc?.csafVersion);
  let distributionText = $derived(appStore.state.webview.doc?.distributionText);
  let published = $derived(appStore.state.webview.doc?.published);
  let lastUpdate = $derived(appStore.state.webview.doc?.lastUpdate);
  let status = $derived(appStore.state.webview.doc?.status);
  let baseSeverity = $derived(appStore.state.webview.doc?.highestScore?.baseSeverity);
  let baseScore = $derived(appStore.state.webview.doc?.highestScore?.baseScore);
  const cellStyleValue = "content-center px-6 py-0 [word-wrap:break-word] hyphens-auto";
  const cellStyleKey = "content-center w-40 py-0";
</script>

<div class="w-full">
  <div class="mb-3">
    <div class="flex flex-row">
      <div>
        <span class="text-xl text-balance">{title} </span>
        {#if appStore.state.webview.doc?.status !== Status.final}
          <span class="ml-3 text-lg text-gray-400">{status}</span>
        {/if}
      </div>
      {#if appStore.state.webview.doc?.highestScore}
        <Cvss {baseScore} baseSeverity={baseSeverity ?? ""}></Cvss>
      {/if}
    </div>
  </div>
  <div class="flex w-full flex-row flex-wrap">
    <div class="grid w-full grid-cols-[auto_minmax(0,_1fr)] gap-1.5 text-sm">
      <div class={cellStyleKey}>Publisher name</div>
      <div class={cellStyleValue}>{publisherName}</div>
      <div class={cellStyleKey}>Publisher namespace</div>
      <div class={cellStyleValue}>
        <a href={publisherNamespace} class="underline">
          <i class="bx bx-link"></i>{publisherNamespace}
        </a>
      </div>
      {#if publisherContactDetails}
        <div class={cellStyleKey}>Publisher contact details</div>
        <div class={cellStyleValue}>{publisherContactDetails}</div>
      {/if}
      {#if publisherIssuingAuthority}
        <div class={cellStyleKey}>Publisher issuing authority</div>
        <div class={cellStyleValue}>{publisherIssuingAuthority}</div>
      {/if}
      <div class={cellStyleKey}>Publisher category</div>
      <div class={cellStyleValue}>{publisherCategory}</div>
      <div class={cellStyleKey}>Published</div>
      <div class={cellStyleValue}>{getReadableDateString(published)}</div>
      <div class={cellStyleKey}>Last update</div>
      <div class={cellStyleValue}>{getReadableDateString(lastUpdate)}</div>
      <div class={cellStyleKey}>CSAF-Version</div>
      <div class={cellStyleValue}>{csafVersion}</div>
      <div class={cellStyleKey}>Category</div>
      <div class={cellStyleValue}>{category}</div>
      {#if distributionText}
        <div class={cellStyleKey}>Distribution</div>
        <div class={cellStyleValue}>{distributionText}</div>
      {/if}
      {#if appStore.state.webview.doc?.aggregateSeverity}
        <div class={cellStyleKey}>Aggregate severity text</div>
        <div class={cellStyleValue}>
          <span>{appStore.state.webview.doc?.aggregateSeverity.text}</span>
        </div>
        {#if appStore.state.webview.doc?.aggregateSeverity.namespace}
          <div class={cellStyleKey}>Aggregate severity namespace</div>
          <div class={cellStyleValue}>
            <a href={appStore.state.webview.doc?.aggregateSeverity.namespace} class="underline">
              <i class="bx bx-link"></i>
              <span>{appStore.state.webview.doc?.aggregateSeverity.namespace}</span>
            </a>
          </div>
        {/if}
      {/if}
      {#if lang}
        <div class={cellStyleKey}>Language</div>
        <div class={cellStyleValue}>{lang}</div>
      {/if}
      {#if sourceLang}
        <div class={cellStyleKey}>Source lang</div>
        <div class={cellStyleValue}>{sourceLang}</div>
      {/if}
      <div class={cellStyleKey}>Tracking Version</div>
      <div class={cellStyleValue}>{trackingVersion}</div>
    </div>
  </div>
  <div class="mt-3 flex flex-row">
    <span class="text-sm text-gray-400">
      Generator:
      {#if generator}
        {appStore.state.webview.doc?.generator?.engine.name}
      {/if}
      {#if generator?.engine?.version}
        {appStore.state.webview.doc?.generator?.engine.version}
      {/if}
      {#if generator?.date}
        · {getReadableDateString(generator.date)}
      {/if}
    </span>
  </div>
</div>
