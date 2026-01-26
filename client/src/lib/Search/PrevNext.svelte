<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { SidebarItem } from "flowbite-svelte";
  import { sidebarItemClass, sidebarItemLinkClass } from "$lib/sidenav";

  const indentedSidebarItemClass = `${sidebarItemLinkClass} !py-1 ps-10`;
</script>

<SidebarItem
  class={sidebarItemClass}
  aClass={indentedSidebarItemClass}
  label="Test"
  href="/#/sources"
>
  {#snippet icon()}{/snippet}
</SidebarItem>
<SidebarItem
  class={sidebarItemClass}
  aClass={indentedSidebarItemClass}
  label="Test 2"
  href="/#/sources"
>
  {#snippet icon()}{/snippet}
</SidebarItem>
