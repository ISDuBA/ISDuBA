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

  let firstChart: any;

  $: if (firstChart) {
    let firstOptions = {
      series: [
        {
          name: "Low crit",
          data: [44, 55, 41, 67, 22, 43]
        },
        {
          name: "Medium crit",
          data: [11, 17, 15, 15, 21, 14]
        },
        {
          name: "High crit",
          data: [21, 7, 25, 13, 22, 8]
        }
      ],
      chart: {
        type: "bar",
        height: 350,
        stacked: true,
        toolbar: {
          show: true
        },
        zoom: {
          enabled: true
        }
      },
      responsive: [
        {
          breakpoint: 480,
          options: {
            legend: {
              position: "bottom",
              offsetX: -10,
              offsetY: 0
            }
          }
        }
      ],
      plotOptions: {
        bar: {
          horizontal: false,
          borderRadiusApplication: "end", // 'around', 'end'
          borderRadiusWhenStacked: "last", // 'all', 'last'
          dataLabels: {
            total: {
              enabled: true,
              style: {
                fontSize: "13px",
                fontWeight: 900
              }
            }
          }
        }
      },
      xaxis: {
        type: "datetime",
        categories: [
          "01/01/2011 GMT",
          "01/02/2011 GMT",
          "01/03/2011 GMT",
          "01/04/2011 GMT",
          "01/05/2011 GMT",
          "01/06/2011 GMT"
        ]
      },
      legend: {
        position: "right",
        offsetY: 40
      },
      fill: {
        opacity: 1
      },
      title: {
        text: "New documents",
        align: "left"
      }
    };

    let chart = new ApexCharts(firstChart, firstOptions);

    chart.render();
  }

  let secondChart: any;

  $: if (secondChart) {
    let secondOptions = {
      series: [
        {
          name: "Imported documents",
          data: [45, 52, 38, 24, 33, 26, 21, 20, 6, 8, 15, 10]
        },
        {
          name: "Signature failed",
          data: [35, 41, 62, 42, 13, 18, 29, 37, 36, 51, 32, 35]
        },
        {
          name: "Download failed",
          data: [87, 57, 74, 99, 75, 38, 62, 47, 82, 56, 45, 47]
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
        categories: [
          "01 Jan",
          "02 Jan",
          "03 Jan",
          "04 Jan",
          "05 Jan",
          "06 Jan",
          "07 Jan",
          "08 Jan",
          "09 Jan",
          "10 Jan",
          "11 Jan",
          "12 Jan"
        ]
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
    let chart = new ApexCharts(secondChart, secondOptions);

    chart.render();
  }

  let thirdChart: any;

  $: if (thirdChart) {
    let thirdOptions = {
      chart: {
        type: "line"
      },
      series: [
        {
          name: "CVE",
          data: [30, 40, 48, 60, 72, 79, 87, 91, 100]
        }
      ],
      xaxis: {
        categories: [2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024]
      },
      title: {
        text: "Total CVE count",
        align: "left"
      }
    };

    let chart = new ApexCharts(thirdChart, thirdOptions);

    chart.render();
  }

  let fourthChart: any;

  $: if (fourthChart) {
    let fourthOptions = {
      chart: {
        type: "line"
      },
      series: [
        {
          name: "Advisories",
          data: [30, 40, 35, 50, 49, 60, 70, 91, 125]
        }
      ],
      xaxis: {
        categories: [2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024]
      },
      title: {
        text: "Total Advisory count",
        align: "left"
      }
    };

    let chart = new ApexCharts(fourthChart, fourthOptions);

    chart.render();
  }
</script>

<svelte:head>
  <title>Statistics</title>
</svelte:head>

<SectionHeader title="Statistics"></SectionHeader>
<hr class="mb-6" />
<div class="grid grid-cols-2 gap-4">
  <div bind:this={firstChart}></div>
  <div bind:this={secondChart}></div>
  <div bind:this={thirdChart}></div>
  <div bind:this={fourthChart}></div>
</div>
