<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { twMerge } from "tailwind-merge";
  import { setContext } from "svelte";
  import type { Snippet } from "svelte";
  import { writable, type Writable } from "svelte/store";
  import type { HTMLButtonAttributes } from "svelte/elements";

  interface AccordionCtxType {
    flush: boolean | undefined;
    activeClass: string;
    inactiveClass: string;
    selected?: Writable<object>;
    classActive?: string;
    classInactive?: string;
  }

  type Props = {
    multiple?: boolean;
    flush?: boolean;
    activeClass?: string;
    inactiveClass?: string;
    defaultClass?: string;
    classActive?: string;
    classInactive?: string;
    children?: Snippet;
  } & HTMLButtonAttributes;

  let {
    multiple = false,
    flush = false,
    activeClass = "bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-800",
    inactiveClass = "text-gray-500 dark:text-gray-400 hover:bg-gray-100 hover:dark:bg-gray-800",
    defaultClass = "text-gray-500 dark:text-gray-400",
    classActive = "",
    classInactive = "",
    children,
    ...restProps
  }: Props = $props();

  const ctx: AccordionCtxType = {
    flush,
    activeClass: twMerge(activeClass, classActive),
    inactiveClass: twMerge(inactiveClass, classInactive),
    selected: multiple ? undefined : writable()
  };

  setContext<AccordionCtxType>("ctx", ctx);

  let frameClass: string = $derived(twMerge(defaultClass, `${restProps.class}`));
</script>

<div class={frameClass} color="none">
  {@render children?.()}
</div>
