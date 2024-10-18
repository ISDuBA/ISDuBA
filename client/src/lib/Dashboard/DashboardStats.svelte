<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import {
    fetchAllImportStatistic,
    fetchStatistic,
    toLocaleISOString,
    type ImportStatistic,
    type Statistic,
    type StatisticEntry
  } from "$lib/Statistics/statistics";
  import Chart from "chart.js/auto";
  import { Button, ButtonGroup, Input, Label, Spinner } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";
  import "chartjs-adapter-moment";

  export let height = "140pt";
  export let stepsInMinutes = 30;
  export let showAllImportStats = false;
  export let showLegend = false;
  export let enableRangeSelection = false;
  export let initialFrom: Date = new Date(Date.now() - 1000 * 60 * 60 * 24 * 2);
  export let updateIntervalInMinutes: number | null = null;
  export let title = `Imports / ${stepsInMinutes} min`;

  let from: string | undefined;
  let to: string | undefined;
  let error: ErrorDetails | null = null;
  let chartComponentRef: any;
  let chart: any;
  let isLoading = false;
  let stats: ImportStatistic = {};
  let intervalID: ReturnType<typeof setInterval> | null;
  let stepsInMilliseconds = 1000 * 60 * stepsInMinutes;
  const updateInterval = 1000 * 60 * (updateIntervalInMinutes ?? 0);
  const colors = [
    "#3D6090",
    "#7EC9EE",
    "#4FB4A3",
    "#2A904D",
    "#AFA11C",
    "#E2D072",
    "#DB667A",
    "#8B084A",
    "#C552B1"
  ];

  $: datasets = Object.keys(stats).map((key: string, index: number) => {
    let label = "";
    if (key === "imports") label = "Imported documents";
    if (key === "signatureFailed") label = "Failed signature checks";
    if (key === "checksumFailed") label = "Failed checksum checks";
    if (key === "filenameFailed") label = "Failed filename checks";
    if (key === "schemaFailed") label = "Failed schema checks";
    if (key === "downloadFailed") label = "Failed downloads";
    if (key === "remoteFailed") label = "Failed remote";
    if (key === "duplicateFailed") label = "Duplicates";
    return {
      label: label,
      data: stats[key]?.map((s) => {
        return { x: s[0], y: s[1] };
      }),
      borderWidth: 1,
      backgroundColor: colors[index]
    };
  });

  const loadStats = async () => {
    if (!from || !to) return;
    error = null;
    let response;
    const newStats: ImportStatistic = {};
    if (showAllImportStats) {
      response = await fetchAllImportStatistic(
        new Date(from),
        new Date(to),
        `${stepsInMilliseconds}ms`
      );
      if (response.ok) {
        const keys = Object.keys(response.value);
        for (let i = 0; i < keys.length; i++) {
          const key = keys[i];
          if (response.value[key]) {
            newStats[key] = fillGaps(new Date(from), new Date(to), response.value[key]);
          }
        }
        if (chart && (!stats.imports || stats.imports?.length === 0)) {
          chart.updateOptions({
            noData: {
              text: "No documents imported recently."
            }
          });
        }
      } else {
        error = response.error;
      }
    } else {
      response = await fetchStatistic(new Date(from), new Date(to), `${stepsInMilliseconds}ms`);
      if (response.ok) {
        newStats.imports = fillGaps(new Date(from), new Date(to), response.value);
        if (newStats.imports.length === 0 && chart) {
          chart.updateOptions({
            noData: {
              text: "No documents imported recently."
            }
          });
        }
      } else {
        error = response.error;
      }
    }
    stats = newStats;
  };

  // Fill gaps with null values so the user can see at which times nothing was imported.
  const fillGaps = (from: Date, to: Date, values: Statistic) => {
    const newStats: Statistic = [];
    for (let i = from.getTime(); i <= to.getTime(); i += stepsInMilliseconds) {
      const foundValue: StatisticEntry | undefined = values.find(
        (v: StatisticEntry) => v[0].getTime() === i
      );
      if (foundValue) {
        newStats.push(foundValue);
      } else {
        newStats.push([new Date(i), null]);
      }
    }
    return newStats;
  };

  const updateOptions = () => {
    chart.data.datasets = datasets;
    chart.update();
  };

  const updateChart = async () => {
    await loadStats();
    updateOptions();
  };

  onMount(async () => {
    from = initialFrom.toISOString().split("T")[0];
    to = new Date().toISOString().split("T")[0];
    await loadStats();
    chart = new Chart(chartComponentRef, {
      type: "bar",
      data: {
        datasets
      },
      options: {
        maintainAspectRatio: false,
        aspectRatio: 1,
        plugins: {
          legend: {
            display: showLegend
          },
          title: {
            align: "start",
            display: true,
            padding: {
              bottom: 20
            },
            text: title
          },
          tooltip: {
            callbacks: {
              label: function (context) {
                if (context.formattedValue && context.dataset.label) {
                  return `${context.formattedValue} ${context.dataset.label}`;
                }
                return "";
              },
              title: (tooltipItems) => {
                const start: any = tooltipItems[0].dataset.data[tooltipItems[0].dataIndex];
                const end: any = tooltipItems[0].dataset.data[tooltipItems[0].dataIndex + 1];
                return `${toLocaleISOString(start.x)}${end ? " - " : ""}${end ? toLocaleISOString(end.x) : ""}`;
              }
            }
          }
        },
        scales: {
          x: {
            type: "time",
            grid: {
              display: true,
              drawOnChartArea: false,
              drawTicks: true,
              tickLength: 6,
              tickWidth: 2
            },
            ticks: {
              callback: (tickValue, index): string => {
                if (index === 0 || (stats.imports?.length && index === stats.imports?.length - 1))
                  return toLocaleISOString(new Date(tickValue));
                else return "";
              }
            },
            time: {
              parser: (v: any): number => {
                return v.getTime();
              }
            },
            afterBuildTicks: (axis) => {
              const newTicks: any[] = [];
              stats.imports?.forEach((stat) => {
                newTicks.push({
                  value: stat[0].getTime(),
                  major: false,
                  label: toLocaleISOString(stat[0])
                });
              });
              axis.ticks = newTicks;
            }
          },
          y: {
            beginAtZero: true
          }
        }
      }
    });
    if (updateIntervalInMinutes) {
      intervalID = setInterval(async () => {
        if (!isLoading) {
          await loadStats();
          updateOptions();
        }
      }, updateInterval);
    }
  });

  onDestroy(() => {
    if (intervalID) {
      clearInterval(intervalID);
    }
  });

  const resetTo = () => {
    const newTo = new Date();
    to = newTo.toISOString().split("T")[0];
  };

  const updateSteps = () => {
    if (!from || !to) return;
    let diff = new Date(to).getTime() - new Date(from).getTime();
    const minute = 1000 * 60;
    const hour = 1000 * 60 * 60;
    const day = hour * 24;
    const week = day * 7;
    const month = week * 4;
    const year = month * 12;
    if (diff >= year) {
      stepsInMilliseconds = Math.floor(diff / year) * month;
    } else if (diff >= month) {
      stepsInMilliseconds = week;
    } else if (diff >= week) {
      stepsInMilliseconds = day;
    } else if (diff >= day) {
      stepsInMilliseconds = hour;
    } else {
      stepsInMilliseconds = minute;
    }
    updateChart();
  };

  const selectPredefinedRange = (range: string) => {
    resetTo();
    const newFrom = new Date();
    let diff = 1;
    stepsInMilliseconds = 1000 * 60 * 60;
    if (range === "month") {
      stepsInMilliseconds = 1000 * 60 * 60 * 24 * 7;
      diff = 30;
    }
    if (range === "year") {
      stepsInMilliseconds = 1000 * 60 * 60 * 24 * 7 * 4;
      diff = 365;
    }
    newFrom.setDate(newFrom.getDate() - diff);
    from = newFrom.toISOString().split("T")[0];
    updateSteps();
  };
