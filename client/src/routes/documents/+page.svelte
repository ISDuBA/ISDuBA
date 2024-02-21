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
    publisher: "",
    title: "",
    trackingID: "",
    version: "",
    activeSortColumn: ""
  };
  const defaultSearchFunction = (attribute: string) => {
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
    id: defaultSearchFunction("id"),
    publisher: defaultSearchFunction("publisher"),
    title: defaultSearchFunction("title"),
    trackingID: defaultSearchFunction("tracking_id"),
    version: defaultSearchFunction("version")
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
  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
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
    }
  });
</script>

<RouteGuard>
  <h1 class="text-lg">Documents</h1>
  {#if documents}
    <TableSearch placeholder="Search by maker name" hoverable={true} bind:inputValue={searchTerm}>
      <Table hoverable={true}>
        <TableHead class="cursor-pointer">
          <TableHeadCell on:click={() => sortDocuments("id")}
            >ID<i
              class:bx={true}
              class:bx-caret-up={sortState["activeSortColumn"] == "id" && sortState["id"] === "asc"}
              class:bx-caret-down={sortState["activeSortColumn"] == "id" &&
                sortState["id"] === "desc"}
            ></i></TableHeadCell
          >
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
              <TableBodyCell>{item.publisher}</TableBodyCell>
              <TableBodyCell>{item.title}</TableBodyCell>
              <TableBodyCell>{item.tracking_id}</TableBodyCell>
              <TableBodyCell>{item.version}</TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    </TableSearch>
  {/if}
</RouteGuard>
