<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import ApexCharts from "apexcharts";

  var options = {
    chart: {
      type: "line"
    },
    series: [
      {
        name: "sales",
        data: [30, 40, 35, 50, 49, 60, 70, 91, 125]
      }
    ],
    xaxis: {
      categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999]
    }
  };

  let chartDiv: any;

  $: if (chartDiv) {
    let chart = new ApexCharts(chartDiv, options);

    chart.render();
  }
</script>

<svelte:head>
  <title>Statistics</title>
</svelte:head>

<SectionHeader title="Statistics"></SectionHeader>
<hr class="mb-6" />

<div bind:this={chartDiv}></div>
