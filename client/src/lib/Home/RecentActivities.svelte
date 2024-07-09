<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Badge } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";
  import Activity from "./Activity.svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";

  let recentActivityQuery = `/api/events?limit=10&count=true&query=$event import_document events != now 168h duration - $time <= $actor me != me mentioned me involved or and`;
  let loadActivityError = "";
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
      return "<1 min ago";
    } else if (passedTime < 3600000) {
      return `${Math.floor(passedTime / 60000)} min ago`;
    } else if (passedTime < 86400000) {
      return `${Math.floor(passedTime / 3600000)} hours ago`;
    } else {
      return `${Math.floor(passedTime / 86400000)} days ago`;
    }
  };

  // const transformActivities = (result: any) => {
  //   let { actor, comments_id, event, event_state, id, time } = result;
  // };

  const fetchData = async () => {
    const response = await request(recentActivityQuery, "GET");
    if (response.ok) {
      //const result = await response.content;
      // transformActivities(result);
    } else if (response.error) {
      loadActivityError = `Could not load Activities. ${getErrorMessage(response.error)}. ${getErrorMessage(response.content)}`;
    }
  };
  onMount(() => {
    fetchData();
  });
</script>

{#if $appStore.app.isUserLoggedIn}
  <div class="flex w-1/2 max-w-[50%] flex-col gap-4">
    <SectionHeader title="Recent activities"></SectionHeader>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if recentActivities?.length && recentActivities.length > 0}
        {#each recentActivities as activity}
          <Activity>
            <div slot="top-right">
              {#if activity.cvss}
                <span>CVSS v3:</span>
                <span class="text-red-500">
                  {activity.cvss?.toFixed(1)}
                </span>
              {:else if activity.ssvc}
                <Badge
                  title={activity.ssvc.vector}
                  style={`background-color: white; color: black; border: 1pt solid ${activity.ssvc.color};`}
                >
                  {activity.ssvc.label}
                </Badge>
              {/if}
            </div>
            <div slot="top-left" class="text-sm text-gray-700">
              {#if activity.name}
                {#if activity.name === "comment"}
                  <div class="flex items-center gap-1">
                    <i class="bx bx-comment"></i>
                    <span>{`New comment from ${activity.user}`}</span>
                  </div>
                {:else if activity.name === "mention"}
                  <div class="flex items-center gap-1">
                    <i class="bx bx-at"></i>
                    <span>{`${activity.user} mentioned you`}</span>
                  </div>
                {:else if activity.name === "state"}
                  <div class="flex items-center gap-1">
                    <i class="bx bx-book-open"></i>
                    <span>{`Now in state ${activity.newState}`}</span>
                  </div>
                {:else if activity.name === "import"}
                  <div class="flex items-center gap-1">
                    <i class="bx bx-import"></i>
                    <span>New import</span>
                  </div>
                {/if}
              {/if}
            </div>
            {#if ["comment", "mention"].includes(activity.name)}
              <div class="italic">{`${activity.content}`}</div>
            {:else if activity.name === "state"}
              <div class="text-black">
                {`${activity.publisher}: ${activity.documentTitle}`}
              </div>
            {:else if activity.name === "import"}
              <div class="text-black">
                {`${activity.publisher}: ${activity.documentTitle}`}
              </div>
            {/if}
            <div slot="bottom-right">
              {#if activity.date}
                <span title={activity.date.toISOString()} class="text-gray-500"
                  >{getRelativeTime(activity.date)}</span
                >
              {/if}
            </div>
          </Activity>
        {/each}
      {/if}
    </div>
    <ErrorMessage message={loadActivityError}></ErrorMessage>
  </div>
{/if}
