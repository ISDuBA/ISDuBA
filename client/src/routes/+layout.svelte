<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script>
  import "../app.pcss";
  import "boxicons/css/boxicons.min.css";
  import SideNav from "$lib/SideNav.svelte";
</script>

<div class="bg-primary-700 flex">
  <SideNav></SideNav>
  <main class="w-full bg-white pl-12 pt-12">
    <slot />
  </main>
</div>
