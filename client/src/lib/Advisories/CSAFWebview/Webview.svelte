<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import General from "$lib/Advisories/CSAFWebview/general/General.svelte";
  import ProductTree from "$lib/Advisories/CSAFWebview/producttree/ProductTree.svelte";
  import Vulnerabilities from "$lib/Advisories/CSAFWebview/vulnerabilities/Vulnerabilities.svelte";
  import ValueList from "./ValueList.svelte";
  import RevisionHistory from "./general/RevisionHistory.svelte";
  import Notes from "./notes/Notes.svelte";
  import Acknowledgments from "./acknowledgments/Acknowledgments.svelte";
  import References from "./references/References.svelte";
  import ProductVulnerabilities from "./productvulnerabilities/ProductVulnerabilities.svelte";
  import FakeButton from "./FakeButton.svelte";

  import { Tabs, TabItem } from "flowbite-svelte";
  import { onMount } from "svelte";

  interface Props {
    position: string;
    basePath: string;
    widthOffset: number;
  }
  let { position = "", basePath = "", widthOffset = 0 }: Props = $props();

  const sideScroll = "w-full overflow-y-auto h-max";
  const tabItemActiveClass =
    "h-7 py-1 px-3 border-gray-300 border text-xs bg-gray-200 dark:bg-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg shadow-sm";
  const tabItemInactiveClass =
    "h-7 py-1 px-3 border-gray-300 border text-xs hover:bg-gray-100 dark:bg-gray-800 dark:hover:bg-gray-700 rounded-lg";
  const webviewDataSections = [
    "vulnerabilitiesOverview",
    "productTree",
    "vulnerabilities",
    "notes",
    "Acknowledgments",
    "references",
    "revisionHistory"
  ] as const;
  type WebviewDataSections = (typeof webviewDataSections)[number];

  let innerWidth = $state(0);
  let screenPhase: number = $derived(Math.max(0, Math.floor((innerWidth - widthOffset) / 550 - 1)));
  let placeToPhase = $state(
    Object.fromEntries(webviewDataSections.map((key) => [key, { show: false, phase: 9 }])) as {
      [place in WebviewDataSections]: { show: boolean; phase: number };
    }
  );

  let tabOpen: { [key in WebviewDataSections]: boolean } = $state({
    vulnerabilitiesOverview: true,
    productTree: false,
    vulnerabilities: false,
    notes: false,
    Acknowledgments: false,
    references: false,
    revisionHistory: false
  });

  const updateUI = async () => {
    // This is a hack
    setTimeout(() => {
      if (position.startsWith("product-")) {
        updateTabOpen("productTree");
        appStore.setSelectedProduct(position.replace("product-", ""));
      }
      if (position.startsWith("cve-")) {
        updateTabOpen("vulnerabilities");
        appStore.setSelectedCVE(position.replace("cve-", ""));
      }
    }, 300);
  };

  const updateTabOpen = (open?: WebviewDataSections) => {
    let openedTab: WebviewDataSections = (Object.entries(tabOpen).find((tab) => tab[1]) ?? [
      "vulnerabilitiesOverview"
    ])[0] as WebviewDataSections;
    const canBeOpened = (tab: WebviewDataSections) =>
      tab === "vulnerabilitiesOverview" ||
      (screenPhase < placeToPhase[tab].phase && placeToPhase[tab].show);
    if (open && canBeOpened(open)) {
      if (open !== openedTab) {
        tabOpen[open] = true;
        tabOpen[openedTab] = false;
      }
    } else if (openedTab !== "vulnerabilitiesOverview" && !canBeOpened(openedTab)) {
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
      appStore.state.webview.doc && appStore.state.webview.doc["isProductTreePresent"]
    );
    increment(
      "vulnerabilities",
      appStore.state.webview.doc && appStore.state.webview.doc["isVulnerabilitiesPresent"]
    );
    increment("notes", appStore.state.webview.doc?.notes);
    increment("Acknowledgments", appStore.state.webview.doc?.acknowledgments);
    increment(
      "references",
      appStore.state.webview.doc && appStore.state.webview.doc.references.length > 0
    );
    increment("revisionHistory", appStore.state.webview.doc?.isRevisionHistoryPresent);
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

  $effect(() => {
    if (position && position != "") {
      updateUI();
    }
  });

  let aliases = $derived(appStore.state.webview.doc?.aliases);

  onMount(() => {
    updatePlaces();
  });
</script>

<svelte:window bind:innerWidth />

<div class="flex w-full flex-col">
  {#if appStore.state.webview.doc}
    <div class="mb-4 w-full">
      <General />
    </div>
  {/if}
  {#if aliases}
    <div class="mb-4">
      <ValueList label="Aliases" values={aliases} />
    </div>
  {/if}
  {#if showTab(placeToPhase.revisionHistory)}
    <Tabs class="mb-2 flex flex-wrap space-x-2 gap-y-2 rtl:space-x-reverse">
      <TabItem
        activeClass={tabItemActiveClass}
        inactiveClass={tabItemInactiveClass}
        bind:open={tabOpen.vulnerabilitiesOverview}
        title="Overview"
      >
        {#if appStore.state.webview.doc?.productVulnerabilities.length > 1}
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
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          bind:open={tabOpen.productTree}
          title="Product tree"
        >
          <div class={sideScroll}>
            <ProductTree {basePath} />
          </div>
        </TabItem>
      {/if}
      {#if showTab(placeToPhase.vulnerabilities)}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          bind:open={tabOpen.vulnerabilities}
          title="Vulnerabilities"
        >
          <div class={sideScroll}>
            <Vulnerabilities {basePath} />
          </div>
        </TabItem>
      {/if}
      {#if showTab(placeToPhase.notes) && appStore.state.webview.doc?.notes}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          bind:open={tabOpen.notes}
          title="Notes"
        >
          <div class={sideScroll}>
            <Notes open notes={appStore.state.webview.doc?.notes} />
          </div>
        </TabItem>
      {/if}
      {#if showTab(placeToPhase.Acknowledgments) && appStore.state.webview.doc?.acknowledgments}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          bind:open={tabOpen.Acknowledgments}
          title="Acknowledgments"
        >
          <div class={sideScroll}>
            <Acknowledgments acknowledgments={appStore.state.webview.doc?.acknowledgments} />
          </div>
        </TabItem>
      {/if}
      {#if showTab(placeToPhase.references)}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          bind:open={tabOpen.references}
          title="References"
        >
          <div class={sideScroll}>
            <References references={appStore.state.webview.doc?.references} />
          </div>
        </TabItem>
      {/if}
      {#if showTab(placeToPhase.revisionHistory)}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          bind:open={tabOpen.revisionHistory}
          title="Revision history"
        >
          <div class={sideScroll}>
            <RevisionHistory />
          </div>
        </TabItem>
      {/if}
    </Tabs>
  {:else}
    <div>
      <FakeButton active>Overview</FakeButton>
      <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
      <div class={sideScroll}>
        {#if appStore.state.webview.doc?.productVulnerabilities.length > 1}
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
    <FakeButton active>Product tree</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <ProductTree {basePath} />
    </div>
  </div>
{/if}
{#if showArea(placeToPhase.vulnerabilities)}
  <div>
    <FakeButton active>Vulnerabilities</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <Vulnerabilities {basePath} />
    </div>
  </div>
{/if}
{#if showArea(placeToPhase.notes) && appStore.state.webview.doc?.notes}
  <div>
    <FakeButton active>Notes</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <Notes open notes={appStore.state.webview.doc?.notes} />
    </div>
  </div>
{/if}

{#if showArea(placeToPhase.Acknowledgments) && appStore.state.webview.doc?.acknowledgments}
  <div>
    <FakeButton active>Acknowledgments</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <Acknowledgments acknowledgments={appStore.state.webview.doc?.acknowledgments} />
    </div>
  </div>
{/if}

{#if showArea(placeToPhase.references)}
  <div>
    <FakeButton active>References</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <References references={appStore.state.webview.doc?.references} />
    </div>
  </div>
{/if}

{#if showArea(placeToPhase.revisionHistory)}
  <div>
    <FakeButton active>Revision history</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <RevisionHistory />
    </div>
  </div>
{/if}
