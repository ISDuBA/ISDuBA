<script lang="ts">
  import { DarkMode } from "flowbite-svelte";
  import { appStore } from "./store";
  import { onMount } from "svelte";
  export let btnClass: string =
    "text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none rounded-lg text-sm p-2.5";

  onMount(() => {
    appStore.updateDarkMode();

    // Make a mutation observer to watch for class changes on ownerDocument.documentElement
    // and update the theme accordingly
    const darkModeObserver = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === "class") {
          appStore.updateDarkMode();
        }
      });
    });
    // Observe changes to the class attribute of the ownerDocument.documentElement
    darkModeObserver.observe(document.documentElement, {
      attributes: true
    });
  });
</script>

<DarkMode {btnClass} />
