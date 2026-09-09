<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { Status } from "$lib/Advisories/types/docmodeltypes";
  import { getReadableDateString } from "../helpers";
  import Cvss from "./CVSS.svelte";
  import { Button } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { getContext } from "svelte";
  import Link from "$lib/Components/Link.svelte";
  import SearchableText from "../SearchableText.svelte";
  import { ArrowOutUpRightSquare, Link as LinkIcon } from "@boxicons/svelte";
  import {
    getInitialReleaseDate,
    getCurrentReleaseDate,
    getPublisher,
    getStatus,
    getTitle,
    getGenerator,
    getTrackingVersion,
    getCategory,
    getCSAFVersion,
    getDistributionText,
    getHighestScore,
    getLang,
    getSourceLang
  } from "$lib/Advisories/docmodel";

  interface Props {
    basePath: string;
  }
  let { basePath = "" }: Props = $props();

  let doc = $derived(appStore.state.webview.doc);
  let trackingVersion = $derived(getTrackingVersion(doc));
  let generator = $derived(getGenerator(doc));
  let publisher = $derived(getPublisher(doc));
  let publisherName = $derived(publisher?.name ?? undefined);
  let publisherCategory = $derived(publisher?.category ?? undefined);
  let publisherNamespace = $derived(publisher?.namespace ?? undefined);
  let publisherIssuingAuthority = $derived(publisher?.issuing_authority ?? undefined);
  let publisherContactDetails = $derived(publisher?.contact_details ?? undefined);
  let aggregateSeverity = $derived(doc?.document.aggregate_severity);
  let category = $derived(getCategory(doc));
  let title = $derived(getTitle(doc));
  let lang = $derived(getLang(doc));
  let sourceLang = $derived(doc != null ? getSourceLang(doc) : undefined);
  let csafVersion = $derived(getCSAFVersion(doc));
  let distributionText = $derived(getDistributionText(doc));
  let published = $derived(doc != null ? getInitialReleaseDate(doc) : undefined);
  let lastUpdate = $derived(doc != null ? getCurrentReleaseDate(doc) : undefined);
  let status = $derived(getStatus(doc));
  let highestScore = $derived(getHighestScore(doc));
  let baseSeverity: string | null = $derived((highestScore?.baseSeverity as string) ?? null);
  let baseScore: string | undefined = $derived(
    highestScore?.baseScore ? `${highestScore?.baseScore}` : undefined
  );
  const cellStyleValue = "content-center px-6 py-0 [word-wrap:break-word] hyphens-auto";
  const cellStyleKey = "content-center w-40 max-w-full py-0 text-balance";

  const openRelatedDocuments = () => {
    // Use push of external router since we want PrevNext to disappear when user navigates to related
    // documents.
    push(`${basePath}${basePath.endsWith("/") ? "" : "/"}related/documents`);
  };

  const relatedDocuments: any = getContext("advisory");
</script>

