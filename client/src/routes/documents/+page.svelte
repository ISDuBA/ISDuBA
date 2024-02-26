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
  import { tablePadding, tdClass } from "$lib/table/defaults";
  let documents: any = [];
  let searchTerm: string = "";
  const sortState: any = {
    id: "",
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
    publisher: defaultSortFunction("publisher"),
    title: defaultSortFunction("title"),
    trackingID: defaultSortFunction("tracking_id"),
    version: defaultSortFunction("version")
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
  $: filteredItems = documents;
  onMount(() => {
    if ($appStore.app.isUserLoggedIn) {
      $appStore.app.keycloak.updateToken(5).then(async () => {
        const response = await fetch("/api/documents", {
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
  <h1 class="text-lg">Documents</h1>
  {#if documents}
    <TableSearch placeholder="Search by maker name" hoverable={true} bind:inputValue={searchTerm}>
      <Table hoverable={true}>
        <TableHead class="cursor-pointer">
          <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("id")}
            >ID<i
              class:bx={true}
              class:bx-caret-up={sortState["activeSortColumn"] == "id" && sortState["id"] === "asc"}
              class:bx-caret-down={sortState["activeSortColumn"] == "id" &&
                sortState["id"] === "desc"}
            ></i></TableHeadCell
          >
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
          <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("version")}
            >Version<i
              class:bx={true}
              class:bx-caret-up={sortState["activeSortColumn"] == "version" &&
                sortState["version"] === "asc"}
              class:bx-caret-down={sortState["activeSortColumn"] == "version" &&
                sortState["version"] === "desc"}
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
              <TableBodyCell {tdClass}>{item.id}</TableBodyCell>
              <TableBodyCell {tdClass}>{item.publisher}</TableBodyCell>
              <TableBodyCell {tdClass}>{item.title}</TableBodyCell>
              <TableBodyCell {tdClass}>{item.tracking_id}</TableBodyCell>
              <TableBodyCell {tdClass}>{item.version}</TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    </TableSearch>
  {/if}
</RouteGuard>
