<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  /* eslint-disable svelte/no-at-html-tags */
  import { tick } from "svelte";
  import { push } from "svelte-spa-router";
  import {
    Button,
    Dropdown,
    Label,
    PaginationItem,
    Select,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Table,
    Img
  } from "flowbite-svelte";
  import { tdClass, tablePadding, title, publisher, searchColumnName } from "$lib/Table/defaults";
  import { Spinner } from "flowbite-svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { ADMIN, EDITOR, IMPORTER, REVIEWER } from "$lib/workflow";
  import { getAllowedWorkflowChanges, isRoleIncluded } from "$lib/permissions";
  import { appStore } from "$lib/store";
  import { getPublisher } from "$lib/publisher";
  import CIconButton from "$lib/Components/CIconButton.svelte";
  import SsvcBadge from "$lib/Advisories/SSVC/SSVCBadge.svelte";
  import { SEARCHTYPES } from "$lib/Queries/query";
  import CCheckbox from "$lib/Components/CCheckbox.svelte";
  import { areArraysEqual } from "$lib/utils";
  import DeleteModal from "./DeleteModal.svelte";
  import { updateMultipleStates } from "$lib/Advisories/advisory";

  let openRow: number | null;
  let abortController: AbortController;
  let requestOngoing = false;
  const toggleRow = (i: number) => {
    openRow = openRow === i ? null : i;
  };
  let limit = 10;
  let offset = 0;
  let count = 0;
  let currentPage = 1;
  let documents: any = null;
  $: documentIDs = documents?.map((d: any) => d.id) ?? [];
  let loading = false;
  let error: ErrorDetails | null;
  let changeWorkflowStateError: ErrorDetails | null;
  let prevQuery = "";
  export let columns: string[];
  export let query: string = "";
  export let searchTerm: string = "";
  export let tableType: SEARCHTYPES;
  export let orderBy: string[] = ["title"];
  export let defaultOrderBy = ["title"];

  $: disableDiffButtons =
    $appStore.app.diff.docA_ID !== undefined && $appStore.app.diff.docB_ID !== undefined;

  let anchorLink: string | null;

  $: areAllSelected =
    documents && areArraysEqual(documentIDs, Array.from($appStore.app.selectedDocumentIDs.keys()));

  $: selectedDocuments =
    $appStore.app.documents?.filter((d: any) => $appStore.app.selectedDocumentIDs.has(d.id)) ?? [];
  $: allowedWorkflowStateChanges = getAllowedWorkflowChanges(
    selectedDocuments?.map((d: any) => d.state) ?? []
  );
  $: workflowOptions = allowedWorkflowStateChanges.map((c) => {
    return { name: c.to, value: c.to };
  });
  $: isMultiSelectionAllowed =
    isRoleIncluded(appStore.getRoles(), [EDITOR, IMPORTER, ADMIN, REVIEWER]) &&
    ((tableType !== SEARCHTYPES.EVENT && appStore.isAdmin()) || tableType === SEARCHTYPES.ADVISORY);

  let selectedState: any;
  let dropdownOpen = false;
  const selectClass =
    "max-w-96 w-fit text-gray-900 disabled:text-gray-400 bg-gray-50 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:disabled:text-gray-500 dark:focus:ring-primary-500 dark:focus:border-primary-500";

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
      fetchData();
      dropdownOpen = false;
      selectedState = undefined;
    } else if (response.error) {
      changeWorkflowStateError = getErrorDetails("Couldn't change state.", response);
    }
  };

  let innerWidth = 0;

  const getColumnDisplayName = (column: string): string => {
    let names: { [key: string]: string } = {
      id: "ID",
      tracking_id: "TRACKING ID",
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

  const getTablePadding = (columns: any, match: string): any => {
    for (let i = 0; i < columns.length; i++) {
      if (columns[i] === match) {
        return [Array(i).fill(0), Array(columns.length - i - 1).fill(0)];
      }
    }
    return [[], Array(columns.length).fill(0)];
  };

  let searchPadding: any[] = [];
  let searchPaddingRight: any[] = [];

  $: if (columns !== undefined) {
    [searchPadding, searchPaddingRight] = getTablePadding(columns, "title");
  }

  const savePosition = () => {
    let position = [offset, currentPage, limit, orderBy];
    sessionStorage.setItem("tablePosition" + query + tableType, JSON.stringify(position));
  };

  let postitionRestored: boolean = false;
  const restorePosition = () => {
    let position = sessionStorage.getItem("tablePosition" + query + tableType);
    if (position) {
      [offset, currentPage, limit, orderBy] = JSON.parse(position);
    } else {
      offset = 0;
      currentPage = 1;
    }
  };

  const setOrderBy = async () => {
    await tick();
    orderBy
      .map((c) => {
        return c.replace("-", "");
      })
      .forEach((c) => {
        if (!orderBy.includes(c)) orderBy = defaultOrderBy;
      });
  };

  $: if (columns) {
    setOrderBy();
  }

  $: if (offset || currentPage || limit || orderBy) {
    if (!postitionRestored) {
      restorePosition();
      postitionRestored = true;
    }
    savePosition();
  }

  $: if (tableType || !tableType) {
    restorePosition();
    savePosition();
  }

  $: isAdmin = isRoleIncluded(appStore.getRoles(), [ADMIN]);

  export async function fetchData(): Promise<void> {
    appStore.setDocuments([]);
    appStore.clearSelectedDocumentIDs();
    openRow = null;
    if (query !== prevQuery) {
      restorePosition();
      savePosition();
      prevQuery = query;
    }
    const searchSuffix = searchTerm ? `"${searchTerm}" search ${searchColumnName} as ` : "";
    const searchColumn = searchTerm ? ` ${searchColumnName}` : "";
    let queryParam = "";
    if (query || searchSuffix) {
      queryParam = `query=${query}${searchSuffix}`;
    }
    let fetchColumns = [...columns];
    let requiredColumns = ["id", "tracking_id", "publisher"];
    for (let c of requiredColumns) {
      if (!fetchColumns.includes(c)) {
        fetchColumns.push(c);
      }
    }
    let documentURL = "";

    if (tableType === SEARCHTYPES.EVENT) {
      documentURL = encodeURI(
        `/api/events?${queryParam}&count=1&orders=${orderBy.join(" ")}&limit=${limit}&offset=${offset}&columns=${fetchColumns.join(" ")}${searchColumn}`
      );
    } else {
      const loadAdvisories = tableType === SEARCHTYPES.ADVISORY;
      documentURL = encodeURI(
        `/api/documents?${queryParam}&advisories=${loadAdvisories}&count=1&orders=${orderBy.join(" ")}&limit=${limit}&offset=${offset}&columns=${fetchColumns.join(" ")}${searchColumn}`
      );
    }

    error = null;
    loading = true;
    if (!requestOngoing) {
      requestOngoing = true;
      abortController = new AbortController();
    } else {
      abortController.abort();
    }
    const response = await request(documentURL, "GET");
    if (response.ok) {
      ({ count, documents } = response.content);
      if (tableType === SEARCHTYPES.EVENT) {
        count = response.content.count;
        documents = response.content.events;
      } else {
        ({ count, documents } = response.content);
      }
      appStore.setDocuments(documents);
    } else if (response.error) {
      error =
        response.error === "400"
          ? getErrorDetails(`Please check your search syntax.`, response)
          : response.content.includes("deadline exceeded")
            ? getErrorDetails(`The server wasn't able to answer your request in time.`)
            : getErrorDetails(`Could not load query.`, response);
    }
    loading = false;
    requestOngoing = false;
  }

  const previous = () => {
    if (offset - limit >= 0) {
      offset = offset - limit > 0 ? offset - limit : 0;
      currentPage -= 1;
    }
    fetchData();
  };
  const next = () => {
    if (offset + limit <= count) {
      offset = offset + limit;
      currentPage += 1;
    }
    fetchData();
  };

  const first = () => {
    offset = 0;
    currentPage = 1;
    fetchData();
  };

  const last = () => {
    offset = (numberOfPages - 1) * limit;
    currentPage = numberOfPages;
    fetchData();
  };

  const switchSort = async (column: string) => {
    let found = orderBy.find((c) => c === column);
    let foundMinus = orderBy.find((c) => c === "-" + column);
    if (foundMinus) {
      orderBy = orderBy.filter((c) => c !== "-" + column);
    }
    if (found) {
      orderBy = orderBy.map((c) => (c === column ? `-${column}` : c));
    }
    if (!found && !foundMinus) {
      orderBy.push(column);
    }
    orderBy = orderBy;
    await tick();
    fetchData();
  };

  const onDeleted = async () => {
    await fetchData();
    first();
  };

  $: numberOfPages = Math.ceil(count / limit);

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

<DeleteModal
  on:deleted={onDeleted}
  documents={$appStore.app.documentsToDelete || []}
  type={tableType}
></DeleteModal>

<div class="flex-grow">
  <div class="mb-2 mt-2 flex flex-row items-baseline justify-between">
    {#if documents?.length > 0}
      <div class="flex flex-row items-baseline gap-8">
        {#if isMultiSelectionAllowed}
          <div class="flex items-center gap-2">
            {#if appStore.isAdmin()}
              <Button
                on:click={() => {
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
                <i class="bx bx-git-commit text-black-700"></i>
              </Button>
              <Dropdown
                bind:open={dropdownOpen}
                on:show={(event) => {
                  if (!event.detail) {
                    changeWorkflowStateError = null;
                  }
                }}
                placement="top-start"
                triggeredBy="#state-icon"
                class="w-full max-w-sm divide-y divide-gray-100 rounded p-4 shadow dark:divide-gray-700 dark:bg-gray-800"
                containerClass="divide-y z-50 border border-gray-300"
              >
                <div class="flex flex-col gap-3">
                  <div class="flex w-fit flex-col gap-3">
                    <Label>
                      <span>New workflow state</span>
                      <Select
                        bind:value={selectedState}
                        items={workflowOptions}
                        placeholder="Choose..."
                        defaultClass={selectClass}
                      ></Select>
                    </Label>
                    <Button
                      on:click={() => {
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
            bind:value={limit}
            on:change={() => {
              offset = 0;
              currentPage = 1;
              fetchData();
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
            <PaginationItem on:click={first}>
              <i class="bx bx-arrow-to-left"></i>
            </PaginationItem>
            <PaginationItem on:click={previous}>
              <i class="bx bx-chevrons-left"></i>
            </PaginationItem>
          </div>
          <div class="flex items-center">
            <input
              class={`${numberOfPages < 10000 ? "w-16" : "w-20"} cursor-pointer border pr-1 text-right`}
              on:change={() => {
                if (!parseInt("" + currentPage)) currentPage = 1;
                currentPage = Math.floor(currentPage);
                if (currentPage < 1) currentPage = 1;
                if (currentPage > numberOfPages) currentPage = numberOfPages;
                offset = (currentPage - 1) * limit;
                fetchData();
              }}
              bind:value={currentPage}
            />
            <span class="ml-2 mr-3 text-nowrap">of {numberOfPages} pages</span>
          </div>
          <div class:invisible={currentPage === numberOfPages} class:flex={true}>
            <PaginationItem on:click={next}>
              <i class="bx bx-chevrons-right"></i>
            </PaginationItem>
            <PaginationItem on:click={last}>
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
      <a href={anchorLink}>
        <Table style="w-auto" hoverable={true} noborder={true}>
          <TableHead class="cursor-pointer">
            {#if isMultiSelectionAllowed}
              <TableHeadCell padding="px-0">
                <CCheckbox
                  checked={areAllSelected}
                  on:click={(event) => {
                    const isChecked = event.detail.target.checked;
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
            <TableHeadCell padding="px-0"></TableHeadCell>
            {#each columns as column}
              {#if column !== searchColumnName}
                <TableHeadCell
                  padding={tablePadding}
                  on:click={() => {
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
            {#each documents as item, i}
              <tr
                class="cursor-pointer odd:bg-white even:bg-gray-100 hover:bg-gray-200"
                on:click={() => {
                  push(`/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`);
                }}
                on:mouseenter={() => {
                  anchorLink = `#/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`;
                }}
                on:mouseleave={() => {
                  anchorLink = null;
                }}
              >
                {#if isMultiSelectionAllowed}
                  <TableBodyCell tdClass="px-0">
                    <CCheckbox
                      checked={$appStore.app.selectedDocumentIDs.has(item.id)}
                      on:click={(event) => {
                        const isChecked = event.detail.target.checked;
                        if (isChecked) {
                          appStore.addSelectedDocumentID(item.id);
                        } else {
                          appStore.removeSelectedDocumentID(item.id);
                        }
                      }}
                    ></CCheckbox>
                  </TableBodyCell>
                {/if}
                <TableBodyCell tdClass="px-0">
                  <div class="flex items-center">
                    {#if isAdmin && tableType !== SEARCHTYPES.EVENT}
                      <CIconButton
                        on:click={() => {
                          appStore.setDocumentsToDelete([item]);
                          appStore.setIsDeleteModalOpen(true);
                        }}
                        title={`delete ${item.tracking_id}`}
                        icon="trash"
                        color="red"
                      ></CIconButton>
                    {/if}
                    <button
                      on:click|stopPropagation={(e) => {
                        if ($appStore.app.diff.docA_ID) {
                          appStore.setDiffDocB_ID(item.id);
                        } else {
                          appStore.setDiffDocA_ID(item.id);
                        }
                        appStore.openToolbox();
                        e.preventDefault();
                      }}
                      class:invisible={!$appStore.app.isToolboxOpen &&
                        $appStore.app.diff.docA_ID === undefined &&
                        $appStore.app.diff.docB_ID === undefined}
                      disabled={$appStore.app.diff.docA_ID === item.id.toString() ||
                        $appStore.app.diff.docB_ID === item.id.toString() ||
                        disableDiffButtons}
                      class="min-w-[26px] p-1"
                      title={`compare ${item.tracking_id}`}
                    >
                      <Img
                        src="plus-minus.svg"
                        class={`${
                          $appStore.app.diff.docA_ID === item.id.toString() ||
                          $appStore.app.diff.docB_ID === item.id.toString() ||
                          disableDiffButtons
                            ? "invert-[70%]"
                            : ""
                        } min-h-4`}
                      />
                    </button>
                  </div>
                </TableBodyCell>
                {#each columns as column}
                  {#if column !== searchColumnName}
                    {#if column === "cvss_v3_score" || column === "cvss_v2_score"}
                      <TableBodyCell {tdClass}
                        ><span class:text-red-500={Number(item[column]) > 5.0}
                          >{item[column] == null ? "" : item[column]}</span
                        ></TableBodyCell
                      >
                    {:else if column === "ssvc"}
                      <TableBodyCell {tdClass}>
                        {#if item[column]}
                          <SsvcBadge vector={item[column]}></SsvcBadge>
                        {/if}
                      </TableBodyCell>
                    {:else if column === "state"}
                      <TableBodyCell {tdClass}
                        ><i
                          title={item[column]}
                          class:bx={true}
                          class:bxs-star={item[column] === "new"}
                          class:bx-show={item[column] === "read"}
                          class:bxs-analyse={item[column] === "assessing"}
                          class:bx-book-open={item[column] === "review"}
                          class:bx-archive={item[column] === "archived"}
                          class:bx-trash={item[column] === "delete"}
                        ></i>
                      </TableBodyCell>
                    {:else if column === "initial_release_date"}
                      <TableBodyCell {tdClass}
                        >{item.initial_release_date?.split("T")[0]}</TableBodyCell
                      >
                    {:else if column === "current_release_date"}
                      <TableBodyCell {tdClass}
                        >{item.current_release_date?.split("T")[0]}</TableBodyCell
                      >
                    {:else if column === "title"}
                      <TableBodyCell tdClass={title}
                        ><span title={item[column]}>{item[column]}</span></TableBodyCell
                      >
                    {:else if column === "publisher"}
                      <TableBodyCell tdClass={publisher}
                        ><span title={item[column]}>{getPublisher(item[column], innerWidth)}</span
                        ></TableBodyCell
                      >
                    {:else if column === "recent"}
                      <TableBodyCell {tdClass}
                        ><span title={item[column]}
                          >{item[column] ? item[column].split("T")[0] : ""}</span
                        ></TableBodyCell
                      >
                    {:else if column === "four_cves"}
                      <TableBodyCell {tdClass}
                        >{#if item[column] && item[column][0]}
                          <!-- svelte-ignore a11y-click-events-have-key-events -->
                          <!-- svelte-ignore a11y-no-static-element-interactions -->
                          {#if item[column].length > 1}
                            <div
                              class="mr-2 flex items-center"
                              on:mouseenter={() => (anchorLink = null)}
                              on:click|stopPropagation={() => toggleRow(i)}
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
                        {/if}</TableBodyCell
                      >
                    {:else if column === "critical"}
                      <TableBodyCell {tdClass}
                        ><span class:text-red-500={Number(item[column]) > 5.0}
                          >{item[column] == null ? "" : item[column]}</span
                        ></TableBodyCell
                      >
                    {:else}
                      <TableBodyCell {tdClass}>{item[column] ?? ""}</TableBodyCell>
                    {/if}
                  {/if}
                {/each}
              </tr>
              {#if item[searchColumnName]}
                <TableBodyRow class="border border-y-indigo-500/100 bg-white">
                  <!-- eslint-disable-next-line  @typescript-eslint/no-unused-vars -->
                  {#each searchPadding as _}
                    <TableBodyCell {tdClass}></TableBodyCell>
                  {/each}
                  <TableBodyCell {tdClass}>{@html item[searchColumnName]}</TableBodyCell>
                  <!-- eslint-disable-next-line  @typescript-eslint/no-unused-vars -->
                  {#each searchPaddingRight as _}
                    <TableBodyCell {tdClass}></TableBodyCell>
                  {/each}
                </TableBodyRow>
              {/if}
            {/each}
          </TableBody>
        </Table>
      </a>
    </div>
  {:else if query}
    No results were found.
  {/if}
</div>
