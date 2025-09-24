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
  import { push } from "svelte-spa-router";

  import { Tabs, TabItem } from "flowbite-svelte";
  import { onMount, tick } from "svelte";

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

  type SingleWebviewDataSection =
    | "vulnerabilitiesOverview"
    | "productTree"
    | "vulnerabilities"
    | "notes"
    | "Acknowledgments"
    | "references"
    | "revisionHistory"
    | undefined;
  type WebviewDataSections = {
    [key: string]: boolean;
    vulnerabilitiesOverview: boolean;
    productTree: boolean;
    vulnerabilities: boolean;
    notes: boolean;
    Acknowledgments: boolean;
    references: boolean;
    revisionHistory: boolean;
  };

  let closeAll = $state(false);
  let tabOpen: WebviewDataSections = $state({
    vulnerabilitiesOverview: true,
    productTree: false,
    vulnerabilities: false,
    notes: false,
    Acknowledgments: false,
    references: false,
    revisionHistory: false
  });

  // When a link in a component inside the Webview is clicked we want the appropriate
  // tab to be opened.
  // Alternatively, it should be possible to drill down the method openTab so the children
  // could use it but then every child has to use that method and that would make the code
  // harder to maintain.
  let forcedTab: SingleWebviewDataSection = $derived.by(() => {
    if (position && position != "") {
      if (position.startsWith("product-")) {
        return "productTree";
      }
      if (position.startsWith("cve-")) {
        return "vulnerabilities";
      }
    }
    return undefined;
  });
  let reallyOpen: WebviewDataSections = $derived.by(() => {
    const open = $state.snapshot(tabOpen);
    if (closeAll || forcedTab) {
      const keys = Object.keys(open);
      for (let i = 0; i < keys.length; i++) {
        const key = keys[i];
        if (closeAll) open[key] = false;
        else open[key] = key === forcedTab;
      }
    }
    return open;
  });

  let innerWidth = $state(0);

  let maxTabs: number = $derived.by(() => {
    return Object.keys(tabOpen).length;
  });
  let tabCount: number = $derived.by(() => {
    let count = 1;
    if (appStore.state.webview.doc && appStore.state.webview.doc["isProductTreePresent"]) count++;
    if (appStore.state.webview.doc && appStore.state.webview.doc["isVulnerabilitiesPresent"])
      count++;
    if (appStore.state.webview.doc?.notes) count++;
    if (appStore.state.webview.doc?.acknowledgments) count++;
    if (appStore.state.webview.doc && appStore.state.webview.doc.references.length > 0) count++;
    if (appStore.state.webview.doc?.isRevisionHistoryPresent) count++;
    return count;
  });
  let missingTabs: number = $derived.by(() => {
    return maxTabs - tabCount;
  });
  // Number of sections that can be shown next to each other
  let screenPhase: number = $derived(
    Math.max(0, Math.floor((innerWidth - widthOffset) / 550 - 1 + missingTabs))
  );

  const openTab = async (tab: string, openRoot = true) => {
    if (openRoot && position && position != "") {
      push(basePath);
    }
    await tick();
    const keys = Object.keys(tabOpen);
    for (let i = 0; i < keys.length; i++) {
      const key = keys[i];
      tabOpen[key] = tab === key;
    }
  };

  const setSelectedItems = () => {
    if (position && position != "") {
      if (position.startsWith("product-")) {
        appStore.setSelectedProduct(position.replace("product-", ""));
      }
      if (position.startsWith("cve-")) {
        appStore.setSelectedCVE(position.replace("cve-", ""));
      }
    }
  };

  // Need variable reallyOpen to return false for every tab once and then to return true
  // for the tab that should be opened. Otherwise the appropriate TabItem would not
  // recognize that it should be opened.
  const tempCloseAll = async () => {
    closeAll = true;
    await tick();
    closeAll = false;
  };

  $effect(() => {
    if (position && position != "") {
      setSelectedItems();
      tempCloseAll();
    }
  });

  let aliases = $derived(appStore.state.webview.doc?.aliases);

  onMount(() => {
    setSelectedItems();
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
  {#if screenPhase < Object.keys(tabOpen).length}
    <Tabs class="mb-2 flex flex-wrap space-x-2 gap-y-2 rtl:space-x-reverse">
      <TabItem
        activeClass={tabItemActiveClass}
        inactiveClass={tabItemInactiveClass}
        open={reallyOpen.vulnerabilitiesOverview}
        onclick={() => openTab("vulnerabilitiesOverview")}
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
      {#if screenPhase < 2}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          open={reallyOpen.productTree}
          onclick={() => openTab("productTree")}
          title="Product tree"
        >
          <div class={sideScroll}>
            <ProductTree {basePath} />
          </div>
        </TabItem>
      {/if}
      {#if screenPhase < 3}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          open={reallyOpen.vulnerabilities}
          onclick={() => openTab("vulnerabilities")}
          title="Vulnerabilities"
        >
          <div class={sideScroll}>
            <Vulnerabilities {basePath} />
          </div>
        </TabItem>
      {/if}
      {#if screenPhase < 4 && appStore.state.webview.doc?.notes}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          open={reallyOpen.notes}
          onclick={() => openTab("notes")}
          title="Notes"
        >
          <div class={sideScroll}>
            <Notes open notes={appStore.state.webview.doc?.notes} />
          </div>
        </TabItem>
      {/if}
      {#if screenPhase < 5 && appStore.state.webview.doc?.acknowledgments}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          open={reallyOpen.Acknowledgments}
          onclick={() => openTab("Acknowledgments")}
          title="Acknowledgments"
        >
          <div class={sideScroll}>
            <Acknowledgments acknowledgments={appStore.state.webview.doc?.acknowledgments} />
          </div>
        </TabItem>
      {/if}
      {#if screenPhase < 6}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          open={reallyOpen.references}
          onclick={() => openTab("references")}
          title="References"
        >
          <div class={sideScroll}>
            <References references={appStore.state.webview.doc?.references} />
          </div>
        </TabItem>
      {/if}
      {#if screenPhase < 7}
        <TabItem
          activeClass={tabItemActiveClass}
          inactiveClass={tabItemInactiveClass}
          open={reallyOpen.revisionHistory}
          onclick={() => openTab("revisionHistory")}
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
{#if screenPhase > 1}
  <div>
    <FakeButton active>Product tree</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <ProductTree {basePath} />
    </div>
  </div>
{/if}
{#if screenPhase > 2}
  <div>
    <FakeButton active>Vulnerabilities</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <Vulnerabilities {basePath} />
    </div>
  </div>
{/if}
{#if screenPhase > 3 && appStore.state.webview.doc?.notes}
  <div>
    <FakeButton active>Notes</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <Notes open notes={appStore.state.webview.doc?.notes} />
    </div>
  </div>
{/if}

{#if screenPhase > 4 && appStore.state.webview.doc?.acknowledgments}
  <div>
    <FakeButton active>Acknowledgments</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <Acknowledgments acknowledgments={appStore.state.webview.doc?.acknowledgments} />
    </div>
  </div>
{/if}

{#if screenPhase > 5}
  <div>
    <FakeButton active>References</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <References references={appStore.state.webview.doc?.references} />
    </div>
  </div>
{/if}

{#if screenPhase > 6}
  <div>
    <FakeButton active>Revision history</FakeButton>
    <div class="mt-2 mb-4 h-px bg-gray-200 dark:bg-gray-700"></div>
    <div class={sideScroll}>
      <RevisionHistory />
    </div>
  </div>
{/if}
