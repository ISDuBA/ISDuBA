<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  interface Props {
    baseScore: string;
    baseSeverity: string;
  }
  let { baseScore, baseSeverity }: Props = $props();

  const getSeverityClass = (severity: string) => {
    return severity.toLowerCase();
  };
</script>

<div class={"score " + getSeverityClass(baseSeverity)}>
  <span class="baseScore">{baseScore}</span>
  <span class="baseSeverity">({baseSeverity})</span>
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
    border: 2px solid #53aa33;
  }

  .score.low {
    background: #ffcb0d;
    border: 2px solid #ffcb0d;
  }

  .score.medium {
    background: #f9a009;
    border: 2px solid #f9a009;
  }

  .score.high {
    background: #df3d03;
    border: 2px solid #df3d03;
  }

  .score.critical {
    background: #cc0500;
    border: 2px solid #cc0500;
  }

  .score {
    display: flex;
    flex-wrap: nowrap;
    gap: 2pt;
    margin: 0 15px;
    border: 2px solid #666666;
    background: #dddddd;
    font-size: small;
    font-weight: bolder;
    padding: 0.25em;
    height: fit-content;
  }
</style>
