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
    fetchStatistic,
    toLocaleISOString,
    type Statistic,
    type StatisticEntry
  } from "$lib/Statistics/statistics";
  import Chart from "chart.js/auto";
  import { Spinner } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";
  import "chartjs-adapter-moment";

  let error: ErrorDetails | null = null;
  let chartComponentRef: any;
  let chart: any;
  let isLoading = false;
  let stats: Statistic = [];
  let intervalID: ReturnType<typeof setInterval> | null;
  const stepsInMinutes = 30;
  const stepsInMilliseconds = 1000 * 60 * stepsInMinutes;
  const title = `Imports / ${stepsInMinutes} min`;
  const updateInterval = 1000 * 60 * 10; // 10 min

  const loadStats = async () => {
    error = null;
    const from = new Date();
    from.setDate(from.getDate() - 2);
    const to = new Date();
    let response = await fetchStatistic(from, new Date(to), `${stepsInMilliseconds}ms`);
    if (response.ok) {
      const mappedValues: Statistic = response.value.map((entry: any) => {
        const date = new Date(entry[0]);
        date.setMinutes(date.getMinutes() + date.getTimezoneOffset());
        return [date, entry[1]];
      });
      stats = fillGaps(from, to, mappedValues);
      if (stats.length === 0 && chart) {
        chart.updateOptions({
          noData: {
            text: "No documents imported recently."
          }
        });
      }
    } else {
      error = response.error;
    }
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
    chart.data.datasets[0].data = stats.map((s) => {
      return { x: s[0], y: s[1] };
    });
    chart.update();
  };

  onMount(async () => {
    await loadStats();
    chart = new Chart(chartComponentRef, {
      type: "bar",
      data: {
        datasets: [
          {
            label: title,
            data: stats.map((s) => {
              return { x: s[0], y: s[1] };
            }),
            borderWidth: 1,
            backgroundColor: "#aaddac"
          }
        ]
      },
      options: {
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false
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
                if (index === 0 || index === stats.length - 1)
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
              stats.forEach((stat) => {
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
    intervalID = setInterval(async () => {
      if (!isLoading) {
        await loadStats();
        updateOptions();
      }
    }, updateInterval);
  });

  onDestroy(() => {
    if (intervalID) {
      clearInterval(intervalID);
    }
  });
</script>

<div class="mb-8 flex w-full max-w-[96%] flex-col 2xl:w-[46%]">
  <ErrorMessage {error}></ErrorMessage>
  <div class="max-h-72 border px-2">
    <canvas bind:this={chartComponentRef}></canvas>
  </div>
  {#if isLoading}
    <div class:invisible={!isLoading} class={isLoading ? "loadingFadeIn" : ""}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
  {/if}
</div>
