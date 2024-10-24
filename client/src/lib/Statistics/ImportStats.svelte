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
    fetchImportFailuresStatistic,
    type StatisticGroup,
    type StatisticType,
    fetchBasicStatistic,
    getCVSSTextualRating,
    type CVSSTextualRating,
    mergeImportFailureStatistics
  } from "$lib/Statistics/statistics";
  import Chart from "chart.js/auto";
  import { Button, ButtonGroup, Input, Label, Spinner } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";
  import "chartjs-adapter-moment";
  import StatsTable from "./StatsTable.svelte";
  import {
    HOUR_MS,
    YEAR_MS,
    MONTH_MS,
    DAY_MS,
    WEEK_MS,
    pad,
    setToEndOfDay,
    toLocaleISOString
  } from "$lib/time";

  export let divContainerClass = "mb-16";
  export let height = "140pt";
  export let stepsInMinutes = 30;
  export let showLegend = false;
  export let showRangeSelection = false;
  export let initialFrom: Date = new Date(Date.now() - DAY_MS * 2);
  export let updateIntervalInMinutes: number | null = null;
  export let title = `Imports / ${stepsInMinutes} min`;
  export let axes: Axis[] = [{ label: "Docs", types: ["imports"] }];
  export let isStacked = false;
  export let showModeToggle = false;
  export let colors: string[] | undefined = undefined;
  export let source: Source | null = null;

  type Axis = {
    label: string;
    types: StatisticType[];
  };
  type StatisticsMode = "diagram" | "table";
  type Source = {
    id: number;
    isFeed: boolean;
  };

  let from: string | undefined;
  let to: string | undefined;
  let error: ErrorDetails | null = null;
  let chartComponentRef: any;
  let chart: any;
  let isLoading = false;
  let stats: StatisticGroup = {};
  let intervalID: ReturnType<typeof setInterval> | null;
  let stepsInMilliseconds = 1000 * 60 * stepsInMinutes;
  let mode: StatisticsMode = "diagram";
  const basicButtonClass = "py-1 px-3";
  const buttonClass = `${basicButtonClass} bg-white hover:bg-gray-100`;
  const pressedButtonClass = `${basicButtonClass} bg-gray-200 text-black hover:!bg-gray-100`;
  const updateInterval = 1000 * 60 * (updateIntervalInMinutes ?? 0);
  const categoryColors = [
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
  const rangeColors = ["#ddd", "#FFEFB0", "#E6A776", "#CD5D3A", "#B41500"];

  $: types = axes.map((axis) => axis.types).flat();
  $: datasets = Object.keys(stats).map((key: string, index: number) => {
    let label = key;
    if (key === "cve") label = "CVEs of documents";
    if (key === "imports") label = "Imports";
    if (key === "importFailuresCombined") label = "Failed imports";
    if (key === "signatureFailed") label = "Failed signature checks";
    if (key === "checksumFailed") label = "Failed checksum checks";
    if (key === "filenameFailed") label = "Failed filename checks";
    if (key === "schemaFailed") label = "Failed schema checks";
    if (key === "downloadFailed") label = "Failed downloads";
    if (key === "remoteFailed") label = "Failed remote";
    if (key === "duplicateFailed") label = "Failures because of duplicates";
    if (key.startsWith("cvss_")) {
      label = key.replace("cvss_", "");
    }
    const yAxisID = axes.findIndex((axis) => axis.types.includes(key as StatisticType));
    return {
      label: label,
      data: stats[key]?.map((s) => {
        return { x: s[0], y: s[1] };
      }),
      borderWidth: 0,
      backgroundColor: getColor(index),
      yAxisID: `y${yAxisID > 0 ? yAxisID : ""}`
    };
  });

  const getColor = (index: number) => {
    return colors
      ? colors[index]
      : types.length === 1 && types.includes("critical")
        ? rangeColors[index]
        : categoryColors[index];
  };

  const isToday = (date: Date) => {
    const today = new Date();
    return (
      date.getDate() === today.getDate() &&
      date.getMonth() === today.getMonth() &&
      date.getFullYear() === today.getFullYear()
    );
  };

  const loadStats = async () => {
    if (!from || !to) return;
    error = null;
    let response: any;
    const toParameter = isToday(new Date(to))
      ? new Date(Date.now() + HOUR_MS)
      : setToEndOfDay(new Date(to));
    const newStats: StatisticGroup = {};
    if (types.includes("imports")) {
      response = await fetchBasicStatistic(
        new Date(from),
        toParameter,
        stepsInMilliseconds,
        "imports",
        source?.id,
        source?.isFeed
      );
      if (response.ok) {
        Object.assign(newStats, response.value);
      } else {
        error = response.error;
      }
    }
    if (types.includes("importFailures") || types.includes("importFailuresCombined")) {
      response = await fetchImportFailuresStatistic(
        new Date(from),
        toParameter,
        stepsInMilliseconds,
        source?.id,
        source?.isFeed
      );
      if (response.ok) {
        if (types.includes("importFailuresCombined")) {
          Object.assign(newStats, mergeImportFailureStatistics(response.value));
        } else {
          Object.assign(newStats, response.value);
        }
      } else {
        error = response.error;
      }
    }
    if (types.includes("cve")) {
      response = await fetchBasicStatistic(
        new Date(from),
        toParameter,
        stepsInMilliseconds,
        "cve",
        source?.id,
        source?.isFeed
      );
      if (response.ok) {
        Object.assign(newStats, response.value);
      } else {
        error = response.error;
      }
    }
    if (types.includes("critical")) {
      response = await fetchBasicStatistic(
        new Date(from),
        toParameter,
        stepsInMilliseconds,
        "critical",
        source?.id,
        source?.isFeed
      );
      if (response.ok) {
        const crit: any = response.value.critical;
        if (crit) {
          const critStats: any = {
            cvss_null: [],
            cvss_None: [],
            cvss_Low: [],
            cvss_Medium: [],
            cvss_High: []
          };
          for (let i = 0; i < crit.length; i++) {
            const date = crit[i][0];
            const counts: any = {
              cvss_null: [],
              cvss_None: [],
              cvss_Low: [],
              cvss_Medium: [],
              cvss_High: []
            };
            const keys = Object.keys(critStats);
            // Iterate through the values of one point of time
            if (crit[i][1]) {
              for (let j = 0; j < crit[i][1].length; j++) {
                type NumberOfDocs = number;
                type CritCount = [number | null, NumberOfDocs];
                const critCount: CritCount = crit[i][1][j];
                const numberOfDocs = critCount[1];
                const cvss = critCount?.[0];
                if (cvss) {
                  const cvssTextialRating: CVSSTextualRating = getCVSSTextualRating(cvss);
                  counts[cvssTextialRating] = counts[`cvss_${cvssTextialRating}`] + numberOfDocs;
                } else {
                  counts["null"] = counts["null"] + numberOfDocs;
                }
              }
              keys.forEach((key) => {
                critStats[key].push([date, counts[key.replace("cvss_", "")]]);
              });
            } else {
              keys.forEach((key) => {
                critStats[key].push([date, 0]);
              });
            }
          }
          Object.assign(newStats, critStats);
        }
      } else {
        error = response.error;
      }
    }
    stats = newStats;
  };

  const setMode = (newMode: StatisticsMode) => {
    mode = newMode;
  };

  const updateOptions = () => {
    if (from) {
      const minFrom = new Date(from);
      minFrom.setHours(0);
      minFrom.setMinutes(0);
      minFrom.setSeconds(0);
      minFrom.setMilliseconds(0);
      chart.options.scales.x.min = minFrom;
    }
    if (to && from) {
      let maxTo = new Date(to);
      if (isToday(maxTo)) {
        maxTo = new Date(Date.now() + HOUR_MS);
      } else {
        maxTo = setToEndOfDay(new Date(to));
      }
      chart.options.scales.x.max = maxTo;
    }
  };

  const updateData = async () => {
    await loadStats();
    chart.data.datasets = datasets;
  };

  const updateChart = async () => {
    updateOptions();
    await updateData();
    chart.update();
  };

  // Source: https://stackoverflow.com/questions/6117814/get-week-of-year-in-javascript-like-in-php/6117889#6117889
  function getWeekNumber(d: Date) {
    // Copy date so don't modify original
    d = new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()));
    // Set to nearest Thursday: current date + 4 - current day number
    // Make Sunday's day number 7
    d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay() || 7));
    // Get first day of year
    var yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
    // Calculate full weeks to nearest Thursday
    var weekNo = Math.ceil(((d.getTime() - yearStart.getTime()) / 86400000 + 1) / 7);
    // Return array of year and week number
    return weekNo;
  }

  const createLabelForXAxis = (date: Date): string | undefined => {
    if (!from || !to) return;
    let label = "";
    const paddedMonth = pad(date.getMonth() + 1);
    const paddedDate = pad(date.getDate());
    const paddedHours = pad(date.getHours());
    const paddedMinutes = pad(date.getMinutes());
    let diff = new Date(to).getTime() - new Date(from).getTime();
    if (diff >= YEAR_MS) {
      label = `${date.getFullYear()}-${paddedMonth}`;
    } else if (diff > MONTH_MS + 3 * DAY_MS) {
      label = `${date.getFullYear()}-${paddedMonth}-${paddedDate}`;
    } else if (diff >= MONTH_MS) {
      label = `${date.getFullYear()}-W${getWeekNumber(date)}`;
    } else if (diff == WEEK_MS) {
      label = `${date.getFullYear()}-${paddedMonth}-${paddedDate}`;
    } else {
      label = `${paddedHours}:${paddedMinutes}`;
    }
    return label;
  };

  const initChart = () => {
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
          tooltip: {
            callbacks: {
              label: function (context: any) {
                if (context.formattedValue && context.dataset.label) {
                  return `${context.dataset.label}: ${context.formattedValue}`;
                }
                return "";
              },
              title: (tooltipItems: any[]) => {
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
            stacked: isStacked,
            ticks: {
              callback: (tickValue: string | number, _index: number, _ticks: any[]): string => {
                return createLabelForXAxis(new Date(tickValue)) ?? "";
              }
            },
            time: {
              // Overwrite to keep exact time.
              parser: (v: unknown): number => {
                if (v instanceof Date) return v.getTime();
                return 0;
              }
            },
            afterBuildTicks: (axis: any) => {
              const newTicks: any[] = [];
              const firstProperty = Object.keys(stats)[0];
              stats?.[firstProperty]?.forEach((stat) => {
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
            beginAtZero: true,
            stacked: isStacked,
            title: { display: axes[0].label.length > 0, text: axes[0].label }
          }
        }
      }
    });
    if (axes[1]) {
      const showLabel = axes[1].label.length > 0;
      chart.options.scales.y1 = {
        beginAtZero: true,
        grid: {
          drawOnChartArea: false // only want the grid lines for one axis to show up
        },
        title: { display: showLabel, text: axes[1].label },
        position: "right"
      };
    }
    // Remove "Crit" from legend labels because otherwise it would appear in front of every crit label
    // which would be too much "noise".
    chart.options.plugins.legend.labels.generateLabels = (chart: any) => {
      const items: any[] = [];
      chart.legend.legendItems.forEach((item: any, index: number) => {
        const datasetMeta = chart.getDatasetMeta(item.datasetIndex);
        const label = datasetMeta.label.replace("cvss_", "");
        items.push({
          text: label,
          datasetIndex: index,
          fillStyle: getColor(index),
          hidden: datasetMeta.hidden
        });
      });
      return items;
    };
  };

  onMount(async () => {
    from = initialFrom.toISOString().split("T")[0];
    to = new Date().toISOString().split("T")[0];
    await loadStats();
    initChart();
    if (updateIntervalInMinutes) {
      intervalID = setInterval(async () => {
        if (!isLoading) {
          updateChart();
        }
      }, updateInterval);
    }
    updateOptions();
    chart.update();
  });

  onDestroy(() => {
    if (intervalID) {
      clearInterval(intervalID);
    }
  });

  // Fit steps to selected time range so the bars don't become to thin.
  const updateSteps = () => {
    if (!from || !to) return;
    let diff = new Date(to).getTime() - new Date(from).getTime();
    if (diff >= YEAR_MS) {
      stepsInMilliseconds = Math.floor(diff / YEAR_MS) * MONTH_MS;
    } else if (diff >= MONTH_MS) {
      stepsInMilliseconds = WEEK_MS;
    } else if (diff >= WEEK_MS) {
      stepsInMilliseconds = DAY_MS;
    } else {
      stepsInMilliseconds = HOUR_MS;
    }
    updateChart();
  };

  // In case of month or year we need some padding so the last month/year is not cut-off.
  const selectPredefinedRange = (range: string) => {
    const newFrom = new Date();
    const newTo = new Date();
    let diff = 1;
    stepsInMilliseconds = HOUR_MS;
    if (range === "month") {
      stepsInMilliseconds = WEEK_MS;
      diff = 30;
      newTo.setDate(newTo.getDate() + 1);
    }
    if (range === "year") {
      stepsInMilliseconds = MONTH_MS;
      diff = 365;
      newTo.setMonth(newTo.getMonth() + 1);
    }
    newFrom.setDate(newFrom.getDate() - diff);
    from = newFrom.toISOString().split("T")[0];
    to = newTo.toISOString().split("T")[0];
    updateSteps();
  };
</script>

<div class={divContainerClass}>
  <div class="flex flex-col gap-4">
    <div class="flex gap-6">
      <h3>{title}</h3>
      {#if showModeToggle}
        <ButtonGroup>
          <Button
            class={mode === "diagram" ? pressedButtonClass : buttonClass}
            on:click={() => setMode("diagram")}><i class="bx bx-bar-chart"></i></Button
          >
          <Button
            class={mode === "table" ? pressedButtonClass : buttonClass}
            on:click={() => setMode("table")}><i class="bx bx-table"></i></Button
          >
        </ButtonGroup>
      {/if}
    </div>
    <ErrorMessage {error}></ErrorMessage>
    {#if isLoading}
      <div class:invisible={!isLoading} class={isLoading ? "loadingFadeIn" : ""}>
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
    {/if}
    <div hidden={mode === "table"} class="border px-2">
      <div style:height>
        <canvas bind:this={chartComponentRef}></canvas>
      </div>
    </div>
    {#if mode === "table"}
      <StatsTable {stats}></StatsTable>
    {/if}
    {#if showRangeSelection}
      <div class="my-2 flex flex-wrap items-end justify-start gap-4 md:justify-center">
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
