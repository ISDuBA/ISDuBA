<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { A, P } from "flowbite-svelte";
  /*global __APP_VERSION__*/
  const version: string = __APP_VERSION__;
</script>

<h1 class="mb-3 text-lg">About ISDuBA</h1>
Versions:
<ul>
  <li>Client: {version}</li>
  <li>Backend: blubb</li>
</ul>
<A href="https://github.com/ISDuBA/" class="underline hover:no-underline"
  >Visit the ISDuBA project on Github</A
>
