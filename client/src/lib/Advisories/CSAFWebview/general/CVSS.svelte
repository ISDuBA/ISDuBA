<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import type { CVSSTextualRating } from "$lib/Statistics/statistics";

  interface Props {
    baseScore: string;
    baseSeverity: string;
  }
  let { baseScore, baseSeverity }: Props = $props();

  const getCVSSTextualRating = (CVSS: number): CVSSTextualRating => {
    if (CVSS === 0) {
      return "None";
    } else if (CVSS <= 3.9) {
      return "Low";
    } else if (CVSS <= 6.9) {
      return "Medium";
    } else if (CVSS <= 8.9) {
      return "High";
    } else {
      return "Critical";
    }
  };

  const getSeverityClass = (severity: string, score: string) => {
    if (severity) {
      return severity.toLowerCase();
    } else if (score) {
      return getCVSSTextualRating(Number(score)).toLowerCase();
    }
  };
</script>

<div class={"score " + getSeverityClass(baseSeverity, baseScore)}>
  <span class="baseScore">{baseScore}</span>
  {#if baseSeverity}
    <span class="baseSeverity">({baseSeverity})</span>
  {/if}
</div>

<style>
  .score.none,
  .score.low,
  .score.medium,
  .score.high,
  .score.critical {
    color: #ffffff;
  }

  .score.none {
    background: #53aa33;
  }

  .score.low {
    background: #ffcb0d;
  }

  .score.medium {
    background: #f9a009;
  }

  .score.high {
    background: #df3d03;
  }

  .score.critical {
    background: #cc0500;
  }

  .score {
    display: flex;
    flex-wrap: nowrap;
    gap: 2pt;
    justify-content: center;
    margin: 0 15px;
    background: #dddddd;
    color: black;
    font-size: small;
    font-weight: bolder;
    padding: 0.25em;
    height: fit-content;
    min-width: fit-content;
    width: 50pt;
  }
</style>
