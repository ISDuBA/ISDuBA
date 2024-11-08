<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { DarkMode } from "flowbite-svelte";
  import { appStore } from "../store";
  import { onMount } from "svelte";
  export let btnClass: string =
    "text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none rounded-lg text-sm p-2.5";

  onMount(() => {
    appStore.updateDarkMode();
    const darkModeObserver = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === "class") {
          appStore.updateDarkMode();
        }
      });
    });
    darkModeObserver.observe(document.documentElement, {
      attributes: true
    });
  });
</script>

<DarkMode {btnClass} />
