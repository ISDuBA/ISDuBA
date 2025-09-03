<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  /**
   * updateUI waits until UI is settled and goes back to the last anchor.
   * @param id
   */
  async function updateUI(id: string) {
    setTimeout(() => {
      const element = document.getElementById(`${id}`);
      const y = element!.getBoundingClientRect().top + window.scrollY - 150;
      window.scrollTo({ top: y, behavior: "smooth" });
    }, 200);
  }
  /**
   * backpressed handles the history when the button was actually clicked.
   * @param e
   */
  const backPressed = (e: Event) => {
    const lastElement = appStore.state.webview.ui.history[0];
    appStore.shiftHistory();
    updateUI(lastElement);
    e.preventDefault();
  };
</script>

{#if appStore.state.webview.ui.history.length > 0}
  <a class="backbutton" href="#top" on:click={backPressed}>Last pos. <i class="bx bx-undo"></i></a>
{/if}
