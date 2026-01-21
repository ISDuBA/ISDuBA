<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { request } from "$lib/request";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { Button, Spinner, TableBodyCell, TableBodyRow, TableHeadCell } from "flowbite-svelte";
  import { onMount, tick } from "svelte";
  import WorkflowStateIcon from "$lib/Advisories/WorkflowStateIcon.svelte";
  import { fetchDocumentSSVC } from "./advisory";
  import SSVCBadge from "./SSVC/SSVCBadge.svelte";
  import { push } from "svelte-spa-router";
  import { appStore } from "$lib/store.svelte";
  import { addSlashes } from "$lib/utils";

  interface Related {
    [key: string]: string[];
  }

  interface Props {
    params?: any;
  }

  let { params = null }: Props = $props();

  let document: any | undefined = $state(undefined);
  let documents: any | undefined = $state(undefined);
  let cves: Related | undefined = $state(undefined);
  let ssvc: string | undefined = $state(undefined);
  let isLoading: boolean = $state(false);
  let advisoryState: string | undefined = $state(undefined);
  let loadError: ErrorDetails | null = $state(null);

  let encodedTrackingID = $derived(
    document?.tracking?.id ? encodeURIComponent(addSlashes(document.tracking?.id)) : undefined
  );
  let encodedPublisherNamespace = $derived(
    document?.publisher?.name ? encodeURIComponent(addSlashes(document.publisher?.name)) : undefined
  );

  const baseClass = "text-center px-2 py-2 w-fit min-w-0";

  const loadAdvisoryState = async () => {
    const response = await request(
      `/api/documents?advisories=true&columns=state&query=$tracking_id ${encodedTrackingID} = $publisher "${encodedPublisherNamespace}" = and`,
      "GET"
    );
    if (response.ok) {
      const result = response.content;
      advisoryState = result.documents?.[0].state;
      return result.documents?.[0].state;
    } else if (response.error) {
      loadError = getErrorDetails(`Couldn't load state.`, response);
    }
  };

  const loadDocument = async () => {
    const response = await request(`/api/documents/${params.id}`, "GET");
    if (response.ok) {
      const result = await response.content;
      ({ document } = result);
    } else if (response.error) {
      loadError = getErrorDetails(`Could not load document.`, response);
    }
  };

  onMount(async () => {
    isLoading = true;
    try {
      await loadDocument();
      if (loadError) return;
      await tick();
      if (!encodedTrackingID || !encodedPublisherNamespace) return;
      const result = await fetchDocumentSSVC(encodedTrackingID, encodedPublisherNamespace);
      if (typeof result === "string") {
        ssvc = result;
      } else if (result?.message) {
        loadError = result;
        return;
      }
      loadAdvisoryState();
      if (loadError) return;
      const response = await request(`/api/documents/${params.id}/cve_related`, "GET");
      if (response.ok) {
        cves = {};
        documents = {};
        response.content.forEach((doc: any) => {
          if (cves && !cves[doc.cve]) {
            cves[doc.cve] = [doc];
          } else if (cves) {
            cves[doc.cve].push(doc);
          }

          if (!documents[doc.document_id]) {
            documents[doc.document_id] = doc;
            documents[doc.document_id].cve = [doc.cve];
          } else if (documents) {
            documents[doc.document_id].cve.push(doc.cve);
          }
        });
      } else if (response.error) {
        loadError = getErrorDetails(`Could not load documents.`, response);
      }
    } finally {
      isLoading = false;
    }
  });

  const compare = (doc: any) => {
    appStore.setDiffDocA_ID(params.id);
    appStore.setDiffDocB_ID(doc.document_id);
    push("/diff");
  };

  // Find out if there is a document of the same advisory with same version number but different tracking status because we show
  // tracking status only if there are at least two documents with same version number.
  const hasDocWithSameVersion = (doc: any) => {
    return (
      Object.values(documents).find(
        (d: any) =>
          d.publisher === doc.publisher &&
          d.tracking_id === doc.tracking_id &&
          d.tracking_version === doc.tracking_version &&
          d.tracking_status !== doc.tracking_status
      ) !== undefined
    );
  };
</script>

{#snippet generalInformation(
  state: string | undefined,
  tracking_version: string,
  ssvc: string | undefined,
  tracking_status: string | undefined,
  sameVersion: boolean
)}
  <div class="flex items-center gap-2">
    {#if state}
      <WorkflowStateIcon advisoryState={state}></WorkflowStateIcon>
    {/if}
    <div
      class="flex h-6 min-w-6 items-center justify-center border-1 border-gray-100 px-1 normal-case dark:border-gray-700"
    >
      {tracking_version}
      {#if sameVersion && tracking_status}
        ({tracking_status})
      {/if}
    </div>
    {#if ssvc}
      <SSVCBadge vector={ssvc}></SSVCBadge>
    {/if}
  </div>
{/snippet}

<div style="max-height: 90vh;">
  <ErrorMessage error={loadError}></ErrorMessage>
  {#if isLoading}
    <div class:invisible={!isLoading} class={isLoading ? "loadingFadeIn" : ""}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
  {:else if documents && cves}
    <CustomTable
      tableClass="h-fit w-fit border-separate border-spacing-0"
      tableContainerClass="h-full"
      containerClass="h-full"
      hoverable={false}
      title={`Documents having the same CVEs as ${params.trackingID ?? document?.tracking?.id}`}
      stickyHeaders={true}
    >
      {#snippet tableHeadSlot()}
        <TableHeadCell class="text-center align-top">
          <div class="flex flex-col items-center gap-2">
            <span>{params.trackingID ?? document?.tracking?.id}</span>
            {@render generalInformation(
              advisoryState,
              document?.tracking.version,
              ssvc,
              undefined,
              false
            )}
          </div>
        </TableHeadCell>
        {#each Object.values(documents) as doc}
          {@const d = doc as any}
          {@const sameVersion = hasDocWithSameVersion(d)}
          <TableHeadCell class="text-center align-top">
            <div class="flex h-full flex-col items-center justify-between gap-2">
              <a
                class="text-primary-700 dark:text-primary-400 hover:underline"
                href={`/#/advisories/${encodeURIComponent(d.publisher)}/${encodeURIComponent(d.tracking_id)}/documents/${d.document_id}`}
                >{d.tracking_id}</a
              >
              {@render generalInformation(
                d.state,
                d.tracking_version,
                d.ssvc,
                d.tracking_status,
                sameVersion
              )}
              <Button color="light" size="xs" class="h-6" onclick={() => compare(d)}>
                Compare
              </Button>
            </div>
          </TableHeadCell>
        {/each}
      {/snippet}
      {#snippet mainSlot()}
        {#each Object.keys(cves as Related) as string[] as cve, index (index)}
          <TableBodyRow
            class={cve && cve === params.cve ? "!bg-primary-100 dark:!bg-primary-800" : ""}
          >
            <TableBodyCell class={baseClass}>{cve}</TableBodyCell>
            {#each Object.values(documents) as doc}
              <TableBodyCell class={baseClass}>
                {#if (doc as any).cve.includes(cve)}
                  <i class="bx bx-check"></i>
                {/if}
              </TableBodyCell>
            {/each}
          </TableBodyRow>
        {/each}
      {/snippet}
    </CustomTable>
  {/if}
</div>
