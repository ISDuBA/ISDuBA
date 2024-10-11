<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Spinner, TableBodyCell } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { push } from "svelte-spa-router";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails, getErrorDetails } from "$lib/Errors/error";
  import { tdClass } from "$lib/Table/defaults";
  import { request } from "$lib/request";
  import { onDestroy, onMount } from "svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import {
    type Source,
    type Statistic,
    StatisticType,
    fetchSources,
    fetchStatistic
  } from "$lib/Sources/source";
  import { appStore } from "$lib/store";
  import ApexCharts from "apexcharts";

  let messageError: ErrorDetails | null;
  let sourcesError: ErrorDetails | null;
  let statisticError: ErrorDetails | null;

  let loadingSources: boolean = false;

  let sources: Source[] = [];
  async function getMessage() {
    const response = await request("api/sources/message", "GET");
    if (response.ok) {
      return response.content;
    } else {
      messageError = getErrorDetails(`Couldn't load default message`, response);
    }
    return new Map<string, [string]>();
  }

  let sourceUpdate = setInterval(async () => {
    if (appStore.isEditor() || appStore.isSourceManager()) {
      getSources();
    }
  }, 30 * 1000);

  const getSources = async () => {
    loadingSources = true;
    const result = await fetchSources(true);
    loadingSources = false;
    if (result.ok) {
      sources = result.value;
    } else {
      sources = [];
      sourcesError = result.error;
    }
  };

  let overviewChart: any;

  type GroupedStatistic = {
    import: Statistic;
    signatureFailed: Statistic;
    checksumFailed: Statistic;
  };
  let stats: GroupedStatistic | undefined = undefined;

  const fixAxis = (stats: Statistic, axis: Date[]) => {
    let newStatistics: Statistic = [] as unknown as Statistic;
    for (let date of axis) {
      let entry = stats.find((i) => i[0] === date) ?? [date, 0];
      newStatistics.push(entry);
    }
    return newStatistics;
  };

  $: if (overviewChart && stats) {
    let importAxis = stats.import.map((a) => a[0]);
    let signatureAxis = stats.signatureFailed.map((a) => a[0]);
    let checksumAxis = stats.checksumFailed.map((a) => a[0]);

    let xAxis = [...new Set([...importAxis, ...signatureAxis, ...checksumAxis])];
    let importData = fixAxis(stats.import, xAxis).map((a) => a[1]);
    let signatureData = fixAxis(stats.signatureFailed, xAxis).map((a) => a[1]);
    let checksumData = fixAxis(stats.checksumFailed, xAxis).map((a) => a[1]);
    let options = {
      series: [
        {
          name: "Imported documents",
          data: importData
        },
        {
          name: "Signature failed",
          data: signatureData
        },
        {
          name: "Checksum failed",
          data: checksumData
        }
      ],
      chart: {
        height: 350,
        type: "line",
        zoom: {
          enabled: false
        }
      },
      dataLabels: {
        enabled: false
      },
      stroke: {
        width: [5, 7, 5],
        curve: "straight",
        dashArray: [0, 8, 5]
      },
      title: {
        text: "Download Statistics",
        align: "left"
      },
      legend: {
        tooltipHoverFormatter: function (val: any, opts: any) {
          return (
            val +
            " - <strong>" +
            opts.w.globals.series[opts.seriesIndex][opts.dataPointIndex] +
            "</strong>"
          );
        }
      },
      markers: {
        size: 0,
        hover: {
          sizeOffset: 6
        }
      },
      xaxis: {
        categories: xAxis
      },
      tooltip: {
        y: [
          {
            title: {
              formatter: function (val: any) {
                return val;
              }
            }
          },
          {
            title: {
              formatter: function (val: any) {
                return val;
              }
            }
          },
          {
            title: {
              formatter: function (val: any) {
                return val;
              }
            }
          }
        ]
      },
      grid: {
        borderColor: "#f1f1f1"
      }
    };
    let chart = new ApexCharts(overviewChart, options);
    chart.render();
  }

  const loadStats = async () => {
    let from = new Date();
    from.setDate(from.getDate() - 30);
    from.setMinutes(0, 0, 0);
    let to = new Date();
    to.setMinutes(0, 0, 0);

    let importStats: Statistic;
    let checksumFailed: Statistic;
    let signatureFailed: Statistic;

    let type = new StatisticType();
    let resp = await fetchStatistic(from, to, "48h", type);
    if (!resp.ok) {
      statisticError = resp.error;
      return;
    } else {
      importStats = resp.value;
    }

    type.signatureFailed = true;
    resp = await fetchStatistic(from, to, "48h", type);
    if (!resp.ok) {
      statisticError = resp.error;
      return;
    } else {
      signatureFailed = resp.value;
    }

    type.signatureFailed = false;
    type.checksumFailed = true;

    resp = await fetchStatistic(from, to, "48h", type);
    if (!resp.ok) {
      statisticError = resp.error;
      return;
    } else {
      checksumFailed = resp.value;
    }

    stats = {
      import: importStats,
      checksumFailed: checksumFailed,
      signatureFailed: signatureFailed
    };
  };

  onMount(async () => {
    if (appStore.isEditor() || appStore.isSourceManager()) {
      await getSources();
    }
    await loadStats();
  });

  onDestroy(() => {
    clearInterval(sourceUpdate);
  });
