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
  import { Badge, Button } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { convertVectorToSSVCObject } from "$lib/Advisories/SSVC/SSVCCalculator";
  import { getRelativeTime } from "./activity";
  import SsvcBadge from "$lib/Advisories/SSVC/SSVCBadge.svelte";

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

  const documentQueryBase = `/api/documents?columns=id title publisher tracking_id ssvc`;
  const pluck = (arr: any, keys: any) => arr.map((i: any) => keys.map((k: any) => i[k]));
  let activityCount = 0;
  let resultingActivities: any;
  let loadActivityError: ErrorDetails | null;
  let loadMentionsError: ErrorDetails | null;
  let loadDocumentsError: ErrorDetails | null;

  const aggregateNewest = (events: any) => {
    return events.reduce((o: any, n: any) => {
      const key = `${n.actor}|${n.comments_id}|${n.event}|${n.id}`; // create natural key (actor, comment, event, document) to identify event
      if (!o[key]) o[key] = n; // if not already aggregated add event
      if (o[key].time < n.time) o[key] = n; // replace older events of the same type with the newest one
      return o;
    }, {});
  };

  const sortByTime = (events: any) => {
    return events.sort((a: any, b: any) => {
      return a.time < b.time;
    });
  };

  const aggregateByMentions = (events: any) => {
    return Object.values(
      events.reduce((o: any, n: any) => {
        const key = `${n.actor}|${n.comments_id}|${n.event}|${n.id}`; // create natural key (actor, comment, event, document) to identify event
        if (!o[key]) o[key] = n; // if not already aggregated add event
        if (!o[key].mention) o[key] = n; // aggregate nonmentions with mentioned one
        return o;
      }, {})
    );
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
      loadActivityError = getErrorDetails(`Could not load Activities.`, activitiesResponse);
      return [];
    }
  };

  const fetchMentions = async () => {
    const mentionedQuery = `/api/events?limit=10&count=true&query=$event import_document events != now 168h duration - $time <=  me mentioned  and`;
    const mentionsResponse = await request(mentionedQuery, "GET");
    if (mentionsResponse.ok) {
      const mentions = await mentionsResponse.content;
      return mentions.events || [];
    } else if (mentionsResponse.error) {
      loadMentionsError = getErrorDetails(`Could not load Activities.`, mentionsResponse);
      return [];
    }
  };

  const getDocumentIDs = (arr: any) => {
    return arr.map((a: any) => {
      return a[0];
    });
  };

  const fetchDocuments = async (documentIDs: number[]) => {
    const query = documentIDs.map((id: number) => {
      return `$id ${id} integer = `;
    });
    const ors = documentIDs.map(() => {
      return "or ";
    });
    ors.shift(); // we need one less or than ids
    const documentQuery = `${documentQueryBase}&query=${query.join("")}${ors.join("")}`;
    const result = await request(documentQuery, "GET");
    if (result.ok) {
      const content = result.content;
      return content.documents;
    } else if (result.error) {
      loadDocumentsError = getErrorDetails(`Could not load Documents.`, result);
      return [];
    }
  };

  const transformDataToActivities = async () => {
    const activities = await fetchActivities();
    const mentions = await fetchMentions();
    let idsActivities = [];
    let idsMentions = [];
    let documentIDs: any[] = [];
    let documents = [];
    if (
      activities.length > 0 &&
      storedQuery.columns.includes("id") &&
      storedQuery.columns.includes("comments_id")
    ) {
      idsActivities = pluck(activities, ["id", "comments_id"]);
      documentIDs = getDocumentIDs(idsActivities);
    }
    if (mentions.length > 0) {
      idsMentions = pluck(mentions, ["id", "comments_id"]);
      const mentionDocumentIDs = getDocumentIDs(idsMentions);
      documentIDs = documentIDs.concat(mentionDocumentIDs);
    }

    documentIDs = [...new Set(documentIDs)];
    if (documentIDs.length > 0) {
      documents = await fetchDocuments(documentIDs);
    }
    const documentsById = documents.reduce((o: any, n: any) => {
      o[n.id] = n;
      return o;
    }, {});

    const activitiesAggregated = aggregateNewest(activities);
    let recentActivities = Object.values(activitiesAggregated);
    recentActivities = recentActivities.map((a: any) => {
      a.mention = a.message && a.message.includes($appStore.app.tokenParsed?.preferred_username);
      a.documentTitle = documentsById[a.id] ? documentsById[a.id]["title"] : "";
      a.documentURL = documentsById[a.id]
        ? `/advisories/${documentsById[a.id]["publisher"]}/${documentsById[a.id]["tracking_id"]}/documents/${a.id}`
        : "";
      if (a.event === "add_ssvc" || a.event === "change_ssvc" || a.event === "change_sscv")
        a.ssvc = documentsById[a.id] ? documentsById[a.id]["ssvc"] : "";
      if (a.ssvc) a.ssvcLabel = convertVectorToSSVCObject(a.ssvc).label;
      return a;
    });

    const mentionsAggregated = aggregateNewest(mentions);
    let recentMentions = Object.values(mentionsAggregated);
    recentMentions = recentMentions.map((a: any) => {
      a.mention = true;
      a.documentTitle = documentsById[a.id] ? documentsById[a.id]["title"] : "";
      a.documentURL = documentsById[a.id]
        ? `/advisories/${documentsById[a.id]["publisher"]}/${documentsById[a.id]["tracking_id"]}/documents/${a.id}`
        : "";
      return a;
    });
    const activitiesAggregatedByMentions = aggregateByMentions([
      ...recentActivities,
      ...recentMentions
    ]);
    resultingActivities = sortByTime(aggregateByChange(activitiesAggregatedByMentions));
  };

  onMount(async () => {
    transformDataToActivities();
  });

  const showMore = () => {
    push(`/search?query=${storedQuery.id}`);
  };
</script>

{#if $appStore.app.isUserLoggedIn && (appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor())}
  <div class="flex flex-col gap-4 md:w-[46%] md:max-w-[46%]">
    <SectionHeader title={storedQuery.description}></SectionHeader>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if resultingActivities}
        {#if resultingActivities.length > 0}
          {#each resultingActivities as activity}
            <Activity
              on:click={() => {
                push(activity.documentURL);
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
                  <span class="italic"
                    >{activity.message && activity.message?.length < 30
                      ? activity.message
                      : (activity.message?.substring(0, 30) ?? "Message undefined")}</span
                  >
                </div>
              {:else}
                <div>
                  {activity.documentTitle ?? "Title undefined"}
                </div>
                <div class="text-sm text-gray-700">
                  {activity.tracking_id ?? "Trackind ID undefined"}
                </div>
              {/if}
              <span class="text-gray-400" slot="bottom-left">
                {activity.event === "add_comment" || activity.event == "change_comment"
                  ? `${activity.documentTitle ?? "Title undefined"}`
                  : ""}
              </span>
              <div slot="bottom-bottom">
                <div class="flex items-center gap-4 text-xs text-gray-500">
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
                  <div class="my-2 rounded-sm border p-2 text-xs text-gray-800">
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
          <Button on:click={showMore} size="xs" color="light" class="h-6 w-fit rounded-md"
            >More...</Button
          >
        {:else}
          <div class="text-gray-600">No matching events found.</div>
        {/if}
      {/if}
      {#if activityCount > 10}<div class="">…There are more events</div>{/if}
    </div>
    <ErrorMessage error={loadActivityError}></ErrorMessage>
    <ErrorMessage error={loadMentionsError}></ErrorMessage>
    <ErrorMessage error={loadDocumentsError}></ErrorMessage>
  </div>
{/if}
