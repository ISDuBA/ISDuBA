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
  import ApexCharts from "apexcharts";
  import { Spinner } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";

  let error: ErrorDetails | null = null;
  let chartComponentRef: any;
  let chart: any;
  let isLoading = false;
  let stats: Statistic = [];
  let intervalID: ReturnType<typeof setInterval> | null;
  const steps = 1000 * 60 * 15; // 15 min
  const title = "Imports / 15 min";
  const updateInterval = 1000 * 60 * 15; // 10 min

  const loadStats = async () => {
    error = null;
    const from = new Date();
    from.setDate(from.getDate() - 2);
    const to = new Date();
    let response = await fetchStatistic(from, new Date(to), `${steps}ms`);
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
    for (let i = from.getTime(); i <= to.getTime(); i += steps) {
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
    const optionsToUpdate = {
      series: [
        {
          name: title,
          data: stats.map((s) => s[1])
        }
      ],
      xaxis: {
        categories: stats.map((s, index) => {
          if (index === 0) return s[0];
          else if (index === stats.length - 1) return "now";
          return "";
        })
      }
    };
    chart.updateOptions(optionsToUpdate, false, false);
  };

  onMount(async () => {
    await loadStats();
    let options = {
      colors: ["#aaddac"],
      chart: {
        height: 150,
        type: "bar",
        zoom: {
          enabled: false,
          autoScaleYaxis: true
        },
        toolbar: {
          show: false
        }
      },
      dataLabels: {
        enabled: false
      },
      grid: {
        borderColor: "#dedede",
        yaxis: {
          lines: {
            show: true
          }
        }
      },
      noData: {
        text: "No information available."
      },
      plotOptions: {
        bar: {
          columnWidth: 4,
          borderRadius: 0
        }
      },
      series: [
        {
          name: title,
          data: stats.map((s) => s[1])
        }
      ],
      title: {
        text: title
      },
      tooltip: {
        x: {
          formatter: (_val: any, details: any) => {
            const { dataPointIndex } = details;
            const start = toLocaleISOString(stats[dataPointIndex][0]);
            const end = stats[dataPointIndex + 1]?.[0]
              ? toLocaleISOString(stats[dataPointIndex + 1]?.[0])
              : undefined;
            return `${start}${end ? " - " + end : ""}`;
          }
        },
        y: {
          formatter: (val: any) => {
            return val;
          }
        }
      },
      xaxis: {
        axisTicks: {
          show: false
        },
        categories: stats.map((s, index) => {
          if (index === 0) return s[0];
          else if (index === stats.length - 1) return "now";
          return "";
        }),
        labels: {
          rotate: 0,
          formatter: (val: Date | string | number) => {
            if (typeof val === "number") return "";
            if (typeof val === "string") return val;
            return toLocaleISOString(val);
          }
        },
        stepSize: 10
      }
    };
    chart = new ApexCharts(chartComponentRef, options);
    chart.render();
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
  <div class="border">
    <div bind:this={chartComponentRef}></div>
  </div>
  {#if isLoading}
    <div class:invisible={!isLoading} class={isLoading ? "loadingFadeIn" : ""}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
  {/if}
</div>
