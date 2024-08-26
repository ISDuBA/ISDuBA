<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { Source } from "$lib/Sources/source";
  import { request } from "$lib/request";
  import {
    Checkbox,
    Input,
    Label,
    Button,
    StepIndicator,
    TableBodyCell,
    Spinner,
    Modal,
    Select,
    Table,
    TableBodyRow
  } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { tdClass } from "$lib/Table/defaults";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails, getErrorDetails } from "$lib/Errors/error";
  import type { CSAFProviderMetadata, DirectoryURL, ROLIEFeed } from "$lib/provider";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { onMount } from "svelte";
  import SourceForm from "./SourceForm.svelte";
  export let params: any = null;

  enum LogLevel {
    debug = "debug",
    info = "info",
    warn = "warn",
    error = "error"
  }

  type Feed = {
    id?: number;
    enable?: boolean;
    url: string;
    label: string;
    rolie: boolean;
    log_level: LogLevel;
  };

  let modalOpen: boolean = false;
  let modalMessage = "";
  let modalTitle = "";
  let modalCallback: any;

  let saveError: ErrorDetails | null;
  let loadError: ErrorDetails | null;
  let feedError: ErrorDetails | null;
  let logError: ErrorDetails | null;
  let pmd: CSAFProviderMetadata;
  let pmdFeeds: Feed[] = [];
  let missingFeeds: Feed[] = [];
  let feeds: Feed[] = [];

  let newFeed: Feed | null;

  let logs: any[] = [];

  let logLevels = [
    { value: LogLevel.error, name: "Error" },
    { value: LogLevel.info, name: "Info" },
    { value: LogLevel.warn, name: "Warn" },
    { value: LogLevel.debug, name: "Debug" }
  ];

  let loadingFeeds: boolean = false;
  let loadingPMD: boolean = false;
  let loadingLogs: boolean = false;

  let currentStep = 0;
  let steps: string[] = ["Source URL", "Source Config", "Feed Selection"];

  let sources: Source[] = [];

  let formClass = "max-w-[800pt]";

  let headers: [string, string][] = [["", ""]];
  let source: Source = {
    name: "",
    url: "",
    active: false,
    rate: 1,
    slots: 2,
    strict_mode: true,
    headers: [""],
    ignore_patterns: [""]
  };

  const saveSource = async (): Promise<boolean> => {
    let method = "POST";
    let path = `/api/sources`;
    const formData = new FormData();
    if (source.id) {
      method = "PUT";
      path += `/${source.id}`;
      formData.append("id", source.id.toString());
    } else {
      formData.append("url", source.url);
    }
    formData.append("name", source.name);
    if (source.active !== undefined) {
      formData.append("active", source.active.toString());
    }
    if (source.rate && source.rate !== 0) {
      formData.append("rate", source.rate.toString());
    }
    if (source.slots && source.slots !== 0) {
      formData.append("slots", source.slots.toString());
    }
    if (source.strict_mode !== undefined) {
      formData.append("strict_mode", source.strict_mode.toString());
    }
    if (source.insecure !== undefined) {
      formData.append("insecure", source.insecure.toString());
    }
    if (source.signature_check !== undefined) {
      formData.append("signature_check", source.signature_check.toString());
    }
    if (source.age != undefined && source.age !== "") {
      formData.append("age", source.age.toString());
    }
    if (source.client_cert_public) {
      formData.append("client_cert_public", source.client_cert_public);
    }
    if (source.client_cert_private) {
      formData.append("client_cert_private", source.client_cert_private);
    }
    if (source.client_cert_passphrase && source.client_cert_passphrase !== "***") {
      formData.append("client_cert_passphrase", source.client_cert_passphrase);
    }
    for (const header of source.headers) {
      if (header != "") {
        formData.append("headers", header);
      }
    }
    for (const pattern of source.ignore_patterns) {
      if (pattern != "") {
        formData.append("ignore_patterns", pattern);
      }
    }
    const resp = await request(path, method, formData);
    if (resp.ok) {
      if (resp.content.id) {
        source.id = resp.content.id;
      }
      saveError = null;
      return true;
    } else if (resp.error) {
      saveError = getErrorDetails(`Could not save source`, resp);
    }
    return false;
  };

  const getLabelName = (feed: ROLIEFeed | DirectoryURL): string => {
    let label = "";
    if (typeof feed === "string") {
      label += feed;
    } else {
      if (feed.summary) {
        label += `${feed.summary} `;
      }
      label += feed.tlp_label;
    }
    return label;
  };

  const parseFeeds = (): Feed[] => {
    let feeds: Feed[] = [];

    let dist = pmd.distributions ?? [];

    for (const entry of dist) {
      if (entry.rolie) {
        for (const feed of entry.rolie.feeds) {
          feeds.push({
            url: feed.url,
            label: getLabelName(feed),
            log_level: LogLevel.error,
            rolie: true,
            enable: true
          });
        }
      }
      if (entry.directory_url) {
        feeds.push({
          url: entry.directory_url,
          label: getLabelName(entry.directory_url),
          log_level: LogLevel.error,
          rolie: false,
          enable: true
        });
      }
    }

    return feeds;
  };

  const saveFeeds = async (feeds: Feed[]) => {
    for (const feed of feeds) {
      if (!feed.enable) {
        continue;
      }
      const formData = new FormData();
      formData.append("url", feed.url);
      formData.append("label", feed.label);
      formData.append("log_level", feed.log_level);
      const resp = await request(`/api/sources/${source.id}/feeds`, "POST", formData);
      if (resp.ok) {
        if (!params?.id) {
          push(`/sources/${source.id}`);
        }
      } else if (resp.error) {
        saveError = getErrorDetails(`Could not save feed.`, resp);
      }
    }
  };

  const getSourceName = async (): Promise<string> => {
    await fetchSources();
    let name = pmd.publisher.name;
    if (!sources.find((s) => s.name === name)) {
      return name;
    }
    for (let i = 1; ; i++) {
      let customName = `${name} #${i}`;
      if (!sources.find((s) => s.name === customName)) {
        return customName;
      }
    }
  };

  const fetchPMD = async () => {
    loadingPMD = true;
    const resp = await request(`/api/pmd?url=${encodeURIComponent(source.url)}`, "GET");
    if (resp.ok) {
      pmd = resp.content;
      source.name = await getSourceName();
      loadError = null;
      currentStep = 1;
      pmdFeeds = parseFeeds();
    } else if (resp.error) {
      loadError = getErrorDetails(`Could not load PMD.`, resp);
    }
    loadingPMD = false;
  };

  const fetchSources = async () => {
    loadingPMD = true;
    const resp = await request(`/api/sources`, "GET");
    loadingPMD = false;
    if (resp.ok) {
      if (resp.content.sources) {
        sources = resp.content.sources;
      }
    } else if (resp.error) {
      loadError = getErrorDetails(`Could not load source`, resp);
    }
  };

  const fetchFeeds = async () => {
    loadingFeeds = true;
    const resp = await request(`/api/sources/${source.id}/feeds`, "GET");
    loadingFeeds = false;
    if (resp.ok) {
      if (resp.content.feeds) {
        feeds = resp.content.feeds;
      }
    } else if (resp.error) {
      feedError = getErrorDetails(`Could not load feed`, resp);
    }
  };

  const deleteFeed = async (id: number) => {
    const resp = await request(`/api/sources/feeds/${id}`, "DELETE");
    if (resp.error) {
      feedError = getErrorDetails(`Could not load feed`, resp);
    }
    await fetchFeeds();
    calculateMissingFeeds();
  };

  const fetchFeedLogs = async (id: number) => {
    loadingLogs = true;
    const resp = await request(`/api/sources/feeds/${id}/log`, "GET");
    loadingLogs = false;
    if (resp.ok) {
      logs = resp.content;
    } else if (resp.error) {
      logError = getErrorDetails(`Could not load feed logs`, resp);
    }
  };

  const isSameFeed = (a: Feed, b: Feed) => a.url === b.url && a.rolie === b.rolie;

  const calculateMissingFeeds = () => {
    missingFeeds = pmdFeeds.filter((a) => !feeds.some((b) => isSameFeed(a, b)));
  };

  const getSource = async (id: number) => {
    await fetchSources();
    let found = sources.find((s) => s.id === id);
    if (found) {
      source = found;
      if (!source.headers) {
        source.headers = [];
      }
      parseHeaders();
      if (!source.ignore_patterns) {
        source.ignore_patterns = [""];
      }
      if (source.client_cert_private === "***") {
        source.client_cert_private = undefined;
      }
      if (source.client_cert_public === "***") {
        source.client_cert_public = undefined;
      }
      loadError = null;
    } else {
      loadError = getErrorDetails(`Could not find source`);
    }
  };

  const onChangedHeaders = () => {
    const lastIndex = headers.length - 1;
    if (
      (headers[lastIndex][0].length > 0 && headers[lastIndex][1].length > 0) ||
      (lastIndex - 1 >= 0 &&
        headers[lastIndex - 1][0].length > 0 &&
        headers[lastIndex - 1][1].length > 0)
    ) {
      headers.push(["", ""]);
      headers = headers;
    }
  };

  const parseHeaders = () => {
    headers = [];
    for (const header of source.headers) {
      let h = header.split(":");
      headers.push([h[0], h[1]]);
    }
    if (headers.length === 0) {
      headers.push(["", ""]);
    }
    onChangedHeaders();
  };

  onMount(async () => {
    let id = params?.id;
    if (id) {
      await getSource(Number(id));
      await fetchPMD();
      await fetchFeeds();
      calculateMissingFeeds();
    }
  });
