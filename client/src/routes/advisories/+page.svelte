<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { goto } from "$app/navigation";
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    TableSearch
  } from "flowbite-svelte";
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
    "/api/documents?columns=id state tracking_id version publisher current_release_date initial_release_date title tlp cvss_v2_score cvss_v3_score four_cves"
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

<RouteGuard>
  <h1 class="mb-3 mt-10 text-lg">Advisories</h1>
  {#if documents}
    <div style="overflow-y:auto; height:55%">
      <TableSearch placeholder="Search by maker name" hoverable={true} bind:inputValue={searchTerm}>
        <Table hoverable={true}>
          <TableHead class="cursor-pointer">
            <TableHeadCell on:click={() => sortDocuments("id")}
              >ID<i
                class:bx={true}
                class:bx-caret-up={sortState["activeSortColumn"] == "id" &&
                  sortState["id"] === "asc"}
                class:bx-caret-down={sortState["activeSortColumn"] == "id" &&
                  sortState["id"] === "desc"}
              ></i></TableHeadCell
            >
            <TableHeadCell on:click={() => sortDocuments("cvss")}
              >CVSS<i
                class:bx={true}
                class:bx-caret-up={sortState["activeSortColumn"] == "cvss" &&
                  sortState["cvss"] === "asc"}
                class:bx-caret-down={sortState["activeSortColumn"] == "cvss" &&
                  sortState["id"] === "desc"}
              ></i></TableHeadCell
            >
            <TableHeadCell>CVEs</TableHeadCell>
            <TableHeadCell on:click={() => sortDocuments("publisher")}
              >Publisher<i
                class:bx={true}
                class:bx-caret-up={sortState["activeSortColumn"] == "publisher" &&
                  sortState["publisher"] === "asc"}
                class:bx-caret-down={sortState["activeSortColumn"] == "publisher" &&
                  sortState["publisher"] === "desc"}
              ></i></TableHeadCell
            >
            <TableHeadCell on:click={() => sortDocuments("title")}
              >Title<i
                class:bx={true}
                class:bx-caret-up={sortState["activeSortColumn"] == "title" &&
                  sortState["title"] === "asc"}
                class:bx-caret-down={sortState["activeSortColumn"] == "title" &&
                  sortState["title"] === "desc"}
              ></i></TableHeadCell
            >
            <TableHeadCell on:click={() => sortDocuments("trackingID")}
              >Tracking ID<i
                class:bx={true}
                class:bx-caret-up={sortState["activeSortColumn"] == "trackingID" &&
                  sortState["trackingID"] === "asc"}
                class:bx-caret-down={sortState["activeSortColumn"] == "trackingID" &&
                  sortState["trackingID"] === "desc"}
              ></i></TableHeadCell
            >
            <TableHeadCell on:click={() => sortDocuments("version")}
              >Version<i
                class:bx={true}
                class:bx-caret-up={sortState["activeSortColumn"] == "version" &&
                  sortState["version"] === "asc"}
                class:bx-caret-down={sortState["activeSortColumn"] == "version" &&
                  sortState["version"] === "desc"}
              ></i></TableHeadCell
            >
            <TableHeadCell on:click={() => sortDocuments("state")}
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
            {#each filteredItems as item}
              <TableBodyRow
                class="cursor-pointer"
                on:click={() => {
                  goto(`/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`);
                }}
              >
                <TableBodyCell>{item.id}</TableBodyCell>
                <TableBodyCell
                  ><span class:text-red-500={Number(item.cvss_v3_score) > 5.0}
                    >{item.cvss_v3_score}</span
                  ></TableBodyCell
                >
                <TableBodyCell>{item.four_cves.replace(/\[|\]|\"/g, "")}</TableBodyCell>
                <TableBodyCell>{item.publisher}</TableBodyCell>
                <TableBodyCell>{item.title}</TableBodyCell>
                <TableBodyCell>{item.tracking_id}</TableBodyCell>
                <TableBodyCell>{item.version}</TableBodyCell>
                <TableBodyCell>{item.state}</TableBodyCell>
              </TableBodyRow>
            {/each}
          </TableBody>
        </Table>
      </TableSearch>
    </div>
  {/if}
</RouteGuard>
