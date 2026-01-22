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
    baseScore?: string;
    baseSeverity?: string;
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

  const getClass = (severity: string | undefined, score: string | undefined) => {
    if (severity && score) {
      if (severity.toLowerCase() !== getCVSSTextualRating(Number(score)).toLowerCase()) {
        return "";
      } else {
        return severity.toLowerCase();
      }
    } else {
      if (severity) {
        return severity.toLowerCase();
      } else if (score !== undefined) {
        return getCVSSTextualRating(Number(score)).toLowerCase();
      }
    }
  };
</script>

{#if (baseScore !== null && baseScore !== undefined) || baseSeverity !== undefined}
  <div class={"score " + getClass(baseSeverity, baseScore)}>
    <span class="baseScore">{baseScore}</span>
    {#if (baseSeverity && baseScore === null) || baseScore === undefined}
      <span class="baseSeverity">{baseSeverity}</span>
    {:else if baseSeverity}
      <span class="baseSeverity">({baseSeverity})</span>
    {/if}
  </div>
{/if}

<style>
  .score.none,
  .score.high,
  .score.critical {
    color: #ffffff;
  }

  .score,
  .score.low,
  .score.medium {
    color: #222;
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
    background: #dddddd;
    font-size: small;
    font-weight: bolder;
    padding: 0.25em;
    height: 28px;
    min-width: fit-content;
    width: 36pt;
  }
</style>