</script>

<svelte:head>
  <title>Sources</title>
</svelte:head>

<div>
  <SectionHeader title="Sources"></SectionHeader>
  {#if appStore.isEditor() || appStore.isSourceManager()}
    <CustomTable
      title="CSAF Provider"
      headers={[
        {
          label: "Name",
          attribute: "name"
        },
        {
          label: "Domain/PMD",
          attribute: "url"
        },
        {
          label: "Active",
          attribute: "active"
        },
        {
          label: "Loading/Queued",
          attribute: "stats"
        }
      ]}
    >
      {#each sources as source, index (index)}
        <tr
          on:click={() => {
            if (appStore.isSourceManager()) {
              push(`/sources/${source.id}`);
            }
          }}
          on:blur={() => {}}
          on:focus={() => {}}
          class={appStore.isSourceManager() ? "cursor-pointer" : ""}
        >
          <TableBodyCell {tdClass}>{source.name}</TableBodyCell>
          <TableBodyCell {tdClass}>{source.url}</TableBodyCell>
          <TableBodyCell {tdClass}
            ><i class={"bx " + (source.active ? "bxs-circle" : "bx-circle")}></i></TableBodyCell
          >
          <TableBodyCell {tdClass}
            >{source.stats?.downloading}/{source.stats?.waiting}</TableBodyCell
          >
        </tr>
      {/each}
      <div slot="bottom">
        <div
          class:invisible={!loadingSources}
          class={loadingSources ? "loadingFadeIn" : ""}
          class:mb-4={true}
        >
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
        {#if appStore.isSourceManager()}
          <Button href="/#/sources/new" class="mb-2" color="primary" size="xs">
            <i class="bx bx-plus"></i>
            <span>Add source</span>
          </Button>
        {/if}
        <ErrorMessage error={sourcesError}></ErrorMessage>
      </div>
    </CustomTable>
  {/if}
  {#await getMessage() then resp}
    {#if resp.message}
      {resp.message}
    {/if}
  {/await}
  <ErrorMessage error={sourcesError}></ErrorMessage>
  <ErrorMessage error={messageError}></ErrorMessage>

  <br />
  {#if appStore.isImporter()}
    <Button href="/#/sources/upload" class="my-2" color="primary" size="xs">
      <i class="bx bx-upload"></i>
      <span>Upload documents</span>
    </Button>
  {/if}
</div>

<SectionHeader title="Statistics"></SectionHeader>
<div bind:this={overviewChart}></div>
<ErrorMessage error={statisticError}></ErrorMessage>