</script>

<svelte:head>
  <title>Sources - {params?.id ? "Edit source" : "New source"}</title>
</svelte:head>

<Modal size="xs" title={modalTitle} bind:open={modalOpen} autoclose outsideclose>
  <div class="text-center">
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {modalMessage}
    </h3>
    <Button
      on:click={() => {
        modalCallback();
      }}
      color="red"
      class="me-2">Yes, I'm sure</Button
    >
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

{#if params?.id}
  <SectionHeader title={source.name}></SectionHeader>
  <div class="flex">
    <div class="flex-auto">
      <Table class="2xl:w-max" noborder>
        <TableBodyRow>
          <TableBodyCell>URL</TableBodyCell>
          <TableBodyCell>{source.url}</TableBodyCell>
        </TableBodyRow>
        {#if pmd}
          <TableBodyRow>
            <TableBodyCell>Canonical URL</TableBodyCell>
            <TableBodyCell>{pmd.canonical_url}</TableBodyCell>
          </TableBodyRow>
          <TableBodyRow>
            <TableBodyCell>Publisher Name</TableBodyCell>
            <TableBodyCell>{pmd.publisher.name}</TableBodyCell>
          </TableBodyRow>
          <TableBodyRow>
            <TableBodyCell>Publisher Contact</TableBodyCell>
            <TableBodyCell>{pmd.publisher.contact_details}</TableBodyCell>
          </TableBodyRow>
          <TableBodyRow>
            <TableBodyCell>Issuing Authority</TableBodyCell>
            <TableBodyCell>{pmd.publisher.issuing_authority}</TableBodyCell>
          </TableBodyRow>
        {/if}
      </Table>
      <div class:hidden={!loadingPMD} class:mb-4={true}>
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
    </div>

    <div class="flex-auto">
      <SourceForm
        {source}
        formSubmit={async () => {
          await saveSource();
        }}
        {formClass}
        enableActive={true}
      ></SourceForm>
    </div>
  </div>
  <CustomTable
    title="Feeds"
    headers={[
      {
        label: "Label",
        attribute: "label"
      },
      {
        label: "URL",
        attribute: "url"
      },
      {
        label: "Rolie",
        attribute: "rolie"
      },
      {
        label: "Log level",
        attribute: "log_level"
      }
    ]}
  >
    {#each feeds as feed, index (index)}
      <tr
        on:click={() => {
          if (feed.id) {
            fetchFeedLogs(feed.id);
          }
        }}
        class="cursor-pointer"
      >
        <TableBodyCell {tdClass}>{feed.label}</TableBodyCell>
        <TableBodyCell {tdClass}>{feed.url}</TableBodyCell>
        <TableBodyCell {tdClass}>{feed.rolie}</TableBodyCell>
        <TableBodyCell {tdClass}>{feed.log_level}</TableBodyCell>
        <td>
          <Button
            on:click={(event) => {
              event.stopPropagation();
              modalCallback = () => {
                if (feed.id) {
                  deleteFeed(feed.id);
                }
              };
              modalMessage = "Are you sure you want to delete this source?";
              modalTitle = `Source ${source.name}`;
              modalOpen = true;
            }}
            title={`Delete source "${source.name}"`}
            class="border-0 p-2"
            color="light"
          >
            <i class="bx bx-trash text-xl text-red-500"></i>
          </Button>
        </td>
      </tr>
    {/each}
    <div slot="bottom">
      <div class:hidden={!loadingFeeds} class:mb-4={true}>
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
      <ErrorMessage error={feedError}></ErrorMessage>
    </div>
  </CustomTable>

  <CustomTable
    title="Missing feeds"
    headers={[
      {
        label: "Label",
        attribute: "label"
      },
      {
        label: "URL",
        attribute: "url"
      },
      {
        label: "Rolie",
        attribute: "rolie"
      },
      {
        label: "Log level",
        attribute: "log_level"
      }
    ]}
  >
    {#each missingFeeds as feed, index (index)}
      <tr
        class="cursor-pointer"
        on:click={() => {
          newFeed = feed;
        }}
      >
        <TableBodyCell {tdClass}>{feed.label}</TableBodyCell>
        <TableBodyCell {tdClass}>{feed.url}</TableBodyCell>
        <TableBodyCell {tdClass}>{feed.rolie}</TableBodyCell>
        <TableBodyCell {tdClass}>{feed.log_level}</TableBodyCell>
        <td>
          <Button
            on:click={(event) => {
              event.stopPropagation();
              modalCallback = () => {
                console.log("TODO");
              };
              modalMessage = "Are you sure you want to delete this source?";
              modalTitle = `Source ${source.name}`;
              modalOpen = true;
            }}
            title={`Delete source "${source.name}"`}
            class="border-0 p-2"
            color="light"
          >
            <i class="bx bx-trash text-xl text-red-500"></i>
          </Button>
        </td>
      </tr>
    {/each}
    <div slot="bottom">
      <div class:hidden={!loadingFeeds && !loadingPMD} class:mb-4={true}>
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
      <ErrorMessage error={feedError}></ErrorMessage>
    </div>
  </CustomTable>
  {#if newFeed}
    <form
      on:submit={async () => {
        if (newFeed) {
          await saveFeeds([newFeed]);
          newFeed = null;
          await fetchFeeds();
          calculateMissingFeeds();
        }
      }}
      class={formClass}
    >
      <Label>URL</Label>
      <Input readonly bind:value={newFeed.url}></Input>
      <Label>Log level</Label>
      <Select items={logLevels} bind:value={newFeed.log_level} />
      <Label>Label</Label>
      <Input bind:value={newFeed.label}></Input>
      <br />
      <Button type="submit" color="light">
        <i class="bx bxs-save me-2"></i>
        <span>Save feed</span>
      </Button>
    </form>{/if}
  <CustomTable
    title="Logs"
    headers={[
      {
        label: "Time",
        attribute: "time"
      },
      {
        label: "level",
        attribute: "level"
      },
      {
        label: "Message",
        attribute: "msg"
      }
    ]}
  >
    {#each logs as log, index (index)}
      <tr>
        <TableBodyCell {tdClass}>{log.time}</TableBodyCell>
        <TableBodyCell {tdClass}>{log.level}</TableBodyCell>
        <TableBodyCell {tdClass}>{log.msg}</TableBodyCell>
      </tr>
    {/each}
    <div slot="bottom">
      <div class:hidden={!loadingLogs} class:mb-4={true}>
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
      <ErrorMessage error={logError}></ErrorMessage>
    </div>
  </CustomTable>
{:else}
  <SectionHeader title="Add new source"></SectionHeader>

  <StepIndicator size="h-1" class={formClass} currentStep={currentStep + 1} {steps} />
  {#if currentStep === 0}
    <form on:submit={fetchPMD} class={formClass}>
      <Label>URL</Label>
      <Input bind:value={source.url}></Input>
      <br />
      <div class:hidden={!loadingPMD} class:mb-4={true}>
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
      <Button type="submit" color="light">
        <i class="bx bx-check me-2"></i>
        <span>Check URL</span>
      </Button>
    </form>
  {/if}
  {#if currentStep === 1}
    <SourceForm
      {formClass}
      {source}
      formSubmit={async () => {
        if (await saveSource()) {
          currentStep = 2;
        }
      }}
    ></SourceForm>
  {/if}
  {#if currentStep === 2}
    <form
      on:submit={() => {
        saveFeeds(pmdFeeds);
      }}
      class={formClass}
    >
      {#each pmdFeeds as feed}
        <Label>URL</Label>
        <Input readonly bind:value={feed.url}></Input>
        <Label>Log level</Label>
        <Select items={logLevels} bind:value={feed.log_level} />
        <Label>Label</Label>
        <Input bind:value={feed.label}></Input>
        <Checkbox bind:checked={feed.enable}>Enable</Checkbox>
        <br />
      {/each}
      <Button type="submit" color="light">
        <i class="bx bxs-save me-2"></i>
        <span>Save feed</span>
      </Button>
    </form>
  {/if}
{/if}

<ErrorMessage error={saveError}></ErrorMessage>
<ErrorMessage error={loadError}></ErrorMessage>
