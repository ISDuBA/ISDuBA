<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { A, P, Li, List } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  const clientVersion: string = __APP_VERSION__;
  const backendVersion: string = __BACKEND_VERSION__;
</script>

<SectionHeader title="About ISDuBA"></SectionHeader>

<P>
  <A href="https://github.com/ISDuBA/" class="underline hover:no-underline"
    >Visit the ISDuBA project on Github</A
  ></P
>
<P class="mt-3">
  Versions:
  <List tag="ul" class="space-y-1" list="none">
    <Li liClass="ml-3">Client: {clientVersion}</Li>
    <Li liClass="ml-3">Backend: {backendVersion}</Li>
  </List>
</P>
