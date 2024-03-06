<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { push } from "svelte-spa-router";
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    TableSearch
  } from "flowbite-svelte";
  import { tdClass, tablePadding } from "$lib/table/defaults";

  let openRow: number | null;

  const toggleRow = (i: number) => {
    openRow = openRow === i ? null : i;
  };
  let documents: any = [];
  let searchTerm: string = "";
  const sortState: any = {
    id: "",
    cvss: "",
    publisher: "",
    title: "",
    trackingID: "",
    version: "",
    activeSortColumn: ""
  };
  const defaultSortFunction = (attribute: string) => {
    return {
      asc: (ad1: any, ad2: any) => {
        if (ad1[attribute] < ad2[attribute]) return -1;
        if (ad1[attribute] > ad2[attribute]) return 1;
        return 0;
      },
      desc: (ad2: any, ad1: any) => {
        if (ad1[attribute] < ad2[attribute]) return -1;
        if (ad1[attribute] > ad2[attribute]) return 1;
        return 0;
      }
    };
  };
  const sortFunctionsByColumn: any = {
    id: defaultSortFunction("id"),
    cvss: defaultSortFunction("cvss_v3_score"),
    publisher: defaultSortFunction("publisher"),
    title: defaultSortFunction("title"),
    trackingID: defaultSortFunction("tracking_id"),
    version: defaultSortFunction("version"),
    state: defaultSortFunction("state")
  };
  const sortDocuments = (column: string) => {
    sortState["activeSortColumn"] = column;
    if (sortState[column] === "asc") {
      documents = [...documents.sort(sortFunctionsByColumn[column]["desc"])];
      sortState[column] = "desc";
    } else {
      documents = [...documents.sort(sortFunctionsByColumn[column]["asc"])];
      sortState[column] = "asc";
    }
  };
  const allUnreadDocuments = encodeURI(
    "/api/documents?columns=id tracking_id version publisher current_release_date initial_release_date title tlp cvss_v2_score cvss_v3_score four_cves"
  );
  $: filteredItems = documents;
  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      $appStore.app.keycloak.updateToken(5).then(async () => {
        const response = await fetch(allUnreadDocuments, {
          headers: {
            Authorization: `Bearer ${$appStore.app.keycloak.token}`
          }
        });
        if (response.ok) {
          ({ documents } = await response.json());
          sortDocuments("id");
        } else {
          // Do errorhandling
        }
      });
    }
  });
</script>

<h1 class="mb-3 text-lg">Advisories</h1>
{#if documents}
  <TableSearch
    placeholder="Search by maker name"
    hoverable={true}
    bind:inputValue={searchTerm}
    noborder={true}
  >
    <TableHead class="cursor-pointer">
      <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("id")}
        >ID<i
          class:bx={true}
          class:bx-caret-up={sortState["activeSortColumn"] == "id" && sortState["id"] === "asc"}
          class:bx-caret-down={sortState["activeSortColumn"] == "id" && sortState["id"] === "desc"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("cvss")}
        >CVSS<i
          class:bx={true}
          class:bx-caret-up={sortState["activeSortColumn"] == "cvss" && sortState["cvss"] === "asc"}
          class:bx-caret-down={sortState["activeSortColumn"] == "cvss" &&
            sortState["id"] === "desc"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding}>CVEs</TableHeadCell>
      <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("publisher")}
        >Publisher<i
          class:bx={true}
          class:bx-caret-up={sortState["activeSortColumn"] == "publisher" &&
            sortState["publisher"] === "asc"}
          class:bx-caret-down={sortState["activeSortColumn"] == "publisher" &&
            sortState["publisher"] === "desc"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("title")}
        >Title<i
          class:bx={true}
          class:bx-caret-up={sortState["activeSortColumn"] == "title" &&
            sortState["title"] === "asc"}
          class:bx-caret-down={sortState["activeSortColumn"] == "title" &&
            sortState["title"] === "desc"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("trackingID")}
        >Tracking ID<i
          class:bx={true}
          class:bx-caret-up={sortState["activeSortColumn"] == "trackingID" &&
            sortState["trackingID"] === "asc"}
          class:bx-caret-down={sortState["activeSortColumn"] == "trackingID" &&
            sortState["trackingID"] === "desc"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding}>Initial Release</TableHeadCell>
      <TableHeadCell padding={tablePadding}>Current Release</TableHeadCell>
      <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("version")}
        >Version<i
          class:bx={true}
          class:bx-caret-up={sortState["activeSortColumn"] == "version" &&
            sortState["version"] === "asc"}
          class:bx-caret-down={sortState["activeSortColumn"] == "version" &&
            sortState["version"] === "desc"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("state")}
        >State<i
          class:bx={true}
          class:bx-caret-up={sortState["activeSortColumn"] == "state" &&
            sortState["state"] === "asc"}
          class:bx-caret-down={sortState["activeSortColumn"] == "state" &&
            sortState["state"] === "desc"}
        ></i></TableHeadCell
      >
    </TableHead>
    <TableBody>
      {#each filteredItems as item, i}
        <TableBodyRow
          class="cursor-pointer"
          on:click={(event) => {
            push(`/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`);
          }}
        >
          <TableBodyCell {tdClass}>{item.id}</TableBodyCell>
          <TableBodyCell {tdClass}
            ><span class:text-red-500={Number(item.cvss_v3_score) > 5.0}>{item.cvss_v3_score}</span
            ></TableBodyCell
          >
          <TableBodyCell {tdClass}
            >{#if item.four_cves[0]}
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <!-- svelte-ignore a11y-no-static-element-interactions -->
              {#if item.four_cves.length > 1}
                <span on:click|stopPropagation={() => toggleRow(i)}>
                  {item.four_cves[0]}

                  {#if openRow === i}
                    <i class="bx bx-minus"></i>
                  {:else}
                    <i class="bx bx-plus"></i>
                  {/if}
                </span>
              {:else}
                <span>{item.four_cves[0]}</span>
              {/if}
            {/if}</TableBodyCell
          >
          <TableBodyCell {tdClass}>{item.publisher}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.title}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.tracking_id}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.current_release_date.split("T")[0]}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.initial_release_date.split("T")[0]}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.version}</TableBodyCell>
          <TableBodyCell {tdClass}
            ><i class:bx={true} class:bxs-star={item.state === "new"}></i></TableBodyCell
          >
        </TableBodyRow>
        {#if openRow === i}
          <TableBodyRow>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}>
              <div>
                {#each item.four_cves as cve, i}
                  {#if i !== 0}
                    <div>{cve}</div>
                  {/if}
                {/each}
              </div>
            </TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
          </TableBodyRow>
        {/if}
      {/each}
    </TableBody>
  </TableSearch>
{/if}
