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
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import Activity from "./Activity.svelte";
  import { Badge, Spinner } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { convertVectorToSSVCObject } from "$lib/Advisories/SSVC/SSVCCalculator";
  import { getRelativeTime } from "$lib/time";
  import SsvcBadge from "$lib/Advisories/SSVC/SSVCBadge.svelte";
  import ShowMoreButton from "./ShowMoreButton.svelte";

  export let storedQuery: any;
  const ignoredColumns = [
    "documentURL",
    "id",
    "event_state",
    "actor",
    "documentTitle",
    "title",
    "event",
    "time",
    "message",
    "mention",
    "comments_id",
    "ssvc",
    "ssvcLabel",
    "comments",
    "versions",
    "tracking_id"
  ];

  let activityCount = 0;
  let resultingActivities: any;
  let loadActivityError: ErrorDetails | null = null;
  let loadCommentsError: ErrorDetails | null = null;
  let loadDocumentsError: ErrorDetails | null = null;
  let isLoading = false;

  const aggregateNewest = (events: any) => {
    return events.reduce((o: any, n: any) => {
      const key = `${n.actor}|${n.comments_id}|${n.event}|${n.id}`; // create natural key (actor, comment, event, document) to identify event
      if (!o[key]) o[key] = n; // if not already aggregated add event
      if (o[key].time < n.time) o[key] = n; // replace older events of the same type with the newest one
      return o;
    }, {});
  };

  const aggregateByChange = (events: any) => {
    return Object.values(
      events.reduce((o: any, n: any) => {
        const key = `${n.actor}|${n.comments_id}|${n.event}|${n.id}`; // create natural key (actor, comment, event, document) to identify event
        const ADD_COMMENT = "add_comment";
        const CHANGE_COMMENT = "change_comment";
        const add_key = `${n.actor}|${n.comments_id}|${ADD_COMMENT}|${n.id}`;
        const change_key = `${n.actor}|${n.comments_id}|${CHANGE_COMMENT}|${n.id}`;
        if (n.event === ADD_COMMENT && o[change_key]) return o;
        if (n.event === ADD_COMMENT) {
          o[key] = n;
          return o;
        }
        if (n.event === CHANGE_COMMENT && !o[add_key]) {
          o[key] = n;
          return o;
        }
        if (n.event === CHANGE_COMMENT && o[add_key]) {
          delete o[add_key];
        }
        o[key] = n;
        return o;
      }, {})
    );
  };

  const fetchActivities = async () => {
    const columns = `${storedQuery.columns ? "columns=" + storedQuery.columns.join(" ") : ""}`;
    const orders = `${storedQuery.orders ? "&orders=" + storedQuery.orders.join(" ") : ""}`;
    const searchParams = `${columns}&query=${storedQuery.query}&limit=6${orders}`;
    const activitiesResponse = await request(`/api/events?${searchParams}`, "GET");
    if (activitiesResponse.ok) {
      const activities = await activitiesResponse.content;
      activityCount = activities.count;
      return activities.events || [];
    } else if (activitiesResponse.error) {
      loadActivityError = getErrorDetails(`Could not load events.`, activitiesResponse);
      return [];
    }
  };

  const loadComments = async (activities: any[]) => {
    loadCommentsError = null;
    for (let i = 0; i < activities.length; i++) {
      const activity = activities[i];
      if (activity.comments_id) {
        const response = await request(`api/comments/post/${activity.comments_id}`, "GET");
        if (response.ok) {
          activity.message = response.content.message;
        } else if (loadCommentsError !== null) {
          loadCommentsError = getErrorDetails("Could not load comment.", response);
        }
      }
    }
  };

  const transformDataToActivities = async () => {
    let activities = await fetchActivities();
    const activitiesAggregated = aggregateNewest(activities);
    let recentActivities = Object.values(activitiesAggregated);
    recentActivities = recentActivities.map((a: any) => {
      a.mention = a.message && a.message.includes($appStore.app.tokenParsed?.preferred_username);
      if (a.ssvc) a.ssvcLabel = convertVectorToSSVCObject(a.ssvc).label;
      if (a.tracking_id && a.publisher) {
        a.documentURL = `/advisories/${a.publisher}/${a.tracking_id}/documents/${a.id}`;
      }
      return a;
    });
    await loadComments(activities);
    resultingActivities = aggregateByChange(recentActivities);
  };

  onMount(async () => {
    isLoading = true;
    await transformDataToActivities();
    isLoading = false;
  });
