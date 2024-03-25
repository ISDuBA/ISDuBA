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
{/if}