</script>

<div class="mb-8 flex w-full max-w-[96%] flex-col 2xl:w-[46%]">
  <ErrorMessage {error}></ErrorMessage>
  {#if isLoading}
    <div class:invisible={!isLoading} class={isLoading ? "loadingFadeIn" : ""}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
  {/if}
  <div class="flex flex-col gap-4 border px-2">
    <div style:height>
      <canvas bind:this={chartComponentRef}></canvas>
    </div>
    {#if enableRangeSelection}
      <div class="my-2 flex items-end justify-center gap-4">
        <Label for="from"
          ><span>From:</span>
          <Input let:props>
            <input on:change={updateSteps} id="from" type="date" {...props} bind:value={from} />
          </Input>
        </Label>
        <Label for="to"
          ><span>To:</span>
          <Input let:props>
            <input on:change={updateSteps} id="to" type="date" {...props} bind:value={to} />
          </Input>
        </Label>
        <ButtonGroup class="h-fit">
          <Button
            on:click={() => {
              selectPredefinedRange("day");
            }}>Day</Button
          >
          <Button
            on:click={() => {
              selectPredefinedRange("month");
            }}>Month</Button
          >
          <Button
            on:click={() => {
              selectPredefinedRange("year");
            }}>Year</Button
          >
        </ButtonGroup>
      </div>
    {/if}
  </div>
</div>
