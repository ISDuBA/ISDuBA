<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { routerState } from "$routes/router.svelte";
  import type { Snippet } from "svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";

  type Props = {
    onclick?: (() => void) | ((e: Event) => void);
    ariaLabel?: string;
    href: string;
    id?: string;
    internal?: boolean;
    children?: Snippet;
  } & HTMLButtonAttributes;

  let { ariaLabel, href, id, internal = true, onclick, children, ...restProps }: Props = $props();

  const onClicked = (event: Event) => {
    if (internal) {
      routerState.didPush = true;
    }
    if (onclick) onclick(event);
  };
</script>

<a onclick={onClicked} class={restProps.class ?? ""} aria-label={ariaLabel} {href} {id}>
  {@render children?.()}
</a>
