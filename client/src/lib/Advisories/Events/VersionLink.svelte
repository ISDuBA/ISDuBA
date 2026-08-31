<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { getContext } from "svelte";
  import type { AdvisoryVersion } from "../advisory.svelte";
  import type { AdvisoryRouteParams } from "$routes/routes";
  import { A } from "flowbite-svelte";

  let advisoryVersions: () => AdvisoryVersion[] = getContext("advisoryVersions");
  let params: () => AdvisoryRouteParams = getContext("params");

  let pathPrefix = $derived(
    "#/advisories/" + params().publisherNamespace + "/" + params().trackingID + "/documents/"
  );

  interface Props {
    documentID: number;
  }
  let { documentID }: Props = $props();

  let versionOfCurrentDoc = $derived(advisoryVersions().find((v) => v.id === Number(params().id)));

  let version = $derived(advisoryVersions().find((v) => v.id === documentID));

  let otherStatusAvailable = $derived(
    advisoryVersions().some(
      (v) => v.version === version?.version && v.tracking_status !== version.tracking_status
    )
  );
</script>

{#if versionOfCurrentDoc?.id === version?.id}
  on version: {version?.version}
  {otherStatusAvailable ? `(${version?.tracking_status})` : ""}
{:else}
  <A class="text-primary-700 dark:text-primary-400" href={`${pathPrefix}${documentID}/`}
    >on version: {version?.version}
    {otherStatusAvailable ? `(${version?.tracking_status})` : ""}</A
  >
{/if}