</script>

{#if $appStore.app.isUserLoggedIn && (appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor() || appStore.isAdmin())}
  <div class="flex flex-col gap-4 md:w-[46%] md:max-w-[46%]">
    <SectionHeader title={storedQuery.description}></SectionHeader>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if isLoading}
        <div class:invisible={!isLoading} class={isLoading ? "loadingFadeIn" : ""}>
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
      {/if}
      {#if resultingActivities}
        {#if resultingActivities.length > 0}
          {#each resultingActivities as activity}
            <Activity
              on:click={() => {
                if (activity.documentURL) push(activity.documentURL);
              }}
            >
              <div slot="top-right">
                <span>{getRelativeTime(new Date(activity.time))}</span>
              </div>
              <span slot="top-left">
                {#if activity.mention}
                  {activity.actor} mentioned you
                {:else if activity.event === "add_comment"}
                  {activity.actor} added a comment
                {:else if activity.event === "add_ssvc"}
                  {activity.actor} added a SSVC "{activity.ssvcLabel}""
                {:else if activity.event === "import_document"}
                  {activity.actor} imported a document
                {:else if activity.event === "change_ssvc" || activity.event === "change_sscv"}
                  {activity.actor} changed a SSVC to "{activity.ssvcLabel}"
                {:else if activity.event === "change_comment"}
                  {activity.actor} changed a comment
                {:else if activity.event === "state_change"}
                  {activity.actor} changed the state to <Badge color="dark"
                    >{activity.event_state}</Badge
                  >
                {/if}
              </span>
              {#if activity.event === "add_comment" || activity.event == "change_comment"}
                <div>
                  <span class="block overflow-hidden text-ellipsis text-nowrap italic"
                    >{activity.message ?? "Message undefined"}</span
                  >
                </div>
              {:else}
                <div>
                  {activity.title ?? "Title undefined"}
                </div>
                <div class="text-sm text-gray-700 dark:text-gray-400">
                  {activity.tracking_id ?? "Trackind ID undefined"}
                </div>
              {/if}
              <span class="text-gray-400" slot="bottom-left">
                {activity.event === "add_comment" || activity.event == "change_comment"
                  ? `${activity.title ?? "Title undefined"}`
                  : ""}
              </span>
              <div slot="bottom-bottom" class="mt-2">
                <div class="flex items-center gap-4 text-xs text-slate-400">
                  {#if activity.comments !== undefined}
                    <div class="flex items-center gap-1">
                      <i class="bx bx-comment"></i>
                      <span>{activity.comments}</span>
                    </div>
                  {/if}
                  {#if activity.versions !== undefined}
                    <div class="flex items-center gap-1">
                      <i class="bx bx-collection"></i>
                      <span>{activity.versions}</span>
                    </div>
                  {/if}
                  {#if activity.ssvc}
                    <SsvcBadge vector={activity.ssvc}></SsvcBadge>
                  {/if}
                </div>
                {#if Object.keys(activity).filter((k) => !ignoredColumns.includes(k)).length > 0}
                  <div class="my-2 rounded-sm border p-2 text-xs text-gray-800 dark:text-gray-200">
                    {#each Object.keys(activity).sort() as key}
                      {#if !ignoredColumns.includes(key) && activity[key] !== undefined && activity[key] !== null}
                        <div>{key}: {activity[key]}</div>
                      {/if}
                    {/each}
                  </div>
                {/if}
              </div>
            </Activity>
          {/each}
        {:else}
          <div class="text-gray-600 dark:text-gray-400">No matching events found.</div>
        {/if}
      {/if}
    </div>
    <ShowMoreButton id={storedQuery.id}></ShowMoreButton>
    {#if activityCount > 10}<div class="">…There are more events</div>{/if}
    <ErrorMessage error={loadActivityError}></ErrorMessage>
    <ErrorMessage error={loadCommentsError}></ErrorMessage>
    <ErrorMessage error={loadDocumentsError}></ErrorMessage>
  </div>
{/if}
