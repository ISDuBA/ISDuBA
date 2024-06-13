<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Badge } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";

  let error = "";
  let recentActivities = [
    {
      name: "mention",
      user: "Max",
      content: "@beate, thank you for the hint. It looks really promising!",
      date: new Date(Date.now() - 3000)
    },
    {
      name: "comment",
      user: "Beate",
      content: "Here is another page where you can find more details about the vulnerabilities...",
      date: new Date(Date.now() - 180000)
    },
    {
      name: "state",
      newState: "review",
      publisher: "BSI",
      documentTitle: "CVRF-CSAF-Converter: XML External Entities Vulnerability",
      ssvc: convertVectorToLabel("SSVCv2/E:P/A:Y/T:T/P:E/B:A/M:H/D:A/2024-06-12T12:42:25Z/"),
      date: new Date(Date.now() - 22000000)
    },
    {
      name: "import",
      publisher: "BSI",
      documentTitle:
        "Stack Buffer Overflow vulnerability in FastStone Image Viewer 7.5 and earlier",
      cvss: 10.0,
      date: new Date(Date.now() - 780000000)
    }
  ];

  const getRelativeTime = (date: Date) => {
    const now = Date.now();
    const passedTime = now - date.getTime();
    if (passedTime < 60000) {
      return "Less than one minute ago";
    } else if (passedTime < 3600000) {
      return `${Math.floor(passedTime / 60000)} min ago`;
    } else if (passedTime < 86400000) {
      return `${Math.floor(passedTime / 3600000)} hours ago`;
    } else {
      return `${Math.floor(passedTime / 86400000)} days ago`;
    }
  };
</script>

{#if $appStore.app.isUserLoggedIn}
  <div class="flex w-1/2 max-w-[50%] flex-col gap-4">
    <SectionHeader title="Recent activities"></SectionHeader>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if recentActivities?.length && recentActivities.length > 0}
        {#each recentActivities as activity}
          <div class="rounded-md border border-solid border-gray-300 p-4 shadow-md">
            <div class="mb-3 flex flex-row justify-between text-xs">
              {#if activity.date}
                <span title={activity.date.toISOString()}>{getRelativeTime(activity.date)}</span>
                {#if activity.name}
                  {#if activity.name === "comment"}
                    <i class="bx bx-comment text-xl"></i>
                  {:else if activity.name === "mention"}
                    <i class="bx bx-at text-xl"></i>
                  {:else if activity.name === "state"}
                    <i class="bx bx-book-open text-xl"></i>
                  {:else if activity.name === "import"}
                    <i class="bx bx-import text-xl"></i>
                  {/if}
                {/if}
              {/if}
            </div>
            <hr class="mb-3" />
            {#if ["comment", "mention"].includes(activity.name)}
              <div class="text-black">{activity.content}</div>
            {:else if activity.name === "state"}
              <div class="text-black">
                {`${activity.publisher}: ${activity.documentTitle}`}
              </div>
              <div class="flex justify-between">
                <Badge
                  class="h-6 w-fit"
                  title={activity.ssvc.vector}
                  style={`color: white; background-color: ${activity.ssvc.color}`}
                >
                  {activity.ssvc.label}
                </Badge>
                <div>{`New state: ${activity.newState}`}</div>
              </div>
            {:else if activity.name === "import"}
              <div class="text-black">
                {`${activity.publisher}: ${activity.documentTitle}`}
              </div>
              <div class="font-bold text-red-700">CVSS: {activity.cvss?.toFixed(1)}</div>
            {/if}
          </div>
        {/each}
      {/if}
    </div>
    <ErrorMessage message={error}></ErrorMessage>
  </div>
{/if}
