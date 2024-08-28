<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { TableBodyCell, Spinner } from "flowbite-svelte";
  import { tdClass } from "$lib/Table/defaults";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";

  let logs: any[] = [];
  let loadingLogs: boolean = false;
  let logError: ErrorDetails | null;
</script>

<CustomTable
  title="Logs"
  headers={[
    {
      label: "Time",
      attribute: "time"
    },
    {
      label: "level",
      attribute: "level"
    },
    {
      label: "Message",
      attribute: "msg"
    }
  ]}
>
  {#each logs as log, index (index)}
    <tr>
      <TableBodyCell {tdClass}>{log.time}</TableBodyCell>
      <TableBodyCell {tdClass}>{log.level}</TableBodyCell>
      <TableBodyCell {tdClass}>{log.msg}</TableBodyCell>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingLogs} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <ErrorMessage error={logError}></ErrorMessage>
  </div>
</CustomTable>
