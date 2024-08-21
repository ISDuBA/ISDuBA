<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import General from "$lib/Advisories/CSAFWebview/general/General.svelte";
  import ProductTree from "$lib/Advisories/CSAFWebview/producttree/ProductTree.svelte";
  import Vulnerabilities from "$lib/Advisories/CSAFWebview/vulnerabilities/Vulnerabilities.svelte";
  import ValueList from "./ValueList.svelte";
  import RevisionHistory from "./general/RevisionHistory.svelte";
  import Notes from "./notes/Notes.svelte";
  import Acknowledgements from "./acknowledgements/Acknowledgements.svelte";
  import References from "./references/References.svelte";
  import ProductVulnerabilities from "./productvulnerabilities/ProductVulnerabilities.svelte";

  export let position = "";
  export let basePath = "";

  const updateUI = async () => {
    // This is a hack
    setTimeout(() => {
      if (position.startsWith("product-")) {
        appStore.setProductTreeSectionVisible();
        appStore.setSelectedProduct(position.replace("product-", ""));
      }
      if (position.startsWith("cve-")) {
        appStore.setSelectedCVE(position.replace("cve-", ""));
        appStore.setVulnerabilitiesSectionVisible();
      }
    }, 300);
  };

  $: if (position && position != "") {
    updateUI();
  }

  $: aliases = $appStore.webview.doc?.aliases;

  $: isCSAF = !(
    !$appStore.webview.doc?.isRevisionHistoryPresent &&
    !$appStore.webview.doc?.isDocPresent &&
    !$appStore.webview.doc?.isProductTreePresent &&
    !$appStore.webview.doc?.isPublisherPresent &&
    !$appStore.webview.doc?.isTLPPresent &&
    !$appStore.webview.doc?.isTrackingPresent &&
    !$appStore.webview.doc?.isVulnerabilitiesPresent
  );
</script>

<div class="flex flex-col">
  {#if isCSAF}
    {#if $appStore.webview.doc}
      <div class="mb-4">
        <General />
      </div>
    {/if}
    {#if $appStore.webview.doc?.productVulnerabilities.length > 1}
      <Collapsible
        header="Vulnerabilities overview"
        open={$appStore.webview.ui.isVulnerabilitiesOverviewVisible}
      >
        <ProductVulnerabilities {basePath} />
      </Collapsible>
    {:else}
      <h2>No Vulnerabilities overview</h2>
      (As no products are connected to vulnerabilities.)
    {/if}
    {#if $appStore.webview.doc && $appStore.webview.doc["isProductTreePresent"]}
      <Collapsible
        header="Product tree"
        onOpen={() => {
          appStore.setProductTreeOpen();
        }}
        open={$appStore.webview.ui.isProductTreeVisible}
        onClose={() => {
          appStore.setProductTreeSectionInVisible();
          appStore.resetSelectedProduct();
          appStore.setProductTreeClosed();
        }}
      >
        <ProductTree {basePath} />
      </Collapsible>
    {/if}
    {#if $appStore.webview.doc && $appStore.webview.doc["isVulnerabilitiesPresent"]}
      <Collapsible
        header="Vulnerabilities"
        open={$appStore.webview.ui.isVulnerabilitiesSectionVisible}
        onClose={() => {
          appStore.setVulnerabilitiesSectionInvisible();
        }}
      >
        <Vulnerabilities {basePath} />
      </Collapsible>
    {/if}
  {/if}

  {#if aliases}
    <ValueList label="Aliases" values={aliases} />
  {/if}
  {#if $appStore.webview.doc?.notes}
    <div>
      <Collapsible header="Notes" level={2}>
        <Notes notes={$appStore.webview.doc?.notes} />
      </Collapsible>
    </div>
  {/if}

  {#if $appStore.webview.doc?.acknowledgements}
    <div>
      <Collapsible header="Acknowledgements" level={2}>
        <Acknowledgements acknowledegements={$appStore.webview.doc?.acknowledgements} />
      </Collapsible>
    </div>
  {/if}

  {#if $appStore.webview.doc && $appStore.webview.doc.references.length > 0}
    <div>
      <Collapsible header="References" level={2}>
        <References references={$appStore.webview.doc?.references} />
      </Collapsible>
    </div>
  {/if}

  {#if $appStore.webview.doc?.isRevisionHistoryPresent}
    <div>
      <Collapsible
        header="Revision history"
        level={2}
        open={$appStore.webview.ui.isRevisionHistoryVisible}
      >
        <RevisionHistory />
      </Collapsible>
    </div>
  {/if}
</div>
