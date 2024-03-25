<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store";
  import { Status, TLP } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
  import { Table, TableBodyCell, TableBodyRow } from "flowbite-svelte";
  import Acknowledgments from "$lib/Advisories/CSAFWebview/acknowledgments/Acknowledgments.svelte";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import Notes from "$lib/Advisories/CSAFWebview/notes/Notes.svelte";
  import References from "$lib/Advisories/CSAFWebview/references/References.svelte";
  import RevisionHistory from "./RevisionHistory.svelte";
  import ValueList from "$lib/Advisories/CSAFWebview/ValueList.svelte";
  let tlpStyle = "";
  $: aliases = $appStore.webview.doc?.aliases;
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
  $: tlp = $appStore.webview.doc?.tlp;
  $: tlpurl = $appStore.webview.doc?.tlp.url;
  $: if (tlp?.label === TLP.WHITE) {
    tlpStyle = "tlpclear";
  } else if (tlp?.label === TLP.RED) {
    tlpStyle = "tlred";
  } else if (tlp?.label === TLP.AMBER) {
    tlpStyle = "tlamber";
  } else if (tlp?.label === TLP.GREEN) {
    tlpStyle = "tlgreen";
  }
  $: id = $appStore.webview.doc?.id;
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
  const cellStyle = "px-6 py-1";
</script>

<div class="w-max">
  <Table noborder>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>ID</TableBodyCell>
      <TableBodyCell class={cellStyle}>{id}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>CSAF-Version</TableBodyCell>
      <TableBodyCell class={cellStyle}>{csafVersion}</TableBodyCell>
    </TableBodyRow>
    {#if $appStore.webview.doc?.aggregateSeverity}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Aggregate severity text</TableBodyCell>
        <TableBodyCell class={cellStyle}>
          <span>{$appStore.webview.doc?.aggregateSeverity.text}</span></TableBodyCell
        >
      </TableBodyRow>
      {#if $appStore.webview.doc?.aggregateSeverity.namespace}
        <TableBodyRow>
          <TableBodyCell class={cellStyle}>Aggregate severity namespace</TableBodyCell>
          <TableBodyCell class={cellStyle}
            ><span>{$appStore.webview.doc?.aggregateSeverity.namespace}</span></TableBodyCell
          >
        </TableBodyRow>
      {/if}
    {/if}
    {#if tlp?.label}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>TLP</TableBodyCell>
        <TableBodyCell class={cellStyle}><span class={tlpStyle}>{tlp?.label}</span></TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if tlp?.url}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>TLP URL</TableBodyCell>
        <TableBodyCell class={cellStyle}><a href={tlpurl}>{tlp?.url}</a></TableBodyCell>
      </TableBodyRow>
    {/if}
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Category</TableBodyCell>
      <TableBodyCell class={cellStyle}>{category}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Title</TableBodyCell>
      <TableBodyCell class={cellStyle}>{title}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Publisher name</TableBodyCell>
      <TableBodyCell class={cellStyle}>{publisherName}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Publisher category</TableBodyCell>
      <TableBodyCell class={cellStyle}>{publisherCategory}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Publisher namespace</TableBodyCell>
      <TableBodyCell class={cellStyle}>{publisherNamespace}</TableBodyCell>
    </TableBodyRow>
    {#if publisherIssuingAuthority}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Publisher issuing authority</TableBodyCell>
        <TableBodyCell class={cellStyle}>{publisherIssuingAuthority}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if publisherContactDetails}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Publisher contact details</TableBodyCell>
        <TableBodyCell class={cellStyle}>{publisherContactDetails}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if lang}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Language</TableBodyCell>
        <TableBodyCell class={cellStyle}>{lang}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if sourceLang}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Source lang</TableBodyCell>
        <TableBodyCell class={cellStyle}>{sourceLang}</TableBodyCell>
      </TableBodyRow>
    {/if}
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Published</TableBodyCell>
      <TableBodyCell class={cellStyle}>{published}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Last update</TableBodyCell>
      <TableBodyCell class={cellStyle}>{lastUpdate}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell class={cellStyle}>Tracking Version</TableBodyCell>
      <TableBodyCell class={cellStyle}>{trackingVersion}</TableBodyCell>
    </TableBodyRow>
    {#if $appStore.webview.doc?.status !== Status.final}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Status</TableBodyCell>
        <TableBodyCell class={cellStyle}>{status}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if generator}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Generator engine</TableBodyCell>
        <TableBodyCell class={cellStyle}
          ><span>{$appStore.webview.doc?.generator?.engine.name}</span></TableBodyCell
        >
      </TableBodyRow>
    {/if}
    {#if generator?.engine?.version}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Generator engine version</TableBodyCell>
        <TableBodyCell class={cellStyle}
          ><span>{$appStore.webview.doc?.generator?.engine.version}</span></TableBodyCell
        >
      </TableBodyRow>
    {/if}
    {#if generator?.date}
      <TableBodyRow>
        <TableBodyCell class={cellStyle}>Generator date</TableBodyCell>
        <TableBodyCell class={cellStyle}><span>{generator?.date}</span></TableBodyCell>
      </TableBodyRow>
    {/if}
  </Table>
</div>

{#if aliases}
  <ValueList label="Aliases" values={aliases} />
{/if}

{#if $appStore.webview.doc?.isRevisionHistoryPresent}
  <div>
    <Collapsible
      header="Revision history"
      level="3"
      open={$appStore.webview.ui.isRevisionHistoryVisible}
    >
      <RevisionHistory />
    </Collapsible>
  </div>
{/if}

{#if $appStore.webview.doc?.notes}
  <div>
    <Collapsible header="Notes" level="3">
      <Notes notes={$appStore.webview.doc?.notes} />
    </Collapsible>
  </div>
{/if}

{#if $appStore.webview.doc?.acknowledgements}
  <div>
    <Collapsible header="Acknowledgements" level="3">
      <Acknowledgments acknowledgements={$appStore.webview.doc?.acknowledgements} />
    </Collapsible>
  </div>
{/if}

{#if $appStore.webview.doc && $appStore.webview.doc.references.length > 0}
  <div>
    <Collapsible header="References" level="3">
      <References references={$appStore.webview.doc?.references} />
    </Collapsible>
  </div>
{/if}

<style>
  .tlpclear {
    background: #000;
    color: #fff;
    padding: 0.3rem;
  }
  .tlpred {
    background: #000;
    color: #ff2b2b;
    padding: 0.3rem;
  }
  .tlpamber {
    background: #000;
    color: #ffc000;
    padding: 0.3rem;
  }
  .tlpgreen {
    background: #000;
    color: #33ff00;
    padding: 0.3rem;
  }
</style>
