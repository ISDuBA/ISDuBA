<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store";
  import { Status } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
  import { getReadableDateString } from "../helpers";

  $: trackingVersion = $appStore.webview.doc?.trackingVersion;
  $: generator = $appStore.webview.doc?.generator;
  $: publisherName = $appStore.webview.doc?.publisher.name;
  $: publisherCategory = $appStore.webview.doc?.publisher.category;
  $: publisherNamespace = $appStore.webview.doc?.publisher.namespace;
  $: publisherIssuingAuthority = $appStore.webview.doc?.publisher.issuing_authority;
  $: publisherContactDetails = $appStore.webview.doc?.publisher.contact_details;
  $: category = $appStore.webview.doc?.category;
  $: title = $appStore.webview.doc?.title;
  $: lang = $appStore.webview.doc?.lang;
  $: sourceLang = $appStore.webview.doc?.sourceLang;
  $: csafVersion = $appStore.webview.doc?.csafVersion;
  $: published = $appStore.webview.doc?.published;
  $: lastUpdate = $appStore.webview.doc?.lastUpdate;
  $: status = $appStore.webview.doc?.status;
  $: if (
    !$appStore.webview.doc?.isRevisionHistoryPresent &&
    !$appStore.webview.doc?.isDocPresent &&
    !$appStore.webview.doc?.isProductTreePresent &&
    !$appStore.webview.doc?.isPublisherPresent &&
    !$appStore.webview.doc?.isTLPPresent &&
    !$appStore.webview.doc?.isTrackingPresent &&
    !$appStore.webview.doc?.isVulnerabilitiesPresent
  ) {
    appStore.setSingleErrorMsg("Are you sure the URL refers to a CSAF document?");
  }
  const cellStyleValue = "content-center px-6 py-0 [word-wrap:break-word] hyphens-auto";
  const cellStyleKey = "content-center w-40 py-0";
</script>

<div class="w-full">
  <div class="mb-3">
    <span class="text-xl">{title} </span>
    {#if $appStore.webview.doc?.status !== Status.final}
      <span class="ml-3 text-lg text-gray-400">{status}</span>
    {/if}
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
      {#if $appStore.webview.doc?.aggregateSeverity}
        <div class={cellStyleKey}>Aggregate severity text</div>
        <div class={cellStyleValue}>
          <span>{$appStore.webview.doc?.aggregateSeverity.text}</span>
        </div>
        {#if $appStore.webview.doc?.aggregateSeverity.namespace}
          <div class={cellStyleKey}>Aggregate severity namespace</div>
          <div class={cellStyleValue}>
            <a href={$appStore.webview.doc?.aggregateSeverity.namespace} class="underline">
              <i class="bx bx-link"></i>
              <span>{$appStore.webview.doc?.aggregateSeverity.namespace}</span>
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
        {$appStore.webview.doc?.generator?.engine.name}
      {/if}
      {#if generator?.engine?.version}
        {$appStore.webview.doc?.generator?.engine.version}
      {/if}
      {#if generator?.date}
        · {getReadableDateString(generator.date)}
      {/if}
    </span>
  </div>
</div>
