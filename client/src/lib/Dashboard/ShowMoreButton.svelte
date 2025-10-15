<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button } from "flowbite-svelte";
  import { push } from "svelte-spa-router";

  interface Props {
    id: string;
  }

  let { id }: Props = $props();

  const showMore = () => {
    push(`/search?query=${id}`);
  };
</script>

<Button onclick={showMore} color="light" class="h-fit w-fit rounded-md !px-2 !py-1">
  <i class="bx bx-spreadsheet text-lg"></i>
</Button>
