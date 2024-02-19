<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
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
      } else {
        // Do errorhandling
      }
    }
  });
</script>

<RouteGuard>
  <h1 class="text-lg">Advisories</h1>
  {#if documents}
    <TableSearch placeholder="Search by maker name" hoverable={true} bind:inputValue={searchTerm}>
      <TableHead>
        <TableHeadCell>ID</TableHeadCell>
        <TableHeadCell>Publisher</TableHeadCell>
        <TableHeadCell>Title</TableHeadCell>
        <TableHeadCell>Tracking ID</TableHeadCell>
        <TableHeadCell>Version</TableHeadCell>
      </TableHead>
      <TableBody>
        {#each filteredItems as item}
          <TableBodyRow>
            <TableBodyCell>{item.id}</TableBodyCell>
            <TableBodyCell>{item.publisher}</TableBodyCell>
            <TableBodyCell>{item.title}</TableBodyCell>
            <TableBodyCell>{item.tracking_id}</TableBodyCell>
            <TableBodyCell>{item.version}</TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </TableSearch>
  {/if}
</RouteGuard>
