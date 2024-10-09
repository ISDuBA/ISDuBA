<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store";
  import General from "$lib/Advisories/CSAFWebview/general/General.svelte";
  import ProductTree from "$lib/Advisories/CSAFWebview/producttree/ProductTree.svelte";
  import Vulnerabilities from "$lib/Advisories/CSAFWebview/vulnerabilities/Vulnerabilities.svelte";
  import ValueList from "./ValueList.svelte";
  import RevisionHistory from "./general/RevisionHistory.svelte";
  import Notes from "./notes/Notes.svelte";
  import Acknowledgements from "./acknowledgements/Acknowledgements.svelte";
  import References from "./references/References.svelte";
  import ProductVulnerabilities from "./productvulnerabilities/ProductVulnerabilities.svelte";

  import { Tabs, TabItem } from "flowbite-svelte";

  export let position = "";
  export let basePath = "";

  const sideScroll = "w-full overflow-y-auto h-max";
  const webviewDataSections = [
    "vulnerabilitiesOverview",
    "productTree",
    "vulnerabilities",
    "notes",
    "acknowledgements",
    "references",
    "revisionHistory"
  ] as const;
  type WebviewDataSections = (typeof webviewDataSections)[number];

  let screenPhase: number = 0;
  let placeToPhase = Object.fromEntries(
    webviewDataSections.map((key) => [key, { show: false, phase: 9 }])
  ) as { [place in WebviewDataSections]: { show: boolean; phase: number } };
  let isCSAF: boolean = false;

  let tabOpen: { [key in WebviewDataSections]: boolean } = {
    vulnerabilitiesOverview: true,
    productTree: false,
    vulnerabilities: false,
    notes: false,
    acknowledgements: false,
    references: false,
    revisionHistory: false
  };

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

  const updateTabOpen = () => {
    let openedTab: WebviewDataSections = (Object.entries(tabOpen).find((tab) => tab[1]) ?? [
      "vulnerabilitiesOverview"
    ])[0] as WebviewDataSections;
    if (
      openedTab !== "vulnerabilitiesOverview" &&
      (screenPhase >= placeToPhase[openedTab].phase || !placeToPhase[openedTab].show)
    ) {
      tabOpen.vulnerabilitiesOverview = true;
      tabOpen[openedTab] = false;
    }
  };

  const updatePlaces = () => {
    let phaseObject: { [key in WebviewDataSections]?: { show: boolean; phase: number } } = {};
    let next = 1;
    let increment = (place: WebviewDataSections, show: any) => {
      phaseObject[place] = { show: !!show, phase: next };
      if (show) {
        next++;
      }
    };
    increment(
      "productTree",
      $appStore.webview.doc && $appStore.webview.doc["isProductTreePresent"]
    );
    increment(
      "vulnerabilities",
      $appStore.webview.doc && $appStore.webview.doc["isVulnerabilitiesPresent"]
    );
    increment("notes", $appStore.webview.doc?.notes);
    increment("acknowledgements", $appStore.webview.doc?.acknowledgements);
    increment("references", $appStore.webview.doc && $appStore.webview.doc.references.length > 0);
    increment("revisionHistory", $appStore.webview.doc?.isRevisionHistoryPresent);
    placeToPhase = {
      vulnerabilitiesOverview: { show: true, phase: next },
      ...phaseObject
    } as { [key in WebviewDataSections]: { show: boolean; phase: number } };
    updateTabOpen();
  };

  const showTab = (place: { show: boolean; phase: number }) =>
    place.show && place.phase > screenPhase;
  const showArea = (place: { show: boolean; phase: number }) =>
    place.show && place.phase <= screenPhase;

  $: if (position && position != "") {
    updateUI();
  }

  $: aliases = $appStore.webview.doc?.aliases;

  $: innerWidth = 0;
  $: {
    let oldPhase = screenPhase;
    screenPhase = Math.max(0, Math.floor(innerWidth / 480 - 2));
    if (oldPhase !== screenPhase) {
      updatePlaces();
    }
  }
  $: {
    isCSAF = !!(
      $appStore.webview.doc?.isRevisionHistoryPresent ||
      $appStore.webview.doc?.isDocPresent ||
      $appStore.webview.doc?.isProductTreePresent ||
      $appStore.webview.doc?.isPublisherPresent ||
      $appStore.webview.doc?.isTLPPresent ||
      $appStore.webview.doc?.isTrackingPresent ||
      $appStore.webview.doc?.isVulnerabilitiesPresent
    );
    updatePlaces();
  }
</script>

<svelte:window bind:innerWidth />

