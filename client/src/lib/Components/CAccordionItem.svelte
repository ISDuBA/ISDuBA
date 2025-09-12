<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { twMerge } from "tailwind-merge";
  import { getContext, onMount } from "svelte";
  import { writable, type Writable } from "svelte/store";
  import { fade, blur, fly, slide } from "svelte/transition";
  import type { HTMLButtonAttributes } from "svelte/elements";
  import type { Snippet } from "svelte";

  type TransitionTypes =
    | "fade"
    | "fly"
    | "slide"
    | "blur"
    | "in:fly"
    | "out:fly"
    | "in:slide"
    | "out:slide"
    | "in:fade"
    | "out:fade"
    | "in:blur"
    | "out:blur";

  interface TransitionParamTypes {
    delay?: number;
    duration?: number;
    easing?: (t: number) => number;
    css?: (t: number, u: number) => string;
    tick?: (t: number, u: number) => void;
  }
  interface AccordionCtxType {
    flush: boolean | undefined;
    activeClass: string;
    inactiveClass: string;
    selected?: Writable<object>;
    classActive?: string;
    classInactive?: string;
  }

  type Props = {
    id?: string;
    tag?: string;
    open?: boolean;
    activeClass?: string;
    inactiveClass?: string;
    defaultClass?: string;
    transitionType?: TransitionTypes;
    transitionParams?: TransitionParamTypes;
    paddingFlush?: string;
    paddingDefault?: string;
    textFlushOpen?: string;
    textFlushDefault?: string;
    borderClass?: string;
    borderOpenClass?: string;
    borderBottomClass?: string;
    borderSharedClass?: string;
    classActive?: string;
    classInactive?: string;
    toggleCallback?: () => Promise<any>;
    headerSlot: Snippet;
    arrowup?: Snippet;
    arrowdown?: Snippet;
    children?: Snippet;
  } & HTMLButtonAttributes;

  let {
    id = undefined,
    tag = "h2",
    open = $bindable(false),
    activeClass = undefined,
    inactiveClass = undefined,
    defaultClass = "flex items-center justify-between w-full font-medium text-left group-first:rounded-t-xl border-gray-200 dark:border-gray-700",
    transitionType = "slide",
    transitionParams = {},
    paddingFlush = "py-5",
    paddingDefault = "p-5",
    textFlushOpen = "text-gray-900 dark:text-white",
    textFlushDefault = "text-gray-500 dark:text-gray-400",
    borderClass = "border-s border-e group-first:border-t",
    borderOpenClass = "border-s border-e",
    borderBottomClass = "border-b",
    borderSharedClass = "border-gray-200 dark:border-gray-700",
    classActive = undefined,
    classInactive = undefined,
    toggleCallback = undefined,
    headerSlot,
    arrowup = undefined,
    arrowdown = undefined,
    children,
    ...restProps
  }: Props = $props();

  let activeCls = twMerge(activeClass, classActive);
  let inactiveCls = twMerge(inactiveClass, classInactive);

  // make a custom transition function that returns the desired transition
  const multiple = (node: HTMLElement, params: any) => {
    switch (transitionType) {
      case "blur":
        return blur(node, params);
      case "fly":
        return fly(node, params);
      case "fade":
        return fade(node, params);
      default:
        return slide(node, params);
    }
  };

  const ctx: AccordionCtxType = getContext("ctx") ?? {};

  // single selection
  const self = {};
  const selected = ctx.selected ?? writable();

  let _open: boolean = open;
  open = false;

  onMount(() => {
    if (_open) $selected = self;

    // this will trigger unsubscribe on destroy
    return selected.subscribe((x: any) => (open = x === self));
  });

  const handleToggle = async (_: Event) => {
    if (toggleCallback) {
      await toggleCallback();
    }
    selected.set(open ? {} : self);
  };

  let buttonClass: string = $derived(
    twMerge([
      defaultClass,
      ctx.flush ? "" : borderClass,
      borderBottomClass,
      borderSharedClass,
      ctx.flush ? paddingFlush : paddingDefault,
      open && (ctx.flush ? textFlushOpen : activeCls || ctx.activeClass),
      !open && (ctx.flush ? textFlushDefault : inactiveCls || ctx.inactiveClass),
      restProps.class ? `${restProps.class}` : ""
    ])
  );

  let contentClass: string = $derived(
    twMerge([
      ctx.flush ? paddingFlush : paddingDefault,
      ctx.flush ? "" : borderOpenClass,
      borderBottomClass,
      borderSharedClass
    ])
  );
</script>

<svelte:element this={tag} class="group">
  <button onclick={handleToggle} type="button" {id} class={buttonClass} aria-expanded={open}>
    {#if headerSlot}
      {@render headerSlot?.()}
    {:else if open}
      {#if arrowup}
        {@render arrowup()}
      {:else}
        <svg
          class="h-3 w-3 text-gray-800 dark:text-white"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 10 6"
        >
          <path
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 5 5 1 1 5"
          />
        </svg>
      {/if}
    {:else if arrowdown}
      {@render arrowdown()}
    {:else}
      <svg
        class="h-3 w-3 text-gray-800 dark:text-white"
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 10 6"
      >
        <path
          stroke="currentColor"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="m1 1 4 4 4-4"
        />
      </svg>
    {/if}
  </button>
</svelte:element>
{#if open}
  <div transition:multiple={transitionParams}>
    <div class={contentClass}>
      {@render children?.()}
    </div>
  </div>
{/if}
