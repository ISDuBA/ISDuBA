<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  /* eslint-disable svelte/no-at-html-tags */
  import {
    Button,
    Dropdown,
    Label,
    PaginationItem,
    Select,
    TableBody,
    TableBodyCell,
    TableHead,
    TableHeadCell,
    Table,
    Img
  } from "flowbite-svelte";
  import DOMPurify from "dompurify";
  import { tablePadding, title, publisher, searchColumnName, tdClass } from "$lib/Table/defaults";
  import { Spinner } from "flowbite-svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { ADMIN, EDITOR, IMPORTER, REVIEWER } from "$lib/workflow";
  import { getAllowedWorkflowChanges, isRoleIncluded } from "$lib/permissions";
  import { appStore } from "$lib/store.svelte";
  import { getPublisher } from "$lib/publisher";
  import CIconButton from "$lib/Components/CIconButton.svelte";
  import SsvcBadge from "$lib/Advisories/SSVC/SSVCBadge.svelte";
  import { SEARCHTYPES } from "$lib/Queries/query";
  import CCheckbox from "$lib/Components/CCheckbox.svelte";
  import { areArraysEqual } from "$lib/utils";
  import DeleteModal from "./DeleteModal.svelte";
  import { updateMultipleStates } from "$lib/Advisories/advisory";
  import CVSS from "$lib/Advisories/CSAFWebview/general/CVSS.svelte";
  import type { SearchParameters } from "$lib/Search/search";

  const toggleRow = (i: number) => {
    openRow = openRow === i ? null : i;
  };
  let changeWorkflowStateError: ErrorDetails | null = $state(null);
  interface Props {
    columns: string[];
    query?: string;
    tableType: SEARCHTYPES;
    orderBy?: string[];
    loading: boolean;
    openRow?: number | null;
    count: number;
    offset: number;
    limit: number;
    currentPage: number;
    numberOfPages: number;
    error: ErrorDetails | null;
    documents: any;
    dataChanged: () => void;
    last: () => void;
    setSearchParameters: (SearchParameters: SearchParameters) => void;
  }

  let {
    columns,
    query = "",
    tableType,
    documents = null,
    error = null,
    loading = false,
    openRow = $bindable(null),
    orderBy = ["title"],
    count = $bindable(0),
    offset = 0,
    limit = 10,
    currentPage = 1,
    numberOfPages = $bindable(0),
    dataChanged = () => {},
    last = () => {},
    setSearchParameters = (_paginationParameters: SearchParameters) => {}
  }: Props = $props();
  const tdClassRelative = `${tdClass} relative`;

  let disableDiffButtons = $derived(
    appStore.state.app.diff.docA_ID !== undefined && appStore.state.app.diff.docB_ID !== undefined
  );

  let documentIDs = $derived(documents?.map((d: any) => d.id) ?? []);

  let areAllSelected = $derived(
    documents &&
      areArraysEqual(documentIDs, Array.from(appStore.state.app.selectedDocumentIDs.keys()))
  );

  let selectedDocuments = $derived(
    appStore.state.app.search.results?.filter((d: any) =>
      appStore.state.app.selectedDocumentIDs.has(d.id)
    ) ?? []
  );
  let allowedWorkflowStateChanges = $derived(
    getAllowedWorkflowChanges(selectedDocuments?.map((d: any) => d.state) ?? [])
  );
  let workflowOptions = $derived(
    allowedWorkflowStateChanges.map((c) => {
      return { name: c.to, value: c.to };
    })
  );
  let isMultiSelectionAllowed = $derived(
    isRoleIncluded(appStore.getRoles(), [EDITOR, IMPORTER, ADMIN, REVIEWER]) &&
      ((tableType !== SEARCHTYPES.EVENT && appStore.isAdmin()) ||
        tableType === SEARCHTYPES.ADVISORY)
  );
  let areThereAnyComments = $derived(
    tableType === SEARCHTYPES.EVENT && documents?.find((d: any) => d.event === "add_comment")
  );

  let selectedState: any = $state(null);
  let dropdownOpen = $state(false);
  const selectClass =
    "max-w-96 w-fit text-gray-900 disabled:text-gray-400 bg-gray-50 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:disabled:text-gray-500 dark:focus:ring-primary-500 dark:focus:border-primary-500";

  const getAdvisoryLink = (item: any) =>
    `/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`;
  const getAdvisoryAnchorLink = (item: any) => "#" + getAdvisoryLink(item);

  const changeWorkflowState = async () => {
    if (!selectedDocuments || selectedDocuments.length < 0) return;
    const changes: any[] = [];
    selectedDocuments?.forEach((doc: any) => {
      changes.push({
        state: selectedState,
        publisher: doc.publisher,
        tracking_id: doc.tracking_id
      });
    });
    changeWorkflowStateError = null;
    const response = await updateMultipleStates(changes);
    if (response.ok) {
      dataChanged();
      dropdownOpen = false;
      selectedState = undefined;
    } else if (response.error) {
      changeWorkflowStateError = getErrorDetails("Couldn't change state.", response);
    }
  };

  let innerWidth = $state(0);

  const getColumnDisplayName = (column: string): string => {
    let names: { [key: string]: string } = {
      id: "ID",
      tracking_id: "TRACKING ID",
      tracking_status: "TRACKING STATUS",
      version: "VERSION",
      publisher: "PUBLISHER",
      current_release_date: "CURRENT RELEASE",
      initial_release_date: "INITIAL RELEASE",
      title: "TITLE",
      tlp: "TLP",
      cvss_v2_score: "CVSS2",
      cvss_v3_score: "CVSS3",
      ssvc: "SSVC",
      four_cves: "CVES",
      state: "STATE"
    };

    return names[column] ?? column;
  };

  let isAdmin = $derived(isRoleIncluded(appStore.getRoles(), [ADMIN]));

  const previous = async () => {
    if (offset - limit >= 0) {
      setSearchParameters({
        currentPage: currentPage - 1
      });
    }
  };
  const next = async () => {
    if (offset + limit <= count) {
      setSearchParameters({
        currentPage: currentPage + 1
      });
    }
  };

  const first = async () => {
    setSearchParameters({
      currentPage: 1
    });
  };

  const switchSort = async (column: string) => {
    let newOrderBy = structuredClone($state.snapshot(orderBy));
    let found = newOrderBy.find((c) => c === column);
    let foundMinus = newOrderBy.find((c) => c === "-" + column);
    if (foundMinus) {
      newOrderBy = newOrderBy.filter((c) => c !== "-" + column);
    }
    if (found) {
      newOrderBy = newOrderBy.map((c) => (c === column ? `-${column}` : c));
    }
    if (!found && !foundMinus) {
      newOrderBy.push(column);
    }
    setSearchParameters({
      orderBy: newOrderBy
    });
  };

  const onDeleted = async () => {
    dataChanged();
  };

  const getColumnOrder = (orderBy: string[], column: string): string => {
    let index = orderBy.indexOf(column);
    let indexMinus = orderBy.indexOf("-" + column);

    if (indexMinus >= 0) {
      return indexMinus + 1 + "";
    }
    if (index >= 0) {
      return index + 1 + "";
    }

    return "";
  };
