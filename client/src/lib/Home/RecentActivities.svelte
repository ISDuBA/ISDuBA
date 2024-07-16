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
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import Activity from "./Activity.svelte";
  import { Badge } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  const recentActivityQuery = `/api/events?limit=10&count=true&query=$event import_document events != now 168h duration - $time <= $actor me != me involved`;
  const mentionedQuery = `/api/events?limit=10&count=true&query=$event import_document events != now 168h duration - $time <=  me mentioned  and`;
  const documentQueryBase = `/api/documents?columns=id title publisher tracking_id`;
  const pluck = (arr: any, keys: any) => arr.map((i: any) => keys.map((k: any) => i[k]));
  let activityCount = 0;
  let resultingActivities: any;
  let loadActivityError = "";
  let loadMentionsError = "";
  let loadDocumentsError = "";
  let loadCommentsError = "";

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
    const activitiesResponse = await request(recentActivityQuery, "GET");
    if (activitiesResponse.ok) {
      const activities = await activitiesResponse.content;
      activityCount = activities.count;
      return activities.events;
    } else if (activitiesResponse.error) {
      loadActivityError = `Could not load Activities. ${getErrorMessage(activitiesResponse.error)}. ${getErrorMessage(activitiesResponse.content)}`;
      return [];
    }
  };

  const fetchMentions = async () => {
    const mentionsResponse = await request(mentionedQuery, "GET");
    if (mentionsResponse.ok) {
      const mentions = await mentionsResponse.content;
      return mentions.events;
    } else if (mentionsResponse.error) {
      loadMentionsError = `Could not load Activities. ${getErrorMessage(mentionsResponse.error)}. ${getErrorMessage(mentionsResponse.content)}`;
      return [];
    }
  };

  const getDocumentIDs = (arr: any) => {
    return arr.map((a: any) => {
      return a[0];
    });
  };

  const getCommentIDs = (arr: any) => {
    return arr
      .map((a: any) => {
        return a[1];
      })
      .filter((id: any) => {
        return id !== null;
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
      loadDocumentsError = `Could not load Documents. ${getErrorMessage(result.error)}. ${getErrorMessage(result.content)}`;
      return [];
    }
  };

  const fetchComments = async (commentIDs: number[]) => {
    if (commentIDs.length > 0) {
      const promises = await Promise.allSettled(
        commentIDs.map(async (v: any) => {
          return request(`/api/comments/post/${v}`, "GET");
        })
      );
      const result = promises
        .filter((p: any) => p.status === "fulfilled" && p.value.ok)
        .map((p: any) => {
          return p.value;
        });
      if (promises.length != result.length) {
        loadCommentsError = `Could not load all comments. An error occured on the server. Please contact an administrator.`;
        return [];
      }
      return result.map((r) => {
        return r.content;
      });
    } else {
      loadCommentsError = `Could not load comments. An error occured on the server. Please contact an administrator.`;
      return [];
    }
  };

  const transformDataToActivities = async () => {
    const activities = await fetchActivities();
    const mentions = await fetchMentions();
    let idsActivities = [];
    let idsMentions = [];
    let documentIDs: any[] = [];
    let commentIDs: any[] = [];
    let documents = [];
    let comments = [];
    if (activities.length > 0) {
      idsActivities = pluck(activities, ["id", "comments_id"]);
      documentIDs = getDocumentIDs(idsActivities);
      commentIDs = getCommentIDs(idsActivities);
    }
    if (mentions.length > 0) {
      idsMentions = pluck(mentions, ["id", "comments_id"]);
      const mentionDocumentIDs = getDocumentIDs(idsMentions);
      documentIDs.concat(mentionDocumentIDs);
      const mentionCommentIDs = getCommentIDs(idsMentions);
      commentIDs = commentIDs.concat(mentionCommentIDs);
    }
    documentIDs = [...new Set(documentIDs)];
    commentIDs = [...new Set(commentIDs)];
    if (documentIDs) {
      documents = await fetchDocuments(documentIDs);
    }
    if (commentIDs) {
      comments = await fetchComments(commentIDs);
    }
    const documentsById = documents.reduce((o: any, n: any) => {
      o[n.id] = n;
      return o;
    }, {});
    const commentsByID = comments.reduce((o: any, n: any) => {
      o[n.id] = n.message;
      return o;
    }, {});
    console.log(commentsByID);
    const activitiesAggregated = aggregateNewest(activities);
    let recentActivities = Object.values(activitiesAggregated);
    recentActivities = recentActivities.map((a: any) => {
      a.mention = false;
      a.documentTitle = documentsById[a.id] ? documentsById[a.id]["title"] : "";
      a.documentURL = documentsById[a.id]
        ? `/advisories/${documentsById[a.id]["publisher"]}/${documentsById[a.id]["tracking_id"]}/documents/${a.id}`
        : "";
      if (a.event === "change_comment" || a.event === "add_comment")
        a.message = commentsByID[a.comments_id];
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
      if (a.event === "change_comment" || a.event === "add_comment")
        a.message = commentsByID[a.comments_id];

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
</script>

{#if $appStore.app.isUserLoggedIn}
  <div class="flex w-1/2 max-w-[50%] flex-col gap-4">
    <SectionHeader title="Recent activities"></SectionHeader>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if resultingActivities}
        {#if resultingActivities.length > 0}
          {#each resultingActivities as activity}
            <Activity
              on:click={() => {
                push(activity.documentURL);
              }}
            >
              <span slot="top-right">{getRelativeTime(new Date(activity.time))}</span>
              <span slot="top-left">
                {#if activity.mention}
                  {activity.actor} mentioned you
                {:else if activity.event === "add_comment"}
                  {activity.actor} commented on {activity.documentTitle}
                {:else if activity.event === "add_ssvc"}
                  {activity.actor} added a SSVC
                {:else if activity.event === "import_document"}
                  {activity.actor} added imported a document
                {:else if activity.event === "change_ssvc" || activity.event === "change_sscv"}
                  {activity.actor} changed a SSVC
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
                  <i class="bx bxs-quote-alt-left"></i>
                  <span class="italic"
                    >{activity.message?.length < 30
                      ? activity.message
                      : activity.message?.substring(0, 30)}</span
                  >
                </div>
              {:else}
                <div>
                  {activity.documentTitle}
                </div>
              {/if}
              <span class="text-gray-400" slot="bottom-left">
                {activity.event === "add_comment" || activity.event == "change_comment"
                  ? `${activity.documentTitle}`
                  : ""}
              </span>
            </Activity>
          {/each}
        {:else}
          No recent activities on advisories you are involved in.
        {/if}
      {/if}
      {#if activityCount > 10}<div class="">…There are more activities</div>{/if}
    </div>
    <ErrorMessage message={loadActivityError}></ErrorMessage>
    <ErrorMessage message={loadMentionsError}></ErrorMessage>
    <ErrorMessage message={loadDocumentsError}></ErrorMessage>
    <ErrorMessage message={loadCommentsError}></ErrorMessage>
  </div>
{/if}