<div class="grid auto-cols-fr grid-flow-col gap-10">
  {#if isCSAF}
    <div class="flex w-full flex-col">
      {#if $appStore.webview.doc}
        <div class="mb-4 w-full">
          <General />
        </div>
      {/if}
      {#if aliases}
        <ValueList label="Aliases" values={aliases} />
      {/if}
      {#if showTab(placeToPhase.revisionHistory)}
        <Tabs tabStyle="pill">
          <TabItem bind:open={tabOpen.vulnerabilitiesOverview} title="Vulnerabilities overview">
            {#if $appStore.webview.doc?.productVulnerabilities.length > 1}
              <div class={sideScroll}>
                <ProductVulnerabilities {basePath} />
              </div>
            {:else}
              <i>
                <h2>No Vulnerabilities overview</h2>
                (As no products are connected to vulnerabilities.)
              </i>
            {/if}
          </TabItem>
          {#if showTab(placeToPhase.productTree)}
            <TabItem bind:open={tabOpen.productTree} title="Product tree">
              <div class={sideScroll}>
                <ProductTree {basePath} />
              </div>
            </TabItem>
          {/if}
          {#if showTab(placeToPhase.vulnerabilities)}
            <TabItem bind:open={tabOpen.vulnerabilities} title="Vulnerabilities">
              <div class={sideScroll}>
                <Vulnerabilities {basePath} />
              </div>
            </TabItem>
          {/if}
          {#if showTab(placeToPhase.notes) && $appStore.webview.doc?.notes}
            <TabItem bind:open={tabOpen.notes} title="Notes">
              <div class={sideScroll}>
                <Notes notes={$appStore.webview.doc?.notes} />
              </div>
            </TabItem>
          {/if}
          {#if showTab(placeToPhase.acknowledgements) && $appStore.webview.doc?.acknowledgements}
            <TabItem bind:open={tabOpen.acknowledgements} title="Acknowledgements">
              <div class={sideScroll}>
                <Acknowledgements acknowledegements={$appStore.webview.doc?.acknowledgements} />
              </div>
            </TabItem>
          {/if}
          {#if showTab(placeToPhase.references)}
            <TabItem bind:open={tabOpen.references} title="References">
              <div class={sideScroll}>
                <References references={$appStore.webview.doc?.references} />
              </div>
            </TabItem>
          {/if}
          {#if showTab(placeToPhase.revisionHistory)}
            <TabItem bind:open={tabOpen.revisionHistory} title="Revision history">
              <div class={sideScroll}>
                <RevisionHistory />
              </div>
            </TabItem>
          {/if}
        </Tabs>
      {:else}
        <div>
          <span class="text-xl">Vulnerabilities Overwiew</span>
          <div class={sideScroll}>
            {#if $appStore.webview.doc?.productVulnerabilities.length > 1}
              <ProductVulnerabilities {basePath} />
            {:else}
              <i>
                <h2>No Vulnerabilities overview</h2>
                (As no products are connected to vulnerabilities.)
              </i>
            {/if}
          </div>
        </div>
      {/if}
    </div>
    {#if showArea(placeToPhase.productTree)}
      <div>
        <span class="text-xl">Product tree</span>
        <div class={sideScroll}>
          <ProductTree {basePath} />
        </div>
      </div>
    {/if}
    {#if showArea(placeToPhase.vulnerabilities)}
      <div>
        <span class="text-xl">Vulnerabilities</span>
        <div class={sideScroll}>
          <Vulnerabilities {basePath} />
        </div>
      </div>
    {/if}
  {/if}
  {#if showArea(placeToPhase.notes) && $appStore.webview.doc?.notes}
    <div>
      <span class="text-xl">Notes</span>
      <div class={sideScroll}>
        <Notes notes={$appStore.webview.doc?.notes} />
      </div>
    </div>
  {/if}

  {#if showArea(placeToPhase.acknowledgements) && $appStore.webview.doc?.acknowledgements}
    <div>
      <span class="text-xl">Acknowledgements</span>
      <div class={sideScroll}>
        <Acknowledgements acknowledegements={$appStore.webview.doc?.acknowledgements} />
      </div>
    </div>
  {/if}

  {#if showArea(placeToPhase.references)}
    <div>
      <span class="text-xl">References</span>
      <div class={sideScroll}>
        <References references={$appStore.webview.doc?.references} />
      </div>
    </div>
  {/if}

  {#if showArea(placeToPhase.revisionHistory)}
    <div>
      <span class="text-xl">Revision history</span>
      <div class={sideScroll}>
        <RevisionHistory />
      </div>
    </div>
  {/if}
</div>
