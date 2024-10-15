<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import ApexCharts from "apexcharts";
  import type { ErrorDetails } from "$lib/Errors/error";
  import { fetchStatistic, StatisticType, type Statistic } from "./source";
  import { onMount } from "svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";

  export let title: string;
  export let id: number | undefined = undefined;
  export let isFeed = false;
  let statisticError: ErrorDetails | null;
  let overviewChart: any;

  type GroupedStatistic = {
    import: Statistic;
    signatureFailed: Statistic;
    checksumFailed: Statistic;
    filenameFailed: Statistic;
    schemaFailed: Statistic;
    downloadFailed: Statistic;
    remoteFailed: Statistic;
    duplicateFailed: Statistic;
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
      grid: {
        borderColor: "#dedede",
        yaxis: {
          lines: {
            show: true
          }
        }
      },
      stroke: {
        width: [5, 7, 5],
        curve: "straight",
        dashArray: [0, 8, 5]
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
      noData: {
        text: "No information available."
      },
      xaxis: {
        categories: xAxis
      },
      theme: {
        mode: "light",
        palette: "palette3"
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

    const statistics: any = {};
    let type: any = new StatisticType();
    let resp = await fetchStatistic(from, to, "48h", type);
    if (!resp.ok) {
      statisticError = resp.error;
      return;
    } else {
      statistics.import = resp.value;
    }

    for (let i = 0; i < Object.getOwnPropertyNames(type).length; i++) {
      const propertyToFetch = Object.getOwnPropertyNames(type)[i];
      for (let j = 0; j < Object.getOwnPropertyNames(type).length; j++) {
        const property = Object.getOwnPropertyNames(type)[j];
        if (property === propertyToFetch) {
          type[property] = true;
        } else {
          type[property] = false;
        }
      }
      let response = await fetchStatistic(from, to, "48h", type, id, isFeed);
      if (!response.ok) {
        statisticError = response.error;
        return;
      } else {
        statistics[propertyToFetch] = response.value;
      }
    }
    stats = statistics;
  };

  onMount(() => {
    loadStats();
  });
</script>

<SectionHeader {title}></SectionHeader>
<div bind:this={overviewChart}></div>
<ErrorMessage error={statisticError}></ErrorMessage>
