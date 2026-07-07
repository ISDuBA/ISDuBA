<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { activeClass, nonActiveClass, sidebarItemLinkClass } from "$lib/sidenav";
  import type { Snippet } from "svelte";
  import { getContext } from "svelte";
  import Link from "./Link.svelte";

  type Props = {
    active?: boolean;
    aClass?: string;
    href: string;
    label: string;
    icon: Snippet;
  };

  let { active, aClass = sidebarItemLinkClass, label, href, icon }: Props = $props();

  const activeURL: () => string = getContext("activeURL");

  const isActive = $derived(active !== undefined ? active : href === activeURL());
</script>

<li>
  <Link {href} class={`${aClass} ${isActive ? activeClass : nonActiveClass}`}>
    {@render icon()}
    <span>{label}</span>
  </Link>
</li>