<div class="w-full">
  <div class="mb-3">
    <div class="flex flex-row gap-2">
      <div>
        <span class="-mt-1 inline-block text-xl text-balance">
          <SearchableText text={title} textPath="/document/title" />
        </span>
        {#if status !== Status.final}
          <span class="ml-3 text-lg text-gray-400">{status}</span>
        {/if}
      </div>
      {#if highestScore}
        <Cvss {baseScore} baseSeverity={baseSeverity ?? ""}></Cvss>
      {/if}
      {#if relatedDocuments?.()}
        {@const len = Object.keys(relatedDocuments()).length}
        {#if len > 0}
          <Button
            onclick={openRelatedDocuments}
            color="light"
            size="xs"
            class="h-7"
            title="Open related documents"
          >
            <div class="flex items-center">
              <LinkIcon rotate={90} />
              {len}
            </div>
          </Button>
        {/if}
      {/if}
    </div>
  </div>
  <div class="flex w-full flex-row flex-wrap">
    <div class="grid w-full grid-cols-[minmax(140px,1fr)_auto] gap-1.5 text-sm">
      <div class={cellStyleKey}>Publisher name</div>
      <div class={cellStyleValue}>
        <SearchableText text={publisherName} textPath="/document/publisher/name" />
      </div>
      <div class={cellStyleKey}>Publisher namespace</div>
      <div class={cellStyleValue}>
        {#if publisherNamespace}
          <Link href={publisherNamespace} class="underline">
            <ArrowOutUpRightSquare />
            <SearchableText text={publisherNamespace} textPath="/document/publisher/namespace" />
          </Link>
        {/if}
      </div>
      {#if publisherContactDetails}
        <div class={cellStyleKey}>Publisher contact details</div>
        <div class={cellStyleValue}>
          <SearchableText
            text={publisherContactDetails}
            textPath="/document/publisher/contact_details"
          />
        </div>
      {/if}
      {#if publisherIssuingAuthority}
        <div class={cellStyleKey}>Publisher issuing authority</div>
        <div class={cellStyleValue}>
          <SearchableText
            text={publisherIssuingAuthority}
            textPath="/document/publisher/issuing_authority"
          />
        </div>
      {/if}
      <div class={cellStyleKey}>Publisher category</div>
      <div class={cellStyleValue}>
        <SearchableText text={publisherCategory} textPath="/document/publisher/category" />
      </div>
      <div class={cellStyleKey}>Published</div>
      <div class={cellStyleValue}>
        <SearchableText
          text={getReadableDateString(published)}
          textPath="/document/tracking/initial_release_date"
        />
      </div>
      <div class={cellStyleKey}>Last update</div>
      <div class={cellStyleValue}>
        <SearchableText
          text={getReadableDateString(lastUpdate)}
          textPath="/document/tracking/current_release_date"
        />
      </div>
      <div class={cellStyleKey}>CSAF-Version</div>
      <div class={cellStyleValue}>
        <SearchableText text={csafVersion} textPath="/document/csaf_version" />
      </div>
      <div class={cellStyleKey}>Category</div>
      <div class={cellStyleValue}>
        <SearchableText text={category} textPath="/document/category" />
      </div>
      {#if distributionText}
        <div class={cellStyleKey}>Distribution</div>
        <div class={cellStyleValue}>
          <SearchableText text={distributionText} textPath="/document/distribution/text" />
        </div>
      {/if}
      {#if aggregateSeverity}
        <div class={cellStyleKey}>Aggregate severity text</div>
        <div class={cellStyleValue}>
          <SearchableText
            text={aggregateSeverity.text}
            textPath="/document/aggregate_severity/text"
          />
        </div>
        {#if aggregateSeverity.namespace}
          <div class={cellStyleKey}>Aggregate severity namespace</div>
          <div class={cellStyleValue}>
            <Link href={aggregateSeverity.namespace} class="underline">
              <ArrowOutUpRightSquare />
              <SearchableText
                text={aggregateSeverity.namespace}
                textPath="/document/aggregate_severity/namespace"
              />
            </Link>
          </div>
        {/if}
      {/if}
      {#if lang}
        <div class={cellStyleKey}>Language</div>
        <div class={cellStyleValue}>
          <SearchableText text={lang} textPath="/document/lang" />
        </div>
      {/if}
      {#if sourceLang}
        <div class={cellStyleKey}>Source lang</div>
        <div class={cellStyleValue}>
          <SearchableText text={sourceLang} textPath="/document/source_lang" />
        </div>
      {/if}
      <div class={cellStyleKey}>Tracking Version</div>
      <div class={cellStyleValue}>{trackingVersion}</div>
    </div>
  </div>
  <div class="mt-3 flex flex-row">
    <span class="text-sm text-gray-400">
      Generator:
      {#if generator}
        <SearchableText
          text={generator?.engine.name}
          textPath="/document/tracking/generator/engine/name"
        />
      {/if}
      {#if generator?.engine?.version}
        <SearchableText
          text={generator?.engine.version}
          textPath="/document/tracking/generator/engine/version"
        />
      {/if}
      {#if generator?.date}
        · <SearchableText
          text={getReadableDateString(generator.date)}
          textPath="/document/tracking/generator/date"
        />
      {/if}
    </span>
  </div>
</div>
