<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store";
  import Back from "$lib/CSAFWebview/Back.svelte";
  import Collapsible from "$lib/CSAFWebview/Collapsible.svelte";
  import General from "$lib/CSAFWebview/general/General.svelte";
  import ProductTree from "$lib/CSAFWebview/producttree/ProductTree.svelte";
  import ProductVulnerabilities from "$lib/CSAFWebview/productvulnerabilities/ProductVulnerabilities.svelte";
  import Vulnerabilities from "$lib/CSAFWebview/vulnerabilities/Vulnerabilities.svelte";
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

{#if isCSAF}
  {#if $appStore.webview.doc}
    <Collapsible header="General" open={$appStore.webview.ui.isGeneralSectionVisible}>
      <General />
    </Collapsible>
  {/if}
  {#if $appStore.webview.doc?.productVulnerabilities.length > 1}
    <Collapsible
      header="Vulnerabilities overview"
      open={$appStore.webview.ui.isVulnerabilitiesOverviewVisible}
    >
      <ProductVulnerabilities />
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
      <ProductTree />
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
      <Vulnerabilities />
    </Collapsible>
  {/if}

  {#if $appStore.webview.ui.history.length > 0}
    <Back />
  {/if}
{/if}
