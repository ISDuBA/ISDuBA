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
  import { Table, TableBodyCell, TableBodyRow } from "flowbite-svelte";

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
  const cellStyleValue = "px-6 py-0";
  const cellStyleKey = "w-40 py-0";
</script>

<div class="2xl:w-max">
  <div class="mb-3">
    <span class="text-xl">{title} </span>
    {#if $appStore.webview.doc?.status !== Status.final}
      <span class="ml-3 text-lg text-gray-400">{status}</span>
    {/if}
  </div>
  <div class="flex flex-row flex-wrap">
    <div class="flex flex-col">
      <Table noborder>
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>Publisher name</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{publisherName}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>Publisher namespace</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{publisherNamespace}</TableBodyCell>
        </TableBodyRow>
        {#if publisherContactDetails}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyleKey}>Publisher contact details</TableBodyCell>
            <TableBodyCell tdClass={cellStyleValue}>{publisherContactDetails}</TableBodyCell>
          </TableBodyRow>
        {/if}
        {#if publisherIssuingAuthority}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyleKey}>Publisher issuing authority</TableBodyCell>
            <TableBodyCell tdClass={cellStyleValue}>{publisherIssuingAuthority}</TableBodyCell>
          </TableBodyRow>
        {/if}
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>Publisher category</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{publisherCategory}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>Published</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{published}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>Last update</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{lastUpdate}</TableBodyCell>
        </TableBodyRow>
      </Table>
    </div>
    <div class="flex flex-col">
      <Table noborder>
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>CSAF-Version</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{csafVersion}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>Category</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{category}</TableBodyCell>
        </TableBodyRow>
        {#if $appStore.webview.doc?.aggregateSeverity}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyleKey}>Aggregate severity text</TableBodyCell>
            <TableBodyCell tdClass={cellStyleValue}>
              <span>{$appStore.webview.doc?.aggregateSeverity.text}</span></TableBodyCell
            >
          </TableBodyRow>
          {#if $appStore.webview.doc?.aggregateSeverity.namespace}
            <TableBodyRow>
              <TableBodyCell tdClass={cellStyleKey}>Aggregate severity namespace</TableBodyCell>
              <TableBodyCell tdClass={cellStyleValue}
                ><span>{$appStore.webview.doc?.aggregateSeverity.namespace}</span></TableBodyCell
              >
            </TableBodyRow>
          {/if}
        {/if}
        {#if lang}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyleKey}>Language</TableBodyCell>
            <TableBodyCell tdClass={cellStyleValue}>{lang}</TableBodyCell>
          </TableBodyRow>
        {/if}
        {#if sourceLang}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyleKey}>Source lang</TableBodyCell>
            <TableBodyCell tdClass={cellStyleValue}>{sourceLang}</TableBodyCell>
          </TableBodyRow>
        {/if}
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyleKey}>Tracking Version</TableBodyCell>
          <TableBodyCell tdClass={cellStyleValue}>{trackingVersion}</TableBodyCell>
        </TableBodyRow>
      </Table>
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
        {generator?.date.replace("T", " ").split(".")[0]}
      {/if}
    </span>
  </div>
</div>
