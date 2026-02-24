<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { A, Li, List, P } from "flowbite-svelte";
  import { loadAdvisoryVersions, type AdvisoryVersion } from "$lib/Advisories/advisory";
  import { onMount } from "svelte";
  import { addSlashes } from "$lib/utils";
  import type { ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

  interface Props {
    params: any;
    document: any;
  }

  let { params = null, document = {} }: Props = $props();

  const uid = $props.id();

  let loadAdvisoryVersionsError: ErrorDetails | null = $state(null);
  let advisoryVersions: AdvisoryVersion[] = $state([]);

  let suggestedPublisherNamespace = $derived(document.publisher.name);
  let suggestedTrackingID = $derived(document.tracking.id);
  let encodedTrackingID = $derived(
    document.tracking?.id ? encodeURIComponent(addSlashes(suggestedTrackingID)) : undefined
  );
  let encodedPublisherNamespace = $derived(
    document.publisher?.name
      ? encodeURIComponent(addSlashes(suggestedPublisherNamespace))
      : undefined
  );
  let suggestedLinkByID = $derived(
    `/advisories/${suggestedPublisherNamespace}/${suggestedTrackingID}/documents/${params.id}`
  );
  let encodedSuggestedLinkByID = $derived(
    `/advisories/${encodedPublisherNamespace}/${encodedTrackingID}/documents/${params.id}`
  );

  const getAdvisoryVersions = async () => {
    if (!encodedTrackingID || !encodedPublisherNamespace) return;
    const result = await loadAdvisoryVersions(
      encodeURIComponent(addSlashes(params.trackingID)),
      encodeURIComponent(addSlashes(params.publisherNamespace))
    );
    if (result) {
      if (result.error) {
        loadAdvisoryVersionsError = result.error;
      } else if (result.advisoryVersions) {
        advisoryVersions = result.advisoryVersions;
      }
    }
  };

  onMount(() => {
    getAdvisoryVersions();
  });
</script>

<div class="flex flex-col gap-6">
  <div class="mb-2 font-bold">
    <i class="bx bx-error-circle" aria-hidden="true"></i>
    <span>The URL doesn't reference any document</span>
  </div>
  <P>
    Do you want to open the document with ID {params.id}?
    <br />
    <A href={`#${encodedSuggestedLinkByID}`}>{suggestedLinkByID}</A>
  </P>
  <ErrorMessage error={loadAdvisoryVersionsError}></ErrorMessage>
  {#if advisoryVersions}
    {#if advisoryVersions.length === 1}
      <P>
        Or do you want to open the following document with publisher {params.publisherNamespace}
        and tracking ID {params.trackingID}?
        <br />
        <A
          href={`#/advisories/${params.publisherNamespace}/${params.trackingID}/documents/${advisoryVersions[0].id}`}
          >/advisories/{params.publisherNamespace}/{params.trackingID}/documents/{advisoryVersions[0]
            .id}
        </A>
      </P>
    {:else if advisoryVersions.length > 1}
      <P>
        Or do you want to open one of the following documents with publisher {params.publisherNamespace}
        and tracking ID {params.trackingID}?
        <List tag="ul" class="space-y-1 text-gray-500 dark:text-gray-400">
          {#each advisoryVersions as version, i (`inconsistencymessage-${uid}-${i}`)}
            <Li>
              <A
                href={`#/advisories/${params.publisherNamespace}/${params.trackingID}/documents/${version.id}`}
                >/advisories/{params.publisherNamespace}/{params.trackingID}/documents/{version.id}
              </A>
            </Li>
          {/each}
        </List>
      </P>
    {/if}
  {/if}
</div>
