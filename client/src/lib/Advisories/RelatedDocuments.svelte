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
  import { tick } from "svelte";
  import WorkflowStateIcon from "$lib/Advisories/WorkflowStateIcon.svelte";
  import { fetchDocumentSSVC, isResultConsistent } from "$lib/Advisories/advisory.svelte";
  import SSVCBadge from "./SSVC/SSVCBadge.svelte";
  import { push } from "$routes/router.svelte";
  import { appStore } from "$lib/store.svelte";
  import { addSlashes } from "$lib/utils";
  import Link from "$lib/Components/Link.svelte";
  import InconsistencyMessage from "./InconsistencyMessage.svelte";
  import { AlertCircle, Check } from "@boxicons/svelte";
  import type { WorkflowState } from "$lib/workflow";

  interface BasicRelatedDocument {
    document_id: number;
    publisher: string;
    ssvc?: string;
    state: WorkflowState;
    title: string;
    tracking_id: string;
    tracking_status: string;
    tracking_version: string;
  }

  interface RelatedDocumentFromBackend extends BasicRelatedDocument {
    cve: string;
  }

  interface RelatedDocument extends BasicRelatedDocument {
    cve: string[];
  }

  interface RelatedCVEResponse {
    error?: string;
    ok: boolean;
    content: RelatedDocumentFromBackend[];
  }

  interface Props {
    params?: any;
  }

  let { params = null }: Props = $props();

  const uid = $props.id();

  let document: any | undefined = $state(undefined);
  let documents: RelatedDocument[] | undefined = $state(undefined);
  let cves: string[] | undefined = $state(undefined);
  let ssvc: string | undefined = $state(undefined);
  let isLoading: boolean = $state(false);
  let advisoryState: string | undefined = $state(undefined);
  let loadError: ErrorDetails | null = $state(null);
  let isInconsistent = $state(false);
  let max = $state(50);

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
    isInconsistent = false;
    const response = await request(`/api/documents/${params.id}`, "GET");
    if (response.ok) {
      const result = await response.content;
      ({ document } = result);
      if (!isResultConsistent(params, document)) {
        isInconsistent = true;
      }
    } else if (response.error) {
      loadError = getErrorDetails(`Could not load document.`, response);
    }
  };

  $effect(() => {
    if (params) {
      setTimeout(async () => {
        isLoading = true;
        document = undefined;
        documents = undefined;
        cves = undefined;
        loadError = null;
        try {
          await loadDocument();
          if (loadError || isInconsistent) {
            isLoading = false;
            return;
          }
          await tick();
          if (!encodedTrackingID || !encodedPublisherNamespace) return;
          const result = await fetchDocumentSSVC(params.id);
          if (typeof result === "string") {
            ssvc = result;
          } else if (result?.message) {
            loadError = result;
            return;
          }
          loadAdvisoryState();
          if (loadError) return;
          const response: RelatedCVEResponse = (await request(
            `/api/documents/${params.id}/cve_related`,
            "GET"
          )) as RelatedCVEResponse;
          if (response.ok) {
            const newCves: string[] = [];
            const newDocuments: RelatedDocument[] = [];

            for (let i = 0; i < response.content.length; i++) {
              const doc = response.content[i];
              const docIndex = newDocuments.findIndex((d: RelatedDocument) => {
                return doc.document_id === d.document_id;
              });
              if (docIndex === -1) {
                newDocuments.push({
                  ...doc,
                  cve: [doc.cve]
                });
              } else {
                newDocuments[docIndex].cve.push(doc.cve);
              }

              if (!newCves.includes(doc.cve)) {
                newCves.push(doc.cve);
              }
            }

            // We are more interested in documents with more CVEs so they should come first
            newDocuments.sort((a: RelatedDocument, b: RelatedDocument) => {
              if (a.cve.length < b.cve.length) {
                return -1;
              } else if (a.cve.length > b.cve.length) {
                return 1;
              }
              return 0;
            });
            documents = newDocuments;
            cves = newCves;
          } else if (response.error) {
            loadError = getErrorDetails(`Could not load documents.`, response);
          }
        } finally {
          isLoading = false;
        }
      }, 0);
    }
  });

  const compare = (doc: any) => {
    appStore.setDiffDocA_ID(params.id);
    appStore.setDiffDocB_ID(doc.document_id);
    push("/diff");
  };

  // Get document as an object that can be compared with the other documents.
  const getComparableDocument = () => {
    return {
      publisher: document.publisher.name,
      tracking_id: document.tracking.id,
      tracking_version: document.tracking.version,
      tracking_status: document.tracking.status
    };
  };

  // Find out if there is a document of the same advisory with same version number but different tracking status because we show
  // tracking status only if there are at least two documents with same version number.
  const hasDocWithSameVersion = (doc: any) => {
    if (!documents) return false;
    const docs = [...documents, getComparableDocument()];
    return (
      docs.find(
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
  {:else if document && isInconsistent}
    <InconsistencyMessage {document} {params} relatedDocuments />
  {:else if document && documents && cves}
    {#if documents?.length === 0}
      <div class="mb-2 font-bold">
        <AlertCircle aria-hidden="true" />
        <span>The document {document?.tracking?.id} has no related documents.</span>
      </div>
    {:else}
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
                document?.tracking.status,
                hasDocWithSameVersion(getComparableDocument())
              )}
            </div>
          </TableHeadCell>
          {#each documents?.slice(0, max) as doc, i (`relateddocuments-1-${uid}-${i}`)}
            {@const d = doc as RelatedDocument}
            {@const sameVersion = hasDocWithSameVersion(d)}
            <TableHeadCell class="text-center align-top">
              <div class="flex h-full flex-col items-center justify-between gap-2">
                <Link
                  class="text-primary-700 dark:text-primary-400 hover:underline"
                  href={`/#/advisories/${encodeURIComponent(d.publisher)}/${encodeURIComponent(d.tracking_id)}/documents/${d.document_id}`}
                  >{d.tracking_id}</Link
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
          {#if documents && max < documents.length}
            <TableHeadCell class="">
              <Button
                onclick={() => {
                  if (!documents) return;
                  max = max + 50;
                }}
                class="h-6 text-nowrap"
                color="light"
              >
                Load more...
              </Button>
            </TableHeadCell>
          {/if}
        {/snippet}
        {#snippet mainSlot()}
          {#each cves as cve, j (`relateddocuments-1-${uid}-${j}`)}
            <TableBodyRow
              class={cve && cve === params.cve ? "!bg-primary-100 dark:!bg-primary-800" : ""}
            >
              <TableBodyCell
                class={`${baseClass} ${cve && cve === params.cve ? "!font-bold" : ""}`}
              >
                {cve}
              </TableBodyCell>
              {#each documents?.slice(0, max) as doc, k (`relateddocuments-1-${uid}-${k}`)}
                <TableBodyCell class={baseClass}>
                  <div class="flex justify-center">
                    {#if (doc as RelatedDocument).cve.includes(cve)}
                      <Check />
                    {/if}
                  </div>
                </TableBodyCell>
              {/each}
            </TableBodyRow>
          {/each}
        {/snippet}
      </CustomTable>
    {/if}
  {/if}
</div>
