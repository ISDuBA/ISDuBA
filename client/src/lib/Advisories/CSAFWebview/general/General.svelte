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
  import Acknowledgements from "$lib/Advisories/CSAFWebview/acknowledgements/Acknowledgements.svelte";
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
    tlpStyle = "tlpred";
  } else if (tlp?.label === TLP.AMBER) {
    tlpStyle = "tlpamber";
  } else if (tlp?.label === TLP.GREEN) {
    tlpStyle = "tlpgreen";
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
  const cellStyle = "px-6 py-0";
</script>

<div class="2xl:w-max">
  <Table noborder>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>ID</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{id}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>CSAF-Version</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{csafVersion}</TableBodyCell>
    </TableBodyRow>
    {#if $appStore.webview.doc?.aggregateSeverity}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Aggregate severity text</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}>
          <span>{$appStore.webview.doc?.aggregateSeverity.text}</span></TableBodyCell
        >
      </TableBodyRow>
      {#if $appStore.webview.doc?.aggregateSeverity.namespace}
        <TableBodyRow>
          <TableBodyCell tdClass={cellStyle}>Aggregate severity namespace</TableBodyCell>
          <TableBodyCell tdClass={cellStyle}
            ><span>{$appStore.webview.doc?.aggregateSeverity.namespace}</span></TableBodyCell
          >
        </TableBodyRow>
      {/if}
    {/if}
    {#if tlp?.label}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>TLP</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}><span class={tlpStyle}>{tlp?.label}</span></TableBodyCell
        >
      </TableBodyRow>
    {/if}
    {#if tlp?.url}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>TLP URL</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}><a href={tlpurl}>{tlp?.url}</a></TableBodyCell>
      </TableBodyRow>
    {/if}
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Category</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{category}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Title</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{title}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Publisher name</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{publisherName}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Publisher category</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{publisherCategory}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Publisher namespace</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{publisherNamespace}</TableBodyCell>
    </TableBodyRow>
    {#if publisherIssuingAuthority}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Publisher issuing authority</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}>{publisherIssuingAuthority}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if publisherContactDetails}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Publisher contact details</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}>{publisherContactDetails}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if lang}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Language</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}>{lang}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if sourceLang}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Source lang</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}>{sourceLang}</TableBodyCell>
      </TableBodyRow>
    {/if}
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Published</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{published}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Last update</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{lastUpdate}</TableBodyCell>
    </TableBodyRow>
    <TableBodyRow>
      <TableBodyCell tdClass={cellStyle}>Tracking Version</TableBodyCell>
      <TableBodyCell tdClass={cellStyle}>{trackingVersion}</TableBodyCell>
    </TableBodyRow>
    {#if $appStore.webview.doc?.status !== Status.final}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Status</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}>{status}</TableBodyCell>
      </TableBodyRow>
    {/if}
    {#if generator}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Generator engine</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}
          ><span>{$appStore.webview.doc?.generator?.engine.name}</span></TableBodyCell
        >
      </TableBodyRow>
    {/if}
    {#if generator?.engine?.version}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Generator engine version</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}
          ><span>{$appStore.webview.doc?.generator?.engine.version}</span></TableBodyCell
        >
      </TableBodyRow>
    {/if}
    {#if generator?.date}
      <TableBodyRow>
        <TableBodyCell tdClass={cellStyle}>Generator date</TableBodyCell>
        <TableBodyCell tdClass={cellStyle}><span>{generator?.date}</span></TableBodyCell>
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
      <Acknowledgements acknowledegements={$appStore.webview.doc?.acknowledgements} />
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
  }
  .tlpred {
    background: #000;
    color: #ff2b2b;
  }
  .tlpamber {
    background: #000;
    color: #ffc000;
  }
  .tlpgreen {
    background: #000;
    color: #33ff00;
  }
</style>