</script>

<svelte:window bind:innerWidth />

<DeleteModal {onDeleted} documents={appStore.state.app.documentsToDelete || []} type={tableType}
></DeleteModal>

<div class="flex-grow">
  <div class="mt-2 mb-2 flex flex-row items-baseline justify-between">
    {#if documents?.length > 0}
      <div class="flex flex-row items-baseline gap-8">
        {#if isMultiSelectionAllowed}
          <div class="flex items-center gap-2">
            {#if appStore.isAdmin()}
              <Button
                onclick={() => {
                  appStore.setDocumentsToDelete(selectedDocuments);
                  appStore.setIsDeleteModalOpen(true);
                }}
                class="!p-2"
                color="light"
                disabled={!selectedDocuments || selectedDocuments.length === 0}
              >
                <i class="bx bx-trash text-red-600"></i>
              </Button>
            {/if}
            {#if tableType === SEARCHTYPES.ADVISORY}
              <Button
                class="!p-2"
                color="light"
                disabled={workflowOptions.length === 0}
                id="state-icon"
              >
                <i class="bx bx-git-commit text-black-700 dark:text-gray-300"></i>
              </Button>
              <Dropdown
                bind:isOpen={dropdownOpen}
                ontoggle={(event) => {
                  if (event.newState) {
                    changeWorkflowStateError = null;
                  }
                }}
                placement="top-start"
                triggeredBy="#state-icon"
                class="z-50 w-full max-w-sm divide-y divide-gray-100 rounded border border-gray-300 p-4 shadow dark:divide-gray-700 dark:bg-gray-800"
              >
                <div class="flex flex-col gap-3">
                  <div class="flex w-fit flex-col gap-3">
                    <Label>
                      <span>New workflow state</span>
                      <Select
                        bind:value={selectedState}
                        items={workflowOptions}
                        placeholder="Choose..."
                        class={selectClass}
                      ></Select>
                    </Label>
                    <Button
                      onclick={() => {
                        changeWorkflowState();
                      }}
                      disabled={!selectedState}
                      class="h-fit">Change</Button
                    >
                  </div>
                  <ErrorMessage error={changeWorkflowStateError}></ErrorMessage>
                </div>
              </Dropdown>
            {/if}
          </div>
        {/if}
        <div class="flex items-baseline gap-2">
          <Select
            size="md"
            id="pagecount"
            class="mt-2 h-8 w-24 !p-2 leading-3"
            items={[
              { name: "10", value: 10 },
              { name: "25", value: 25 },
              { name: "50", value: 50 },
              { name: "100", value: 100 }
            ]}
            value={limit}
            onchange={(event: any) => {
              setSearchParameters({
                currentPage: 1,
                limit: event.target.value
              });
            }}
          ></Select>
          <Label class="mr-3 text-nowrap"
            >{query
              ? "Matches per page"
              : tableType === SEARCHTYPES.ADVISORY
                ? "Advisories per page"
                : tableType === SEARCHTYPES.DOCUMENT
                  ? "Documents per page"
                  : "Events per page"}</Label
          >
        </div>
      </div>
      <div>
        <div class="mx-3 flex flex-row">
          <div class:invisible={currentPage === 1} class:flex={true} class:mr-3={true}>
            <PaginationItem onclick={first}>
              <i class="bx bx-arrow-to-left"></i>
            </PaginationItem>
            <PaginationItem onclick={previous}>
              <i class="bx bx-chevrons-left"></i>
            </PaginationItem>
          </div>
          <div class="flex items-center">
            <input
              class={`${numberOfPages < 10000 ? "w-16" : "w-20"} cursor-pointer border pr-1 text-right dark:bg-gray-800`}
              onchange={() => {
                let tmpCurrentPage = currentPage;
                if (!parseInt("" + tmpCurrentPage)) tmpCurrentPage = 1;
                tmpCurrentPage = Math.floor(tmpCurrentPage);
                if (tmpCurrentPage < 1) tmpCurrentPage = 1;
                if (tmpCurrentPage > numberOfPages) tmpCurrentPage = numberOfPages;
                setSearchParameters({
                  currentPage: tmpCurrentPage
                });
              }}
              value={currentPage}
            />
            <span class="mr-3 ml-2 text-nowrap">of {numberOfPages} pages</span>
          </div>
          <div class:invisible={currentPage === numberOfPages} class:flex={true}>
            <PaginationItem onclick={next}>
              <i class="bx bx-chevrons-right"></i>
            </PaginationItem>
            <PaginationItem onclick={last}>
              <i class="bx bx-arrow-to-right"></i>
            </PaginationItem>
          </div>
        </div>
      </div>
      <div class="mr-3 text-nowrap">
        {#if query}
          {count} matches found
        {:else if tableType === SEARCHTYPES.ADVISORY}
          {count} advisories in total
        {:else if tableType === SEARCHTYPES.DOCUMENT}
          {count} documents in total
        {:else}
          {count} events in total
        {/if}
      </div>
    {/if}
  </div>
  <div class:invisible={!loading} class:mb-4={true} class={loading ? "loadingFadeIn" : ""}>
    Loading ...
    <Spinner color="gray" size="4"></Spinner>
  </div>

  <ErrorMessage {error}></ErrorMessage>
  {#if documents?.length > 0}
    <div class="w-auto">
      <Table style="w-auto" hoverable={true} border={false}>
        <TableHead class="cursor-pointer dark:bg-gray-800">
          {#if isMultiSelectionAllowed}
            <TableHeadCell class="px-1">
              <CCheckbox
                checked={areAllSelected}
                onClicked={(event) => {
                  const isChecked = event.target.checked;
                  if (isChecked) {
                    for (let i = 0; i < documentIDs.length; i++) {
                      appStore.addSelectedDocumentID(documentIDs[i]);
                    }
                  } else {
                    appStore.clearSelectedDocumentIDs();
                  }
                }}
              ></CCheckbox>
            </TableHeadCell>
          {/if}
          <TableHeadCell class="px-0"></TableHeadCell>
          {#if areThereAnyComments}
            <TableHeadCell class={`${tablePadding} cursor-default`}>Comment</TableHeadCell>
          {/if}
          {#each columns as column}
            {#if column !== searchColumnName}
              <TableHeadCell
                class={tablePadding}
                onclick={() => {
                  switchSort(column);
                }}
                >{getColumnDisplayName(column)}<i
                  class:bx={true}
                  class:bx-caret-up={orderBy.find((c) => {
                    return c === column;
                  }) !== undefined}
                  class:bx-caret-down={orderBy.find((c) => {
                    return c === `-${column}`;
                  }) !== undefined}
                ></i>{getColumnOrder(orderBy, column)}</TableHeadCell
              >
            {/if}
          {/each}
        </TableHead>
        <TableBody>
          {#each documents as doc, i}
            {@const item = [SEARCHTYPES.ADVISORY, SEARCHTYPES.DOCUMENT].includes(tableType)
              ? {
                  id: doc.id,
                  ...doc.data[0]
                }
              : doc}
            <tr
              class={i % 2 == 1
                ? "bg-white hover:bg-gray-200 dark:bg-gray-800 dark:hover:bg-gray-600"
                : "bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600"}
            >
              {#if isMultiSelectionAllowed}
                <TableBodyCell class="px-1">
                  <CCheckbox
                    checked={appStore.state.app.selectedDocumentIDs.has(item.id)}
                    onClicked={(event) => {
                      const isChecked = event.target.checked;
                      if (isChecked) {
                        appStore.addSelectedDocumentID(item.id);
                      } else {
                        appStore.removeSelectedDocumentID(item.id);
                      }
                    }}
                  ></CCheckbox>
                </TableBodyCell>
              {/if}
              <TableBodyCell class="px-0">
                <div class="flex items-center">
                  {#if isAdmin && tableType !== SEARCHTYPES.EVENT}
                    <CIconButton
                      onClicked={() => {
                        appStore.setDocumentsToDelete([item]);
                        appStore.setIsDeleteModalOpen(true);
                      }}
                      title={`delete ${item.tracking_id}`}
                      icon="trash"
                      color="red"
                    ></CIconButton>
                  {/if}
                  <button
                    onclick={(e) => {
                      e.stopPropagation();
                      if (appStore.state.app.diff.docA_ID) {
                        appStore.setDiffDocB_ID(item.id);
                      } else {
                        appStore.setDiffDocA_ID(item.id);
                      }
                      appStore.openToolbox();
                      e.preventDefault();
                    }}
                    class:invisible={!appStore.state.app.isToolboxOpen &&
                      appStore.state.app.diff.docA_ID === undefined &&
                      appStore.state.app.diff.docB_ID === undefined}
                    disabled={appStore.state.app.diff.docA_ID === item.id.toString() ||
                      appStore.state.app.diff.docB_ID === item.id.toString() ||
                      disableDiffButtons}
                    class="min-w-[26px] p-1"
                    title={`Add to comparison: ${item.tracking_id}`}
                  >
                    <Img
                      src="plus-minus.svg"
                      class={`${
                        appStore.state.app.diff.docA_ID === item.id.toString() ||
                        appStore.state.app.diff.docB_ID === item.id.toString() ||
                        disableDiffButtons
                          ? "invert-[70%]"
                          : "dark:invert"
                      } min-h-4`}
                    />
                  </button>
                </div>
              </TableBodyCell>
              {#if areThereAnyComments}
                <TableBodyCell class={tdClassRelative}
                  ><a
                    class="absolute top-0 right-0 bottom-0 left-0"
                    href={getAdvisoryAnchorLink(item)}
                    aria-label="View advisory details"
                  >
                  </a>
                  <div class="m-2 table w-full text-wrap">
                    {#if item.comments_id}
                      {#await request(`api/comments/post/${item.comments_id}`, "GET")}
                        <Spinner color="gray" size="4"></Spinner>
                      {:then response}
                        {#if response.ok}
                          <div class="w-[120pt] max-w-[140pt] text-wrap">
                            {response.content.message}
                          </div>
                        {:else}
                          <span class="text-red-700">Couldn't load comment.</span>
                        {/if}
                      {/await}
                    {/if}
                  </div></TableBodyCell
                >
              {/if}
              {#each columns as column}
                {#if column !== searchColumnName}
                  {#if column === "cvss_v3_score" || column === "cvss_v2_score"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                        aria-label="View advisory details"
                      >
                      </a>
                      <CVSS baseScore={item[column]}></CVSS>
                    </TableBodyCell>
                  {:else if column === "ssvc"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                        aria-label="View advisory details"
                      >
                      </a>
                      <div class="m-2 table w-16 text-wrap">
                        {#if item[column]}
                          <SsvcBadge vector={item[column]}></SsvcBadge>
                        {/if}
                      </div></TableBodyCell
                    >
                  {:else if column === "state"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class="m-2 table w-full text-wrap">
                        <i
                          title={item[column]}
                          class:bx={true}
                          class:bxs-certification={item[column] === "new"}
                          class:bx-show={item[column] === "read"}
                          class:bxs-analyse={item[column] === "assessing"}
                          class:bx-book-open={item[column] === "review"}
                          class:bx-archive={item[column] === "archived"}
                          class:bx-trash={item[column] === "delete"}
                        ></i>
                      </div></TableBodyCell
                    >
                  {:else if column === "initial_release_date"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class="m-2 table w-full text-wrap">
                        {item.initial_release_date?.split("T")[0]}
                      </div></TableBodyCell
                    >
                  {:else if column === "current_release_date"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class="m-2 table w-full text-wrap">
                        {item.current_release_date?.split("T")[0]}
                      </div></TableBodyCell
                    >
                  {:else if column === "title"}
                    <TableBodyCell class={title + " relative"}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class="m-2 table w-[min(250px)] text-wrap">
                        <span title={item[column]}>{item[column]}</span>
                      </div></TableBodyCell
                    >
                  {:else if column === "publisher"}
                    <TableBodyCell class={publisher + " relative"}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class={publisher + " m-2"}>
                        <span title={item[column]}>{getPublisher(item[column], innerWidth)}</span>
                      </div></TableBodyCell
                    >
                  {:else if column === "recent"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class="m-2 table w-full text-wrap">
                        <span title={item[column]}
                          >{item[column] ? item[column].split("T")[0] : ""}</span
                        >
                      </div></TableBodyCell
                    >
                  {:else if column === "four_cves"}
                    <TableBodyCell class={tdClassRelative}>
                      {#if !(item[column] && item[column][0] && item[column].length > 1)}
                        <a
                          aria-label="View advisory details"
                          class="absolute top-0 right-0 bottom-0 left-0"
                          href={getAdvisoryAnchorLink(item)}
                        >
                        </a>
                      {/if}
                      <div class="w-32">
                        <div class="z-50 table p-2 text-wrap">
                          {#if item[column] && item[column][0]}
                            <!-- svelte-ignore a11y_click_events_have_key_events -->
                            <!-- svelte-ignore a11y_no_static_element_interactions -->
                            {#if item[column].length > 1}
                              <div
                                class="mr-2 flex cursor-pointer items-center"
                                onclick={(event) => {
                                  event.stopPropagation();
                                  toggleRow(i);
                                }}
                              >
                                <div class="flex-grow">
                                  {item[column][0]}
                                  {#if openRow === i}
                                    <div>
                                      {#each item.four_cves as cve, i}
                                        {#if i !== 0}
                                          <p>{cve}</p>
                                        {/if}
                                      {/each}
                                    </div>
                                  {/if}
                                </div>
                                <span>
                                  {#if openRow === i}
                                    <i class="bx bx-minus"></i>
                                  {:else}
                                    <i class="bx bx-plus"></i>
                                  {/if}
                                </span>
                              </div>
                            {:else}
                              <span>{item[column][0]}</span>
                            {/if}
                          {/if}
                        </div>
                      </div></TableBodyCell
                    >
                  {:else if column === "critical"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <CVSS baseScore={item[column]}></CVSS>
                    </TableBodyCell>
                  {:else if column === "tracking_id"}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class="m-2 table w-40 text-wrap">
                        {item[column] ?? ""}
                      </div></TableBodyCell
                    >
                  {:else}
                    <TableBodyCell class={tdClassRelative}
                      ><a
                        aria-label="View advisory details"
                        class="absolute top-0 right-0 bottom-0 left-0"
                        href={getAdvisoryAnchorLink(item)}
                      >
                      </a>
                      <div class="m-2 table w-full text-wrap">
                        {item[column] ?? ""}
                      </div></TableBodyCell
                    >
                  {/if}
                {/if}
              {/each}
            </tr>
            {#if [SEARCHTYPES.ADVISORY, SEARCHTYPES.DOCUMENT].includes(tableType)}
              {#each doc.data as result}
                {#if result[searchColumnName]}
                  <tr
                    class={i % 2 == 1
                      ? "border-t border-t-gray-200 bg-white dark:border-t-gray-700 dark:bg-gray-800"
                      : "border-t border-t-gray-300 bg-gray-100 dark:border-t-gray-600 dark:bg-gray-700"}
                  >
                    <TableBodyCell colspan={columns.length} class={tdClassRelative}>
                      {@html DOMPurify.sanitize(result[searchColumnName], {
                        USE_PROFILES: { html: true }
                      })}
                    </TableBodyCell>
                  </tr>
                {/if}
              {/each}
            {/if}
          {/each}
        </TableBody>
      </Table>
    </div>
  {:else if query}
    No results were found.
  {/if}
</div>
